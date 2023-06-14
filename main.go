package main

import (
	"net/http"

	"github.com/zxzixuanwang/image-syncer/config"
	"github.com/zxzixuanwang/image-syncer/pkg/log"
	syncjob "github.com/zxzixuanwang/image-syncer/pkg/sync-job"
	"github.com/zxzixuanwang/image-syncer/web"
)

func main() {
	l := log.NewFileLogger(config.Conf.Sync.FilePosition.Log, config.Conf.App.Env)

	if err := syncjob.Tickerjob(l, config.Conf.Sync.Interval); err != nil {
		panic(err)
	}
	l.Infof("listening port %s", config.Conf.App.Port)
	if err := http.ListenAndServe(config.Conf.App.Port, web.Route(l)); err != nil {
		panic(err)
	}
}
