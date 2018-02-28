package main;

import (
	_"net/rpc"
	_"net/http"
	_"log"
	"os/exec"
	"io/ioutil"
	"fmt"
	"context"
	"os"
	"time"
	_"reflect"
	"strconv"
	"strings"
	"net/rpc"
	"log"
	"net/http"
	"syscall"
	_"golang.org/x/sys/unix"
)

//go对RPC的支持，支持三个级别：TCP、HTTP、JSONRPC
//go的RPC只支持GO开发的服务器与客户端之间的交互，因为采用了gob编码

//注意字段必须可导出
type Params struct {
	Width, Height int;
}

type CommandParam struct {
	Commandname string
	Commandargs []string
}
type Rect struct{}
//进程map缓存：{pid:"pname"}
var processmap =make(map[int]CommandParam)



//函数必须是导出的
//必须有两个导出类型参数
//第一个参数是接收参数
//第二个参数是返回给客户端参数，必须是指针类型
//函数还要有一个返回值error
func (r *Rect) Area(p Params, ret *int) error {
	*ret = p.Width * p.Height;
	return nil;
}

func (r *Rect) Perimeter(p Params, ret *int) error {
	*ret = (p.Width + p.Height) * 2;
	return nil;
}
/*
*执行cmd命令的函数Run
*/
func (r *Rect) Run(params CommandParam,ret *string) error{
	fmt.Println(params)
	go func(){
		cmd :=exec.Command(params.Commandname,params.Commandargs...)
		//fmt.Println(cmd.Process.Pid)
		stdout, err := cmd.StdoutPipe()
		if err != nil {

			*ret =err.Error()
		}else {
			// 运行命令
			if err := cmd.Start(); err != nil {
				*ret =err.Error()
			}else {
				// 读取输出结果
				opBytes, err := ioutil.ReadAll(stdout)
				if err != nil {
					*ret =err.Error()
				}else {
					*ret =string(opBytes)
				}
			}
		}
		// 保证关闭输出流
		defer stdout.Close()
	}()
	return nil
}

func (r *Rect) RunBack(params CommandParam,ret *string) error {
	//不能放到后台真正执行
		ctx,_ := context.WithCancel(context.Background())
		cmd :=exec.CommandContext(ctx,params.Commandname,params.Commandargs...)
		cmd.Stdout = os.Stdout
		cmd.Start()
		processid :=cmd.Process.Pid+3
		processmap[processid]=params
		fmt.Println("程序id是", processid)
		//cancel()   //是否杀死进程
		//cmd.Wait() //是否等待进程结束
	return nil
}

func (r *Rect) Runcmd(params CommandParam,ret *string) error {
	cmd :=exec.Command(params.Commandname,params.Commandargs...)
	//将其他命令传入生成出的进程
	//给新进程设置文件描述符，可以重定向到文件中
	cmd.Stdin=os.Stdin
	cmd.Stdout=os.Stdout
	cmd.Stderr=os.Stderr
	//开始执行新进程，不等待新进程退出
	cmd.Start()
	return nil
}
func ProcessIsAlive(pid int,value interface{}) bool{
	/*
	判断进程是否存活
	*/
	cmd := fmt.Sprintf("ps -ef|grep %d|grep -v grep|wc -l",pid)
	fmt.Println(cmd)
	out ,err :=exec.Command("bash","-c",cmd).Output()
	if err !=nil {
		return false
	}
	nums,err:=strconv.Atoi(strings.Replace(string(out),"\n","",-1))
	if err!=nil{
		return false
	}
	if nums==0{
		return false
	}else if nums>0{
		return true
	}
	fmt.Println(nums)
	return true
}

func start(pid int,params CommandParam) bool{
	/*
	启动函数
	*/
	ctx,_ := context.WithCancel(context.Background())
	cmd :=exec.CommandContext(ctx,params.Commandname,params.Commandargs...)
	//cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}//linux  not compat
	cmd.Stdout = os.Stdout
	cmd.Start()
	processid :=cmd.Process.Pid+3
	_,ok :=processmap[pid]
	if ok{
		delete(processmap,pid)
	}
	processmap[processid]=params
	return true
}
func stop(pid int){
	/*
	停止进程函数
	*/
	syscall.Kill(int(pid), syscall.SIGKILL)
}
func restart(){
	/*
	重启进程函数
	*/
}
func Healthcheck(){
	/*
	*进程状态健康检测函数
	*/
	for{
		for k,v :=range processmap{
			fmt.Println(k,v)
			status :=ProcessIsAlive(k,v)
			fmt.Println(">>>>>",status)
			if status==false{
				//没有存活
				start(k,v)
			}else {
				//存活了
				fmt.Println("%d is running",k)
				continue
			}
		}
		time.Sleep(10*time.Second)
	}
}
func main() {
	//err1:=unix.Chdir("/home")
	//if err1!=nil{
	//	fmt.Println(err1)
	//}else{
	//	fmt.Println("切换目录成功")
	//}
	rect := new(Rect);
	//注册一个rect服务
	rpc.Register(rect);
	//把服务处理绑定到http协议上
	rpc.HandleHTTP();
	log.Println("start rpc server on 8081!")
	go Healthcheck()
	err := http.ListenAndServe(":8081", nil);
	if err != nil {
		log.Fatal(err);
	}
	//status :=ProcessIsAlive(12,"a")
	//fmt.Println(status)


}