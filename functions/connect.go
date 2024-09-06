package actions

import (
	"time"

	"github.com/gorilla/websocket"
)

func internalConnect(tkn, act_description, act_status string, act_type any) {
	switch act_type {
	case "Game":
		{
			act_type = 0
			break
		}
	case "Listening":
		{
			act_type = 2
			break
		}
	case "Watching":
		{
			act_type = 3
			break
		}
	case "Competing":
		{
			act_type = 5
			break
		}
	}
	ws, _, _ := websocket.DefaultDialer.Dial("wss://gateway.discord.gg/?v=10&encoding=json", nil)
	payload := map[string]interface{}{
		"op": 2,
		"d": map[string]interface{}{
			"token":   tkn,
			"intents": 0,
			"properties": map[string]interface{}{
				"os":      "linux",
				"browser": "http_requests",
				"device":  "discord",
			},
			"presence": map[string]interface{}{
				"activities": []map[string]interface{}{
					{
						"name": act_description,
						"type": act_type,
					},
				},
				"status": act_status,
			},
		},
	}
	ws.WriteJSON(payload)
}

func Connect(tkn, act_name, act_status string, act_type any, timer *time.Ticker, stop chan struct{}) {
	internalConnect(tkn, act_name, act_status, act_type)
	for {
		select {
		case <-stop:
			{
				timer.Stop()
				return
			}
		case <-timer.C:
			{
				internalConnect(tkn, act_name, act_status, act_type)
			}
		}
	}
}
