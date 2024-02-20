package main

import (
	"fmt"
	"sync"
)

var wg3 sync.WaitGroup
var mux1 sync.Mutex
var ch1 = make(chan int, 50)

func main() {
	//VerifyHttp("127.0.0.1:10809")

	fmt.Println("           ___                     ___            _ " +
		"\n ___  ___ | . " + "\\ _ _  ___ __   _ _ | . \\ ___  ___ | |" +
		"\n/ . |/ . \\|  _/| '_>/ . \\\\ \\/| | ||  _// . \\/ . \\| |" +
		"\n\\_. |\\___/|_|  |_|  \\___//\\_\\`_. ||_|  \\___/\\___/|_|" +
		"\n<___'                        <___'                  ")
	Init()
	////开启隧道代理
	//go httpSRunTunnelProxyServer()
	//go socket5RunTunnelProxyServer()
	//启动webAPi
	Run()
}