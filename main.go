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

	for _, v := range [][3]string{
		[3]string{`'2025-09-22'`, `'14:30:00'`, `'2025-09-22 14:30:00'`},
		[3]string{`'2025-09-22'`, `time('14:30:00')`, `'2025-09-22 14:30:00'`},
		[3]string{`'2025/09/22'`, `'14:30'`, `'2025/09/22 14:30'`},
	} {
		sql := fmt.Sprintf(`
			INSERT INTO t_datetime
			(d_date, d_time, d_datetime )
			VALUES
			(%s, %s, %s)`, v[0], v[1], v[2])

		fmt.Println(sql)
		rc, err := conn.Exec(sql)
		if err != nil {
			return err
		}
		if count, err := rc.RowsAffected(); err != nil {
			return err
		} else {
			println(count, "record(s) updated.")
		}
	}

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
		for _, v := range r {
			fmt.Printf("%#v as %T\n", v, v)
		}
		fmt.Println()
	}
	rows, err = conn.Query(`SELECT * from t_datetime`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		r := make([]sql.RawBytes, 4)

		err := rows.Scan(&r[0], &r[1], &r[2], &r[3])
		if err != nil {
			return err
		}
		for _, v := range r {
			fmt.Println(string(v))
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
