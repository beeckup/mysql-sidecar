package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"time"
)

func backupMysql(mysqlBackupConfiguration mysqlBackupConfiguration) filename {

	testConnection(mysqlBackupConfiguration.host, mysqlBackupConfiguration.port)
	fmt.Printf("Executing mysqldump on %s database %s...\n", mysqlBackupConfiguration.host, mysqlBackupConfiguration.database)
	cmd := exec.Command("mysqldump",
		"--complete-insert",
		"--skip_add_locks",
		"--skip-lock-tables",
		"--max_allowed_packet=1500M",
		"-P"+mysqlBackupConfiguration.port,
		"-h"+mysqlBackupConfiguration.host,
		"-u"+mysqlBackupConfiguration.user,
		"-p"+mysqlBackupConfiguration.password,
		mysqlBackupConfiguration.database)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Fatal(err)
	}

	currentTime := time.Now().Local()
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

	fileDump := "tempdump/" + currentTime.Format("2006-01-02") + "-" + timestamp + "_dump.sql"

	err = ioutil.WriteFile(fileDump, bytes, 0644)
	if err != nil {
		panic(err)
	}

	//zip file before uploading
	filesToZip := []string{fileDump}
	fileDumpZip := "tempdump/" + currentTime.Format("2006-01-02") + "-" + timestamp + "_dump.zip"

	if err := zipFiles(fileDumpZip, filesToZip); err != nil {
		panic(err)
	}

	filename(fileDump).delete()

	return filename(fileDumpZip)

}
