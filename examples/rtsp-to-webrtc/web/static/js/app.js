let suuid = "demo1";
let baseUrl = "rtsp-api";
// let baseUrl = "";

let config = {
    iceServers: [{
        'url': 'turn:101.200.83.51:3478?transport=udp',
        'credential': 'admin123',
        'username': 'admin'
    }]
};
$('#upload').change(function(){
    // 获取FileList的第一个元素
    var f = document.getElementById('upload').files[0];
    src = window.URL.createObjectURL(f);
    document.getElementById('preview2').src = src
})

function draw(preview){
    var canvas = document.getElementById(preview);
    if (!canvas.getContext) return;
    var ctx = canvas.getContext("2d");
    var img=document.getElementById('preview2');
    ctx.drawImage(img, 0, 0);
}
// document.querySelector("img").onclick = function (){
//
// }

const pc = new RTCPeerConnection(config);
pc.onnegotiationneeded = handleNegotiationNeededEvent;

let log = msg => {
    document.getElementById('div').innerHTML += msg + '<br>'
}

pc.ontrack = function (event) {
    log(event.streams.length + ' track is delivered')
    var el = document.createElement(event.track.kind)
    el.style="width: 650px; height: 300px; background-color:black;"
    el.srcObject = event.streams[0]
    el.muted = true
    el.autoplay = true
    el.controls = true
    el.width = 650
    console.log(el)
    document.getElementById('remoteVideos').appendChild(el)
}

pc.oniceconnectionstatechange = e => log("======" + pc.iceConnectionState)


async function handleNegotiationNeededEvent() {
    let offer = await pc.createOffer();
    await pc.setLocalDescription(offer);
    getRemoteSdp();
}

$(document).ready(function () {
    // $('#' + suuid).addClass('active');
    // connectRtsp();
});

// pc.ondatachannel = function(event) {
//   receiveChannel = event.channel;
//   receiveChannel.onmessage = function(event){
//     alert(event.data);
//   };
// };
//发送弹幕
function send() {
    // var preview = document.getElementById("preview").src;
    //创建一个弹幕模板
    var _html = "";
    //获取发送的弹幕颜色
    var textColor = $("input[name='text_color']:checked").val()
    //获取弹幕内容标签
    var bilibiliTxt = $("#bilibili_content");
    //获取弹幕数据
    var content = bilibiliTxt.val();
    //创建一个p标签ID
    var content_id = "bilibili" + Math.ceil(Math.random() * 500);
    //创建一个置顶div标签ID
    var div_id = "ban2_" + Math.ceil(Math.random() * 500);
    //随机创建一个随机数，使弹幕位置随机生成的效果
    var topVal = Math.ceil(Math.random() * 100);
    //获取发送弹幕的位置
    var textLoc = $("input[name='text_loc']:checked").val();

    //如果选择了置顶弹幕
    if (textLoc == "top") {
        //计算出上下弹幕的间隔
        var a = $(".ban2").size() * 10;
        //实例模板
        _html += "<div id=" + div_id + " class='ban2' style='position: relative; text-align:center; padding-bottom:6px;'>";
        _html += "<p id='" + content_id + "' class='bilibili_txt_2' style='top:" + a + "px; margin-top:0px; width:100%; color:" + textColor + ";'>" + content + "</p>";
        _html += "</div>";
    }
    //如果选择了正常弹幕
    else if (textLoc == "normal") {
        //实例模板
        _html = "<p id='" + content_id + "' class='bilibili_txt_1' style='top:" + topVal + "px; color: white;position:absolute;right:0px; margin-top:0px; color:" + textColor + ";'>" + content + "</p>";
        // _html += "<img  id='" + content_id + "_img' class='bilibili_txt_1'  src='"+preview+"' style='top:" + topVal + "px; color: white;position:absolute;width: 6rem;height: 3rem;right:0px; margin-top:20px; color:" + textColor + ";'>";
        _html += "<canvas  id='" + content_id + "_img' class='bilibili_txt_1'  src='' style='top:" + topVal + "px; color: white;position:absolute;width: 6rem;height: 3rem;right:0px; margin-top:20px; color:" + textColor + ";'>";
    }
    //清空弹幕输入框内容
    bilibiliTxt.val("");
    //添加到弹幕板上
    $("#remoteVideos").append(_html)
    draw(content_id+"_img");
    //调用启动动画效果
    bilibiliAnimation(content_id,div_id);
}
//启动动画
function bilibiliAnimation(id,divid){
  //开始向左前行动画
  $(".bilibili_txt_1").animate({left:'0px'},8000,function(){
    $("#"+id).remove();
    $("#"+id+"_img").remove();
  })
}



function getCodecInfo() {
    $.get(baseUrl + "/codec/" + suuid, function (data) {
        try {
            if (data == "") {
                alert("ffmpeg not start,rtsp server not data");
                return
            }
            data = JSON.parse(data);
            if (data.length > 1) {
                log('add audio Transceiver')
                pc.addTransceiver('audio', {
                    'direction': 'sendrecv'
                })
            }
        } catch (e) {
            console.log(e);
        } finally {

            log('add video Transceiver')
            pc.addTransceiver('video', {
                'direction': 'sendrecv'
            });

            //send ping becouse PION not handle RTCSessionDescription.close()
            sendChannel = pc.createDataChannel('foo');
            console.log('foo channel has start');
            console.log("=======", sendChannel.binaryType);
            sendChannel.onclose = () => console.log('sendChannel has closed');
            sendChannel.onerror = function (error) {
                console.log("dataChannel.OnError:", error);
            };
            sendChannel.onopen = function (event) {
                console.log('sendChannel has opened', event);
                sendChannel.send('ping');
                setInterval(() => {
                    sendChannel.send('ping');
                }, 1000)
            }
            sendChannel.onmessage = e => log(`Message from DataChannel '${sendChannel.label}' payload '${e.data}'`);
        }
    });
}

let sendChannel = null;

function connectRtsp() {
    $.post(baseUrl + "/connectRtsp", {
        url: $("#rtspUrl").val()
    }, function (data) {
        try {
        } catch (e) {
            console.warn(e);
        }
    });
}

function getRemoteSdp() {
    $.post(baseUrl + "/recive", {
        suuid: suuid,
        data: btoa(pc.localDescription.sdp),
    }, function (data) {
        try {
            $('#remoteSessionDescription').val(data);
            if (data == "") {
                $('#localSessionDescription').val("");
            } else {
                pc.setRemoteDescription(new RTCSessionDescription({
                    type: 'answer',
                    sdp: atob(data)
                }))
                $('#localSessionDescription').val(btoa(pc.localDescription.sdp));
            }
        } catch (e) {
            console.warn(e);
        }
    });
}
