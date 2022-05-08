package repository

// AddTopic 新增topic
func AddTopic(id int64, topic *Topic) {
	topicIndexMap[id] = topic
}

// AddPost 新增post
func AddPost(topicId int64, post *Post) {
	posts := postIndexMap[topicId]
	posts = append(posts, post)
	postIndexMap[topicId] = posts
}

// IsExist 判断id对应topic是否存在
func IsExist(topicId int64) bool {
	return topicIndexMap[topicId] != nil
}

// QueryTopicIds 获取全部的id数组
func (*TopicDao) QueryTopicIds() []int {
	// 返回数据
	ids := make([]int, 0, len(topicIndexMap))
	for id, _ := range topicIndexMap {
		// 将64位int转化为int
		//idStr := strconv.FormatInt(id, 10)
		//id, _ := strconv.Atoi(idStr)
		ids = append(ids, int(id))
	}
	return ids
}

// QueryPostIds 获取全部的id数组
func (*PostDao) QueryPostIds() []int {
	// 定义可变长数据
	var ids []int
	//ids := make([]int, 0, len(postIndexMap))
	for _, posts := range postIndexMap {
		for _, post := range posts {
			// 将64位int转化为int
			ids = append(ids, int(post.Id))
		}
	}
	return ids
}
