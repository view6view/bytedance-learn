package main

import (
	"gopkg.in/gin-gonic/gin.v1"
	"homework4/controller"
	"homework4/repository"
	"os"
)

func main() {
	// 初始化数据索引
	if err := Init("./data/"); err != nil {
		os.Exit(-1)
	}
	// 初始化gin框架引擎配置
	r := gin.Default()
	// 构建路由
	// 查询帖子
	r.GET("/community/page/get/:id", func(c *gin.Context) {
		topicId := c.Param("id")
		data := controller.QueryPageInfo(topicId)
		c.JSON(200, data)
	})
	// 新增帖子
	r.POST("/community/topic/add", func(c *gin.Context) {
		// 传入控制器处理
		data := controller.AddTopic(c)
		c.JSON(200, data)
	})
	// 新增评论
	r.POST("/community/post/add", func(c *gin.Context) {
		// 传入控制器处理
		data := controller.AddPost(c)
		c.JSON(200, data)
	})
	// 启动服务
	err := r.Run()
	if err != nil {
		return
	}
}

func Init(filePath string) error {
	if err := repository.Init(filePath); err != nil {
		return err
	}
	return nil
}
