package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func handleConn(conn net.Conn) {
	defer conn.Close()
	message := make([]byte, 4096)
	length, err := conn.Read(message)

	if err != nil {
		conn.Close()
		return
	}

	if length > 0 {
		fmt.Println("SENT: " + string(message))
	}

	conn.Write([]byte("World"))

}

func startServer(port string) {

	fmt.Println("Checking for the addr")
	portAddr, err := net.ResolveTCPAddr("tcp4", "localhost:"+port)
	checkError(err)

	fmt.Println("listening at the addr")
	listener, err := net.ListenTCP("tcp", portAddr)
	checkError(err)

	for {

		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go handleConn(conn)

	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Ran into an error: " + err.Error())
		os.Exit(1)
	}
}

func test(port string) {
	fmt.Println("Starting client...")
	connection, error := net.Dial("tcp", "localhost:"+port)
	if error != nil {
		fmt.Println(error)
	}

	for {
		// Send message
		message := "Hello"
		connection.Write([]byte(strings.TrimRight(message, "\n")))

		// receive message
		recieve := make([]byte, 4096)
		length, err := connection.Read(recieve)
		if err != nil {
			connection.Close()
			break
		}
		if length > 0 {
			fmt.Println("RECEIVED: " + string(recieve))
		}
	}

}

func main() {
	args := os.Args

	if len(args) >= 2 {
		fmt.Println("Listening on TCP Port: ", args[1])

		go startServer(args[1])

		for {
			test(args[1])
		}
	} else {
		fmt.Println("Pass a port as parameter. Try again.")
	}

}
