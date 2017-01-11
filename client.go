/*
Implements the solution to assignment 1 for UBC CS 416 2016 W2.

Usage:
$ go run client.go [local UDP ip:port] [server UDP ip:port]

Example:
$ go run client.go 127.0.0.1:2020 127.0.0.1:7070

*/

package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"net"
	"os"
	"time"
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

	// TODO

	LocalIpAndPort, err := net.ResolveUDPAddr("udp", local_ip_port)
	CheckError(err)

	ServerIpAndPort, err := net.ResolveUDPAddr("udp", remote_ip_port)
	CheckError(err)

	var minimum uint32 = 0
	var maximum uint32 = math.MaxUint32

	for {
		var guess uint32 = ComputeGuess(minimum, maximum)
		buf, err := Marshall(guess)
		CheckError(err)

		// send guess, if timeout, resend
		// else compute next guess or print result
		isTimeout := true

		for isTimeout {

			Conn, err := net.DialUDP("udp", LocalIpAndPort, ServerIpAndPort)
			CheckError(err)
			_, err = Conn.Write(buf)
			Conn.Close()

			// get response
			ServerConn, err := net.ListenUDP("udp", LocalIpAndPort)
			CheckError(err)
			ServerConn.SetReadDeadline(time.Now().Add(time.Second * 2))

			buffer := make([]byte, 1024)
			n, _, err := ServerConn.ReadFromUDP(buffer)

			// timed out
			if err, ok := err.(net.Error); ok && err.Timeout() {
				ServerConn.Close()

			} else if err != nil {
				// TODO ????
				CheckError(err)

			} else {
				isTimeout = false

				ServerConn.Close()

				switch n {
				case 4:
					maximum = guess
				case 3:
					minimum = guess
				default:
					fmt.Println(string(buffer[0:n]))
					return
				}
			}
		}
	}
}

func Marshall(guess uint32) ([]byte, error) {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	err := enc.Encode(guess)
	return network.Bytes(), err
}

// TODO reference???   https://varshneyabhi.wordpress.com/2014/12/23/simple-udp-clientserver-in-golang/
func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

func ComputeGuess(min uint32, max uint32) uint32 {
	var result uint32 = min + (max-min)/2
	return result
}
