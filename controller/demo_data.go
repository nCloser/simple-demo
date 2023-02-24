package controller

var DemoVideos = []Video{
	{
		Id:            1,
		Author:        DemoUser,
		PlayUrl:       "http://219.216.99.9:10111/1.mp4",
		CoverUrl:      "http://219.216.99.9:10111/bg.jpg", // CoverUrl:      "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg",
		FavoriteCount: 1,
		CommentCount:  0,
		IsFavorite:    true,
		Title:         "熊霸天下",
	},
	{
		Id:            2,
		Author:        DemoUser,
		PlayUrl:       "http://219.216.99.9:10111/0.mp4",
		CoverUrl:      "http://219.216.99.9:10111/bg.jpg",
		FavoriteCount: 0,
		CommentCount:  0,
		IsFavorite:    false,
	},
}

// var DemoComments = []Comment{
// 	{
// 		Id:         1,
// 		User:       DemoUser,
// 		Content:    "Test Comment",
// 		CreateDate: "05-01",
// 	},
// 	{
// 		Id:         2,
// 		User:       DemoUser,
// 		Content:    "Test Comment",
// 		CreateDate: "05-01",
// 	},
// }

var DemoUser = User{
	Id:            1,
	Name:          "TestUser",
	FollowCount:   0,
	FollowerCount: 0,
	IsFollow:      true,
}
