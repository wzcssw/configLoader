package model

import (
	"configLoader/lib"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Prometheus struct {
	Path      string `yaml:"path"`
	ReloadURL string `yaml:"reload_url"`
}

func (p Prometheus) Reload() error {
	return lib.Post(p.ReloadURL)
}

func (p Prometheus) WriteConifg(fileName, content string) error {
	err := lib.ReFreshFile(fmt.Sprintf("%s/%s", p.Path, fileName), content)
	if err != nil {
		log.Error("prometheus文件创建失败", fileName, err)
		return err
	}
	return nil
}
