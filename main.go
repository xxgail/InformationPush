package main

import (
	"InformationPush/lib/mysqllib"
	"InformationPush/lib/redislib"
	"InformationPush/routers"
	"InformationPush/timer/task"
	worker2 "InformationPush/worker"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	initConfig()
	initRedis()
	initMysql()
	initFile()

	// 初始化路由
	router := gin.Default()
	routers.Init(router)

	task.PushTaskInit()

	go open()
	go work()

	httpPort := viper.GetString("app.httpPort")
	http.ListenAndServe(":"+httpPort, router)
}

func initConfig() {
	viper.SetConfigName("config/app")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func initMysql() {
	mysqllib.InitDB()
}

func initRedis() {
	redislib.InitClient()
}

func initFile() {
	gin.DisableConsoleColor()

	logFile := viper.GetString("app.logFile")
	fmt.Println(logFile)
	f, _ := os.Create(logFile)
	gin.DefaultErrorWriter = io.MultiWriter(f)
}

func open() {
	time.Sleep(1000 * time.Millisecond)

	httpUrl := viper.GetString("app.httpUrl")
	httpUrl = "http://" + httpUrl + "/home/index"

	fmt.Println("访问页面体验：", httpUrl)
	fmt.Errorf("ffffffffff")

	cmd := exec.Command("open", httpUrl)
	cmd.Output()
}

func work() {
	_ = worker2.NewAsyncTaskWorker()
}
