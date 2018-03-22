
package main

import (
	"fmt"
	"time"
	"os"
	"os/exec"
	_"strings"
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
var cmdwithpath string
var err error
func init() {
	cmdwithpath, err = exec.LookPath("bash")
	if err != nil {
		fmt.Println("not find bash.")
		os.Exit(5)
	}
}
func test(rpc Rpcparams){
	fmt.Println(rpc)
}
var cmdActionpath string
func init(){
	cmdActionpath,err := exec.LookPath("bash")
	if err != nil{
		fmt.Println("not find bash.")
		os.Exit(5)
	}
	fmt.Println(cmdActionpath)
}
func  Run (rpcparams Rpcparams) error {
	if rpcparams.Commandparams==""{
		fmt.Println("xxxxxx")
		return nil
	}
	//mkdir
	taskpath :="/data/"+rpcparams.Path
	fmt.Println(taskpath)
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
	}
	fmt.Println("11112")
	//svn代码拉取Comm
	err =os.Chdir(taskpath)
	if err != nil {
		return err
	}
	fmt.Println("cmdActionpath=",cmdActionpath)
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
	x := 1
	y := 1
	change(&x)
	changeWithoutPointer(y)
	fmt.Printf("x is %d, y is %d\n", x, y)
	addtimes :=time.Now().Unix()
	fmt.Println(addtimes)
	t := time.Now()
	tm1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	tm2 := tm1.AddDate(0, 0, 1)
	fmt.Println(tm2)
	c:=32
	d:=c >> 4
	fmt.Println(d)
	dir,_:=os.Getwd()
	fmt.Println(dir)
	err:=os.Chdir("/home")
	if err !=nil{
		fmt.Println("xxx")
		fmt.Println(err.Error())
	}
	//a:=[]string{"svn","checkout","http://172.16.56.11/svn/archnews/Mainline/NewsPlatform","--username","chenhuachao","--password","6694d7602e48!@$%","--no-auth-cache","--non-interactive"}
	//svncommand:=fmt.Sprintf("svn checkout  %s  --username %s --password '%s' --no-auth-cache --non-interactive","http://172.16.56.11/svn/archnews/Mainline/NewsPlatform","chenhuachao","6694d7602e48!@$%")
	//cmd:=exec.Command(cmdwithpath,"-c",svncommand)
	//out, err := cmd.CombinedOutput()
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//s:=string(out)
	//fmt.Println(s)
	fmt.Println(os.Getwd())
	var aaa Rpcparams
	aaa =Rpcparams{"root","python test002.py","111.txt","error.log","chenhuachao","svn checkout  -r *  http://172.16.56.11/svn/archnews/Mainline/NewsPlatform/test  --username chenhuachao --password 6694d7602e48!@$% --no-auth-cache --non-interactive","test"}

	err=Run(aaa)
	fmt.Println(err)
}