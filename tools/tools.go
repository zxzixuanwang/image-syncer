package tools

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/zxzixuanwang/image-syncer/config"
)

func HandleImageName(name, tagName string, l *logrus.Logger, fun F) ([]string, []int, error) {
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

type F func(tag string) error

func HandleDefaultTag(tag string) error {
	tagLow := strings.ToLower(tag)
	if strings.HasPrefix(tagLow, "snapshot") || strings.HasSuffix(tagLow, "snapshot") {
		return fmt.Errorf("has invalid tag")
	}
	return nil
}
