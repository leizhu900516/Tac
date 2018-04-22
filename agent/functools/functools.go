//定义常用的工具函数=>mysql的操作函数+常用工具集


package functools

import (
	"net"
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)


func GetHostIp() string{
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

func  Select(db *sql.DB,sql string) []map[string]string{
	result :=make([]map[string]string,0)
	fmt.Println("///////",sql)
	rows,err:=db.Query(sql)
	if err!=nil{
		fmt.Println("jkgjkagjg",err)
		return result
	}
	defer db.Close()
	cols,_:=rows.Columns()
	vals :=make([][]byte,len(cols))
	scans:=make([]interface{},len(cols))
	for k,_ :=range vals {
		scans[k] = &vals[k]
	}

	for rows.Next(){
		rows.Scan(scans...)
		row :=make(map[string]string)
		for k,v :=range vals{
			key:=cols[k]
			row[key]=string(v)
		}
		result=append(result,row)

	}
	return result
}

func Update(db *sql.DB,sql string){
	//update sql
	stmt,err :=db.Prepare(sql)
	defer db.Close()
	if err!=nil{
		fmt.Println(err)
	}
	ret,err:=stmt.Exec()
	if err!=nil{
		fmt.Println(err)
	}
	affect,err:=ret.RowsAffected()
	if err!=nil{
		fmt.Println(err)
	}
	fmt.Println("affect=",affect)
}