package config

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	Conf           *Configs
	DefaultConfigs Configs = Configs{
		Registry: Registrys{},
		Auth:     Auth{},
		App: App{
			Port: ":8080",
			Env:  "test",
		},
		Sync: Sync{
			FilePosition: FilePosition{
				Auth:   "./auth.yaml",
				Images: "./images.json",
				Log:    "./app.log",
			},
			RoutineNum: 1,
			RetryCount: 2,
			Interval:   10,
		},
	}
)

type App struct {
	Port string
	Env  string
}

type Configs struct {
	App      App       `json:"app,omitempty"`
	Registry Registrys `json:"registry,omitempty"`
	Auth     Auth      `json:"auth,omitempty"`
	Sync     Sync      `json:"sync,omitempty"`
}

type Sync struct {
	FilePosition FilePosition `json:"file_position,omitempty"`
	RoutineNum   int          `json:"routineNum,omitempty"`
	RetryCount   int          `json:"retryCount,omitempty"`
	Filter       Filter       `json:"filter,omitempty"`
	Interval     int          `json:"interval,omitempty"`
	MaxSyncDes   int          `json:"maxSyncDes,omitempty"`
}

type Filter struct {
	OsFilterList   []string `json:"osFilterList,omitempty"`
	ArchFilterList []string `json:"archFilterList,omitempty"`
}

type FilePosition struct {
	Auth   string `json:"auth,omitempty"`
	Images string `json:"images,omitempty"`
	Log    string `json:"log,omitempty"`
}

type Registrys struct {
	Reg []Registry `json:"reg,omitempty"`
}

type Registry struct {
	SourceName      string `json:"sourceName,omitempty"`
	DestinationName string `json:"destinationName,omitempty"`
	Namespace       string `json:"namespace,omitempty"`
}

type Auth struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

func init() {
	var (
		cfgFile = pflag.StringP("config", "c", "", "config file")
	)
	pflag.Parse()
	if *cfgFile != "" {
		viper.SetConfigFile(*cfgFile)
	} else {
		viper.AddConfigPath("configs")
		viper.SetConfigName("sync")
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	Conf = &DefaultConfigs
	err := viper.Unmarshal(&Conf)
	if err != nil {
		panic(err)
	}
	max := 0
	length := len(Conf.Registry.Reg)

	for i := 0; i < length; i++ {
		Conf.Registry.Reg[i].DestinationName = strings.TrimRight(Conf.Registry.Reg[i].DestinationName, "/")
		Conf.Registry.Reg[i].SourceName = strings.TrimRight(Conf.Registry.Reg[i].SourceName, "/")

		reg := Conf.Registry.Reg[i]

		if strings.Contains(reg.DestinationName, ",") {
			des := strings.Split(reg.DestinationName, ",")
			desR := 0
			for _, ns := range des {
				if desR == 0 {
					Conf.Registry.Reg[i].DestinationName = ns
					desR++
					continue
				}

				Conf.Registry.Reg = append(Conf.Registry.Reg, Registry{
					SourceName:      reg.SourceName,
					DestinationName: ns,
					Namespace:       reg.Namespace,
				})
				desR++
			}
			if desR > max {
				max = desR
			}
		}
	}
	Conf.Sync.MaxSyncDes = max
	fmt.Println(">>>>", Conf.Registry.Reg)
}
