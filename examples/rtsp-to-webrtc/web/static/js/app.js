let suuid = "demo1";
let baseUrl = "rtsp-api";
// let baseUrl = "";

let config = {
  iceServers: [{
    urls: ["stun:101.200.83.51:3478"]
  }]
};

const pc = new RTCPeerConnection(config);
pc.onnegotiationneeded = handleNegotiationNeededEvent;

let log = msg => {
  document.getElementById('div').innerHTML += msg + '<br>'
}

pc.ontrack = function(event) {
  log(event.streams.length + ' track is delivered')
  var el = document.createElement(event.track.kind)
  el.srcObject = event.streams[0]
  el.muted = true
  el.autoplay = true
  el.controls = true
  el.width = 600
  document.getElementById('remoteVideos').appendChild(el)
}

pc.oniceconnectionstatechange = e => log(pc.iceConnectionState)



async function handleNegotiationNeededEvent() {
  let offer = await pc.createOffer();
  await pc.setLocalDescription(offer);
  console.log("+++++++++++++++++++++++++3444444444")
  getRemoteSdp();
}

$(document).ready(function() {
  // $('#' + suuid).addClass('active');
  connectRtsp();
});

// pc.ondatachannel = function(event) {
//   receiveChannel = event.channel;
//   receiveChannel.onmessage = function(event){
//     alert(event.data);
//   };
// };

function getCodecInfo() {
  $.get(baseUrl+"/codec/" + suuid, function(data) {
    try {
      if(data == ""){
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
      var dataChannelOptions = {
        reliable: true,
        maxRetransmitTime: "2000"
      };

      //send ping becouse PION not handle RTCSessionDescription.close()
      sendChannel = pc.createDataChannel('foo',dataChannelOptions);
      console.log('foo channel has start');
      sendChannel.onclose = () => console.log('sendChannel has closed');
      sendChannel.onopen = () => {
        console.log('sendChannel has opened');
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
  $.post(baseUrl+"/connectRtsp", {
    url: $("#rtspUrl").val()
  }, function(data) {
    try {
    } catch (e) {
      console.warn(e);
    }
  });
}

function getRemoteSdp() {
  $.post(baseUrl+"/recive", {
    suuid: suuid,
    data: btoa(pc.localDescription.sdp),
  }, function(data) {
    try {
      pc.setRemoteDescription(new RTCSessionDescription({
        type: 'answer',
        sdp: atob(data)
      }))
      $('#remoteSessionDescription').val(data);
      if(data == ""){
        $('#localSessionDescription').val("");
      }else{
        $('#localSessionDescription').val(btoa(pc.localDescription.sdp));
      }
    } catch (e) {
      console.warn(e);
    }
  });
}
