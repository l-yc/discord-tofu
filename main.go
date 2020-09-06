package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/BurntSushi/toml"
	"github.com/l-yc/discord-tofu/answer"
)

type Config struct {
	Token				string
	ClientID		string
}

type Flags struct {
	ConfigFile	string
}

var (
	flags	Flags
)

func init() {
	//flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&flags.ConfigFile, "c", "config.toml", "Config File")
	flag.Parse()
}

// Reads info from config file
func ReadConfig() Config {
	configfile := flags.ConfigFile
	log.Println("Reading from config file:", configfile)
	_, err := os.Stat(configfile)
	if err != nil {
		log.Fatal("Config file is missing:", configfile)
	}

	var config Config
	if _, err := toml.DecodeFile(configfile, &config); err != nil {
		log.Fatal(err)
	}
	return config
}

func main() {
	config := ReadConfig()

	discord, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		log.Println("Error creating discord session:", err)
		return
	}
	// Invite the bot
	fmt.Println("Invite the bot to your server by visiting the following link:")
	fmt.Printf("https://discord.com/oauth2/authorize?client_id=%s&scope=bot", config.ClientID)
	fmt.Println()

	// Register the MessageCreate func as a callback for MessageCreate events.
	discord.AddHandler(answer.MessageCreate)

	// In this example, we only care about receiving message events.
	discord.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	err = discord.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}
