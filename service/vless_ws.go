package service

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/arlettebrook/serv00-ct8/models"

	"github.com/google/uuid"

	"github.com/gorilla/websocket"
)

// 全局的 websocket Upgrader
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 验证消息中的 UUID
func verifyUUID(msg []byte, expectedID uuid.UUID) bool {
	if len(msg) < 17 {
		return false
	}
	id := msg[1:17]
	return uuid.UUID(id) == expectedID
}

func parseTarget(msg []byte) (string, int, int, error) {
	i := int(msg[17]) + 19
	if len(msg) < i+2 {
		return "", 0, 0, fmt.Errorf("message too short")
	}
	targetPort := int(binary.BigEndian.Uint16(msg[i:]))
	i += 2
	if len(msg) <= i {
		return "", 0, 0, fmt.Errorf("message too short")
	}
	ATYP := msg[i]
	i++

	var host string
	switch ATYP {
	case 1:
		if len(msg) < i+4 {
			return "", 0, 0, fmt.Errorf("message too short")
		}
		host = fmt.Sprintf("%d.%d.%d.%d", msg[i], msg[i+1], msg[i+2], msg[1+3])
		i += 4
	case 2:
		if len(msg) <= i {
			return "", 0, 0, fmt.Errorf("message too short")
		}
		length := int(msg[i])
		i++
		if len(msg) < i+length {
			return "", 0, 0, fmt.Errorf("message too short")
		}
		host = string(msg[i : i+length])
		i += length
	case 3:
		if len(msg) < i+16 {
			return "", 0, 0, fmt.Errorf("message too short")
		}
		host = net.IP(msg[i : i+16]).String()
		i += 16
	default:
		return "", 0, 0, fmt.Errorf("unknown address type")
	}

	return host, targetPort, i, nil
}

// 从 TCP 转发到 WebSocket
func forwardTCPToWebSocket(conn *websocket.Conn, tcpConn net.Conn) {
	defer conn.Close()
	defer tcpConn.Close()

	buf := make([]byte, 1024)
	for {
		n, err := tcpConn.Read(buf)
		if err != nil {
			log.Printf("Read from TCP connection error: %s", err)
			return
		}

		if err := conn.WriteMessage(websocket.BinaryMessage, buf[:n]); err != nil {
			log.Printf("Write to WebSocket error: %s", err)
			return
		}
	}
}

// 从 WebSocket 转发到 TCP
func forwardWebSocketToTCP(conn *websocket.Conn, tcpConn net.Conn) {
	defer tcpConn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Read from WebSocket error: %s", err)
			return
		}

		if _, err := tcpConn.Write(msg); err != nil {
			log.Printf("Write to tcp connection error: %s", err)
			return
		}
	}
}

func handleWebSocket(conn *websocket.Conn, expectedID uuid.UUID) {
	/*defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			log.Println("Disconnect!")
		}
	}(conn)*/

	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Printf("Read message error: %s", err)
		return
	}

	if !verifyUUID(msg, expectedID) {
		log.Println("UUID验证失败！")
		return
	}

	host, targetPort, offset, err := parseTarget(msg)
	if err != nil {
		log.Printf("Error parsing target: %s", err)
		return
	}

	log.Printf("conn: %s:%d", host, targetPort)

	err = conn.WriteMessage(websocket.BinaryMessage, []byte{msg[0], 0})
	if err != nil {
		log.Printf("WriteMessage error: %s", err)
	}

	tcpConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", host, targetPort))
	if err != nil {
		log.Printf("Connection error: %s", err)
		return
	}
	//defer tcpConn.Close()

	tcpConn.Write(msg[offset:])

	go forwardTCPToWebSocket(conn, tcpConn)
	go forwardWebSocketToTCP(conn, tcpConn)

}

func handleRequest(expectedID uuid.UUID) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("Failed to set websocket upgrade: %v", err)
			return
		}
		log.Println("Connecting...")
		handleWebSocket(conn, expectedID)
	}
}

// RunVlessWs todo: 暂不可用
func RunVlessWs(config models.Config) {
	expectedID, err := uuid.Parse(config.UUID)
	if err != nil {
		log.Fatalf("Invalid UUID format: %s", err)
	}

	log.Printf("Listening on port: %d", config.Port)

	http.HandleFunc("/", handleRequest(expectedID))
	err = http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
	if err != nil {
		log.Fatalf("启动http服务失败: %s", err)
	}
}
