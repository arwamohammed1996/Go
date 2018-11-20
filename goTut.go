package main 

import "fmt"
import "io/ioutil"
import "net/http"
import "encoding/json"
import "bytes"
//import "os"
//import "log"
//import"net/http"

func main() {
	//Request the responce from the api by HTTP.GET
 responce, err :=http.Get("https://free.currencyconverterapi.com/api/v6/convert?q=EUR_PHP&compact=y")
  if err !=nil {
  	fmt.Printf("The requset faild -_- %s",err)
  }else{
  	//Read the json responce
  	data,_ :=ioutil.ReadAll(responce.Body)
  	fmt.Println(string(data))
  }

  //post req ,create the data 
  jsonData :=map[string]string{"firstname":"Arwa","lastname":"Ali"}
  jsonValue,_:=json.Marshal(jsonData)
  //post req 
  responce,err =http.Post("https://httpbin.org/post","application/json",bytes.NewBuffer(jsonValue))
  if err !=nil {
  	fmt.Printf("The requset faild -_- %s",err)
  }else{
  	//Read the json responce
  	data,_ :=ioutil.ReadAll(responce.Body)
  	fmt.Println(string(data))
  }

}