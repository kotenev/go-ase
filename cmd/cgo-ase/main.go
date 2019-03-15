package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strings"

	_ "github.com/SAP/go-ase/cgo"
	"github.com/SAP/go-ase/libase"
	"github.com/bgentry/speakeasy"
)

var (
	fHost         = flag.String("H", "localhost", "database hostname")
	fPort         = flag.String("P", "4901", "database sql port")
	fUser         = flag.String("u", "sa", "database user name")
	fPass         = flag.String("p", "", "database user password")
	fUserstorekey = flag.String("k", "", "userstorekey")
	fDatabase     = flag.String("D", "", "database")
)

func exec(db *sql.DB, q string) error {
	result, err := db.Exec(q)
	if err != nil {
		return fmt.Errorf("Executing the statement failed: %v", err)
	}

	return processResult(result)
}

func query(db *sql.DB, q string) error {
	rows, err := db.Query(q)
	if err != nil {
		return fmt.Errorf("Query failed: %v", err)
	}
	defer rows.Close()

	return processRows(rows)
}

func subcmd(db *sql.DB, part string) error {
	partS := strings.Split(part, " ")
	cmd := partS[0]
	q := strings.Join(partS[1:], " ")

	switch cmd {
	case "exec":
		err := exec(db, q)
		if err != nil {
			return fmt.Errorf("Exec errored: %v", err)
		}
	case "query":
		err := query(db, q)
		if err != nil {
			return fmt.Errorf("Query errored: %v", err)
		}
	default:
		log.Printf("Unknown command: %s", cmd)
	}

	return nil
}

func main() {
	flag.Parse()
	pass := *fPass
	var err error
	if len(pass) == 0 && len(*fUserstorekey) == 0 {
		pass, err = speakeasy.Ask("Please enter the password of user " + *fUser + ": ")
		if err != nil {
			log.Println(err)
			return
		}
	}

	dsn := libase.DsnInfo{
		Host:         *fHost,
		Port:         *fPort,
		Username:     *fUser,
		Password:     pass,
		Userstorekey: *fUserstorekey,
		Database:     *fDatabase,
	}

	db, err := sql.Open("ase", dsn.AsSimple())
	if err != nil {
		log.Printf("Error opening database connection: %v", err)
		return
	}
	defer db.Close()

	// test the database connection
	err = db.Ping()
	if err != nil {
		log.Printf("Pinging the server failed: %v", err)
		return
	}

	if len(flag.Args()) == 0 {
		return
	}

	subcmds := strings.Split(strings.Join(flag.Args(), " "), "--")
	for _, s := range subcmds {
		err = subcmd(db, strings.TrimSpace(s))
		if err != nil {
			log.Printf("Execution of '%s' resulted in error: %v", s, err)
			return
		}
	}
}
