package main

import (
	"go-web/model"
	"log"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	// password:SFtBfmxbCPY8EMt3
	// password: PUgAUwhg35ANx9Le
	log.Println("DB Init ...")
	db := model.ConnectToDB()
	defer db.Close()
	// 替代删表建表操作
	db.AutoMigrate(&model.User{}, &model.Post{}, &model.Comment{})
	// 	model.SetDB(db)

	// 	db.DropTableIfExists(model.User{}, model.Post{}, "follower")
	// 	db.CreateTable(model.User{}, model.Post{})

	// 	model.AddUser("bonfy", "abc123", "i@bonfy.im")
	// 	model.AddUser("rene", "abc123", "rene@test.com")

	// 	u1, _ := model.GetUserByUsername("bonfy")
	// 	u1.CreatePost("Beautiful day in Portland!")
	// 	model.UpdateAboutMe(u1.Username, `I'm the author of Go-Mega Tutorial you are reading now!`)

	// 	u2, _ := model.GetUserByUsername("rene")
	// 	u2.CreatePost("The Avengers movie was so cool!")
	// 	u2.CreatePost("Sun shine is beautiful")

	// 	u1.Follow(u2.Username)
}
