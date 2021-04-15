// Bot的主要功能将更新于这个文件
package pichumod

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//! 以下均为Example 可以全部删掉 也可以提供参考

func PrivateParse(eventinfo MessagePrivate) { // 建立事件 解析内容
	if match, _ := regexp.MatchString(`(?m)^[!|！].*`, eventinfo.Message); match { // 判断消息内容是否为Command
		var command []string = strings.Fields(eventinfo.Message[1:])
		if len(command) > 0 {

			// fmt.Println(command[0])
			switch command[0] {
			case "test", "helloworld":
				var m string
				if time.Now().Format("PM") == "PM" {
					m = "下午"
				} else {
					m = "上午"
				}
				nowtime := time.Now().Format(fmt.Sprintf("2006年1月2日%s3:04", m))
				SendPrivateMsg(fmt.Sprintf("这是一条测试消息！\n你好，我的朋友%s！\n现在是：%s", eventinfo.Sender.Nickname, nowtime), eventinfo.UserID)
			}
		}
	}
}

func GroupParse(eventinfo MessageGroup) {
	if match, _ := regexp.MatchString(`(?m)^[!！].*`, eventinfo.Message); match {
		var command []string = strings.Fields(eventinfo.Message[1:])
		if len(command) > 0 {
			// commandParse
			switch command[0] {
			// 长事件案例
			case "复读机":
				// 新建长事件方法
				ch, eid := NewEvent(eventinfo.Sender.UserID, eventinfo.GroupID, "fdj")
				defer close(*ch)

				SendGroupMsg("我，就是人类的本质", eventinfo.GroupID)

				for {
					r := <-*ch // 获取传入channel的值(string)
					if r == "关闭复读机" {
						SendGroupMsg("复读机已关闭", eventinfo.GroupID)
						break
					} else {
						SendGroupMsg(r, eventinfo.GroupID)
					}
				}
				delete(Events, eid) // 结束事件
			}
		}
	} else {
		// 判断非命令消息
		for _, value := range Events { // 匹配事件池内所有事件
			if value.GroupID == eventinfo.GroupID { // 匹配事件发生的群号
				switch value.Eventtype {
				case "fdj":
					*value.Channel <- eventinfo.Message // 传入数据(string)
				}
			}
		}
	}
}

func NewEvent(userid float64, groupid float64, etype string) (*chan string, string) {
	var eventid string = strconv.FormatInt(time.Now().UnixNano(), 10) // 事件ID 以时间戳定义
	ch := make(chan string)
	Events[eventid] = Event{Channel: &ch, UserID: userid, GroupID: groupid, Eventtype: etype}
	fmt.Println(Events)
	return &ch, eventid
}
