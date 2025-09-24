```main.go
package main

import (
    "database/sql"
    "fmt"
    "os"

    "github.com/mattn/go-sqlite3"
)

func mains() error {
    sqlite3.SQLiteTimestampFormats = sqlite3.SQLiteTimestampFormats[:0]
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
            d_datetime   DATETIME,
            d_text       TEXT)`)
    if err != nil {
        return err
    }

    for _, v := range [][4]string{
        [4]string{`'2025-09-22'`, `'14:30:00'`, `'2025-09-22 14:30:00'`, `'壱'`},
        [4]string{`'2025-09-22'`, `time('14:30:00')`, `'2025-09-22 14:30:00'`, `'弐'`},
        [4]string{`'2025/09/22'`, `'14:30'`, `'2025/09/22 14:30'`, `'参'`},
    } {
        sql := fmt.Sprintf(`
            INSERT INTO t_datetime
            (d_date, d_time, d_datetime, d_text)
            VALUES
            (%s, %s, %s, %s)`, v[0], v[1], v[2], v[3])

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

    fmt.Println("(any)")
    rows, err := conn.Query(`SELECT * from t_datetime`)
    if err != nil {
        return err
    }
    defer rows.Close()
    for rows.Next() {
        r := make([]any, 5)

        err := rows.Scan(&r[0], &r[1], &r[2], &r[3], &r[4])
        if err != nil {
            return err
        }
        for _, v := range r {
            fmt.Printf("%#v as %T\n", v, v)
        }
        fmt.Println()
    }

    fmt.Println("(string)")
    rows, err = conn.Query(`SELECT * from t_datetime`)
    if err != nil {
        return err
    }
    defer rows.Close()
    for rows.Next() {
        r := make([]string, 5)

        err := rows.Scan(&r[0], &r[1], &r[2], &r[3], &r[4])
        if err != nil {
            return err
        }
        for _, v := range r {
            fmt.Println(v)
        }
        fmt.Println()
    }

    fmt.Println("(RawBytes)")
    rows, err = conn.Query(`SELECT * from t_datetime`)
    if err != nil {
        return err
    }
    defer rows.Close()
    for rows.Next() {
        r := make([]sql.RawBytes, 5)

        err := rows.Scan(&r[0], &r[1], &r[2], &r[3], &r[4])
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
```

```./mattn |
open
create table

            INSERT INTO t_datetime
            (d_date, d_time, d_datetime, d_text)
            VALUES
            ('2025-09-22', '14:30:00', '2025-09-22 14:30:00', '壱')
1 record(s) updated.

            INSERT INTO t_datetime
            (d_date, d_time, d_datetime, d_text)
            VALUES
            ('2025-09-22', time('14:30:00'), '2025-09-22 14:30:00', '弐')
1 record(s) updated.

            INSERT INTO t_datetime
            (d_date, d_time, d_datetime, d_text)
            VALUES
            ('2025/09/22', '14:30', '2025/09/22 14:30', '参')
1 record(s) updated.
(any)
1 as int64
time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC) as time.Time
"14:30:00" as string
time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC) as time.Time
"壱" as string

2 as int64
time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC) as time.Time
"14:30:00" as string
time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC) as time.Time
"弐" as string

3 as int64
time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC) as time.Time
"14:30" as string
time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC) as time.Time
"参" as string

(string)
1
0001-01-01T00:00:00Z
14:30:00
0001-01-01T00:00:00Z
壱

2
0001-01-01T00:00:00Z
14:30:00
0001-01-01T00:00:00Z
弐

3
0001-01-01T00:00:00Z
14:30
0001-01-01T00:00:00Z
参

(RawBytes)
1
0001-01-01T00:00:00Z
14:30:00
0001-01-01T00:00:00Z
壱

2
0001-01-01T00:00:00Z
14:30:00
0001-01-01T00:00:00Z
弐

3
0001-01-01T00:00:00Z
14:30
0001-01-01T00:00:00Z
参

```
