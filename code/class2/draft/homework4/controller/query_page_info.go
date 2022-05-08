package controller

import (
	"homework4/service"
	"strconv"
)

// QueryPageInfo 查询帖子信息控制器
func QueryPageInfo(topicIdStr string) *PageData {
	topicId, err := strconv.ParseInt(topicIdStr, 10, 64)
	// 如果参数错误返回参数错误码
	if err != nil {
		return &PageData{
			Code: -1,
			Msg:  err.Error(),
		}
	}
	pageInfo, err := service.QueryPageInfo(topicId)
	// 如果业务错误返回业务错误码
	if err != nil {
		return &PageData{
			Code: -1,
			Msg:  err.Error(),
		}
	}
	// 返回封装实体消息
	return &PageData{
		Code: 0,
		Msg:  "success",
		Data: pageInfo,
	}
}
