package main

import (
    "encoding/json"
    "math/rand"
    "net/http"
    //"time"
    "fmt"
)

type Result struct {
	Num1 int `json:"first_num"`
	Num2 int `json:"second_num"`
    Result int `json:"result"`
}

type floatResult struct {
	Num1 int `json:"first_num"`
	Num2 int `json:"second_num"`
	FloatResult float32 `json:"result"`
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
    response := map[string]string{"info": "This is a simple math API"}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func firstHandler(w http.ResponseWriter, r *http.Request) {
    // rand.Seed(time.Now().UnixNano())
    number := rand.Intn(100)
    result := Result{Num1: number}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

func secondHandler(w http.ResponseWriter, r *http.Request) {
    // rand.Seed(time.Now().UnixNano())
    number := rand.Intn(100)

    result := Result{Num2: number}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

func addHandler(w http.ResponseWriter, r *http.Request) {
    // rand.Seed(time.Now().UnixNano())
    number1 := rand.Intn(100)
    number2 := rand.Intn(100)
    result := number1 + number2

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(Result{Num1: number1, Num2: number2, Result: result})
}

func subHandler(w http.ResponseWriter, r *http.Request) {
    // rand.Seed(time.Now().UnixNano())
    number1 := rand.Intn(100)
    number2 := rand.Intn(100)
    result := number1 - number2

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(Result{Num1: number1, Num2: number2, Result: result})
}

func mulHandler(w http.ResponseWriter, r *http.Request) {
    // rand.Seed(time.Now().UnixNano())
    number1 := rand.Intn(100)
    number2 := rand.Intn(100)
    result := number1 * number2

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(Result{Num1: number1, Num2: number2, Result: result})
}

func divHandler(w http.ResponseWriter, r *http.Request) {
    // rand.Seed(time.Now().UnixNano())
    number1 := rand.Intn(100)
    number2 := rand.Intn(100)
    result := float32(number1) / float32(number2)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(floatResult{Num1: number1, Num2: number2, FloatResult: result})
}

func main() {
    http.HandleFunc("/info", infoHandler)
    http.HandleFunc("/first", firstHandler)
    http.HandleFunc("/second", secondHandler)
    http.HandleFunc("/add", addHandler)
    http.HandleFunc("/sub", subHandler)
    http.HandleFunc("/mul", mulHandler)
    http.HandleFunc("/div", divHandler)

    fmt.Println("Server is running on port 1234...")
    http.ListenAndServe(":1234", nil)
}