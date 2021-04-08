package main

import (
	"flag"
	"fmt"
	"net/url"
	"time"

	"github.com/0ojixueseno0/go-Pichubot-base/pichumod"

	"github.com/gorilla/websocket"
	"github.com/hjson/hjson-go"
)

// var host = flag.String("ip", "localhost:6700", "CQbot program's host (with port)")
var Fastmode = flag.Bool("faststart", false, "是否快速启动")

func main() {
	flag.Parse() // 解析fastmode
	pichumod.ReadSettings()
	defer pichumod.LogFile.Close()
	// readSettings() // 读取设置 已经挪到了Pichumod
	// linkLog() 已经挪到了Pichumod
	if !*Fastmode {
		fmt.Println("Program will start in 5 seconds (You can Type Ctrl+C to Cancel)")
		time.Sleep(5 * time.Second)
	}
	for {
		core(pichumod.Config.WSIP)
	}
}

func core(host string) {
	//连接到CQhttp的Websocket服务
	url := url.URL{Scheme: "ws", Host: host, Path: "/"}
	var dailer *websocket.Dialer
	c, _, err := dailer.Dial(url.String(), nil)
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("成功连接到 Websocket 服务器" + pichumod.Config.WSIP)
		// tools.connect = c
		pichumod.Connect = c
		go pichumod.SendPrivateMsg("机器人成功启动", 2773173293)
		// time.Sleep(1 * time.Second)
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
		// time := time.Now().Format(time.Kitchen)

		go func() {

			switch receive["post_type"] {
			// 消息事件
			case "message":
				sender := receive["sender"].(map[string]interface{})
				switch receive["message_type"] {
				// 私聊信息
				case "private":
					pichumod.PrintLog(2, fmt.Sprintf("[↓][私聊][%s(%.f)]: %s", sender["nickname"], sender["user_id"].(float64), receive["message"]))
					var eventinfo pichumod.MessagePrivate = parsePrivate(receive)
					// fmt.Println(eventinfo)
					pichumod.PrivateParse(eventinfo)

					// 群聊信息
				case "group":
					pichumod.PrintLog(2, fmt.Sprintf("[↓][群聊(%.f)][%s(%.f)]: %s", receive["group_id"].(float64), sender["nickname"], sender["user_id"].(float64), receive["message"]))
					var eventinfo pichumod.MessageGroup = parseGroup(receive)
					pichumod.GroupParse(eventinfo)

				default:
					pichumod.PrintLog(3, fmt.Sprintf("Cannot Parse 'message' event -> %s", receive))
				}

				// 通知事件
			case "notice":
				switch receive["notice_type"] {
				// 群文件上传
				case "group_upload":

					// 群管理员变动
				case "group_admin":

					// 群成员减少
				case "group_decrease":

					// 群成员增加
				case "group_increase":

					// 群禁言
				case "group_ban":

					// 好友添加
				case "friend_add":

					// 群消息撤回
				case "group_recall":

					// 好友消息撤回
				case "friend_recall":

					// 群内戳一戳 群红包运气王 群成员荣誉变更
				case "notify":

				default:
					pichumod.PrintLog(3, fmt.Sprintf("Cannot Parse 'notice' event -> %s", receive))
				}

				// 请求事件
			case "request":
				switch receive["request_type"] {
				// 添加好友申请
				case "friend":

					// 加群邀请
				case "group":
					// fmt.Println(receive)
					pichumod.SetGroupInviteRequest(receive["flag"].(string), true, "") // 自动同意加群
				default:
					pichumod.PrintLog(3, fmt.Sprintf("Cannot Parse 'request' event -> %s", receive))
				}
				// 元事件
			case "meta_event":
				switch receive["meta_event_type"] {
				// 生命周期
				case "lifecycle":

					// 心跳包
				case "heartbeat":
					pichumod.PrintLog(1, "Received a heartbeat package.")
				default:
					pichumod.PrintLog(3, fmt.Sprintf("Cannot Parse 'meta_event' event -> %s", receive))
				}
			default:
				// 回调
				// fmt.Println(receive)
				if _, ok := receive["echo"]; ok {
					if _, ok := pichumod.RawEvents[receive["echo"].(string)]; ok {
						*pichumod.RawEvents[receive["echo"].(string)].Channel <- receive
					}
				} else {
					pichumod.PrintLog(3, fmt.Sprintf("Got Error Package -> %s", receive))
				}
			}
		}()
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

func parsePrivate(receive map[string]interface{}) pichumod.MessagePrivate {
	var eventinfo pichumod.MessagePrivate
	eventinfo.SelfID = receive["self_id"].(float64)
	eventinfo.SubType = receive["sub_type"].(string)
	eventinfo.MessageID = receive["message_id"].(float64)
	eventinfo.UserID = receive["user_id"].(float64)
	eventinfo.Message = receive["message"].(string)
	eventinfo.RawMessage = receive["raw_message"].(string)
	eventinfo.Sender.UserID = receive["sender"].(map[string]interface{})["user_id"].(float64)
	eventinfo.Sender.Nickname = receive["sender"].(map[string]interface{})["nickname"].(string)
	eventinfo.Sender.Sex = receive["sender"].(map[string]interface{})["sex"].(string)
	eventinfo.Sender.Age = receive["sender"].(map[string]interface{})["age"].(float64)
	return eventinfo
}
func parseGroup(receive map[string]interface{}) pichumod.MessageGroup {
	var eventinfo pichumod.MessageGroup
	eventinfo.SelfID = receive["self_id"].(float64)
	eventinfo.SubType = receive["sub_type"].(string)
	eventinfo.MessageID = receive["message_id"].(float64)
	eventinfo.GroupID = receive["group_id"].(float64)
	eventinfo.UserID = receive["user_id"].(float64)
	eventinfo.Message = receive["message"].(string)
	eventinfo.RawMessage = receive["raw_message"].(string)
	eventinfo.Sender.UserID = receive["sender"].(map[string]interface{})["user_id"].(float64)
	eventinfo.Sender.Nickname = receive["sender"].(map[string]interface{})["nickname"].(string)
	eventinfo.Sender.Card = receive["sender"].(map[string]interface{})["card"].(string)
	eventinfo.Sender.Sex = receive["sender"].(map[string]interface{})["sex"].(string)
	eventinfo.Sender.Age = receive["sender"].(map[string]interface{})["age"].(float64)
	eventinfo.Sender.Area = receive["sender"].(map[string]interface{})["area"].(string)
	eventinfo.Sender.Level = receive["sender"].(map[string]interface{})["level"].(string)
	eventinfo.Sender.Role = receive["sender"].(map[string]interface{})["role"].(string)
	eventinfo.Sender.Title = receive["sender"].(map[string]interface{})["title"].(string)
	return eventinfo
}

// func Inlist(str string, list []string) bool {
// 	for _, item := range list {
// 		if item == str {
// 			return true
// 		}
// 	}
// 	return false
// }
