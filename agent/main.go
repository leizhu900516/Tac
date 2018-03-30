package main

import (
	"fmt"
	conf "./parserconfig"
)

func main() {
	myConfig := new(conf.Config)
	myConfig.InitConfig("./parserconfig/config")
	fmt.Println(myConfig.Read( "mysqluser"))
	//fmt.Printf("%v", myConfig.Mymap)
}