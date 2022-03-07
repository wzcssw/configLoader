package model

import (
	"configLoader/lib"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Prometheus struct {
	Path      string
	ReloadURL string
}

func (p Prometheus) Reload() error {
	return lib.Get(p.Path)
}

func (p Prometheus) WriteConifg(fileName, content string) error {
	err := lib.ReFreshFile(fmt.Sprintf("%s/%s", p.Path, fileName), content)
	if err != nil {
		log.Error("prometheus文件创建失败", fileName, err)
		return err
	}
	return nil
}
