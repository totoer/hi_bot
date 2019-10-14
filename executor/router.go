package executor

import (
	"bufio"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"

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

type Router struct {
	Handlers map[string]string
}

func (r *Router) Append(template, script string) {
	r.Handlers[template] = script
}

func (r Router) Route(message *Message) {
	for template, script := range r.Handlers {
		re := regexp.MustCompile(template)

		if submatch := re.FindSubmatch([]byte(message.Message)); len(submatch) > 1 {
			response, ok := ExecHandler(message.Author, string(submatch[1]), script)
			if ok {
				message.Result = append(message.Result, response)
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

func (r *Router) updateBots() {
	for {
		botFiles, err := ioutil.ReadDir(viper.GetString("bots_path"))
		if err != nil {
			panic("Bots folder is missing")
		}

		re := regexp.MustCompile("--! (.+)")
		var existsTemplates []string

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

				existsTemplates = append(existsTemplates, template)

				if _, ok := r.Handlers[template]; !ok {
					r.Handlers[template] = f.Name()
				}
			}
		}

		for template, _ := range r.Handlers {
			ok := false
			for _, existsTemplate := range existsTemplates {
				if template == existsTemplate {
					ok = true
					break
				}
			}

			if !ok {
				delete(r.Handlers, template)
			}
		}

		time.Sleep(5 * time.Second)
	}
}

func NewRouter() *Router {
	router := new(Router)
	router.Handlers = make(map[string]string)

	go router.updateBots()

	return router
}
