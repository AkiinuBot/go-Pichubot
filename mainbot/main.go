package mainbot

import (
	"flag"
	"fmt"
	"net/url"
	"time"

	"github.com/0ojixueseno0/go-Pichubot-base/pichumod"
	"github.com/gorilla/websocket"
	"github.com/hjson/hjson-go"
)

var Fastmode = flag.Bool("faststart", false, "是否快速启动")

func Run() {
	flag.Parse()
	pichumod.ReadSettings()
	defer pichumod.LogFile.Close()
	if !*Fastmode {
		fmt.Println("Program will start in 5 seconds (You can Type Ctrl+C to Cancel)")
		time.Sleep(5 * time.Second)
	}
	//TODO: WebFrame

	//
	for {
		func(host string) {
			// 连接到CQhttp的Websocket服务
			url := url.URL{Scheme: "ws", Host: host, Path: "/"}
			var dailer *websocket.Dialer
			c, _, err := dailer.Dial(url.String(), nil)
			if err != nil {
				fmt.Println(err)
				return
			} else {
				fmt.Println("成功连接到 Websocket 服务器" + pichumod.Config.WSIP)
				pichumod.Connect = c
				go pichumod.SendPrivateMsg("机器人成功启动", pichumod.Config.MasterQQ)
			}

			// 处理收到的所有POST
			for {
				_, message, err := c.ReadMessage()
				if err != nil {
					fmt.Println(err)
				}
				receive, err := jsonParse(message)
				if err != nil {
					fmt.Println(err)
				}
				go msgParse(receive)
			}
		}(pichumod.Config.WSIP)
	}
}

func msgParse(receive map[string]interface{}) {

	switch receive["post_type"] {
	// 消息事件
	case "message":
		switch receive["message_type"] {
		// 私聊信息
		case "private":
			var eventinfo pichumod.MessagePrivate = parsePrivate(receive)
			pichumod.PrintLog(2, fmt.Sprintf("[↓][私聊][%s(%.f)]: %s", eventinfo.Sender.Nickname, eventinfo.Sender.UserID, eventinfo.Message))
			OnPrivateMsg(eventinfo)

		// 群聊信息
		case "group":
			var eventinfo pichumod.MessageGroup = parseGroup(receive)
			pichumod.PrintLog(2, fmt.Sprintf("[↓][群聊(%.f)][%s(%.f)]: %s", eventinfo.GroupID, eventinfo.Sender.Nickname, eventinfo.Sender.UserID, eventinfo.Message))
			OnGroupMsg(eventinfo)

		default:
			pichumod.PrintLog(3, fmt.Sprintf("Cannot Parse 'message' event -> %s", receive))
		}

		// 通知事件
	case "notice":
		switch receive["notice_type"] {
		// 群文件上传
		case "group_upload":
			var eventinfo pichumod.GroupUpload = parseGroupupload(receive)
			pichumod.PrintLog(2, fmt.Sprintf("[N][群文件(%.f)][%.f]: %s", eventinfo.Group_id, eventinfo.User_id, eventinfo.File.Name))
			OnGroupUpload(eventinfo)

			// 群管理员变动
		case "group_admin":
			var eventinfo pichumod.GroupAdmin = parseGroupadmin(receive)
			var x string = "null"
			if eventinfo.Sub_type == "set" {
				x = "+"
			} else {
				x = "-"
			}
			pichumod.PrintLog(2, fmt.Sprintf("[N][群(%.f)管理][%s %.f]", eventinfo.Group_id, x, eventinfo.User_id))
			OnGroupAdmin(eventinfo)

			// 群成员减少
		case "group_decrease":
			var eventinfo pichumod.GroupDecrease = parseGroupdecrease(receive)
			pichumod.PrintLog(2, fmt.Sprintf("[N][成员退群(%.f)][%.f] Type: %s", eventinfo.Group_id, eventinfo.User_id, eventinfo.Sub_type))
			OnGroupDecrease(eventinfo)

			// 群成员增加
		case "group_increase":
			var eventinfo pichumod.GroupIncrease = parseGroupincrease(receive)
			pichumod.PrintLog(2, fmt.Sprintf("[N][成员入群(%.f)][%.f -> %.f] Type: %s", eventinfo.Group_id, eventinfo.Operator_id, eventinfo.User_id, eventinfo.Sub_type))
			OnGroupIncrease(eventinfo)

			// 群禁言
		case "group_ban":
			var eventinfo pichumod.GroupBan = parseGroupban(receive)
			pichumod.PrintLog(2, fmt.Sprintf("[N][群聊(%.f)] %.f 禁言/解禁了 %.f for %.fs", eventinfo.Group_id, eventinfo.Operator_id, eventinfo.User_id, eventinfo.Duration))
			OnGroupBan(eventinfo)

			// 好友添加
		case "friend_add":
			var eventinfo pichumod.FriendAdd = parseFriendAdd(receive)
			pichumod.PrintLog(2, fmt.Sprintf("[N][成功添加好友]%.f", eventinfo.User_id))
			OnFriendAdd(eventinfo)

			// 群消息撤回
		case "group_recall":
			var eventinfo pichumod.GroupRecall = parseGrouprecall(receive)
			pichumod.PrintLog(2, fmt.Sprintf("[N][群聊(%.f)][%.f] 撤回了消息(id): %.f", eventinfo.Group_id, eventinfo.User_id, eventinfo.Message_id))
			OnGroupRecall(eventinfo)

			// 好友消息撤回
		case "friend_recall":
			var eventinfo pichumod.FriendRecall = parseFriendrecall(receive)
			pichumod.PrintLog(2, fmt.Sprintf("[N][私聊][%.f] 撤回了消息(id): %.f", eventinfo.User_id, eventinfo.Message_id))
			OnFriendRecall(eventinfo)

			// 群内戳一戳 群红包运气王 群成员荣誉变更
		case "notify":
			var eventinfo pichumod.Notify = parseNotify(receive)
			pichumod.PrintLog(2, fmt.Sprintf("[N][Notify][Group:%.f] %.f -> %s", eventinfo.Group_id, eventinfo.User_id, eventinfo.Sub_type))
			OnNotify(eventinfo)

		default:
			pichumod.PrintLog(3, fmt.Sprintf("Cannot Parse 'notice' event -> %s", receive))
		}

		// 请求事件
	case "request":
		switch receive["request_type"] {
		// 添加好友申请
		case "friend":
			var eventinfo pichumod.FriendRequest = parseFriendrequest(receive)
			pichumod.PrintLog(2, fmt.Sprintf("[↓][好友申请] %.f 申请加你为好友 -> %s", eventinfo.User_id, eventinfo.Comment))
			OnFriendRequest(eventinfo)

			// 加群邀请
		case "group":
			// pichumod.SetGroupInviteRequest(receive["flag"].(string), true, "") // 自动同意加群
			var eventinfo pichumod.GroupRequest = parseGrouprequest(receive)
			pichumod.PrintLog(2, fmt.Sprintf("[↓][加群/邀请] %.f %s -> %.f(验证信息: %s)", eventinfo.User_id, eventinfo.Sub_type, eventinfo.Group_id, eventinfo.Comment))
			OnGroupRequest(eventinfo)

		default:
			pichumod.PrintLog(3, fmt.Sprintf("Cannot Parse 'request' event -> %s", receive))
		}
		// 元事件
	case "meta_event":
		switch receive["meta_event_type"] {
		// 生命周期
		case "lifecycle":
			var eventinfo pichumod.MetaLifecycle = parseMetalifecycle(receive)
			pichumod.PrintLog(1, fmt.Sprintf("[↓][Lifecycle][%.f] Type: %s", eventinfo.Self_id, eventinfo.Sub_type))
			OnMetaLifecycle(eventinfo)

			// 心跳包
		case "heartbeat":
			var eventinfo pichumod.MetaHeartbeat = parseMetaheartbeat(receive)
			pichumod.PrintLog(1, fmt.Sprintf("[↓][Heartbeat][%.f] Type: %s", eventinfo.Self_id, eventinfo.Status))
			OnMetaHeartbeat(eventinfo)

			// pichumod.PrintLog(1, "Received a heartbeat package.")
		default:
			pichumod.PrintLog(3, fmt.Sprintf("Cannot Parse 'meta_event' event -> %s", receive))
		}
	default:
		// 短事件回调
		if _, ok := receive["echo"]; ok {
			if _, ok := pichumod.RawEvents[receive["echo"].(string)]; ok {
				*pichumod.RawEvents[receive["echo"].(string)].Channel <- receive
			}
		} else {
			pichumod.PrintLog(3, fmt.Sprintf("Got Error Package -> %s", receive))
		}
	}
}

func jsonParse(input []byte) (map[string]interface{}, error) {
	var output map[string]interface{}
	if err := hjson.Unmarshal(input, &output); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return output, nil
}

func parsePrivate(r map[string]interface{}) pichumod.MessagePrivate {
	e := pichumod.MessagePrivate{
		SelfID:     r["self_id"].(float64),
		SubType:    r["sub_type"].(string),
		MessageID:  r["message_id"].(float64),
		UserID:     r["user_id"].(float64),
		Message:    r["message"].(string),
		RawMessage: r["raw_message"].(string),
		Sender: struct {
			UserID   float64
			Nickname string
			Sex      string
			Age      float64
		}{
			UserID:   r["sender"].(map[string]interface{})["user_id"].(float64),
			Nickname: r["sender"].(map[string]interface{})["nickname"].(string),
			Sex:      r["sender"].(map[string]interface{})["sex"].(string),
			Age:      r["sender"].(map[string]interface{})["age"].(float64),
		}}
	return e
}
func parseGroup(r map[string]interface{}) pichumod.MessageGroup {
	e := pichumod.MessageGroup{
		SelfID:     r["self_id"].(float64),
		SubType:    r["sub_type"].(string),
		MessageID:  r["message_id"].(float64),
		GroupID:    r["group_id"].(float64),
		UserID:     r["user_id"].(float64),
		Message:    r["message"].(string),
		RawMessage: r["raw_message"].(string)}
	switch e.SubType {
	case "normal":
		e.Sender = struct {
			UserID   float64
			Nickname string
			Card     string
			Sex      string
			Age      float64
			Area     string
			Level    string
			Role     string
			Title    string
		}{
			UserID:   r["sender"].(map[string]interface{})["user_id"].(float64),
			Nickname: r["sender"].(map[string]interface{})["nickname"].(string),
			Card:     r["sender"].(map[string]interface{})["card"].(string),
			Sex:      r["sender"].(map[string]interface{})["sex"].(string),
			Age:      r["sender"].(map[string]interface{})["age"].(float64),
			Area:     r["sender"].(map[string]interface{})["area"].(string),
			Level:    r["sender"].(map[string]interface{})["level"].(string),
			Role:     r["sender"].(map[string]interface{})["role"].(string),
			Title:    r["sender"].(map[string]interface{})["title"].(string)}
	case "anoymous":
		e.Anonymous = struct {
			Id   float64
			Name string
			Flag string
		}{
			Id:   r["anonymous"].(map[string]interface{})["id"].(float64),
			Name: r["anonymous"].(map[string]interface{})["name"].(string),
			Flag: r["anonymous"].(map[string]interface{})["flag"].(string)}
	}
	return e
}
func parseGroupupload(r map[string]interface{}) pichumod.GroupUpload {
	e := pichumod.GroupUpload{
		Time:     r["time"].(float64),
		Self_id:  r["self_id"].(float64),
		Group_id: r["group_id"].(float64),
		User_id:  r["user_id"].(float64),
		File: struct {
			Id    string
			Name  string
			Size  float64
			Busid float64
		}{
			Id:    r["file"].(map[string]string)["id"],
			Name:  r["file"].(map[string]string)["name"],
			Size:  r["file"].(map[string]float64)["size"],
			Busid: r["file"].(map[string]float64)["busid"],
		}}
	return e
}
func parseGroupadmin(r map[string]interface{}) pichumod.GroupAdmin {
	e := pichumod.GroupAdmin{
		Time:     r["time"].(float64),
		Self_id:  r["self_id"].(float64),
		Sub_type: r["sub_type"].(string),
		Group_id: r["group_id"].(float64),
		User_id:  r["user_id"].(float64),
	}
	return e
}
func parseGroupdecrease(r map[string]interface{}) pichumod.GroupDecrease {
	e := pichumod.GroupDecrease{
		Time:        r["time"].(float64),
		Self_id:     r["self_id"].(float64),
		Sub_type:    r["sub_type"].(string),
		Group_id:    r["group_id"].(float64),
		Operator_id: r["operator_id"].(float64),
		User_id:     r["user_id"].(float64),
	}
	return e
}
func parseGroupincrease(r map[string]interface{}) pichumod.GroupIncrease {
	e := pichumod.GroupIncrease{
		Time:        r["time"].(float64),
		Self_id:     r["self_id"].(float64),
		Sub_type:    r["sub_type"].(string),
		Group_id:    r["group_id"].(float64),
		Operator_id: r["operator_id"].(float64),
		User_id:     r["user_id"].(float64),
	}
	return e
}
func parseGroupban(r map[string]interface{}) pichumod.GroupBan {
	e := pichumod.GroupBan{
		Time:        r["time"].(float64),
		Self_id:     r["self_id"].(float64),
		Sub_type:    r["sub_type"].(string),
		Group_id:    r["group_id"].(float64),
		Operator_id: r["operator_id"].(float64),
		User_id:     r["user_id"].(float64),
		Duration:    r["duration"].(float64),
	}
	return e
}
func parseFriendAdd(r map[string]interface{}) pichumod.FriendAdd {
	e := pichumod.FriendAdd{
		Time:    r["time"].(float64),
		Self_id: r["self_id"].(float64),
		User_id: r["user_id"].(float64),
	}
	return e
}

func parseGrouprecall(r map[string]interface{}) pichumod.GroupRecall {
	e := pichumod.GroupRecall{
		Time:        r["time"].(float64),
		Self_id:     r["self_id"].(float64),
		Group_id:    r["group_id"].(float64),
		User_id:     r["user_id"].(float64),
		Operator_id: r["operator_id"].(float64),
		Message_id:  r["message_id"].(float64),
	}
	return e
}
func parseFriendrecall(r map[string]interface{}) pichumod.FriendRecall {
	e := pichumod.FriendRecall{
		Time:       r["time"].(float64),
		Self_id:    r["self_id"].(float64),
		User_id:    r["user_id"].(float64),
		Message_id: r["message_id"].(float64),
	}
	return e
}
func parseNotify(r map[string]interface{}) pichumod.Notify {
	e := pichumod.Notify{
		Time:     r["time"].(float64),
		Self_id:  r["self_id"].(float64),
		Sub_type: r["sub_type"].(string),
		Group_id: r["group_id"].(float64),
		User_id:  r["user_id"].(float64),
	}
	if e.Sub_type == "honor" {
		e.Honor_type = r["honor_type"].(string)
	} else {
		e.Target_id = r["target_id"].(float64)
	}
	return e
}
func parseFriendrequest(r map[string]interface{}) pichumod.FriendRequest {
	e := pichumod.FriendRequest{
		Time:    r["time"].(float64),
		Self_id: r["self_id"].(float64),
		User_id: r["user_id"].(float64),
		Comment: r["comment"].(string),
		Flag:    r["flag"].(string),
	}
	return e
}
func parseGrouprequest(r map[string]interface{}) pichumod.GroupRequest {
	e := pichumod.GroupRequest{
		Time:     r["time"].(float64),
		Self_id:  r["self_id"].(float64),
		Sub_type: r["sub_type"].(string),
		Group_id: r["group_id"].(float64),
		User_id:  r["user_id"].(float64),
		Comment:  r["comment"].(string),
		Flag:     r["flag"].(string),
	}
	return e
}

func parseMetalifecycle(r map[string]interface{}) pichumod.MetaLifecycle {
	e := pichumod.MetaLifecycle{
		Time:     r["time"].(float64),
		Self_id:  r["self_id"].(float64),
		Sub_type: r["sub_type"].(string),
	}
	return e
}
func parseMetaheartbeat(r map[string]interface{}) pichumod.MetaHeartbeat {
	e := pichumod.MetaHeartbeat{
		Time:     r["time"].(float64),
		Self_id:  r["self_id"].(float64),
		Status:   r["status"],
		Interval: r["interval"].(float64),
	}
	return e
}
