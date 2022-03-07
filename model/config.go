package model

type Config struct {
	Prometheus   Prometheus   `yaml:"prometheus"`
	AlertManager AlertManager `yaml:"alertmanager"`
}
