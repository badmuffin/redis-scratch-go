package main

import (
	"bufio"
	"errors"
	"io"
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
		rd:   bufio.NewReader(f),
	}

	// A goroutine is like a lightweight thread — it runs concurrently with the rest of the program.
	// Background goroutine that syncs the file to disk every second to reduce data loss
	go func() {
		for {
			aof.mu.Lock()
			aof.file.Sync()
			aof.mu.Unlock()
			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

func (aof *Aof) CloseFile() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	return aof.file.Close()
}

func (aof *Aof) Write(value Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	// write the command to the file in the same RESP format that we receive
	_, err := aof.file.Write(value.Marshal())
	if err != nil {
		return err
	}

	return nil
}

// Read - read commands from the AOF file and for each command it finds, it calls your
// callback function with that command.
func (aof *Aof) Read(callback func(value Value)) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	resp := NewResp(aof.file)

	for {
		value, err := resp.Read()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		callback(value)
	}

	return nil
}
