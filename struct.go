// 存放所有结构体
package Pichubot

type Config struct {
	Loglvl   int
	Host     string
	MasterQQ string
	Path     string
}

type Bot struct {
	Config Config
}

//* Message Event Part

// 私聊消息事件所使用的结构体
type MessagePrivate struct {
	SelfID     int64  // 收到事件的机器人 QQ 号
	SubType    string // 消息子类型(friend,group,other)
	MessageID  int64  // 消息id
	UserID     int64  // 发送者 QQ 号
	Message    string // 消息内容
	RawMessage string // 原始消息内容
	Sender     struct {
		UserID   int64  // 发送者 QQ 号
		Nickname string // 昵称
		Sex      string // 性别，male 或 female 或 unknown
		Age      int64  // 年龄
	} // 发送人信息
}

// 群聊消息事件所使用的结构体
type MessageGroup struct {
	SelfID    int64  // 收到事件的机器人 QQ 号
	SubType   string // 消息子类型 正常消息是 normal 匿名消息是 anonymous 系统提示是 notice
	MessageID int64  // 消息 ID
	GroupID   int64  // 群号
	UserID    int64  // 发送者 QQ 号
	Anonymous struct {
		Id   int64  // 匿名用户 ID
		Name string // 匿名用户名称
		Flag string // 匿名用户 flag，在调用禁言 API 时需要传入
	} //	匿名信息，如果不是匿名消息则为 null
	Message    string // 消息内容
	RawMessage string //	原始消息内容
	Sender     struct {
		UserID   int64  // 发送者 QQ 号
		Nickname string // 昵称
		Card     string // 群名片／备注
		Sex      string // 性别，male 或 female 或 unknown
		Age      int64  // 年龄
		Area     string // 地区
		Level    string // 成员等级
		Role     string // 角色，owner 或 admin 或 member
		Title    string // 专属头衔
	} // 发送人信息
}

//* Notice Event Part

type GroupUpload struct {
	Time     int64 //	事件发生的时间戳
	Self_id  int64 //	收到事件的机器人 QQ 号
	Group_id int64 //	群号
	User_id  int64 //	发送者 QQ 号
	File     struct {
		Id    string //	文件 ID
		Name  string //	文件名
		Size  int64  //	文件大小（字节数）
		Busid int64  //	busid（目前不清楚有什么作用）
	} //	文件信息
}

type GroupAdmin struct {
	Time     int64  // 事件发生的时间戳
	Self_id  int64  // 收到事件的机器人 QQ 号
	Sub_type string // set|unset	事件子类型，分别表示设置和取消管理员
	Group_id int64  // 群号
	User_id  int64  // 管理员 QQ 号
}

type GroupDecrease struct {
	Time        int64  //	事件发生的时间戳
	Self_id     int64  //	收到事件的机器人 QQ 号
	Sub_type    string //	leave|kick|kick_me	事件子类型，分别表示主动退群、成员被踢、登录号被踢
	Group_id    int64  //	群号
	Operator_id int64  //	操作者 QQ 号（如果是主动退群，则和 user_id 相同）
	User_id     int64  //	离开者 QQ 号
}

type GroupIncrease struct {
	Time        int64  //	事件发生的时间戳
	Self_id     int64  //	收到事件的机器人 QQ 号
	Sub_type    string // approve|invite	事件子类型，分别表示管理员已同意入群、管理员邀请入群
	Group_id    int64  //	群号
	Operator_id int64  //	操作者 QQ 号
	User_id     int64  //	加入者 QQ 号
}

type GroupBan struct {
	Time        int64  //	事件发生的时间戳
	Self_id     int64  //	收到事件的机器人 QQ 号
	Sub_type    string //	ban、lift_ban	事件子类型，分别表示禁言、解除禁言
	Group_id    int64  //	群号
	Operator_id int64  //	操作者 QQ 号
	User_id     int64  //	被禁言 QQ 号
	Duration    int64  //	禁言时长，单位秒
}

type FriendAdd struct {
	Time    int64 // 事件发生的时间戳
	Self_id int64 // 收到事件的机器人 QQ 号
	User_id int64 // 新添加好友 QQ 号
}

type GroupRecall struct {
	Time        int64 // 事件发生的时间戳
	Self_id     int64 // 收到事件的机器人 QQ 号
	Group_id    int64 // 群号
	User_id     int64 // 消息发送者 QQ 号
	Operator_id int64 // 操作者 QQ 号
	Message_id  int64 // 被撤回的消息 ID
}
type FriendRecall struct {
	Time       int64 //	事件发生的时间戳
	Self_id    int64 //	收到事件的机器人 QQ 号
	User_id    int64 //	消息发送者 QQ 号
	Message_id int64 //	被撤回的消息 ID
}

type Notify struct {
	Time       int64  // 事件发生的时间戳
	Self_id    int64  // 收到事件的机器人 QQ 号
	Sub_type   string // poke|lucky_king|honor	提示类型:戳一戳|红包运气王|荣誉
	Group_id   int64  // 群号
	User_id    int64  // 发送者 QQ 号|红包发送者 QQ 号|成员 QQ 号
	Target_id  int64  // 被戳者 QQ 号|运气王 QQ 号 //!仅在戳一戳|运气王事件中有值
	Honor_type string // talkative|performer|emotion	荣誉类型:龙王|群聊之火|快乐源泉 //!仅在荣誉事件中有值
}

type FriendRequest struct {
	Time    int64  //	事件发生的时间戳
	Self_id int64  //	收到事件的机器人 QQ 号
	User_id int64  //	发送请求的 QQ 号
	Comment string //	验证信息
	Flag    string //	请求 flag，在调用处理请求的 API 时需要传入
}

type GroupRequest struct {
	Time     int64  //	事件发生的时间戳
	Self_id  int64  //	收到事件的机器人 QQ 号
	Sub_type string //	add、invite	请求子类型，分别表示加群请求、邀请登录号入群
	Group_id int64  //	群号
	User_id  int64  //	发送请求的 QQ 号
	Comment  string //	验证信息
	Flag     string //	请求 flag，在调用处理请求的 API 时需要传入
}

//* Meta Event Part
type MetaLifecycle struct {
	Time     int64  //	事件发生的时间戳
	Self_id  int64  //	收到事件的机器人 QQ 号
	Sub_type string // enable、disable、connect	事件子类型，分别表示 OneBot 启用、停用、WebSocket 连接成功
}

//! ONEBOT: 注意，目前生命周期元事件中，只有 HTTP POST 的情况下可以收到 enable 和 disable，只有正向 WebSocket 和反向 WebSocket 可以收到 connect。

type MetaHeartbeat struct {
	Time     int64       //	事件发生的时间戳
	Self_id  int64       //	收到事件的机器人 QQ 号
	Status   interface{} //	状态信息
	Interval int64       //	到下次心跳的间隔，单位毫秒
}
