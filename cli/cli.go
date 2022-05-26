package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"encoding/json"
)

// Structure that the Webserver is returning, Json tags help represent the keys
type Greeting struct {
    Language string `json:"lang"`
    Hello string `json:"greeting"`
}

type Color string
// ANSI Escape Sequences 
// Setting constant color variables
const (
	ColorBlack  Color = "\u001b[30m"
	ColorRed          = "\u001b[31m"
	ColorGreen        = "\u001b[32m"
	ColorYellow       = "\u001b[33m"
	ColorBlue         = "\u001b[34m"
	ColorReset        = "\u001b[0m"
)

// Prints a message in the given color
func colorize(color Color, message string) {
	fmt.Println(string(color), message, string(ColorReset))
}

// User interaction enabled color picker, defaults to terminal color
func pickColor() Color{
	fmt.Printf("What is your favorite color?\n")
	var color string
	fmt.Scanf("%s", &color)
	switch{
	case color == "Blue" || color == "blue" || color == "b":
		return ColorBlue
	case color == "Green" || color == "green" || color == "g":
		return ColorGreen
	case color == "Red" || color == "red" || color == "r":
		return ColorRed
	case color == "Yellow" || color == "yellow" || color == "y":
		return ColorYellow
	default:
		fmt.Printf("%s is not a valid color...\n", color)
	}
	return ColorReset
}

// Requests the golang web server for a json object of message greetings, returns the object in a mutable form
func requestServer() ([]Greeting, error){
	// I know the expected return of the request, so I create an object of greetings to accept that return
	var greetings []Greeting
	// Send a GET request
	resp, err := http.Get("http://localhost:8080/hellojson")
	if err != nil {
		return greetings, err
    }
	// Close the connection 
	defer resp.Body.Close()
	// We have to decode the JSON GET request so we can use it 
	json.NewDecoder(resp.Body).Decode(&greetings)
	// If we made it this far, we have an object and no errors ocurred 
    return greetings, nil
}

// Finds the correct message to display
func getLangOption(lang string, arr []Greeting) string{
	for _, x := range arr{
		if x.Language == lang{
			return x.Hello
		}
	}
	return ""
}

func main(){
	
	// flag.{{variable type}}({{name of flag}}, {{default value}}, {{-h output for help}})
	// str := flag.String("username", "root", "Specify login username")
	// int_flag := flag.Int("port", 3006, "specify the port")
	useColor := flag.Bool("color", false, "display colorized output\nValid colors are: Red, Green, Blue or Yellow")
	lang := flag.String("lang", "en", "language that the message is displayed in\nen, de, sp")
	flag.Parse()
	
	if *useColor {
		// Gets the color option from the user
		color := pickColor()
		// Requests the server for messages
		body, err := requestServer()
		if err != nil{
			log.Fatal(fmt.Sprintf(string(color), err))
		}
		// Finds the correct message based off the language set in the flag
		message := getLangOption(*lang, body)
		// Prints prettily 
		colorize(color, message)
	}else{
		// Requests the server for messages
		body, err := requestServer()
		if err != nil{
			log.Fatal(err)
		}
		// Finds the correct message based off the language set in the flag
		message := getLangOption(*lang, body)
		fmt.Println(message)
	}
}