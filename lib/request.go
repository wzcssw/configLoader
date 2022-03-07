package lib

import (
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func Get(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		log.Error(err)
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("调用reload接口", string(body))
	return nil
}
