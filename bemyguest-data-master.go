package main

import (
        "database/sql" 
        "encoding/json"
        "os"
        "fmt"
        "io/ioutil"
        "net/http"
        _ "github.com/go-sql-driver/mysql"
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



// Locations Data Struct

type LocationsData struct {

     Data []ContinentssRec `json:"data"`

}

type ContinentssRec struct { 

     Continent string `json:"continent"`
     Uuid  string `json:"uuid"`
     Code  string `json:"code"` 
     Countries Countries `json:"countries"`

}


type Countries struct {

     Data []CountriesRec `json:"data"`

}

type CountriesRec struct {

     Name string `json:"name"` 
     Code string `json:"code"` 
     Uuid string `json:"uuid"` 
     States States `json:"states"` 
}

type States struct {

     Data []StatesRec `json:"data"`

}

type StatesRec struct {

     Name string `json:"name"` 
     Uuid string `json:"uuid"` 
     Cities Cities `json:"cities"` 
}

type Cities struct {

     Data []CitiesRec `json:"data"`

}

type CitiesRec struct {

     Name string `json:"name"` 
     Uuid string `json:"uuid"` 
}

// End Locations



// Types Data Struct
type TypesData struct {

     Data []TypesRec `json:"data"`

}

type TypesRec struct {

     Name string `json:"name"`
     Uuid  string `json:"uuid"`

}
// End types 


// Currency Data Struct
type CurrenciesData struct {

     Data []CurrenciesRec `json:"data"`

}

type CurrenciesRec struct {

     Name string `json:"name"`
     Uuid  string `json:"uuid"`
     Code  string `json:"code"`
     Symbol  string `json:"symbol"`

}
// End currencies 

// Category Data Struct 
type CategoryData struct {

         Data []ParentCategory `json:"data"`

}

type ChildCategory struct {
    Name string `json:"name"`
    Uuid  string `json:"uuid"`
}

type ParentCategory struct {
    Name string `json:"name"`
    Uuid  string `json:"uuid"`
    Children []ChildCategory `json:"children"`
}
// end category Data Struct

var mysqlObj MysqlConfig
var bmgObj BMGConfig

// BemyGuest API Pull 
func API_Data_master_Pull() {

     /*type Message struct {
          name, uuid string
     }*/
     
    
    ApiUrl := bmgObj.ApiUrl + "config" 
    
     client := &http.Client{}

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

    
    category_insert, _ := db.Prepare("INSERT INTO deals_categories(category_uuid, root_uuid, category_name) VALUES(?,?,?) ")
    currency_insert, _ := db.Prepare("INSERT INTO deals_currencies(currency_uuid, code, symbol, name) VALUES(?,?,?,?) ")
    types_insert, _ := db.Prepare("INSERT INTO deals_type(type_uuid, type_name) values(?,?)")
    continent_insert, _ := db.Prepare("INSERT INTO deals_location_continents(continent_uuid, continent_code,continent_name) values(?,?,?)")
    countries_insert, _ := db.Prepare("INSERT INTO deals_location_countries(country_uuid, continent_uuid, country_code, country_name) values(?,?,?,?)")
    states_insert, _ := db.Prepare("INSERT INTO deals_location_states(state_uuid, country_uuid, state_name) values(?,?,?)")
    cities_insert, _ := db.Prepare("INSERT INTO deals_location_cities(city_uuid, state_uuid, city_name) values(?,?,?)")

        var data map[string]map[string]json.RawMessage       
        var categoriesStruct CategoryData
        var currenciesStruct CurrenciesData
        var typesStruct TypesData
        var locationsStruct LocationsData


        json.Unmarshal(resp_body, &data)

        json.Unmarshal(data["data"]["categories"], &categoriesStruct) // get categories structure
        json.Unmarshal(data["data"]["currencies"], &currenciesStruct) // get currencies structure
        json.Unmarshal(data["data"]["types"], &typesStruct) // get types structure
        json.Unmarshal(data["data"]["locations"], &locationsStruct) // get types structure
        

        for i := range currenciesStruct.Data {

             itemy := currenciesStruct.Data[i]

             //fmt.Println(itemy.Name)

             currency_insert.Exec(itemy.Uuid, itemy.Code, itemy.Symbol, itemy.Name)

        }


        for j := range typesStruct.Data {

             itemz := typesStruct.Data[j]

             fmt.Println(itemz.Name)

             types_insert.Exec(itemz.Uuid, itemz.Name)

        }

        
        for g := range categoriesStruct.Data {
               itemx := categoriesStruct.Data[g]
               fmt.Println(string(itemx.Uuid), "0", string(itemx.Name))
               category_insert.Exec(string(itemx.Uuid), "0", string(itemx.Name))

                for h := range itemx.Children {
                     item_children := itemx.Children[h]
                     //fmt.Println(item_children.Name)
                     category_insert.Exec(string(item_children.Uuid), string(itemx.Uuid), string(item_children.Name))
                }
                
        }

        for k := range locationsStruct.Data {
               item_loc := locationsStruct.Data[k]
               fmt.Println(item_loc.Continent)
               continent_insert.Exec(item_loc.Uuid, item_loc.Code, item_loc.Continent)
               for l := range item_loc.Countries.Data {
                     item_country := item_loc.Countries.Data[l]
                     fmt.Println(item_country.Name)
                     countries_insert.Exec(item_country.Uuid, item_loc.Uuid, item_country.Code, item_country.Name)

                        for m:= range item_country.States.Data {
                             item_states := item_country.States.Data[m]
                             fmt.Println(item_states.Name)
                             states_insert.Exec(item_states.Uuid, item_country.Uuid, item_states.Name)
                                 for n:= range item_states.Cities.Data {
                                        item_cities := item_states.Cities.Data[n]
                                        fmt.Println(item_cities.Name)
                                        cities_insert.Exec(item_cities.Uuid, item_states.Uuid, item_cities.Name)
                                 }
                        }
                    fmt.Println("##################")
               }
              fmt.Println("---------")  
        }

        defer category_insert.Close()
        defer currency_insert.Close()
        defer types_insert.Close()
        defer continent_insert.Close()
        defer countries_insert.Close()
        defer states_insert.Close()
        defer cities_insert.Close()


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

    API_Data_master_Pull()

}
