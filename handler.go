package main

import "sync"

// key : value :: map[string] : func([]Value) Value {...}
// where key - command name like PING, SET, GET
// and value - function
var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
}

func ping(args []Value) Value {
	if len(args) > 0 {
		return Value{typ: "string", str: args[0].bulk}
	}
	return Value{typ: "string", str: "PONG"}

}

// hash map for storing key value
var SETs = map[string]string{}

// a mutex for concurrency - this prevents race conditions
// when multiple client send requests at the same time
var SETsMu = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERROR wrong number of arguments for SET command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERROR wrong number of arguments for GET command"}
	}

	key := args[0].bulk

	SETsMu.Lock()
	value, ok := SETs[key]
	SETsMu.Unlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}
