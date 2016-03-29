package main

import (
        "database/sql" 
        "encoding/json"
        "os"  
        "io"      
        "log"
        "fmt"
        "io/ioutil"
        "net/http"
        _ "github.com/go-sql-driver/mysql"
        "strings"
        "strconv"
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

type DataDetail struct {

     Data DealsRec `json:"data"`

}

type DealsRec struct {

     Uuid string  `json:"uuid"`
     Published bool `json:"published"`
     Title string  `json:"title"`
     Description string  `json:"description"`
     Highlights string  `json:"highlights"`
     AdditionalInfo string  `json:"additionalInfo"`
     PriceIncludes string  `json:"priceIncludes"`
     PriceExcludes string  `json:"priceExcludes"`
     Itinerary string  `json:"itinerary"`
     Warnings string  `json:"warnings"`
     Safety string  `json:"safety"`
     Lat string  `json:"latitude"`
     Lon string  `json:"longitude"`
     MinPax int  `json:"minPax"`
     MaxPax int  `json:"maxPax"`
     BasePrice  float32  `json:"basePrice"`
     IsFlatPaxPrice bool  `json:"isFlatPaxPrice"`
     ReviewCount int  `json:"reviewCount"`
     ReviewAverageScore int  `json:"reviewAverageScore"`
     TypeName string  `json:"typeName"`
     TypeUuid string  `json:"typeUuid"`
     PhotosUrl string  `json:"photosUrl"`
     BusinessHoursFrom string `json:"businessHoursFrom"`
     BusinessHoursTo string  `json:"businessHoursTo"`
     MeetingTime string  `json:"meetingTime"`
     HotelPickup string  `json:"hotelPickup"`
     MeetingLocation string  `json:"meetingLocation"`
     Url string  `json:"url"`
     StaticUrl string  `json:"staticUrl"`
     Currency Currencies `json:"currency"` 
     Photos []Photos `json:"photos"`
     Categories []Categories `json:"categories"`
     Locations []Locations `json:"locations"`
     ProductTypes []ProductTypes `json:"productTypes"`     

}

type Currencies struct {

      Code string  `json:"code"`
      Symbol string  `json:"symbol"`
      Uuid string  `json:"uuid"`

}

type Photosarr struct {
      
      Photos []Photos  `json:"photos"` 

}

type Photos struct {
      
      Caption string  `json:"caption"` 
      Uuid string  `json:"uuid"`
      Paths PhotoPathInfo `json:"paths"`

}

type PhotoPathInfo struct {

     Original string  `json:"original"`
     Size75x50 string  `json:"75x50"`
     Size175x112 string  `json:"175x112"`
     Size680x325 string  `json:"680x325"`
     Size1280x720 string  `json:"1280x720"`

}

type Categories struct {

      Name string  `json:"name"` 
      Uuid string  `json:"uuid"` 

}

type Locations struct {

      City string  `json:"city"`
      CityUuid string  `json:"cityUuid"`
      State string  `json:"state"`
      StateUuid string  `json:"stateUuid"`
      Country string  `json:"country"`
      CountryUuid string  `json:"countryUuid"` 

}


type ProductTypes struct {

      Uuid string  `json:"uuid"`
      Title string  `json:"title"`
      Description int `json:"description"`
      DurationDays int `json:"durationDays"`
      DurationHours int `json:"durationHours"`
      DurationMinutes int `json:"durationMinutes"` 
      MinPax int  `json:"paxMin"`
      MaxPax int  `json:"paxMax"`
      DaysInAdvance int `json:"daysInAdvance"`
      IsNonRefundable bool `json:"isNonRefundable"`
      HasChildPrice bool `json:"hasChildPrice"`
      MinAdultAge int `json:"minAdultAge"`
      MaxAdultAge int `json:"maxAdultAge"`
      AllowChildren bool `json:"allowChildren"`
      MinChildAge int `json:"minChildAge"`
      MaxChildAge int `json:"maxChildAge"`
      InstantConfirmation bool `json:"instantConfirmation"`
      VoucherUse string `json:"voucherUse"`
      VoucherRedemptionAddress string `json:"voucherRedemptionAddress"`
      Prices map[string]string `json:"prices"`

}

var mysqlObj MysqlConfig
var bmgObj BMGConfig

// End Flat Data Struct



func getImageByUrl(url string, newfilename string) {
    
    // don't worry about errors
    response, e := http.Get(url)
    if e != nil {
        log.Fatal(e)
    }

    defer response.Body.Close()

    //open a file for writing
    file, err := os.Create("assets/images/deals/" + newfilename)
    if err != nil {
        log.Fatal(err)
    }
    // Use io.Copy to just dump the response body to the file. This supports huge files
    _, err = io.Copy(file, response.Body)
    if err != nil {
        log.Fatal(err)
    }
    file.Close()
}


// BemyGuest API Pull 
func API_Pull(PageNum int) {

        
        client := &http.Client{}

        ApiUrl := bmgObj.ApiUrl + "products?page=" + strconv.Itoa(PageNum) + "&per_page=" + strconv.Itoa(bmgObj.PerPage) + "&published=published&date_start=" + bmgObj.StartDate + "&date_end=" + bmgObj.EndDate + "&currency=" + bmgObj.SelectedCurrency

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

      
    stmtFlat, _ := db.Prepare("UPDATE deals_flat_data SET published=?, title=?, descriptions=?, highlights=?, additional_info=?, price_includes=?, price_excludes=?, itenerary=?, warnings=?, safety=?, latitude=?, longitude=?, min_pax=?, max_pax=?, base_price=?, review_count=?, review_average_score=?, meeting_location=?, country_uuid=?, state_uuid=?, city_uuid=?, type_uuid=?, currency_code=?, currency_symbol=?, type_name=?, country_name=?, state_name=?, city_name=?, url=?, static_url=?, business_hours_from=?, business_hours_to=?, meeting_time=? where deal_uuid = ?")

    stmt, _ := db.Prepare("UPDATE deals SET type_uuid = ?, country_uuid = ?, state_uuid = ?, city_uuid = ?, title = ?, description = ?, highlights = ?, additional_info = ?, price_includes = ?, price_excludes = ?, itinerary = ?, warnings = ?, safety = ?, meeting_location = ?, min_pax = ?, max_pax = ?,  business_hours_from = ?, business_hours_to = ?, meeting_time = ?, latitude = ?, longitude = ?, photosUrl = ?, url = ?, staticUrl = ?, currency_uuid = ?, base_price = ? where deal_uuid = ?")     
    
    stmt2, _ := db.Prepare("UPDATE deals_product_type SET title = ?, description = ?, duration_days = ?, duration_hours = ?, duration_minutes = ?, pax_min = ?, pax_max = ?, days_in_advance = ?, is_non_refundable = ?, has_child_price = ?, min_adult_age = ?, max_adult_age = ?, allow_children = ?, min_child_age = ?, max_child_age = ?, instant_confirmation = ?, voucher_use = ?, voucher_redemption_address = ? WHERE product_type_uuid = ?")
 

        var DealsFlatStruct FlatData
        var DealsDetil DataDetail

        //var data map[string]map[string]json.RawMessage 

        json.Unmarshal(resp_body, &DealsFlatStruct)
        //json.Unmarshal(resp_body["photos"], &Photosarr)
        //json.Unmarshal(data, DealsFlatStruct)

        //fmt.Println(string(resp_body))
        

        for i := range DealsFlatStruct.Data {

             item_flat := DealsFlatStruct.Data[i]

             
             /*for k := range item_flat.Locations {
                   item_location := item_flat.Locations[k]

                   fmt.Println(item_location.City)
             }*/

             for l := range item_flat.Categories {
                   item_category := item_flat.Categories[l]

                   stmt7.Exec(item_flat.Uuid, item_category.Uuid)
             }
             
             // get Deal Detil

             ApiUrl2 := bmgObj.ApiUrl + "products/" + item_flat.Uuid + "/?currency=" + bmgObj.SelectedCurrency + "&date_start=" + bmgObj.StartDate + "&date_end=" + bmgObj.EndDate

                 req2, _ := http.NewRequest("GET", ApiUrl2, nil)
                 req2.Header.Add("X-Authorization", bmgObj.ApiKey)

                 resp2, err2 := client.Do(req2)

                 if err2 != nil {
                       fmt.Println(err2)
                       fmt.Println("Errored in deal detail")
                       return
                 }

                 defer resp2.Body.Close()

                 resp_body2, _ := ioutil.ReadAll(resp2.Body)
                 

                 json.Unmarshal(resp_body2, &DealsDetil)
                 
                 //stmtFlat, _ := db.Prepare("INSERT INTO deals_flat_data(deal_uuid, published, title, descriptions, highlights, additional_info, price_includes, price_excludes, itenerary, warnings, safety, latitude, longitude, min_pax, max_pax, base_price, selling_price, review_count, review_average_score, meeting_location, country_uuid, state_uuid, city_uuid, type_uuid, currency_code, currency_symbol, type_name, country_name, state_name, city_name, url, static_url, business_hours_from, business_hours_to, meeting_time, hotel_pickup) VALUES(?, ?, ?, ?, ?,?, ?, ?, ?, ?,?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)") 

                 fmt.Println(DealsDetil.Data.Uuid)

                 stmtFlat.Exec("0", DealsDetil.Data.Title, DealsDetil.Data.Description, DealsDetil.Data.Highlights, DealsDetil.Data.AdditionalInfo, DealsDetil.Data.PriceIncludes, DealsDetil.Data.PriceExcludes, DealsDetil.Data.Itinerary, DealsDetil.Data.Warnings, DealsDetil.Data.Safety, DealsDetil.Data.Lat, DealsDetil.Data.Lon, DealsDetil.Data.MinPax, DealsDetil.Data.MaxPax, DealsDetil.Data.MinPax, DealsDetil.Data.BasePrice, DealsDetil.Data.ReviewCount, DealsDetil.Data.ReviewAverageScore, DealsDetil.Data.MeetingLocation, DealsDetil.Data.Locations[0].CountryUuid, DealsDetil.Data.Locations[0].StateUuid, DealsDetil.Data.Locations[0].CityUuid, DealsDetil.Data.TypeUuid, DealsDetil.Data.Currency.Code, DealsDetil.Data.Currency.Symbol, DealsDetil.Data.TypeName, DealsDetil.Data.Locations[0].Country, DealsDetil.Data.Locations[0].State, DealsDetil.Data.Locations[0].City, DealsDetil.Data.Url, DealsDetil.Data.StaticUrl, DealsDetil.Data.BusinessHoursFrom, DealsDetil.Data.BusinessHoursTo, DealsDetil.Data.MeetingTime, DealsDetil.Data.HotelPickup, DealsDetil.Data.Uuid )

   

                 stmt.Exec(DealsDetil.Data.TypeUuid, DealsDetil.Data.Locations[0].CountryUuid, DealsDetil.Data.Locations[0].StateUuid, DealsDetil.Data.Locations[0].CityUuid, DealsDetil.Data.Title, DealsDetil.Data.Description, DealsDetil.Data.Highlights, DealsDetil.Data.AdditionalInfo, DealsDetil.Data.PriceIncludes, DealsDetil.Data.PriceExcludes, DealsDetil.Data.Itinerary, DealsDetil.Data.Warnings, DealsDetil.Data.Safety, DealsDetil.Data.MeetingLocation,  DealsDetil.Data.MinPax, DealsDetil.Data.MaxPax, DealsDetil.Data.BusinessHoursFrom, DealsDetil.Data.BusinessHoursTo, DealsDetil.Data.MeetingTime, DealsDetil.Data.Lat, DealsDetil.Data.Lon, DealsDetil.Data.PhotosUrl, DealsDetil.Data.Url, DealsDetil.Data.StaticUrl, DealsDetil.Data.Currency.Uuid, DealsDetil.Data.BasePrice, DealsDetil.Data.Uuid )

                   
                 for m:= range DealsDetil.Data.ProductTypes {

                      item_producttype := DealsDetil.Data.ProductTypes[m]

                      stmt2.Exec(item_producttype.Title, item_producttype.Description, item_producttype.DurationDays, item_producttype.DurationHours, item_producttype.DurationMinutes, item_producttype.MinPax, item_producttype.MaxPax, item_producttype.DaysInAdvance, item_producttype.IsNonRefundable, item_producttype.HasChildPrice, item_producttype.MinAdultAge, item_producttype.MaxAdultAge, item_producttype.AllowChildren, item_producttype.MinChildAge, item_producttype.MaxChildAge, item_producttype.InstantConfirmation, item_producttype.VoucherUse, item_producttype.VoucherRedemptionAddress, item_producttype.Uuid )

                 }


             //currency_insert.Exec(item_flat.Uuid)

        }      
        
       defer stmtFlat.Close()
       defer stmt2.Close() 
       defer stmt.Close() 

    defer db.Close()


}


func getfile(filename string) string {
     
     takestr := len(filename) 
     to := strings.LastIndex(filename, "/") + 1
     
     return filename[to:takestr] 
}



func main() {
   //a := getfile("/images/content/original/33082i1ae1e21d3cd38ef3f16ca97eb9ef6b12.jpg")

   //fmt.Println(a)

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
    
    for i := 0; i < 100; i++ {

        fmt.Println("Loop = ", i)

        API_Pull(i)

    }    




}