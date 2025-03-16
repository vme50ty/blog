package model

import "time"

// Comment 评论结构
type Comment struct {
	ID        int `gorm:"primary_key"`
	PostID    int `gorm:"index:idx_post_id"` // 外键关联帖子
	UserID    int `gorm:"index:idx_user_id"`
	User      User
	Body      string     `gorm:"type:varchar(180)"`
	Timestamp *time.Time `sql:"DEFAULT:current_timestamp"`
}

// 获取单个帖子（含评论）
func GetPostByID(id int) (*Post, error) {
	var post Post
	if err := db.Preload("User").
		Preload("Comments").
		Preload("Comments.User"). // 加载评论的用户信息
		First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func (c *Comment) FormattedTimeAgo() string {
	// 实现你的时间格式化逻辑，和Post相同的实现方式
	return FromTime(*c.Timestamp)
}
