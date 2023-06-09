package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/zxzixuanwang/image-syncer/config"
	syncjob "github.com/zxzixuanwang/image-syncer/pkg/sync-job"
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
		num, order, err := handleImageName(imageName, imageTag, l, handleDefaultTag)
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
func handleImageName(name, tagName string, l *logrus.Logger, fun f) ([]string, []int, error) {
	var err error
	defer func() {
		if err != nil {
			l.Errorf(err.Error()+" image name is %s:%s", name, tagName)
		}
	}()

	formatSlice := strings.SplitN(name, "/", 3)
	if len(formatSlice) > 3 || len(formatSlice) < 2 {
		err = fmt.Errorf("error image name format")
		return nil, nil, err
	}

	if err = fun(tagName); err != nil {
		return nil, nil, err
	}
	sourceName := formatSlice[0]
	namespace := formatSlice[1]
	order := make([]int, 0, config.Conf.Sync.MaxSyncDes)
	var des = make([]string, 0, config.Conf.Sync.MaxSyncDes)
	i := 0
	for _, v := range config.Conf.Registry.Reg {

		if v.SourceName == sourceName && v.Namespace == namespace {
			des = append(des, v.DestinationName)
			order = append(order, i)
		}
		i++
	}

	return des, order, nil
}

type f func(tag string) error

func handleDefaultTag(tag string) error {
	tagLow := strings.ToLower(tag)
	if strings.HasPrefix(tagLow, "snapshot") || strings.HasSuffix(tagLow, "snapshot") {
		return fmt.Errorf("has invalid tag")
	}
	return nil
}

func authHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, p, ok := r.BasicAuth()
		if !ok || !(u == config.Conf.Auth.Username && p == config.Conf.Auth.Password) {
			w.WriteHeader(http.StatusUnauthorized)
			response(w, []byte(StatusUnauthorized))
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
