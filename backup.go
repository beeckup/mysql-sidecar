package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"bytes"
	"github.com/robfig/cron"
	"os/exec"
)

const (
	CONN_PORT = "3333"
	CONN_TYPE = "tcp"
)

func main() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, ":"+CONN_PORT)
	if err != nil {
		fmt.Println("Cannot bind! ", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Useless listen " + CONN_PORT)

	c := cron.New()
	c.AddFunc(os.Getenv("SCHEDULE"), func() { fmt.Println("Cron schedule: "+os.Getenv("SCHEDULE")) })
	c.AddFunc(os.Getenv("SCHEDULE"), func() { exe_cmd("./cron_script.sh") })
	c.Start()

	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}

		//logs an incoming message
		fmt.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())

		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)
	// Read the incoming connection into the buffer.
	reqLen, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	// Builds the message.
	message := "Hi, I received your message! It was "
	message += strconv.Itoa(reqLen)
	message += " bytes long and that's what it said: \""
	n := bytes.Index(buf, []byte{0})
	message += string(buf[:n-1])
	message += "\" ! Honestly I have no clue about what to do with your messages, so Bye Bye!\n"

	// Write the message in the connection channel.
	conn.Write([]byte(message));
	// Close the connection when you're done with it.
	conn.Close()
}

func exe_cmd(cmd string) {

	out, err := exec.Command("sh","-c",cmd).Output()
	if err != nil {
		fmt.Printf("%s", err)
	}
	fmt.Printf("%s", out)
}