package controllers

import (
	"github.com/astaxie/beego"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"log"
	"fmt"
	"strconv"
	"encoding/json"
	"time"
	"strings"
)
/*
主页
*/
type MainController struct {
	beego.Controller
}
type LoginController struct {
	beego.Controller
}
type QuitController struct {
	beego.Controller
}
type AuthController struct {
	beego.Controller
}
type BackgroundtaskManageDeleteController struct{
	beego.Controller
}
type BackgroundtaskManageGetController struct{
	beego.Controller
}
type BackgroundtaskManagePostController struct{
	beego.Controller
}
type DelgroundtaskManageGetController struct{
	beego.Controller
}
type BackgroundtaskController struct {
	beego.Controller
}
func (self *QuitController) Get(){
	self.Ctx.SetCookie("u-tac","")
	self.TplName="login.html"
}
func (self *AuthController) Post(){
	var data =make(map[string]interface{})
	var params map[string]string
	json.Unmarshal(self.Ctx.Input.RequestBody,&params)
	username:=params["username"]
	password:=params["password"]
	db,err:=sql.Open("mysql",beego.AppConfig.String("mysqlurl"))
	if err!=nil{
		fmt.Println(err)
	}
	var id int
	err=db.QueryRow("select id from userinfo where username=? and password=?",username,password).Scan(&id)
	if err!=nil{
		data["code"]=1
		data["msg"]=err.Error()
		data["data"]=""
	}
	if id >=0{
		data["code"]=0
		data["msg"]="登陆成功"
		data["data"]=""
		self.Ctx.SetCookie("u-tac",username,3600*24,"/")
	}
	self.Data["json"]=data
	fmt.Println(data)
	self.ServeJSON()

}
func ( self *LoginController) Get() {
	self.TplName = "login.html"
}
func (self *MainController) Get() {
	name :=self.Ctx.GetCookie("u-tac")
	if name== ""{
		self.Ctx.Redirect(302,"/login",)
	}
	db,err := sql.Open("mysql",beego.AppConfig.String("mysqlurl"))
	if err !=nil{
		log.Fatal(err)
	}
	sql := fmt.Sprintf("select * from backgroundtask")
	ipsql :=fmt.Sprintf("select * from taskip where status = 1")
	result,count := selectSqlData(db,sql)
	iplist,ipcount :=selectSqlData(db,ipsql)
	defer db.Close()
	self.Data["Website"] = "beego.me"
	self.Data["Email"] = "astaxie@gmail.com"
	self.Data["data"]=result
	self.Data["count"]=count
	self.Data["name"]=name
	self.Data["iplist"]=iplist
	self.Data["ipcount"]=ipcount
	self.TplName = "index.html"
}
/*
后台任务获取
*/

func (self *BackgroundtaskController) Get(){
	db,err := sql.Open("mysql",beego.AppConfig.String("mysqlurl"))
	if err !=nil{
		log.Fatal(err)
	}
	data :=make(map[string]interface{})
	sql := fmt.Sprintf("select * from backgroundtask")
	result,count := selectSqlData(db,sql)
	data["code"]=0
	data["msg"]="true"
	data["data"]=result
	data["count"]=count
	self.Data["json"]=data
	self.ServeJSON()

}

func (self *BackgroundtaskManageDeleteController) Delete(){
	taskid :=self.Ctx.Input.Param("taskid")
	fmt.Println(taskid)
}
func (self *BackgroundtaskManageGetController) Get(){
	var (
		taskname string
		author string
		hostip string
	)

	data:=make(map[string]interface{})
	taskid := self.Input().Get("taskid")
	pid,_:=strconv.Atoi(taskid)
	db, err := sql.Open("mysql", beego.AppConfig.String("mysqluser")+":"+beego.AppConfig.String("mysqlpass")+"@tcp("+beego.AppConfig.String("mysqlurls")+":"+beego.AppConfig.String("mysqlport")+")/"+beego.AppConfig.String("mysqldb")+"?charset=utf8")
	defer db.Close()
	if err!=nil{
		fmt.Println(err)
	}

	sql :=fmt.Sprintf("select taskname,author,hostip from backgroundtask where pid=%d",pid)
	err1 :=db.QueryRow(sql).Scan(&taskname,&author,&hostip)
	if err1!=nil{
		fmt.Println(err1)
	}
	detail :=make([]string,0)
	detail=append(detail,taskname)
	detail=append(detail,author)
	detail=append(detail,hostip)
	data["code"]=0
	data["data"]=detail
	data["msg"]="success"
	self.Data["json"]=data
	self.ServeJSON()
}
/*添加后台任务*/
func (self *BackgroundtaskManagePostController) Post(){
	data := make(map[string]interface{})
	var params map[string]string
	json.Unmarshal(self.Ctx.Input.RequestBody,&params)
	fmt.Println(params)
	taskname :=params["taskname"]
	ipaddress :=params["ipaddress"]
	url :=params["url"]
	svnuser :=params["svnuser"]
	svnpasswd :=params["svnpasswd"]
	svn_number :=params["svn_number"]
	action_cmd :=params["action_cmd"]
	fmt.Println("action_cmd=",action_cmd)
	addtimes := time.Now().Unix()
	db, err := sql.Open("mysql", beego.AppConfig.String("mysqlurl"))
	if err != nil {
		fmt.Println(err)
	}
	svnurl_split_slice :=strings.Split(url,"/")
	svnpath:=svnurl_split_slice[len(svnurl_split_slice)-1]
	var svncommand string
	if svn_number == "*"{
		svncommand = fmt.Sprintf("svn checkout  %s  --username %s --password %s --no-auth-cache --non-interactive",url,svnuser,svnpasswd)
	}else{
		svncommand = fmt.Sprintf("svn checkout  -r %s  %s  --username %s --password %s --no-auth-cache --non-interactive",svn_number,url,svnuser,svnpasswd)
	}
	rpcparams :=Rpcparams{"root",action_cmd,"111.txt","error.log",taskname,svncommand,svnpath}
	pid :=Rpcclient(ipaddress,rpcparams)
	fmt.Println(">>>>>pid=",pid)
	mysqlparam := MysqlParams{"insert into backgroundtask(taskname,ipaddress,url,svnuser,svnpasswd,svn_number,action_cmd,addtimes) values(?,?,?,?,?,?,?,?)", []string{taskname,ipaddress,url,svnuser,svnpasswd,svn_number,action_cmd},db}
	insertid := mysqlparam.Insert(taskname,ipaddress,url,svnuser,svnpasswd,svn_number,action_cmd,addtimes)
	err = mysqlparam.Update(fmt.Sprintf("update backgroundtask set pid=%s,status=1 where id=%d",pid,insertid))
	checkErr(err)
	fmt.Println(insertid)
	defer db.Close()
	data["code"]=0
	data["msg"]=""
	data["data"]=""
	self.Data["json"]=data
	self.ServeJSON()
}

func (self *DelgroundtaskManageGetController) Get(){
	//删除任务
	data := make(map[string]interface{})
	taskid:=self.Input().Get("taskid")
	ip:=self.Input().Get("ip")
	pid,_:=strconv.Atoi(taskid)
	result:=RpcDelTaskClient(ip,pid)
	if result=="0"{
		data["code"]=0
		data["msg"]=""
		data["data"]=""

	}else{
		data["code"]=1
		data["msg"]=result
		data["data"]=""
	}
	self.Data["json"]=data
	self.ServeJSON()
}
func selectSqlData(db *sql.DB,sql string) ([]map[string]string,int){
	/*sql查询返回数据*/
	result :=make([]map[string]string,0)
	rows2, err := db.Query(sql)
	if err!=nil{
		fmt.Println("xxxxx",err)
		return result,0
	}
	//返回所有列
	cols, _ := rows2.Columns()
	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(cols))
	//这里表示一行填充数据
	scans := make([]interface{}, len(cols))
	//这里scans引用vals，把数据填充到[]byte里
	for k, _ := range vals {
		scans[k] = &vals[k]
	}
	i := 0
	for rows2.Next() {
		//填充数据
		rows2.Scan(scans...)
		//每行数据
		row := make(map[string]string)
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := cols[k]
			//这里把[]byte数据转成string
			row[key] = string(v)
		}
		//放入结果集
		result = append(result,row)
		i++
	}
	fmt.Println(result)
	return result,i
}