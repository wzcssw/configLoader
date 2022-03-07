package model

type Service interface {
	Reload() error                              // 配置重载
	WriteConifg(fileName, content string) error // 写配置文件
}

type AlertManager Prometheus
