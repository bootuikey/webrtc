package main

import (
	"fmt"
	"github.com/deepch/vdk/format/rtsp"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	flag := false
	var rtspUrlName string
	for idx, args := range os.Args {
		fmt.Println("参数"+strconv.Itoa(idx)+":", args)
		if idx == 1 {
			flag = true
			rtspUrlName = args
		}
	}
	if rtspUrlName == "" {
		fmt.Println("请输入rtspurl地址")
		return
	}
	if flag {
		go getRtspServeStreams(rtspUrlName)
		go serveHTTP()
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
	log.Println(name, "connect", url)
	rtsp.DebugRtsp = true
	session, err := rtsp.Dial(url)
	if err != nil {
		log.Println(name, err)
		return false
		time.Sleep(5 * time.Second)
	}
	session.RtpKeepAliveTimeout = 10 * time.Second
	if err != nil {
		log.Println(name, err)
		time.Sleep(5 * time.Second)
		return false
	}
	codec, err := session.Streams()
	if err != nil {
		log.Println(name, err)
		time.Sleep(5 * time.Second)
	}
	Config.coAd(name, codec)
	for {
		pkt, err := session.ReadPacket()
		if err != nil {
			log.Println(name, err)
			break
		}
		Config.cast(name, pkt)
	}
	err = session.Close()
	if err != nil {
		log.Println("session Close error", err)
	}
	time.Sleep(5 * time.Second)
	return true
}
