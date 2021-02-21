package file

import (
	"github.com/hanson/textfinder/config"
	"io/ioutil"
	"log"
	"strings"
	"sync"
)

type DirLevel struct {
	Path  string
	Level int64
}

func OpenDir(wg *sync.WaitGroup, path string, level int64, dirCh chan *DirLevel, fileCh chan string) {

	if config.Cfg.Debug {
		log.Printf("opening dir %s", path)
	}

	rd, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, fi := range rd {
		if fi.IsDir() {
			if config.Cfg.Level == -1 || config.Cfg.Level >= level {
				dl := &DirLevel{
					Path:  path + "/" + fi.Name(),
					Level: level + 1,
				}
				wg.Add(1)
				dirCh <- dl
			}
		} else {
			exceptFlag := false
			if len(config.Cfg.ExceptList) > 0 {
				for _, except := range config.Cfg.ExceptList {
					if strings.Contains(fi.Name(), except) {
						exceptFlag = true
						break
					}
				}
			}

			onlyFlag := true
			if len(config.Cfg.OnlyList) > 0 {
				onlyFlag = false
				for _, only := range config.Cfg.OnlyList {
					if strings.Contains(fi.Name(), only) {
						onlyFlag = true
						break
					}
				}
			}

			if exceptFlag {
				continue
			}

			if !onlyFlag {
				continue
			}

			fileCh <- path + "/" + fi.Name()
			wg.Add(1)
		}
	}
}

func OperateFile(path string) bool {

	if config.Cfg.Debug {
		log.Printf("opening file %s", path)
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	if config.Cfg.Action == "find" {
		if strings.Contains(string(content), config.Cfg.Find) {
			return true
		}
	} else if config.Cfg.Action == "replace" {
		newReader := strings.ReplaceAll(string(content), config.Cfg.Find, config.Cfg.Replace)
		err = ioutil.WriteFile(path, []byte(newReader), 0777)
		if err != nil {
			log.Fatal(err)
		}

		return true
	}

	return false
}
