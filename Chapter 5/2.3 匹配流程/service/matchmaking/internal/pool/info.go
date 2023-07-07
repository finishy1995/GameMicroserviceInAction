package pool

// info 匹配信息
type info struct {
	ticketId    string // 当前匹配的匹配 ticket ID
	createdTime int64  // 创建时间
	endTime     int64  // 结束时间
	status      Status // 状态
}
