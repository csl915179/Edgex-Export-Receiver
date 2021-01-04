package Controller

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

const (
	ReadBufferSize  = 1024
	WriteBufferSize = 1024
)

var Upgrader = websocket.Upgrader{
	//读取存储空间的大小
	ReadBufferSize: ReadBufferSize,
	//写入存储空间的大小
	WriteBufferSize: WriteBufferSize,
	//允许跨域
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	Upgrader.Upgrade(w, r, nil)

	//u := r.URL.RawQuery
	//m, _ := url.ParseQuery(u)
	//token := m[common.SessionTokenKey][0]

	//wsConn.clientmapping[token] = ws
	w.Write([]byte("Hello"))
	log.Printf(" %d ws client connected success!")
}

func WsClientSend(token string, message []byte) {

}
