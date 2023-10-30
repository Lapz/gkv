package main

import (
	"fmt"
	"gkv/message"
	"net"
)

const PORT int = 8080

type KVStore = map[string]Entry

type Entry struct {
	Value []byte
}

var store = make(map[string]Entry)

func handleConnection(conn net.Conn) error {
	var rawBytes []byte = make([]byte, 1024)

	_, err := conn.Read(rawBytes)

	if err != nil {
		fmt.Println("Error reading a msg", err)
	}

	defer conn.Close() // always close the connection

	msg, err := message.New(rawBytes)

	if err != nil {
		return err
	}

	switch msg.GetKind() {
	case message.Get:
		{
			key, err := msg.GetArgument(1)

			if err != nil {
				return err
			}

			entry, present := store[string(key)]

			if present {
				conn.Write(entry.Value)
			} else {
				conn.Write([]byte(fmt.Sprintf("Unknown key %s", key)))
			}

		}
	case message.Set:
		{
			key, err := msg.GetArgument(1)

			if err != nil {
				return err
			}

			value, err := msg.GetArgument(2)

			if err != nil {
				return err
			}

			e := Entry{
				value,
			}

			store[string(key)] = e

			conn.Write(value)
		}
	}

	return nil

}

func main() {

	fmt.Println("Starting a server on port 8080")
	ln, err := net.Listen("tcp", ":8080")

	if err != nil {
		fmt.Println("A server already exists on port :8080")
		return
	}

	for {
		conn, err := ln.Accept()

		println("New connection accepted")

		if err != nil {
			// handle error
		}
		go handleConnection(conn)
	}

}
