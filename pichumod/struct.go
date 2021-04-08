// 存放所有结构体
package pichumod

var RawEvents = make(map[string]RawEvent) // 短事件ID列表
var Events = make(map[string]Event)       // 事件ID列表

type RawEvent struct {
	Channel *chan map[string]interface{}
}

type Event struct {
	UserID    float64
	GroupID   float64
	Channel   *chan string
	Eventtype string
}

type ConfigStu struct {
	Logmode  string  // settings.hjson - logmode
	Loglvl   int     // 由Logmod产生
	WSIP     string  // settings.hjson - ws-ip
	MasterQQ float64 // settings.hjson - master-qq
}

// 私聊消息事件所使用的结构体
type MessagePrivate struct {
	SelfID     float64       // 收到事件的机器人 QQ 号
	SubType    string        // 消息子类型(friend,group,other)
	MessageID  float64       // 消息id
	UserID     float64       // 发送者 QQ 号
	Message    string        // 消息内容
	RawMessage string        // 原始消息内容
	Sender     SenderPrivate // 发送人信息
}

type SenderPrivate struct {
	UserID   float64 // 发送者 QQ 号
	Nickname string  // 昵称
	Sex      string  // 性别，male 或 female 或 unknown
	Age      float64 // 年龄
}

// 群聊消息事件所使用的结构体
type MessageGroup struct {
	SelfID    float64 // 收到事件的机器人 QQ 号
	SubType   string  // 消息子类型 正常消息是 normal 匿名消息是 anonymous 系统提示是 notice
	MessageID float64 // 消息 ID
	GroupID   float64 // 群号
	UserID    float64 // 发送者 QQ 号
	// 匿名消息暂不做处理
	Message    string      // 消息内容
	RawMessage string      //	原始消息内容
	Sender     SenderGroup // 发送人信息
}

type SenderGroup struct {
	UserID   float64 // 发送者 QQ 号
	Nickname string  // 昵称
	Card     string  // 群名片／备注
	Sex      string  // 性别，male 或 female 或 unknown
	Age      float64 // 年龄
	Area     string  // 地区
	Level    string  // 成员等级
	Role     string  // 角色，owner 或 admin 或 member
	Title    string  // 专属头衔
}
