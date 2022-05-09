package service

import (
	"bufio"
	"encoding/json"
	"errors"
	"homework4/repository"
	"homework4/util"
	"os"
	"sync"
	"time"
)

var (
	addTopicLock sync.Mutex
)

// AddTopic 新增帖子
func AddTopic(topic *repository.Topic) (bool, error) {
	return NewAddTopicFlow(topic).Do()
}

func NewAddTopicFlow(topic *repository.Topic) *AddTopicFlow {
	return &AddTopicFlow{
		topic: topic,
	}
}

type AddTopicFlow struct {
	topic *repository.Topic
}

// Do 业务流程编排
func (f *AddTopicFlow) Do() (bool, error) {
	// 参数校验
	if err := f.checkParam(); err != nil {
		return false, err
	}
	// 加锁，防止并发造成数据异常
	addTopicLock.Lock()
	// 完善数据,内部系统稳定运行,暂不考虑异常
	f.prepareInfo()
	// 存储数据,内部系统稳定运行.暂不考虑异常
	f.solveData()
	// 解锁
	addTopicLock.Unlock()
	return true, nil
}

// 对新增topic的参数进行校验
func (f *AddTopicFlow) checkParam() error {
	if len(f.topic.Title) > 20 {
		return errors.New("topic title length must be less than 20")
	}
	if len(f.topic.Content) > 500 {
		return errors.New("topic content length must be less than 500")
	}
	return nil
}

// 对数据进行封装
func (f *AddTopicFlow) prepareInfo() {
	// 设置创建时间
	f.topic.CreateTime = time.Now().Unix()
	// 获取全部id数组
	ids := repository.NewTopicDaoInstance().QueryTopicIds()
	// 设置id
	f.topic.Id = int64(util.GetUniId(ids))
}

// 处理数据
func (f *AddTopicFlow) solveData() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	// 向map里面新增
	go func() {
		defer wg.Done()
		repository.AddTopic(f.topic.Id, f.topic)
	}()
	// 向文件里面写入
	go func() {
		defer wg.Done()
		open, _ := os.OpenFile(repository.FilePath+"topic", os.O_WRONLY|os.O_APPEND, 0666)
		defer open.Close()
		writer := bufio.NewWriter(open)
		buf, _ := json.Marshal(f.topic)
		topicStr := string(buf)
		writer.WriteString("\n" + topicStr)
		writer.Flush()
	}()
	wg.Wait()
}
