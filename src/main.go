package main

import (
	"flag"
	"github.com/hanson/textfinder/config"
	"github.com/hanson/textfinder/file"
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

type info struct {
	begin int64
}

func main() {
	i := &info{
		begin: time.Now().UnixNano() / int64(time.Millisecond),
	}

	flag.StringVar(&config.Cfg.Action, "a", "find", "find replace")
	flag.StringVar(&config.Cfg.Find, "f", "", "text want to find")
	flag.StringVar(&config.Cfg.Replace, "r", "", "text want to replace")
	flag.StringVar(&config.Cfg.Path, "p", ".", "path")
	flag.StringVar(&config.Cfg.Except, "e", "", "file name except")
	flag.StringVar(&config.Cfg.Only, "o", "", "only file name include")
	flag.Int64Var(&config.Cfg.Level, "l", 0, "how deep of dir that you operate, -1 for all")
	flag.BoolVar(&config.Cfg.Debug, "d", false, "echo info")

	flag.Parse()

	if config.Cfg.Find == "" {
		log.Fatal("finding text is null")
	}

	if config.Cfg.Except != "" {
		config.Cfg.ExceptList = strings.Split(config.Cfg.Except, ",")
	}

	if config.Cfg.Only != "" {
		config.Cfg.OnlyList = strings.Split(config.Cfg.Only, ",")
	}

	fi, err := os.Stat(config.Cfg.Path)
	if err != nil {
		log.Fatal(err)
	}

	dirCh := make(chan *file.DirLevel, 1000)
	fileCh := make(chan string, 1000)

	var wg sync.WaitGroup

	if fi.IsDir() {
		file.OpenDir(&wg, config.Cfg.Path, 0, dirCh, fileCh)
	} else {
		file.OperateFile(config.Cfg.Path)
	}

	go func() {
		for c := range dirCh {
			file.OpenDir(&wg, c.Path, c.Level, dirCh, fileCh)
			wg.Done()
		}
	}()

	var fileList []string

	go func() {
		for c := range fileCh {
			if file.OperateFile(c) {
				fileList = append(fileList, c)
			}
			wg.Done()
		}
	}()

	wg.Wait()

	printResult(i, fileList)

	return
}

func printResult(i *info, fileList []string) {

	end := time.Now().UnixNano() / int64(time.Millisecond)

	for _, file := range fileList {
		log.Println("file path:" + file)
	}

	log.Printf("total: %d", len(fileList))

	spent := end - i.begin
	log.Printf("spent %d ms\n", spent)

	if config.Cfg.Action == "find" {
		log.Println("find success!")
	} else if config.Cfg.Action == "replace" {
		log.Println("replace success!")
	}
}
