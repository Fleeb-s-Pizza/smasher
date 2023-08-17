package utils

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
	path := GetFilePath(endpoint, domain, hash)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func DownloadFile(endpoint string, domain string, hash string, url string) error {
	filepath := GetFilePath(endpoint, domain, hash)

	out, err := os.Create(filepath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}

	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error downloading file:", err)
		return err
	}

	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return err
	}

	return nil
}

func GetFileExtension(path string) string {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	fileBytes := make([]byte, 512)
	file.Read(fileBytes)
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

func GetFilePath(endpoint string, domain string, hash string) string {
	return "./cache/" + endpoint + "/" + domain + "/" + hash
}

func Contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}

	return false
}
