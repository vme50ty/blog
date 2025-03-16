package model

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt"
)

type User struct {
	ID           int    `gorm:"primary_key"`
	Username     string `gorm:"type:varchar(64)"`
	Email        string `gorm:"type:varchar(120)"`
	PasswordHash string `gorm:"type:varchar(128)"`
	LastSeen     *time.Time
	AboutMe      string `gorm:"type:varchar(140)"`
	Avatar       string `gorm:"type:varchar(200)"`
	Posts        []Post
	Followers    []*User `gorm:"many2many:follower;association_jointable_foreignkey:follower_id"`
}

func (u *User) SetPassword(password string) {
	u.PasswordHash = GeneratePasswordHash(password)
}

// SetAvatar func
func (u *User) SetAvatar(email string) {
	u.Avatar = fmt.Sprintf("https://www.gravatar.com/avatar/%s?d=identicon", Md5(email))
}

// CheckPassword func
func (u *User) CheckPassword(password string) bool {
	return GeneratePasswordHash(password) == u.PasswordHash
}

// GetUserByUsername func
func GetUserByUsername(username string) (*User, error) {
	var user User
	if err := db.Where("username=?", username).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// AddUser func
func AddUser(username, password, email string) error {
	user := User{Username: username, Email: email}
	user.SetPassword(password)
	user.SetAvatar(email)
	if err := db.Create(&user).Error; err != nil {
		return err
	}
	return user.FollowSelf()
}

// UpdateUserByUsername func
func UpdateUserByUsername(username string, contents map[string]interface{}) error {
	item, err := GetUserByUsername(username)
	if err != nil {
		return err
	}
	return db.Model(item).Updates(contents).Error
}

// UpdateLastSeen func
func UpdateLastSeen(username string) error {
	contents := map[string]interface{}{"last_seen": time.Now()}
	return UpdateUserByUsername(username, contents)
}

// UpdateAboutMe func
func UpdateAboutMe(username, text string) error {
	contents := map[string]interface{}{"about_me": text}
	return UpdateUserByUsername(username, contents)
}

// Follow func
// follow someone usr_id other.id follow_id u.id
func (u *User) Follow(username string) error {
	other, err := GetUserByUsername(username)
	if err != nil {
		return err
	}
	return db.Model(other).Association("Followers").Append(u).Error
}

// Unfollow func
func (u *User) Unfollow(username string) error {
	other, err := GetUserByUsername(username)
	if err != nil {
		return err
	}
	return db.Model(other).Association("Followers").Delete(u).Error
}

// FollowSelf func
func (u *User) FollowSelf() error {
	return db.Model(u).Association("Followers").Append(u).Error
}

// 获取当前用户的粉丝数量
func (u *User) FollowersCount() int {
	return db.Model(u).Association("Followers").Count()
}

// 获取当前用户u关注的人的ID列表  "user_id, follower_id"表示用户ID和粉丝ID的关联表
func (u *User) FollowingIDs() []int {
	var ids []int
	rows, err := db.Table("follower").Where("follower_id = ?", u.ID).Select("user_id, follower_id").Rows()
	if err != nil {
		log.Println("Counting Following error:", err)
		return ids
	}
	defer rows.Close()
	for rows.Next() {
		var id, followerID int
		rows.Scan(&id, &followerID)
		ids = append(ids, id)
	}
	return ids
}

// 获取当前用户u关注的人的数量
func (u *User) FollowingCount() int {
	ids := u.FollowingIDs()
	return len(ids)
}

// 获取当前用户u关注的人的文章列表
func (u *User) FollowingPosts() (*[]Post, error) {
	var posts []Post
	ids := u.FollowingIDs()
	if err := db.Preload("User").Order("timestamp desc").Where("user_id in (?)", ids).Find(&posts).Error; err != nil {
		return nil, err
	}
	return &posts, nil
}

// 用户u是否被username关注了
func (u *User) IsFollowedByUser(username string) bool {
	user, _ := GetUserByUsername(username)
	ids := user.FollowingIDs()
	// user关注的ID列表
	for _, id := range ids {
		if u.ID == id {
			return true
		}
	}
	return false
}

// CreatePost func
func (u *User) CreatePost(body string) error {
	post := Post{Body: body, UserID: u.ID}
	return db.Create(&post).Error
}

func (u *User) CreateComment(body string, postID int) error {
	comment := Comment{
		PostID: postID,
		UserID: u.ID,
		Body:   body,
	}
	return db.Create(&comment).Error
}

// FollowingPostsByPageAndLimit func
func (u *User) FollowingPostsByPageAndLimit(page, limit int) (*[]Post, int, error) {
	var total int
	var posts []Post
	offset := (page - 1) * limit
	ids := u.FollowingIDs()
	if err := db.Preload("User").Order("timestamp desc").Where("user_id in (?)", ids).Offset(offset).Limit(limit).Find(&posts).Error; err != nil {
		return nil, total, err
	}
	db.Model(&Post{}).Where("user_id in (?)", ids).Count(&total)
	return &posts, total, nil
}

// GenerateToken func
func (u *User) GenerateToken() (string, error) {
	// 使用 HS256 签名算法生成一个新的 token，包含用户的 username 和 exp（过期时间）
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": u.Username,                           // 将 username 放入 token 中
		"exp":      time.Now().Add(time.Hour * 2).Unix(), // 过期时间为当前时间加两小时
	})

	// 使用 secret 进行签名，并返回签名后的 token 字符串
	return token.SignedString([]byte("secret"))
}

// CheckToken func
func CheckToken(tokenString string) (string, error) {
	// 使用 jwt.Parse 函数解析 tokenString，并传入一个函数用于验证签名的密钥
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 检查签名算法是否是我们期望的 HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// 返回密钥 "secret" 用于验证签名
		return []byte("secret"), nil
	})

	// 如果 token 有效并且解析成功，提取 claims 中的 "username" 并返回
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["username"].(string), nil
	} else {
		// 如果 token 无效或解析失败，返回错误
		return "", err
	}
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	if err := db.Where("email=?", email).Find(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdatePassword func
func UpdatePassword(username, password string) error {
	contents := map[string]interface{}{"password_hash": Md5(password)}
	return UpdateUserByUsername(username, contents)
}
