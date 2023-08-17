package main

import (
	"fmt"
	"github.com/Fleeb-s-Pizza/smasher/web/handlers"
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

	handlers.LoadBuildInfo()

	http.HandleFunc("/info", func(writer http.ResponseWriter, request *http.Request) {
		handlers.HandleInfoRequest(writer, request)
	})

	http.HandleFunc("/image", func(writer http.ResponseWriter, request *http.Request) {
		handlers.HandleImageRequest(writer, request)
	})

	// UI Section
	http.HandleFunc("/ui", func(writer http.ResponseWriter, request *http.Request) {
		handlers.HandleUIRequest(writer, request, "/ui")
	})

	http.HandleFunc("/css/style.css", func(writer http.ResponseWriter, request *http.Request) {
		handlers.HandleUIRequest(writer, request, "/css/style.css")
	})

	http.HandleFunc("/js/app.js", func(writer http.ResponseWriter, request *http.Request) {
		handlers.HandleUIRequest(writer, request, "/js/app.js")
	})

	fmt.Println("Server started at " + os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT"))

	err = http.ListenAndServe(os.Getenv("SERVER_HOST")+":"+os.Getenv("SERVER_PORT"), nil)
	if err != nil {
		fmt.Println("Error starting server: ", err)
		panic(err)
		return
	}
}
