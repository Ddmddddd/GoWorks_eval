package doctor

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"
)

func Deletetag(w http.ResponseWriter,r *http.Request){
	var ID int
	r.ParseForm()
	if r.Method == "GET" {
		patientId= strings.Join(r.Form["patientId"], "")
		id := strings.Join(r.Form["id"],"")
		ID,_=strconv.Atoi(id)
	}
	db, err := sql.Open("mysql", "huang:zjubio2019@tcp(mysql.zjubiomedit.com:3306)/cdms_biomedit?charset=utf8")
	CheckErr(err)
	if(ID==1){
		db.Exec("update eval_types set drink=0 where patientId = (?)",patientId)
	}
	if(ID==2){
		db.Exec("update eval_types set smoke=0 where patientId = (?)",patientId)
	}
	if(ID==3){
		db.Exec("update eval_types set tc=0 where patientId = (?)",patientId)
	}
	if(ID==4){
		db.Exec("update eval_types set sbp=0 where patientId = (?)",patientId)
	}
	if(ID==5){
		db.Exec("update eval_types set Hdlc=0 where patientId = (?)",patientId)
	}
	if(ID==6){
		db.Exec("update eval_types set dia=0 where patientId = (?)",patientId)
	}
	if(ID==7){
		db.Exec("update eval_types set bmi=0 where patientId = (?)",patientId)
	}
	db.Close()
}
