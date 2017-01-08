/*
Implements the solution to assignment 1 for UBC CS 416 2016 W2.

Usage:
$ go run client.go [local UDP ip:port] [server UDP ip:port]

Example:
$ go run client.go 127.0.0.1:2020 127.0.0.1:7070

*/

package main

import (
	"encoding/gob"
	"fmt"
	"os"
	"bytes"
)

// Main workhorse method.
func main() {
	args := os.Args[1:]

	// Missing command line args.
	if len(args) != 2 {
		fmt.Println("Usage: client.go [local UDP ip:port] [server UDP ip:port]")
		return
	}

	// Extract the command line args.
	local_ip_port := args[0]
	remote_ip_port := args[1]

	// for compilation...
	fmt.Println("local_ip_port : %d", local_ip_port)
	fmt.Println("remote_ip_port : %d", remote_ip_port)

	// TODO
}

func Marshall(guess uint32) ([]byte, error) {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(guess)
	return network.Bytes(), err
}
