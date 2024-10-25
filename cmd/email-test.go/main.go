package main

import (
	"crypto/tls"
	"fmt"
	"go-web/config"
	"net"
	"time"

	"gopkg.in/gomail.v2"
)

// password:SFtBfmxbCPY8EMt3
// password: PUgAUwhg35ANx9Le

func main() {

	// email := "2537012722@qq.com"
	// SendEmail(email, "test email", "<h1>测试邮件内容</h1>")
	time1 := time.Now()
	loc, _ := time.LoadLocation("Asia/Shanghai")
	time2 := time1.In(loc)
	fmt.Println(time2.Format("2006-01-02 15:04:05"))
}

func SendEmail(target, subject, content string) {
	server, port, usr, pwd := config.GetSMTPConfig()
	fmt.Println("SMTP服务器:", server, "端口:", port, "用户名:", usr, "密码:", pwd)
	// 创建一个新的邮件发送器
	d := gomail.NewDialer(server, port, usr, pwd)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// 尝试连接并获取 SendCloser
	sendCloser, err := d.Dial()
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			fmt.Println("连接超时:", err)
		} else if opErr, ok := err.(*net.OpError); ok {
			fmt.Println("操作错误:", opErr)
		} else {
			fmt.Println("连接错误:", err)
		}
		return
	}
	defer sendCloser.Close() // 确保在函数结束时关闭连接

	fmt.Println("成功连接到SMTP服务器")

	// 创建邮件
	m := gomail.NewMessage()
	m.SetHeader("From", usr)
	m.SetHeader("To", target)
	m.SetAddressHeader("Cc", usr, "admin")
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", content)

	// 发送邮件
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("邮件发送错误:", err)
		return
	}
	fmt.Println("邮件发送成功！")
}
