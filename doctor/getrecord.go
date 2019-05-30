package doctor

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type result struct {
	Wrisk string `json:"wrisk"'`
	Frisk string `json:"frisk"`
	Irisk string `json:"irisk"`
}

func Getrecord(w http.ResponseWriter,r *http.Request){
	var R result
	db, err := sql.Open("mysql", "huang:zjubio2019@tcp(mysql.zjubiomedit.com:3306)/cdms_biomedit?charset=utf8")
	CheckErr(err)
	r.ParseForm()
	if r.Method == "GET" {
		patientId= strings.Join(r.Form["patientId"], "")
	}
	rows1, _ := db.Query("select WHO from eval_results where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows1.Next(){
		var wrisk string
		err = rows1.Scan(&wrisk)
		CheckErr(err)
		R.Wrisk=wrisk
	}
	rows2, _ := db.Query("select Frammingham from eval_results where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows2.Next(){
		var frisk string
		err = rows2.Scan(&frisk)
		CheckErr(err)
		R.Frisk=frisk
	}
	rows3, _ := db.Query("select Frammingham from eval_results where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows3.Next(){
		var irisk string
		err = rows3.Scan(&irisk)
		CheckErr(err)
		R.Irisk=irisk
	}
	b, _ := json.Marshal(R)
	fmt.Fprintf(w,string(b))
	db.Close()
}
