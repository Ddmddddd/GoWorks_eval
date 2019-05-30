package doctor

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var patientId string

type risklist struct {
	Tbc string `json:"tbc"`
	Smoke string `json:"smoke"`
	Hdlc string `json:"Hdlc"`
	Bp string `json:"sbp"`
	Dia string `json:"dia"`
	Bmi string `json:"bmi"`
	Drink string `json:"drink"`
}

func Types(w http.ResponseWriter,r *http.Request){
	var R risklist
	db, err := sql.Open("mysql", "huang:zjubio2019@tcp(mysql.zjubiomedit.com:3306)/cdms_biomedit?charset=utf8")
	CheckErr(err)
	r.ParseForm()
	if r.Method == "GET" {
		patientId= strings.Join(r.Form["patientId"], "")
	}
	rows1, _ := db.Query("select tc from eval_riskfactors where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows1.Next(){
		var tbc string
		err = rows1.Scan(&tbc)
		CheckErr(err)
		R.Tbc=tbc
	}
	rows2, _ := db.Query("select smoke from eval_riskfactors where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows2.Next(){
		var smoke string
		err = rows2.Scan(&smoke)
		CheckErr(err)
		R.Smoke=smoke
	}
	rows3, _ := db.Query("select Hdlc from eval_riskfactors where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows3.Next(){
		var Hdlc string
		err = rows3.Scan(&Hdlc)
		CheckErr(err)
		R.Hdlc=Hdlc
	}
	rows4, _ := db.Query("select sbp from eval_riskfactors where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows4.Next(){
		var bp string
		err = rows4.Scan(&bp)
		CheckErr(err)
		R.Bp=bp
	}
	rows5, _ := db.Query("select dia from eval_riskfactors where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows5.Next(){
		var dia string
		err = rows5.Scan(&dia)
		CheckErr(err)
		R.Dia=dia
	}
	rows6, _ := db.Query("select bmi from eval_riskfactors where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows6.Next(){
		var bmi string
		err = rows6.Scan(&bmi)
		CheckErr(err)
		R.Bmi=bmi
	}
	rows7, _ := db.Query("select drink from eval_riskfactors where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows7.Next(){
		var drink string
		err = rows7.Scan(&drink)
		CheckErr(err)
		R.Drink=drink
	}
	b, _ := json.Marshal(R)
	fmt.Fprintf(w,string(b))
	db.Close()
}