package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"sync"
)

var (
	matches   []string
	waitGroup = sync.WaitGroup{}
	lock      = sync.Mutex{}
)

func fileSearch(r string, s string) {
	fmt.Printf("Searchnig in", r)
	files, _ := ioutil.ReadDir(r)
	for _, f := range files {
		if strings.Contains(f.Name(), s) {
			lock.Lock()
			matches = append(matches, filepath.Join(r, f.Name()))
			lock.Unlock()
		}

		if f.IsDir() {
			waitGroup.Add(1)
			go fileSearch(filepath.Join(r, f.Name()), s)
		}
	}
	waitGroup.Done()
}

func main() {
	waitGroup.Add(1)
	go fileSearch("/Users/andy/Dev/srv/src/rekroot", "readme.md")
	waitGroup.Wait()
	for _, file := range matches {
		fmt.Println("Matched", file)
	}
}
