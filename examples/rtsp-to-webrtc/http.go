package main

import (
	"encoding/base64"
	"encoding/json"
	"github.com/deepch/vdk/format/rtsp"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/deepch/vdk/av"

	"github.com/deepch/vdk/codec/h264parser"
	"github.com/gin-gonic/gin"
	"github.com/pion/webrtc/v2"
	"github.com/pion/webrtc/v2/pkg/media"
)

func serveHTTP() {
	router := gin.Default()
	//router.LoadHTMLGlob("web/template/*")
	//router.GET("/index", func(c *gin.Context) {
	//	fi, all := Config.list()
	//	sort.Strings(all)
	//	port := Config.Server.HTTPPort
	//	time := time.Now().String()
	//	fmt.Println("============="+fi);
	//	//c.HTML(http.StatusOK, "index.html", gin.H{
	//	//	"port":     Config.Server.HTTPPort,
	//	//	"suuid":    fi,
	//	//	"suuidMap": all,
	//	//	"version":  time.Now().String(),
	//	//})
	//	var dict map[string]string //定义dict为map类型
	//	dict = make(map[string]string) //让dict可编辑
	//	dict["port"] = port
	//	dict["suuid"] = fi
	//	//dict["suuidMap"] = all
	//	dict["version"] = time
	//	c.Writer.Write(buildJson(dict));
	//})
	router.GET("/player/:suuid", func(c *gin.Context) {
		_, all := Config.list()
		sort.Strings(all)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"port":     Config.Server.HTTPPort,
			"suuid":    c.Param("suuid"),
			"suuidMap": all,
			"version":  time.Now().String(),
		})
	})
	router.POST("/connectRtsp", connectRtspUrl)
	router.POST("/recive", reciver)
	router.GET("/codec/:uuid", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		if Config.ext(c.Param("uuid")) {
			codecs := Config.coGe(c.Param("uuid"))
			if codecs == nil {
				return
			}
			b, err := json.Marshal(codecs)
			if err == nil {
				_, err = c.Writer.Write(b)
				if err == nil {
					log.Println("Write Codec Info error", err)
					return
				}
			}
		}
	})
	//router.StaticFS("/static", http.Dir("web/static"))
	err := router.Run(Config.Server.HTTPPort)
	if err != nil {
		log.Fatalln("Start HTTP Server error", err)
	}
}

/**
启动连接rtsp
*/
func connectRtspUrl(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	url := c.PostForm("url")
	serveStreams1(url)
}

func reciver(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	data := c.PostForm("data")
	suuid := c.PostForm("suuid")

	log.Println("Request", suuid)
	if Config.ext(suuid) {
		/*

			Get Codecs INFO

		*/
		codecs := Config.coGe(suuid)
		if codecs == nil {
			log.Println("Codec error")
			return
		}
		sps := codecs[0].(h264parser.CodecData).SPS()
		pps := codecs[0].(h264parser.CodecData).PPS()
		/*

			Recive Remote SDP as Base64

		*/
		sd, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			log.Println("DecodeString error", err)
			return
		}
		/*

			Create Media MediaEngine

		*/

		mediaEngine := webrtc.MediaEngine{}
		offer := webrtc.SessionDescription{
			Type: webrtc.SDPTypeOffer,
			SDP:  string(sd),
		}
		err = mediaEngine.PopulateFromSDP(offer)
		if err != nil {
			log.Println("PopulateFromSDP error", err)
			return
		}

		var payloadType uint8
		for _, videoCodec := range mediaEngine.GetCodecsByKind(webrtc.RTPCodecTypeVideo) {
			if videoCodec.Name == "H264" && strings.Contains(videoCodec.SDPFmtpLine, "packetization-mode=1") {
				payloadType = videoCodec.PayloadType
				break
			}
		}
		if payloadType == 0 {
			log.Println("Remote peer does not support H264")
			return
		}
		if payloadType != 126 {
			log.Println("Video might not work with codec", payloadType)
		}
		log.Println("Work payloadType", payloadType)
		api := webrtc.NewAPI(webrtc.WithMediaEngine(mediaEngine))

		peerConnection, err := api.NewPeerConnection(webrtc.Configuration{
			ICEServers: []webrtc.ICEServer{
				{
					URLs: []string{"turn:101.200.83.51:3478"}, Username: "admin", Credential: "admin123",
				},
			},
		})
		if err != nil {
			log.Println("NewPeerConnection error", err)
			return
		}
		log.Println("===============peerConnection state", peerConnection.ConnectionState())

		/*

			ADD KeepAlive Timer

		*/
		timer1 := time.NewTimer(time.Second * 2)
		peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
			// Register text message handling
			d.OnMessage(func(msg webrtc.DataChannelMessage) {
				//fmt.Printf("Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
				timer1.Reset(2 * time.Second)
			})
		})
		/*

			ADD Video Track

		*/
		videoTrack, err := peerConnection.NewTrack(payloadType, rand.Uint32(), "video", suuid+"_pion")
		if err != nil {
			log.Fatalln("NewTrack", err)
		}
		_, err = peerConnection.AddTransceiverFromTrack(videoTrack,
			webrtc.RtpTransceiverInit{
				Direction: webrtc.RTPTransceiverDirectionSendonly,
			},
		)
		if err != nil {
			log.Println("AddTransceiverFromTrack error", err)
			return
		}
		_, err = peerConnection.AddTrack(videoTrack)
		if err != nil {
			log.Println("AddTrack error", err)
			return
		}
		/*

			ADD Audio Track

		*/
		var audioTrack *webrtc.Track
		if len(codecs) > 1 && (codecs[1].Type() == av.PCM_ALAW || codecs[1].Type() == av.PCM_MULAW) {
			switch codecs[1].Type() {
			case av.PCM_ALAW:
				audioTrack, err = peerConnection.NewTrack(webrtc.DefaultPayloadTypePCMA, rand.Uint32(), "audio", suuid+"audio")
			case av.PCM_MULAW:
				audioTrack, err = peerConnection.NewTrack(webrtc.DefaultPayloadTypePCMU, rand.Uint32(), "audio", suuid+"audio")
			}
			if err != nil {
				log.Println(err)
				return
			}
			_, err = peerConnection.AddTransceiverFromTrack(audioTrack,
				webrtc.RtpTransceiverInit{
					Direction: webrtc.RTPTransceiverDirectionSendonly,
				},
			)
			if err != nil {
				log.Println("AddTransceiverFromTrack error", err)
				return
			}
			_, err = peerConnection.AddTrack(audioTrack)
			if err != nil {
				log.Println(err)
				return
			}
		}
		if err := peerConnection.SetRemoteDescription(offer); err != nil {
			log.Println("SetRemoteDescription error", err, offer.SDP)
			return
		}
		answer, err := peerConnection.CreateAnswer(nil)
		if err != nil {
			log.Println("CreateAnswer error", err)
			return
		}

		if err = peerConnection.SetLocalDescription(answer); err != nil {
			log.Println("SetLocalDescription error", err)
			return
		}
		_, err = c.Writer.Write([]byte(base64.StdEncoding.EncodeToString([]byte(answer.SDP))))
		if err != nil {
			log.Println("Writer SDP error", err)
			return
		}
		control := make(chan bool, 10)
		//peerConnection.OnICEConnectionStateChange(func(connectionState webrtc.ICEConnectionState) {
		//	log.Printf("Connection State has changed %s \n", connectionState.String())
		//	if connectionState != webrtc.ICEConnectionStateConnected {
		//		log.Println("Client Close Exit")
		//		err := peerConnection.Close()
		//		if err != nil {
		//			log.Println("peerConnection Close error", err)
		//		}
		//		control <- true
		//		return
		//	}
		//	if connectionState == webrtc.ICEConnectionStateConnected {
		go func() {
			cuuid, ch := Config.clAd(suuid)
			log.Println("start stream", suuid, "client", cuuid)
			defer func() {
				log.Println("stop stream", suuid, "client", cuuid)
				defer Config.clDe(suuid, cuuid)
			}()
			var Vpre time.Duration
			var start bool
			timer1.Reset(5 * time.Second)
			for {
				select {
				case <-timer1.C:
					log.Println("Client Close Keep-Alive Timer")
					peerConnection.Close()
				case <-control:
					return
				case pck := <-ch:
					//timer1.Reset(2 * time.Second)
					if pck.IsKeyFrame {
						start = true
					}
					if !start {
						continue
					}
					if pck.IsKeyFrame {
						pck.Data = append([]byte("\000\000\001"+string(sps)+"\000\000\001"+string(pps)+"\000\000\001"), pck.Data[4:]...)

					} else {
						pck.Data = pck.Data[4:]
					}
					var Vts time.Duration
					if pck.Idx == 0 && videoTrack != nil {
						if Vpre != 0 {
							Vts = pck.Time - Vpre
						}
						samples := uint32(90000 / 1000 * Vts.Milliseconds())
						err := videoTrack.WriteSample(media.Sample{Data: pck.Data, Samples: samples})
						if err != nil {
							return
						}
						Vpre = pck.Time
					} else if pck.Idx == 1 && audioTrack != nil {
						err := audioTrack.WriteSample(media.Sample{Data: pck.Data, Samples: uint32(len(pck.Data))})
						if err != nil {
							return
						}
					}
				}
			}

		}()
		//	}
		//})
		return
	}
}

func serveStreams1(url string) bool {
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
