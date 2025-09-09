package main

// key : value :: map[string] : func([]Value) Value {...}
// where key - command name like PING, SET, GET
// and value - function
var Handlers = map[string]func([]Value) Value{
	"PING": ping,
}

func ping(args []Value) Value {
	if len(args) > 0 {
		return Value{typ: "string", str: args[0].bulk}
	}
	return Value{typ: "string", str: "PONG"}

}
