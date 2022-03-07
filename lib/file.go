package lib

import (
	"io/ioutil"
)

func ReFreshFile(fileName, content string) error {
	return ioutil.WriteFile(fileName, []byte(content), 0666)
}
