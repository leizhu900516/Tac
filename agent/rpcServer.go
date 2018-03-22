package main;

import (
	_"net/rpc"
	_"net/http"
	_"log"
	"os/exec"
	_"io/ioutil"
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
	_"golang.org/x/sys/unix"
)

//go对RPC的支持，支持三个级别：TCP、HTTP、JSONRPC
//go的RPC只支持GO开发的服务器与客户端之间的交互，因为采用了gob编码
var cmdActionpath string
//注意字段必须可导出
type Params struct {
	Width, Height int;
}
type Rpcparams struct {
	User string
	Commandparams string
	Log  string
	Error_log string
	Path string
	Svncheckcommand string
	Svnpath string

}
type CommandParam struct {
	Commandname string
	Commandargs []string
}
type Rect struct{}
//进程map缓存：{pid:"pname"}
var processmap =make(map[int]CommandParam)


func init(){
	cmdActionpath,err := exec.LookPath("bash")
	if err != nil{
		fmt.Println("not find bash.")
		os.Exit(5)
	}
	fmt.Println(cmdActionpath)
}
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
func (r *Rect) Run (rpcparams Rpcparams,ret *string) error {
	if rpcparams.Commandparams==""{
		fmt.Println("xxxxxx")
		return nil
	}
	//mkdir
	taskpath :="/data/"+rpcparams.Path
	fmt.Println(taskpath)
	_,err := os.Stat(taskpath)
	if os.IsNotExist(err){
		err :=Mkdir(taskpath)
		if err!=nil{
			fmt.Println(err)
			return err
		}
	}
	fmt.Println("11112")
	//svn代码拉取Comm
	cmd:=exec.Command(cmdActionpath,"-c",rpcparams.Svncheckcommand)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	s:=string(out)
	fmt.Println("222",s)

	//执行任务
	err =os.Chdir(taskpath+"/"+rpcparams.Svnpath)
	if err != nil {
		return err
	}
	var cmdString string
	if rpcparams.User != "" && rpcparams.User != "root"{
		cmdString = fmt.Sprintf(cmdActionpath, "-c", "su %s -c '%s %s >> %s 2>>%s &'", rpcparams.User, cmdActionpath, rpcparams.Commandparams, rpcparams.Log, rpcparams.Error_log)

	}else{
		cmdString = fmt.Sprintf("setsid  %s >> %s 2>>%s &", rpcparams.Commandparams, rpcparams.Log, rpcparams.Error_log)
	}
	cmd = exec.Command(cmdActionpath,"-c",cmdString)
	_,err = cmd.CombinedOutput()
	if err !=nil{
		return err
	}
	return nil
}

func (r *Rect) RunBack(params CommandParam,ret *string,path string) error {
	//不能放到后台真正执行
		err :=Mkdir(path)
		if err!=nil{
			return err
		}
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

func Mkdir(path string)error{
	cmdString :=fmt.Sprintf(`mkdir -p %s >/dev/null 2>&1`,path)
	cmd := exec.Command(cmdActionpath,"-c",cmdString)
	_,err :=cmd.CombinedOutput()
	if err !=nil{
		return err
	}
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
func stop(pid int) error{
	/*
	停止进程函数
	*/
	cmdString := fmt.Sprintf("kill %s",strconv.Itoa(pid))
	cmd := exec.Command(cmdActionpath,"-c",cmdString)
	_,err := cmd.CombinedOutput()
	if err != nil{
		return err
	}
	return nil
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