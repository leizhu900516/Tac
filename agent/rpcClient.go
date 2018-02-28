package main;

import (
	"net/rpc"
	"log"
	"fmt"
)

type Params struct {
	Width, Height int;
}
type Commandparam struct {
	Commandname string
	Commandargs []string
}

func main() {
	//连接远程rpc服务
	rpc, err := rpc.DialHTTP("tcp", "127.0.0.1:8081");
	if err != nil {
		log.Fatal(err);
	}
	//调用远程方法
	//注意第三个参数是指针类型

	//var result string
	//err4 := rpc.Call("Rect.Run", Commandparam{"python",[]string{"/home/GOPATH/src/Tac/tests/sleeptest.py"}}, &result);
	//if err4 != nil {
	//	log.Fatal(err4);
	//}
	//fmt.Println(result)

	var result1 string
	err5 := rpc.Call("Rect.RunBack", Commandparam{"runuser",[]string{"-l","root","-c","python /home/GOPATH/src/Tac/tests/sleeptest.py 2>&1 &"}}, &result1);	if err5 != nil {
		log.Fatal(err5);
	}
	fmt.Println(result1)

	//var result1 string
	//err5 := rpc.Call("Rect.Runcmd", Commandparam{"python",[]string{"/home/GOPATH/src/Tac/tests/sleeptest.py"}}, &result1);
	//if err5 != nil {
	//	log.Fatal(err5);
	//}
	//fmt.Println(result1)
}