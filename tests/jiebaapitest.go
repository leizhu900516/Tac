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
	title:=`肖捷在回答记者有关支持中小微企业发展问题时表示，中小微企业是吸纳就业的主力军，是创新、带动投资、促进消费的生力军。今年财政部将继续营造有利于中小微企业发展和创业创新的良好环境。

	　　首先，继续加大财税政策支持力度，特别是继续对中小微企业实施减税措施。其次，推动缓解小微企业融资难融资慢问题，今年要设立国家融资担保基金，完善普惠金融发展的专项资金政策。第三，今年要在政府采购政策中支持中小微企业发展，采取向中小微企业预留政府采购份额、鼓励企业组成联合体投标等方式，进一步优化中小企业参与政府采购活动的市场环境。

	　　另外，针对今年调低赤字率的安排，财政部部长肖捷说，“尽管今年财政预算赤字率比上年有所降低，但是明确告诉大家，积极财政政策的方向没有变。”

	　　2018年政府工作报告提出，今年赤字率拟按2.6%安排，比去年预算低0.4个百分点。对此，肖捷表示，这也是近年来中国财政赤字率首次降低。这与中国经济稳中向好，财政状况不断改善是相吻合的。赤字率的下降，也将为中国经济的长远发展和实施有效的宏观调控留下更大空间。

	　　肖捷说，今年实施积极的财政政策，仍然保持较强力度。首先，今年财政支出的规模，仍然在继续扩大，将达到约21万亿元，比去年增长7.6%。“我们做的是加法，不是减法”。第二，今年的减税降费的力度在继续加大。这也是实施积极财政政策的应有之义。

	　　房地产税立法会以中国国情出发

	　　房地产税的征收一直是百姓关心的与自身利益密切相关的话题，今年的《政府工作报告》中也提出，稳妥推进房地产税立法。

	　　针对这一热点，史耀斌介绍，目前，全国人大常委会预算工作委员会、财政部和其他有关方面正在抓紧起草和完善房地产税法草案。总体思路就是立法先行、充分授权、分步实施。

	　　史耀斌说，房地产税作为世界通行的税种，都有共性制度安排：对所有工商业住房和个人住房都会按照评估值征税；各国都有税收优惠，对困难家庭、低收入家庭给予税收减免；这个税属于地方税，地方政府用这些收入满足教育、治安等公共基础设施的支出；房地产税的税基确定比较复杂，需要建立完备的税收征管模式。

	　　史耀斌介绍，中国现在没有房地产税制度，但是进行过试点，积累了一些经验。他表示，推进这项改革过程中，会注意参考国际上共性制度性安排的一些特点，同时按照中国的国情，从国情出发，合理设计房地产税制度。比如会合并整合税种，会合理降低房地产在建设交易环节的税费负担，以此让房地产税更加公平。`
	urllist=append(urllist,"http://artnlp.eastmoney.com/cutword/keyword")
	urllist=append(urllist,"http://artnlp.eastmoney.com/cutword/cword")
	type Param struct{
		Content string `json:"content"`
		Title string `json:"title"`
	}
	p :=&Param{}
	p.Title=title
	p.Content=text
	var i int
	data,_:=json.Marshal(p)
	result := make(map[string]interface{})
	for i=0;i<10;i++{
		for _,url :=range urllist{
			fmt.Println(url)
			respone ,err :=http.Post(url,body_type,bytes.NewBuffer(data))
			fmt.Println(respone.StatusCode)
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