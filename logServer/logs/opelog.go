package logs

import (
	"fmt"
    "log"
	"io/ioutil"
	"encoding/json"
	"net/http"
	_ "github.com/lib/pq"

)

type opelog struct {
    Logid           string    `json:"logid"`
    Logtime         string    `json:"logtime"`
    Admin           string    `json:"admin"`
    Module      	string    `json:"moudle"`
    Logsql   		string    `json:"logsql"`
}

type classify struct {
    Category           string    `json:"category"`
    Content          string    `json:"content"`
}

var nowcontent string;
var nowcategory string;

func checkErr(err error) {   //报错
    if err != nil {
        log.Println("出错啦!")
        panic(err)
    }
}

func Category(w http.ResponseWriter, r *http.Request){          //获取分类内容
    w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
    w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
    w.Header().Set("content-type", "application/json")
    defer r.Body.Close()
    con, err := ioutil.ReadAll(r.Body) //获取post的数据
    checkErr(err)
    su := &classify{}         //把json转换回来
    json.Unmarshal([]byte(con), &su)
    fmt.Println("客户端输入：")
    fmt.Println("\tmodule:", su.Category)
    fmt.Println("\tcontent:", su.Content)
    nowcategory=su.Category;
    nowcontent=su.Content;
}

func Records(admin string, module string,logsql string){
	log.Println("正在添加日志...")
    rows, err := db.Prepare("insert into t_opelog (admin,module,logsql) values('$1','$2','$3')")
    checkErr(err)
    //_,err := stmt.Exec(admin, module, logsql)
    //checkErr(err)
    _,err = rows.Exec(admin, module, logsql)
    checkErr(err)
    log.Println("日志添加成功！")
}

func Display(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
    w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
    w.Header().Set("content-type", "application/json")
    var opelog []opelog
    opelog = Getopelog()
	data, _ := json.Marshal(opelog)
    fmt.Fprintf(w,string(data))
}

func Getopelog()[]opelog{
    var Opelog []opelog
    if nowcategory!="" {
        rows, err := db.Query("select logid,logtime,admin,module,logsql from opesql where $1=$2",nowcategory,nowcontent) 
        checkErr(err)
        for rows.Next(){
		var opelog opelog
		err = rows.Scan(&opelog.Logid,&opelog.Logtime,&opelog.Admin,&opelog.Module,&opelog.Logsql)
		if err != nil {
			fmt.Println("showscan error:",err)
		}
        Opelog=append(Opelog,opelog)
	}
    }else {
        rows, err := db.Query("select logid,logtime,admin,module,logsql from opesql") 
        checkErr(err)
        for rows.Next(){
		var opelog opelog
		err = rows.Scan(&opelog.Logid,&opelog.Logtime,&opelog.Admin,&opelog.Module,&opelog.Logsql)
		if err != nil {
			fmt.Println("showscan error:",err)
		}
        Opelog=append(Opelog,opelog)
	}
    }
    
    return Opelog
}
