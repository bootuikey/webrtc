package main

import (
	_ "github.com/deepch/vdk/format/rtsp"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	_ "time"
)

func main() {
	// 获取日志文件句柄
	// 已 只写入文件|没有时创建|文件尾部追加 的形式打开这个文件
	logFile, err := os.OpenFile(`go.log`, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	// 设置存储位置
	log.SetOutput(logFile)

	flag := false
	var rtspUrlName string
	for idx, args := range os.Args {
		log.Println("参数"+strconv.Itoa(idx)+":", args)
		if idx == 2 {
			flag = true
			rtspUrlName = args
		}
	}
	if rtspUrlName == "" {
		log.Println("请输入rtspurl地址!!")
		return
	}
	if flag {
		go serveHTTP()
		go serveStreams()
		sigs := make(chan os.Signal, 1)
		done := make(chan bool, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			sig := <-sigs
			log.Println(sig)
			done <- true
		}()
		log.Println("Server Start Awaiting Signal")
		<-done
		log.Println("Exiting")
	}
}
