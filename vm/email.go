package vm

import (
	"go-web/config"
	"go-web/model"
)

// EmailViewModel struct
type EmailViewModel struct {
	Username string
	Token    string
	Server   string
}

// EmailViewModelOp struct
type EmailViewModelOp struct{}

// GetVM func
func (EmailViewModelOp) GetVM(email string) EmailViewModel {
	v := EmailViewModel{}
	u, _ := model.GetUserByEmail(email)
	v.Username = u.Username
	v.Token, _ = u.GenerateToken()
	v.Server = config.GetServerURL()
	return v
}
