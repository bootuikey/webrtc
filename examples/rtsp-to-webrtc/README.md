# RTSPtoWebRTC

RTSP Stream to WebBrowser over WebRTC based on Pion

full native! not use ffmpeg or gstreamer

if you need RTSPtoWSMP4f use https://github.com/deepch/RTSPtoWSMP4f

## Team

Deepch - https://github.com/deepch streaming developer

Dmitry - https://github.com/vdalex25 web developer

Now test work on (chrome, safari, firefox) no MAC OS

![RTSPtoWebRTC image](doc/demo4.png)

## Installation
1.
```bash
go get github.com/deepch/RTSPtoWebRTC
```
2.
```bash
cd src/github.com/deepch/RTSPtoWebRTC
```
3.
```bash
go run .
```
4.
```bash
open web browser http://127.0.0.1:8083 work chrome, safari, firefox
```

## Configuration

### Edit file config.json

format:

```bash
{
  "server": {
    "http_port": ":8083"
  },
  "streams": {
    "demo1": {
      "url": "rtsp://170.93.143.139/rtplive/470011e600ef003a004ee33696235daa"
    },
    "demo2": {
      "url": "rtsp://170.93.143.139/rtplive/470011e600ef003a004ee33696235daa"
    },
    "demo3": {
      "url": "rtsp://170.93.143.139/rtplive/470011e600ef003a004ee33696235daa"
    }
  }
}
```

## Limitations

Video Codecs Supported: H264

Audio Codecs Supported: pcm alaw and pcm mulaw 





执行命令将静态文件打包成go文件

1:go-bindata -o=./asset/asset.go -pkg=asset config/...
2:打包命令: go build
3:启动参数：RTSPtoWebRTC.exe -rtspurl rtsp://101.200.83.51/test

#config 配置文件打包命令
go-bindata -pkg=asset -o=./asset/bindata.go ./config/

#ffmpeg 发送视频流命令
ffmpeg -re -i D:\a.mp4 -rtsp_transport tcp -vcodec h264 -f rtsp rtsp://101.200.83.51/test
ffmpeg -re -i D:\a.webm -rtsp_transport tcp -vcodec h264 -f rtsp rtsp://101.200.83.51/test
