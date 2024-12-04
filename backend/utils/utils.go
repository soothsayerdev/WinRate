package utils

import (
	"fmt"
	"database/sql"
	"log"
	"os"


	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func ConnectDB() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"))

	var err error 
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados: ", err)
	}
	//defer db.Close()

	// verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Erro ao verificar a conexão: ", err)
	}

	fmt.Println("Conectado ao banco de dados com sucesso!")
}