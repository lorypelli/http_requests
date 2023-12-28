package main

import (
	b "bytes"
	j "encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func connect(tkn string, act_description string, act_type any, act_status string) {
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

func checkStatus(tkn string, act_name string, act_type any, act_status string, timer *time.Ticker, stop chan struct{}) {
	connect(tkn, act_name, act_type, act_status)
	for {
		select {
		case <-stop:
			{
				timer.Stop()
				return
			}
		case <-timer.C:
			{
				connect(tkn, act_name, act_type, act_status)
			}
		}
	}
}

func main() {
	a := app.New()
	login := a.NewWindow("Login")
	program := a.NewWindow("http_requests")
	bot := a.NewWindow("bot_info")
	msg_list := a.NewWindow("msg_list")
	act_type := widget.NewSelect([]string{"Game", "Listening", "Watching", "Competing"}, nil)
	act_type.SetSelected("Watching")
	act_description := widget.NewEntry()
	act_description.Text = "http_requests"
	act_description.SetPlaceHolder("Insert detailed activity")
	act_status := widget.NewSelect([]string{"online", "idle", "dnd"}, nil)
	act_status.SetSelected("dnd")
	req, _ := http.NewRequest("GET", "https://raw.githubusercontent.com/lorypelli/http_requests/main/icon.png", nil)
	c := &http.Client{}
	res, _ := c.Do(req)
	_, ok := a.(desktop.App)
	if ok {
		homeDirectory, _ := os.UserHomeDir()
		path := fmt.Sprintf("%s/http_requests/", homeDirectory)
		os.MkdirAll(path, os.ModePerm)
		filePath := filepath.Join(path, "icon.png")
		file, _ := os.Create(filePath)
		io.Copy(file, res.Body)
		file, _ = os.Open(filePath)
		stats, _ := os.Stat(filePath)
		size := stats.Size()
		fileBytes := make([]byte, size)
		fileSlice := fileBytes[:]
		file.Read(fileSlice)
		icon := fyne.NewStaticResource("icon.png", fileBytes)
		login.SetIcon(icon)
		program.SetIcon(icon)
		bot.SetIcon(icon)
		msg_list.SetIcon(icon)
	}
	login.Resize(fyne.NewSize(670, 165))
	program.Resize(fyne.NewSize(400, 240))
	bot.Resize(fyne.NewSize(400, 200))
	msg_list.Resize(fyne.NewSize(1280, 720))
	login.SetFixedSize(true)
	program.SetFixedSize(true)
	bot.SetFixedSize(true)
	msg_list.SetFixedSize(true)
	login.CenterOnScreen()
	program.CenterOnScreen()
	bot.CenterOnScreen()
	msg_list.CenterOnScreen()
	show := false
	login.SetCloseIntercept(func() {
		a.Quit()
	})
	bot.SetCloseIntercept(func() {
		bot.Hide()
	})
	msg_list.SetCloseIntercept(func() {
		msg_list.Hide()
	})
	tkn := widget.NewPasswordEntry()
	tkn.SetPlaceHolder("Insert bot token")
	validate := widget.NewButton("Validate", func() {
		req, err := http.NewRequest("POST", "https://discord.com/api/v10/auth/login", nil)
		if err != nil {
			dialog.ShowError(err, login)
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
		c := &http.Client{}
		res, err := c.Do(req)
		timer := time.NewTicker(120 * time.Second)
		if err != nil {
			dialog.ShowError(err, login)
		} else if res.StatusCode != 200 {
			var body struct {
				Message string
			}
			bytes, _ := io.ReadAll(res.Body)
			j.Unmarshal(bytes, &body)
			dialog.ShowInformation("Error", body.Message, login)
		} else {
			login.Hide()
			logout := func(b bool) {
				if b {
					stop := make(chan struct{})
					go checkStatus(tkn.Text, act_description.Text, act_type.Selected, act_status.Selected, timer, stop)
					close(stop)
					login.Show()
					program.Hide()
					msg_list.Hide()
					bot.Hide()
				} else {
					show = false
				}
			}
			program.SetCloseIntercept(func() {
				if !show {
					dialog.ShowConfirm("Logout", "Are you sure you want to logout?", logout, program)
					show = true
				}
			})
			go checkStatus(tkn.Text, act_description.Text, act_type.Selected, act_status.Selected, timer, nil)
			bot_info := widget.NewButton("Bot Info", func() {
				req, err := http.NewRequest("GET", "https://discord.com/api/v10/users/@me", nil)
				if err != nil {
					dialog.ShowError(err, program)
				}
				req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
				c := &http.Client{}
				res, err := c.Do(req)
				if err != nil {
					dialog.ShowError(err, program)
				} else if res.StatusCode != 200 {
					var body struct {
						Message string
					}
					bytes, _ := io.ReadAll(res.Body)
					j.Unmarshal(bytes, &body)
					dialog.ShowInformation("Error", body.Message, program)
				} else {
					var bots struct {
						Id       string
						Username string
						Avatar   string
					}
					bytes, _ := io.ReadAll(res.Body)
					j.Unmarshal(bytes, &bots)
					req, err = http.NewRequest("GET", "https://discord.com/api/v10/users/@me/guilds", nil)
					if err != nil {
						dialog.ShowError(err, program)
					}
					req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
					c = &http.Client{}
					res, err = c.Do(req)
					if err != nil {
						dialog.ShowError(err, program)
					} else if res.StatusCode != 200 {
						var body struct {
							Message string
						}
						bytes, _ := io.ReadAll(res.Body)
						j.Unmarshal(bytes, &body)
						dialog.ShowInformation("Error", body.Message, program)
					} else {
						var guilds []struct{}
						bytes, _ = io.ReadAll(res.Body)
						j.Unmarshal(bytes, &guilds)
						var img fyne.Resource
						if bots.Avatar == "" {
							img, err = fyne.LoadResourceFromURLString("https://cdn.discordapp.com/embed/avatars/0.png")
						} else {
							img, err = fyne.LoadResourceFromURLString(fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", bots.Id, bots.Avatar))
						}
						if err != nil {
							dialog.ShowError(err, program)
						} else {
							img_box := canvas.NewImageFromResource(img)
							img_box.FillMode = canvas.ImageFillContain
							img_box.SetMinSize(fyne.NewSquareSize(32))
							bot.SetContent(container.NewCenter(container.NewVBox(img_box, widget.NewLabel(fmt.Sprintf("Username: %s", bots.Username)), widget.NewLabel(fmt.Sprintf("ID: %s", bots.Id)), widget.NewLabel(fmt.Sprintf("Server Count: %d", len(guilds))))))
							bot.Show()
						}
					}
				}
			})
			logout_btn := widget.NewButton("Logout", func() {
				dialog.ShowConfirm("Logout", "Are you sure you want to logout?", logout, program)
				show = true
			})
			navbar := container.NewHBox(bot_info, layout.NewSpacer(), logout_btn)
			chn_id := widget.NewEntry()
			chn_id.SetPlaceHolder("Insert channel ID")
			show_msg_list := widget.NewButton("Show Message List", func() {
				req, err := http.NewRequest("GET", fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages?limit=100", chn_id.Text), nil)
				if err != nil {
					dialog.ShowError(err, program)
				}
				req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
				c := &http.Client{}
				res, err := c.Do(req)
				if err != nil {
					dialog.ShowError(err, program)
				} else if res.StatusCode != 200 {
					var body struct {
						Message string
					}
					bytes, _ := io.ReadAll(res.Body)
					j.Unmarshal(bytes, &body)
					dialog.ShowInformation("Error", body.Message, program)
				} else {
					type user struct {
						Id       string
						Username string
						Avatar   string
						Bot      bool
					}
					type msg struct {
						Author  user
						Content string
					}
					bytes, _ := io.ReadAll(res.Body)
					msgs := []msg{}
					j.Unmarshal(bytes, &msgs)
					msgs_container := container.NewVBox()
					users_container := container.NewVBox()
					var urls []string
					var users []user
					bot_logo, _ := fyne.LoadResourceFromURLString("https://cdn.emojidex.com/emoji/seal/Bot_tag.png")
					for i := len(msgs) - 1; i >= 0; i-- {
						var avatar fyne.Resource
						url := fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", msgs[i].Author.Id, msgs[i].Author.Avatar)
						foundAvatar := false
						foundUser := false
						for c := 0; c < len(urls); c++ {
							if urls[c] == url {
								foundAvatar = true
							}
						}
						urls = append(urls, url)
						if !foundAvatar {
							if msgs[i].Author.Avatar == "" {
								avatar, _ = fyne.LoadResourceFromURLString("https://cdn.discordapp.com/embed/avatars/0.png")
							} else {
								avatar, _ = fyne.LoadResourceFromURLString(url)
							}
						}
						avatar_box := canvas.NewImageFromResource(avatar)
						avatar_box.FillMode = canvas.ImageFillContain
						avatar_box.SetMinSize(fyne.NewSquareSize(32))
						content := widget.NewLabel(msgs[i].Content)
						if len(msgs[i].Content) == 0 {
							content = widget.NewLabel("No Content!")
							content.TextStyle.Bold = true
							content.TextStyle.Italic = true
						}
						for c := 0; c < len(users); c++ {
							if users[c].Username == msgs[i].Author.Username {
								foundUser = true
							}
						}
						msgs_container.Add(container.NewHBox(widget.NewLabel(fmt.Sprintf("%s :", msgs[i].Author.Username)), content))
						user := container.NewHBox(avatar_box, widget.NewLabel(msgs[i].Author.Username))
						if !foundUser {
							if msgs[i].Author.Bot {
								bot_logo_box := canvas.NewImageFromResource(bot_logo)
								bot_logo_box.FillMode = canvas.ImageFillContain
								bot_logo_box.SetMinSize(fyne.NewSquareSize(32))
								user.Add(bot_logo_box)
							}
							users_container.Add(container.NewBorder(nil, nil, user, nil))
						}
						users = append(users, msgs[i].Author)
					}
					msgs_scroll := container.NewScroll(msgs_container)
					users_scroll := container.NewScroll(users_container)
					split_container := container.NewHSplit(msgs_scroll, users_scroll)
					split_container.SetOffset(0.8)
					msg_list.SetContent(split_container)
					msg_list.Show()
				}
			})
			navbar_edit := container.NewHBox(bot_info, layout.NewSpacer(), show_msg_list, layout.NewSpacer(), logout_btn)
			msg := widget.NewMultiLineEntry()
			msg.Wrapping = fyne.TextWrapWord
			msg.SetPlaceHolder("Insert message")
			count := widget.NewLabel("0")
			msg_box := container.NewStack(msg, container.NewHBox(layout.NewSpacer(), container.NewBorder(nil, count, nil, nil)))
			msg_id := widget.NewEntry()
			msg_id.SetPlaceHolder("Insert message ID")
			guild_id := widget.NewEntry()
			guild_id.SetPlaceHolder("Insert guild ID")
			chn_type := widget.NewSelect([]string{"Text", "Voice", "Stage", "Announcement", "Forum", "Media"}, nil)
			chn_type.SetSelected("Text")
			chn_name := widget.NewEntry()
			chn_name.SetPlaceHolder("Insert channel name")
			usr_id := widget.NewEntry()
			usr_id.SetPlaceHolder("Insert user ID")
			thread_name := widget.NewEntry()
			thread_name.SetPlaceHolder("Insert thread name")
			role_id := widget.NewEntry()
			role_id.SetPlaceHolder("Insert role ID")
			role_name := widget.NewEntry()
			role_name.SetPlaceHolder("Insert role name")
			confirm_action := widget.NewButton("Send", func() {
				body := map[string]interface{}{
					"content": msg.Text,
				}
				json, _ := j.Marshal(body)
				req, err := http.NewRequest("POST", fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages", chn_id.Text), b.NewBuffer(json))
				if err != nil {
					dialog.ShowError(err, program)
				}
				req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
				req.Header.Add("Content-Type", "application/json")
				c := &http.Client{}
				res, err := c.Do(req)
				if err != nil {
					dialog.ShowError(err, program)
				} else if res.StatusCode != 200 {
					var body struct {
						Message string
					}
					bytes, _ := io.ReadAll(res.Body)
					j.Unmarshal(bytes, &body)
					dialog.ShowInformation("Error", body.Message, program)
				} else {
					dialog.ShowInformation("Success", "The message has been successfully sent!", program)
				}
			})
			actions := widget.NewSelect([]string{"Write a message", "Edit a message", "Pin a message", "Create a channel", "Edit a channel", "Create a thread", "Delete a channel", "Delete a message", "Unpin a message", "Kick a user", "Ban a user", "Unban a user", "Create a role", "Edit a role", "Delete a role", "Add a role to a member", "Remove a role from a member"}, nil)
			msg.OnChanged = func(s string) {
				count.SetText(fmt.Sprint(len(s)))
				if len(msg.Text) > 4096 {
					confirm_action.Disable()
				}
			}
			actions.OnChanged = func(s string) {
				switch s {
				case "Write a message":
					{
						program.SetContent(container.NewBorder(navbar_edit, nil, nil, nil, container.NewVBox(chn_id, actions, msg_box, confirm_action)))
						program.Resize(fyne.NewSize(400, 240))
						confirm_action.SetText("Send")
						confirm_action.OnTapped = func() {
							body := map[string]interface{}{
								"content": msg.Text,
							}
							json, _ := j.Marshal(body)
							req, err := http.NewRequest("POST", fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages", chn_id.Text), b.NewBuffer(json))
							if err != nil {
								dialog.ShowError(err, program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							req.Header.Add("Content-Type", "application/json")
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, program)
							} else if res.StatusCode != 200 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, program)
							} else {
								dialog.ShowInformation("Success", "The message has been successfully sent!", program)
							}
						}
						break
					}
				case "Edit a message":
					{
						program.Resize(fyne.NewSize(400, 270))
						program.SetContent(container.NewBorder(navbar_edit, nil, nil, nil, container.NewVBox(chn_id, actions, msg_id, msg_box, confirm_action)))
						confirm_action.SetText("Edit")
						confirm_action.OnTapped = func() {
							body := map[string]interface{}{
								"content": msg.Text,
							}
							json, _ := j.Marshal(body)
							req, err := http.NewRequest("PATCH", fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages/%s", chn_id.Text, msg_id.Text), b.NewBuffer(json))
							if err != nil {
								dialog.ShowError(err, program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							req.Header.Add("Content-Type", "application/json")
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, program)
							} else if res.StatusCode != 200 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, program)
							} else {
								dialog.ShowInformation("Success", "The message has been successfully edited!", program)
							}
						}
						break
					}
				case "Pin a message":
					{
						program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(chn_id, actions, msg_id, confirm_action)))
						program.Resize(fyne.NewSize(400, 200))
						confirm_action.SetText("Pin")
						confirm_action.OnTapped = func() {
							req, err := http.NewRequest("PUT", fmt.Sprintf("https://discord.com/api/v10/channels/%s/pins/%s", chn_id.Text, msg_id.Text), nil)
							if err != nil {
								dialog.ShowError(err, program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, program)
							} else if res.StatusCode != 204 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, program)
							} else {
								dialog.ShowInformation("Success", "The message has been successfully pinned!", program)
							}
						}
						break
					}
				case "Create a channel":
					{
						program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(guild_id, actions, chn_type, chn_name, confirm_action)))
						program.Resize(fyne.NewSize(400, 240))
						confirm_action.SetText("Create")
						confirm_action.OnTapped = func() {
							body := map[string]interface{}{
								"name": chn_name.Text,
								"type": 0,
							}
							json, _ := j.Marshal(body)
							req, err := http.NewRequest("POST", fmt.Sprintf("https://discord.com/api/v10/guilds/%s/channels", guild_id.Text), b.NewBuffer(json))
							if err != nil {
								dialog.ShowError(err, program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							req.Header.Add("Content-Type", "application/json")
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, program)
							} else if res.StatusCode != 201 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, program)
							} else {
								dialog.ShowInformation("Success", "The channel has been successfully created!", program)
							}
						}
						chn_type.OnChanged = func(s string) {
							var choice int
							switch chn_type.Selected {
							case "Text":
								{
									choice = 0
									break
								}
							case "Voice":
								{
									choice = 2
									break
								}
							case "Announcement":
								{
									choice = 5
									break
								}
							case "Stage":
								{
									choice = 13
									break
								}
							case "Forum":
								{
									choice = 15
									break
								}
							case "Media":
								{
									choice = 16
									break
								}
							}
							confirm_action.OnTapped = func() {
								body := map[string]interface{}{
									"name": chn_name.Text,
									"type": choice,
								}
								json, _ := j.Marshal(body)
								req, err := http.NewRequest("POST", fmt.Sprintf("https://discord.com/api/v10/guilds/%s/channels", guild_id.Text), b.NewBuffer(json))
								if err != nil {
									dialog.ShowError(err, program)
								}
								req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
								req.Header.Add("Content-Type", "application/json")
								c := &http.Client{}
								res, err := c.Do(req)
								if err != nil {
									dialog.ShowError(err, program)
								} else if res.StatusCode != 201 {
									var body struct {
										Message string
									}
									bytes, _ := io.ReadAll(res.Body)
									j.Unmarshal(bytes, &body)
									dialog.ShowInformation("Error", body.Message, program)
								} else {
									dialog.ShowInformation("Success", "The channel has been successfully created!", program)
								}
							}
						}
						break
					}
				case "Edit a channel":
					{
						program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(chn_id, actions, chn_name, confirm_action)))
						program.Resize(fyne.NewSize(400, 200))
						confirm_action.SetText("Edit")
						confirm_action.OnTapped = func() {
							body := map[string]interface{}{
								"name": chn_name.Text,
							}
							json, _ := j.Marshal(body)
							req, err := http.NewRequest("PATCH", fmt.Sprintf("https://discord.com/api/v10/channels/%s", chn_id.Text), b.NewBuffer(json))
							if err != nil {
								dialog.ShowError(err, program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							req.Header.Add("Content-Type", "application/json")
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, program)
							} else if res.StatusCode != 200 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, program)
							} else {
								dialog.ShowInformation("Success", "The channel has been successfully edited!", program)
							}
						}
						break
					}
				case "Create a thread":
					{
						program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(chn_id, actions, msg_id, thread_name, confirm_action)))
						program.Resize(fyne.NewSize(400, 220))
						confirm_action.SetText("Create")
						confirm_action.OnTapped = func() {
							body := map[string]interface{}{
								"name": thread_name.Text,
							}
							json, _ := j.Marshal(body)
							req, err := http.NewRequest("POST", fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages/%s/threads", chn_id.Text, msg_id.Text), b.NewBuffer(json))
							if err != nil {
								dialog.ShowError(err, program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							req.Header.Add("Content-Type", "application/json")
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, program)
							} else if res.StatusCode != 201 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, program)
							} else {
								dialog.ShowInformation("Success", "The thread has been successfully created!", program)
							}
						}
						break
					}
				case "Delete a channel":
					{
						program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(chn_id, actions, confirm_action)))
						program.Resize(fyne.NewSize(400, 150))
						confirm_action.SetText("Delete")
						confirm_action.OnTapped = func() {
							req, err := http.NewRequest("DELETE", fmt.Sprintf("https://discord.com/api/v10/channels/%s", chn_id.Text), nil)
							if err != nil {
								dialog.ShowError(err, program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, program)
							} else if res.StatusCode != 200 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, program)
							} else {
								dialog.ShowInformation("Success", "The channel has been successfully deleted!", program)
							}
						}
					}
				case "Delete a message":
					{
						program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(chn_id, actions, msg_id, confirm_action)))
						program.Resize(fyne.NewSize(400, 200))
						confirm_action.SetText("Delete")
						confirm_action.OnTapped = func() {
							body := map[string]interface{}{
								"content": msg.Text,
							}
							json, _ := j.Marshal(body)
							req, err := http.NewRequest("DELETE", fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages/%s", chn_id.Text, msg_id.Text), b.NewBuffer(json))
							if err != nil {
								dialog.ShowError(err, program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, program)
							} else if res.StatusCode != 204 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, program)
							} else {
								dialog.ShowInformation("Success", "The message has been successfully deleted!", program)
							}
						}
						break
					}
				case "Unpin a message":
					{
						program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(chn_id, actions, msg_id, confirm_action)))
						program.Resize(fyne.NewSize(400, 200))
						confirm_action.SetText("Unpin")
						confirm_action.OnTapped = func() {
							req, err := http.NewRequest("DELETE", fmt.Sprintf("https://discord.com/api/v10/channels/%s/pins/%s", chn_id.Text, msg_id.Text), nil)
							if err != nil {
								dialog.ShowError(err, program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, program)
							} else if res.StatusCode != 204 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, program)
							} else {
								dialog.ShowInformation("Success", "The message has been successfully unpinned!", program)
							}
						}
						break
					}
				case "Kick a user":
					{
						program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(guild_id, actions, usr_id, confirm_action)))
						program.Resize(fyne.NewSize(400, 200))
						confirm_action.SetText("Kick")
						confirm_action.OnTapped = func() {
							req, err := http.NewRequest("DELETE", fmt.Sprintf("https://discord.com/api/v10/guilds/%s/members/%s", guild_id.Text, usr_id.Text), nil)
							if err != nil {
								dialog.ShowError(err, program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, program)
							} else if res.StatusCode != 204 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, program)
							} else {
								dialog.ShowInformation("Success", "The user has been successfully kicked!", program)
							}
						}
						break
					}
				case "Ban a user":
					{
						program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(guild_id, actions, usr_id, confirm_action)))
						program.Resize(fyne.NewSize(400, 200))
						confirm_action.SetText("Ban")
						confirm_action.OnTapped = func() {
							req, err := http.NewRequest("PUT", fmt.Sprintf("https://discord.com/api/v10/guilds/%s/bans/%s", guild_id.Text, usr_id.Text), nil)
							if err != nil {
								dialog.ShowError(err, program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, program)
							} else if res.StatusCode != 204 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, program)
							} else {
								dialog.ShowInformation("Success", "The user has been successfully banned!", program)
							}
						}
						break
					}
				case "Unban a user":
					{
						program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(guild_id, actions, usr_id, confirm_action)))
						program.Resize(fyne.NewSize(400, 200))
						confirm_action.SetText("Unban")
						confirm_action.OnTapped = func() {
							req, err := http.NewRequest("DELETE", fmt.Sprintf("https://discord.com/api/v10/guilds/%s/bans/%s", guild_id.Text, usr_id.Text), nil)
							if err != nil {
								dialog.ShowError(err, program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, program)
							} else if res.StatusCode != 204 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, program)
							} else {
								dialog.ShowInformation("Success", "The user has been successfully unbanned!", program)
							}
						}
						break
					}
				case "Create a role":
					{
						program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(guild_id, actions, role_name, confirm_action)))
						program.Resize(fyne.NewSize(400, 200))
						confirm_action.SetText("Create")
						confirm_action.OnTapped = func() {
							body := map[string]interface{}{
								"name": role_name.Text,
							}
							json, _ := j.Marshal(body)
							req, err := http.NewRequest("POST", fmt.Sprintf("https://discord.com/api/v10/guilds/%s/roles", guild_id.Text), b.NewBuffer(json))
							if err != nil {
								dialog.ShowError(err, program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							req.Header.Add("Content-Type", "application/json")
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, program)
							} else if res.StatusCode != 200 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, program)
							} else {
								dialog.ShowInformation("Success", "The role has been successfully created!", program)
							}
						}
						break
					}
				case "Edit a role":
					{
						program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(guild_id, actions, role_id, role_name, confirm_action)))
						program.Resize(fyne.NewSize(400, 240))
						confirm_action.SetText("Edit")
						confirm_action.OnTapped = func() {
							body := map[string]interface{}{
								"name": role_name.Text,
							}
							json, _ := j.Marshal(body)
							req, err := http.NewRequest("PATCH", fmt.Sprintf("https://discord.com/api/v10/guilds/%s/roles/%s", guild_id.Text, role_id.Text), b.NewBuffer(json))
							if err != nil {
								dialog.ShowError(err, program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							req.Header.Add("Content-Type", "application/json")
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, program)
							} else if res.StatusCode != 200 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, program)
							} else {
								dialog.ShowInformation("Success", "The role has been successfully edited!", program)
							}
						}
						break
					}
				case "Delete a role":
					{
						program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(guild_id, actions, role_id, confirm_action)))
						program.Resize(fyne.NewSize(400, 200))
						confirm_action.SetText("Delete")
						confirm_action.OnTapped = func() {
							req, err := http.NewRequest("DELETE", fmt.Sprintf("https://discord.com/api/v10/guilds/%s/roles/%s", guild_id.Text, role_id.Text), nil)
							if err != nil {
								dialog.ShowError(err, program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, program)
							} else if res.StatusCode != 204 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, program)
							} else {
								dialog.ShowInformation("Success", "The role has been successfully deleted!", program)
							}
						}
						break
					}
				case "Add a role to a member":
					{
						program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(guild_id, actions, usr_id, role_id, confirm_action)))
						program.Resize(fyne.NewSize(400, 240))
						confirm_action.SetText("Add")
						confirm_action.OnTapped = func() {
							req, err := http.NewRequest("PUT", fmt.Sprintf("https://discord.com/api/v10/guilds/%s/members/%s/roles/%s", guild_id.Text, usr_id.Text, role_id.Text), nil)
							if err != nil {
								dialog.ShowError(err, program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, program)
							} else if res.StatusCode != 204 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, program)
							} else {
								dialog.ShowInformation("Success", "The role has been successfully added to the provided member!", program)
							}
						}
						break
					}
				case "Remove a role from a member":
					{
						program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(guild_id, actions, usr_id, role_id, confirm_action)))
						program.Resize(fyne.NewSize(400, 240))
						confirm_action.SetText("Remove")
						confirm_action.OnTapped = func() {
							req, err := http.NewRequest("DELETE", fmt.Sprintf("https://discord.com/api/v10/guilds/%s/members/%s/roles/%s", guild_id.Text, usr_id.Text, role_id.Text), nil)
							if err != nil {
								dialog.ShowError(err, program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, program)
							} else if res.StatusCode != 204 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, program)
							} else {
								dialog.ShowInformation("Success", "The role has been successfully removed from the provided member!", program)
							}
						}
						break
					}
				}
			}
			actions.SetSelected("Write a message")
			program.SetContent(container.NewBorder(navbar_edit, nil, nil, nil, container.NewVBox(chn_id, actions, msg_box, confirm_action)))
			program.Show()
		}
	})
	custom_activity := widget.NewCheck("Custom Activity", nil)
	activity_box := container.NewGridWithRows(1)
	box := container.NewBorder(nil, nil, container.NewHBox(act_type, act_status), nil, act_description)
	custom_activity.OnChanged = func(b bool) {
		if b {
			activity_box.Add(box)
			activity_box.Refresh()
			act_description.OnChanged = func(s string) {
				if len(s) > 32 {
					validate.Disable()
				} else {
					validate.Enable()
				}
			}
		} else {
			activity_box.RemoveAll()
			activity_box.Refresh()
			act_type.SetSelected("Watching")
			act_description.Text = "http_requests"
		}
	}
	login.SetContent(container.NewVBox(container.NewCenter(custom_activity), tkn, validate, activity_box))
	login.Show()
	a.Run()
}
