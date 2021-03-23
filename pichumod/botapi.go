package pichumod

import (
	"errors"
	"fmt"

	"github.com/gorilla/websocket"
)

// ws发包
func sendwspack(message string) error {
	err := Connect.WriteMessage(websocket.TextMessage, []byte(message))
	return err
}

// 发送API
func apiSend(apiType string, params string) error {
	err := sendwspack(fmt.Sprintf(
		`{"action": "%s", "params": %s}`,
		apiType, params))
	return err
}

// onebot API: https://git.io/Jmy1B
// cqhttp API: https://github.com/Mrs4s/go-cqhttp/blob/master/docs/cqhttp.md#api-1

// 发送私聊信息
// message - 消息内容 自动解析CQ码
// user_id - 对方QQ号
// funcR:error wsR:message_id
func SendPrivateMsg(message string, user_id int) error {
	err := apiSend("send_private_msg", fmt.Sprintf(`{"user_id": %d, "message": "%s"}`, user_id, message))
	PrintLog(2, fmt.Sprintf("[↑][私聊][%d]: %s", user_id, message))
	return err
}

// 发送群聊消息
// message  - 要发送的内容
// group_id - 群号
// funcR:error wsR:message_id
func SendGroupMsg(message string, group_id int) error {
	err := apiSend("send_group_msg", fmt.Sprintf(`{"group_id": %d, "message": "%s"}`, group_id, message))
	PrintLog(2, fmt.Sprintf("[↑][群聊][%d] %s", group_id, message))
	return err
}

// 发送消息
// msgtype - 消息类型 group/private
// message - 消息内容
// toid    - 群号/QQ号
// 本条API并不是 Onebot/CQhttp 原生API
// funcR:error wsR:message_id
func SendMsg(msgtype string, message string, toid int) error {
	var err error = nil
	switch msgtype {
	case "group":
		err = SendGroupMsg(message, toid)
	case "private":
		err = SendPrivateMsg(message, toid)
	default:
		return errors.New("an error using function pichumod.SendMsg: msgtype should be group or private")
	}
	return err
}

// 撤回消息
// message_id - 消息id 发出时的返回值
// funcR:error wsR:nil
func DeleteMsg(message_id int32) error {
	err := apiSend("delete_msg", fmt.Sprintf(`{"message_id": %d}`, message_id))
	return err
}

// 获取消息
// message_id - 获取消息
// 获取到的消息会被解析为正常收到的消息
// 后续建议使用echo区分这些消息
// funcR:error wsR:NormalMessage+echo(wip)
func GetMsg(message_id int32) error {
	err := apiSend("get_msg", fmt.Sprintf(`{"message_id": %d}`, message_id))
	return err
}

// 获取合并转发消息
// id - 合并转发 ID
// funcR:error wsR:消息内容，使用 消息的数组格式 表示，数组中的消息段全部为 node 消息段
func GetForwardMsg(id string) error {
	err := apiSend("get_forward_msg", fmt.Sprintf(`{"id":"%s"}`, id))
	return err
}

// 发送好友赞
// user_id - 对方QQ号
// times - 点赞次数(每个好友每天最多 10 次)
// funcR:err wsR: nil
func SendLike(user_id int, times int) error {
	err := apiSend("send_like", fmt.Sprintf(`{"user_id": %d, "times": %d}`, user_id, times))
	return err
}

// 群组踢人
// group_id - 群号
// user_id - 要踢的 QQ 号
// reject_add_request - 是否拒绝再次入群
// funcR:err wsR: nil
func SetGroupKick(group_id int, user_id int, reject_add_request bool) error {
	err := apiSend("set_group_kick", fmt.Sprintf(`{"group_id": %d, "user_id": %d, "reject_add_request": %s}`, group_id, user_id, BoolToStr(reject_add_request)))
	return err
}

// 群组单人禁言
// group_id - 群号
// user_id - 要禁言的QQ号
// duration - 禁言时长(s) 0表示取消禁言
// funcR:err wsR:nil
func SetGroupBan(group_id int, user_id int, duration int) error {
	err := apiSend("set_group_ban", fmt.Sprintf(`{"group_id": %d, "user_id": %d, "duration": "%d"}`, group_id, user_id, duration))
	return err
}

// 群组匿名用户禁言
// group_id - 群号
// anymous_flag - 匿名用户的 flag（需从群消息上报的数据中获得）
// duration - 禁言时长(s) 0表示取消禁言
// funcR:err wsR: nil
func SetGroupAnonymousBan(group_id int, anymous_flag string, duration int) error {
	err := apiSend("set_group_ban", fmt.Sprintf(`{"group_id": %d, "anymous_flag": "%s", "duration": "%d"}`, group_id, anymous_flag, duration))
	return err
}

// 群全员禁言
// group_id 群号
// enable 是否禁言
// funcR:err wsR: nil
func SetGroupWholeBan(group_id int, enable bool) error {
	err := apiSend("set_group_kick", fmt.Sprintf(`{"group_id": %d, "enable": %s}`, group_id, BoolToStr(enable)))
	return err
}

// 群组设置管理员(需要机器人为群主)
// group_id 群号
// user_id QQ号
// enable true 为设置，false 为取消
// funcR:err wsR: nil
func SetGroupAdmin(group_id int, user_id int, enable bool) error {
	err := apiSend("set_group_admin", fmt.Sprintf(`{"group_id": %d, "user_id": %d , "enable": %s}`, group_id, user_id, BoolToStr(enable)))
	return err
}

// 群组匿名
// group_id 群号
// enable 是否允许匿名聊天
// funcR:err wsR: nil
func SetGroupAnonymous(group_id int, enable bool) error {
	err := apiSend("set_group_anonymous", fmt.Sprintf(`{"group_id": %d, "enable": %s}`, group_id, BoolToStr(enable)))
	return err
}

// 设置群名片
// group_id 群号
// user_id 成员QQ
// card 空字符串表示删除群名片
// funcR:err wsR: nil
func SetGroupCard(group_id int, user_id int, card string) error {
	err := apiSend("set_group_card", fmt.Sprintf(`{"group_id": %d, "user_id": %d , "card": "%s"}`, group_id, user_id, card))
	return err
}

// 设置群名
// group_id 群号
// group_name 新群名
// funcR:err wsR: nil
func SetGroupName(group_id int, group_name string) error {
	err := apiSend("set_group_name", fmt.Sprintf(`{"group_id": %d, "group_name": "%s"}`, group_id, group_name))
	return err
}

// 退群
// group_id 群号
// is_dismiss 是否解散，如果登录号是群主，则仅在此项为 true 时能够解散
func SetGroupLeave(group_id int, is_dismiss bool) error {
	err := apiSend("set_group_leave", fmt.Sprintf(`{"group_id": %d, "is_dismiss": %s}`, group_id, BoolToStr(is_dismiss)))
	return err
}

// 设置群组专属头衔
// group_id 群号
// user_id 成员QQ
// special_title 空字符串表示删除专属头衔
// funcR:err wsR: nil
func SetGroupSpecialTitle(group_id int, user_id int, special_title string) error {
	err := apiSend("set_group_special_title", fmt.Sprintf(`{"group_id": %d, "user_id": %d , "special_title": "%s"}`, group_id, user_id, special_title))
	return err
}

// 处理加好友请求
// flag 加好友请求的 flag（需从上报的数据中获得）
// approve 是否同意请求
// funcR:err wsR: nil
func SetFriendAddRequest(flag string, approve bool) error {
	err := apiSend("set_friend_add_request", fmt.Sprintf(`{"flag": "%s", "approve": %s}`, flag, BoolToStr(approve)))
	return err
}

// 处理加群请求
// flag 加群请求的 flag
// approve 是否同意请求
// reason 拒绝理由(只有在拒绝时有效)
func SetGroupAddRequest(flag string, approve bool, reason string) error {
	err := apiSend("set_group_add_request", fmt.Sprintf(`{"flag": "%s", "sub_type": "add", "approve": %s, "reason": "%s"}`, flag, BoolToStr(approve), reason))
	return err
}

// 处理加群邀请
// flag 加群邀请的 flag
// approve 是否同意邀请
// reason 拒绝理由(只有在拒绝时有效)
func SetGroupInviteRequest(flag string, approve bool, reason string) error {
	err := apiSend("set_group_add_request", fmt.Sprintf(`{"flag": "%s", "sub_type": "invite", "approve": %s, "reason": "%s"}`, flag, BoolToStr(approve), reason))
	return err
}

// 获取登录号信息
// funcR: 错误处理 wsR: user_id QQ号, nickname QQ昵称
func GetLoginInfo() error {
	err := apiSend("get_login_info", "")
	return err
}

// 以下为 CQhttp 的API

// 获取图片信息
// file string
// get_image
// wsR: size			int32		图片源文件大小
//			filename	string	图片文件原名
//			url				string	图片下载地址
func GetImage(file string) error {
	err := apiSend("get_image", fmt.Sprintf(`{"file": "%s"}`, file))
	return err
}

// 图片OCR
// image file - string
// wsR: texts			TextDetection[]	OCR结果
//      -	text				string	文本
// 			-	confidence	int32		置信度
// 			-	coordinates	vector2	坐标
// 			language	string					语言
func OCRImage(image string) error {
	err := apiSend("ocr_image", fmt.Sprintf(`{"image": "%s"}`, image))
	return err
}
