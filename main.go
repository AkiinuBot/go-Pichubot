package main

import (
	"github.com/0ojixueseno0/go-Pichubot-base/mainbot"
	"github.com/0ojixueseno0/go-Pichubot-base/pichumod"
)

func init() {
	mainbot.OnPrivateMsg = pichumod.PrivateParse
	mainbot.OnGroupMsg = pichumod.GroupParse
}

func main() {
	mainbot.Run()
}
