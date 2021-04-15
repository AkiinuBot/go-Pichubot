//! 整合大部分常用api于函数内
package pichumod

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

// ws发包
func sendwspack(message string) error {
	err := Connect.WriteMessage(websocket.TextMessage, []byte(message))
	return err
}

// 发送API
func apiSend(apiType string, params string) (map[string]interface{}, error) {
	eventid := strconv.FormatInt(time.Now().UnixNano(), 10)
	ch := make(chan map[string]interface{})
	defer close(ch)

	chinfo := RawEvent{Channel: &ch}

	RawEvents[eventid] = chinfo

	PrintLog(1, fmt.Sprintf("[↑][EID:%s][Type:%s]S:%s", eventid, apiType, params))
	err := sendwspack(fmt.Sprintf(`{"action": "%s", "params": %s, "echo": "%s"}`, apiType, params, eventid))
	var receive map[string]interface{}
	if err == nil {
		select {
		case receive = <-ch:
			// fmt.Println(receive)
			PrintLog(1, fmt.Sprintf("[↓][EID:%s][Type:%s]R:%s", eventid, apiType, receive))
		case <-time.After(5 * time.Second):
			PrintLog(3, fmt.Sprintf("[↓][EID:%s][Type:%s]Timeout", eventid, apiType))
			err = errors.New("timout in func apiSend")
		}
	} else {
		fmt.Println(err)
	}
	delete(RawEvents, eventid)
	return receive, err
}

// onebot API: https://git.io/Jmy1B
// cqhttp API: https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#api-1

// 发送私聊信息
// message - 消息内容 自动解析CQ码
// user_id - 对方QQ号
// return message_id error
func SendPrivateMsg(message string, user_id float64) (map[string]interface{}, error) {
	res, err := apiSend("send_private_msg", fmt.Sprintf(`{"user_id": %.f, "message": "%s"}`, user_id, message))
	// fmt.Println(res)
	if err != nil {
		return nil, err
	}
	if res["status"].(string) == "ok" {
		PrintLog(2, fmt.Sprintf("[↑][私聊][%.f]: %s", user_id, message))
	} else {
		PrintLog(3, fmt.Sprintf("[↑][发送失败][私聊][%.f]: %s", user_id, message))
	}
	return res, err
}

// 发送群聊消息
// message  - 要发送的内容
// group_id - 群号
// return message_id error
func SendGroupMsg(message string, group_id float64) (map[string]interface{}, error) {
	res, err := apiSend("send_group_msg", fmt.Sprintf(`{"group_id": %.f, "message": "%s"}`, group_id, message))
	if err != nil {
		return nil, err
	}
	if res["status"].(string) == "ok" {
		PrintLog(2, fmt.Sprintf("[↑][群聊][%.f]: %s", group_id, message))
	} else {
		PrintLog(3, fmt.Sprintf("[↑][发送失败][群聊][%.f]: %s", group_id, message))
	}
	// PrintLog(2, fmt.Sprintf("[↑][群聊][%.f] %s", group_id, message))
	return res, err
}

// 发送消息
// msgtype - 消息类型 group/private
// message - 消息内容
// toid    - 群号/QQ号
// 本条API并不是 Onebot/CQhttp 原生API
// return message_id error
func SendMsg(msgtype string, message string, toid float64) (map[string]interface{}, error) {
	var err error = nil
	var res map[string]interface{}
	switch msgtype {
	case "group":
		res, err = SendGroupMsg(message, toid)
	case "private":
		res, err = SendPrivateMsg(message, toid)
	default:
		return nil, errors.New("an error using function pichumod.SendMsg: msgtype should be group or private")
	}
	return res, err
}

// 撤回消息
// message_id - 消息id 发出时的返回值
// return error
func DeleteMsg(message_id int32) error {
	_, err := apiSend("delete_msg", fmt.Sprintf(`{"message_id": %d}`, message_id))
	return err
}

// 获取消息
// message_id - 获取消息
// return {time message_type message_id real_id sender message} error
func GetMsg(message_id int32) (map[string]interface{}, error) {
	res, err := apiSend("get_msg", fmt.Sprintf(`{"message_id": %d}`, message_id))
	return res, err
}

// 获取合并转发消息
// id - 合并转发 ID
// return message err
func GetForwardMsg(id string) (map[string]interface{}, error) {
	res, err := apiSend("get_forward_msg", fmt.Sprintf(`{"id":"%s"}`, id))
	return res, err
}

// 发送好友赞
// user_id - 对方QQ号
// times - 点赞次数(每个好友每天最多 10 次)
// return err
func SendLike(user_id int, times int) error {
	_, err := apiSend("send_like", fmt.Sprintf(`{"user_id": %d, "times": %d}`, user_id, times))
	return err
}

// 群组踢人
// group_id - 群号
// user_id - 要踢的 QQ 号
// reject_add_request - 是否拒绝再次入群
// return err
func SetGroupKick(group_id int, user_id int, reject_add_request bool) error {
	_, err := apiSend("set_group_kick", fmt.Sprintf(`{"group_id": %d, "user_id": %d, "reject_add_request": %s}`, group_id, user_id, BoolToStr(reject_add_request)))
	return err
}

// 群组单人禁言
// group_id - 群号
// user_id - 要禁言的QQ号
// duration - 禁言时长(s) 0表示取消禁言
// return err
func SetGroupBan(group_id int, user_id int, duration int) error {
	_, err := apiSend("set_group_ban", fmt.Sprintf(`{"group_id": %d, "user_id": %d, "duration": "%d"}`, group_id, user_id, duration))
	return err
}

// 群组匿名用户禁言
// group_id - 群号
// anymous_flag - 匿名用户的 flag（需从群消息上报的数据中获得）
// duration - 禁言时长(s) 0表示取消禁言
// return err
func SetGroupAnonymousBan(group_id int, anymous_flag string, duration int) error {
	_, err := apiSend("set_group_ban", fmt.Sprintf(`{"group_id": %d, "anymous_flag": "%s", "duration": "%d"}`, group_id, anymous_flag, duration))
	return err
}

// 群全员禁言
// group_id 群号
// enable 是否禁言
// return err
func SetGroupWholeBan(group_id int, enable bool) error {
	_, err := apiSend("set_group_kick", fmt.Sprintf(`{"group_id": %d, "enable": %s}`, group_id, BoolToStr(enable)))
	return err
}

// 群组设置管理员(需要机器人为群主)
// group_id 群号
// user_id QQ号
// enable true 为设置，false 为取消
// return err
func SetGroupAdmin(group_id int, user_id int, enable bool) error {
	_, err := apiSend("set_group_admin", fmt.Sprintf(`{"group_id": %d, "user_id": %d , "enable": %s}`, group_id, user_id, BoolToStr(enable)))
	return err
}

// 群组匿名
// group_id 群号
// enable 是否允许匿名聊天
// return err
func SetGroupAnonymous(group_id int, enable bool) error {
	_, err := apiSend("set_group_anonymous", fmt.Sprintf(`{"group_id": %d, "enable": %s}`, group_id, BoolToStr(enable)))
	return err
}

// 设置群名片
// group_id 群号
// user_id 成员QQ
// card 空字符串表示删除群名片
// return err
func SetGroupCard(group_id int, user_id int, card string) error {
	_, err := apiSend("set_group_card", fmt.Sprintf(`{"group_id": %d, "user_id": %d , "card": "%s"}`, group_id, user_id, card))
	return err
}

// 设置群名
// group_id 群号
// group_name 新群名
// return err
func SetGroupName(group_id int, group_name string) error {
	_, err := apiSend("set_group_name", fmt.Sprintf(`{"group_id": %d, "group_name": "%s"}`, group_id, group_name))
	return err
}

// 退群
// group_id 群号
// is_dismiss 是否解散，如果登录号是群主，则仅在此项为 true 时能够解散
// return err
func SetGroupLeave(group_id int, is_dismiss bool) error {
	_, err := apiSend("set_group_leave", fmt.Sprintf(`{"group_id": %d, "is_dismiss": %s}`, group_id, BoolToStr(is_dismiss)))
	return err
}

// 设置群组专属头衔
// group_id 群号
// user_id 成员QQ
// special_title 空字符串表示删除专属头衔
// return err
func SetGroupSpecialTitle(group_id int, user_id int, special_title string) error {
	_, err := apiSend("set_group_special_title", fmt.Sprintf(`{"group_id": %d, "user_id": %d , "special_title": "%s"}`, group_id, user_id, special_title))
	return err
}

// 处理加好友请求
// flag 加好友请求的 flag（需从上报的数据中获得）
// approve 是否同意请求
// return err
func SetFriendAddRequest(flag string, approve bool) error {
	_, err := apiSend("set_friend_add_request", fmt.Sprintf(`{"flag": "%s", "approve": %s}`, flag, BoolToStr(approve)))
	return err
}

// 处理加群请求
// flag 加群请求的 flag
// approve 是否同意请求
// reason 拒绝理由(只有在拒绝时有效)
// return err
func SetGroupAddRequest(flag string, approve bool, reason string) error {
	_, err := apiSend("set_group_add_request", fmt.Sprintf(`{"flag": "%s", "sub_type": "add", "approve": %s, "reason": "%s"}`, flag, BoolToStr(approve), reason))
	return err
}

// 处理加群邀请
// flag 加群邀请的 flag
// approve 是否同意邀请
// reason 拒绝理由(只有在拒绝时有效)
// return err
func SetGroupInviteRequest(flag string, approve bool, reason string) error {
	_, err := apiSend("set_group_add_request", fmt.Sprintf(`{"flag": "%s", "sub_type": "invite", "approve": %s, "reason": "%s"}`, flag, BoolToStr(approve), reason))
	return err
}

// 获取登录号信息
// return {user_id nickname} err
func GetLoginInfo() (map[string]interface{}, error) {
	res, err := apiSend("get_login_info", "")
	return res, err
}

// 以下为 CQhttp 的API

// 获取图片信息
// file string
// get_image
// return size			int32		图片源文件大小
//				filename	string	图片文件原名
//				url				string	图片下载地址
// 				err
func GetImage(file string) (map[string]interface{}, error) {
	res, err := apiSend("get_image", fmt.Sprintf(`{"file": "%s"}`, file))
	return res, err
}

// 图片OCR
// image file - string
// return texts			TextDetection[]	OCR结果
//      	-	text				string	文本
// 				-	confidence	int32		置信度
// 				-	coordinates	vector2	坐标
// 			  language	string					语言
//				err
func OCRImage(image string) (map[string]interface{}, error) {
	res, err := apiSend("ocr_image", fmt.Sprintf(`{"image": "%s"}`, image))
	return res, err
}
