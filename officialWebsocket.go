package main

import (
	"BiliLine_Windows/Global"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/vtb-link/bianka/live"
	"github.com/vtb-link/bianka/proto"
	"log"
	"time"
)

const (
	MaxBodySize     = int32(1 << 11)
	CmdSize         = 4
	PackSize        = 4
	HeaderSize      = 2
	VerSize         = 2
	OperationSize   = 4
	SeqIdSize       = 4
	HeartbeatSize   = 4
	RawHeaderSize   = PackSize + HeaderSize + VerSize + OperationSize + SeqIdSize
	MaxPackSize     = MaxBodySize + int32(RawHeaderSize)
	PackOffset      = 0
	HeaderOffset    = PackOffset + PackSize
	VerOffset       = HeaderOffset + HeaderSize
	OperationOffset = VerOffset + VerSize
	SeqIdOffset     = OperationOffset + OperationSize
	HeartbeatOffset = SeqIdOffset + SeqIdSize
)

const (
	OP_HEARTBEAT       = int32(2)
	OP_HEARTBEAT_REPLY = int32(3)
	OP_SEND_SMS_REPLY  = int32(5)
	OP_AUTH            = int32(7)
	OP_AUTH_REPLY      = int32(8)
)

type WebsocketClient struct {
	conn       *websocket.Conn
	msgBuf     chan *Proto
	sequenceId int32
	dispather  map[int32]protoLogic
	authed     bool
}

type protoLogic func(p *Proto) (err error)

type Proto struct {
	PacketLength int32
	HeaderLength int16
	Version      int16
	Operation    int32
	SequenceId   int32
	Body         []byte
	BodyMuti     [][]byte
}

type AuthRespParam struct {
	Code int64 `json:"code,omitempty"`
}

// StartWebsocket 启动长连
func StartWebsocket(wsAddr, authBody string) (err error) {
	// 建立连接
	conn, _, err := websocket.DefaultDialer.Dial(wsAddr, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	wc := &WebsocketClient{
		conn:      conn,
		msgBuf:    make(chan *Proto, 1024),
		dispather: make(map[int32]protoLogic),
	}

	// 注册分发处理函数
	wc.dispather[OP_AUTH_REPLY] = wc.authResp
	wc.dispather[OP_HEARTBEAT_REPLY] = wc.heartBeatResp
	wc.dispather[OP_SEND_SMS_REPLY] = wc.msgResp

	// 发送鉴权信息
	err = wc.sendAuth(authBody)
	if err != nil {
		return
	}

	// 读取信息
	go wc.ReadMsg()

	// 处理信息
	go wc.DoEvent()

	return
}

// ReadMsg 读取长连信息
func (wc *WebsocketClient) ReadMsg() {
	for {
		retProto := &Proto{}
		_, buf, err := wc.conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			continue
		}
		retProto.PacketLength = int32(binary.BigEndian.Uint32(buf[PackOffset:HeaderOffset]))
		retProto.HeaderLength = int16(binary.BigEndian.Uint16(buf[HeaderOffset:VerOffset]))
		retProto.Version = int16(binary.BigEndian.Uint16(buf[VerOffset:OperationOffset]))
		retProto.Operation = int32(binary.BigEndian.Uint32(buf[OperationOffset:SeqIdOffset]))
		retProto.SequenceId = int32(binary.BigEndian.Uint32(buf[SeqIdOffset:]))
		if retProto.PacketLength < 0 || retProto.PacketLength > MaxPackSize {
			continue
		}
		if retProto.HeaderLength != RawHeaderSize {
			continue
		}
		if bodyLen := int(retProto.PacketLength - int32(retProto.HeaderLength)); bodyLen > 0 {
			retProto.Body = buf[retProto.HeaderLength:retProto.PacketLength]
		} else {
			continue
		}
		retProto.BodyMuti = [][]byte{retProto.Body}
		if len(retProto.BodyMuti) > 0 {
			retProto.Body = retProto.BodyMuti[0]
		}
		wc.msgBuf <- retProto
	}
}

// DoEvent 处理信息
func (wc *WebsocketClient) DoEvent() {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case p := <-wc.msgBuf:
			if p == nil {
				continue
			}
			if wc.dispather[p.Operation] == nil {
				continue
			}
			err := wc.dispather[p.Operation](p)
			if err != nil {
				continue
			}
		case <-ticker.C:
			wc.sendHeartBeat()
		}
	}
}

// sendAuth 发送鉴权
func (wc *WebsocketClient) sendAuth(authBody string) (err error) {
	p := &Proto{
		Operation: OP_AUTH,
		Body:      []byte(authBody),
	}
	return wc.sendMsg(p)
}

// sendHeartBeat 发送心跳
func (wc *WebsocketClient) sendHeartBeat() {
	if !wc.authed {
		return
	}
	msg := &Proto{}
	msg.Operation = OP_HEARTBEAT
	msg.SequenceId = wc.sequenceId
	wc.sequenceId++
	err := wc.sendMsg(msg)
	if err != nil {
		return
	}
	log.Println("[WebsocketClient | sendHeartBeat] seq:", msg.SequenceId)
}

// sendMsg 发送信息
func (wc *WebsocketClient) sendMsg(msg *Proto) (err error) {
	dataBuff := &bytes.Buffer{}
	packLen := int32(RawHeaderSize + len(msg.Body))
	msg.HeaderLength = RawHeaderSize
	binary.Write(dataBuff, binary.BigEndian, packLen)
	binary.Write(dataBuff, binary.BigEndian, int16(RawHeaderSize))
	binary.Write(dataBuff, binary.BigEndian, msg.Version)
	binary.Write(dataBuff, binary.BigEndian, msg.Operation)
	binary.Write(dataBuff, binary.BigEndian, msg.SequenceId)
	binary.Write(dataBuff, binary.BigEndian, msg.Body)
	err = wc.conn.WriteMessage(websocket.BinaryMessage, dataBuff.Bytes())
	if err != nil {
		log.Println("[WebsocketClient | sendMsg] send msg err:", msg)
		return
	}
	return
}

// authResp 鉴权处理函数
func (wc *WebsocketClient) authResp(msg *Proto) (err error) {
	resp := &AuthRespParam{}
	err = json.Unmarshal(msg.Body, resp)
	if err != nil {
		return
	}
	if resp.Code != 0 {
		return
	}
	wc.authed = true
	log.Println("[WebsocketClient | authResp] auth success")
	return
}

// heartBeatResp  心跳结果
func (wc *WebsocketClient) heartBeatResp(msg *Proto) (err error) {
	log.Println("[WebsocketClient | heartBeatResp] get HeartBeat resp", msg.Body)
	return
}

// msgResp 可以这里做回调
func (wc *WebsocketClient) msgResp(msg *Proto) (err error) {
	for index, cmd := range msg.BodyMuti {
		log.Printf("[WebsocketClient | msgResp] recv MsgResp index:%d ver:%d cmd:%s", index, msg.Version, string(cmd))
		cmd, data, err := proto.AutomaticParsingMessageCommand(cmd)
		if err != nil {
			log.Println(err)
			return err
		}
		switch cmd {
		case proto.CmdLiveOpenPlatformDanmu:
			if globalConfiguration.IsOnlyGift {
				break
			}
			DanmuData := data.(*proto.CmdDanmuData)

			Broadcast.Broadcast(DanmuData)

		case proto.CmdLiveOpenPlatformSendGift:
			GiftData := data.(*proto.CmdSendGiftData)

			if !globalConfiguration.AutoJoinGiftLine {
				break
			}
			if line.GiftIndex[GiftData.OpenID] != 0 {
				break
			}

			if !GiftData.Paid {
				break
			}

			lineTemp := GlobalType.GiftLine{
				OpenID: GiftData.OpenID,
				// OpenID:     strconv.Itoa(GiftData.Uid),
				UserName:   GiftData.Uname,
				Avatar:     GiftData.Uface,
				PrintColor: globalConfiguration.GiftPrintColor,
				GiftPrice:  float64(GiftData.Price * GiftData.GiftNum / 1000),
			}
			if (float64(GiftData.GiftNum*GiftData.Price))/1000 >= globalConfiguration.GiftLinePrice {
				line.GiftLine = append(line.GiftLine, lineTemp)
				line.GiftIndex[GiftData.OpenID] = len(line.GiftLine)
				SendLineToWs(GlobalType.Line{}, lineTemp, GlobalType.GiftLineType)
				SetLine(line)
			}
		case live.CmdLiveOpenPlatformGuard:
			log.Println(cmd, data.(*proto.CmdGuardData))
		}
	}
	return
}
