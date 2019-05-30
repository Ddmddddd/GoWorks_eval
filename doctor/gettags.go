package doctor

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type tags struct {
	Smoke uint8 `json:"smoke"`
	Drink uint8 `json:"drink"`
	Bmi uint8 `json:"bmi"`
	Dia uint8 `json:"dia"`
	Tc uint8 `json:"tc"`
	Hdlc uint8 `json:"Hdlc"`
	Sbp uint8 `json:"sbp"`
}

func Gettags(w http.ResponseWriter,r *http.Request){
	var T tags
	r.ParseForm()
	if r.Method == "GET" {
		patientId= strings.Join(r.Form["patientId"], "")
	}
	db, err := sql.Open("mysql", "huang:zjubio2019@tcp(mysql.zjubiomedit.com:3306)/cdms_biomedit?charset=utf8")
	CheckErr(err)
	rows1, _ := db.Query("select tc from eval_types where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows1.Next() {
		var tbc []uint8
		err = rows1.Scan(&tbc)
		CheckErr(err)
		T.Tc = tbc[0]
	}
	rows2, _ := db.Query("select smoke from eval_types where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows2.Next(){
		var smoke []uint8
		err = rows2.Scan(&smoke)
		CheckErr(err)
		T.Smoke=smoke[0]
	}
	rows3, _ := db.Query("select Hdlc from eval_types where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows3.Next(){
		var Hdlc []uint8
		err = rows3.Scan(&Hdlc)
		CheckErr(err)
		T.Hdlc=Hdlc[0]
	}
	rows4, _ := db.Query("select sbp from eval_types where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows4.Next(){
		var bp []uint8
		err = rows4.Scan(&bp)
		CheckErr(err)
		T.Sbp=bp[0]
	}
	rows5, _ := db.Query("select dia from eval_types where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows5.Next(){
		var dia []uint8
		err = rows5.Scan(&dia)
		CheckErr(err)
		T.Dia=dia[0]
	}
	rows6, _ := db.Query("select bmi from eval_types where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows6.Next(){
		var bmi []uint8
		err = rows6.Scan(&bmi)
		CheckErr(err)
		T.Bmi=bmi[0]
	}
	rows7, _ := db.Query("select drink from eval_types where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows7.Next(){
		var drink []uint8
		err = rows7.Scan(&drink)
		CheckErr(err)
		T.Drink=drink[0]
	}
	b, _ := json.Marshal(T)
	fmt.Fprintf(w,string(b))
	db.Close()
}
