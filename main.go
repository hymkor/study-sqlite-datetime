package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/glebarez/go-sqlite/compat"
)

func mains() error {
	println("open")
	conn, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return err
	}
	defer conn.Close()

	println("create table")
	_, err = conn.Exec(`
		CREATE TABLE t_datetime (
			id           INTEGER PRIMARY KEY,
			d_date       DATE,
			d_time       TIME,
			d_datetime   DATETIME)`)
	if err != nil {
		return err
	}

	println("insert table-1")
	rc, err := conn.Exec(`
		INSERT INTO t_datetime
		(d_date, d_time, d_datetime )
		VALUES
		('2025-09-22', '14:30:00', '2025-09-22 14:30:00')`)
	if err != nil {
		return err
	}
	if count, err := rc.RowsAffected(); err != nil {
		return err
	} else {
		println(count, "record(s) updated.")
	}

	println("insert table-2")
	rc, err = conn.Exec(`
		INSERT INTO t_datetime
		(d_date, d_time, d_datetime )
		VALUES
		('2025/09/22', '14:30', '2025/09/22 14:30')`)
	if err != nil {
		return err
	}
	if count, err := rc.RowsAffected(); err != nil {
		return err
	} else {
		println(count, "record(s) updated.")
	}

	println("query table")
	rows, err := conn.Query(`SELECT * from t_datetime`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		r := make([]any, 4)

		err := rows.Scan(&r[0], &r[1], &r[2], &r[3])
		if err != nil {
			return err
		}
		for i, v := range r {
			fmt.Printf("%d: %#v %T\n", i, v, v)
		}
		fmt.Println()
	}
	return nil
}

func main() {
	if err := mains(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
