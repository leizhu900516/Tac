
package main

import (
	"fmt"
	"time"
	"os"
	"os/exec"
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
var cmdwithpath string
var err error
func init() {
	cmdwithpath, err = exec.LookPath("bash")
	if err != nil {
		fmt.Println("not find bash.")
		os.Exit(5)
	}
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
	svncommand:=fmt.Sprintf("svn checkout  %s  --username %s --password '%s' --no-auth-cache --non-interactive","http://172.16.56.11/svn/archnews/Mainline/NewsPlatform","chenhuachao","6694d7602e48!@$%")
	cmd:=exec.Command(cmdwithpath,"-c",svncommand)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err.Error())
	}
	s:=string(out)
	fmt.Println(s)
	fmt.Println(os.Getwd())
}