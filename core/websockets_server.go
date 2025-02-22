package core

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
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
			log.Println("Error reading from websocket:", err)
			return
		}
		msg := buff[:buffSize]

		s.manager.SendEvent(PlayerEvent{
			PlayerId: string(msg),
			GameEvent: GameEvent{
				EventId: string(msg),
				Data:    string(msg),
			},
		})

		playerView := s.manager.access.playerView(120, 12)
		jsonData, err := json.Marshal(playerView)
		if err != nil {
			log.Println("Failed to marshal json")
		}

		_, err = ws.Write(jsonData)
		if err != nil {
			log.Println("Failed to write response to socket")
		}
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
