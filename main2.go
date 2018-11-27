package main

import (
        "gopkg.in/mgo.v2"
        "gopkg.in/mgo.v2/bson"
        "net/http"
        "github.com/gorilla/context"
        "fmt"
        "encoding/json"
        "time"
        "log"
        "net"
        "io/ioutil"
)

func main() {

	 // connect to the database
  db, err := mgo.Dial("localhost")
  if err != nil {
    log.Fatal("cannot dial mongo", err)
  }
  defer db.Close() // clean up when we’re done
  // Adapt our handle function using withDB
  h := Adapt(http.HandlerFunc(handle), withDB(db))
  // add the handler
  http.Handle("/comments", context.ClearHandler(h))
  // start the server
  if err := http.ListenAndServe(":8080", nil); err != nil {
    log.Fatal(err)
  }

}


//We’re going to use it to create a session that connects to MongoDB 
//before the request handler code runs, and then clean up 
//that session once the handler code has finished.
type Adapter func(http.Handler) http.Handler
func Adapt(h http.Handler, adapters ...Adapter) http.Handler {
  for _, adapter := range adapters {
    h = adapter(h)
  }
  return h
}


//e’re going to write a function that returns an Adapter that will setup 
//the database session for our handlers and store it in a context 
func withDB(db *mgo.Session) Adapter {
  // return the Adapter
  return func(h http.Handler) http.Handler {
    // the adapter (when called) should return a new handler
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      // copy the database session      
      dbsession := db.Copy()
      defer dbsession.Close() // clean up 
      // save it in the mux context
      context.Set(r, "database", dbsession)
      // pass execution to the original handler
      h.ServeHTTP(w, r)
    })
  }
}

//handle the get post requests
func handle(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
  case "GET":
    handleRead(w, r)
  case "POST":
    //I can makr for loop as many cyrrencies I want and send it to the function 
  
//  for i:=0; i<6 ;i++{
    handleInsert(w, r)
 // }
  default:
    http.Error(w, "Not supported", http.StatusMethodNotAllowed)
  }
}


//Our `comment` struct will hold the data that represents a single comment. 
type comment struct {
  ID     bson.ObjectId `json:"id,string" bson:"_id"`
  Cur string        `json:"cur" bson:"cur"`
  Rate   string        `json:"rate" bson:"rate"`
  When   time.Time     `json:"when" bson:"when"`
}


// will take the data from an http.Request, and insert it into the database. 
func handleInsert(w http.ResponseWriter, r *http.Request) {
  db := context.Get(r, "database").(*mgo.Session)
  // decode the request body
  var c comment
  if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }

  // give the comment a unique ID and set the time
 var a [4]string
  a[0] = "KYD"
  a[1] = "EUR"
  a[2] = "QAR"
  a[3] = "SAR"
  
//responce, err :=http.Get("https://free.currencyconverterapi.com/api/v6/convert?q="+curr+"_EGP&compact=ultra") 

    //Read the json responce
 
  for i:=0; i<5 ;i++{
  responce, err :=http.Get("https://free.currencyconverterapi.com/api/v6/convert?q="+a[i]+"_EGP&compact=ultra") 
  if err !=nil {
    fmt.Printf("The requset faild -_- %s",err)
  }else{
  
    data,_ :=ioutil.ReadAll(responce.Body)
    host, port, err := net.SplitHostPort(string(data))
   c.Cur=host
   c.Rate=port
    fmt.Println(host)
    fmt.Println(port)
    fmt.Println(err)
    //return fmt.Sprintf(string(data))
  }

   

  

  c.ID = bson.NewObjectId()
  c.When = time.Now()
  
  // insert it into the database
  if err := db.DB("commentsapp").C("comments").Insert(&c); err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
  }
  // redirect to it
  http.Redirect(w, r, "/comments/"+c.ID.Hex(), http.StatusTemporaryRedirect)
}
}

//Red
func handleRead(w http.ResponseWriter, r *http.Request) {
  db := context.Get(r, "database").(*mgo.Session)
  // load the comments
  var comments []*comment
  if err := db.DB("commentsapp").C("comments").
    Find(nil).Sort("-when").Limit(100).All(&comments); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  // write it out
  if err := json.NewEncoder(w).Encode(comments); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
}



