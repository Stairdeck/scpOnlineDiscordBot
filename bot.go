package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"github.com/andersfylling/disgord"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Response struct {
	Success  *bool        `json:"Success"`
	Error    *string      `json:"Error"`
	Cooldown *int         `json:"Cooldown"`
	Servers  []ServerInfo `json:"Servers"`
}

type ServerInfo struct {
	ID      *int    `json:"ID"`
	Players *string `json:"Players"`
}

type ServerConfigInfo struct {
	Name     *string `json:"name"`
	ServerId *int    `json:"serverId"`
	BotToken *string `json:"discordUserBotToken"`
}

type ConfigFile struct {
	Logger    *bool              `json:"logger"`
	AccountId *int               `json:"accountId"`
	ApiKey    *string            `json:"apiKey"`
	Servers   []ServerConfigInfo `json:"servers"`
}

var config ConfigFile

func main() {
	var serversConfig = initConfig()

	if serversConfig == nil {
		fmt.Print("Press 'Enter' to exit...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}

	go initServerInfo()

	for _, info := range serversConfig {
		go initBots(info)
	}

	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return
}

func initServerInfo() {
	for {
		file, err := os.Create("cache.json")
		if err != nil {
			fmt.Println("Unable to create cache file:", err)
			os.Exit(1)
			return
		}

		resp, err := http.Get("https://api.scpslgame.com/serverinfo.php?key=" + *config.ApiKey + "&id=" + strconv.Itoa(*config.AccountId) + "&players=true")
		if err != nil {
			log.Println("Failed when getting data from server:" + err.Error())
			log.Println("Try again in 10 sec")
			time.Sleep(time.Second * 10)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Failed when getting data from server:" + err.Error())
			log.Println("Try again in 10 sec")
			time.Sleep(time.Second * 10)
			continue
		}

		resp.Body.Close()

		data := body

		var response Response

		err = json.Unmarshal(data, &response)
		if err != nil {
			log.Println("Error while parsing response: " + err.Error())
			log.Println("Try again in 15 sec")
			time.Sleep(time.Second * 15)
			continue
		}

		if !*response.Success {
			log.Println("Error while getting response: " + *response.Error)
			log.Println("Try again in 15 sec")
			time.Sleep(time.Second * 15)
			continue
		}

		_, err = file.Write(data)
		file.Close()

		if err != nil {
			log.Fatalln(err)
		}

		if *config.Logger {
			log.Println("Successfully getting data from server")
		}

		time.Sleep(time.Second * time.Duration(*response.Cooldown+1))
	}
}

func getServerInfo() []ServerInfo {
	data, err := ioutil.ReadFile("cache.json")
	if err != nil {
		fmt.Println("Error while reading cache: " + err.Error())
		return nil
	}

	if len(data) == 0 {
		return nil
	}

	var response Response

	err = json.Unmarshal(data, &response)
	if err != nil {
		log.Println("Error while parsing cache: " + err.Error())
		return nil
	}

	return response.Servers
}

func initBots(info ServerConfigInfo) {
	log.Println("Start connection to bot of server " + *info.Name)

	var discordConfig disgord.Config

	if *config.Logger {
		discordConfig = disgord.Config{
			BotToken: *info.BotToken,
			Logger:   disgord.DefaultLogger(false),
		}
	} else {
		discordConfig = disgord.Config{
			BotToken: *info.BotToken,
		}
	}

	client, err := disgord.NewClient(discordConfig)

	defer client.StayConnectedUntilInterrupted(context.Background())

	if err != nil {
		log.Println(err)
	}

	client.Ready(func() {
		log.Println("Successfully connected to " + *info.Name + "!")
		go setStatus(*client, info)
	})
}

func setStatus(client disgord.Client, info ServerConfigInfo) {
	for {
		var activity = disgord.Activity{}
		var botStatus string

		var isOnline = true
		var players = ""

		var serverArray = getServerInfo()

		if serverArray == nil {
			if *config.Logger {
				log.Println("Empty cache. Maybe info isn't updated yet? Trying again in 1 sec")
			}
			time.Sleep(time.Second * 1)
			continue
		}

		for _, element := range serverArray {
			if *info.ServerId == *element.ID {
				if element.Players != nil {
					players = *element.Players
					if *config.Logger {
						log.Println(*info.Name + " found and online " + players)

					}
					break
				} else {
					players = ""
				}
			}
		}

		if players == "" {
			if *config.Logger {
				log.Printf(*info.Name + " is offline")
			}
			isOnline = false
		}

		/* TODO
		var EmojiActivity = disgord.ActivityEmoji{
			Name: "",
			//ID:            disgord.NewSnowflake(),
		}*/

		if isOnline {
			activity = disgord.Activity{
				Name: players,
				Type: 0,
			}
			if strings.Split(players, "/")[0] == "0" {
				botStatus = disgord.StatusIdle
			} else {
				botStatus = disgord.StatusOnline
			}
		} else {
			activity = disgord.Activity{
				Name: "offline",
				Type: 3,
			}
			botStatus = disgord.StatusDnd
		}

		var status = disgord.UpdateStatusPayload{
			AFK:    false,
			Game:   &activity,
			Status: botStatus,
		}

		err := client.UpdateStatus(&status)

		if err != nil {
			log.Println("Fail when updating bot status of " + *info.Name + "")
			log.Println("Trying again in 15 sec")
		}
		time.Sleep(time.Second * time.Duration(15))
	}
}

func initConfig() []ServerConfigInfo {

	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		fmt.Println("Error while reading config.json: " + err.Error())
		return nil
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		log.Println("Error while parsing json: " + err.Error())
		return nil
	}

	if config.Logger == nil {
		log.Println("You need to set logger to true/false in config.json")

		return nil
	}

	for _, info := range config.Servers {
		if info.Name == nil || *info.Name == "" {
			log.Println("You need to set name for all servers")

			return nil
		}

		if info.ServerId == nil {
			log.Println("You need to set server id for " + *info.Name + " in config.json")
			return nil
		}

		if info.BotToken == nil || *info.BotToken == "" {
			log.Println("You need to set bot token for " + *info.Name + " in config.json")

			return nil
		}
	}

	return config.Servers
}
