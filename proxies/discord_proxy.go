package proxies

// https://discordapp.com/developers/applications - for append application
// https://discordapp.com/oauth2/authorize?&client_id=CLIENTID&scope=bot&permissions=8 - for auth bot to chanel

import (
	"fmt"
	"hi_bot/executor"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type DiscordProxy struct {
	Token        string
	messageChan  chan *executor.Message
	responseChan chan []string
}

func (dp *DiscordProxy) onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	fmt.Println("DiscordProxy send message")
	dp.messageChan <- executor.NewMessage(m.Author.Username, m.Content)
	responseMessages := <-dp.responseChan

	for _, rMessage := range responseMessages {
		s.ChannelMessageSend(m.ChannelID, rMessage)
	}
}

func (dp *DiscordProxy) Run(messageChan chan *executor.Message, responseChan chan []string) {
	dp.messageChan = messageChan
	dp.responseChan = responseChan

	dg, err := discordgo.New("Bot " + dp.Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(dp.onMessage)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	fmt.Println("Airhorn is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func NewDiscordProxy(token string) *DiscordProxy {
	dp := new(DiscordProxy)
	dp.Token = token
	return dp
}
