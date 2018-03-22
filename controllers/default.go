package controllers

import (
	"github.com/astaxie/beego"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"log"
	"fmt"
	_"reflect"
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

func (c *MainController) Get() {
	db,err := sql.Open("mysql",beego.AppConfig.String("mysqlurl"))
	if err !=nil{
		log.Fatal(err)
	}
	sql := fmt.Sprintf("select * from backgroundtask")
	ipsql :=fmt.Sprintf("select * from taskip where status = 1")
	result,count := selectSqlData(db,sql)
	iplist,ipcount :=selectSqlData(db,ipsql)
	defer db.Close()
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.Data["data"]=result
	c.Data["count"]=count
	c.Data["iplist"]=iplist
	c.Data["ipcount"]=ipcount
	c.TplName = "index.html"
}
/*
后台任务获取
*/
type BackgroundtaskController struct {
	beego.Controller
}
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



type BackgroundtaskManageDeleteController struct{
	beego.Controller
}
type BackgroundtaskManageGetController struct{
	beego.Controller
}
type BackgroundtaskManagePostController struct{
	beego.Controller
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
	action_cmd :=fmt.Sprintf(params["action_cmd"])
	fmt.Println("action_cmd=",action_cmd)
	addtimes := time.Now().Unix()
	db, err := sql.Open("mysql", beego.AppConfig.String("mysqlurl"))
	if err != nil {
		fmt.Println(err)
	}
	svnurl_split_slice :=strings.Split(url,"/")
	svnpath:=svnurl_split_slice[len(svnurl_split_slice)-1]
	svncommand := fmt.Sprintf("svn checkout  -r %s  %s  --username %s --password %s --no-auth-cache --non-interactive",svn_number,url,svnuser,svnpasswd)
	fmt.Println("svncommand=",svncommand)
	rpcparams :=Rpcparams{"root",action_cmd,"111.txt","error.log",taskname,svncommand,svnpath}
	fmt.Println(rpcparams)
	pid :=Rpcclient(ipaddress,rpcparams)
	fmt.Println(pid)
	mysqlparam := MysqlParams{"insert into backgroundtask(taskname,ipaddress,url,svnuser,svnpasswd,svn_number,action_cmd,addtimes) values(?,?,?,?,?,?,?,?)", []string{taskname,ipaddress,url,svnuser,svnpasswd,svn_number,action_cmd},db}
	insertid := mysqlparam.Insert(taskname,ipaddress,url,svnuser,svnpasswd,svn_number,action_cmd,addtimes)
	fmt.Println("insertid=>>>",insertid)
	err = mysqlparam.Update(fmt.Sprintf("update backgroundtask set pid=%d",insertid))
	checkErr(err)
	fmt.Println(insertid)
	defer db.Close()
	data["code"]=0
	data["msg"]=""
	data["data"]=""
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