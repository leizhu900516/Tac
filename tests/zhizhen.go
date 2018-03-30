
package main

import (
	"fmt"
	"time"
	"os"
	"os/exec"
	_"strings"
	"math/rand"
	"strconv"
	"strings"
	"net"
)

func change(x *int){
	*x += 1
}

func changeWithoutPointer(y int){
	y += 1
	fmt.Printf("y without pointer is %d\n", y)
}
type CommandParam struct {
	Commandname string
	Commandargs []string
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


func test(rpc Rpcparams){
	fmt.Println(rpc)
}
var err error
var cmdActionpath string
var ipaddress string
func init(){
	cmdActionpath,err = exec.LookPath("bash")
	if err != nil{
		fmt.Println("not find bash.")
		os.Exit(5)
	}
	fmt.Println("cmdActionpath=xxx",cmdActionpath)
}
func gethostIp() string{
	//get host ip address
	addr,err:=net.InterfaceAddrs()
	if err!=nil{
		fmt.Println(err)
	}
	for _,address :=range addr{
		if ipnet,ok :=address.(*net.IPNet);ok && !ipnet.IP.IsLoopback(){
			if ipnet.IP.To4() != nil{
				return  ipnet.IP.String()
				break
			}

		}
	}
	return ""
}
func getProgramId(programname string) (int,error){
	//获取程序pid
	cmdstring :=fmt.Sprintf("ps -ef|grep %s|grep -v grep|awk '{print $2}'",programname)
	cmd := exec.Command(cmdActionpath,"-c",cmdstring)
	out,err :=cmd.CombinedOutput()
	if err!=nil{
		return 0,err
	}
	s :=strings.Replace(string(out),"\n","",-1)
	pid,err :=strconv.Atoi(s)
	if err!=nil {
		return 0,err
	}
	return pid,nil
}
func  Run (rpcparams Rpcparams) error {
	if rpcparams.Commandparams==""{
		fmt.Println("xxxxxx")
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
		return err
	}
	fmt.Println("Svncheckcommand=",rpcparams.Svncheckcommand)
	cmd:=exec.Command(cmdActionpath,"-c",rpcparams.Svncheckcommand)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("eee",err)
		return err
	}
	s:=string(out)
	fmt.Println("222",s)

	//执行任务
	err =os.Chdir(taskpath+"/"+rpcparams.Svnpath)
	if err != nil {
		local_dir,_:=os.Getwd()
		fmt.Println("当前文件目录=",local_dir)
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
func main() {
	rnd :=rand.New(rand.NewSource(time.Now().UnixNano()))
	vcode := fmt.Sprintf("%06v",rnd.Int31n(1000000))
	fmt.Println(vcode)
	addtimes :=time.Now().Unix()
	fmt.Println(addtimes)
	t := time.Now()
	tm1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	tm2 := tm1.AddDate(0, 0, 1)
	fmt.Println(tm2)
	//dir,_:=os.Getwd()
	//fmt.Println(dir)
	//err:=os.Chdir("/home")
	//if err !=nil{
	//	fmt.Println("xxx")
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(cmdActionpath)
	//svncommand:=fmt.Sprintf("svn checkout  %s  --username %s --password '%s' --no-auth-cache --non-interactive","http://172.16.56.11/svn/archnews/Mainline/NewsPlatform","chenhuachao","6694d7602e48!@$%")
	//cmd:=exec.Command(cmdActionpath,"-c",svncommand)
	//out, err := cmd.CombinedOutput()
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//s:=string(out)
	//fmt.Println(s)
	//fmt.Println(os.Getwd())
	//var aaa Rpcparams
	//aaa =Rpcparams{"root","python test002.py","111.txt","error.log","pythontest","svn checkout  http://172.16.56.11/svn/archnews/Mainline/NewsPlatform/test  --username chenhuachao --password 6694d7602e48!@$% --no-auth-cache --non-interactive","test"}
	//
	//err=Run(aaa)
	//fmt.Println(err)
	//pid,err :=getProgramId("test002.py")
	//if err != nil{
	//	fmt.Println(err)
	//}
	//fmt.Println("pid=",pid)
	//var programname []string
	//programname = strings.Split("python test002.py"," ")
	//fmt.Println(programname[len(programname)-1])
	var v interface{}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10; i++{
		v = i
		if (r.Intn(100) % 2) == 0 {
			v = "hello"
		}

		if _, ok := v.(int); ok {
			fmt.Printf("%d\n", v)
		}
	}

}