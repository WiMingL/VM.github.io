package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/naoina/toml"
	"io/ioutil"
	"os"
	"sync"
)

type AppConfig struct {
	DownloadHost string
	LocalDir     string
	Port         string
	Domain       string
	Appid		 string
	Appsecret    string
}
type Token struct {
	key string
	timeout int64
	sync.RWMutex
}

var App AppConfig
var token *Token


func initConfig(file string) error {
	// 读取配置文件
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	if err := toml.Unmarshal(buf, &App); err != nil {
		return err
	}
	return nil
}

func main() {
	token = &Token{}

	err := initConfig("../weboffice-demo.conf")
	if err != nil {
		fmt.Println("init config faild: %v", err.Error())
		os.Exit(2)
	}
	rDefault := gin.Default()
	InitRouter(rDefault)
	rDefault.Run(":" + App.Port)
}
