package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	// setup tcp listener for any client to communicate with this
	ln, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Server is running on port 6379...")

	aof, err := NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer aof.CloseFile()

	aof.Read(func(value Value) {
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid Command ", command)
			return
		}

		handler(args)
	})

	// start receiving requests
	conn, err := ln.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	// close connection once finished
	defer conn.Close()

	for {
		// create a new RESP reader for the client's TCP connection
		resp := NewResp(conn)
		value, err := resp.Read()

		if err != nil {
			fmt.Println(err)
			return
		}

		if value.typ != "array" {
			fmt.Println("Invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		// first element - PING, SET, etc
		command := strings.ToUpper(value.array[0].bulk)
		// rest of the element
		args := value.array[1:]

		// create a RESP writer bound to client's connection
		writer := NewWriter(conn)

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		if command == "SET" || command == "HSET" {
			aof.Write(value)
		}

		result := handler(args)
		writer.Write(result)
	}
}
