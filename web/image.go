package main

import (
	"encoding/json"
	"fmt"
	"github.com/h2non/bimg"
	"net/http"
	"os"
	"strconv"
	"time"
)

func HandleImageRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")

	if r.Method != http.MethodGet {
		errorJson, _ := json.Marshal(Error{
			Message: "Method not allowed",
			Status:  http.StatusMethodNotAllowed,
		})
		http.Error(w, string(errorJson), http.StatusMethodNotAllowed)
		return
	}

	if !r.URL.Query().Has("url") || r.URL.Query().Get("url") == "" {
		errorJson, _ := json.Marshal(Error{
			Message: "Missing or empty url parameter",
			Status:  http.StatusBadRequest,
		})
		http.Error(w, string(errorJson), http.StatusBadRequest)
		return
	}

	url := r.URL.Query().Get("url")
	if !CheckIfStringUrl(url) {
		errorJson, _ := json.Marshal(Error{
			Message: "Invalid url parameter",
			Status:  http.StatusBadRequest,
		})
		http.Error(w, string(errorJson), http.StatusBadRequest)
		return
	}

	var err error
	width, height := 0, 0

	if r.URL.Query().Has("width") {
		width, err = strconv.Atoi(r.URL.Query().Get("width"))
		if err != nil {
			errorJson, _ := json.Marshal(Error{
				Message: "Invalid width parameter",
				Status:  http.StatusBadRequest,
			})
			http.Error(w, string(errorJson), http.StatusBadRequest)
			return
		}
	}

	if r.URL.Query().Has("height") {
		height, err = strconv.Atoi(r.URL.Query().Get("height"))
		if err != nil {
			errorJson, _ := json.Marshal(Error{
				Message: "Invalid height parameter",
				Status:  http.StatusBadRequest,
			})
			http.Error(w, string(errorJson), http.StatusBadRequest)
			return
		}
	}

	if width > 2048 || height > 2048 {
		errorJson, _ := json.Marshal(Error{
			Message: "Maximum width or height is 2048",
			Status:  http.StatusBadRequest,
		})
		http.Error(w, string(errorJson), http.StatusBadRequest)
		return
	}

	webp := false
	if r.URL.Query().Has("webp") {
		webp, err = strconv.ParseBool(r.URL.Query().Get("webp"))
		if err != nil {
			errorJson, _ := json.Marshal(Error{
				Message: "Invalid webp parameter",
				Status:  http.StatusBadRequest,
			})
			http.Error(w, string(errorJson), http.StatusBadRequest)
			return
		}
	}

	quality := 100
	if r.URL.Query().Has("quality") {
		quality, err = strconv.Atoi(r.URL.Query().Get("quality"))
		if err != nil {
			errorJson, _ := json.Marshal(Error{
				Message: "Invalid quality parameter",
				Status:  http.StatusBadRequest,
			})
			http.Error(w, string(errorJson), http.StatusBadRequest)
			return
		}
	}

	// url to md5
	hashedUrl := HashUrl(url) + "-" + strconv.Itoa(width) + "-" + strconv.Itoa(height) + "-" + strconv.FormatBool(webp) + "-" + strconv.Itoa(quality)
	zeroHashedUrl := HashUrl(url) + "-0-0-" + strconv.FormatBool(webp) + "-" + strconv.Itoa(quality)
	domain := ExtractDomainFromUrl(url)

	// create folder
	err = CreateFolder("image", domain)
	if err != nil {
		errorJson, _ := json.Marshal(Error{
			Message: "Error creating folder",
			Status:  http.StatusInternalServerError,
		})
		http.Error(w, string(errorJson), http.StatusInternalServerError)
		panic(err)
		return
	}

	// check if file exist
	if !CheckIfFileExist("image", domain, zeroHashedUrl) {
		err := DownloadFile("image", domain, zeroHashedUrl, url)
		if err != nil {
			errorJson, _ := json.Marshal(Error{
				Message: "Error downloading file",
				Status:  http.StatusInternalServerError,
			})

			http.Error(w, string(errorJson), http.StatusInternalServerError)
			panic(err)
			return
		}
	}

	buffer, err := os.ReadFile(GetFilePath("image", domain, hashedUrl))
	if err != nil {
		panic(err)
		return
	}

	img := bimg.NewImage(buffer)

	if width != 0 && height != 0 {
		if width < 1 || height < 1 {
			errorJson, _ := json.Marshal(Error{
				Message: "Minimum width or height is 1",
				Status:  http.StatusBadRequest,
			})
			http.Error(w, string(errorJson), http.StatusBadRequest)
			return
		}

		buffer, err := img.Resize(width, height)
		if err != nil {
			errorJson, _ := json.Marshal(Error{
				Message: "Error resizing image",
				Status:  http.StatusInternalServerError,
			})

			http.Error(w, string(errorJson), http.StatusInternalServerError)
			fmt.Println(hashedUrl)
			panic(err)
			return
		}

		img = bimg.NewImage(buffer)
	}

	if webp {
		imgBuffer, err := img.Convert(bimg.WEBP)
		if err != nil {
			errorJson, _ := json.Marshal(Error{
				Message: "Error converting to webp",
				Status:  http.StatusInternalServerError,
			})

			http.Error(w, string(errorJson), http.StatusInternalServerError)
			fmt.Println(hashedUrl)
			panic(err)
			return
		}

		img = bimg.NewImage(imgBuffer)
	}

	processed, err := img.Process(bimg.Options{Quality: quality})
	if err != nil {
		panic(err)
		return
	}

	err = os.WriteFile(GetFilePath("image", domain, hashedUrl), processed, 0644)
	if err != nil {
		panic(err)
		return
	}

	http.ServeFile(w, r, GetFilePath("image", domain, hashedUrl))
	err = os.Chtimes(GetFilePath("image", domain, hashedUrl), time.Now(), time.Now())
}
