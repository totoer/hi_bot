package executor

import (
	"regexp"
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

func NewRouter() *Router {
	router := new(Router)
	router.Handlers = make(map[string]string)
	return router
}
