package main

import (
        "encoding/json"
        "os"  
        "io"      
        "log"
        "fmt"
        "io/ioutil"
        "net/http"
        "strings"
        "strconv"
)


// Config Load


type BMGConfig struct {
     
     // API Env
     ApiUrl string `json:"api_url"`
     ApiKey string `json:"api_key"`
     SelectedCurrency string `json:"currency"`
     StartDate string `json:"start_date"`
     EndDate string `json:"end_date"`
     PerPage int `json:"per_page"`   
     StartCount int `json:"start_count"`
     EndCount int `json:"end_count"`   
     
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
     Photos []Photos `json:"photos"`  
     PhotosUrl string  `json:"photosUrl"`

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

   
        var DealsFlatStruct FlatData
        var DealsDetil DataDetail


        json.Unmarshal(resp_body, &DealsFlatStruct)

        

        for i := range DealsFlatStruct.Data {

             item_flat := DealsFlatStruct.Data[i]

             
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
                 
                 fmt.Println(DealsDetil.Data.Uuid)

                 for j := range DealsDetil.Data.Photos {
                   item_photo := DealsDetil.Data.Photos[j]

                      // get image 

                       fmt.Println(item_photo.Paths.Original)


                       if item_photo.Paths.Size75x50 != "" {
                           getImageByUrl(DealsDetil.Data.PhotosUrl+item_photo.Paths.Size75x50, "75x50/"+getfile(item_photo.Paths.Size75x50))
                       }
                       if item_photo.Paths.Size175x112 != "" {                       
                           getImageByUrl(DealsDetil.Data.PhotosUrl+item_photo.Paths.Size175x112, "175x112/"+getfile(item_photo.Paths.Size175x112))
                       }
                       if item_photo.Paths.Size680x325 != "" {
                           getImageByUrl(DealsDetil.Data.PhotosUrl+item_photo.Paths.Size680x325, "680x325/"+getfile(item_photo.Paths.Size680x325))
                       }


                 }
   
        }      


}


func getfile(filename string) string {
     
     takestr := len(filename) 
     to := strings.LastIndex(filename, "/") + 1
     
     return filename[to:takestr] 
}



func main() {

    
    file2, e2 := ioutil.ReadFile("configs/bemyguest.json")
    if e2 != nil {
        fmt.Printf("File2 error: %v\n", e2)
        os.Exit(1)
    }
    
    json.Unmarshal(file2, &bmgObj)
    
    for i := bmgObj.StartCount; i < bmgObj.EndCount; i++ {

        fmt.Println("Loop = ", i)

        API_Pull(i)

    }



}