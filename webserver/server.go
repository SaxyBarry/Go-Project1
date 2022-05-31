package main

import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

type Greeting struct {
    Language string `json:"lang"`
    Hello string `json:"greeting"`
}

// Creates a Path for /hello at http://localhost:8080/hellojson that writes Hello! to the page and adds security
func helloJsonHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusCreated)
    w.Header().Set("Content-Type", "application/json")
	// Any page other than /hello returns a 404
    if r.URL.Path != "/hellojson" {
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }
	// Only Accepts GET requests: curl http://localhost:8080/hellojson
    if r.Method != "GET" {
        http.Error(w, "Method is not supported.", http.StatusNotFound)
        return
    }
    // A map of maps
    // resp := make(map[string]map[string]string)
    // greetingsMap := map[string]string{"english":"Hello, World!", "german" : "Hallo, Welt!", "spanish" : "¡Hola, Mundo!"}
	// resp["Greeting"] = greetingsMap 
    query, err := queryAll()
    if err != nil{
        fmt.Fprintf(w, "An Error has occurred %s", err)
        log.Fatal(err)
    }
	json.NewEncoder(w).Encode(query)
}


func queryAll() ([]Greeting, error){
    db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/go_p1")
     if err != nil {
        return nil, err
    }

    defer db.Close()
    var greetings []Greeting
    selectS, err := db.Query("SELECT * FROM greetings")
    if err != nil {
        return nil, err
    }

    for selectS.Next(){
        var greeting Greeting
        selectS.Scan(&greeting.Language, &greeting.Hello)
        greetings = append(greetings, greeting)
    }
    return greetings, nil
}

func insertGreeting(lang string, greet string) (error){
    db, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/go_p1")
     if err != nil {
        return err
    }

    defer db.Close()
    ins := fmt.Sprintf("INSERT INTO greetings(language, greetings) VALUES('%s', '%s')", lang, greet)
    _, errr := db.Exec(ins)
    
    if errr != nil {
        return errr
    }
    fmt.Printf("Inserted")
    return nil
}

// Handles the submit from http://localhost:8080/form.html
func formHandler(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        fmt.Fprintf(w, "ParseForm() err: %v", err)
        return
    }
    // fmt.Fprintf(w, "POST request successful\n")
    lang := r.FormValue("lang")
    greet := r.FormValue("greet")

    err := insertGreeting(lang, greet)

    if err != nil {
        log.Fatal(err)
    }

    greetingQ, err := queryAll()
    if err != nil {
        log.Fatal(err)
    }
    for _, x := range greetingQ{
        fmt.Fprintf(w, "Language: %s, Greeting: %s\n", x.Language, x.Hello)
    }


}

func main(){
    // Where the files are located that will be served
    fileServer := http.FileServer(http.Dir("./static"))
    // The web path for the files being served
    http.Handle("/", fileServer) 
    http.HandleFunc("/form", formHandler)
    http.HandleFunc("/hellojson", helloJsonHandler)

	fmt.Printf("Starting server at port 8080\n")
	// Starts a server at Port 8080 if possible, else displays an error.
	// http://localhost:8080 
	if err := http.ListenAndServe(":8080", nil); err != nil{
		log.Fatal(err)
	}
}