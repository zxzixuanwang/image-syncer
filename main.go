package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/zxzixuanwang/image-syncer/config"
	mq "github.com/zxzixuanwang/image-syncer/pkg/Mq"
	"github.com/zxzixuanwang/image-syncer/pkg/log"
	syncjob "github.com/zxzixuanwang/image-syncer/pkg/sync-job"
	"github.com/zxzixuanwang/image-syncer/tools"
	"github.com/zxzixuanwang/image-syncer/web"
)

func main() {
	l := log.NewFileLogger(config.Conf.Sync.FilePosition.Log, config.Conf.App.Env)
	defer func() {
		if err := recover(); err != nil {
			debug.PrintStack()
		}
		syncjob.Save(l)
	}()

	list := make([]string, 0, len(config.Conf.Rabbitmq.Info))
	for _, v := range config.Conf.Rabbitmq.Info {
		list = append(list, mq.RabbitMqUrl(v.Username, v.Password, v.Address))
	}
	r, err := mq.NewRabbitMqSvc(config.Conf.Rabbitmq.QueueName, list, l)
	if err != nil {
		panic(err)
	}

	mqHandler, err := mq.NewMqHandlerSvc(r)
	if err != nil {
		panic(err)
	}
	err = mqHandler.Receive()
	if err != nil {
		panic(err)
	}
	mqHandler.Handle(func(b []byte) error {
		data := new(mq.RabbitMqData)
		err = json.Unmarshal(b, data)
		if err != nil {
			l.Error("unmarshal mq data err", err)
			return err
		}
		imageWholeName := fmt.Sprintf("%s:%s", data.Name, data.Tag)

		l.Debugf("Gotted image is %s", imageWholeName)
		num, order, err := tools.HandleImageName(data.Name, data.Tag, l, tools.HandleDefaultTag)
		if err != nil {
			l.Error("Handle image name err", err)
			return err
		}
		for i := 0; i < len(num); i++ {
			syncjob.Set(fmt.Sprintf("%s&%s", imageWholeName, num[i]), order[i])
		}
		l.Debug("get images", syncjob.Get())

		return nil
	})
	if err := syncjob.Tickerjob(l, config.Conf.Sync.Interval); err != nil {
		panic(err)
	}
	l.Infof("listening port %s", config.Conf.App.Port)
	if err := http.ListenAndServe(config.Conf.App.Port, web.Route(l)); err != nil {
		panic(err)
	}
}
