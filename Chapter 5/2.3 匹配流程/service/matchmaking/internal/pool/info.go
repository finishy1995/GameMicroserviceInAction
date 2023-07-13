package pool

// Info 匹配信息
type Info struct {
	ticketId    string // 当前匹配的匹配 ticket ID
	userId      string // 当前匹配的用户 ID
	createdTime int64  // 创建时间
	endTime     int64  // 结束时间
	status      Status // 状态
}

func (i *Info) Output() (int32, int64) {
	return int32(i.status), i.endTime
}
