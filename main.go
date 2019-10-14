package main

import (
	"flag"
	"fmt"

	"hi_bot/executor"
	"hi_bot/proxies"

	"github.com/spf13/viper"
)

type Proxy interface {
	Run(chan *executor.Message, chan []string)
}

func process(r chan []string, router *executor.Router, message *executor.Message) {
	pipe := make(chan *executor.Message)
	go router.Run(pipe)
	pipe <- message
	select {
	case <-pipe:
		r <- message.Result
	}
}

func main() {
	var err error

	configFile := flag.String("config", "config", "Config file")
	flag.Parse()

	viper.SetConfigName(*configFile)
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	err = viper.ReadInConfig()
	if err != nil {
		panic("Config not readed")
	}

	router := executor.NewRouter()

	// botFiles, err := ioutil.ReadDir("./bots")
	// if err != nil {
	// 	panic("Bots folder is missing")
	// }

	// re := regexp.MustCompile("--! (.+)")

	// for _, f := range botFiles {
	// 	file, _ := os.Open(filepath.Join("./bots", f.Name()))
	// 	defer file.Close()
	// 	reader := bufio.NewReader(file)
	// 	line, err := reader.ReadString('\n')
	// 	if err != nil {
	// 		continue
	// 	}

	// 	if submatch := re.FindSubmatch([]byte(line)); len(submatch) > 1 {
	// 		fmt.Println("Append bot:", string(submatch[1]), f.Name())
	// 		router.Append(string(submatch[1]), f.Name())
	// 	}
	// }

	discordMessageChan := make(chan *executor.Message)
	discordResponseChan := make(chan []string)
	discordQuitChan := make(chan int)
	discordProxy := proxies.NewDiscordProxy(viper.GetString("discord_bot_token"))
	go discordProxy.Run(discordMessageChan, discordResponseChan, discordQuitChan)

	fmt.Println("Start listen proxies")
	for {
		select {
		case message := <-discordMessageChan:
			fmt.Println("Receive message from DiscordProxy")
			process(discordResponseChan, router, message)
		}
	}

	defer func() {
		discordQuitChan <- 0
	}()
}
