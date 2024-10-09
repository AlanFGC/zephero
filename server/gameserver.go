package server

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
)

type GameServer struct {
	conns map[*websocket.Conn]bool
}

func makeGameServer() *GameServer {
	return &GameServer{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *GameServer) addConn(conn *websocket.Conn) {
	s.conns[conn] = true
	s.connectionLoop(conn)
	s.removeConn(conn)
}

func (s *GameServer) removeConn(conn *websocket.Conn) {
	delete(s.conns, conn)
}

func (s *GameServer) connectionLoop(ws *websocket.Conn) {
	buff := make([]byte, 1024)
	for {
		buffSize, err := ws.Read(buff)
		if err != nil {
			fmt.Println(err)
			continue
		}
		msg := buff[:buffSize]
		fmt.Println(string(msg))
		if string(msg) == "ping" {
			fmt.Println("sending back message:")
			ws.Write([]byte("pong"))
		}
	}
}

func Start() {
	fmt.Println("GameServer starting...")
	server := makeGameServer()
	http.Handle("/ws", websocket.Handler(server.addConn))
	if err := http.ListenAndServe(":3000", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
