package repository

import (
	"sync"
)

type Topic struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
}
type TopicDao struct {
}

var (
	topicDao *TopicDao
	// Once 是将只执行一个操作的对象。
	topicOnce sync.Once
)

// NewTopicDaoInstance 初始化TopicDao的单例，为外部调用提供一个接收者对象
func NewTopicDaoInstance() *TopicDao {
	topicOnce.Do(
		func() {
			topicDao = &TopicDao{}
		})
	return topicDao
}

// QueryTopicById 通过*TopicDao函数接收者调用执行查询方法
func (*TopicDao) QueryTopicById(id int64) *Topic {
	// 返回数据
	return topicIndexMap[id]
}
