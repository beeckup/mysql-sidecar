package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
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
	folder    string
	cleanDays int64
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

func getDifferenceDays(t1 *time.Time, t2 *time.Time) int64 {

	t2Tmp := *t2
	diff := t1.Sub(t2Tmp)
	return int64(diff.Hours() / 24)

}

func getDifferenceMinutes(t1 *time.Time, t2 *time.Time) int64 {

	t2Tmp := *t2
	diff := t1.Sub(t2Tmp)
	return int64(diff.Minutes())

}
