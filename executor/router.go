package executor

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/spf13/viper"
)

type Message struct {
	Author  string
	Message string
	Result  []string
}

func NewMessage(author, message string) *Message {
	m := new(Message)
	m.Author = author
	m.Message = message
	m.Result = make([]string, 1)
	return m
}

type Script struct {
	Template string `json:"template"`
	Name     string `json:"name"`
}

type Handlers struct {
	Scripts []Script `json:"scripts"`
}

type Router struct{}

func (r Router) Route(message *Message) {
	if jsonFile, err := os.Open(viper.GetString("bots_config")); err == nil {
		var handlers Handlers
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &handlers)

		for _, script := range handlers.Scripts {
			re := regexp.MustCompile(script.Template)

			if submatch := re.FindSubmatch([]byte(message.Message)); len(submatch) > 1 {
				response, ok := ExecHandler(message.Author, string(submatch[1]), script.Name)
				if ok {
					message.Result = append(message.Result, response)
				}
			}
		}
	}
}

func (r *Router) Run(pipe chan *Message) {
	for {
		select {
		case message := <-pipe:
			r.Route(message)
			close(pipe)
			return
		}
	}
}

func NewRouter() *Router {
	router := new(Router)
	return router
}
