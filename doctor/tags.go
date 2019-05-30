package doctor

import (
	"database/sql"
	"net/http"
	"strings"
)

func Tags(w http.ResponseWriter,r *http.Request){
	var patientId string
	drink:=1
	smoke:=1
	bmi:=1
	dia:=1
	sbp:=1
	tc:=1
	Hdlc:=1
	r.ParseForm()
	if r.Method == "GET" {
		patientId= strings.Join(r.Form["patientId"], "")
	}
	db, err := sql.Open("mysql", "huang:zjubio2019@tcp(mysql.zjubiomedit.com:3306)/cdms_biomedit?charset=utf8")
	CheckErr(err)
	db.Exec("delete from eval_types where patienId=(?)",patientId)
	db.Exec("insert into eval_types(patientId,drink,smoke,bmi,dia,sbp,tc,Hdlc) values (?,?,?,?,?,?,?,?)",patientId,drink,smoke,bmi,dia,sbp,tc,Hdlc)
	db.Close()
}