package vm

import (
	"go-web/model"
	"strconv"
)

// PostViewModel 帖子详情页视图模型
type PostViewModel struct {
	BaseViewModel
	Post     model.Post
	Comments []model.Comment
	Flash    string
}

// PostViewModelOp 帖子视图操作结构
type PostViewModelOp struct{}

// GetVM 获取帖子视图模型
func (PostViewModelOp) GetVM(username, flash string, postID string) (PostViewModel, error) {
	v := PostViewModel{}
	v.SetTitle("Post Detail")

	// 转换postID为整型
	id, err := strconv.Atoi(postID)
	if err != nil {
		return v, err
	}

	// 获取帖子及评论数据
	post, err := model.GetPostByID(id)
	if err != nil {
		return v, err
	}

	v.Post = *post
	v.Comments = post.Comments
	v.Flash = flash
	v.SetCurrentUser(username)
	return v, nil
}

// CreateComment 创建评论方法
func (PostViewModelOp) CreateComment(username, body string, postID int) error {
	user, err := model.GetUserByUsername(username)
	if err != nil {
		return err
	}

	return user.CreateComment(body, postID)
}
