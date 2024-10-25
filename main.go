package main

import (
	"go-web/controller"
	"go-web/model"
	"log"
	"net/http"

	"github.com/gorilla/context"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Setup DB
	log.Println("DB Init ...")
	db := model.ConnectToDB()
	defer db.Close()
	model.SetDB(db)

	controller.Startup()

	http.ListenAndServe(":8888", context.ClearHandler(http.DefaultServeMux))
}
