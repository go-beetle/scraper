package scraper

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
)

func WriteFile(document []byte, u *URL) {
	var folderPath, filename string
	timePostfix := time.Now().Format("2006-01-02_15:04:05")

	if len(u.Path) == 1 {
		folderPath = "data/"
		filename = fmt.Sprintf("%s%s(%s)", folderPath, u.Hostname, timePostfix)

	} else {
		suburls := strings.Split(u.Path, "/")
		folder := strings.Join(suburls[:len(suburls)-1], "/")

		folderPath = fmt.Sprintf("data/%s%s", u.Hostname, folder)
		filename = fmt.Sprintf("data/%s%s(%s)", u.Hostname, u.Path, timePostfix)
	}
	zap.S().Debugf("Making folder if not made already %s", folderPath)
	err := os.MkdirAll(folderPath, 0700)
	if err != nil {
		zap.S().Warnf("Make director error %v", err)
	}

	zap.S().Infof("Writing file %s", filename)
	err = ioutil.WriteFile(filename, document, 0700)
	if err != nil {
		zap.S().Warnf("writing error %v, filename %s", err, filename)
	}
}