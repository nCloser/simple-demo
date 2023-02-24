package controller

import (
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	title := c.PostForm("title") // 作品描述

	var user UserDto
	res := DB.Where(&UserDto{Token: token}).Find(&user)
	if res.Error != nil {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "数据库错误"})
	}
	if res.RowsAffected == 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户不存在，无法上传"})
		return
	}

	data, err := c.FormFile("data") // TODO 视频的时间可以限制一下
	//同一用户上传同一个文件会有错误，由于文件名相同会出现网络错误的提示，已解决
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	filesuffix := filepath.Ext(data.Filename)
	finalName := fmt.Sprintf("%d_%d_%s", time.Now().Unix(), user.Id, filename)
	saveFile := filepath.Join("./public/", finalName)
	videoName := finalName[0 : len(finalName)-len(filesuffix)] //去掉后缀

	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "上传错误： " + err.Error(),
		})
		return
	}

	coverFile := videoName + "_thumbnail.png"
	outputPath := filepath.Join("./public/", videoName+"_thumbnail.png")
	if err := GenerateThumbnail(saveFile, outputPath); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "生成缩略图错误： " + err.Error(),
		})
	}

	//存入数据库
	if err := SaveVideoInfo(user, finalName, coverFile, title); err != nil {
		fmt.Println("存入数据库失败")
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "上传视频错误： " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " 上传成功！",
	})

	// 用户作品数更新
	var videoDtoList []VideoDto
	if user.Id == 0 {
		DB.Where(VideoDto{Token: token}).Find(&videoDtoList)
	} else {
		DB.Where(VideoDto{AuthorId: user.Id}).Find(&videoDtoList)
	}
	DB.Where(UserDto{Id: user.Id}).First(&user)
	user.WorkCount = int64(len(videoDtoList))
	DB.Save(&user)
}

func GenerateThumbnail(videoName string, output string) error {
	cmd := "ffmpeg -i " + videoName + " -f image2 -t 0.001 " + output
	// err := exec.Command("cmd", "/c", cmd).Run()
	err := exec.Command("/bin/bash", "-c", cmd).Run()
	if err != nil {
		return err
	}
	return nil
}

func SaveVideoInfo(user UserDto, videoFile string, coverFile string, title string) error {
	videoDto := VideoDto{
		AuthorId:      user.Id,
		PlayUrl:       ServerAddress + videoFile,
		CoverUrl:      ServerAddress + coverFile,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
		ReleaseTime:   time.Now().Unix(), //TODO 这个时间可以再考虑一下
	}

	res := DB.Omit("id", "is_favorite").Create(&videoDto).Error
	return res
}

// PublishList all users have same publish video list
func PublishList(c *gin.Context) {
	/*
		用户的视频发布列表，直接列出用户所有投稿过的视频
	*/
	token := c.Query("token")
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)

	videoList, err := GetPublishedVideoList(token, userId)
	if err != nil {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 1,
			},
		})
		return
	}
	// FIXME 上传作品后，不会马上更新作品数，需要重新登录才会更新
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}

func GetPublishedVideoList(token string, userId int64) ([]Video, error) {
	var videoDtoList []VideoDto
	var err error
	if userId == 0 {
		err = DB.Where(VideoDto{Token: token}).Find(&videoDtoList).Error
	} else {
		err = DB.Where(VideoDto{AuthorId: userId}).Find(&videoDtoList).Error
	}
	if err != nil {
		return make([]Video, 0), err
	}

	// 用户作品数更新
	var user UserDto
	DB.Where(UserDto{Id: userId}).First(&user)
	user.WorkCount = int64(len(videoDtoList))
	DB.Save(&user)

	videoList := make([]Video, len(videoDtoList))
	for i := 0; i < len(videoDtoList); i++ {
		videoList[i] = VideoDto2Video(videoDtoList[i])
	}
	return videoList, err
}
