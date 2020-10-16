package mysqllib

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	username := viper.GetString("mysql.username")
	password := viper.GetString("mysql.password")
	host := viper.GetString("mysql.host")
	port := viper.GetString("mysql.port")
	database := viper.GetString("mysql.database")

	path := strings.Join([]string{username, ":", password, "@tcp(", host, ":", port, ")/", database, "?charset=utf8&loc=Asia%2FShanghai&parseTime=true"}, "")
	fmt.Println(path)

	DB, _ = sql.Open("mysql", path)

	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)

	if err := DB.Ping(); err != nil {
		fmt.Println("open database fail")
		return
	}
	fmt.Println("connect success")
}

func GetMysqlConn() (d *sql.DB) {
	return DB
}
