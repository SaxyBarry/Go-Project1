package main

import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
)

type Greeting struct {
    Language string `json:"lang"`
    Hello string `json:"greeting"`
}

// Creates a Path for /hello at http://localhost:8080/hello that writes Hello! to the page and adds security
func helloHandler(w http.ResponseWriter, r *http.Request) {
	// Any page other than /hello returns a 404
    if r.URL.Path != "/hello" {
        http.Error(w, "404 not found.", http.StatusNotFound)
        return
    }
	// Only Accepts GET requests: curl http://localhost:8080/hello
    if r.Method != "GET" {
        http.Error(w, "Method is not supported.", http.StatusNotFound)
        return
    }


    fmt.Fprintf(w, "Hello!")
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
    en := &Greeting{
        Language: "en",
        Hello: "Hello, World!",
    }
    sp := &Greeting{
        Language: "sp",
        Hello: "¡Hola, Mundo!",
    }
    de := &Greeting{
        Language: "de",
        Hello: "Hallo, Welt!",
    }
	json.NewEncoder(w).Encode([]*Greeting{en, sp, de})
	return
}

// Handles the submit from http://localhost:8080/form.html
func formHandler(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        fmt.Fprintf(w, "ParseForm() err: %v", err)
        return
    }
    fmt.Fprintf(w, "POST request successful\n")
    name := r.FormValue("name")
    address := r.FormValue("address")

    fmt.Fprintf(w, "Name = %s\n", name)
    fmt.Fprintf(w, "Address = %s\n", address)
}

func main(){
    // Where the files are located that will be served
    fileServer := http.FileServer(http.Dir("./static"))
    // The web path for the files being served
    http.Handle("/", fileServer) 
    http.HandleFunc("/form", formHandler)
	http.HandleFunc("/hello", helloHandler)
    http.HandleFunc("/hellojson", helloJsonHandler)

	fmt.Printf("Starting server at port 8080\n")
	// Starts a server at Port 8080 if possible, else displays an error.
	// http://localhost:8080 
	if err := http.ListenAndServe(":8080", nil); err != nil{
		log.Fatal(err)
	}
}