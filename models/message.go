package models

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
)

// 消息
type Message struct {
	gorm.Model
	FromId   int64  // 发送者
	TargetId int64  // 接收者
	Type     int    // 发送类型：群聊、私聊
	Media    int    // 消息类型：文字、图片、音频
	Content  string // 消息体
	Pic      string
	Url      string
	Desc     string
	Amount   int // 其他数字统计
}

func (table *Message) TableName() string {
	return "message" // 数据库里表的命名
}

type Node struct {
	Conn      *websocket.Conn
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射关系
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwLocker sync.RWMutex

func Chat(writer http.ResponseWriter, request *http.Request) {
	// 1. 获取参数 并检验 token

	//token := query.Get("token")
	query := request.URL.Query()
	id := query.Get("userId")
	userId, _ := strconv.ParseInt(id, 10, 64)
	//targetId := query.Get("targetId")
	//context := query.Get("context")
	//msgType := query.Get("type")
	isvalida := true // checkToken() 待做

	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer, request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 2. 获取conn 连接
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.NonThreadSafe),
	}

	// 3. 用户关系

	// 4. userId 和 Node 绑定关系并加锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	// 5. 完成发送的逻辑
	go sendProc(node)

	// 6. 完成接收的逻辑
	go recvProc(node)
	sendMsg(userId, []byte("欢迎进入聊天系统👏🏻"))

}

func sendProc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

func recvProc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			return
		}
		broadMsg(data)
		fmt.Println("websocket<------", data)
	}
}

var udpSendChan chan []byte = make(chan []byte, 1024)

func broadMsg(data []byte) {
	udpSendChan <- data
}

func init() {
	go udpSendProc()
	go udpRecvProc()
}

// 完成 udp 数据发送协程
func udpSendProc() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 255),
		Port: 3000,
	})
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
	}

	for {
		select {
		case data := <-udpSendChan:
			_, err := conn.Write(data)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}

// 完成 udp 数据接收协程
func udpRecvProc() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   net.IPv4zero,
		Port: 3000,
	})
	defer func(conn *net.UDPConn) {
		_ = conn.Close()
	}(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		var buf [512]byte
		n, err := conn.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
			return
		}
		dispatch(buf[0:n])
	}
}

// 后端调度逻辑处理
func dispatch(data []byte) {
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch msg.Type {
	case 1:
		sendMsg(msg.TargetId, data) // 私信
		// case 2:
		// 	sendGroupMsg() // 群发
		// case 3:
		// 	sendAllMsg() // 广播
		// case 4:

	}
}

func sendMsg(userId int64, msg []byte) {
	// rwLocker.RLocker()  // fatal error: sync: RUnlock of unlocked RWMutex
	rwLocker.RLock()
	node, ok := clientMap[userId]
	rwLocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}

// func sendGroupMsg() {

// }

// func sendAllMsg() {

// }
