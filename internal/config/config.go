package config

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var C *Config

type Config struct {
	HTTPServer HTTPServer  `yaml:"http_server"`
	Database   SQLDatabase `yaml:"database"`
	User       User        `yaml:"user"`
}

type HTTPServer struct {
	Address string `yaml:"address"`
}

type SQLDatabase struct {
	Driver   string `yaml:"driver"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	DB       string `yaml:"db"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

// String returns SQLDatabase formatted DSN
func (d *SQLDatabase) String() string {
	switch d.Driver {
	case "mysql":
		return d.mysqlDSN()
	}

	panic("SQLDatabase driver is not supported")
}

func (d *SQLDatabase) mysqlDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&multiStatements=true&interpolateParams=true&collation=utf8mb4_general_ci&clientFoundRows=true", d.User, d.Password, d.Host, d.Port, d.DB)
}

type User struct {
	Secret string `yaml:"secret"`
}

func Init(filename string) {
	c := new(Config)
	v := viper.New()
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetEnvPrefix(envPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	v.AutomaticEnv()

	if err := v.ReadConfig(bytes.NewReader(builtinConfig)); err != nil {
		log.Fatalf("failed on config initialization: %v", err)
	}

	if filename != "" {
		v.SetConfigFile(filename)
		if err := v.MergeInConfig(); err != nil {
			log.Fatalf("opening config file [%s] failed: %v", filename, err)
		} else {
			log.Infof("config file [%s] opened successfully", filename)
		}
	}

	err := v.Unmarshal(c, func(config *mapstructure.DecoderConfig) {
		config.TagName = "yaml"
		config.DecodeHook = mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
		)
	})
	if err != nil {
		log.Fatalf("failed on config unmarshal: %s: %v", filename, err)
	}

	C = c
}
