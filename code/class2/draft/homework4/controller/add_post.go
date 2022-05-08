package controller

import (
	"encoding/json"
	"gopkg.in/gin-gonic/gin.v1"
	"homework4/repository"
	"homework4/service"
)

// AddPost 新增帖子控制器
func AddPost(c *gin.Context) *PageData {
	// 获取请求body
	data, err := c.GetRawData()
	// 如果获取请求体错误就返回错误码
	if err != nil {
		return &PageData{
			Code: -1,
			Msg:  err.Error(),
		}
	}
	// 解析json实体
	var post = repository.Post{}
	err = json.Unmarshal(data, &post)
	// 如果解析json错误返回错误码
	if err != nil {
		return &PageData{
			Code: -1,
			Msg:  err.Error(),
		}
	}
	ok, err := service.AddPost(&post)
	// 如果业务错误返回业务错误码
	if err != nil {
		return &PageData{
			Code: -1,
			Msg:  err.Error(),
		}
	}
	// 返回封装实体消息
	if ok {
		return &PageData{
			Code: 0,
			Msg:  "success",
			Data: nil,
		}
	} else {
		return &PageData{
			Code: 500,
			Msg:  "fail",
			Data: "新增topic，未知错误!",
		}
	}
}
