package controllers

import (
	"github.com/astaxie/beego"
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"log"
	"fmt"
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
	result,count := selectSqlData(db,sql)
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.Data["data"]=result
	c.Data["count"]=count
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
	//fmt.Fprintln(result)
	data["code"]=0
	data["msg"]="true"
	data["data"]=result
	data["count"]=count
	self.Data["json"]=data
	self.ServeJSON()

}
func selectSqlData(db *sql.DB,sql string) ([]map[string]string,int){
	/*sql查询返回数据*/
	result :=make([]map[string]string,0)
	defer db.Close()
	fmt.Println(sql)
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
	//result := make(map[int]map[string]string)

	//is := make([]interface{}, 0)
	for rows2.Next() {
		//填充数据
		rows2.Scan(scans...)
		//每行数据
		row := make(map[string]string)
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := cols[k]
			//fmt.Printf(string(v))
			//这里把[]byte数据转成string
			row[key] = string(v)
		}
		//放入结果集
		result = append(result,row)
		//result[i] = row
		i++
	}
	fmt.Println(result)
	return result,i
}