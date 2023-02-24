package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	token := c.Query("token")
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	actionType := c.Query("action_type")

	var user UserDto
	err := DB.Where(&UserDto{Token: token}).First(&user).Error
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户鉴权失败"})
		return
	}

	var favoriteVideoDto VideoDto
	err = DB.Where(VideoDto{Id: videoId}).First(&favoriteVideoDto).Error
	if err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "数据库读取错误/数据为空"})
		return
	}

	var favoriteInfo Favorite

	switch actionType {
	case "1":
		favoriteVideoDto.IsFavorite = true //TODO bool类型有原子性操作可以用吗
		// 喜欢列表更新
		favoriteInfo = Favorite{
			UserId:  user.Id,
			VideoId: favoriteVideoDto.Id,
		}
		DB.Create(&favoriteInfo)
		// 用户点赞数更新
		user.FavoriteCount = DB.Where(Favorite{UserId: user.Id}).Find(&favoriteInfo).RowsAffected
		// 视频点赞数更新
		favoriteVideoDto.FavoriteCount = DB.Where(Favorite{VideoId: favoriteVideoDto.Id}).Find(&favoriteInfo).RowsAffected
	case "2":
		favoriteVideoDto.IsFavorite = false

		// 喜欢列表更新
		favoriteInfo = Favorite{
			UserId:  user.Id,
			VideoId: favoriteVideoDto.Id,
		}
		DB.Model(Favorite{}).Where("user_id = ? AND video_id = ?", user.Id, favoriteVideoDto.Id).Delete(&favoriteInfo)
		// 用户点赞数更新
		user.FavoriteCount = DB.Where(Favorite{UserId: user.Id}).Find(&favoriteInfo).RowsAffected
		// 视频点赞数更新
		favoriteVideoDto.FavoriteCount = DB.Where(Favorite{VideoId: favoriteVideoDto.Id}).Find(&favoriteInfo).RowsAffected
	}
	fmt.Println(user.FavoriteCount)
	DB.Save(&favoriteVideoDto)
	DB.Save(&user)

	c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "点赞数据更新"})
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	token := c.Query("token")
	// 鉴权
	var user UserDto

	if err := DB.Where(&UserDto{Token: token, Id: userId}).First(&user).Error; err != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户鉴权失败"})
		return
	}
	//
	videoList, err := GetFavoriteVideoList(token, userId)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 1,
			},
		})
		return
	}
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}

func GetFavoriteVideoList(token string, userId int64) ([]Video, error) {
	var videoDto VideoDto
	var favoriteList []Favorite

	var err error
	if userId != 0 {
		err = DB.Where(Favorite{UserId: userId}).Find(&favoriteList).Error
	} else {
		err = gorm.ErrEmptySlice
		return make([]Video, 0), err
	}

	videoList := make([]Video, len(favoriteList))
	for i := 0; i < len(favoriteList); i++ {
		err = DB.Where(VideoDto{Id: favoriteList[i].VideoId}).First(&videoDto).Error
		if err != nil {
			return make([]Video, 0), err
		}
		videoList[i] = VideoDto2Video(videoDto)
	}
	return videoList, err
}
