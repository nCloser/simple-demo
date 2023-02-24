package controller

import (
	"fmt"

	"gorm.io/gorm"
)

var ServerAddress = "http://219.216.99.9:10111/"

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

// 用于查询的Video结构体
type Video struct {
	Id            int64  `json:"id,omitempty"`
	Author        User   `json:"author,omitempty"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"title,omitempty"`
}

// 用于存储的Video结构体
type VideoDto struct {
	Id            int64  `json:"id,omitempty" gorm:"AUTO_INCREMENT"`
	AuthorId      int64  `json:"author_id"`
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty" gorm:"default:false"`
	Title         string `json:"title,omitempty"`
	ReleaseTime   int64  `json:"release_time,omitempty"`
	Token         string `json:"token,omitempty"`
}

// 评论id和创建时间由gorm.Model自动生成
type Comment struct {
	gorm.Model
	UserId      int64   `gorm:"not null;comment:创作用户ID"`
	VideoId     int64   `gorm:"not null;comment:视频ID"`
	CommentText string  `gorm:"type: text;not null;comment:评论内容"`
	UserInfo    UserDto `gorm:"foreignKey:UserId; references:ID; comment:评论所属用户"`
}

type Favorite struct {
	gorm.Model
	UserId  int64 `gorm:"user_id"`
	VideoId int64 `gorm:"video_id"`
}

// 用于查询的User结构体
type User struct {
	Id              int64  `json:"id,omitempty"`
	Name            string `json:"name,omitempty"`
	FollowCount     int64  `json:"follow_count,omitempty"`
	FollowerCount   int64  `json:"follower_count,omitempty"`
	IsFollow        bool   `json:"is_follow,omitempty"`        // true-已关注，false-未关注
	Avatar          string `json:"avatar,omitempty"`           // 用户头像
	BackgroundImage string `json:"background_image,omitempty"` // 用户个人页顶部大图
	Signature       string `json:"signature,omitempty"`        // 个人简介
	TotalFavorited  int64  `json:"total_favorited,omitempty"`  // 获赞数量
	WorkCount       int64  `json:"work_count,omitempty"`       // 作品数
	FavoriteCount   int64  `json:"favorite_count,omitempty"`   // 喜欢数
}

// 用于插入的User结构体
type UserDto struct {
	Id              int64  `json:"id,omitempty" gorm:"AUTO_INCREMENT"`
	Username        string `json:"username,omitempty"`
	Password        string `json:"password,omitempty"`
	Token           string `json:"token,omitempty"`
	FollowCount     int64  `json:"follow_count,omitempty"`
	FollowerCount   int64  `json:"follower_count,omitempty"`
	IsFollow        bool   `json:"is_follow,omitempty" gorm:"default:false"` // true-已关注，false-未关注
	Avatar          string `json:"avatar,omitempty"`                         // 用户头像
	BackgroundImage string `json:"background_image,omitempty"`               // 用户个人页顶部大图
	Signature       string `json:"signature,omitempty"`                      // 个人简介
	TotalFavorited  int64  `json:"total_favorited,omitempty"`                // 获赞数量
	WorkCount       int64  `json:"work_count,omitempty"`                     // 作品数
	FavoriteCount   int64  `json:"favorite_count,omitempty"`                 // 喜欢数
}

type Message struct {
	Id         int64  `json:"id,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateTime string `json:"create_time,omitempty"`
}

type MessageSendEvent struct {
	UserId     int64  `json:"user_id,omitempty"`
	ToUserId   int64  `json:"to_user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

type MessagePushEvent struct {
	FromUserId int64  `json:"user_id,omitempty"`
	MsgContent string `json:"msg_content,omitempty"`
}

func UserDto2User(userDto UserDto) User {
	user := User{
		Id:              userDto.Id,
		Name:            userDto.Username,
		FollowCount:     userDto.FollowCount,
		FollowerCount:   userDto.FollowerCount,
		IsFollow:        userDto.IsFollow,
		Avatar:          userDto.Avatar,
		BackgroundImage: userDto.BackgroundImage,
		Signature:       userDto.Signature,
		TotalFavorited:  userDto.TotalFavorited,
		WorkCount:       userDto.WorkCount,
		FavoriteCount:   userDto.FavoriteCount,
	}
	return user
}

func VideoDto2Video(videoDto VideoDto) Video {
	var user UserDto
	res := DB.Where(UserDto{Id: videoDto.AuthorId}).Find(&user)
	if res.Error != nil {
		fmt.Println(res.Error) // TODO 改成log更好
	}
	video := Video{
		Id:            videoDto.Id,
		Author:        UserDto2User(user),
		PlayUrl:       videoDto.PlayUrl,
		CoverUrl:      videoDto.CoverUrl,
		FavoriteCount: videoDto.FavoriteCount,
		CommentCount:  videoDto.CommentCount,
		IsFavorite:    videoDto.IsFavorite,
		Title:         videoDto.Title,
	}
	return video
}
