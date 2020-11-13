package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "io"
    "os"
    "time"
    "encoding/json"
    "strconv"
    "strings"
    "math/rand"
)

//setting up mongo
var session, _ = mgo.Dial(os.Getenv("MONGO_HOST"))
var c = session.DB("TutDb").C("ToDo")

type book struct {
    ISBN        int           `json:"isbn" bson:"isbn,omitempty"`
    Title       string        `json:"title"`
    Author      string        `json:"author"`
    Date        time.Time     `json:"date"`
    Available   bool          `json:available`
}

  func getBooks(w http.ResponseWriter, r *http.Request) {
    var res []book

    _ = c.Find(nil).All(&res)

    json.NewEncoder(w).Encode(res)
}

 func getBooksISBN(w http.ResponseWriter, r *http.Request) {
    var res []book
    vars := mux.Vars(r)
    str_isbn := vars["isbn"]

    if str_isbn != "" {
        var result book
        isbn, _ := strconv.Atoi(str_isbn)

        _ = c.Find(bson.M{"isbn": isbn}).One(&result)
        res = append(res, result)
    }

    json.NewEncoder(w).Encode(res)
}

func getBooksTitle(w http.ResponseWriter, r *http.Request) {
   var res []book
   vars := mux.Vars(r)
   title := vars["title"]

   title = strings.Replace(title,"\"","",-2) //remove the quotes

   if title != "" {
       _ = c.Find (bson.M{"title": title}).All(&res)

   }

   json.NewEncoder(w).Encode(res)
}

func addBook(w http.ResponseWriter, r *http.Request) {
    isbn := rand.Intn(9999999999)

    _ = c.Insert(book{
        isbn,
        r.FormValue("title"),
        r.FormValue("author"),
        time.Now(),
        false,
    })

    result := book{}
    _ = c.Find(bson.M{"title": r.FormValue("title")}).One(&result)
    json.NewEncoder(w).Encode(result)
}

func delBook(w http.ResponseWriter, r*http.Request) {
    vars := mux.Vars(r)
    str_isbn := vars["isbn"]
    isbn, _ := strconv.Atoi(str_isbn)

    err := c.Remove(bson.M{"isbn":isbn})
    if err != nil  {
        json.NewEncoder(w).Encode(err.Error())
    } else {
        io.WriteString(w, "{result: 'OK'}")
    }
}

func changeAvailability(w http.ResponseWriter, r*http.Request) {
    vars := mux.Vars(r)

    str_isbn := vars["isbn"]
    isbn, _ := strconv.Atoi(str_isbn)

    var book book
    _ = c.Find(bson.M{"isbn": isbn}).One(&book)
    
    err := c.Update(bson.M{"isbn": isbn}, bson.M{"$set": bson.M{"available": !book.Available}})
    if err != nil {
        io.WriteString(w, `{"updated": false, "error": `+err.Error()+`}`)
    } else {
        io.WriteString(w, `{"updated": true}`)
    }
}

func main() {
    session.SetMode(mgo.Monotonic, true)
    defer session.Close()
    router := mux.NewRouter()
    router.HandleFunc("/grc/", addBook).Methods("POST", "PUT")

    router.HandleFunc("/grc/", getBooks).Methods("GET")
    router.HandleFunc("/grc/SearchByISBN/{isbn}", getBooksISBN).Methods("GET")
    router.HandleFunc("/grc/SearchByTitle/{title}", getBooksTitle).Methods("GET")

    router.HandleFunc("/grc/{isbn}", changeAvailability).Methods("PATCH")

    router.HandleFunc("/grc/{isbn}", delBook).Methods("DELETE")

    log.Fatal(http.ListenAndServe(":8000", router))
}
