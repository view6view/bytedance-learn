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
	AddPostLock sync.Mutex
)

// AddPost 新增帖子
func AddPost(post *repository.Post) (bool, error) {
	return NewAddPostFlow(post).Do()
}

func NewAddPostFlow(post *repository.Post) *AddPostFlow {
	return &AddPostFlow{
		post: post,
	}
}

type AddPostFlow struct {
	post *repository.Post
}

// Do 业务流程编排
func (f *AddPostFlow) Do() (bool, error) {
	// 参数校验
	if err := f.checkParam(); err != nil {
		return false, err
	}
	// 加锁，防止并发造成数据异常
	AddPostLock.Lock()
	// 完善数据,内部系统稳定运行,暂不考虑异常
	f.prepareInfo()
	// 存储数据,内部系统稳定运行,暂不考虑异常
	f.solveData()
	// 解锁
	AddPostLock.Unlock()
	return true, nil
}

// 对新增post的参数进行校验
func (f *AddPostFlow) checkParam() error {
	// 校验是否存在对应的topic
	if !repository.IsExist(f.post.ParentId) {
		return errors.New("no corresponding topic exists")
	}
	if len(f.post.Content) > 500 {
		return errors.New("topic content length must be less than 500")
	}
	return nil
}

// 对数据进行封装
func (f *AddPostFlow) prepareInfo() {
	// 设置创建时间
	f.post.CreateTime = time.Now().Unix()
	// 获取全部id数组
	ids := repository.NewPostDaoInstance().QueryPostIds()
	// 设置id
	f.post.Id = int64(util.GetUniId(ids))
}

// 处理数据
func (f *AddPostFlow) solveData() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	// 向map里面新增
	go func() {
		defer wg.Done()
		repository.AddPost(f.post.ParentId, f.post)
	}()
	// 向文件里面写入
	go func() {
		defer wg.Done()
		open, _ := os.OpenFile(repository.FilePath+"post", os.O_WRONLY|os.O_APPEND, 0666)
		defer open.Close()
		writer := bufio.NewWriter(open)
		buf, _ := json.Marshal(f.post)
		topicStr := string(buf)
		writer.WriteString("\n" + topicStr)
		writer.Flush()
	}()
	wg.Wait()
}
