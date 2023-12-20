package main

import (
	b "bytes"
	j "encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	login := a.NewWindow("Login")
	program := a.NewWindow("http_requests")
	bot := a.NewWindow("bot_info")
	req, _ := http.NewRequest("GET", "https://raw.githubusercontent.com/lorypelli/http_requests/main/icon.png", nil)
	c := &http.Client{}
	res, _ := c.Do(req)
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
	login.Resize(fyne.NewSize(670, 170))
	program.Resize(fyne.NewSize(400, 240))
	bot.Resize(fyne.NewSize(400, 200))
	login.SetFixedSize(true)
	program.SetFixedSize(true)
	bot.SetFixedSize(true)
	login.CenterOnScreen()
	program.CenterOnScreen()
	bot.CenterOnScreen()
	login.SetIcon(icon)
	program.SetIcon(icon)
	bot.SetIcon(icon)
	show := false
	program.SetCloseIntercept(func() {
		if !show {
			dialog.ShowConfirm("Logout", "Are you sure you want to logout?", func(b bool) {
				if b {
					login.Show()
					program.Hide()
				} else {
					show = false
				}
			}, program)
			show = true
		}
	})
	bot.SetCloseIntercept(func() {
		bot.Hide()
	})
	tkn := widget.NewPasswordEntry()
	tkn.SetPlaceHolder("Insert bot token")
	login.SetContent(container.NewVBox(layout.NewSpacer(), tkn, widget.NewButton("Validate", func() {
		req, err := http.NewRequest("POST", "https://discord.com/api/v10/auth/login", nil)
		if err != nil {
			dialog.ShowError(err, login)
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
		c := &http.Client{}
		res, err := c.Do(req)
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
			navbar := container.NewHBox(widget.NewButton("Bot Info", func() {
				req, err := http.NewRequest("GET", "https://discord.com/api/v10/users/@me", nil)
				if err != nil {
					dialog.ShowError(err, login)
				}
				req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
				c := &http.Client{}
				res, err := c.Do(req)
				if err != nil {
					dialog.ShowError(err, login)
				}
				var bots struct {
					Id       string
					Username string
				}
				bytes, _ := io.ReadAll(res.Body)
				j.Unmarshal(bytes, &bots)
				req, err = http.NewRequest("GET", "https://discord.com/api/v10/users/@me/guilds", nil)
				if err != nil {
					dialog.ShowError(err, login)
				}
				req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
				c = &http.Client{}
				res, err = c.Do(req)
				if err != nil {
					dialog.ShowError(err, login)
				}
				var guilds []struct{}
				bytes, _ = io.ReadAll(res.Body)
				j.Unmarshal(bytes, &guilds)
				bot.SetContent(container.NewCenter(container.NewVBox(widget.NewLabel(fmt.Sprintf("Username: %s", bots.Username)), widget.NewLabel(fmt.Sprintf("ID: %s", bots.Id)), widget.NewLabel(fmt.Sprintf("Server Count: %d", len(guilds))))))
				bot.Show()
			}), layout.NewSpacer(), widget.NewButton("Logout", func() {
				dialog.ShowConfirm("Logout", "Are you sure you want to logout?", func(b bool) {
					if b {
						login.Show()
						program.Hide()
					} else {
						show = false
					}
				}, program)
				show = true
			}))
			chn_id := widget.NewEntry()
			chn_id.SetPlaceHolder("Insert channel ID")
			msg := widget.NewMultiLineEntry()
			msg.SetPlaceHolder("Insert message")
			msg_id := widget.NewEntry()
			msg_id.SetPlaceHolder("Insert message ID")
			guild_id := widget.NewEntry()
			guild_id.SetPlaceHolder("Insert guild ID")
			chn_type := widget.NewSelect([]string{"Text", "Voice", "Stage", "Announcement", "Forum", "Media"}, func(s string) {})
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
			actions.OnChanged = func(s string) {
				switch s {
				case "Write a message":
					{
						program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(chn_id, actions, msg, confirm_action)))
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
						program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(chn_id, actions, msg_id, msg, confirm_action)))
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
			program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(chn_id, actions, msg, confirm_action)))
			program.Show()
		}
	}), layout.NewSpacer()))
	login.Show()
	a.Run()
}
