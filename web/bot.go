package web

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"net/http"

	"github.com/spf13/viper"
)

func botGet(w http.ResponseWriter, r *http.Request) {
	responseBody := make(map[string]interface{})
	botName := r.URL.Query().Get("bot")

	if botName != "" {
		rf, _ := ioutil.ReadFile(filepath.Join(viper.GetString("bots_path"), botName))
		responseBody["text"] = string(rf)

	} else {
		responseBody["bots"] = make(map[string]string)

		botFiles, err := ioutil.ReadDir(viper.GetString("bots_path"))
		if err != nil {
			panic("Bots folder is missing")
		}

		re := regexp.MustCompile("--! (.+)")

		for _, f := range botFiles {
			file, _ := os.Open(filepath.Join(viper.GetString("bots_path"), f.Name()))
			defer file.Close()
			reader := bufio.NewReader(file)
			line, err := reader.ReadString('\n')
			if err != nil {
				continue
			}

			if submatch := re.FindSubmatch([]byte(line)); len(submatch) > 1 {
				template := string(submatch[1])

				responseBody["bots"].(map[string]string)[template] = f.Name()
			}
		}
	}

	jsonResponse(w, responseBody)
}

func botPost(w http.ResponseWriter, r *http.Request) {
	var requestData map[string]string
	json.NewDecoder(r.Body).Decode(&requestData)

	filename, filenameOk := requestData["filename"]
	source, sourceOk := requestData["source"]

	if filenameOk && sourceOk {
		filepath := filepath.Join(viper.GetString("bots_path"), filename)
		ioutil.WriteFile(filepath, []byte(source), 0644)
		jsonResponse(w, map[string]interface{}{"ok": true})
		return
	}

	jsonResponse(w, map[string]interface{}{"ok": false})
}

func botDelete(w http.ResponseWriter, r *http.Request) {}

func botHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		botGet(w, r)
	case "POST":
		botPost(w, r)
	case "DELETE":
		botDelete(w, r)
	}
}
