package main

import (
	"encoding/json"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
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

	// url to md5
	hashedUrl := HashUrl(url) + "-" + strconv.Itoa(width) + "-" + strconv.Itoa(height)
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
	if !CheckIfFileExist("image", domain, hashedUrl) {
		err := DownloadFile("image", domain, hashedUrl, url)
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

	if width != 0 && height != 0 {
		err := resizeImage("./cache/image/"+domain+"/"+hashedUrl, "./cache/image/"+domain+"/"+hashedUrl, uint(width), uint(height))
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
	}

	http.ServeFile(w, r, "./cache/image/"+domain+"/"+hashedUrl)
	err = os.Chtimes("./cache/image/"+domain+"/"+hashedUrl, time.Now(), time.Now())
}

func resizeImage(srcImg, destImg string, width, height uint) error {
	file, err := os.Open(srcImg)
	if err != nil {
		panic(err)
		return err
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
		return err
	}

	// Resize the image
	img = resize.Resize(width, height, img, resize.Lanczos3)

	// Create the destination file
	dest, err := os.Create(destImg)
	if err != nil {
		panic(err)
		return err
	}
	defer dest.Close()

	ext := GetFileExtension(srcImg)
	switch ext {
	case "jpg", "jpeg":
		err = jpeg.Encode(dest, img, nil)
	case "png":
		err = png.Encode(dest, img)
	case "gif":
		err = gif.Encode(dest, img, nil)
	}
	if err != nil {
		panic(err)
		return err
	}

	return nil
}
