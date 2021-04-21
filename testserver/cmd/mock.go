package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
)

func setupMockServer() {
	list, err := net.Listen("tcp", ":8081")
	if err != nil {
		panic(err)
	}
	conn, err := list.Accept()
	if err != nil {
		panic(err)
	}
	fmt.Println("Mock server is serving at 8081")
	rConn := bufio.NewReader(conn)
	go func() {
		for {
			_, _, err := rConn.ReadLine()
			if errors.Is(err, io.EOF) {
				return
			}
			_, err = conn.Write([]byte("something\n"))
			if err != nil {
				panic(err)
			}
		}
	}()
}

func mockSocket(readChan chan []byte, writeChan chan []byte) {
	conn, err := net.Dial("tcp", ":8081")
	if err != nil {
		panic(err)
	}
	go writeLoop(conn, writeChan)
	rChan := bufio.NewReader(conn)
	go readLoop(rChan, readChan)
	fmt.Println("Dialing to 8081")

}

func writeLoop(conn net.Conn, writeChan chan []byte) {
	for {
		msg, ok := <-writeChan
		if !ok {
			return
		}
		conn.Write(msg)
	}
}
func readLoop(conn *bufio.Reader, readChan chan []byte) {
	for {
		ln, _, err := conn.ReadLine()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}
			panic(err)
		}
		readChan <- ln
	}
}
