package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

type filename string

func (f filename) delete() {
	fmt.Printf("Deleting file %s.... \n", f)
	err := os.Remove(string(f))
	if err != nil {
		fmt.Println("error", err)
		os.Exit(1)
	}
}

type runConfiguration struct {
	onlyOnce bool
}

type uploadConfiguration struct {
	awsAccessKeyId     string
	awsSecretAccesKey  string
	awsDefaultRegion   string
	awsBucket          string
	minioUrl           string
	minioEnabled       bool
	minioSsl           bool
	targetFolderPrefix string
}

type mysqlBackupConfiguration struct {
	host     string
	port     string
	database string
	user     string
	password string
	allDbs   bool
}

type cleanConfiguration struct {
	cleanDays int
}

func testConnection(host string, port string) {
	fmt.Printf("Testing connection to %s...\n", host+":"+port)

	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	} else {
		fmt.Println("Connection OK")
	}
	defer conn.Close()
}
