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
	// "time"
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

	// loop here??

	for {

		Conn, err := net.DialUDP("udp", LocalIpAndPort, ServerIpAndPort)
		CheckError(err)

		//defer Conn.Close()

		var guess uint32 = ComputeGuess(minimum, maximum)

		fmt.Println("Guess is: %d", guess)
		buf, err := Marshall(guess)
		_, err = Conn.Write(buf)

		x := Conn.Close()
		fmt.Println("the result of Conn.Close() is: %v", x)

		// now read...
		ServerConn, err := net.ListenUDP("udp", LocalIpAndPort)
		CheckError(err)

		//defer ServerConn.Close()

		buffer := make([]byte, 1024)

		n, addr, err := ServerConn.ReadFromUDP(buffer)
		fmt.Println("Received ", string(buffer[0:n]), " from ", addr, " where the value of n is: ", n)

		y := ServerConn.Close()
		fmt.Println("the result of ServerConn.Close() is: %v", y)

		switch n {
		case 4:
			maximum = guess
			//guess = ComputeGuess(minimum, maximum)
			// return to loop
		case 3:
			minimum = guess
			//guess = ComputeGuess(minimum, maximum)
			// return to loop
		default:
			// should be done...
			fmt.Println("Should be the result.....")
			return
		}
	}

}

// guess is the number we are going to send
// returns byte slice/array with the guess converted to sendable state
// similar to serialization...
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
