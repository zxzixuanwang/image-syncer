package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/zxzixuanwang/image-syncer/config"
	syncjob "github.com/zxzixuanwang/image-syncer/pkg/sync-job"
	"github.com/zxzixuanwang/image-syncer/tools"
)

func Route(l *logrus.Logger) *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/images/sync/hook", func(w http.ResponseWriter, r *http.Request) {

		imageName := r.URL.Query().Get("name")
		if imageName == "" {
			w.WriteHeader(http.StatusOK)
			return
		}
		imageTag := r.URL.Query().Get("tag")
		if imageTag == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		imageWholeName := fmt.Sprintf("%s:%s", imageName, imageTag)

		l.Debugf("Gotted image is %s", imageWholeName)
		num, order, err := tools.HandleImageName(imageName, imageTag, l, tools.HandleDefaultTag)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}
		for i := 0; i < len(num); i++ {
			syncjob.Set(fmt.Sprintf("%s&%s", imageWholeName, num[i]), order[i])
		}
		l.Debug("get images", syncjob.Get())

		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)
	r.HandleFunc("/images/sync/do", func(w http.ResponseWriter, r *http.Request) {

	}).Methods(http.MethodGet)
	r.Use(authHandler)
	r.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			l.Infof("host:%s,uri:%s,remoteip:%s", r.Host, r.RequestURI, r.RemoteAddr)
			h.ServeHTTP(w, r)
		})
	})
	return r
}

func authHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok || !(u == config.Conf.Auth.Username && p == config.Conf.Auth.Password) {
			w.WriteHeader(http.StatusUnauthorized)
			_ = response(w, []byte(StatusUnauthorized))
			return
		}

		h.ServeHTTP(w, r)
	})
}

const (
	StatusUnauthorized = "no auth or wrong pasword or wrong username"
)

func response(w http.ResponseWriter, content []byte) error {
	_, err := w.Write(content)
	return err
}
