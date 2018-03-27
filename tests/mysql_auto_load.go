package main

import (
	"database/sql"
	_"github.com/go-sql-driver/mysql"
	"fmt"
)

type Newsinfo struct{
	DataSourceName string
	Content_forget string
	Title string
	Publishtime int
	Image string
	Summary string
	Cateid int8
}
func insertformal(newsinfo Newsinfo){
	//插入正式环境
	var title string
	db,err:=sql.Open("mysql","wordpress_hc:gc895316@tcp(113.10.195.169:3306)/wechat?charset=utf8")
	if err!=nil{
		fmt.Println(err)
	}
	defer db.Close()
	err=db.QueryRow("select title from wx_fb_news where title='?'",newsinfo.Title).Scan(&title)
	if err!=nil{
		fmt.Println(err)
	}
	if title!=""{
		fmt.Println("没有数据，可以插入")

	}else {
		stmt,_:=db.Prepare("insert into wx_fb_news (title,content,pubtime,show_url,abstract,cateid) values('?','?',?,'?','?',?)")
		ret,err :=stmt.Exec(newsinfo.Title,newsinfo.Content_forget,newsinfo.Publishtime,newsinfo.Image,newsinfo.Summary,newsinfo.Cateid)
		if err!=nil{
			fmt.Println(err)
		}
		if lastinsertid,err:=ret.LastInsertId();err==nil{
			fmt.Println(lastinsertid)
		}
	}


}
func main(){
	newsinfo :=Newsinfo{}
	localdb,err :=sql.Open("mysql","wordpress_hc:gc895316@tcp(113.10.195.169:3306)/wechat?charset=utf8")
	if err!=nil{
		fmt.Println(err)
	}
	defer localdb.Close()
	//rows,err :=db.Query("select title,content_forged,publishtime,image,summary,type from wx_fb_news where online_id=0")
	rows,err :=localdb.Query("select title,content,pubtime,show_url,abstract,cateid from wx_fb_news where online_id=0")
	if err!=nil{
		fmt.Println(err)
	}
	for rows.Next(){
		if err:=rows.Scan(&newsinfo.Title,&newsinfo.Content_forget,&newsinfo.Publishtime,&newsinfo.Image,&newsinfo.Summary,&newsinfo.Cateid);err!=nil{
			fmt.Println(err)
		}
		fmt.Println(newsinfo.Summary)
	}


}