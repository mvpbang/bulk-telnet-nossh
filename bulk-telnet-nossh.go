package main

import (
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

// init
func init() {
	file, err := os.OpenFile("nossh.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	handleErr(err)

	multiWriter := io.MultiWriter(file, os.Stdout)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	log.SetOutput(multiWriter)
}

// handle error
func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	Target []string `yaml:"target"`
}

var config Config

func readYaml() {
	yfile, err := os.ReadFile("ips.yml")
	handleErr(err)

	// decode
	err = yaml.Unmarshal(yfile, &config)
	handleErr(err)
}

/*
基于tcp连通性判断
*/
func tryConnaddr(addr string, wg *sync.WaitGroup) {
	defer wg.Done()

	// 设置连接超时时间为6秒
	conn, err := net.DialTimeout("tcp", addr, 6*time.Second)
	if err != nil {
		log.Printf("失败 %s\n", addr)
		return
	} else {
		log.Printf("成功 %s\n", addr)
	}

	defer conn.Close()
}

func main() {
	readYaml()
	var wg sync.WaitGroup

	for _, addr := range config.Target {
		wg.Add(1)
		go tryConnaddr(addr, &wg)
	}

	wg.Wait()
}
