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

	http.Handle("/ui", http.FileServer(http.Dir("./ui")))

	fmt.Println("Server started at " + os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT"))

	err = http.ListenAndServe(os.Getenv("SERVER_HOST")+":"+os.Getenv("SERVER_PORT"), nil)
	if err != nil {
		fmt.Println("Error starting server: ", err)
		panic(err)
		return
	}
}
