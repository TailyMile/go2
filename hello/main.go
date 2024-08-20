package main


import (
	"fmt"
	"log"
	"net/http"
	"time"
	"io"
	"os"
)

func main()  {
	fmt.Println("Hello World")

	//EZ user HTTP client
	client := &http.Client{Timeout: time.Second}
	resOne, err := client.Get("http://www.google.ru")
	if err != nil {
		log.Fatal(err)
	}
	// res.Status return text status message
	// res.StatusCode return code status message
	// GET = 200, POST = 201, PUT\PATCH = 202, DELETE = 204
	fmt.Println(resOne.Status, resOne.StatusCode, resOne.Request.URL)

	body, _ := io.ReadAll(resOne.Body)
	resOne.Body.Close()
	file, err := os.Create("out.txt")
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.Write(body)
	if err != nil {
		log.Fatal(err)
	}
}