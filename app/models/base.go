package models

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"go-app/config"
	"log"
	"os"
)

var DbConnection *sql.DB

// init DBの初期化を行う
func init() {
	// AWSLambda実行の場合はスキップ
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") != "" {
		return
	}

	var err error
	DbConnection, err = sql.Open(config.Config.SqlDriver, config.Config.DbName)
	if err != nil {
		log.Fatalln(err)
	}

	initProgram()
}
