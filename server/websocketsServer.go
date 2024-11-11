package server

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
)

type WebSockServer struct {
	conns   map[*websocket.Conn]bool
	manager *GameManager
}

func makeWebSockServer(manager *GameManager) *WebSockServer {
	return &WebSockServer{
		conns:   make(map[*websocket.Conn]bool),
		manager: manager,
	}
}

func (s *WebSockServer) addConn(conn *websocket.Conn) {
	fmt.Println("NEW CONNECTION")
	s.conns[conn] = true
	s.connectionLoop(conn)
	s.removeConn(conn)
}

func (s *WebSockServer) removeConn(conn *websocket.Conn) {
	delete(s.conns, conn)
}

func (s *WebSockServer) connectionLoop(ws *websocket.Conn) {
	buff := make([]byte, 1024)
	for {
		buffSize, err := ws.Read(buff)
		if err != nil {
			fmt.Println(err)
			continue
		}
		msg := buff[:buffSize]
		fmt.Println("sending new event via web sockets	")
		s.manager.SendEvent(PlayerEvent{
			PlayerId: string(msg[:3]),
			EventId:  string(msg[3:]),
		})
		ws.Write([]byte("Game manager event received"))
	}
}

func RunWebSocketsServer(gameManager *GameManager) {
	fmt.Println("WebSockServer starting...")
	wsServer := makeWebSockServer(gameManager)
	http.Handle("/ws", websocket.Handler(wsServer.addConn))
	if err := http.ListenAndServe(":3000", nil); err != nil {
		fmt.Println("Server error:", err)
	}
}
