package main

import (
	"github.com/deepch/vdk/format/rtsp"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
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
		go getRtspServeStreams(rtspUrlName)
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

func getRtspServeStreams(url string) bool {
	name := "demo1"
	for {
		log.Println(name, "connect", url)
		rtsp.DebugRtsp = true
		session, err := rtsp.DialTimeout(url, 5*time.Second)
		if err != nil {
			log.Println(name, err)
			time.Sleep(5 * time.Second)
			continue
		}
		session.RtpKeepAliveTimeout = 10 * time.Second
		if err != nil {
			log.Println(name, err)
			time.Sleep(5 * time.Second)
			continue
		}
		codec, err := session.Streams()
		if err != nil {
			log.Println(name, err)
			time.Sleep(5 * time.Second)
			continue
		}
		Config.coAd(name, codec)
		for {
			pkt, err := session.ReadPacket()
			if err != nil {
				log.Println(name, err)
				break
			}
			Config.cast(name, pkt)

			log.Println(name, "=====pkt size", len(pkt.Data))
		}
		err = session.Close()
		if err != nil {
			log.Println("session Close error", err)
		}
		log.Println(name, "reconnect wait 20s")
		time.Sleep(20 * time.Second)
	}
}
