package controller

import (
	"html/template"

	"github.com/gorilla/sessions"
)

var (
	homeController home
	templates      map[string]*template.Template
	//一个模板映射（字典），键是字符串（例如模板文件的名称），值是指向 template.Template 的指针。这个字典存储了多个页面的模板，供后续渲染时使用。
	sessionName string
	//代表会话（session）的名称，用于标识用户的会话。
	flashName string
	//代表 Flash 消息的名称，用于在页面间传递消息。
	store     *sessions.CookieStore
	pageLimit int
)

func init() {
	templates = PopulateTemplates()
	store = sessions.NewCookieStore([]byte("something-very-secret"))
	// 创建一个新的 Cookie 存储，用于会话管理。
	//[]byte("something-very-secret") 是一个加密密钥，用于安全地存储和解密会话信息。这个密钥应当足够复杂以保证安全性。
	sessionName = "go-web"
	flashName = "go-flash"
	pageLimit = 5
}

// Startup func
func Startup() {
	homeController.registerRoutes()
}
