package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout:  0,
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	WriteBufferPool:   nil,
	Subprotocols:      nil,
	Error:             nil,
	CheckOrigin:       nil,
	EnableCompression: false,
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "homePage")
}

func reader(conn *websocket.Conn) {
	// var resp Response
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			logrus.Error(err)
			return
		}

		logrus.Info("Message from client : " + string(p))

		for i := 0; i < 10; i++ {

			//resp.Code = 1
			//resp.Message = "hi from server golang!"
			//resp.Data = nil

			//if err := conn.WriteJSON(resp); err != nil {
			//	logrus.Error(err)
			//	return
			//}

			time.Sleep(5 * time.Second)

			message := []byte(strconv.Itoa(i))

			if err := conn.WriteMessage(1, message); err != nil {
				logrus.Error(err)
				return
			}
		}
	}
}

func wsEndPoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error(err)
		return
	}

	defer conn.Close()

	logrus.Info("Client successfully connected ... ")
	reader(conn)
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndPoint)
}

func main() {
	fmt.Println("Go Websockets")
	setupRoutes()
	logrus.Fatal(http.ListenAndServe(":8080", nil))
}
