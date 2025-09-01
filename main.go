package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	// setup tcp listener for any client to communicate with this
	ln, err := net.Listen("tcp", ":6379")

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Server is running on port 6379...")

	// start receiving requests
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	// close connection once finished
	defer conn.Close()

	// create an infinite loop and receive commands from clients and respond to them
	for {
		// this is to store the data recieved from the client (temporarily)
		buff := make([]byte, 1024) // declare this outside the loop and reuse it

		//read message from client
		_, err := conn.Read(buff)
		if err != nil {
			if err != io.EOF {
				break
			}
			fmt.Println("error reading from client: ", err.Error())
			os.Exit(1) // its better to just return instead of killing whole server
		}

		// ignore request and send back a respond
		conn.Write([]byte("+OK\r\n"))
	}
}
