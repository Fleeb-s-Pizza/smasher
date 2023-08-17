package handlers

import (
	"encoding/json"
	"net/http"
	"os"
)

var build BuildInfo

func LoadBuildInfo() {
	file, err := os.Open("build.json")
	if err != nil {
		panic(err)
		return
	}

	defer file.Close()

	var buildInfo BuildInfo
	err = json.NewDecoder(file).Decode(&buildInfo)
	if err != nil {
		panic(err)
		return
	}

	build = buildInfo
}

func HandleInfoRequest(w http.ResponseWriter, r *http.Request) {
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

	projectInfo := ProjectInfo{
		Name:   "Smasher",
		Author: "Vladimír Urík",
		Build:  build,
	}

	projectInfoJson, _ := json.Marshal(projectInfo)
	_, err := w.Write(projectInfoJson)
	if err != nil {
		panic(err)
		return
	}
}
