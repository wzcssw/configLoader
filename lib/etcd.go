package lib

import (
	"time"

	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
)

var (
	ETCDCli *clientv3.Client
)

func InitETCDConn(addr string) error {
	//配置
	config := clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: time.Second * 5,
	}
	//连接创建一个客户端
	client, err := clientv3.New(config)
	if err != nil {
		log.Panic(err)
		panic(err)
	}

	ETCDCli = client
	return nil
}
