package controllers

import (
	"database/sql"
	"fmt"
	"reflect"
	_"strconv"
)

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
	fmt.Println(reflect.TypeOf(param))
	res,err := insert.Exec(param...)
	checkErr(err)
	id,err :=res.LastInsertId()
	checkErr(err)
	fmt.Println(id,reflect.TypeOf(id))
	return id
}





func checkErr(errs error){
	if errs != nil {
		fmt.Println(errs)
	}
	return
}