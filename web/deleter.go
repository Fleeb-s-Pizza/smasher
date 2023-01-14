package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func RemoveOldFiles() {
	for {
		filepath.Walk("./cache", func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			if time.Since(info.ModTime()) > 7*24*time.Hour {
				err := os.Remove(path)
				if err != nil {
					fmt.Println("Error deleting file: ", path, err)
				} else {
					fmt.Println("Deleted file: ", path)
				}
			}
			return nil
		})
		time.Sleep(30 * time.Minute)
	}
}
