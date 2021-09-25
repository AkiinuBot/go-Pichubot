package pichubot

import (
	"strconv"
	"time"
)

//* Events
type ShortEvent struct {
	Channel *chan map[string]interface{}
}
type LongEvent struct {
	UserID   int64
	GroupID  int64
	Channel  *chan string
	EventKey string
	EventID  string
}

var ShortEvents = make(map[string]ShortEvent) // 短事件容器
var LongEvents = make(map[string]LongEvent)   // 长事件容器

// 新建长事件
func NewEvent(userid int64, groupid int64, key string) LongEvent {
	var eventid string = strconv.FormatInt(time.Now().UnixNano(), 10) // 事件ID 以时间戳定义
	ch := make(chan string)
	event := LongEvent{Channel: &ch, UserID: userid, GroupID: groupid, EventKey: key, EventID: eventid}
	LongEvents[eventid] = event
	return event
}

// 关闭事件
func (event LongEvent) Close() {
	close(*event.Channel)
	delete(LongEvents, event.EventID)
}
