package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func backupMysql(mysqlBackupConfiguration mysqlBackupConfiguration) filename {

	testConnection(mysqlBackupConfiguration.host, mysqlBackupConfiguration.port)

	var filesToZip []string

	//List Databases

	if mysqlBackupConfiguration.allDbs {

		db, err := sql.Open("mysql", mysqlBackupConfiguration.user+":"+mysqlBackupConfiguration.password+"@tcp("+mysqlBackupConfiguration.host+":"+mysqlBackupConfiguration.port+")/")
		if err != nil {
			log.Fatal(err)
			os.Exit(0)
		}

		rows, err := db.Query("SHOW DATABASES;")
		if err != nil {
			log.Fatal(err)
			os.Exit(0)
		}

		// Get column names
		columns, err := rows.Columns()
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		// Make a slice for the values
		values := make([]sql.RawBytes, len(columns))

		// rows.Scan wants '[]interface{}' as an argument, so we must copy the
		// references into such a slice
		// See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
		scanArgs := make([]interface{}, len(values))
		for i := range values {
			scanArgs[i] = &values[i]
		}

		for rows.Next() {
			// get RawBytes from data
			err = rows.Scan(scanArgs...)
			if err != nil {
				panic(err.Error()) // proper error handling instead of panic in your app
			}

			// Now do something with the data.
			// Here we just print each column as a string.
			var value string
			for i, col := range values {
				// Here we can check if the value is nil (NULL value)
				if col == nil {
					value = "NULL"
				} else {
					value = string(col)
				}

				if columns[i] == "Database" {
					fmt.Printf("Backupping %s \n", columns[i]+" "+value)
					filesToZip = append(filesToZip, string(singleMysqlBackup(mysqlBackupConfiguration.user, mysqlBackupConfiguration.password, mysqlBackupConfiguration.host, mysqlBackupConfiguration.port, value)))
				}

			}
		}
		if err = rows.Err(); err != nil {
			log.Fatal(err)
			os.Exit(0)
		}
		defer db.Close()

	} else {
		filesToZip = append(filesToZip, string(singleMysqlBackup(mysqlBackupConfiguration.user, mysqlBackupConfiguration.password, mysqlBackupConfiguration.host, mysqlBackupConfiguration.port, mysqlBackupConfiguration.database)))
	}

	currentTime := time.Now().Local()
	timestamp := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)
	//zip file before uploading

	fileDumpZip := "tempdump/" + currentTime.Format("2006-01-02") + "-" + timestamp + ".zip"

	if err := zipFiles(fileDumpZip, filesToZip); err != nil {
		panic(err)
	}

	for _, file := range filesToZip {
		filename(file).delete()
	}

	return filename(fileDumpZip)

}

func singleMysqlBackup(user string, password string, host string, port string, database string) filename {

	fmt.Printf("Executing mysqldump on %s database %s...\n", host, database)
	cmd := exec.Command("mysqldump",
		"--complete-insert",
		"--skip_add_locks",
		"--skip-lock-tables",
		"--max_allowed_packet=1500M",
		"-P"+port,
		"-h"+host,
		"-u"+user,
		"-p"+password,
		database)

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

	fileDump := "tempdump/" + currentTime.Format("2006-01-02") + "-" + timestamp + "_" + database + ".sql"

	err = ioutil.WriteFile(fileDump, bytes, 0644)
	if err != nil {
		panic(err)
	}

	return filename(fileDump)

}
