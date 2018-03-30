package main;
/*
@author:chenhuachao
@time: 2018-02-09
@ 计划任务的agent-server,配置文件在/parserconfig/config中
*/
import (
	_"net/rpc"
	_"net/http"
	_"log"
	"os/exec"
	_"io/ioutil"
	"fmt"
	"os"
	"time"
	_"reflect"
	"strconv"
	"strings"
	"net/rpc"
	"log"
	"net/http"
	_"golang.org/x/sys/unix"
	conf "./parserconfig"
	tools "./functools"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
)


var (
	cmdActionpath string
	err error
	AgentConfig conf.Config
	db   *sql.DB
)


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

type Rect struct{}
//进程map缓存：{pid:"pname"}
var processmap =make(map[string]Rpcparams)


func init(){
	cmdActionpath,err = exec.LookPath("bash")
	if err != nil{
		fmt.Println("not find bash.")
		os.Exit(5)
	}
	AgentConfig.InitConfig("./parserconfig/config")
	db,err =sql.Open("mysql",AgentConfig.Read("mysqlurl"))
	if err!=nil{
		fmt.Println(">>>>",err)
	}
	localhostIp:=tools.GetHostIp()
	sql:=fmt.Sprintf("select * from backgroundtask where ipaddress='%s'",localhostIp)
	result :=tools.Select(db,sql)
	fmt.Println(result)
	if len(result)>0{
		for _,values :=range result{
			fmt.Println(">>>>>>>>>>",values)
			pid:=values["pid"]
			taskname:=values["taskname"]
			url:=values["url"]
			svnurl_split_slice :=strings.Split(url,"/")
			svnpath:=svnurl_split_slice[len(svnurl_split_slice)-1]
			action_cmd:=values["action_cmd"]
			rpcparams:=Rpcparams{"root",action_cmd,"","",taskname,"",svnpath}
			fmt.Println(rpcparams)
			processmap[pid]=rpcparams
		}
	}
}

//函数必须是导出的
//必须有两个导出类型参数
//第一个参数是接收参数
//第二个参数是返回给客户端参数，必须是指针类型
//函数还要有一个返回值error

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
		*ret="-1"
		return nil
	}
	//mkdir
	taskpath :="/data/"+rpcparams.Path
	fmt.Println("taskpath=",taskpath)
	_,err := os.Stat(taskpath)
	fmt.Println("......",err)
	if os.IsNotExist(err){
		fmt.Println("vvvv")
		err :=os.MkdirAll(taskpath,0777)
		if err!=nil{
			fmt.Println(err)
			*ret="-1"
			return err
		}else {
			fmt.Println("create dir success")
		}
	}else{
		fmt.Println("dir exits")
	}
	fmt.Println("11112")
	//svn代码拉取Comm
	err =os.Chdir(taskpath)
	if err != nil {
		*ret="-1"
		return err
	}
	cmd:=exec.Command(cmdActionpath,"-c",rpcparams.Svncheckcommand)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("eee",err)
		*ret="-1"
		return err
	}
	s:=string(out)
	fmt.Println(s)
	//执行任务
	err =os.Chdir(taskpath+"/"+rpcparams.Svnpath)
	if err != nil {
		local_dir,_:=os.Getwd()
		fmt.Println("当前文件目录=",local_dir)
		*ret="-1"
		return err
	}
	var cmdString string
	var programlist []string
	programlist = strings.Split(rpcparams.Commandparams," ")
	progranname := programlist[len(programlist)-1]
	if rpcparams.User != "" && rpcparams.User != "root"{
		cmdString = fmt.Sprintf(cmdActionpath, "-c", "su %s -c '%s %s >> %s 2>>%s &'", rpcparams.User, cmdActionpath, rpcparams.Commandparams, rpcparams.Log, rpcparams.Error_log)

	}else{
		cmdString = fmt.Sprintf("setsid  %s >> %s 2>>%s &", rpcparams.Commandparams, rpcparams.Log, rpcparams.Error_log)
	}
	fmt.Println("cmdString=",cmdString)
	cmd = exec.Command(cmdActionpath,"-c",cmdString)
	_,err = cmd.CombinedOutput()
	if err !=nil{
		*ret="-1"
		return err
	}
	pid,err := getProgramId(progranname)
	if err != nil{
		*ret="-1"
		return err
	}
	*ret=pid
	processmap[pid]=rpcparams
	return nil
}
/*
func (r *Rect) RunBack(params CommandParam,ret *string,path string) error {
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
*/
func getProgramId(programname string) (string,error){
	//获取程序pid
	cmdstring :=fmt.Sprintf("ps -ef|grep %s|grep -v grep|awk '{print $2}'",programname)
	cmd := exec.Command(cmdActionpath,"-c",cmdstring)
	out,err :=cmd.CombinedOutput()
	if err!=nil{
		return "-1",err
	}
	pid :=strings.Replace(string(out),"\n","",-1)
	return pid,nil
}
func ProcessIsAlive(pid string,value interface{}) bool{
	/*
	判断进程是否存活
	*/
	procid,err:=strconv.Atoi(pid)
	if err!=nil{
		fmt.Println(err)
	}
	cmd := fmt.Sprintf("ps -ef|grep %d|grep -v grep|wc -l",procid)
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

func start(pid string,rpcparams Rpcparams) bool{
	/*
	启动函数
	*/
	taskpath :="/data/"+rpcparams.Path
	err =os.Chdir(taskpath+"/"+rpcparams.Svnpath)
	if err != nil {
		local_dir,_:=os.Getwd()
		fmt.Println("当前文件目录=",local_dir)
		return false
	}
	programlist := strings.Split(rpcparams.Commandparams," ")
	progranname := programlist[len(programlist)-1]
	var cmdString string
	if rpcparams.User != "" && rpcparams.User != "root"{
		cmdString = fmt.Sprintf(cmdActionpath, "-c", "su %s -c '%s %s >> %s 2>>%s &'", rpcparams.User, cmdActionpath, rpcparams.Commandparams, rpcparams.Log, rpcparams.Error_log)

	}else{
		cmdString = fmt.Sprintf("setsid  %s >> %s 2>>%s &", rpcparams.Commandparams, rpcparams.Log, rpcparams.Error_log)
	}
	fmt.Println("cmdString=",cmdString)
	cmd := exec.Command(cmdActionpath,"-c",cmdString)
	out,err := cmd.CombinedOutput()
	if err !=nil{
		fmt.Println(err)
		return false
	}
	newpid,err := getProgramId(progranname)
	if err != nil{
		fmt.Println(err)
		return false
	}
	delete(processmap,pid)
	processmap[newpid]=rpcparams
	s:=string(out)
	fmt.Println(s)
	//更新数据库的状态和pid
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
	*每十秒检测一次
	*/
	for{
		fmt.Println("healthchecking...")
		for k,v :=range processmap{
			fmt.Println(k,v)
			status :=ProcessIsAlive(k,v)
			fmt.Println(">>>>>",status)
			if status==false{
				//死亡了。。。
				var flag=0
				sql:=fmt.Sprintf("update backgroundtask set status=%d where pid=%s",flag,k)
				tools.Update(db,sql)
				status:=start(k,v)
				if status==true{
					flag =1
					tools.Update(db,sql)
				}
			}else {
				//存活中。。。
				fmt.Println("%d is running",k)
				continue
			}
		}
		time.Sleep(10*time.Second)
	}
}
func main() {
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

}