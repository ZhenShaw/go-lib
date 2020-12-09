package main

import (
	"github.com/zhenshaw/go-lib/logs"
)

func main() {

	z := logs.NewZapLog(logs.DefaultConsole(), logs.DefaultFile())
	L := z.Sugar()

	L.Debug("%s:%s", "aaa", "vvv")

}
