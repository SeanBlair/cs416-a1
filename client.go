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
	// "strconv"
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
		fmt.Println("Guess is: ", guess)
		buf, err := Marshall(guess)
		CheckError(err)

		// send guess with timeout feature

		isTimeout := true

		for isTimeout {

			fmt.Println("Guess is: ", guess)

			Conn, err := net.DialUDP("udp", LocalIpAndPort, ServerIpAndPort)
			CheckError(err)
			_, err = Conn.Write(buf)

			x := Conn.Close()
			fmt.Println("the result of Conn.Close() is: ", x)

			// now read...
			// TODO: set timeout...
			ServerConn, err := net.ListenUDP("udp", LocalIpAndPort)
			CheckError(err)

			fmt.Println("listening.....")

			ServerConn.SetReadDeadline(time.Now().Add(time.Second * 2))

			buffer := make([]byte, 1024)

			n, addr, err := ServerConn.ReadFromUDP(buffer)

			if err, ok := err.(net.Error); ok && err.Timeout() {
				fmt.Println("============================================================")
				// This was a timeout
				// //isTimeout = true
				// fmt.Println("this was a timeout!!!!!!!!")
				z := ServerConn.Close()
				fmt.Println("the result of ServerConn.Close() is: ", z)

			} else if err != nil {
				// This was an error, but not a timeout
				// TODO
				fmt.Println("got into the error check of the timeout...")

			} else {
				isTimeout = false

				fmt.Println("made it past ServerConn.ReadFromUDP(buffer)")
				CheckError(err)
				fmt.Println("Received ", string(buffer[0:n]), " from ", addr, " where the value of n is: ", n)

				y := ServerConn.Close()
				fmt.Println("the result of ServerConn.Close() is: ", y)

				switch n {
				case 0:
					// do nothing???
					fmt.Println("was the zero case......")
				case 4:
					maximum = guess
				case 3:
					minimum = guess
				default:
					fmt.Println("Should be the result.....")
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
