package db

import (
	"github.com/aberic/gnomon"
)

const (
	GDBUrl  = "DB_URL"  // DBUrl 数据库 URL
	GDBName = "DB_NAME" // DBName 数据库名称
	GDBUser = "DB_USER" // DBUser 数据库用户名
	GDBPass = "DB_PASS" // DBPass 数据库用户密码
)

var (
	// sql 数据库操作对象
	sql *gnomon.SQLCommon
)

func init() {
	sql = gnomon.SQL()
	connect()
}

func connect() {
	dbURL := gnomon.Env().GetD(GDBUrl, "127.0.0.1:3306")
	dbUser := gnomon.Env().GetD(GDBUser, "root")
	dbPass := gnomon.Env().GetD(GDBPass, "root")
	dbName := gnomon.Env().GetD(GDBName, "baas")
	if err := sql.Connect(dbURL, dbUser, dbPass, dbName, false, 5, 20); err != nil {
		connect()
	}
	// 全局禁用表名复数
	sql.DB.SingularTable(true) // 如果设置为true,`User`的默认表名为`user`,使用`TableName`设置的表名不受影响
}
