package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func init() {
	projectName := "go-web"
	getConfig(projectName)
}

func getConfig(projectName string) {
	viper.SetConfigName("config.yml.sample") // 指定完整的配置文件名
	viper.SetConfigType("yml")

	// 添加多个配置文件查找路径
	viper.AddConfigPath(".")                                   // 当前目录
	viper.AddConfigPath("..")                                  // 父目录
	viper.AddConfigPath("../..")                               // 祖父目录
	viper.AddConfigPath(fmt.Sprintf("$HOME/.%s", projectName)) // 自定义目录
	viper.AddConfigPath(fmt.Sprintf("/data/docker/config/%s", projectName))

	// 读取配置文件
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}
}

func GetMysqlConnectingString() string {
	usr := viper.GetString("mysql.user")
	pwd := viper.GetString("mysql.password")
	host := viper.GetString("mysql.host")
	db := viper.GetString("mysql.db")
	charset := viper.GetString("mysql.charset")

	return fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=%s&parseTime=true", usr, pwd, host, db, charset)
}

// GetSMTPConfig func
func GetSMTPConfig() (server string, port int, user, pwd string) {
	server = viper.GetString("mail.smtp")
	port = viper.GetInt("mail.smtp-port")
	user = viper.GetString("mail.user")
	pwd = viper.GetString("mail.password")
	return
}

func GetServerURL() (url string) {
	url = viper.GetString("server.url")
	return
}
