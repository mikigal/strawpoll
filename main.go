package main

import (
	"./mysql"
	"./web"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	mysql.Connect("root@tcp(127.0.0.1:3306)/strawpoll")
	web.Start()
}
