package config

import (
	"configLoader/model"
	"io/ioutil"

	log "github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Prefix   string                      `yaml:"prefix"`
	Etcd     string                      `yaml:"etcd"`
	Services map[string]model.Prometheus `yaml:"services"` // 这里的Prometheus需要抽象一下
}

var (
	Conf = new(Config)
	// Prefix          = "/monitor"
	// ETCDAddr string = "10.69.77.193:9379" // etcd 地址
)

func LoadConfig() {
	ymlfile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	err = yaml.Unmarshal(ymlfile, Conf)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	log.Infof("config >>> %+v", Conf)
}
