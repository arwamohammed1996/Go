package main

import (
    "fmt"
    "html/template"
    "log"
    "net/http"
   "io/ioutil"
//    "net"
    "encoding/json"
//     "strconv"
)



    
    

//type cur struct {
//    Name string
//    value float64
//}    


func RootHandler(w http.ResponseWriter, r *http.Request) {


    //parse the html file 
     tmpl, err := template.ParseFiles("index.html")
      if err != nil {
        fmt.Println("Index Template Parse Error: ", err)
      }


   //get request to take the data from the api
      var v[6]float64
      var a [6]string
      a[0] = "KWD"
      a[1] = "BHD"
      a[2] = "OMR"
      a[3] = "JOD"
      a[4] ="GBP"
      a[5] ="EUR"
     


     for i:=0; i<6 ;i++{ 
     responce, err1 :=http.Get("https://free.currencyconverterapi.com/api/v6/convert?q="+a[i]+"_EGP&compact=ultra")
     //("https://free.currencyconverterapi.com/api/v6/convert?q=USD_EGP,EGP_USD")
     // 
      if err1 !=nil {
      fmt.Printf("The requset faild -_- %s",err1)
      }else{

       //Read the json responce
        data,_ :=ioutil.ReadAll(responce.Body)
         var dat map[string]float64

         if err := json.Unmarshal(data, &dat); err != nil {
         panic(err)
         }
         fmt.Println(dat)
         
       

                v[i]=(dat[a[i]+"_EGP"])
          // t:=strconv.FormatFloat(v, 'f', 6, 64)
           //fmt.Println(t)
          

        
    }

    currency := struct {
            Name string
            Value float64
        }{
            Name: (a[i]+"_EGP"),
            Value: v[i],
        }
         
         fmt.Println(v[i])
         err = tmpl.Execute(w,currency)
        }
   
       if err != nil {
        fmt.Println("Index Template Execution Error: ", err)
       }

}



func main() {
    http.HandleFunc("/", RootHandler) // sets router
    err := http.ListenAndServe(":4000", nil) // set listen port
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}

