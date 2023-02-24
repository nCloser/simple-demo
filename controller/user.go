package controller

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// usersLoginInfo use map to store user info, and key is username+password for demo
// user data will be cleared every time the server starts
// test data: username=zhanglei, password=douyin
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

// var userIdSequence int64

type UserLoginResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
	UserId     int64  `json:"user_id,omitempty"`
	Token      string `json:"token"`
}

type UserResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
	User       User   `json:"user"`
}

func Register(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password // token 是否可以加密，
	var users UserDto
	res := DB.Where(&UserDto{Username: username}).Find(&users)
	// fmt.Println(res.RowsAffected)
	if res.RowsAffected == 1 { //用户存在且只有一个
		c.JSON(http.StatusOK, UserLoginResponse{
			StatusCode: 1,
			StatusMsg:  "用户已经存在",
		})
	} else if res.RowsAffected == 0 { //用户不存在
		// atomic.AddInt64(&userIdSequence, 1) // 原子
		newUser, res := CreateNewUser(username, password, token)
		if res.Error != nil {
			c.JSON(http.StatusOK, UserLoginResponse{
				StatusCode: 1,
				StatusMsg:  "创建用户失败",
			})
		} else {
			c.JSON(http.StatusOK, UserLoginResponse{
				StatusCode: 0,
				StatusMsg:  "注册成功",
				UserId:     newUser.Id,
				Token:      username + password,
			})
		}
	} else if res.Error != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			StatusCode: 1,
			StatusMsg:  "数据库错误",
		})
	}
}

func CreateNewUser(username string, password string, token string) (UserDto, *gorm.DB) {

	default_avatar, _ := ioutil.ReadFile("public/campus4-c0-01999.png")
	default_avatar_str := base64.StdEncoding.EncodeToString(default_avatar)
	default_bgimg, _ := ioutil.ReadFile("public/bg.jpg")
	default_bgimg_str := base64.StdEncoding.EncodeToString(default_bgimg)

	newUser := UserDto{
		Username:        username,
		Password:        password,
		Token:           token,
		FollowCount:     0,
		FollowerCount:   0,
		Avatar:          "data:image/png;base64," + default_avatar_str,
		BackgroundImage: "data:image/jpg;base64," + default_bgimg_str,
		Signature:       "gogogo",
		TotalFavorited:  0,
		WorkCount:       0,
		FavoriteCount:   0,
	}

	res := DB.Omit("Id", "IsFollow").Create(&newUser)
	return newUser, res
}

func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	var userFinded UserDto
	if res := DB.Where(&UserDto{Username: username}).First(&userFinded); res.Error == nil {
		if userFinded.Password == password { //密码正确
			c.JSON(http.StatusOK, UserLoginResponse{
				StatusCode: 0,
				StatusMsg:  "登录成功",
				UserId:     userFinded.Id,
				Token:      token,
			})
		} else { //TODO 可以加多次重试禁止登录等
			c.JSON(http.StatusOK, UserLoginResponse{
				StatusCode: 1,
				StatusMsg:  "密码错误，请重试",
			})
		}
	} else if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, UserLoginResponse{
			StatusCode: 1,
			StatusMsg:  "用户不存在，请注册",
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			StatusCode: 1,
			StatusMsg:  "数据读取错误",
		})
	}
}

func UserInfo(c *gin.Context) {
	token := c.Query("token")
	id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)

	var userFinded UserDto
	if res := DB.Where(&UserDto{Token: token, Id: id}).First(&userFinded); res.Error == nil { //验证token和id， TODO 如果使用Find找到多个怎么办
		c.JSON(http.StatusOK, UserResponse{
			StatusCode: 0,
			StatusMsg:  "用户信息读取成功",
			User: User{
				Id:              userFinded.Id,
				Name:            userFinded.Username,
				FollowCount:     userFinded.FollowCount,
				FollowerCount:   userFinded.FollowerCount,
				IsFollow:        userFinded.IsFollow,
				Avatar:          userFinded.Avatar,
				BackgroundImage: userFinded.BackgroundImage,
				Signature:       userFinded.Signature,
				TotalFavorited:  userFinded.TotalFavorited,
				WorkCount:       userFinded.WorkCount,
				FavoriteCount:   userFinded.FavoriteCount,
			},
		})
	} else if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusOK, UserLoginResponse{
			StatusCode: 1,
			StatusMsg:  "用户不存在，请注册",
		})
	} else {
		c.JSON(http.StatusOK, UserLoginResponse{
			StatusCode: 1,
			StatusMsg:  "数据读取错误",
		})
	}
}
