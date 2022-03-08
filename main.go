package main

import (
	"configLoader/config"
	"configLoader/lib"
	"context"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
)

func main() {
	config.LoadConfig()

	var jobChan = make(chan *clientv3.Event, 1024)

	lib.InitETCDConn(config.Conf.Etcd)

	watcher := clientv3.NewWatcher(lib.ETCDCli)
	watchRespChan := watcher.Watch(context.TODO(), config.Conf.Prefix, clientv3.WithPrefix())

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
				// TODO 完善删除操作
				log.Info("key:", string(event.Kv.Key), "删除", "Revision:", event.Kv.ModRevision)
			}
		}
	}

}

func DoWork(key, value string) { // key例子 /monitor/bj/prometheus/prometheus.yml
	strs := strings.Split(strings.TrimPrefix(key, "/"), "/")
	if len(strs) < 4 {
		log.Error("无效的key: ", key)
		return
	}

	serviceName := strs[2]

	service, exist := config.Conf.Services[serviceName]
	if !exist {
		log.Errorf("服务(%s)不存在", serviceName)
	}

	fileName := strs[len(strs)-1]

	service.WriteConifg(fileName, value)

	service.Reload()
}

// etcdctl put --endpoints=http://10.69.77.193:9379 /cron/jobs/job1 "hello world"
