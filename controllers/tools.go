package controllers

import (
	"database/sql"
	"fmt"
	"reflect"
	_"strconv"
	rpc2 "net/rpc"
	_"log"
)
/*
*mysql操作接口
*/
type MysqlOperate interface {
	Select()([]map[string]string,int)
	Insert() int64
	Update()
	Delete()
}

type MysqlParams struct {
	sql string
	param []string
	db *sql.DB
}

func (mp MysqlParams) Select()([]map[string]string,int){
	result := make([]map[string]string,0)
	rows2 ,err :=mp.db.Query(mp.sql)
	defer mp.db.Close()
	if err != nil {
		fmt.Println(err)
		return result,0
	}
	cols ,_ := rows2.Columns()
	vals := make([][]byte,len(cols))
	scans := make([]interface{},len(cols))
	for k,_ :=range vals {
		scans[k] = &vals[k]
	}
	i:=0
	for rows2.Next() {
		rows2.Scan(scans...)
		row := make(map[string]string)
		for k,v := range vals{
			key := cols[k]
			row[key] = string(v)
		}
		result = append(result,row)
		i++
	}
	return result,i
}
func (mp MysqlParams) Insert(param ...interface{}) int64{
	insert, err := mp.db.Prepare(mp.sql)
	checkErr(err)
	//defer mp.db.Close()
	fmt.Println(reflect.TypeOf(param))
	res,err := insert.Exec(param...)
	checkErr(err)
	id,err :=res.LastInsertId()
	checkErr(err)
	fmt.Println(id,reflect.TypeOf(id))
	return id
}

func (mp MysqlParams) Update(sql string) error {
	/*mysql update function*/
	stmt,err := mp.db.Prepare(sql)
	checkErr(err)
	//defer mp.db.Close()
	res,err := stmt.Exec()
	lastId,err := res.LastInsertId()
	if err != nil{
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println(lastId,rowCnt)
	return nil

}




func checkErr(errs error){
	if errs != nil {
		fmt.Println(errs)
	}
	return
}
//rpc公共类部分
type Commandparam struct {
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
//func Rpcclient(ip string,cmdline Commandparam) string{
//	/*
//	*rpc客户端连接程序
//	*/
//	rpcClient,err := rpc2.DialHTTP("tcp",ip)
//	if err != nil {
//		log.Fatal(err)
//	}
//	var result1 string
//	err5 := rpcClient.Call("Rect.RunBack", cmdline, &result1);	if err5 != nil {
//		log.Fatal(err5)
//	}
//	fmt.Println(result1)
//	return result1
//}
func Rpcclient(ip string,rpcparams Rpcparams) string{
	/*
	*rpc客户端连接程序
	*/
	fmt.Println(rpcparams)
	ipaddress :=ip+":8081"
	fmt.Println("ipaddress=",ipaddress)
	rpcClient,err := rpc2.DialHTTP("tcp",ipaddress)
	if err != nil {
		fmt.Println(err)
	}
	var result1 string
	err5 := rpcClient.Call("Rect.Run", &rpcparams, &result1);if err5 != nil {
		fmt.Println(err5)
	}
	fmt.Println(result1)
	return result1
}
func RpcDelTaskClient(ip string,pid int) string{
	//删除任务rpc客户端
	ipaddress :=ip+":8081"
	rpcClient,err := rpc2.DialHTTP("tcp",ipaddress)
	if err !=nil{
		fmt.Println(err)
	}
	var result string
	err1 :=rpcClient.Call("Rect.Deltask",pid,&result);if err1!=nil{
		fmt.Println(err)
	}
	return result
}
