package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func main() {
	begin := time.Now().UnixNano() / int64(time.Millisecond)
	var path string
	var action string
	var find string
	var replace string

	flag.StringVar(&action, "a", "find", "find replace")
	flag.StringVar(&find, "f", "", "text want to find")
	flag.StringVar(&replace, "r", "", "text want to replace")
	flag.StringVar(&path, "p", ".", "path")

	flag.Parse()

	if find == "" {
		fmt.Printf("finding text is null")
		return
	}

	fileInfo, err := fileInfo(path)
	if err != nil {
		fmt.Println(err)
		return
	}

	var pathList []string

	var wg sync.WaitGroup

	if fileInfo.IsDir() {
		err = filepath.Walk(path, func(subPath string, info os.FileInfo, err error) error {
			if info == nil {
				return err
			}
			if !info.IsDir() {
				wg.Add(1)
				go func() {
					reader, err := ioutil.ReadFile(subPath)
					if err != nil {
						fmt.Println("err: ", err)
						return
					}

					if action == "find" {
						if strings.Contains(string(reader), find) {
							pathList = append(pathList, subPath)
						}
					} else if action == "replace" {
						newReader := strings.ReplaceAll(string(reader), find, replace)
						err = ioutil.WriteFile(subPath, []byte(newReader), 0777)
						if err != nil {
							fmt.Println(err)
							return
						}
					}
					wg.Done()
				}()
			}
			return nil
		})

		if err != nil {
			fmt.Printf("err: %v", err)
			return
		}
	}

	wg.Wait()

	end := time.Now().UnixNano() / int64(time.Millisecond)
	if action == "find" {
		fmt.Printf("find success! total : %d . spend %d ms \n", len(pathList), end-begin)
		for _, path := range pathList {
			fmt.Println(path)
		}
	} else if action == "replace" {
		fmt.Println("replace success!")
	}
}

func fileInfo(path string) (os.FileInfo, error) {
	f, err := os.Stat(path)
	if err == nil {
		return f, nil
	}

	return nil, fmt.Errorf("err: %v", err)
}
