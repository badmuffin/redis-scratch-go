package main

import (
	"sync"
)

// key : value :: map[string] : func([]Value) Value {...}
// where key - command name like PING, SET, GET
// and value - function
var Handlers = map[string]func([]Value) Value{
	"PING":    ping,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetall,
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

// SET key value
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

	SETsMu.RLock() // many clients can read concurrently
	value, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

// HSET and HGET
// {
// 	"users": {
// 		"u1": "Gandalf",
// 		"u2": "Bilbo",
// 	},
// 	"posts": {
// 		"p1": "Hello World",
// 		"p2": "Welcome to my blog",
// 	},
// }

// map of maps
var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

// HSET hash key value
func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERROR wrong number of arguments for HSET command"}
	}

	hash := args[0].bulk  // hash name, e.g "users"
	key := args[1].bulk   // field, e.g "u1"
	value := args[2].bulk // value, e.g "Gandalf"

	HSETsMu.Lock()
	_, ok := HSETs[hash] // check if hash exists

	if !ok {
		// create a empty map
		HSETs[hash] = map[string]string{}
	}

	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERROR wrong number of arguments for HGET command"}
	}

	hash := args[0].bulk // hash name, e.g "users"
	key := args[1].bulk  // field, e.g "u1"

	HSETsMu.RLock() // many clients can read concurrently
	value, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}
	return Value{typ: "bulk", bulk: value}
}

func hgetall(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERROR wrong number of arguments for HGETALL command"}
	}

	hash := args[0].bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash]
	HSETsMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	values := []Value{}
	for k, v := range value {
		values = append(values, Value{typ: "bulk", bulk: k})
		values = append(values, Value{typ: "bulk", bulk: v})
	}

	return Value{typ: "array", array: values}
}
