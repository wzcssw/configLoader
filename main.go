package main

import (
	"configLoader/lib"
	"configLoader/model"
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
	"gopkg.in/yaml.v2"
)

var (
	Config          = make(map[string]model.Prometheus)
	Prefix          = "monitor"
	ETCDAddr string = "10.69.77.193:9379" // etcd 地址
	Region   string = "bj"                // Region etcd前缀 例北京：bj

	serviceMap = map[string]model.Service{ // 可以写到配置文件中
		"prometheus": model.Prometheus{
			Path:      "/Users/oushisei/Desktop/go_workspace/src/configLoader",
			ReloadURL: "http://127.0.0.1:9090/-/reload",
		},
		"alertmanager": model.Prometheus{
			Path:      "/home/wangzhicheng1/alertmanager",
			ReloadURL: "http://127.0.0.1:9093/-/reload",
		},
	}
)

func LoadConfig() {
	ymlfile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	err = yaml.Unmarshal(ymlfile, Config)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

func main() {
	LoadConfig()
	////
	var jobChan = make(chan *clientv3.Event, 1024)

	lib.InitETCDConn(ETCDAddr)

	watcher := clientv3.NewWatcher(lib.ETCDCli)
	watchRespChan := watcher.Watch(context.TODO(), fmt.Sprintf("/%s/%s", Prefix, Region), clientv3.WithPrefix())

	go func() {
		for {
			event := <-jobChan
			DoWork(string(event.Kv.Key), string(event.Kv.Value))
			time.Sleep(3 * time.Second)
		}
	}()

	// 处理kv变化事件
	for watchResp := range watchRespChan {
		for _, event := range watchResp.Events {
			switch event.Type {
			case 0:
				log.Info("key:", string(event.Kv.Key), "更新", string(event.Kv.Value))
				jobChan <- event
			case 1:
				log.Info("key:", string(event.Kv.Key), "删除", "Revision:", event.Kv.ModRevision)
			}
		}
	}

}

func DoWork(key, value string) { // key例子 /monitor/bj/prometheus/prometheus.yml
	strs := strings.Split(strings.TrimPrefix(key, "/"), "/")
	if len(strs) < 4 {
		log.Error("无效的key", key)
	}

	serviceName := strs[2]

	service, exist := Config[serviceName]
	if !exist {
		log.Errorf("服务(%s)不存在", serviceName)
	}

	fileName := strs[len(strs)-1]

	service.WriteConifg(fileName, value)

	service.Reload()
}

// etcdctl put --endpoints=http://10.69.77.193:9379 /cron/jobs/job1 "hello world"
