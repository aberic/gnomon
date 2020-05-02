package db

import (
	"github.com/aberic/gnomon"
	"github.com/aberic/gnomon/db"
)

const (
	GDBUrl  = "DB_URL"  // DBUrl 数据库 URL
	GDBName = "DB_NAME" // DBName 数据库名称
	GDBUser = "DB_USER" // DBUser 数据库用户名
	GDBPass = "DB_PASS" // DBPass 数据库用户密码
)

var (
	// sql 数据库操作对象
	sql *db.MySQL
)

func init() {
	sql = &db.MySQL{}
	connect()
}

func connect() {
	dbURL := gnomon.EnvGetD(GDBUrl, "127.0.0.1:3306")
	dbUser := gnomon.EnvGetD(GDBUser, "root")
	dbPass := gnomon.EnvGetD(GDBPass, "root")
	dbName := gnomon.EnvGetD(GDBName, "baas")
	if err := sql.Connect(dbURL, dbUser, dbPass, dbName, false, 5, 20); err != nil {
		connect()
	}
	// 全局禁用表名复数
	sql.DB.SingularTable(true) // 如果设置为true,`User`的默认表名为`user`,使用`TableName`设置的表名不受影响
}
