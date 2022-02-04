package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	// "strconv"
	"strings"
	"syscall"
	"time"

	// "strings"

	"github.com/bwmarrin/discordgo"
	"github.com/enescakir/emoji"
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
	lines, err = readLines("test.txt")
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
	if m.Author.ID == session.State.User.ID || string(m.Content[0]) != botPrefix {
		return
	}


	if parsedMessage[0] == "?word" {
		sendDailyWord(wordOfTheDay, session, m.Message )
	}
 
	if parsedMessage[0] == "?guess" {
		result := calculateGuess(parsedMessage[1])
		embed := generateEmbed(*m.Author, result, parsedMessage[1])
		_, err := session.ChannelMessageSendEmbed(m.ChannelID, embed)

		if err != nil{
			fmt.Print(err)
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
  
type results struct {
	incorrect bool
	guessMap []int
}

func calculateGuess(guess string) results{
	guessArr := strings.Split(guess, "")
	guessMap := validateLetters(guessArr)

	res := results{incorrect : guessMap.incorrect, guessMap: guessMap.stringPositionResults }

	return res
}

type wordPositions struct {
		positions []int
}

func generateWOODMap() map[string]wordPositions {
	split := strings.Split(wordOfTheDay, "")
	m := make(map[string]wordPositions)

	// loop over word, push positions of letters to slice 
	for i, letter := range split {
		_, keyExists := m[letter]
		if(!keyExists){
		positionsArr := make([]int, 0)
		positionsArr = append(positionsArr, i)
			m[letter] = struct {
				positions []int
			}{
				positions: positionsArr,
			}
		} else {
			positionsArr := m[letter].positions
			positionsArr = append(positionsArr, i)
			m[letter] = struct {
				positions []int
			}{
				positions: positionsArr,
			}
		}
	}
	return m
}

type validationResults struct {
	stringPositionResults []int
	incorrect bool
}

// TODO: if there's a guess with two of the same letters and only one instance exists, check for length of positions array. if it doesn't match the current, FAIL that position
func validateLetters(guessArr []string) validationResults {
	validationMap := validationResults{
		stringPositionResults: make([]int, 0),
		incorrect: false,
	}
	// this should create a map showing the different positions in a word
	w := generateWOODMap()

// 2 is assigned when the correct placement and letter
// 1 is assigned when it's the correct letter
// 0 is assigned when the letter is not present at all
	for currentLetterIndex, letter := range guessArr {
		_, letterExists := w[letter]
		if(letterExists){
			// loop over letter's positions
			for _, position := range w[letter].positions {
				if(currentLetterIndex == position){
					validationMap.stringPositionResults = append(validationMap.stringPositionResults, 2)
					break
				} else {
					// If the letter does not match the position
					validationMap.stringPositionResults = append(validationMap.stringPositionResults, 1)
					break
				}
			}
		} else {
			validationMap.stringPositionResults = append(validationMap.stringPositionResults, 0)
			validationMap.incorrect = true
		}
	}
	return validationMap
}

func generateEmbed(user discordgo.User, result results, userGuess string) *discordgo.MessageEmbed {
	var color int
	if result.incorrect {
		color = 0xff0000
	} else {
		color = 0x00ff00
	}

	emojis := createEmojiString(result.guessMap)

	embed := &discordgo.MessageEmbed{
    Author:      &discordgo.MessageEmbedAuthor{},
    Color:       color,
    Description: user.Mention(),
    Fields: []*discordgo.MessageEmbedField{
			 {
            Name:   "Your Guess",
            Value:  userGuess,
            Inline: false,
        },
        {
            Name:   "Result",
            Value:  emojis,
            Inline: false,
        },
    },
    Timestamp: time.Now().Format(time.RFC3339), 
	}
	return embed
}

func createEmojiString(guessMap []int) string {
	var emojiString = ""
	for _, value := range guessMap {
		if(value == 2){
			emojiString += emoji.GreenCircle.String()
		} else if (value == 1){
			emojiString += emoji.YellowCircle.String()
		} else {
			emojiString += emoji.WhiteCircle.String()
		}
	}
	return emojiString
}

func contains(s []string, l string) bool {
    for _, letter := range s {
        if letter == l {
            return true
        }
    }
    return false
}


