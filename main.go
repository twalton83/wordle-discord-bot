package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	// "strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var botPrefix = "?"
var lines []string
var wordOfTheDay string

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

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	lines, err = readLines("sgb-words.txt")
	if err != nil {
		fmt.Println("Error reading lines", err)
		return
	}
	wordOfTheDay = pickWord(lines)
	fmt.Print(wordOfTheDay)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	discord.Close()
}

func messageCreate(session *discordgo.Session, m *discordgo.MessageCreate) {
	parsedMessage := strings.Split(string(m.Content), " ")
	fmt.Print(parsedMessage)
	if m.Author.ID == session.State.User.ID || string(m.Content[0]) != botPrefix {
		return
	}

	if parsedMessage[0] == "?word" {
		sendDailyWord(wordOfTheDay, session, m.Message )
	}
 
	if parsedMessage[0] == "?guess" {
		result := calculateGuess(parsedMessage[1])

		if result {
			session.ChannelMessageSend(m.ChannelID, "Correct!")
		} else {
			session.ChannelMessageSend(m.ChannelID, "Try again!")
		}
	}
}


func pickWord(lines []string) string {
	randsource := rand.NewSource(time.Now().UnixNano())
	randgenerator := rand.New(randsource)
	randNum := randgenerator.Intn(len(lines))

  return lines[randNum]
}

func readLines(path string) ([]string, error) {
	// This file is so large that it may not make sense to create an array in memory
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

func sendDailyWord(word string, s *discordgo.Session, m *discordgo.Message){
	s.ChannelMessageSend(m.ChannelID, word)
}
  
func calculateGuess(guess string) bool{
	guessArr := strings.Split(guess, "")
	guessMap := validateLetters(guessArr)

	fmt.Print(guessMap)

	if guessMap["incorrect"] == true {

		return false
	} else {
		return true
	}
}

func validateLetters(guessArr []string) map[string]interface{} {
	validationMap := make(map[string]interface{})

	wordOfTheDayArr := strings.Split(wordOfTheDay, "")
// 2 is assigned when the correct placement and letter
// 1 is assigned when it's the correct letter
// 0 is assigned when the letter is not present at all
	for i, letter := range guessArr {
		if(wordOfTheDayArr[i] != letter){
			validationMap[letter] = 2
		} else if (contains(wordOfTheDayArr, letter)) {
			validationMap[letter] = 1
		} else {
			validationMap[letter] =  0
			validationMap["incorrect"] = true
		}
	}
	return validationMap
}

func generateEmbed(){
	
}

func contains(s []string, l string) bool {
    for _, letter := range s {
        if letter == l {
            return true
        }
    }
    return false
}

