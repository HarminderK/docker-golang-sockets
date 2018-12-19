package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

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
		return
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
