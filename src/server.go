package main

import (
	"fmt"
	"net"
)

type chat struct {
	clients map[string]net.Conn
	users   []string
	// Use channels for more parallel communication
	// broadcast chan []byte
	// history   chan []byte
}

type chatRooms struct {
	// Chat ID -> chat
	chats     map[string]chat
	handShake map[string]bool
}

// Message struct
type Message struct {
	// Name, Chatid, AuthToken, Text string
	Name, Chatid, Text string
}

func (rooms *chatRooms) handleConn(conn net.Conn) {
	defer conn.Close()
	message := make([]byte, 4096)
	length, err := conn.Read(message)

	if err != nil {
		return
	}
	fmt.Println("SENT: " + string(message))

	var m Message
	err := json.Unmarshal(message, &m)

	if err != nil {
		fmt.println(err)
	}

	if err != nil && length > 0 {

		i, ok := rooms.chats[m.Chatid]

		if ok {
			for _, user := range i.users {
				if user == m.Name {
					i.clients[user] = conn
					continue
				}
				recepConn, ok := i.clients[user]
				if ok {
					_, err := recepConn.Write([]byte(m.Text))
					// Clean up any closed connections
					if err != nil {
						fmt.Println("User: " + user + " is offline.")
					}
				}
			}

		} else {
			return
		}

	} else {
		conn.Write([]byte("World!"))
	}

}
func (rooms *chatRooms) addChatRoom(chatid string, users []string) {

	chatroom := chat{
		clients: make(map[string]net.Conn),
		users:   users,
	}

	rooms.chats[chatid] = chatroom

}

func startServer(port string) {

	fmt.Println("Checking for the addr")
	portAddr, err := net.ResolveTCPAddr("tcp4", "localhost:"+port)
	checkError(err)

	fmt.Println("Listening at the addr")
	listener, err := net.ListenTCP("tcp", portAddr)
	checkError(err)

	rooms := chatRooms{
		chats: make(map[string]chat),
	}

	for {

		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go rooms.handleConn(conn)
	}
}
