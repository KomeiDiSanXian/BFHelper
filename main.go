// Package main 主程序
package main

import (
	_ "github.com/KomeiDiSanXian/BFHelper/bfhelper"
	_ "github.com/KomeiDiSanXian/BFHelper/console"

	"github.com/KomeiDiSanXian/BFHelper/botsetting"

	"github.com/FloatTech/floatbox/process"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
)

var c zero.Config

func init() {
	s := botsetting.Read()
	// 使用正向ws
	c.Driver = []zero.Driver{driver.NewWebSocketClient(s.WSClient, s.AccessToken)}

	c.NickName = s.BotNames
	c.CommandPrefix = s.CommandPrefix
	c.SuperUsers = s.SuperUsers
	c.RingLen = s.RingLen
	c.MarkMessage = s.MarkMessage
}

func main() {
	zero.RunAndBlock(&c, process.GlobalInitMutex.Unlock)
}
