package doctor

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var tip1= struct {
	Tip string `json:"tip"`
	Type string `json:"type"`
	Id int `json:"id"`
	Source string `json:"source"`
}{
	Tip:"每天快走1小时或慢走1个半小时。",
	Type:"运动",
	Id:1,
	Source:"超重人群每日需额外消耗能量控制体重，消耗约450kcal较为合适，慢走等轻度活动消耗约为255kcal/小时，快走、跳操等中度活动消耗约为500kcal/小时，减轻体重10kg可获得收缩压下降5-20mmHg，具体运动内容可以在随访时根据患者意愿进行修改。",
}
var tip2= struct {
	Tip string `json:"tip"`
	Type string `json:"type"`
	Id int `json:"id"`
	Source string `json:"source"`
}{
	Tip:"三餐30分钟后运动，洗碗、慢走等合计半小时。",
	Type:"运动",
	Id:1,
	Source:"糖尿病人需要在三餐后运动辅助降低血糖，用餐30分钟后进行洗碗、慢走等轻度活动约30分钟对血糖水平下降效果显著。",
}
var tip3= struct {
	Tip string `json:"tip"`
	Type string `json:"type"`
	Id int `json:"id"`
	Source string `json:"source"`
}{
	Tip:"每周额外运动3小时，可分4-6次进行。",
	Type:"运动",
	Id:1,
	Source:"根据国家高血压防治指南，每人每周需要进行中等强度运动每次30分钟，总共约5-7次。根据中国心血管病预防指南，适宜的有氧运动可降低安静时的血压，改善心肺功能，同时调节紧张情绪。",
}
var tip4= struct {
	Tip string `json:"tip"`
	Type string `json:"type"`
	Id int `json:"id"`
	Source string `json:"source"`
}{
	Tip:"戒烟。",
	Type:"吸烟",
	Id:4,
	Source:"大量研究表明，吸烟是各类慢病的独立风险因素之一，存在风险的人群需要科学戒烟，避免被动吸烟。",
}
var tip5= struct {
	Tip string `json:"tip"`
	Type string `json:"type"`
	Id int `json:"id"`
	Source string `json:"source"`
}{
	Tip:"少喝酒，每日白酒少于1两，葡萄酒少于2两，啤酒少于1听。",
	Type:"饮酒",
	Id:2,
	Source:"酒精是各类慢病的独立风险因素之一，上述饮酒量依据来源国家基层高血压管理指南。",
}
var tip6= struct {
	Tip string `json:"tip"`
	Type string `json:"type"`
	Id int `json:"id"`
	Source string `json:"source"`
}{
	Tip:"三餐少油少盐。",
	Type:"饮食",
	Id:3,
	Source:"根据国家基层高血压管理指南，每人每日食盐摄入量不超过6克，可以获得收缩压下降2-8mmHg的效果。",
}
var tip7= struct {
	Tip string `json:"tip"`
	Type string `json:"type"`
	Id int `json:"id"`
	Source string `json:"source"`
}{
	Tip:"建议每天摄入以奶制品、果蔬类为主。",
	Type:"饮食",
	Id:3,
	Source:"根据中国心血管病预防指南，增加膳食中的非精制米面、果蔬摄入量，减少膳食中的总脂肪，可降低血脂和改善心血管健康。",
}
var tip8= struct {
	Tip string `json:"tip"`
	Type string `json:"type"`
	Id int `json:"id"`
	Source string `json:"source"`
}{
	Tip:"减轻精神压力，可以到专业医疗机构就诊。",
	Type:"心理",
	Id:9,
	Source:"精神压力容易引起血压升高，风险患者需要通过一定途径减轻精神压力。",
}
var tip9= struct {
	Tip string `json:"tip"`
	Type string `json:"type"`
	Id int `json:"id"`
	Source string `json:"source"`
}{
	Tip:"药物治疗。是否采取药物治疗需要获得医嘱后确认！",
	Type:"服药",
	Id:5,
	Source:"必要时可以配合药物进行治疗，请修改具体内容后再勾选该项！",
}
var tip10= struct {
	Tip string `json:"tip"`
	Type string `json:"type"`
	Id int `json:"id"`
	Source string `json:"source"`
}{
	Tip:"保证每天睡眠6-8小时。",
	Type:"心理",
	Id:9,
	Source:"保持身心愉悦，保证充足睡眠是减轻精神压力的有效方式之一。",
}
var tip11= struct {
	Tip string `json:"tip"`
	Type string `json:"type"`
	Id int `json:"id"`
	Source string `json:"source"`
}{
	Tip:"每餐主食不超过2两（约1碗米饭）。",
	Type:"饮食",
	Id:3,
	Source:"糖尿病病人需要控制饮食中主食的摄入，米饭中含有大量的糖分，提倡少食多餐。",
}
var tip12= struct {
	Tip string `json:"tip"`
	Type string `json:"type"`
	Id int `json:"id"`
	Source string `json:"source"`
}{
	Tip:"每日测量1-2次血压并上传记录。",
	Type:"自测",
	Id:6,
	Source:"根据中国心血管病预防指南，高血压是心血管疾病独立的、最重要的危险因素，收缩压从115mmHg开始与心血管病风险呈正相关，需要控制血压。高血压患者调整治疗期间每日至少测量2次进行血压，持续观测变化情况，血压平稳后每周监测血压2次。",
}
var tip13= struct {
	Tip string `json:"tip"`
	Type string `json:"type"`
	Id int `json:"id"`
	Source string `json:"source"`
}{
	Tip:"每周测量一次体重。",
	Type:"自测",
	Id:6,
	Source:"观测体重变化。",
}
var tip14= struct {
	Tip string `json:"tip"`
	Type string `json:"type"`
	Id int `json:"id"`
	Source string `json:"source"`
}{
	Tip:"每天阅读血压助手中推送内容10分钟。",
	Type:"自测",
	Id:6,
	Source:"患者对高血压风险认识不足。考虑到患者文化程度可能偏大，处于此阶段的患者需要在随访时对其进行反复的认知教育，告知其高血压、高血糖等疾病的危害。",
}


func Tips(w http.ResponseWriter,r *http.Request){
	var tip []struct {
		Tip  string `json:"tip"`
		Type string `json:"type"`
		Id int `json:"id"`
		Source string `json:"source"`
	}
	db, err := sql.Open("mysql", "huang:zjubio2019@tcp(mysql.zjubiomedit.com:3306)/cdms_biomedit?charset=utf8")
	CheckErr(err)
	r.ParseForm()
	var sex string
	if r.Method == "GET" {
		patientId= strings.Join(r.Form["patientId"], "")
		sex=strings.Join(r.Form["sex"], "")
	}
	if sex=="女"{
		tip5.Tip="少喝酒，每日白酒少于半两，葡萄酒少于100ml，啤酒少于200ml。"
	}
	var tbc,Hdlc,smoke,bp,dia,bmi,drink,attitude string
	rows1, _ := db.Query("select tc from eval_riskfactors where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows1.Next(){
		err = rows1.Scan(&tbc)
		CheckErr(err)
	}
	rows2, _ := db.Query("select smoke from eval_riskfactors where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows2.Next(){
		err = rows2.Scan(&smoke)
		CheckErr(err)
	}
	rows3, _ := db.Query("select Hdlc from eval_riskfactors where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows3.Next(){
		err = rows3.Scan(&Hdlc)
		CheckErr(err)
	}
	rows4, _ := db.Query("select sbp from eval_riskfactors where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows4.Next(){
		err = rows4.Scan(&bp)
		CheckErr(err)
	}
	rows5, _ := db.Query("select dia from eval_riskfactors where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows5.Next(){
		err = rows5.Scan(&dia)
		CheckErr(err)
	}
	rows6, _ := db.Query("select bmi from eval_riskfactors where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows6.Next(){
		err = rows6.Scan(&bmi)
		CheckErr(err)
	}
	rows7, _ := db.Query("select drink from eval_riskfactors where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows7.Next() {
		err = rows7.Scan(&drink)
		CheckErr(err)
	}
	rows8, _ := db.Query("select attitude from eval_userdata where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows8.Next() {
		err = rows8.Scan(&attitude)
		CheckErr(err)
	}
	if bmi=="肥胖"{
		if attitude=="无"||attitude=="有意图"{
			tip1.Tip = "每天合计快走30分钟或慢走1小时，并在一个月内逐渐提高到快走75分钟或慢走2小时。"
			tip1.Source = "先设定简单目标，达成后向下一阶段逐步过度。具体运动内容可以在随访时根据患者意愿进行修改。"
		}else {
			tip1.Tip = "每天快走75分钟或慢走2小时。"
			tip1.Source = "肥胖人群每日需额外消耗能量控制体重，消耗约600kcal较为合适，慢走等轻度活动消耗约为255kcal/小时，快走、跳操等中度活动消耗约为500kcal/小时，减轻体重10kg可获得收缩压下降5-20mmHg，具体运动内容可以在随访时根据患者意愿进行修改。"
		}
	}
	if smoke=="吸烟"{
		if attitude=="无"{
			tip4.Tip="每天阅读抽烟对健康影响的内容5分钟。"
			tip4.Source="用于提高患者对于吸烟危害的认识，促使他们进入下一心理阶段。"
		}else if attitude=="有意图"{
			tip4.Tip="今日起控制每日抽烟量在半包以下，在想抽烟时进行1-3分钟的随地运动。"
			tip4.Source="逐步降低抽烟量，运动改变注意力缓解抽烟欲望。"
		}else if attitude=="有准备"{
			tip4.Tip="今日起控制吸烟量，2个月内逐步达成戒烟的目标。"
			tip4.Source="逐步控制吸烟量以达成戒烟目标。"
		}else if attitude=="已在改变"{
			tip4.Tip="从接收到此计划开始停止抽烟。"
			tip4.Source="大量研究表明，吸烟是各类慢病的独立风险因素之一，存在风险的人群需要科学戒烟，避免被动吸烟。"
		}
		tip=append(tip,tip4)
	}
	if drink == "饮酒" {
		tip = append(tip, tip5)
	}
	if bmi=="超重"||bmi=="肥胖"{
		tip=append(tip,tip1)
		tip=append(tip,tip13)
	}
	if bmi=="超重"||bmi=="肥胖"||bp=="血压过高"||tbc=="血总胆固醇异常"||Hdlc=="高密度脂蛋白胆固醇异常" {
		tip = append(tip, tip3)
		tip = append(tip, tip6)
		tip = append(tip, tip7)
	}
	if bp=="血压过高"{
		tip=append(tip,tip8)
		tip=append(tip,tip12)
		if attitude=="无"{
			tip=append(tip,tip14)
		}
	}
	if bp=="血压过高"||tbc=="血总胆固醇异常"||Hdlc=="高密度脂蛋白胆固醇异常"{
		tip=append(tip,tip9)
	}
	if dia=="血糖过高"{
		if attitude=="无"||attitude=="有意图"{
			tip2.Tip = "三餐30分钟后进行运动，洗碗、慢走等合计5分钟，并在1个月内逐渐提升到半小时。"
			tip2.Source = "先设定简单的目标，在患者达成后逐渐向下一阶段过度。"
		}
		tip=append(tip,tip2)
		tip=append(tip,tip11)
	}
	tip=append(tip,tip10)
	b, _ := json.Marshal(tip)
	fmt.Fprintf(w,string(b))
	db.Close()
}
