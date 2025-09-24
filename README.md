Pure Go の SQLite3 ドライバでの日時系カラムの取り扱い
=====================================================

検証内容
--------

SQLite3 で `DATE`型、`TIME`型、`DATETIME`型は実質文字列カラムとされているが、Pure Go の SQLite3 ドライバ [glebarez/go-sqlite] で any 型で受けとった時、time.Time 型に格納されることが分かっている。が、イレギュラーな表現だった場合、どうなることだろうか？

また `string`, `RawBytes` に格納した場合は、どうなるか？

この検証結果により、SELECT 文で得た値をありのままデータベースに再格納する方法を検討する。

検証プログラム
--------------

```main.go
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

検証プログラムの実行結果
------------------------

```./study-sqlite-datetime |
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
time.Date(2025, time.September, 22, 0, 0, 0, 0, time.UTC) as time.Time
"14:30:00" as string
time.Date(2025, time.September, 22, 14, 30, 0, 0, time.UTC) as time.Time
"壱" as string

2 as int64
time.Date(2025, time.September, 22, 0, 0, 0, 0, time.UTC) as time.Time
"14:30:00" as string
time.Date(2025, time.September, 22, 14, 30, 0, 0, time.UTC) as time.Time
"弐" as string

3 as int64
"2025/09/22" as string
"14:30" as string
"2025/09/22 14:30" as string
"参" as string

(string)
1
2025-09-22T00:00:00Z
14:30:00
2025-09-22T14:30:00Z
壱

2
2025-09-22T00:00:00Z
14:30:00
2025-09-22T14:30:00Z
弐

3
2025/09/22
14:30
2025/09/22 14:30
参

(RawBytes)
1
2025-09-22T00:00:00Z
14:30:00
2025-09-22T14:30:00Z
壱

2
2025-09-22T00:00:00Z
14:30:00
2025-09-22T14:30:00Z
弐

3
2025/09/22
14:30
2025/09/22 14:30
参

```

結果
---

- `any` 型変数に格納しようとした場合:
    - `DATE`型, `DATETIME`型:
        - 正規の書式で登録された値は `time.Time` 型に変換される
        - イレギュラーな値は `string` 型となる
    - `TIME`型:
        - 常に`string` 型に変換される
- `string` , `RawBytes` 型変数に格納させた場合:
    - `DATE`型, `DATETIME`型:
        - 正規の書式で登録された値は `2006-01-02T15:04:05Z` 形式となる
        - イレギュラーな値は、ありのまま格納される
    - `TIME`型:
        - 格納された時のまま

課題
----

どの範囲までが正規の書式であるかを確認しなくてはいけない

- [mattn/go-sqlite3] の場合 `SQLiteTimestampFormats` で定義されている。  
    → `sqlite3.SQLiteTimestampFormats = sqlite3.SQLiteTimestampFormats[:0]`
       とかにすると常に文字列で扱ってくれそう
- [glebarez/go-sqlite] には該当のものは、なさそう

[mattn/go-sqlite3]: https://github.com/mattn/go-sqlite3
[glebarez/go-sqlite]: https://github.com/glebarez/go-sqlite
