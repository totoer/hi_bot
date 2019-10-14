package proxies

// https://discordapp.com/developers/applications - for append application
// https://discordapp.com/oauth2/authorize?&client_id=CLIENTID&scope=bot&permissions=8 - for auth bot to chanel

import (
	"hi_bot/executor"
	"log"

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

	log.Println("DiscordProxy send message")
	dp.messageChan <- executor.NewMessage(m.Author.Username, m.Content)
	responseMessages := <-dp.responseChan

	for _, rMessage := range responseMessages {
		s.ChannelMessageSend(m.ChannelID, rMessage)
	}
}

func (dp *DiscordProxy) Run(messageChan chan *executor.Message, responseChan chan []string, quitChan chan int) {
	dp.messageChan = messageChan
	dp.responseChan = responseChan

	dg, err := discordgo.New("Bot " + dp.Token)
	if err != nil {
		log.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(dp.onMessage)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}

	<-quitChan

	// Cleanly close down the Discord session.
	dg.Close()
}

func NewDiscordProxy(token string) *DiscordProxy {
	dp := new(DiscordProxy)
	dp.Token = token
	return dp
}
