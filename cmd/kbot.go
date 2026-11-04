/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	telebot "gopkg.in/telebot.v4"
)

var (
	TeleToken = os.Getenv("TELE_TOKEN")
)

// TrafficSignal represents a simulated traffic light signal.
type TrafficSignal struct {
	On bool
}

var kbotCmd = &cobra.Command{
	Use:     "kbot",
	Aliases: []string{"start"},
	Short:   "Telegram bot for controlling traffic light signals",
	Long: `A Telegram bot that simulates traffic light signals.

Usage:
  /s red    - Toggle red light
  /s amber  - Toggle amber light
  /s green  - Toggle green light
  hello     - Get a greeting from the bot`,
	Run: func(cmd *cobra.Command, args []string) {

		if TeleToken == "" {
			log.Fatal("TELE_TOKEN environment variable is not set")
		}

		fmt.Printf("kbot %s started\n", AppVersion)

		kbot, err := telebot.NewBot(telebot.Settings{
			Token:  TeleToken,
			Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		})

		if err != nil {
			log.Fatalf("Please check TELE_TOKEN env variable. %s", err)
		}

		// Initialize simulated traffic signals.
		trafficSignals := map[string]*TrafficSignal{
			"red":   {On: false},
			"amber": {On: false},
			"green": {On: false},
		}

		var mu sync.Mutex

		kbot.Handle(telebot.OnText, func(m telebot.Context) error {

			log.Printf("Received message: %s", m.Text())

			payload := strings.TrimSpace(m.Message().Payload)

			if payload == "" {
				fields := strings.Fields(strings.TrimSpace(m.Text()))

				if len(fields) > 1 && strings.HasPrefix(fields[0], "/s") {
					payload = fields[1]
				} else if len(fields) == 1 {
					payload = fields[0]
				}
			}

			switch payload {

			case "hello":
				return m.Send(fmt.Sprintf("Hello I'm Kbot %s!", AppVersion))

			case "red", "amber", "green":

				mu.Lock()
				signal := trafficSignals[payload]
				signal.On = !signal.On
				state := signal.On
				mu.Unlock()

				status := map[bool]string{
					true:  "on",
					false: "off",
				}[state]

				log.Printf("Traffic signal %s switched %s", payload, status)

				return m.Send(fmt.Sprintf(
					"Switched %s light %s",
					payload,
					status,
				))

			default:
				return m.Send("Usage: /s red|amber|green or hello")
			}
		})

		kbot.Start()
	},
}

func init() {
	rootCmd.AddCommand(kbotCmd)
}