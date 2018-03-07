package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"bytes"
	_"reflect"
)
var urllist []string

func main() {
	body_type := "application/json;charset=utf-8"
	text :="习近平强调，建设现代化经济体系，事关我们能否引领世界科技革命和产业变革潮流、赢得国际竞争的主动，事关我们能否顺利实现“两个一百年”奋斗目标。要更加重视发展实体经济，把新一代信息技术、高端装备制造、绿色低碳、生物医药、数字经济、新材料、海洋经济等战略性新兴产业发展作为重中之重，构筑产业体系新支柱。要以壮士断腕的勇气，果断淘汰那些高污染、高排放的产业和企业，为新兴产业发展腾出空间。科技创新是建设现代化产业体系的战略支撑。要着眼国家战略需求，主动承接国家重大科技项目，引进国内外顶尖科技人才，加强对中小企业创新支持，培育更多具有自主知识产权和核心竞争力的创新型企业。"
	urllist=append(urllist,"http://artnlp.eastmoney.com/cutword/keyword")
	urllist=append(urllist,"http://artnlp.eastmoney.com/cutword/cword")
	type Param struct{
		Content string `json:"content"`
		Title string `json:"title"`
	}
	p :=&Param{}
	p.Title=text
	p.Content=text
	var i int
	data,_:=json.Marshal(p)
	result := make(map[string]interface{})
	for i=0;i<1000;i++{
		for _,url :=range urllist{
			fmt.Println(url)
			respone ,err :=http.Post(url,body_type,bytes.NewBuffer(data))
			if err!=nil{
				fmt.Println(err.Error())
			}
			defer respone.Body.Close()
			body,_:=ioutil.ReadAll(respone.Body)

			errs :=json.Unmarshal(body,&result)
			if err !=nil{
				fmt.Println(errs)
			}
			fmt.Println(result)
		}
	}

}