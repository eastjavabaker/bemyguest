package main

import (
        "database/sql" 
        "encoding/json"
        "os"  
        "fmt"
        "io/ioutil"
        "net/http"
        _ "github.com/go-sql-driver/mysql"
        "bytes"
)


// Config Load
type MysqlConfig struct {

     MySqlHost string `json:"host"`
     MySqlUsername string `json:"username"`
     MySqlPassword string `json:"password"`
     MysqlPort string `json:"port"`
     MysqlDb string `json:"database"`
    
}

type BMGConfig struct {
     
     // API Env
     ApiUrl string `json:"api_url"`
     ApiKey string `json:"api_key"`
     SelectedCurrency string `json:"currency"`
     StartDate string `json:"start_date"`
     EndDate string `json:"end_date"`
     PerPage int `json:"per_page"`     
     
}


// Flat Data Struct 
type FlatData struct {
     Data []DealsRec `json:"data"`
}


type DealsRec struct {
     Uuid string  `json:"uuid"`
     Title string `json:"title"`
     DeletedDate string  `json:"deletedAt"`    
}


var mysqlObj MysqlConfig
var bmgObj BMGConfig

// End Flat Data Struct

// BemyGuest API Pull 
func API_UPDATE() {

        
        client := &http.Client{}

        ApiUrl := bmgObj.ApiUrl + "products?deleted=true" 
        
        fmt.Println(ApiUrl)

        req, _ := http.NewRequest("GET", ApiUrl, nil)
        req.Header.Add("X-Authorization", bmgObj.ApiKey)

        resp, err := client.Do(req)

        if err != nil {
                fmt.Println(err)
                fmt.Println("Errored when sending request to the server")
                return
        }

        defer resp.Body.Close()

        resp_body, _ := ioutil.ReadAll(resp.Body)

        // DB Insert, temporary method    


    db, errdb := sql.Open("mysql", mysqlObj.MySqlUsername+":"+mysqlObj.MySqlPassword+"@tcp("+mysqlObj.MySqlHost+":"+mysqlObj.MysqlPort+")/"+mysqlObj.MysqlDb)
    if errdb != nil {
        panic(errdb.Error())  
    }



    stmt7, _ := db.Prepare("INSERT INTO deals_deleted(deal_uuid, title, deleted_date) values(?,?,?)")  
      
   
    var buffer bytes.Buffer
    var DealsFlatStruct FlatData

        json.Unmarshal(resp_body, &DealsFlatStruct)
        
    buffer.WriteString("(")
        for i := range DealsFlatStruct.Data {

             item_flat := DealsFlatStruct.Data[i]

             fmt.Println(item_flat.Uuid)

             stmt7.Exec(item_flat.Uuid, item_flat.Title, item_flat.DeletedDate)   
             
             buffer.WriteString("'" + item_flat.Uuid + "',")

             
        }      
       
       buffer.WriteString("'0')")
       
       fmt.Println(buffer.String())
       
     stmtDelete, _ := db.Prepare("UPDATE deals_flat_data SET published = 2 WHERE deal_uuid IN " + buffer.String())  
     
     stmtDelete.Exec()
       
       defer stmt7.Close() 
       defer stmtDelete.Close() 

    defer db.Close()


}



func main() {

    file, e := ioutil.ReadFile("configs/mysql.json")
    if e != nil {
        fmt.Printf("File error: %v\n", e)
        os.Exit(1)
    }
    
    file2, e2 := ioutil.ReadFile("configs/bemyguest.json")
    if e2 != nil {
        fmt.Printf("File2 error: %v\n", e2)
        os.Exit(1)
    }
    
    json.Unmarshal(file, &mysqlObj)
    json.Unmarshal(file2, &bmgObj)

    API_UPDATE()

}