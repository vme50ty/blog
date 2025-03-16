package main

import (
	"go-web/controller"
	"go-web/model"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/context"

	_ "github.com/go-sql-driver/mysql"

	_ "net/http/pprof"
)

func main() {
	go func() {
		if err := http.ListenAndServe("0.0.0.0:6060", nil); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	// Setup DB
	log.Println("DB Init ...")
	db := model.ConnectToDB()
	defer db.Close()
	model.SetDB(db)

	router := controller.Startup()

	http.ListenAndServe("0.0.0.0:8888", context.ClearHandler(router))
}
