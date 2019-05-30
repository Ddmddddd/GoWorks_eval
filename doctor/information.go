package doctor

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type information struct {
	Drink string `json:"drink"`
	Tc string `json:"tc"`
	Smoke string `json:"smoke"`
	Hdlc string `json:"Hdlc"`
	Sbp string `json:"sbp"`
	Dia string `json:"dia"`
	Bmi string `json:"bmi"`
	Height string `json:"height"`
	Weight string `json:"weight"`
}

func Information(w http.ResponseWriter,r *http.Request){
	var I information
	db, err := sql.Open("mysql", "huang:zjubio2019@tcp(mysql.zjubiomedit.com:3306)/cdms_biomedit?charset=utf8")
	CheckErr(err)
	r.ParseForm()
	if r.Method == "GET" {
		patientId= strings.Join(r.Form["patientId"], "")
	}
	rows1, _ := db.Query("select drink from eval_userdata where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows1.Next(){
		var drink string
		err = rows1.Scan(&drink)
		CheckErr(err)
		I.Drink=drink
	}
	rows2, _ := db.Query("select smoke from eval_userdata where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows2.Next(){
		var smoke string
		err = rows2.Scan(&smoke)
		CheckErr(err)
		I.Smoke=smoke
	}
	rows3, _ := db.Query("select tc from eval_userdata where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows3.Next(){
		var tc string
		err = rows3.Scan(&tc)
		CheckErr(err)
		I.Tc=tc
	}
	rows4, _ := db.Query("select Hdlc from eval_userdata where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows4.Next(){
		var Hdlc string
		err = rows4.Scan(&Hdlc)
		CheckErr(err)
		I.Hdlc=Hdlc
	}
	rows5, _ := db.Query("select height from eval_userdata where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	var height float64
	for rows5.Next(){
		err = rows5.Scan(&height)
		CheckErr(err)
		I.Height=FloatToString(height)
	}
	rows6, _ := db.Query("select weight from eval_userdata where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	var weight float64
	for rows6.Next(){
		err = rows6.Scan(&weight)
		CheckErr(err)
		I.Weight=FloatToString(weight)
	}
	var bmi float64
		bmi=weight/height/height*10000
		I.Bmi=FloatToString(bmi)
		rows7, _ := db.Query("select sbp from eval_userdata where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
		for rows7.Next(){
			var sbp string
			err = rows7.Scan(&sbp)
			CheckErr(err)
			I.Sbp=sbp
		}
		rows8, _ := db.Query("select dia from eval_userdata where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
		for rows8.Next(){
			var dia string
		err = rows8.Scan(&dia)
		CheckErr(err)
		I.Dia=dia
	}
	b, _ := json.Marshal(I)
	fmt.Fprintf(w,string(b))
	db.Close()
}

func FloatToString(input_num float64) string {
	return strconv.FormatFloat(float64(input_num), 'f', 1, 64)
}
