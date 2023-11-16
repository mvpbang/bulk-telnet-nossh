package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
	"net"
	"os"
	"sync"
	"time"
)

// error handle
func errHandle(err error) {
	if err != nil {
		panic(err)
	}
}

// logger
func setLogger() *zap.Logger {
	// create file
	file, err := os.OpenFile("telnet.log", os.O_CREATE|os.O_WRONLY, 0666)
	errHandle(err)
	//defer file.Close()

	// file writer && console writer
	fileWrite := zapcore.AddSync(file)
	consoleWrite := zapcore.AddSync(os.Stdout)
	encoderWrite := zapcore.NewMultiWriteSyncer(fileWrite, consoleWrite)

	// encoder
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{MessageKey: "msg"})

	// core
	core := zapcore.NewCore(encoder, encoderWrite, zap.InfoLevel)

	// logger
	return zap.New(core)
}

// yml文件解析

type Config struct {
	Target []string `yaml:"target"`
}

var config Config

func readYaml() {
	// read yaml
	yfile, err := os.ReadFile("ips.yml")
	errHandle(err)

	// decode
	err = yaml.Unmarshal(yfile, &config)
	errHandle(err)
}

/*
基于tcp连通性判断
*/
func tryConnaddr(addr string, wg *sync.WaitGroup, logger *zap.Logger) {
	defer wg.Done()

	// 设置连接超时时间为5秒
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		//fmt.Printf("连接失败：%s \n", addr)
		logger.Error("失败", zap.String("addr", addr))
		return
	} else {
		//fmt.Printf("连接成功: %s \n", addr)
		logger.Error("成功", zap.String("addr", addr))
	}

	defer conn.Close() // 关闭连接

}

func main() {
	readYaml()
	var wg sync.WaitGroup
	logger := setLogger()

	for _, addr := range config.Target {
		wg.Add(1)
		go tryConnaddr(addr, &wg, logger)
	}
	wg.Wait()
}
