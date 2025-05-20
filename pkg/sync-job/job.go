package syncjob

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/zxzixuanwang/image-syncer/config"
	"github.com/zxzixuanwang/image-syncer/pkg/client"
	"github.com/zxzixuanwang/image-syncer/tools"
)

func Tickerjob(l *logrus.Logger, interval int) error {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan,
		os.Interrupt,
		syscall.SIGALRM,
		syscall.SIGHUP,
		// syscall.SIGINFO, this causes windows to fail
		syscall.SIGINT,
		// syscall.SIGQUIT, // Quit from keyboard, "kill -3"
	)
	jsonImages, err := readFile(config.Conf.Sync.FilePosition.Images)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		l.Error("read file error", err)
		return err
	}
	l.Debug("got json images", string(jsonImages))

	if len(jsonImages) > 0 {
		var images map[string]int
		err = json.Unmarshal(jsonImages, &images)
		if err != nil {
			l.Error("unmarshal init config err", err)
			return err
		}
		reOrder(images, l)
		RangeSet(images)
		l.Debug("set ", Get())

	}
	doTick := time.NewTicker(time.Minute * time.Duration(interval))
	saveTick := time.NewTicker(time.Minute * 3)
	go func() {
		for {
			select {
			case <-doTick.C:
				if err = doJob(l); err != nil {
					l.Error("do job err", err)
				}
			case <-saveTick.C:
				if err = save(l); err != nil {
					l.Error("save job err", err)
				}
			}
		}
	}()
	go func() {
		for range signalChan {
			if err = save(l); err != nil {
				l.Error("save job err when close", err)
			}
			os.Exit(0)
		}
	}()

	return nil
}

func reOrder(images map[string]int, l *logrus.Logger) {
	l.Infoln("reordering...")
	for k := range images {
		desCheckSlice := strings.Split(k, "&")
		if len(desCheckSlice) < 2 {
			l.Warn("get order source err,images is ", k)
			continue
		}
		for i := 0; i < len(config.Conf.Registry.Reg); i++ {
			if desCheckSlice[1] == config.Conf.Registry.Reg[i].DestinationName && strings.HasPrefix(k, config.Conf.Registry.Reg[i].SourceName) {
				images[k] = i
				break
			}

		}
	}
}

func Save(l *logrus.Logger) {
	err := save(l)
	if err != nil {
		l.Error("out save err", err)
	}
}

func save(l *logrus.Logger) error {
	var (
		err     error
		content []byte
	)

	images := Get()
	defer func() {
		if err != nil {
			if len(images) > 0 {
				l.Warn("sync error ,retry doing, get is", images)
			}
			return
		}
		if err := recover(); err != nil {
			debug.PrintStack()
		}
	}()
	l.Debug("get images ", images)

	content, err = json.Marshal(images)
	if err != nil {
		l.Error("marsh error", err)
		return err
	}

	err = writeFile(config.Conf.Sync.FilePosition.Images, content)
	if err != nil {
		return err
	}
	return nil
}

func doJob(l *logrus.Logger) error {
	l.Debug("doing job")
	getAll := Get()
	l.Debug("get all images", getAll)
	if len(getAll) < 1 {
		return nil
	}
	tempCheck := make(map[string]bool, config.Conf.Sync.MaxSyncDes)
	tempChange := make(map[string]string, len(getAll))
	for k, v := range getAll {
		desName := config.Conf.Registry.Reg[v].DestinationName
		tempCheck[desName] = true
		tempSlice := strings.Split(k, desName)
		if len(tempSlice) != 2 {
			l.Error("get slice error", k)
			continue
		}

		tempNameSlice := strings.SplitN(tempSlice[0], "/", 2)
		if len(tempNameSlice) != 2 {
			l.Error("get temp slice error", tempSlice[0])
			continue
		}
		tempChange[fmt.Sprintf("%s/%s", desName, strings.TrimRight(tempNameSlice[1], "&"))] = strings.TrimRight(tempSlice[0], "&")
	}

	l.Debug("i get change ", tempChange)

	for k := range tempCheck {
		syncMap := make(map[string]interface{}, len(tempChange))
		for ck, v := range tempChange {
			desOrSlice := strings.SplitN(ck, "/", 2)
			if len(desOrSlice) != 2 {
				l.Error("get temp slice error", v)
				continue
			}
			if k == desOrSlice[0] {
				syncMap[v] = ck
			}
		}

		l.Info("doing sync: map is ", syncMap)
		syncclient, err := client.NewSyncClient(&client.SyncClientIn{
			AuthFile:             config.Conf.Sync.FilePosition.Auth,
			ImageFile:            config.Conf.Sync.FilePosition.Images,
			LogFile:              tools.Point(config.Conf.Sync.FilePosition.Log),
			RoutineNum:           config.Conf.Sync.RoutineNum,
			Retries:              config.Conf.Sync.RetryCount,
			OsFilterList:         config.Conf.Sync.Filter.OsFilterList,
			ArchFilterList:       config.Conf.Sync.Filter.ArchFilterList,
			L:                    l,
			ConfigFile:           "",
			DefaultDestRegistry:  "",
			DefaultDestNamespace: "",
			ImageList:            syncMap,
		})
		if err != nil {
			l.Error("new client error", err)
			return err
		}
		err = syncclient.Run()
		if err != nil {
			l.Error("run sync err", err)
			return err
		}
	}
	Clean(getAll)
	return nil
}

func writeFile(fileName string, content []byte) error {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.Write(content); err != nil {
		return err
	}
	return nil
}

func readFile(fileName string) ([]byte, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}
