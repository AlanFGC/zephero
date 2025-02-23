package core

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

type WebSockServer struct {
	conns   map[*websocket.Conn]string
	manager *GameManager
}

func makeWebSockServer(manager *GameManager) *WebSockServer {
	return &WebSockServer{
		conns:   make(map[*websocket.Conn]string),
		manager: manager,
	}
}

func (s *WebSockServer) addConn(conn *websocket.Conn) {
	// TODO add username
	username := conn.RemoteAddr().String()
	s.conns[conn] = username
	playerChan := s.manager.registerPlayer(username, func() {
		log.Println("removing player: ", username)
		s.removeConn(conn)
	})
	s.connectionLoop(conn)
	log.Println("Connection terminated")
	close(playerChan)
}

func (s *WebSockServer) removeConn(conn *websocket.Conn) {
	delete(s.conns, conn)
	err := conn.Close()
	if err != nil {
		fmt.Println("Error closing connection", conn, err)
	}
}

func (s *WebSockServer) connectionLoop(ws *websocket.Conn) {
	buff := make([]byte, 1024)
	for {
		// Input
		buffSize, err := ws.Read(buff)
		if err != nil {
			log.Println("Websockets connection interrupted: ", err)
			return
		}
		msg := buff[:buffSize]

		s.manager.SendEvent(PlayerEvent{
			PlayerId: s.conns[ws],
			GameEvent: GameEvent{
				EventId: string(msg),
				Data:    string(msg),
			},
		})

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
