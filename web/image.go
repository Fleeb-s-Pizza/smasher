package main

import (
	"github.com/nfnt/resize"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func HandleImageRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Query().Has("url") || r.URL.Query().Get("url") == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	url := r.URL.Query().Get("url")
	if !CheckIfStringUrl(url) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var err error
	width, height := 0, 0

	if r.URL.Query().Has("width") {
		width, err = strconv.Atoi(r.URL.Query().Get("width"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if r.URL.Query().Has("height") {
		height, err = strconv.Atoi(r.URL.Query().Get("height"))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if width > 2048 || height > 2048 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// url to md5
	hashedUrl := HashUrl(url) + "-" + strconv.Itoa(width) + "-" + strconv.Itoa(height)
	domain := ExtractDomainFromUrl(url)

	// create folder
	err = CreateFolder("image", domain)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// check if file exist
	if !CheckIfFileExist("image", domain, hashedUrl) {
		err := DownloadFile("image", domain, hashedUrl, url)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	if width != 0 && height != 0 {
		err := resizeImage("./cache/image/"+domain+"/"+hashedUrl, "./cache/image/"+domain+"/"+hashedUrl, uint(width), uint(height))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	http.ServeFile(w, r, "./cache/image/"+domain+"/"+hashedUrl)
}

func resizeImage(srcImg, destImg string, width, height uint) error {
	file, err := os.Open(srcImg)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode the image
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Resize the image
	img = resize.Resize(width, height, img, resize.Lanczos3)

	// Create the destination file
	dest, err := os.Create(destImg)
	if err != nil {
		return err
	}
	defer dest.Close()

	// Save the resized image
	ext := filepath.Ext(destImg)
	switch ext {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(dest, img, nil)
	case ".png":
		err = png.Encode(dest, img)
	case ".gif":
		err = gif.Encode(dest, img, nil)
	}
	if err != nil {
		return err
	}
	return nil
}
