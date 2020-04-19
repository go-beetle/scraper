package scraper

import (
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

func Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	zap.S().Infof("Code: %s\t URL: %s", resp.Status, url)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}