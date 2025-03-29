package main

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	// 导入了一个包，但不直接使用它的任何函数或类型，只是为了触发它的初始化逻辑。
	/* MySQL 驱动的特殊之处：
	   MySQL 驱动并不是你直接调用它的某些方法，
	   而是：它在 init() 里向 database/sql 注册了驱动名 mysql */)

var DB *sql.DB
var MySQLServer string
var MySQLPasswd string

type MYSQL_CON struct{}

func initDB() bool {
	// MySQLServer = "rockage.net"
	MySQLServer = "localhost"
	MySQLPasswd = ""

	path := "root:" + MySQLPasswd + "@tcp(" + MySQLServer + ":3306)/stoneren-bbs?charset=utf8mb4"
	DB, _ = sql.Open("mysql", path)

	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)

	//验证连接
	if err := DB.Ping(); err != nil {
		return false
	}
	return true
}

func (s MYSQL_CON) Exec(SQL string) (string, error) {
	var ins string = ""
	defer DB.Close()
	if initDB() {
		ret, _ := DB.Exec(SQL)
		defer DB.Close()
		if ret != nil {
			if LastInsertId, err := ret.LastInsertId(); nil == err {
				ins = strconv.FormatInt(LastInsertId, 10)
			}
		}
	} else {
		return "", fmt.Errorf("mysql connection error") // 直接返回错误
	}
	return ins, nil
}

func QueryDB(query string) ([]map[string]interface{}, error) {

	// 运行 SQL 查询

	var results []map[string]interface{}
	if initDB() {
		rows, err := DB.Query(query)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		// **获取列名**
		columns, err := rows.Columns() //columns获取了查询返回的列名
		if err != nil {
			return nil, err
		}
		// **存储查询结果**

		for rows.Next() { // for -> 直到rows为EOF
			// **创建一个切片，用于存放字段值**
			values := make([]interface{}, len(columns))
			valuePtrs := make([]interface{}, len(columns))

			for i := range values {
				valuePtrs[i] = &values[i] // 指向值的指针
			}

			// **扫描数据**
			if err := rows.Scan(valuePtrs...); err != nil {
				return nil, err
			}

			// **创建 map 并存储数据**
			rowMap := make(map[string]interface{})
			for i, colName := range columns {
				val := values[i]

				// **处理 NULL 值**
				if b, ok := val.([]byte); ok {
					rowMap[colName] = string(b) // 转换 []byte 到 string
				} else {
					rowMap[colName] = val
				}
			}

			results = append(results, rowMap)
		}
	}

	return results, nil
}
