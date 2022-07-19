package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/bwmarrin/discordgo"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/tripolious/discogo"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

// params
var (
	token   = flag.String("token", "", "bot token")
	channel = flag.String("channel", "", "channel bot listens to new messages")
)

var wg sync.WaitGroup

type Config struct {
	Debug   bool     `json:"debug"`
	Version string   `json:"version"`
	Amount  *big.Int `json:"amount"`
}

var config Config

func main() {
	log.Println("booting")
	// boot-up
	flag.Parse()
	ctx, cancelFunc, cancelChan := createLaunchContext()
	defer cancelFunc()

	// load current config
	var err error
	config, err = loadConfig()
	if err != nil {
		log.Fatalf("cant load config %s", err)
	}

	// start discord bot
	err = discogo.Boot(ctx, &wg, *token)
	if err != nil {
		log.Fatalf("booting discord bot failed %s", err)
	}

	// add handler to discord bot to consume messages and respond to them
	var handlers = []interface{}{
		consumeMessage,
	}
	err = discogo.AddHandlers(handlers)
	if err != nil {
		log.Printf("unable to add handlers: %s", err)
	}

	// start a fake logger, to print a message every 10 seconds if debugging is active
	go func() {
		for {
			if config.Debug {
				// send a message to discord channel
				err := discogo.SendMessage(*channel, "debugging active !")
				if err != nil {
					log.Printf("error sending message to discord %s", err)
				}
			}
			time.Sleep(10 * time.Second)
		}
	}()

	log.Println("successfully booted")

	for {
		select {
		case <-cancelChan:
			log.Println("shutting down...")
			wg.Wait()
			return
		}
	}
}

func consumeMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// we skip messages the bot did
	if m.Author.ID == s.State.User.ID {
		return
	}

	// you could add here also a check if the message is from a trusted user

	// we only want to consume messages from the defined channel
	if m.Message.ChannelID != *channel {
		return
	}

	// we ignore messages that not start with "!"
	if string(m.Content[0]) != "!" {
		return
	}

	stripped := strings.TrimSpace(string(m.Content[1:]))
	cmd := strings.Split(stripped, ".")
	switch cmd[0] {
	case "show": // handle !show
		var res []byte
		res, err := json.Marshal(config)
		if err != nil {
			sendMessageAndLogIfFailed("unable to load config: " + err.Error())
			return
		}
		sendMessageAndLogIfFailed(string(res))
	case "reload": // handle !reload
		var err error
		config, err = loadConfig()
		if err != nil {
			sendMessageAndLogIfFailed("unable to reload config: " + err.Error())
			return
		}
		sendMessageAndLogIfFailed("successfully reloaded config file!")
	case "config": // handle config.key newVal
		if len(cmd) < 2 {
			sendMessageAndLogIfFailed("you need to define a key you want to update (!config.xxx newVal)")
			return
		}
		updCmd := strings.Split(cmd[1], " ")
		if len(updCmd) < 2 {
			sendMessageAndLogIfFailed("you need to set a value for they key (!config.xxx newVal)")
			return
		}

		updField := strings.ToUpper(string(updCmd[0][0])) + updCmd[0][1:]

		// use the pointer so we can update config directly
		reflectValue := reflect.ValueOf(&config).Elem()
		field := reflectValue.FieldByName(updField)
		if field == (reflect.Value{}) {
			sendMessageAndLogIfFailed("field " + updField + " doesnt exist in config")
		}
		varType := field.Kind()

		// we handle bool, string and *big.Int - you can add here more types if you want to add different kinds to your config
		switch varType.String() {
		case "bool":
			newValue, err := strconv.ParseBool(updCmd[1])
			if err != nil {
				sendMessageAndLogIfFailed("cant update " + updField + " with value " + updCmd[1] + " - " + err.Error())
				return
			}
			field.Set(reflect.ValueOf(newValue))
			sendMessageAndLogIfFailed("updated " + updField + " to " + updCmd[1])
		case "string":
			newValue := updCmd[1]
			for i := 2; i < len(updCmd); i++ {
				newValue += " " + updCmd[i]
			}
			field.Set(reflect.ValueOf(newValue))
			sendMessageAndLogIfFailed("updated " + updField + " to " + newValue)
		case "ptr":
			ptrR := &field
			switch ptrR.String() {
			case "<*big.Int Value>":
				newValue, success := math.ParseBig256(updCmd[1])
				if !success {
					sendMessageAndLogIfFailed("cant update " + updField + " with value " + updCmd[1])
				}
				field.Set(reflect.ValueOf(newValue))
				sendMessageAndLogIfFailed("updated " + updField + " to " + newValue.String())
			}
		}
	default:
		sendMessageAndLogIfFailed("unknown command - available: !show / !reload / !config.key value")
	}
}

func sendMessageAndLogIfFailed(message string) {
	err := discogo.SendMessage(*channel, message)
	if err != nil {
		log.Printf("unable to send message: %s", message)
	}
}

func loadConfig() (Config, error) {
	var loadedConfig Config

	jsonConfig, err := os.Open("./config.json")
	if err != nil {
		return loadedConfig, err
	}

	jsonByteValue, err := ioutil.ReadAll(jsonConfig)
	if err != nil {
		return loadedConfig, err
	}

	err = json.Unmarshal(jsonByteValue, &loadedConfig)
	if err != nil {
		return loadedConfig, err
	}

	return loadedConfig, nil
}

func createLaunchContext() (context.Context, func(), chan bool) {
	interruptChan := make(chan os.Signal, 1)
	canceledChanChan := make(chan bool, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)
	ctx, cancelCtx := context.WithCancel(context.Background())
	go func() {
		defer close(interruptChan)
		<-interruptChan
		cancelCtx()
		canceledChanChan <- true
	}()
	cancel := func() {
		cancelCtx()
		close(canceledChanChan)
	}
	return ctx, cancel, canceledChanChan
}
