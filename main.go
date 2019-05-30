package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql" //_表示只使用init方法
	"time"

	//"html"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"demo/doctor"
)

var Username,Password string
var score1,score2,score int
var sex,age,drink,smoke,diabetes,tbc,Hdlc,Bheight,Bweight,SBP,patientId,measureTime,phy string
var Fscore,Iscore int
var Fscore1,Fscore2,Fscore3,Fscore4,Fscore5 int
var Iscore1,Iscore2,Iscore3,Iscore4,Iscore5,Iscore6 int
var Frisk,Wrisk int
var Irisk float64
var fsmoke,fdiabetes,ftbc,fHdlc,fbmi,fsbp,fdrink string
var readsex,readage string
var tbclab,Hdlclab,bmilab int
var Nowtime string

type factorlist struct{
	FHdlc string `json:"fHdlc"`
	Fbmi string `json:"fbmi"`
	Fdiabetes string `json:"fdiabetes"`
	Fdrink string `json:"fdrink"`
	Fsbp string `json:"fsbp"`
	Fsmoke string `json:"fsmoke"`
	Ftbc string `json:"ftbc"`
}

type mymux struct {
}

func(p *mymux)ServeHTTP(w http.ResponseWriter, r *http.Request){
    if r.URL.Path == "/Frisk"{
    	evaluateF(w,r)
	}
    if r.URL.Path == "/Wrisk"{
    	evaluateW(w)
	}
    if r.URL.Path == "/Irisk"{
    	evaluateI(w)
	}
    if r.URL.Path =="/factor"{
    	insertfactor()
	}
    if r.URL.Path == "/showfactor" {
		showfactor(w)
	}
    if r.URL.Path == "/sexrecord"{
    	sexrecord(w,r)
	}
    if r.URL.Path =="/agerecord"{
    	agerecord(w,r)
	}
	if r.URL.Path == "/doctor/types"{
		doctor.Types(w,r)
	}
	if r.URL.Path == "/doctor/evaluatelab"{
		doctor.Lab(w,r)
	}
	if r.URL.Path == "/doctor/getrecord" {
		doctor.Getrecord(w,r)
	}
	if r.URL.Path == "/doctor/patient/information" {
		doctor.Information(w,r)
	}
	if r.URL.Path == "/doctor/tags" {
		doctor.Tags(w,r)
	}
	if r.URL.Path == "/doctor/gettags" {
		doctor.Gettags(w,r)
	}
	if r.URL.Path == "/doctor/deletetag" {
		doctor.Deletetag(w,r)
	}
	if r.URL.Path == "/doctor/tips" {
		doctor.Tips(w,r)
	}
	return
}

func main() {
	mux := &mymux{}
	err := http.ListenAndServe(":9010", mux) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}


func insertevluate(db *sql.DB,item1,item2,item3,item4,item5,item6,item7,item8,item9,item10,item11,item12,item13 string){
	db.Exec("INSERT into eval_userdata(sex,age,drink,smoke,dia,tc,Hdlc,height,weight,sbp,patientId,measureTime,attitude) values (?,?,?,?,?,?,?,?,?,?,?,?,?)",item1,item2,item3,item4,item5,item6,item7,item8,item9,item10,item11,item12,item13)
}

func evaluateF(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "huang:zjubio2019@tcp(mysql.zjubiomedit.com:3306)/cdms_biomedit?charset=utf8")
	CheckErr(err)
	r.ParseForm()

	if r.Method == "POST" {
		result, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()
		fmt.Printf("%s\n", result)

		//未知类型的推荐处理方法
		var f interface{}
		json.Unmarshal(result, &f)
		m := f.(map[string]interface{})
		fmt.Printf("%s\n", m)

		for k, v := range m {

			switch k {
			case "sex":
				sex = v.(string)
			case "age":
				age = v.(string)
			case "drink":
				drink = v.(string)
			case "smoke":
				smoke = v.(string)
			case "diabetes":
				diabetes = v.(string)
			case "tbc":
				tbc = v.(string)
			case "Hdlc":
				Hdlc = v.(string)
			case "Bheight":
				Bheight = v.(string)
			case "Bweight":
				Bweight = v.(string)
			case "SBP":
				SBP=v.(string)
			case "patientId":
				patientId=v.(string)
			case "measureTime":
				measureTime=v.(string)
			case "phy":
				phy=v.(string)
			}
		}
		insertevluate(db, sex, age, drink, smoke, diabetes, tbc, Hdlc, Bheight, Bweight,SBP,patientId,measureTime,phy)
		db.Close()
	}
	calculateF()
	sFrisk()
	fmt.Fprintf(w,  strconv.Itoa(Frisk))
}

func calculateF(){
	db, err := sql.Open("mysql", "huang:zjubio2019@tcp(mysql.zjubiomedit.com:3306)/cdms_biomedit?charset=utf8")
	CheckErr(err)
	rows1, _ := db.Query("select sex from eval_userdata ORDER BY id DESC LIMIT 1")
	rows2, _ := db.Query("select age from eval_userdata ORDER BY id DESC LIMIT 1")
	rows3, _ := db.Query("select tc from eval_userdata ORDER BY id DESC LIMIT 1")
	rows4,_:=db.Query("select smoke from eval_userdata ORDER BY id DESC LIMIT 1")
	rows5,_:=db.Query("select Hdlc from eval_userdata ORDER BY id DESC LIMIT 1")
	rows6,_:=db.Query("select sbp from eval_userdata ORDER BY id DESC LIMIT 1")

	for rows1.Next() {
		var sex string

		err = rows1.Scan(&sex)
		if err != nil {
			panic(err)
		}
		switch sex {
		case "男":
			for rows2.Next() {
				var age int

				err = rows2.Scan(&age)
				if err != nil {
					panic(err)
				}
				switch {
				case age >= 20 && age < 35:
					Fscore1=-9
					for rows3.Next(){
						var tc float32
						err=rows3.Scan(&tc)
						if err!=nil{
							panic(err)
						}
						tc=tc*38
						switch {
						case tc<160:
							Fscore2=0
						case tc>=160&&tc<200:
							Fscore2=4
						case tc>=200&&tc<240:
							Fscore2=7
						case tc>=240&&tc<280:
							Fscore2=9
						case tc>=280:
							Fscore2=11
						}
					}
					for rows4.Next(){
						var smoke string
						err=rows4.Scan(&smoke)
						if err!=nil{
							panic(err)
						}
						switch smoke {
						case "吸烟":
							Fscore3=8
						case "不吸烟":
							Fscore3=0
						}
					}
				case age>=35&&age<40:
					Fscore1=-4
					for rows3.Next(){
						var tc float32
						err=rows3.Scan(&tc)
						if err!=nil{
							panic(err)
						}
						tc=tc*38
						switch {
						case tc<160:
							Fscore2=0
						case tc>=160&&tc<200:
							Fscore2=4
						case tc>=200&&tc<240:
							Fscore2=7
						case tc>=240&&tc<280:
							Fscore2=9
						case tc>=280:
							Fscore2=11
						}
					}
					for rows4.Next(){
						var smoke string
						err=rows4.Scan(&smoke)
						if err!=nil{
							panic(err)
						}
						switch smoke {
						case "吸烟":
							Fscore3=8
						case "不吸烟":
							Fscore3=0
						}
					}
				case age>=40&&age<45:
					Fscore1=0
					for rows3.Next(){
						var tc float32
						err=rows3.Scan(&tc)
						if err!=nil{
							panic(err)
						}
						tc=tc*38
						switch  {
						case tc<160:
							Fscore2=0
						case tc>=160&&tc<200:
							Fscore2=3
						case tc>=200&&tc<240:
							Fscore2=5
						case tc>=240&&tc<280:
							Fscore2=6
						case tc>=280:
							Fscore2=8
						}
					}
					for rows4.Next(){
						var smoke string
						err=rows4.Scan(&smoke)
						if err!=nil{
							panic(err)
						}
						switch smoke {
						case "吸烟":
							Fscore3=5
						case "不吸烟":
							Fscore3=0
						}
					}
				case age>=45&&age<50:
					Fscore1=3
					for rows3.Next(){
						var tc float32
						err=rows3.Scan(&tc)
						if err!=nil{
							panic(err)
						}
						tc=tc*38
						switch  {
						case tc<160:
							Fscore2=0
						case tc>=160&&tc<200:
							Fscore2=3
						case tc>=200&&tc<240:
							Fscore2=5
						case tc>=240&&tc<280:
							Fscore2=6
						case tc>=280:
							Fscore2=8
						}
					}
					for rows4.Next(){
						var smoke string
						err=rows4.Scan(&smoke)
						if err!=nil{
							panic(err)
						}
						switch smoke {
						case "吸烟":
							Fscore3=5
						case "不吸烟":
							Fscore3=0
						}
					}
				case age>=50&&age<55:
					Fscore1=6
					for rows3.Next(){
						var tc float32
						err=rows3.Scan(&tc)
						if err!=nil{
							panic(err)
						}
						tc=tc*38
						switch  {
						case tc<160:
							Fscore2=0
						case tc>=160&&tc<200:
							Fscore2=2
						case tc>=200&&tc<240:
							Fscore2=3
						case tc>=240&&tc<280:
							Fscore2=4
						case tc>=280:
							Fscore2=5
						}
					}
					for rows4.Next(){
						var smoke string
						err=rows4.Scan(&smoke)
						if err!=nil{
							panic(err)
						}
						switch smoke {
						case "吸烟":
							Fscore3=3
						case "不吸烟":
							Fscore3=0
						}
					}
				case age>=55&&age<60:
					Fscore1=8
					for rows3.Next(){
						var tc float32
						err=rows3.Scan(&tc)
						if err!=nil{
							panic(err)
						}
						tc=tc*38
						switch  {
						case tc<160:
							Fscore2=0
						case tc>=160&&tc<200:
							Fscore2=2
						case tc>=200&&tc<240:
							Fscore2=3
						case tc>=240&&tc<280:
							Fscore2=4
						case tc>=280:
							Fscore2=5
						}
					}
					for rows4.Next(){
						var smoke string
						err=rows4.Scan(&smoke)
						if err!=nil{
							panic(err)
						}
						switch smoke {
						case "吸烟":
							Fscore3=3
						case "不吸烟":
							Fscore3=0
						}
					}
				case age>=60&&age<65:
					Fscore1=10
					for rows3.Next(){
						var tc float32
						err=rows3.Scan(&tc)
						if err!=nil{
							panic(err)
						}
						tc=tc*38
						switch  {
						case tc<160:
							Fscore2=0
						case tc>=160&&tc<200:
							Fscore2=1
						case tc>=200&&tc<240:
							Fscore2=1
						case tc>=240&&tc<280:
							Fscore2=2
						case tc>=280:
							Fscore2=3
						}
					}
					for rows4.Next(){
						var smoke string
						err=rows4.Scan(&smoke)
						if err!=nil{
							panic(err)
						}
						switch smoke {
						case "吸烟":
							Fscore3=1
						case "不吸烟":
							Fscore3=0
						}
					}
				case age>=65&&age<70:
					Fscore1=11
					for rows3.Next(){
						var tc float32
						err=rows3.Scan(&tc)
						if err!=nil{
							panic(err)
						}
						tc=tc*38
						switch  {
						case tc<160:
							Fscore2=0
						case tc>=160&&tc<200:
							Fscore2=1
						case tc>=200&&tc<240:
							Fscore2=1
						case tc>=240&&tc<280:
							Fscore2=2
						case tc>=280:
							Fscore2=3
						}
					}
					for rows4.Next(){
						var smoke string
						err=rows4.Scan(&smoke)
						if err!=nil{
							panic(err)
						}
						switch smoke {
						case "吸烟":
							Fscore3=1
						case "不吸烟":
							Fscore3=0
						}
					}
				case age>=70&&age<75:
					Fscore1=12
					for rows3.Next(){
						var tc float32
						err=rows3.Scan(&tc)
						if err!=nil{
							panic(err)
						}
						tc=tc*38
						switch  {
						case tc<160:
							Fscore2=0
						case tc>=160&&tc<200:
							Fscore2=0
						case tc>=200&&tc<240:
							Fscore2=0
						case tc>=240&&tc<280:
							Fscore2=1
						case tc>=280:
							Fscore2=1
						}
					}
					for rows4.Next(){
						var smoke string
						err=rows4.Scan(&smoke)
						if err!=nil{
							panic(err)
						}
						switch smoke {
						case "吸烟":
							Fscore3=1
						case "不吸烟":
							Fscore3=0
						}
					}
				case age>=75&&age<80:
					Fscore1=13
					for rows3.Next(){
						var tc float32
						err=rows3.Scan(&tc)
						if err!=nil{
							panic(err)
						}
						tc=tc*38
						switch  {
						case tc<160:
							Fscore2=0
						case tc>=160&&tc<200:
							Fscore2=0
						case tc>=200&&tc<240:
							Fscore2=0
						case tc>=240&&tc<280:
							Fscore2=1
						case tc>=280:
							Fscore2=1
						}
					}
					for rows4.Next(){
						var smoke string
						err=rows4.Scan(&smoke)
						if err!=nil{
							panic(err)
						}
						switch smoke {
						case "吸烟":
							Fscore3=1
						case "不吸烟":
							Fscore3=0
						}
					}
				}
			}//年龄 得分1 血总胆固醇 得分2 是否吸烟 得分3
	/*		for rows2.Next(){
			var age int
			err=rows2.Scan(&age)
			if err!=nil{
				panic(err)
			}
			switch {
			case age>=20&&age<40:
				for rows3.Next(){
					var tc float32
					err=rows3.Scan(&tc)
					if err!=nil{
						panic(err)
					}
					tc=tc*18
					switch {
					case tc<160:
						Fscore2=0
					case tc>=160&&tc<200:
						Fscore2=4
					case tc>=200&&tc<240:
						Fscore2=7
					case tc>=240&&tc<280:
						Fscore2=9
					case tc>=280:
						Fscore2=11
					}
				}
			case age>=40&&age<50:
				for rows3.Next(){
					var tc float32
					err=rows3.Scan(&tc)
					if err!=nil{
						panic(err)
					}
					tc=tc*18
					switch  {
					case tc<160:
						Fscore2=0
					case tc>=160&&tc<200:
						Fscore2=3
					case tc>=200&&tc<240:
						Fscore2=5
					case tc>=240&&tc<280:
						Fscore2=6
					case tc>=280:
						Fscore2=8
					}
				}
			case age>=50&&age<60:
				for rows3.Next(){
					var tc float32
					err=rows3.Scan(&tc)
					if err!=nil{
						panic(err)
					}
					tc=tc*18
					switch  {
					case tc<160:
						Fscore2=0
					case tc>=160&&tc<200:
						Fscore2=2
					case tc>=200&&tc<240:
						Fscore2=3
					case tc>=240&&tc<280:
						Fscore2=4
					case tc>=280:
						Fscore2=5
					}
				}
			case age>=60&&age<70:
				for rows3.Next(){
					var tc float32
					err=rows3.Scan(&tc)
					if err!=nil{
						panic(err)
					}
					tc=tc*18
					switch  {
					case tc<160:
						Fscore2=0
					case tc>=160&&tc<200:
						Fscore2=1
					case tc>=200&&tc<240:
						Fscore2=1
					case tc>=240&&tc<280:
						Fscore2=2
					case tc>=280:
						Fscore2=3
					}
				}
			case age>=70&&age<80:
				for rows3.Next(){
					var tc float32
					err=rows3.Scan(&tc)
					if err!=nil{
						panic(err)
					}
					tc=tc*18
					switch  {
					case tc<160:
						Fscore2=0
					case tc>=160&&tc<200:
						Fscore2=0
					case tc>=200&&tc<240:
						Fscore2=0
					case tc>=240&&tc<280:
						Fscore2=1
					case tc>=280:
						Fscore2=1
					}
				}
			}

		}//血总胆固醇 得分2 */
	/*		for rows2.Next(){
				var age int
				err=rows2.Scan(&age)
				if err!=nil{
					panic(err)
				}
				switch {
				case age>=20&&age<40:
					for rows4.Next(){
						var smoke string
						err=rows4.Scan(&smoke)
						if err!=nil{
							panic(err)
						}
						switch smoke {
						case "吸烟":
							Fscore3=8
						case "不吸烟":
							Fscore3=0
						}
					}
				case age>=40&&age<50:
					for rows4.Next(){
						var smoke string
						err=rows4.Scan(&smoke)
						if err!=nil{
							panic(err)
						}
						switch smoke {
						case "吸烟":
							Fscore3=5
						case "不吸烟":
							Fscore3=0
						}
					}
				case age>=50&&age<60:
					for rows4.Next(){
						var smoke string
						err=rows4.Scan(&smoke)
						if err!=nil{
							panic(err)
						}
						switch smoke {
						case "吸烟":
							Fscore3=3
						case "不吸烟":
							Fscore3=0
						}
					}
				case age>=60&&age<70:
					for rows4.Next(){
						var smoke string
						err=rows4.Scan(&smoke)
						if err!=nil{
							panic(err)
						}
						switch smoke {
						case "吸烟":
							Fscore3=1
						case "不吸烟":
							Fscore3=0
						}
					}
				case age>=70&&age<80:
					for rows4.Next(){
						var smoke string
						err=rows4.Scan(&smoke)
						if err!=nil{
							panic(err)
						}
						switch smoke {
						case "吸烟":
							Fscore3=1
						case "不吸烟":
							Fscore3=0
						}
					}
				}
			}//是否吸烟 得分3 */
			for rows5.Next(){
				var Hdlc float32

				err = rows5.Scan(&Hdlc)
				if err != nil {
					panic(err)
				}
				Hdlc=Hdlc*38
				switch  {
				case Hdlc>=60:
					Fscore4=-1
				case Hdlc>=50&&Hdlc<60:
					Fscore4=0
				case Hdlc>=40&&Hdlc<50:
					Fscore4=1
				case Hdlc<40:
					Fscore4=2
				}
			}//高密度脂蛋白胆固醇 得分4
			for rows6.Next(){
				var SBP int

				err = rows6.Scan(&SBP)
				if err != nil {
					panic(err)
				}
				switch {
				case SBP<120:
					Fscore5=0
				case SBP>=120&&SBP<130:
					Fscore5=0
				case SBP>=130&&SBP<140:
					Fscore5=1
				case SBP>=140&&SBP<160:
					Fscore5=1
				case SBP>=160:
					Fscore5=2
				}
			}//收缩压 得分5
		case "女":
			for rows2.Next() {
				var age int

				err = rows2.Scan(&age)
				if err != nil {
					panic(err)
				}
				switch {
				case age >= 20 && age < 35:
					Fscore1=-7
					for rows3.Next(){
						var tc float32
						err=rows3.Scan(&tc)
						if err!=nil{
							panic(err)
						}
						tc=tc*38
						switch {
						case tc<160:
							Fscore2=0
						case tc>=160&&tc<200:
							Fscore2=4
						case tc>=200&&tc<240:
							Fscore2=8
						case tc>=240&&tc<280:
							Fscore2=11
						case tc>=280:
							Fscore2=13
						}
					}
					for rows4.Next(){
						var smoke string
						err=rows4.Scan(&smoke)
						if err!=nil{
							panic(err)
						}
						switch smoke {
						case "吸烟":
							Fscore3=9
						case "不吸烟":
							Fscore3=0
						}
					}
				case age>=35&&age<40:
					Fscore1=-3
					for rows3.Next(){
						var tc float32
						err=rows3.Scan(&tc)
						if err!=nil{
							panic(err)
						}
						tc=tc*38
						switch {
						case tc<160:
							Fscore2=0
						case tc>=160&&tc<200:
							Fscore2=4
						case tc>=200&&tc<240:
							Fscore2=8
						case tc>=240&&tc<280:
							Fscore2=11
						case tc>=280:
							Fscore2=13
						}
					}
					for rows4.Next(){
						var smoke string
						err=rows4.Scan(&smoke)
						if err!=nil{
							panic(err)
						}
						switch smoke {
						case "吸烟":
							Fscore3=9
						case "不吸烟":
							Fscore3=0
						}
					}
				case age>=40&&age<45:
					Fscore1=0
					for rows3.Next(){
						var tc float32
						err=rows3.Scan(&tc)
						if err!=nil{
							panic(err)
						}
						tc=tc*38
						switch  {
						case tc<160:
							Fscore2=0
						case tc>=160&&tc<200:
							Fscore2=3
						case tc>=200&&tc<240:
							Fscore2=6
						case tc>=240&&tc<280:
							Fscore2=8
						case tc>=280:
							Fscore2=10
						}
					}
					for rows4.Next(){
						var smoke string
						err=rows4.Scan(&smoke)
						if err!=nil{
							panic(err)
						}
						switch smoke {
						case "吸烟":
							Fscore3=7
						case "不吸烟":
							Fscore3=0
						}
					}
				case age>=45&&age<50:
					Fscore1=3
					for rows3.Next(){
						var tc float32
						err=rows3.Scan(&tc)
						if err!=nil{
							panic(err)
						}
						tc=tc*38
						switch  {
						case tc<160:
							Fscore2=0
						case tc>=160&&tc<200:
							Fscore2=3
						case tc>=200&&tc<240:
							Fscore2=6
						case tc>=240&&tc<280:
							Fscore2=8
						case tc>=280:
							Fscore2=10
						}
					}
					for rows4.Next(){
						var smoke string
						err=rows4.Scan(&smoke)
						if err!=nil{
							panic(err)
						}
						switch smoke {
						case "吸烟":
							Fscore3=7
						case "不吸烟":
							Fscore3=0
						}
					}
				case age>=50&&age<55:
					Fscore1=6
					for rows3.Next(){
						var tc float32
						err=rows3.Scan(&tc)
						if err!=nil{
							panic(err)
						}
						tc=tc*38
						switch  {
						case tc<160:
							Fscore2=0
						case tc>=160&&tc<200:
							Fscore2=2
						case tc>=200&&tc<240:
							Fscore2=4
						case tc>=240&&tc<280:
							Fscore2=5
						case tc>=280:
							Fscore2=7
						}
					}
					for rows4.Next(){
						var smoke string
						err=rows4.Scan(&smoke)
						if err!=nil{
							panic(err)
						}
						switch smoke {
						case "吸烟":
							Fscore3=4
						case "不吸烟":
							Fscore3=0
						}
					}
				case age>=55&&age<60:
					Fscore1=8
					for rows3.Next(){
						var tc float32
						err=rows3.Scan(&tc)
						if err!=nil{
							panic(err)
						}
						tc=tc*38
						switch  {
						case tc<160:
							Fscore2=0
						case tc>=160&&tc<200:
							Fscore2=2
						case tc>=200&&tc<240:
							Fscore2=4
						case tc>=240&&tc<280:
							Fscore2=5
						case tc>=280:
							Fscore2=7
						}
					}
					for rows4.Next(){
						var smoke string
						err=rows4.Scan(&smoke)
						if err!=nil{
							panic(err)
						}
						switch smoke {
						case "吸烟":
							Fscore3=4
						case "不吸烟":
							Fscore3=0
						}
					}
				case age>=60&&age<65:
					Fscore1=10
					for rows3.Next(){
						var tc float32
						err=rows3.Scan(&tc)
						if err!=nil{
							panic(err)
						}
						tc=tc*38
						switch  {
						case tc<160:
							Fscore2=0
						case tc>=160&&tc<200:
							Fscore2=1
						case tc>=200&&tc<240:
							Fscore2=2
						case tc>=240&&tc<280:
							Fscore2=3
						case tc>=280:
							Fscore2=4
						}
					}
					for rows4.Next(){
						var smoke string
						err=rows4.Scan(&smoke)
						if err!=nil{
							panic(err)
						}
						switch smoke {
						case "吸烟":
							Fscore3=2
						case "不吸烟":
							Fscore3=0
						}
					}
				case age>=65&&age<70:
					Fscore1=12
					for rows3.Next(){
						var tc float32
						err=rows3.Scan(&tc)
						if err!=nil{
							panic(err)
						}
						tc=tc*38
						switch  {
						case tc<160:
							Fscore2=0
						case tc>=160&&tc<200:
							Fscore2=1
						case tc>=200&&tc<240:
							Fscore2=2
						case tc>=240&&tc<280:
							Fscore2=3
						case tc>=280:
							Fscore2=4
						}
					}
					for rows4.Next(){
						var smoke string
						err=rows4.Scan(&smoke)
						if err!=nil{
							panic(err)
						}
						switch smoke {
						case "吸烟":
							Fscore3=2
						case "不吸烟":
							Fscore3=0
						}
					}
				case age>=70&&age<75:
					Fscore1=14
					for rows3.Next(){
						var tc float32
						err=rows3.Scan(&tc)
						if err!=nil{
							panic(err)
						}
						tc=tc*38
						switch  {
						case tc<160:
							Fscore2=0
						case tc>=160&&tc<200:
							Fscore2=1
						case tc>=200&&tc<240:
							Fscore2=1
						case tc>=240&&tc<280:
							Fscore2=2
						case tc>=280:
							Fscore2=2
						}
					}
					for rows4.Next(){
						var smoke string
						err=rows4.Scan(&smoke)
						if err!=nil{
							panic(err)
						}
						switch smoke {
						case "吸烟":
							Fscore3=1
						case "不吸烟":
							Fscore3=0
						}
					}
				case age>=75&&age<80:
					Fscore1=16
					for rows3.Next(){
						var tc float32
						err=rows3.Scan(&tc)
						if err!=nil{
							panic(err)
						}
						tc=tc*38
						switch  {
						case tc<160:
							Fscore2=0
						case tc>=160&&tc<200:
							Fscore2=1
						case tc>=200&&tc<240:
							Fscore2=1
						case tc>=240&&tc<280:
							Fscore2=2
						case tc>=280:
							Fscore2=2
						}
					}
					for rows4.Next(){
						var smoke string
						err=rows4.Scan(&smoke)
						if err!=nil{
							panic(err)
						}
						switch smoke {
						case "吸烟":
							Fscore3=1
						case "不吸烟":
							Fscore3=0
						}
					}
				}
			}//年龄 得分1 血总胆固醇 得分2 是否吸烟 得分3
			for rows5.Next(){
				var Hdlc float32

				err = rows5.Scan(&Hdlc)
				if err != nil {
					panic(err)
				}
				Hdlc=Hdlc*38
				switch  {
				case Hdlc>=60:
					Fscore4=-1
				case Hdlc>=50&&Hdlc<60:
					Fscore4=0
				case Hdlc>=40&&Hdlc<50:
					Fscore4=1
				case Hdlc<40:
					Fscore4=2
				}
			}//高密度脂蛋白胆固醇 得分4
			for rows6.Next(){
				var SBP int

				err = rows6.Scan(&SBP)
				if err != nil {
					panic(err)
				}
				switch {
				case SBP<120:
					Fscore5=0
				case SBP>=120&&SBP<130:
					Fscore5=1
				case SBP>=130&&SBP<140:
					Fscore5=2
				case SBP>=140&&SBP<160:
					Fscore5=3
				case SBP>=160:
					Fscore5=4
				}
			}//收缩压 得分5
		}
		Fscore=Fscore1+Fscore2+Fscore3+Fscore4+Fscore5
	}

	fmt.Println(Fscore1,Fscore2,Fscore3,Fscore4,Fscore5)
	fmt.Println("弗兰明翰总得分为",Fscore)
	db.Close()
}

func sFrisk() {
	db, err := sql.Open("mysql", "huang:zjubio2019@tcp(mysql.zjubiomedit.com:3306)/cdms_biomedit?charset=utf8")
	CheckErr(err)
	rows1, _ := db.Query("select sex from eval_userdata ORDER BY id DESC LIMIT 1")
	for rows1.Next() {
		var sex string

		err = rows1.Scan(&sex)
		if err != nil {
			panic(err)
		}
		switch sex {
		case "男":
			switch {
			case Fscore < 0:
				Frisk = 0
			case Fscore >= 0 && Fscore < 5:
				Frisk = 1
			case Fscore >= 5 && Fscore < 7:
				Frisk = 2
			case Fscore == 7:
				Frisk = 3
			case Fscore == 8:
				Frisk = 4
			case Fscore == 9:
				Frisk = 5
			case Fscore == 10:
				Frisk = 6
			case Fscore == 11:
				Frisk = 8
			case Fscore == 12:
				Frisk = 10
			case Fscore == 13:
				Frisk = 12
			case Fscore == 14:
				Frisk = 16
			case Fscore == 15:
				Frisk = 20
			case Fscore == 16:
				Frisk = 25
			case Fscore >= 17:
				Frisk = 30
			}
		case "女":
			switch {
			case Fscore < 9:
				Frisk = 0
			case Fscore >= 9 && Fscore < 13:
				Frisk = 1
			case Fscore >= 13 && Fscore < 15:
				Frisk = 2
			case Fscore == 15:
				Frisk = 3
			case Fscore == 16:
				Frisk = 4
			case Fscore == 17:
				Frisk = 5
			case Fscore == 18:
				Frisk = 6
			case Fscore == 19:
				Frisk = 8
			case Fscore == 20:
				Frisk = 11
			case Fscore == 21:
				Frisk = 14
			case Fscore == 22:
				Frisk = 17
			case Fscore == 23:
				Frisk = 22
			case Fscore == 24:
				Frisk = 27
			case Fscore >= 25:
				Frisk = 30
			}
		}
	}
	db.Close()
	fmt.Println("弗兰明翰风险为：",Frisk)
}

func evaluateW(w http.ResponseWriter){

	calculateW()
	fmt.Println("WHO风险为：",Wrisk)
	fmt.Fprintf(w,  strconv.Itoa(Wrisk))
}

func calculateW(){
	db, err := sql.Open("mysql", "huang:zjubio2019@tcp(mysql.zjubiomedit.com:3306)/cdms_biomedit?charset=utf8")
	CheckErr(err)
	rows1, _ := db.Query("select dia from eval_userdata ORDER BY id DESC LIMIT 1")
	rows2, _ := db.Query("select sex from eval_userdata ORDER BY id DESC LIMIT 1")
	rows3, _ := db.Query("select smoke from eval_userdata ORDER BY id DESC LIMIT 1")
	rows4,_:=db.Query("select age from eval_userdata ORDER BY id DESC LIMIT 1")
	rows5,_:=db.Query("select sbp from eval_userdata ORDER BY id DESC LIMIT 1")
	rows6,_:=db.Query("select tc from eval_userdata ORDER BY id DESC LIMIT 1")

	for rows1.Next(){
		var diabetes string
		err = rows1.Scan(&diabetes)
		if err != nil {
			panic(err)
		}
		switch diabetes{
		case "有糖尿病":
			for rows2.Next(){
				var sex string
				err=rows2.Scan(&sex)
				if err!=nil{ panic(err)}
				switch sex {
				case "男":
					for rows3.Next(){
						var smoke string
						err=rows3.Scan(&smoke)
						if err!=nil{ panic(err)}
						switch  smoke{
						case "不吸烟":
							for rows4.Next(){
								var age int
								err=rows4.Scan(&age)
								if err!=nil{ panic(err)}
								switch {
								case age>=70:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=40
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=30
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=30
												case tbc>=6&&tbc<7:
													Wrisk=30
												case tbc>=5&&tbc<6:
													Wrisk=20
												case tbc<5:
													Wrisk=20
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=20
												case tbc>=7&&tbc<8:
													Wrisk=20
												case tbc>=6&&tbc<7:
													Wrisk=10
												case tbc>=5&&tbc<6:
													Wrisk=10
												case tbc<5:
													Wrisk=10
												}
											}
										}
									}
								case age>=60&&age<70:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=40
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=30
												case tbc>=5&&tbc<6:
													Wrisk=30
												case tbc<5:
													Wrisk=20
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=30
												case tbc>=7&&tbc<8:
													Wrisk=20
												case tbc>=6&&tbc<7:
													Wrisk=10
												case tbc>=5&&tbc<6:
													Wrisk=10
												case tbc<5:
													Wrisk=10
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=10
												case tbc>=7&&tbc<8:
													Wrisk=10
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								case age>=50&&age<60:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=30
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=30
												case tbc>=6&&tbc<7:
													Wrisk=10
												case tbc>=5&&tbc<6:
													Wrisk=10
												case tbc<5:
													Wrisk=10
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=10
												case tbc>=7&&tbc<8:
													Wrisk=10
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=10
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								case age>=40&&age<50:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=30
												case tbc<5:
													Wrisk=20
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=30
												case tbc>=7&&tbc<8:
													Wrisk=20
												case tbc>=6&&tbc<7:
													Wrisk=10
												case tbc>=5&&tbc<6:
													Wrisk=10
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=10
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=10
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								}
							}
						case "吸烟":
							for rows4.Next(){
								var age int
								err=rows4.Scan(&age)
								if err!=nil{ panic(err)}
								switch {
								case age>=70:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=40
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=40
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=30
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=30
												case tbc>=6&&tbc<7:
													Wrisk=20
												case tbc>=5&&tbc<6:
													Wrisk=20
												case tbc<5:
													Wrisk=10
												}
											}
										}
									}
								case age>=60&&age<70:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=40
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=30
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=30
												case tbc>=5&&tbc<6:
													Wrisk=20
												case tbc<5:
													Wrisk=10
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=30
												case tbc>=7&&tbc<8:
													Wrisk=20
												case tbc>=6&&tbc<7:
													Wrisk=10
												case tbc>=5&&tbc<6:
													Wrisk=10
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								case age>=50&&age<60:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=40
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=30
												case tbc<5:
													Wrisk=20
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=30
												case tbc>=7&&tbc<8:
													Wrisk=20
												case tbc>=6&&tbc<7:
													Wrisk=10
												case tbc>=5&&tbc<6:
													Wrisk=10
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=10
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								case age>=40&&age<50:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=40
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=30
												case tbc>=6&&tbc<7:
													Wrisk=30
												case tbc>=5&&tbc<6:
													Wrisk=20
												case tbc<5:
													Wrisk=10
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=20
												case tbc>=7&&tbc<8:
													Wrisk=10
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=10
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								}
							}
						}
					}
				case "女":
					for rows3.Next(){
						var smoke string
						err=rows3.Scan(&smoke)
						if err!=nil{ panic(err)}
						switch  smoke{
						case "不吸烟":
							for rows4.Next(){
								var age int
								err=rows4.Scan(&age)
								if err!=nil{ panic(err)}
								switch {
								case age>=70:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=40
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=30
												case tbc>=5&&tbc<6:
													Wrisk=20
												case tbc<5:
													Wrisk=20
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=30
												case tbc>=7&&tbc<8:
													Wrisk=20
												case tbc>=6&&tbc<7:
													Wrisk=20
												case tbc>=5&&tbc<6:
													Wrisk=10
												case tbc<5:
													Wrisk=10
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=10
												case tbc>=7&&tbc<8:
													Wrisk=10
												case tbc>=6&&tbc<7:
													Wrisk=10
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								case age>=60&&age<70:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=40
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=30
												case tbc>=5&&tbc<6:
													Wrisk=20
												case tbc<5:
													Wrisk=20
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=30
												case tbc>=7&&tbc<8:
													Wrisk=20
												case tbc>=6&&tbc<7:
													Wrisk=10
												case tbc>=5&&tbc<6:
													Wrisk=10
												case tbc<5:
													Wrisk=10
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=10
												case tbc>=7&&tbc<8:
													Wrisk=10
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								case age>=50&&age<60:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=30
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=20
												case tbc>=6&&tbc<7:
													Wrisk=10
												case tbc>=5&&tbc<6:
													Wrisk=10
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=20
												case tbc>=7&&tbc<8:
													Wrisk=10
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=1
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								case age>=40&&age<50:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=20
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=20
												case tbc>=6&&tbc<7:
													Wrisk=10
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=20
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=1
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								}
							}
						case "吸烟":
							for rows4.Next(){
								var age int
								err=rows4.Scan(&age)
								if err!=nil{ panic(err)}
								switch {
								case age>=70:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=40
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=30
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=30
												case tbc>=5&&tbc<6:
													Wrisk=20
												case tbc<5:
													Wrisk=20
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=20
												case tbc>=7&&tbc<8:
													Wrisk=20
												case tbc>=6&&tbc<7:
													Wrisk=20
												case tbc>=5&&tbc<6:
													Wrisk=10
												case tbc<5:
													Wrisk=10
												}
											}
										}
									}
								case age>=60&&age<70:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=40
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=30
												case tbc<5:
													Wrisk=20
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=20
												case tbc>=6&&tbc<7:
													Wrisk=20
												case tbc>=5&&tbc<6:
													Wrisk=10
												case tbc<5:
													Wrisk=10
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=20
												case tbc>=7&&tbc<8:
													Wrisk=10
												case tbc>=6&&tbc<7:
													Wrisk=10
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								case age>=50&&age<60:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=40
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=30
												case tbc>=5&&tbc<6:
													Wrisk=20
												case tbc<5:
													Wrisk=10
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=30
												case tbc>=7&&tbc<8:
													Wrisk=20
												case tbc>=6&&tbc<7:
													Wrisk=10
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=10
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								case age>=40&&age<50:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=40
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=20
												case tbc>=5&&tbc<6:
													Wrisk=20
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=30
												case tbc>=7&&tbc<8:
													Wrisk=10
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=1
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		case "无糖尿病":
			for rows2.Next(){
				var sex string
				err=rows2.Scan(&sex)
				if err!=nil{ panic(err)}
				switch sex {
				case "男":
					for rows3.Next(){
						var smoke string
						err=rows3.Scan(&smoke)
						if err!=nil{ panic(err)}
						switch  smoke{
						case "不吸烟":
							for rows4.Next(){
								var age int
								err=rows4.Scan(&age)
								if err!=nil{ panic(err)}
								switch {
								case age>=70:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=30
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=30
												case tbc>=7&&tbc<8:
													Wrisk=30
												case tbc>=6&&tbc<7:
													Wrisk=20
												case tbc>=5&&tbc<6:
													Wrisk=20
												case tbc<5:
													Wrisk=10
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=20
												case tbc>=7&&tbc<8:
													Wrisk=10
												case tbc>=6&&tbc<7:
													Wrisk=10
												case tbc>=5&&tbc<6:
													Wrisk=10
												case tbc<5:
													Wrisk=10
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=10
												case tbc>=7&&tbc<8:
													Wrisk=10
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								case age>=60&&age<70:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=30
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=30
												case tbc>=7&&tbc<8:
													Wrisk=20
												case tbc>=6&&tbc<7:
													Wrisk=20
												case tbc>=5&&tbc<6:
													Wrisk=10
												case tbc<5:
													Wrisk=10
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=10
												case tbc>=7&&tbc<8:
													Wrisk=10
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=1
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								case age>=50&&age<60:
									for rows5.Next(){
										var sbp float32
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc int
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=30
												case tbc>=5&&tbc<6:
													Wrisk=20
												case tbc<5:
													Wrisk=20
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=30
												case tbc>=7&&tbc<8:
													Wrisk=10
												case tbc>=6&&tbc<7:
													Wrisk=10
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=10
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=1
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								case age>=40&&age<50:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=30
												case tbc>=5&&tbc<6:
													Wrisk=20
												case tbc<5:
													Wrisk=10
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=20
												case tbc>=7&&tbc<8:
													Wrisk=10
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=10
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=1
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								}
							}
						case "吸烟":
							for rows4.Next(){
								var age int
								err=rows4.Scan(&age)
								if err!=nil{ panic(err)}
								switch {
								case age>=70:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=40
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=30
												case tbc<5:
													Wrisk=20
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=30
												case tbc>=7&&tbc<8:
													Wrisk=30
												case tbc>=6&&tbc<7:
													Wrisk=30
												case tbc>=5&&tbc<6:
													Wrisk=20
												case tbc<5:
													Wrisk=10
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=20
												case tbc>=7&&tbc<8:
													Wrisk=10
												case tbc>=6&&tbc<7:
													Wrisk=10
												case tbc>=5&&tbc<6:
													Wrisk=10
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								case age>=60&&age<70:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=40
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=30
												case tbc>=5&&tbc<6:
													Wrisk=20
												case tbc<5:
													Wrisk=20
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=30
												case tbc>=7&&tbc<8:
													Wrisk=20
												case tbc>=6&&tbc<7:
													Wrisk=10
												case tbc>=5&&tbc<6:
													Wrisk=10
												case tbc<5:
													Wrisk=10
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=10
												case tbc>=7&&tbc<8:
													Wrisk=10
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								case age>=50&&age<60:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=30
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=30
												case tbc>=6&&tbc<7:
													Wrisk=20
												case tbc>=5&&tbc<6:
													Wrisk=10
												case tbc<5:
													Wrisk=10
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=20
												case tbc>=7&&tbc<8:
													Wrisk=10
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=10
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								case age>=40&&age<50:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=30
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=20
												case tbc>=6&&tbc<7:
													Wrisk=10
												case tbc>=5&&tbc<6:
													Wrisk=10
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=20
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=1
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								}
							}
						}
					}
				case "女":
					for rows3.Next(){
						var smoke string
						err=rows3.Scan(&smoke)
						if err!=nil{ panic(err)}
						switch  smoke{
						case "不吸烟":
							for rows4.Next(){
								var age int
								err=rows4.Scan(&age)
								if err!=nil{ panic(err)}
								switch {
								case age>=70:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=30
												case tbc>=5&&tbc<6:
													Wrisk=30
												case tbc<5:
													Wrisk=20
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=20
												case tbc>=7&&tbc<8:
													Wrisk=10
												case tbc>=6&&tbc<7:
													Wrisk=10
												case tbc>=5&&tbc<6:
													Wrisk=10
												case tbc<5:
													Wrisk=10
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=10
												case tbc>=7&&tbc<8:
													Wrisk=10
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=1
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								case age>=60&&age<70:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=30
												case tbc>=5&&tbc<6:
													Wrisk=20
												case tbc<5:
													Wrisk=20
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=20
												case tbc>=7&&tbc<8:
													Wrisk=10
												case tbc>=6&&tbc<7:
													Wrisk=10
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=10
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=1
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								case age>=50&&age<60:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=30
												case tbc>=5&&tbc<6:
													Wrisk=20
												case tbc<5:
													Wrisk=10
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=20
												case tbc>=7&&tbc<8:
													Wrisk=10
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=10
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=1
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								case age>=40&&age<50:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=30
												case tbc>=5&&tbc<6:
													Wrisk=20
												case tbc<5:
													Wrisk=10
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=20
												case tbc>=7&&tbc<8:
													Wrisk=10
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=10
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=1
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								}
							}
						case "吸烟":
							for rows4.Next(){
								var age int
								err=rows4.Scan(&age)
								if err!=nil{ panic(err)}
								switch {
								case age>=70:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=30
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=30
												case tbc>=6&&tbc<7:
													Wrisk=20
												case tbc>=5&&tbc<6:
													Wrisk=20
												case tbc<5:
													Wrisk=10
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=20
												case tbc>=7&&tbc<8:
													Wrisk=10
												case tbc>=6&&tbc<7:
													Wrisk=10
												case tbc>=5&&tbc<6:
													Wrisk=10
												case tbc<5:
													Wrisk=10
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=10
												case tbc>=7&&tbc<8:
													Wrisk=10
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								case age>=60&&age<70:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=30
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=30
												case tbc>=5&&tbc<6:
													Wrisk=20
												case tbc<5:
													Wrisk=20
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=20
												case tbc>=7&&tbc<8:
													Wrisk=10
												case tbc>=6&&tbc<7:
													Wrisk=10
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=10
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								case age>=50&&age<60:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=40
												case tbc<5:
													Wrisk=30
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=20
												case tbc>=6&&tbc<7:
													Wrisk=10
												case tbc>=5&&tbc<6:
													Wrisk=10
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=20
												case tbc>=7&&tbc<8:
													Wrisk=10
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=1
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								case age>=40&&age<50:
									for rows5.Next(){
										var sbp int
										err=rows5.Scan(&sbp)
										if err!=nil{ panic(err)}
										switch {
										case sbp>=180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=40
												case tbc>=6&&tbc<7:
													Wrisk=40
												case tbc>=5&&tbc<6:
													Wrisk=30
												case tbc<5:
													Wrisk=20
												}
											}
										case sbp>=160&&sbp<180:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=40
												case tbc>=7&&tbc<8:
													Wrisk=20
												case tbc>=6&&tbc<7:
													Wrisk=10
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp>=140&&sbp<160:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=20
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										case sbp<140:
											for rows6.Next(){
												var tbc float32
												err=rows6.Scan(&tbc)
												if err!=nil{ panic(err)}
												switch {
												case tbc>=8:
													Wrisk=1
												case tbc>=7&&tbc<8:
													Wrisk=1
												case tbc>=6&&tbc<7:
													Wrisk=1
												case tbc>=5&&tbc<6:
													Wrisk=1
												case tbc<5:
													Wrisk=1
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	db.Close()
}

func evaluateI(w http.ResponseWriter){
	calculateI()
	sIrisk()
	db, err := sql.Open("mysql", "huang:zjubio2019@tcp(mysql.zjubiomedit.com:3306)/cdms_biomedit?charset=utf8")
	CheckErr(err)
	db.Exec("insert into eval_results(patientId,time,WHO,Frammingham,ICVD) values (?,?,?,?,?)",patientId,measureTime,Wrisk,Frisk,Irisk)
	db.Close()
	fmt.Fprintf(w,  FloatToString(Irisk))
}

func FloatToString(input_num float64) string {
	return strconv.FormatFloat(float64(input_num), 'f', 1, 64)
}

func calculateI(){
	db, err := sql.Open("mysql", "huang:zjubio2019@tcp(mysql.zjubiomedit.com:3306)/cdms_biomedit?charset=utf8")
	CheckErr(err)
	rows1, _ := db.Query("select sex from eval_userdata ORDER BY id DESC LIMIT 1")
	rows2, _ := db.Query("select age from eval_userdata ORDER BY id DESC LIMIT 1")
	rows3, _ := db.Query("select sbp from eval_userdata ORDER BY id DESC LIMIT 1")
	rows4,_:=db.Query("select height from eval_userdata ORDER BY id DESC LIMIT 1")
	rows5,_:=db.Query("select weight from eval_userdata ORDER BY id DESC LIMIT 1")
	rows6,_:=db.Query("select tc from eval_userdata ORDER BY id DESC LIMIT 1")
	rows7,_:=db.Query("select smoke from eval_userdata ORDER BY id DESC LIMIT 1")
	rows8,_:=db.Query("select dia from eval_userdata ORDER BY id DESC LIMIT 1")

	for rows1.Next(){
		var sex string
		err = rows1.Scan(&sex)
		if err != nil {panic(err)}
		switch sex {
		case "男":
			for rows2.Next(){
				var age int
				err = rows2.Scan(&age)
				if err != nil {panic(err)}
				switch {
				case age<40:
					Iscore1=0
				case age>=40&&age<45:
					Iscore1=1
				case age>=45&&age<50:
					Iscore1=2
				case age>=50&&age<55:
					Iscore1=3
				case age>=55&&age<60:
					Iscore1=4
				case age>=60&&age<65:
					Iscore1=5
				case age>=65&&age<70:
					Iscore1=6
				case age>=70&&age<75:
					Iscore1=7
				case age>=75&&age<80:
					Iscore1=8
				case age>=80&&age<85:
					Iscore1=9
				case age>=85&&age<90:
					Iscore1=10
				case age>=90&&age<100:
					Iscore1=11
				}
			}//年龄 得分1
			for rows3.Next(){
				var sbp int
				err = rows3.Scan(&sbp)
				if err != nil {panic(err)}
				switch {
				case sbp<120:
					Iscore2=-2
				case sbp>=120&&sbp<130:
					Iscore2=0
				case sbp>=130&&sbp<140:
					Iscore2=1
				case sbp>=140&&sbp<160:
					Iscore2=2
				case sbp>=160&&sbp<180:
					Iscore2=5
				case sbp>=180:
					Iscore2=8
				}
			}//血压 得分2
		    for rows4.Next(){
				var height,weight float32
				err = rows4.Scan(&height)
				if err != nil {panic(err)}
				height = height/100
				for rows5.Next(){
					err = rows5.Scan(&weight)
					if err != nil {panic(err)}
				}
				var bmi float32
				bmi = weight/(height*height)
				fmt.Println("BMI为：",bmi)
				switch {
				case bmi<24:
					Iscore3=0
				case bmi>=24&&bmi<28:
					Iscore3=1
				case bmi>=28:
					Iscore3=2
				}
			}//BMI  得分3
			for rows6.Next(){
				var tbc float32
				err =rows6.Scan(&tbc)
				if err!=nil{panic(err)}
				switch {
				case tbc<=5:
					Iscore4=0
				case tbc>5:
					Iscore4=1
				}
			}//血总胆固醇 得分4
			for rows7.Next(){
				var smoke string
				err =rows7.Scan(&smoke)
				if err!=nil{panic(err)}
				switch smoke {
				case "吸烟":
					Iscore5=2
				case "不吸烟":
					Iscore5=0
				}
			}//吸烟 得分5
			for rows8.Next(){
				var diabetes string
				err =rows8.Scan(&diabetes)
				if err!=nil{panic(err)}
				switch diabetes {
				case "有糖尿病":
					Iscore6=1
				case "无糖尿病":
					Iscore6=0
				}
			}//糖尿病 得分6
		case"女":
			for rows2.Next(){
				var age int
				err = rows2.Scan(&age)
				if err != nil {panic(err)}
				switch {
				case age<40:
					Iscore1=0
				case age>=40&&age<45:
					Iscore1=1
				case age>=45&&age<50:
					Iscore1=2
				case age>=50&&age<55:
					Iscore1=3
				case age>=55&&age<60:
					Iscore1=4
				case age>=60&&age<65:
					Iscore1=5
				case age>=65&&age<70:
					Iscore1=6
				case age>=70&&age<75:
					Iscore1=7
				case age>=75&&age<80:
					Iscore1=8
				case age>=80&&age<85:
					Iscore1=9
				case age>=85&&age<90:
					Iscore1=10
				case age>=90&&age<100:
					Iscore1=11
				}
			}
			for rows3.Next(){
				var sbp int
				err = rows3.Scan(&sbp)
				if err != nil {panic(err)}
				switch {
				case sbp<120:
					Iscore2=-2
				case sbp>=120&&sbp<130:
					Iscore2=0
				case sbp>=130&&sbp<140:
					Iscore2=1
				case sbp>=140&&sbp<160:
					Iscore2=2
				case sbp>=160&&sbp<180:
					Iscore2=3
				case sbp>=180:
					Iscore2=4
				}
			}
			for rows4.Next(){
				var height,weight float32
				err = rows4.Scan(&height)
				if err != nil {panic(err)}
				height = height/100
				for rows5.Next(){
					err = rows5.Scan(&weight)
					if err != nil {panic(err)}
				}
				var bmi float32
				bmi = weight/(height*height)
				fmt.Println("BMI为：",bmi)
				switch {
				case bmi<24:
					Iscore3=0
				case bmi>=24&&bmi<28:
					Iscore3=1
				case bmi>=28:
					Iscore3=2
				}
			}
			for rows6.Next(){
				var tbc int
				err =rows6.Scan(&tbc)
				if err!=nil{panic(err)}
				switch {
				case tbc<=5:
					Iscore4=0
				case tbc>5:
					Iscore4=1
				}
			}
			for rows7.Next(){
				var smoke string
				err =rows7.Scan(&smoke)
				if err!=nil{panic(err)}
				switch smoke {
				case "吸烟":
					Iscore5=1
				case "不吸烟":
					Iscore5=0
				}
			}
			for rows8.Next(){
				var diabetes string
				err =rows8.Scan(&diabetes)
				if err!=nil{panic(err)}
				switch diabetes {
				case "有糖尿病":
					Iscore6=2
				case "无糖尿病":
					Iscore6=0
				}
			}
		}
	}
	Iscore=Iscore1+Iscore2+Iscore3+Iscore4+Iscore5+Iscore6

	fmt.Println(Iscore1,Iscore2,Iscore3,Iscore4,Iscore5,Iscore6)
	fmt.Println("ICVD总得分为",Iscore)
	db.Close()
}

func sIrisk() {
	db, err := sql.Open("mysql", "huang:zjubio2019@tcp(mysql.zjubiomedit.com:3306)/cdms_biomedit?charset=utf8")
	CheckErr(err)
	rows1, _ := db.Query("select sex from eval_userdata ORDER BY id DESC LIMIT 1")
	for rows1.Next() {
		var sex string

		err = rows1.Scan(&sex)
		if err != nil {
			panic(err)
		}
		switch sex {
		case "男":
			switch {
			case Iscore<=-1:
				Irisk=0.3
			case Iscore==0:
				Irisk=0.5
			case Iscore==1:
				Irisk=0.6
			case Iscore==2:
				Irisk=0.8
			case Iscore==3:
				Irisk=1.1
			case Iscore==4:
				Irisk=1.5
			case Iscore==5:
				Irisk=2.1
			case Iscore==6:
				Irisk=2.9
			case Iscore==7:
				Irisk=3.9
			case Iscore==8:
				Irisk=5.4
			case Iscore==9:
				Irisk=7.3
			case Iscore==10:
				Irisk=9.7
			case Iscore==11:
				Irisk=12.8
			case Iscore==12:
				Irisk=16.8
			case Iscore==13:
				Irisk=21.7
			case Iscore==14:
				Irisk=27.7
			case Iscore==15:
				Irisk=35.3
			case Iscore==16:
				Irisk=44.3
			case Iscore>=17:
				Irisk=52.6
			}
		case "女":
			switch {
			case Iscore==-2:
				Irisk=0.1
			case Iscore==-1:
				Irisk=0.2
			case Iscore==0:
				Irisk=0.2
			case Iscore==1:
				Irisk=0.2
			case Iscore==2:
				Irisk=0.3
			case Iscore==3:
				Irisk=0.5
			case Iscore==4:
				Irisk=1.5
			case Iscore==5:
				Irisk=2.1
			case Iscore==6:
				Irisk=2.9
			case Iscore==7:
				Irisk=3.9
			case Iscore==8:
				Irisk=5.4
			case Iscore==9:
				Irisk=7.3
			case Iscore==10:
				Irisk=9.7
			case Iscore==11:
				Irisk=12.8
			case Iscore==12:
				Irisk=16.8
			case Iscore==13:
				Irisk=21.7
			case Iscore==14:
				Irisk=27.7
			case Iscore==15:
				Irisk=35.3
			case Iscore==16:
				Irisk=44.3
			case Iscore>=17:
				Irisk=52.6
			}
		}
	}
	db.Close()
	fmt.Println("ICVD风险为：",Irisk)
}

func factor(){
	if Fscore2>0&&Iscore3>0{
		ftbc="TC高"
		tbclab=1
	} else {
		ftbc="TC正常"
		tbclab=0
	}
	if Fscore3>0&&Iscore5>0{
		fsmoke="吸烟"
	} else {
		fsmoke="不吸烟"
	}
	if Fscore4>0{
		fHdlc="HDL-C异常"
		Hdlclab=1
	}else{
		fHdlc="HDL-C正常"
		Hdlclab=0
	}
	if Fscore5>0&&Iscore2>1{
		fsbp="血压过高"
	} else{
		fsbp="血压正常"
	}
	db, err := sql.Open("mysql", "huang:zjubio2019@tcp(mysql.zjubiomedit.com:3306)/cdms_biomedit?charset=utf8")
	CheckErr(err)
	rows1, _ := db.Query("select dia from eval_userdata ORDER BY id DESC LIMIT 1")
	rows2, _ := db.Query("select drink from eval_userdata ORDER BY id DESC LIMIT 1")

	for rows1.Next(){
		var diabetes string
		err = rows1.Scan(&diabetes)
		if err != nil {
			panic(err)
		}
		switch diabetes{
		case "有糖尿病":
			fdiabetes="血糖过高"
		case "无糖尿病":
			fdiabetes="血糖正常"
		}
	}
	for rows2.Next(){
		var drink string
		err = rows2.Scan(&drink)
		if err != nil {
			panic(err)
		}
		switch drink{
		case "饮酒":
			fdrink="饮酒"
		case "不饮酒":
			fdrink="不饮酒"
		}
	}
	db.Close()
	if Iscore3==1{
		fbmi="超重"
		bmilab=1
	}else if Iscore3==2{
		fbmi="肥胖"
		bmilab=1
	}else{
		fbmi="体重正常"
		bmilab=0
	}
}

func insertfactor(){
	factor()
	db, err := sql.Open("mysql", "huang:zjubio2019@tcp(mysql.zjubiomedit.com:3306)/cdms_biomedit?charset=utf8")
	CheckErr(err)
	db.Exec("insert into eval_riskfactors(patientId,tc,smoke,Hdlc,sbp,dia,bmi,drink,measureTime) values (?,?,?,?,?,?,?,?,?)",patientId,ftbc,fsmoke,fHdlc,fsbp,fdiabetes,fbmi,fdrink,measureTime)
	db.Close()
}

func showfactor(w http.ResponseWriter){
	factor()
	var F factorlist
	F.Ftbc=ftbc
	F.Fsmoke=fsmoke
	F.FHdlc=fHdlc
	F.Fsbp=fsbp
	F.Fdiabetes=fdiabetes
	F.Fbmi=fbmi
	F.Fdrink=fdrink
	b, _ := json.Marshal(F)
	fmt.Fprintf(w,string(b))
}

func sexrecord(w http.ResponseWriter, r *http.Request){
	var server = "120.27.141.50"
	var port = "zjubiomedit"
	var user = "sa"
	var password = "BiomedIT@ZJU2015"
	var database = "HypertensionDB"

	//连接字符串
	connString := fmt.Sprintf("server=%s;port%s;database=%s;user id=%s;password=%s;encrypt=disable", server, port, database, user, password)
	//建立连接
	db, err := sql.Open("mssql", connString)

	r.ParseForm()
	if r.Method == "POST" {
		result, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()

		//未知类型的推荐处理方法
		var f interface{}
		json.Unmarshal(result, &f)
		m := f.(map[string]interface{})

		for k, v := range m {

			switch k {
			case "patientId":
				patientId=v.(string)
			}
		}
	}

	rows1,_:=db.Query("select SexCode from dbo.PersonPatient where PatientIdentifier = (?)",patientId)

	for rows1.Next(){
		var sexrecord string
		err = rows1.Scan(&sexrecord)
		if err != nil {
			panic(err)
		}
		switch sexrecord {
		case "M":
			sex="男"
		case "F":
			sex="女"
		}
	}
	db.Close()
	fmt.Fprintf(w,sex)
}

func agerecord(w http.ResponseWriter, r *http.Request){
	var server = "120.27.141.50"
	var port = "zjubiomedit"
	var user = "sa"
	var password = "BiomedIT@ZJU2015"
	var database = "HypertensionDB"

	//连接字符串
	connString := fmt.Sprintf("server=%s;port%s;database=%s;user id=%s;password=%s;encrypt=disable", server, port, database, user, password)
	//建立连接
	db, err := sql.Open("mssql", connString)
	r.ParseForm()
	if r.Method == "POST" {
		result, _ := ioutil.ReadAll(r.Body)
		r.Body.Close()

		//未知类型的推荐处理方法
		var f interface{}
		json.Unmarshal(result, &f)
		m := f.(map[string]interface{})

		for k, v := range m {

			switch k {
			case "patientId":
				patientId=v.(string)
			}
		}
	}
	rows1,_:=db.Query("select BirthDate from dbo.PersonPatient where PatientIdentifier = (?)",patientId)
	var mage int

	for rows1.Next(){
		var birthdate string
		var birthyear,birthmonth int
		err = rows1.Scan(&birthdate)
		if err!=nil{panic(err)}
		birthdate=strings.Split(birthdate,"T")[0]
		birthyear,_=strconv.Atoi(strings.Split(birthdate,"-")[0])
		birthmonth,_=strconv.Atoi(strings.Split(birthdate,"-")[1])
		mage=time.Now().Year()-birthyear
		if int(time.Now().Month())<birthmonth{
			mage--
		}
	}
	db.Close()

	fmt.Fprintf(w,strconv.Itoa(mage))
}