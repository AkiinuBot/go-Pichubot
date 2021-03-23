// 存放所有结构体
package pichumod

type ConfigStu struct {
	Logmode  string
	Loglvl   int
	WSIP     string
	MasterQQ int
}

type SenderPrivate struct {
	UserID   int
	Nickname string
	Sex      string
	Age      string
}

type SenderGroup struct {
	UserID   int
	Nickname string
	Card     string
	Sex      string
	Age      int
	Area     string
	Level    string
	Role     string
	Title    string
}
