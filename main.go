package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var botPrefix = "?"

func main() {
	envErr := godotenv.Load()
  if envErr != nil {
    log.Fatal("Error loading .env file")
  }

	botKey := os.Getenv("BOT_TOKEN")

	discord, err := discordgo.New("Bot " + botKey)
	if err != nil {
    log.Fatal("Error connecting to Discord")
  }

	discord.AddHandler(messageCreate)
	discord.Identify.Intents = discordgo.IntentsGuildMessages

	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening connection", err)
		return
	}
	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discord.Close()
}

func messageCreate(session *discordgo.Session, m *discordgo.MessageCreate){
	if m.Author.ID == session.State.User.ID || string(m.Content[0]) != botPrefix {
        return
  }
	session.ChannelMessageSend(m.ChannelID, "Hello World!")
}


