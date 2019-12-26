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
	"strings"
	"time"
)

type ServerInfo struct {
	ServerId *int `json:"serverId"`
	//AccountId    int     `json:"accountId"`
	Ip      *string `json:"ip"`
	Port    *int    `json:"port"`
	Players *string `json:"players"`
	/* Info         string  `json:"info"`
	Pastebin     string  `json:"pastebin"`
	Version      string  `json:"version"`
	PrivateBeta  bool    `json:"privateBeta"`
	FriendlyFire bool    `json:"friendlyFire"`
	Modded       bool    `json:"modded"`
	Whitelist    bool    `json:"whitelist"`
	IsoCode      string  `json:"isoCode"`
	CountryCode  string  `json:"countryCode"`
	Latitude     float32 `json:"latitude"`
	Longitude    float32 `json:"longitude"`
	Official     string  `json:"official"`
	OfficialCode int     `json:"officialCode"`
	DisplaySection int     `json:"displaySection"` */
}

type ServerConfigInfo struct {
	Name     *string `json:"name"`
	Ip       *string `json:"ip"`
	Port     *int    `json:"port"`
	ServerId *int    `json:"serverId"`
	BotToken *string `json:"discordUserBotToken"`
}

type ConfigFile struct {
	UpdateTime *int               `json:"updateTime"`
	Logger     *bool              `json:"logger"`
	Servers    []ServerConfigInfo `json:"servers"`
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

		resp, err := http.Get("https://api.scpslgame.com/lobbylist.php?format=json")
		if err != nil {
			log.Println("Failed when getting data from server:" + err.Error())
			log.Println("Try again in 10 sec")
			time.Sleep(time.Second * 10)
			continue
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("Failed when getting data from server:" + err.Error())
			log.Println("Try again in 10 sec")
			time.Sleep(time.Second * 10)
			continue
		}

		data := body

		defer file.Close()
		_, err = file.Write(data)

		if err != nil {
			log.Fatalln(err)
		}

		if *config.Logger {
			log.Println("Successfully getting data from server")
		}

		time.Sleep(time.Second * time.Duration(*config.UpdateTime))
	}
}

func getServerInfo() []ServerInfo {
	data, err := ioutil.ReadFile("cache.json")
	if err != nil {
		fmt.Println("Error while reading cache: " + err.Error())
		return nil
	}

	var serverInfo []ServerInfo

	err = json.Unmarshal(data, &serverInfo)
	if err != nil {
		log.Println("Error while parsing cache: " + err.Error())
		return nil
	}

	return serverInfo
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
			log.Println("Try again in 10 sec")
			time.Sleep(time.Second * 10)
			continue
		}

		for _, element := range serverArray {
			if info.ServerId == nil {
				if *element.Ip == *info.Ip && *element.Port == *info.Port {
					players = *element.Players
					if *config.Logger {
						log.Println(*info.Name + " found and online " + players)
					}
					break
				}
			} else if *info.ServerId == *element.ServerId {
				players = *element.Players
				if *config.Logger {
					log.Println(*info.Name + " found and online " + players)
				}
				break
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
				Name: "server startup",
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
			log.Println(err)
		}
		time.Sleep(time.Second * time.Duration(*config.UpdateTime))
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

	if config.UpdateTime == nil {
		log.Println("You need to set update time in config.json")

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
			if (info.Ip == nil || info.Port == nil) || (*info.Ip == "") {
				log.Println("You need to set server id or ip and port for " + *info.Name + " in config.json")

				return nil
			}
		}

		if info.BotToken == nil || *info.BotToken == "" {
			log.Println("You need to set bot token for " + *info.Name + " in config.json")

			return nil
		}
	}

	return config.Servers
}
