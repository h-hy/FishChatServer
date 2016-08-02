package mysql_store
 
import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "fmt"
    "strings"
    // "net/http"
    // "github.com/oikomi/FishChatServer/provider"
    // "github.com/oikomi/FishChatServer/log"
)
type MysqlStore struct {
    db *sql.DB
}

func NewMysqlStore(ip string, port string, user string, password string,database string,maxOpenConn int,maxOIdleConn int) *MysqlStore {
    var db *sql.DB
    db, _ = sql.Open("mysql", ""+user+":"+password+"@tcp("+ip+":"+port+")/"+database+"?charset=utf8")
    db.SetMaxOpenConns(maxOpenConn)
    db.SetMaxIdleConns(maxOIdleConn)
    db.Ping()
    return &MysqlStore{
        db: db,
    }
}

func (self *MysqlStore) DeviceStore(IMEI,query string, args ...interface{}) error {
    var argSql string
    for i := 0; i < len(args); i++ {
        argSql+=",?"
    }
    if (len(args)>0){
        query=", "+query
    }
    stmt, err := self.db.Prepare("INSERT into `device` ( IMEI "+query+" ) values (? "+argSql+")")
    if err != nil {
        return err
    }
    result := make([]interface{}, len(args)+1) 
    result[0]=IMEI
    copy(result[1:], args) 

    _, err = stmt.Exec(result...)
    if err != nil {
        return err
    }
    return nil
}


func (self *MysqlStore) DeviceUpdate(IMEI,query string, args ...interface{}) error {
    stmt, err := self.db.Prepare("UPDATE `device` SET "+query+" WHERE IMEI = ?")
    if err != nil {
        return err
    }
    args = append(args,IMEI)
    _, err = stmt.Exec(args...)
    if err != nil {
        return err
    }
    return nil
}

func (self *MysqlStore) GetDeviceFromIMEI(IMEI string) (map[string]string,error) {
    columns := []string{"id", "name"}
    sql := fmt.Sprintf("select IMEI from `device` where IMEI = ?", strings.Join(columns, ", "))
    row := self.db.QueryRow(sql, IMEI)
    scanArgs := make([]interface{}, len(columns))
    values := make([]interface{}, len(columns))
    for j := range values {
        scanArgs[j] = &values[j]
    }
 
    record := make(map[string]string)
        //将行数据保存到record字典
    err := row.Scan(scanArgs...)
    if (err != nil){
        return nil,err
    }
    for i, col := range values {
        if col != nil {
            record[columns[i]] = string(col.([]byte))
        }
    }
    fmt.Println(record)
    return record,nil
}


 

// func main() {
//     var device NewMysqlStore
//     device.IMEI = "12333333"
//     err := device.Store("ICCID","123")
//     if err != nil {
//         log.Error(err.Error())
//     }
//     // device["ICCID"] = "123"

//     device.Update("ICCID= ? ","123")
//     // startHttpServer()
// }
 
// func startHttpServer() {
//     http.HandleFunc("/pool", do)
//     err := http.ListenAndServe(":9090", nil)
//     if err != nil {
//         log.Fatal("ListenAndServe: ", err)
//     }
// }
// func do(w http.ResponseWriter, r *http.Request){
//     go pool(w,r)
// }
// func pool(w http.ResponseWriter, r *http.Request) {

//     rows, err := db.Query("SELECT * FROM watch limit 1")
//     defer rows.Close()
//     rows, err = db.Query("SELECT * FROM watch limit 1")
//     defer rows.Close()
//     checkErr(err)
//     fmt.Println("ok")
 
//     columns, _ := rows.Columns()
//     scanArgs := make([]interface{}, len(columns))
//     values := make([]interface{}, len(columns))
//     for j := range values {
//         scanArgs[j] = &values[j]
//     }
 
//     record := make(map[string]string)
//     for rows.Next() {
//         //将行数据保存到record字典
//         err = rows.Scan(scanArgs...)
//         for i, col := range values {
//             if col != nil {
//                 record[columns[i]] = string(col.([]byte))
//             }
//         }
//     }
 
//     fmt.Println(record)
//     fmt.Fprintln(w, "finish")
// }
 
// func checkErr(err error) {
//     if err != nil {
//         fmt.Println(err)
//         panic(err)
//     }
// }