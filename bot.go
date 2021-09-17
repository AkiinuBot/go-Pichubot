package Pichubot

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/0ojixueseno0/go-Pichubot/utils"
	"github.com/gorilla/websocket"
	"github.com/wonderivan/logger"
)

//* Message Event
//? 消息事件
var OnPrivateMsg []func(eventinfo MessagePrivate) = []func(eventinfo MessagePrivate){} // 私聊消息事件
var OnGroupMsg []func(eventinfo MessageGroup) = []func(eventinfo MessageGroup){}       // 群聊消息事件

//* Notice Event
//? 提醒事件
var OnGroupUpload []func(eventinfo GroupUpload) = []func(eventinfo GroupUpload){}       // 群文件上传
var OnGroupAdmin []func(eventinfo GroupAdmin) = []func(eventinfo GroupAdmin){}          // 群管理员变动
var OnGroupDecrease []func(eventinfo GroupDecrease) = []func(eventinfo GroupDecrease){} // 群成员减少
var OnGroupIncrease []func(eventinfo GroupIncrease) = []func(eventinfo GroupIncrease){} // 群成员增加
var OnGroupBan []func(eventinfo GroupBan) = []func(eventinfo GroupBan){}                // 群聊禁言
var OnFriendAdd []func(eventinfo FriendAdd) = []func(eventinfo FriendAdd){}             // 已经添加好友后的事件
var OnGroupRecall []func(eventinfo GroupRecall) = []func(eventinfo GroupRecall){}       // 群消息撤回(群聊)
var OnFriendRecall []func(eventinfo FriendRecall) = []func(eventinfo FriendRecall){}    // 好友消息撤回(私聊)
var OnNotify []func(eventinfo Notify) = []func(eventinfo Notify){}                      // 群内戳一戳 群红包运气王 群成员荣誉变更

//* Request Event
//? 请求事件
var OnFriendRequest []func(eventinfo FriendRequest) = []func(eventinfo FriendRequest){} // 加好友请求
var OnGroupRequest []func(eventinfo GroupRequest) = []func(eventinfo GroupRequest){}    // 加群请求/邀请

//* Meta Event
//? 元事件
var OnMetaLifecycle []func(eventinfo MetaLifecycle) = []func(eventinfo MetaLifecycle){} // 生命周期
var OnMetaHeartbeat []func(eventinfo MetaHeartbeat) = []func(eventinfo MetaHeartbeat){} // 心跳包

//* Events
type ShortEvent struct {
	Channel *chan map[string]interface{}
}
type LongEvent struct {
	UserID    int64
	GroupID   int64
	Channel   *chan string
	Eventtype string
}

var ShortEvents = make(map[string]ShortEvent) // 短事件容器
var LongEvents = make(map[string]LongEvent)   // 长事件容器
var Connect *websocket.Conn

func NewBot() *Bot {
	return &Bot{}
}

func (bot *Bot) Run() {
	utils.InitLogger(bot.Config.Loglvl) // 初始化日志文件
	for {
		func(host string, path string) {
			//connet to websocket server
			url := url.URL{Scheme: "ws", Host: host, Path: path}
			var dailer *websocket.Dialer
			c, _, err := dailer.Dial(url.String(), nil)
			if err != nil {
				logger.Error(err)
				return
			} else {
				Connect = c // 传出接口
				for {
					_, message, err := c.ReadMessage()
					if err != nil {
						logger.Error(err)
					}
					m := make(map[string]interface{})
					if err := json.Unmarshal([]byte(message), &m); err != nil {
						logger.Error(err)
					}
					go msgParse(m)
				}
			}
		}(bot.Config.Host, bot.Config.Path)
		logger.Info("Websocket will reconnect in 5s")
		time.Sleep(5 * time.Second)
	}
}

func msgParse(receive map[string]interface{}) {

	switch receive["post_type"] {
	// 消息事件
	case "message":
		switch receive["message_type"] {
		// 私聊信息
		case "private":
			var eventinfo MessagePrivate = parsePrivate(receive)
			logger.Info(fmt.Sprintf("[↓][私聊][%s(%d)]: %s", eventinfo.Sender.Nickname, eventinfo.Sender.UserID, eventinfo.Message))
			for _, function := range OnPrivateMsg {
				function(eventinfo)
			}

		// 群聊信息
		case "group":
			var eventinfo MessageGroup = parseGroup(receive)
			logger.Info(fmt.Sprintf("[↓][群聊(%d)][%s(%d)]: %s", eventinfo.GroupID, eventinfo.Sender.Nickname, eventinfo.Sender.UserID, eventinfo.Message))
			for _, function := range OnGroupMsg {
				function(eventinfo)
			}

		default:
			logger.Warn(fmt.Sprintf("Cannot Parse 'message' event -> %s", receive))
		}

		// 通知事件
	case "notice":
		switch receive["notice_type"] {
		// 群文件上传
		case "group_upload":
			var eventinfo GroupUpload = parseGroupupload(receive)
			logger.Info(fmt.Sprintf("[N][群文件(%d)][%d]: %s", eventinfo.Group_id, eventinfo.User_id, eventinfo.File.Name))
			for _, function := range OnGroupUpload {
				function(eventinfo)
			}

			// 群管理员变动
		case "group_admin":
			var eventinfo GroupAdmin = parseGroupadmin(receive)
			var x string = "null"
			if eventinfo.Sub_type == "set" {
				x = "+"
			} else {
				x = "-"
			}
			logger.Info(fmt.Sprintf("[N][群(%d)管理][%s %d]", eventinfo.Group_id, x, eventinfo.User_id))
			for _, function := range OnGroupAdmin {
				function(eventinfo)
			}

			// 群成员减少
		case "group_decrease":
			var eventinfo GroupDecrease = parseGroupdecrease(receive)
			logger.Info(fmt.Sprintf("[N][成员退群(%d)][%d] Type: %s", eventinfo.Group_id, eventinfo.User_id, eventinfo.Sub_type))
			for _, function := range OnGroupDecrease {
				function(eventinfo)
			}

			// 群成员增加
		case "group_increase":
			var eventinfo GroupIncrease = parseGroupincrease(receive)
			logger.Info(fmt.Sprintf("[N][成员入群(%d)][%d -> %d] Type: %s", eventinfo.Group_id, eventinfo.Operator_id, eventinfo.User_id, eventinfo.Sub_type))
			for _, function := range OnGroupIncrease {
				function(eventinfo)
			}

			// 群禁言
		case "group_ban":
			var eventinfo GroupBan = parseGroupban(receive)
			logger.Info(fmt.Sprintf("[N][群聊(%d)] %d 禁言/解禁了 %d for %ds", eventinfo.Group_id, eventinfo.Operator_id, eventinfo.User_id, eventinfo.Duration))
			for _, function := range OnGroupBan {
				function(eventinfo)
			}

			// 好友添加
		case "friend_add":
			var eventinfo FriendAdd = parseFriendAdd(receive)
			logger.Info(fmt.Sprintf("[N][成功添加好友]%d", eventinfo.User_id))
			for _, function := range OnFriendAdd {
				function(eventinfo)
			}

			// 群消息撤回
		case "group_recall":
			var eventinfo GroupRecall = parseGrouprecall(receive)
			logger.Info(fmt.Sprintf("[N][群聊(%d)][%d] 撤回了消息(id): %d", eventinfo.Group_id, eventinfo.User_id, eventinfo.Message_id))
			for _, function := range OnGroupRecall {
				function(eventinfo)
			}

			// 好友消息撤回
		case "friend_recall":
			var eventinfo FriendRecall = parseFriendrecall(receive)
			logger.Info(fmt.Sprintf("[N][私聊][%d] 撤回了消息(id): %d", eventinfo.User_id, eventinfo.Message_id))
			for _, function := range OnFriendRecall {
				function(eventinfo)
			}

			// 群内戳一戳 群红包运气王 群成员荣誉变更
		case "notify":
			var eventinfo Notify = parseNotify(receive)
			logger.Info(fmt.Sprintf("[N][Notify][Group:%d] %d -> %s", eventinfo.Group_id, eventinfo.User_id, eventinfo.Sub_type))
			for _, function := range OnNotify {
				function(eventinfo)
			}

		default:
			logger.Warn(fmt.Sprintf("Cannot Parse 'notice' event -> %s", receive))
		}

		// 请求事件
	case "request":
		switch receive["request_type"] {
		// 添加好友申请
		case "friend":
			var eventinfo FriendRequest = parseFriendrequest(receive)
			logger.Info(fmt.Sprintf("[↓][好友申请] %d 申请加你为好友 -> %s", eventinfo.User_id, eventinfo.Comment))
			for _, function := range OnFriendRequest {
				function(eventinfo)
			}

			// 加群邀请
		case "group":
			// SetGroupInviteRequest(receive["flag"].(string), true, "") // 自动同意加群
			var eventinfo GroupRequest = parseGrouprequest(receive)
			logger.Info(fmt.Sprintf("[↓][加群/邀请] %d %s -> %d(验证信息: %s)", eventinfo.User_id, eventinfo.Sub_type, eventinfo.Group_id, eventinfo.Comment))
			for _, function := range OnGroupRequest {
				function(eventinfo)
			}

		default:
			logger.Warn(fmt.Sprintf("Cannot Parse 'request' event -> %s", receive))
		}
		// 元事件
	case "meta_event":
		switch receive["meta_event_type"] {
		// 生命周期
		case "lifecycle":
			var eventinfo MetaLifecycle = parseMetalifecycle(receive)
			logger.Debug(fmt.Sprintf("[↓][Lifecycle][%d] Type: %s", eventinfo.Self_id, eventinfo.Sub_type))
			for _, function := range OnMetaLifecycle {
				function(eventinfo)
			}

			// 心跳包
		case "heartbeat":
			var eventinfo MetaHeartbeat = parseMetaheartbeat(receive)
			logger.Debug(fmt.Sprintf("[↓][Heartbeat][%d] Type: %s", eventinfo.Self_id, eventinfo.Status))
			for _, function := range OnMetaHeartbeat {
				function(eventinfo)
			}

			// logger.Debug("Received a heartbeat package.")
		default:
			logger.Warn(fmt.Sprintf("Cannot Parse 'meta_event' event -> %s", receive))
		}
	default:
		// 短事件回调
		if _, ok := receive["echo"]; ok {
			if _, ok := ShortEvents[receive["echo"].(string)]; ok {
				*ShortEvents[receive["echo"].(string)].Channel <- receive
			}
		} else {
			logger.Warn(fmt.Sprintf("Got Error Package -> %s", receive))
		}
	}
}

func parsePrivate(r map[string]interface{}) MessagePrivate {
	e := MessagePrivate{
		SelfID:     r["self_id"].(int64),
		SubType:    r["sub_type"].(string),
		MessageID:  r["message_id"].(int64),
		UserID:     r["user_id"].(int64),
		Message:    r["message"].(string),
		RawMessage: r["raw_message"].(string),
		Sender: struct {
			UserID   int64
			Nickname string
			Sex      string
			Age      int64
		}{
			UserID:   r["sender"].(map[string]interface{})["user_id"].(int64),
			Nickname: r["sender"].(map[string]interface{})["nickname"].(string),
			Sex:      r["sender"].(map[string]interface{})["sex"].(string),
			Age:      r["sender"].(map[string]interface{})["age"].(int64),
		}}
	return e
}

func parseGroup(r map[string]interface{}) MessageGroup {
	e := MessageGroup{
		SelfID:     r["self_id"].(int64),
		SubType:    r["sub_type"].(string),
		MessageID:  r["message_id"].(int64),
		GroupID:    r["group_id"].(int64),
		UserID:     r["user_id"].(int64),
		Message:    r["message"].(string),
		RawMessage: r["raw_message"].(string)}
	switch e.SubType {
	case "normal":
		e.Sender = struct {
			UserID   int64
			Nickname string
			Card     string
			Sex      string
			Age      int64
			Area     string
			Level    string
			Role     string
			Title    string
		}{
			UserID:   r["sender"].(map[string]interface{})["user_id"].(int64),
			Nickname: r["sender"].(map[string]interface{})["nickname"].(string),
			Card:     r["sender"].(map[string]interface{})["card"].(string),
			Sex:      r["sender"].(map[string]interface{})["sex"].(string),
			Age:      r["sender"].(map[string]interface{})["age"].(int64),
			Area:     r["sender"].(map[string]interface{})["area"].(string),
			Level:    r["sender"].(map[string]interface{})["level"].(string),
			Role:     r["sender"].(map[string]interface{})["role"].(string),
			Title:    r["sender"].(map[string]interface{})["title"].(string)}
	case "anoymous":
		e.Anonymous = struct {
			Id   int64
			Name string
			Flag string
		}{
			Id:   r["anonymous"].(map[string]interface{})["id"].(int64),
			Name: r["anonymous"].(map[string]interface{})["name"].(string),
			Flag: r["anonymous"].(map[string]interface{})["flag"].(string)}
	}
	return e
}
func parseGroupupload(r map[string]interface{}) GroupUpload {
	e := GroupUpload{
		Time:     r["time"].(int64),
		Self_id:  r["self_id"].(int64),
		Group_id: r["group_id"].(int64),
		User_id:  r["user_id"].(int64),
		File: struct {
			Id    string
			Name  string
			Size  int64
			Busid int64
		}{
			Id:    r["file"].(map[string]interface{})["id"].(string),
			Name:  r["file"].(map[string]interface{})["name"].(string),
			Size:  r["file"].(map[string]interface{})["size"].(int64),
			Busid: r["file"].(map[string]interface{})["busid"].(int64),
		}}
	return e
}
func parseGroupadmin(r map[string]interface{}) GroupAdmin {
	e := GroupAdmin{
		Time:     r["time"].(int64),
		Self_id:  r["self_id"].(int64),
		Sub_type: r["sub_type"].(string),
		Group_id: r["group_id"].(int64),
		User_id:  r["user_id"].(int64),
	}
	return e
}
func parseGroupdecrease(r map[string]interface{}) GroupDecrease {
	e := GroupDecrease{
		Time:        r["time"].(int64),
		Self_id:     r["self_id"].(int64),
		Sub_type:    r["sub_type"].(string),
		Group_id:    r["group_id"].(int64),
		Operator_id: r["operator_id"].(int64),
		User_id:     r["user_id"].(int64),
	}
	return e
}
func parseGroupincrease(r map[string]interface{}) GroupIncrease {
	e := GroupIncrease{
		Time:        r["time"].(int64),
		Self_id:     r["self_id"].(int64),
		Sub_type:    r["sub_type"].(string),
		Group_id:    r["group_id"].(int64),
		Operator_id: r["operator_id"].(int64),
		User_id:     r["user_id"].(int64),
	}
	return e
}
func parseGroupban(r map[string]interface{}) GroupBan {
	e := GroupBan{
		Time:        r["time"].(int64),
		Self_id:     r["self_id"].(int64),
		Sub_type:    r["sub_type"].(string),
		Group_id:    r["group_id"].(int64),
		Operator_id: r["operator_id"].(int64),
		User_id:     r["user_id"].(int64),
		Duration:    r["duration"].(int64),
	}
	return e
}
func parseFriendAdd(r map[string]interface{}) FriendAdd {
	e := FriendAdd{
		Time:    r["time"].(int64),
		Self_id: r["self_id"].(int64),
		User_id: r["user_id"].(int64),
	}
	return e
}

func parseGrouprecall(r map[string]interface{}) GroupRecall {
	e := GroupRecall{
		Time:        r["time"].(int64),
		Self_id:     r["self_id"].(int64),
		Group_id:    r["group_id"].(int64),
		User_id:     r["user_id"].(int64),
		Operator_id: r["operator_id"].(int64),
		Message_id:  r["message_id"].(int64),
	}
	return e
}
func parseFriendrecall(r map[string]interface{}) FriendRecall {
	e := FriendRecall{
		Time:       r["time"].(int64),
		Self_id:    r["self_id"].(int64),
		User_id:    r["user_id"].(int64),
		Message_id: r["message_id"].(int64),
	}
	return e
}
func parseNotify(r map[string]interface{}) Notify {
	e := Notify{
		Time:     r["time"].(int64),
		Self_id:  r["self_id"].(int64),
		Sub_type: r["sub_type"].(string),
		Group_id: r["group_id"].(int64),
		User_id:  r["user_id"].(int64),
	}
	if e.Sub_type == "honor" {
		e.Honor_type = r["honor_type"].(string)
	} else {
		e.Target_id = r["target_id"].(int64)
	}
	return e
}
func parseFriendrequest(r map[string]interface{}) FriendRequest {
	e := FriendRequest{
		Time:    r["time"].(int64),
		Self_id: r["self_id"].(int64),
		User_id: r["user_id"].(int64),
		Comment: r["comment"].(string),
		Flag:    r["flag"].(string),
	}
	return e
}
func parseGrouprequest(r map[string]interface{}) GroupRequest {
	e := GroupRequest{
		Time:     r["time"].(int64),
		Self_id:  r["self_id"].(int64),
		Sub_type: r["sub_type"].(string),
		Group_id: r["group_id"].(int64),
		User_id:  r["user_id"].(int64),
		Comment:  r["comment"].(string),
		Flag:     r["flag"].(string),
	}
	return e
}

func parseMetalifecycle(r map[string]interface{}) MetaLifecycle {
	e := MetaLifecycle{
		Time:     r["time"].(int64),
		Self_id:  r["self_id"].(int64),
		Sub_type: r["sub_type"].(string),
	}
	return e
}
func parseMetaheartbeat(r map[string]interface{}) MetaHeartbeat {
	e := MetaHeartbeat{
		Time:     r["time"].(int64),
		Self_id:  r["self_id"].(int64),
		Status:   r["status"],
		Interval: r["interval"].(int64),
	}
	return e
}
