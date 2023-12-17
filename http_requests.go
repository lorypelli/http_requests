package main

import (
	b "bytes"
	j "encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

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
	fileName := "icon.png"
	file, err := os.Open(fileName)
	if err == nil {
		stats, _ := os.Stat(fileName)
		size := stats.Size()
		fileByte := make([]byte, size)
		fileSlice := fileByte[:]
		file.Read(fileSlice)
		login.SetIcon(fyne.NewStaticResource(fileName, fileByte))
		program.SetIcon(fyne.NewStaticResource(fileName, fileByte))
	}
	login.Resize(fyne.NewSize(640, 170))
	program.Resize(fyne.NewSize(400, 240))
	login.SetFixedSize(true)
	program.SetFixedSize(true)
	login.CenterOnScreen()
	program.CenterOnScreen()
	tkn_textbox := widget.NewEntry()
	tkn_textbox.SetPlaceHolder("Insert bot token")
	login.SetContent(container.NewVBox(tkn_textbox, widget.NewButton("Validate", func() {
		req, err := http.NewRequest("POST", "https://discord.com/api/v10/auth/login", nil)
		if err != nil {
			dialog.ShowError(err, login)
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn_textbox.Text))
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
			req, err := http.NewRequest("GET", "https://discord.com/api/v10/users/@me", nil)
			if err != nil {
				dialog.ShowError(err, login)
			}
			req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn_textbox.Text))
			c := &http.Client{}
			res, err := c.Do(req)
			if err != nil {
				dialog.ShowError(err, login)
			}
			var body struct {
				Id       string
				Username string
			}
			bytes, _ := io.ReadAll(res.Body)
			j.Unmarshal(bytes, &body)
			botId := body.Id
			botUsername := body.Username
			login.Hide()
			chn_id_textbox := widget.NewEntry()
			chn_id_textbox.SetPlaceHolder("Insert channel ID")
			msg_textbox := widget.NewMultiLineEntry()
			msg_textbox.SetPlaceHolder("Insert message")
			confirm_action := widget.NewButton("Send", func() {
				body := map[string]interface{}{
					"content": msg_textbox.Text,
				}
				json, _ := j.Marshal(body)
				req, err := http.NewRequest("POST", fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages", chn_id_textbox.Text), b.NewBuffer(json))
				if err != nil {
					dialog.ShowError(err, program)
				}
				req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn_textbox.Text))
				req.Header.Set("Content-Type", "application/json")
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
						program.SetContent(container.NewBorder(container.NewHBox(widget.NewLabel(botId), layout.NewSpacer(), widget.NewButton("Logout", func() {
							dialog.ShowConfirm("Logout", "Are you sure you want to logout?", func(b bool) {
								if b {
									login.Show()
									program.Hide()
								}
							}, program)
						}), layout.NewSpacer(), widget.NewLabel(botUsername)), nil, nil, nil, container.NewVBox(chn_id_textbox, actions, msg_textbox, confirm_action)))
						program.Resize(fyne.NewSize(400, 240))
						confirm_action.SetText("Send")
						confirm_action.OnTapped = func() {
							body := map[string]interface{}{
								"content": msg_textbox.Text,
							}
							json, _ := j.Marshal(body)
							req, err := http.NewRequest("POST", fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages", chn_id_textbox.Text), b.NewBuffer(json))
							if err != nil {
								dialog.ShowError(err, program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn_textbox.Text))
							req.Header.Set("Content-Type", "application/json")
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
						msg_id_textbox := widget.NewEntry()
						msg_id_textbox.SetPlaceHolder("Insert message ID")
						program.Resize(fyne.NewSize(400, 270))
						program.SetContent(container.NewBorder(container.NewHBox(widget.NewLabel(botId), layout.NewSpacer(), widget.NewButton("Logout", func() {
							dialog.ShowConfirm("Logout", "Are you sure you want to logout?", func(b bool) {
								if b {
									login.Show()
									program.Hide()
								}
							}, program)
						}), layout.NewSpacer(), widget.NewLabel(botUsername)), nil, nil, nil, container.NewVBox(chn_id_textbox, actions, msg_id_textbox, msg_textbox, confirm_action)))
						confirm_action.SetText("Edit")
						confirm_action.OnTapped = func() {
							body := map[string]interface{}{
								"content": msg_textbox.Text,
							}
							json, _ := j.Marshal(body)
							req, err := http.NewRequest("PATCH", fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages/%s", chn_id_textbox.Text, msg_id_textbox.Text), b.NewBuffer(json))
							if err != nil {
								dialog.ShowError(err, program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn_textbox.Text))
							req.Header.Set("Content-Type", "application/json")
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
						msg_id_textbox := widget.NewEntry()
						msg_id_textbox.SetPlaceHolder("Insert message ID")
						program.SetContent(container.NewBorder(container.NewHBox(widget.NewLabel(botId), layout.NewSpacer(), widget.NewButton("Logout", func() {
							dialog.ShowConfirm("Logout", "Are you sure you want to logout?", func(b bool) {
								if b {
									login.Show()
									program.Hide()
								}
							}, program)
						}), layout.NewSpacer(), widget.NewLabel(botUsername)), nil, nil, nil, container.NewVBox(chn_id_textbox, actions, msg_id_textbox, confirm_action)))
						program.Resize(fyne.NewSize(400, 200))
						confirm_action.SetText("Pin")
						confirm_action.OnTapped = func() {
							req, err := http.NewRequest("PUT", fmt.Sprintf("https://discord.com/api/v10/channels/%s/pins/%s", chn_id_textbox.Text, msg_id_textbox.Text), nil)
							if err != nil {
								dialog.ShowError(err, program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn_textbox.Text))
							req.Header.Set("Content-Type", "application/json")
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
						guild_id_textbox := widget.NewEntry()
						guild_id_textbox.SetPlaceHolder("Insert guild ID")
						chn_type := widget.NewSelect([]string{"Text", "Voice", "Stage", "Announcement", "Forum", "Media"}, func(s string) {})
						chn_type.SetSelected("Text")
						chn_name := widget.NewEntry()
						chn_name.SetPlaceHolder("Insert channel name")
						program.SetContent(container.NewBorder(container.NewHBox(widget.NewLabel(botId), layout.NewSpacer(), widget.NewButton("Logout", func() {
							dialog.ShowConfirm("Logout", "Are you sure you want to logout?", func(b bool) {
								if b {
									login.Show()
									program.Hide()
								}
							}, program)
						}), layout.NewSpacer(), widget.NewLabel(botUsername)), nil, nil, nil, container.NewVBox(guild_id_textbox, actions, chn_type, chn_name, confirm_action)))
						program.Resize(fyne.NewSize(400, 240))
						confirm_action.SetText("Create")
						choice := 0
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
							req, err := http.NewRequest("POST", fmt.Sprintf("https://discord.com/api/v10/guilds/%s/channels", guild_id_textbox.Text), b.NewBuffer(json))
							if err != nil {
								dialog.ShowError(err, program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn_textbox.Text))
							req.Header.Set("Content-Type", "application/json")
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
								dialog.ShowInformation("Success", "The message has been successfully unpinned!", program)
							}
						}
						break
					}
				case "Edit a channel":
					{
						confirm_action.SetText("Edit")
						break
					}
				case "Create a thread":
					{
						confirm_action.SetText("Create")
						break
					}
				case "Delete a channel":
					{
						confirm_action.SetText("Delete")
						break
					}
				case "Delete a message":
					{
						confirm_action.SetText("Delete")
						break
					}
				case "Unpin a message":
					{
						msg_id_textbox := widget.NewEntry()
						msg_id_textbox.SetPlaceHolder("Insert message ID")
						program.SetContent(container.NewBorder(container.NewHBox(widget.NewLabel(botId), layout.NewSpacer(), widget.NewButton("Logout", func() {
							dialog.ShowConfirm("Logout", "Are you sure you want to logout?", func(b bool) {
								if b {
									login.Show()
									program.Hide()
								}
							}, program)
						}), layout.NewSpacer(), widget.NewLabel(botUsername)), nil, nil, nil, container.NewVBox(chn_id_textbox, actions, msg_id_textbox, confirm_action)))
						program.Resize(fyne.NewSize(400, 200))
						confirm_action.SetText("Unpin")
						confirm_action.OnTapped = func() {
							req, err := http.NewRequest("DELETE", fmt.Sprintf("https://discord.com/api/v10/channels/%s/pins/%s", chn_id_textbox.Text, msg_id_textbox.Text), nil)
							if err != nil {
								dialog.ShowError(err, program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn_textbox.Text))
							req.Header.Set("Content-Type", "application/json")
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
						confirm_action.SetText("Kick")
						break
					}
				case "Ban a user":
					{
						confirm_action.SetText("Ban")
						break
					}
				case "Unban a user":
					{
						confirm_action.SetText("Unban")
						break
					}
				case "Create a role":
					{
						confirm_action.SetText("Create")
						break
					}
				case "Edit a role":
					{
						confirm_action.SetText("Edit")
						break
					}
				case "Delete a role":
					{
						confirm_action.SetText("Delete")
						break
					}
				case "Add a role to a member":
					{
						confirm_action.SetText("Add")
						break
					}
				case "Remove a role from a member":
					{
						confirm_action.SetText("Remove")
						break
					}
				}
			}
			actions.SetSelected("Write a message")
			program.SetContent(container.NewBorder(container.NewHBox(widget.NewLabel(botId), layout.NewSpacer(), widget.NewButton("Logout", func() {
				dialog.ShowConfirm("Logout", "Are you sure you want to logout?", func(b bool) {
					if b {
						login.Show()
						program.Hide()
					}
				}, program)
			}), layout.NewSpacer(), widget.NewLabel(botUsername)), nil, nil, nil, container.NewVBox(chn_id_textbox, actions, msg_textbox, confirm_action)))
			program.Show()
		}
	})))
	login.Show()
	a.Run()
}
