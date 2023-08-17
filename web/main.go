package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
		return
	}

	// Remove old files (than 7 days) at startup (every 30 minutes)
	go RemoveOldFiles()

	http.HandleFunc("/image", func(writer http.ResponseWriter, request *http.Request) {
		HandleImageRequest(writer, request)
	})

	// UI Section
	http.HandleFunc("/ui", func(writer http.ResponseWriter, request *http.Request) {
		HandleUIRequest(writer, request)
	})

	http.HandleFunc("/css/style.css", func(writer http.ResponseWriter, request *http.Request) {
		HandleUIRequest(writer, request)
	})

	http.HandleFunc("/js/app.js", func(writer http.ResponseWriter, request *http.Request) {
		HandleUIRequest(writer, request)
	})

	fmt.Println("Server started at " + os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT"))

	err = http.ListenAndServe(os.Getenv("SERVER_HOST")+":"+os.Getenv("SERVER_PORT"), nil)
	if err != nil {
		fmt.Println("Error starting server: ", err)
		panic(err)
		return
	}
}
