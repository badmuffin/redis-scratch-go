package main

import (
	"bufio"
	"os"
	"sync"
	"time"
)

// Aof - Append only file method - Redis records each command in the file as RESP.
type Aof struct {
	file *os.File
	rd   *bufio.Reader // to read from the file
	mu   sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	// os.O_CREATE - create the file if it doesn’t exist.
	// os.O_RDWR - open for both reading and writing.
	// 0666 → Unix file permissions (rw-rw-rw-), so anyone can read/write
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)

	if err != nil {
		return nil, err
	}

	// create a new Aof object
	aof := &Aof{
		file: f,
		rd: bufio.NewReader(f),
	}

	// A goroutine is like a lightweight thread — it runs concurrently with the rest of the program.
	// Background goroutine that syncs the file to disk every second to reduce data loss
	go func(){
		for {
			aof.mu.Lock()
			aof.file.Sync()
			aof.mu.Unlock()
			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}
