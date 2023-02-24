package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedResponse struct {
	Response
	NextTime  int64   `json:"next_time,omitempty"`
	VideoList []Video `json:"video_list,omitempty"`
}

// 返回按投稿时间倒序的视频列表
const MAX_VIDEO_NUM int = 30

func getVideoList(time_node int64, token string) []Video {
	var user UserDto
	var tmp Favorite
	DB.Model(UserDto{}).Where(UserDto{Token: token}).First(&user)
	var videoDtoList []VideoDto
	err := DB.Model(VideoDto{}).Where("release_time < ?", time_node).Order("release_time desc").Limit(MAX_VIDEO_NUM).Find(&videoDtoList).Error
	if err != nil {
		fmt.Println(err)
	}
	videoList := make([]Video, len(videoDtoList))
	for i := 0; i < len(videoDtoList); i++ {

		if err := DB.Model(Favorite{}).Where("user_id = ? AND video_id = ?", user.Id, videoDtoList[i].Id).First(&tmp); err != nil { //当前用户是否喜欢
			videoDtoList[i].IsFavorite = false
		} else {
			videoDtoList[i].IsFavorite = true
		}

		videoList[i] = VideoDto2Video(videoDtoList[i])
	}
	return videoList
}

func Feed(c *gin.Context) {
	/*
		不限制登录状态，返回按投稿时间倒序的视频列表，视频数由服务端控制，单次最多30个
	*/
	timeNow := time.Now().Unix()

	// get request params
	latestTime := c.Query("latest_time") // 可选参数，限制返回视频的最新投稿时间戳，精确到秒，不填表示当前时间
	token := c.Query("token")            // 用户登录状态下设置

	//get time node
	var timeNode int64

	if latestTime != "" {
		timeNode, _ = strconv.ParseInt(latestTime, 10, 64)
	} else {
		timeNode = timeNow
	}
	fmt.Println(timeNode)

	videoList := getVideoList(timeNode, token)

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0, StatusMsg: "feed成功"},
		NextTime:  time.Now().Unix(),
		VideoList: videoList,
	})

}
