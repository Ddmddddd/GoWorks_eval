package doctor

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	 "time"
)

func Lab(w http.ResponseWriter,r *http.Request){
	db, err := sql.Open("mysql", "huang:zjubio2019@tcp(mysql.zjubiomedit.com:3306)/cdms_biomedit?charset=utf8&parseTime=true")
	CheckErr(err)
	r.ParseForm()
	if r.Method == "GET" {
		patientId= strings.Join(r.Form["patientId"], "")
	}
	var evaluatelab int
	var d int64 =20190527101010
	t:=time.Unix(d,0)
	now:=time.Now()
	rows1, _ := db.Query("select measureTime from eval_userdata where patientId=(?) ORDER BY id DESC LIMIT 1",patientId)
	for rows1.Next(){

		err = rows1.Scan(&t)
		CheckErr(err)
	}
	subM:=now.Sub(t)
	time_hour:=int(subM.Hours())+8
	if( time_hour >128   ){
		evaluatelab=0
	}else {
		evaluatelab=1
	}
	fmt.Fprintf(w,strconv.Itoa(evaluatelab))
}
