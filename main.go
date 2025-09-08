package main

import (
	"fmt"
	"net"
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

	for {
		resp := NewResp(conn)
		value, err := resp.Read()

		if err != nil {
			fmt.Println(err)
			return
		}

		_ = value

		writer := NewWriter(conn)
		writer.Write(Value{typ: "string", str: "OK"})
	}
}
