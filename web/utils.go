package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

func HashUrl(url string) string {
	data := []byte(url)
	return fmt.Sprintf("%x", md5.Sum(data))
}

func CheckIfStringUrl(url string) bool {
	re := regexp.MustCompile(`^(?:http|https):\/\/([^\/]+)`)
	match := re.FindStringSubmatch(url)

	if len(match) < 2 {
		return false
	}

	return true
}

func ExtractDomainFromUrl(url string) string {
	re := regexp.MustCompile(`^(?:http|https):\/\/([^\/]+)`)
	match := re.FindStringSubmatch(url)

	return match[1]
}

func CreateFolder(endpoint string, domain string) error {
	path := "./cache/" + endpoint + "/" + domain
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			panic(err)
			return err
		}
	}

	return nil
}

func CheckIfFileExist(endpoint string, domain string, hash string) bool {
	path := "./cache/" + endpoint + "/" + domain + "/" + hash
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func DownloadFile(endpoint string, domain string, hash string, url string) (error, string) {
	filepath := "./cache/" + endpoint + "/" + domain + "/" + hash

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return err, ""
	}

	ext := GetFileExtension(resp.Body)
	filepath = filepath + "." + ext

	out, err := os.Create(filepath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err, ""
	}

	defer resp.Body.Close()
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err, ""
	}

	return nil, filepath
}

func GetFileExtension(read io.ReadCloser) string {
	fileBytes := make([]byte, 512)
	fileType := http.DetectContentType(fileBytes)

	switch fileType {
	case "image/jpeg":
		return "jpg"
	case "image/gif":
		return "gif"
	case "image/png":
		return "png"
	case "image/bmp":
		return "bmp"
	case "image/tiff":
		return "tiff"
	case "image/webp":
		return "webp"
	case "application/pdf":
		return "pdf"
	default:
		return "unknown"
	}
}
