package main

import (
	"encoding/json"
	"fmt"
	"github.com/discord/lilliput"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var EncodeOptions = map[string]map[int]int{
	".jpeg": map[int]int{lilliput.JpegQuality: 85},
	".png":  map[int]int{lilliput.PngCompression: 7},
	".webp": map[int]int{lilliput.WebpQuality: 85},
}

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
	zeroHashedUrl := HashUrl(url) + "-0-0"
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

	if width != 0 && height != 0 {
		if width < 10 || height < 10 {
			errorJson, _ := json.Marshal(Error{
				Message: "Minimum width or height is 10",
				Status:  http.StatusBadRequest,
			})
			http.Error(w, string(errorJson), http.StatusBadRequest)
			return
		}

		err := resizeImage(GetFilePath("image", domain, zeroHashedUrl), GetFilePath("image", domain, hashedUrl), width, height)
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

	http.ServeFile(w, r, GetFilePath("image", domain, hashedUrl))
	err = os.Chtimes(GetFilePath("image", domain, hashedUrl), time.Now(), time.Now())
}

func resizeImage(srcImg, destImg string, width, height int) error {
	inputBuf, err := os.ReadFile(srcImg)
	if err != nil {
		panic(err)
		return err
	}

	decoder, err := lilliput.NewDecoder(inputBuf)
	if err != nil {
		panic(err)
		return err
	}
	defer decoder.Close()

	ops := lilliput.NewImageOps(8192)
	defer ops.Close()

	outputType := "." + strings.ToLower(decoder.Description())
	resizeOps := &lilliput.ImageOptions{
		Width:                width,
		Height:               height,
		ResizeMethod:         lilliput.ImageOpsResize,
		NormalizeOrientation: true,
		EncodeOptions:        EncodeOptions[outputType],
	}

	outputImg := make([]byte, 50*1024*1024)

	outputImg, err = ops.Transform(decoder, resizeOps, outputImg)
	if err != nil {
		panic(err)
		return err
	}

	dest, err := os.Create(destImg)
	if err != nil {
		panic(err)
		return err
	}
	defer dest.Close()

	err = os.WriteFile(destImg, outputImg, 0400)
	if err != nil {
		fmt.Printf("error writing out resized image, %s\n", err)
		os.Exit(1)
	}

	if err != nil {
		panic(err)
		return err
	}

	return nil
}
