package main

import (
	b "bytes"
	j "encoding/json"
	"fmt"
	act "http_requests/functions"
	"http_requests/globals"
	"io"
	"net/http"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	windows.Setup()
	act_type := widget.NewSelect([]string{"Game", "Listening", "Watching", "Competing"}, nil)
	act_type.SetSelected("Watching")
	act_description := widget.NewEntry()
	act_description.Text = "http_requests"
	act_description.SetPlaceHolder("Insert detailed activity")
	act_status := widget.NewSelect([]string{"online", "idle", "dnd"}, nil)
	act_status.SetSelected("dnd")
	show := false
	tkn := widget.NewPasswordEntry()
	tkn.SetPlaceHolder("Insert bot token")
	validate := widget.NewButton("Validate", func() {
		req, err := http.NewRequest("POST", "https://discord.com/api/v10/auth/windows.Login", nil)
		if err != nil {
			dialog.ShowError(err, windows.Login)
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
		c := &http.Client{}
		res, err := c.Do(req)
		timer := time.NewTicker(120 * time.Second)
		if err != nil {
			dialog.ShowError(err, windows.Login)
		} else if res.StatusCode != 200 {
			var body struct {
				Message string
			}
			bytes, _ := io.ReadAll(res.Body)
			j.Unmarshal(bytes, &body)
			dialog.ShowInformation("Error", body.Message, windows.Login)
		} else {
			windows.Login.Hide()
			logout := func(b bool) {
				if b {
					stop := make(chan struct{})
					go act.Connect(tkn.Text, act_description.Text, act_type.Selected, act_status.Selected, timer, stop)
					close(stop)
					windows.Login.Show()
					windows.Program.Hide()
					windows.Msg_list.Hide()
					windows.Bot.Hide()
				} else {
					show = false
				}
			}
			windows.Program.SetCloseIntercept(func() {
				if !show {
					dialog.ShowConfirm("Logout", "Are you sure you want to logout?", logout, windows.Program)
					show = true
				}
			})
			go act.Connect(tkn.Text, act_description.Text, act_type.Selected, act_status.Selected, timer, nil)
			bot_info := widget.NewButton("Bot Info", func() {
				act.BotInfo(tkn)
			})
			logout_btn := widget.NewButton("Logout", func() {
				dialog.ShowConfirm("Logout", "Are you sure you want to logout?", logout, windows.Program)
				show = true
			})
			navbar := container.NewHBox(bot_info, layout.NewSpacer(), logout_btn)
			chn_id := widget.NewEntry()
			chn_id.SetPlaceHolder("Insert channel ID")
			show_msg_list := widget.NewButton("Show Message List", func() {
				req, err := http.NewRequest("GET", fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages?limit=100", chn_id.Text), nil)
				if err != nil {
					dialog.ShowError(err, windows.Program)
				}
				req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
				c := &http.Client{}
				res, err := c.Do(req)
				if err != nil {
					dialog.ShowError(err, windows.Program)
				} else if res.StatusCode != 200 {
					var body struct {
						Message string
					}
					bytes, _ := io.ReadAll(res.Body)
					j.Unmarshal(bytes, &body)
					dialog.ShowInformation("Error", body.Message, windows.Program)
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
						content.Wrapping = fyne.TextWrapWord
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
						msgs_container.Add(container.NewBorder(nil, nil, widget.NewLabel(fmt.Sprintf("%s :", msgs[i].Author.Username)), nil, content))
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
					split_container := container.NewHSplit(container.NewScroll(msgs_container), container.NewScroll(users_container))
					split_container.SetOffset(0.8)
					windows.Msg_list.SetContent(split_container)
					windows.Msg_list.Show()
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
				act.Send(msg, chn_id, tkn, windows.Program)
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
						act.WriteMessage(navbar_edit, msg_box, msg, chn_id, tkn, actions, confirm_action)
						break
					}
				case "Edit a message":
					{
						act.EditMessage(navbar_edit, msg_box, tkn, msg, chn_id, msg_id, actions, confirm_action)
						break
					}
				case "Pin a message":
					{
						act.PinMessage(navbar, chn_id, tkn, msg_id, actions, confirm_action)
						break
					}
				case "Create a channel":
					{
						windows.Program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(guild_id, actions, chn_type, chn_name, confirm_action)))
						windows.Program.Resize(fyne.NewSize(400, 240))
						confirm_action.SetText("Create")
						confirm_action.OnTapped = func() {
							body := map[string]interface{}{
								"name": chn_name.Text,
								"type": 0,
							}
							json, _ := j.Marshal(body)
							req, err := http.NewRequest("POST", fmt.Sprintf("https://discord.com/api/v10/guilds/%s/channels", guild_id.Text), b.NewBuffer(json))
							if err != nil {
								dialog.ShowError(err, windows.Program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							req.Header.Add("Content-Type", "application/json")
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, windows.Program)
							} else if res.StatusCode != 201 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, windows.Program)
							} else {
								dialog.ShowInformation("Success", "The channel has been successfully created!", windows.Program)
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
									dialog.ShowError(err, windows.Program)
								}
								req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
								req.Header.Add("Content-Type", "application/json")
								c := &http.Client{}
								res, err := c.Do(req)
								if err != nil {
									dialog.ShowError(err, windows.Program)
								} else if res.StatusCode != 201 {
									var body struct {
										Message string
									}
									bytes, _ := io.ReadAll(res.Body)
									j.Unmarshal(bytes, &body)
									dialog.ShowInformation("Error", body.Message, windows.Program)
								} else {
									dialog.ShowInformation("Success", "The channel has been successfully created!", windows.Program)
								}
							}
						}
						break
					}
				case "Edit a channel":
					{
						windows.Program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(chn_id, actions, chn_name, confirm_action)))
						windows.Program.Resize(fyne.NewSize(400, 200))
						confirm_action.SetText("Edit")
						confirm_action.OnTapped = func() {
							body := map[string]interface{}{
								"name": chn_name.Text,
							}
							json, _ := j.Marshal(body)
							req, err := http.NewRequest("PATCH", fmt.Sprintf("https://discord.com/api/v10/channels/%s", chn_id.Text), b.NewBuffer(json))
							if err != nil {
								dialog.ShowError(err, windows.Program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							req.Header.Add("Content-Type", "application/json")
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, windows.Program)
							} else if res.StatusCode != 200 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, windows.Program)
							} else {
								dialog.ShowInformation("Success", "The channel has been successfully edited!", windows.Program)
							}
						}
						break
					}
				case "Create a thread":
					{
						windows.Program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(chn_id, actions, msg_id, thread_name, confirm_action)))
						windows.Program.Resize(fyne.NewSize(400, 220))
						confirm_action.SetText("Create")
						confirm_action.OnTapped = func() {
							body := map[string]interface{}{
								"name": thread_name.Text,
							}
							json, _ := j.Marshal(body)
							req, err := http.NewRequest("POST", fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages/%s/threads", chn_id.Text, msg_id.Text), b.NewBuffer(json))
							if err != nil {
								dialog.ShowError(err, windows.Program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							req.Header.Add("Content-Type", "application/json")
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, windows.Program)
							} else if res.StatusCode != 201 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, windows.Program)
							} else {
								dialog.ShowInformation("Success", "The thread has been successfully created!", windows.Program)
							}
						}
						break
					}
				case "Delete a channel":
					{
						windows.Program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(chn_id, actions, confirm_action)))
						windows.Program.Resize(fyne.NewSize(400, 150))
						confirm_action.SetText("Delete")
						confirm_action.OnTapped = func() {
							req, err := http.NewRequest("DELETE", fmt.Sprintf("https://discord.com/api/v10/channels/%s", chn_id.Text), nil)
							if err != nil {
								dialog.ShowError(err, windows.Program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, windows.Program)
							} else if res.StatusCode != 200 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, windows.Program)
							} else {
								dialog.ShowInformation("Success", "The channel has been successfully deleted!", windows.Program)
							}
						}
					}
				case "Delete a message":
					{
						windows.Program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(chn_id, actions, msg_id, confirm_action)))
						windows.Program.Resize(fyne.NewSize(400, 200))
						confirm_action.SetText("Delete")
						confirm_action.OnTapped = func() {
							body := map[string]interface{}{
								"content": msg.Text,
							}
							json, _ := j.Marshal(body)
							req, err := http.NewRequest("DELETE", fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages/%s", chn_id.Text, msg_id.Text), b.NewBuffer(json))
							if err != nil {
								dialog.ShowError(err, windows.Program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, windows.Program)
							} else if res.StatusCode != 204 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, windows.Program)
							} else {
								dialog.ShowInformation("Success", "The message has been successfully deleted!", windows.Program)
							}
						}
						break
					}
				case "Unpin a message":
					{
						windows.Program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(chn_id, actions, msg_id, confirm_action)))
						windows.Program.Resize(fyne.NewSize(400, 200))
						confirm_action.SetText("Unpin")
						confirm_action.OnTapped = func() {
							req, err := http.NewRequest("DELETE", fmt.Sprintf("https://discord.com/api/v10/channels/%s/pins/%s", chn_id.Text, msg_id.Text), nil)
							if err != nil {
								dialog.ShowError(err, windows.Program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, windows.Program)
							} else if res.StatusCode != 204 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, windows.Program)
							} else {
								dialog.ShowInformation("Success", "The message has been successfully unpinned!", windows.Program)
							}
						}
						break
					}
				case "Kick a user":
					{
						windows.Program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(guild_id, actions, usr_id, confirm_action)))
						windows.Program.Resize(fyne.NewSize(400, 200))
						confirm_action.SetText("Kick")
						confirm_action.OnTapped = func() {
							req, err := http.NewRequest("DELETE", fmt.Sprintf("https://discord.com/api/v10/guilds/%s/members/%s", guild_id.Text, usr_id.Text), nil)
							if err != nil {
								dialog.ShowError(err, windows.Program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, windows.Program)
							} else if res.StatusCode != 204 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, windows.Program)
							} else {
								dialog.ShowInformation("Success", "The user has been successfully kicked!", windows.Program)
							}
						}
						break
					}
				case "Ban a user":
					{
						windows.Program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(guild_id, actions, usr_id, confirm_action)))
						windows.Program.Resize(fyne.NewSize(400, 200))
						confirm_action.SetText("Ban")
						confirm_action.OnTapped = func() {
							req, err := http.NewRequest("PUT", fmt.Sprintf("https://discord.com/api/v10/guilds/%s/bans/%s", guild_id.Text, usr_id.Text), nil)
							if err != nil {
								dialog.ShowError(err, windows.Program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, windows.Program)
							} else if res.StatusCode != 204 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, windows.Program)
							} else {
								dialog.ShowInformation("Success", "The user has been successfully banned!", windows.Program)
							}
						}
						break
					}
				case "Unban a user":
					{
						windows.Program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(guild_id, actions, usr_id, confirm_action)))
						windows.Program.Resize(fyne.NewSize(400, 200))
						confirm_action.SetText("Unban")
						confirm_action.OnTapped = func() {
							req, err := http.NewRequest("DELETE", fmt.Sprintf("https://discord.com/api/v10/guilds/%s/bans/%s", guild_id.Text, usr_id.Text), nil)
							if err != nil {
								dialog.ShowError(err, windows.Program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, windows.Program)
							} else if res.StatusCode != 204 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, windows.Program)
							} else {
								dialog.ShowInformation("Success", "The user has been successfully unbanned!", windows.Program)
							}
						}
						break
					}
				case "Create a role":
					{
						windows.Program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(guild_id, actions, role_name, confirm_action)))
						windows.Program.Resize(fyne.NewSize(400, 200))
						confirm_action.SetText("Create")
						confirm_action.OnTapped = func() {
							body := map[string]interface{}{
								"name": role_name.Text,
							}
							json, _ := j.Marshal(body)
							req, err := http.NewRequest("POST", fmt.Sprintf("https://discord.com/api/v10/guilds/%s/roles", guild_id.Text), b.NewBuffer(json))
							if err != nil {
								dialog.ShowError(err, windows.Program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							req.Header.Add("Content-Type", "application/json")
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, windows.Program)
							} else if res.StatusCode != 200 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, windows.Program)
							} else {
								dialog.ShowInformation("Success", "The role has been successfully created!", windows.Program)
							}
						}
						break
					}
				case "Edit a role":
					{
						windows.Program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(guild_id, actions, role_id, role_name, confirm_action)))
						windows.Program.Resize(fyne.NewSize(400, 240))
						confirm_action.SetText("Edit")
						confirm_action.OnTapped = func() {
							body := map[string]interface{}{
								"name": role_name.Text,
							}
							json, _ := j.Marshal(body)
							req, err := http.NewRequest("PATCH", fmt.Sprintf("https://discord.com/api/v10/guilds/%s/roles/%s", guild_id.Text, role_id.Text), b.NewBuffer(json))
							if err != nil {
								dialog.ShowError(err, windows.Program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							req.Header.Add("Content-Type", "application/json")
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, windows.Program)
							} else if res.StatusCode != 200 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, windows.Program)
							} else {
								dialog.ShowInformation("Success", "The role has been successfully edited!", windows.Program)
							}
						}
						break
					}
				case "Delete a role":
					{
						windows.Program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(guild_id, actions, role_id, confirm_action)))
						windows.Program.Resize(fyne.NewSize(400, 200))
						confirm_action.SetText("Delete")
						confirm_action.OnTapped = func() {
							req, err := http.NewRequest("DELETE", fmt.Sprintf("https://discord.com/api/v10/guilds/%s/roles/%s", guild_id.Text, role_id.Text), nil)
							if err != nil {
								dialog.ShowError(err, windows.Program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, windows.Program)
							} else if res.StatusCode != 204 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, windows.Program)
							} else {
								dialog.ShowInformation("Success", "The role has been successfully deleted!", windows.Program)
							}
						}
						break
					}
				case "Add a role to a member":
					{
						windows.Program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(guild_id, actions, usr_id, role_id, confirm_action)))
						windows.Program.Resize(fyne.NewSize(400, 240))
						confirm_action.SetText("Add")
						confirm_action.OnTapped = func() {
							req, err := http.NewRequest("PUT", fmt.Sprintf("https://discord.com/api/v10/guilds/%s/members/%s/roles/%s", guild_id.Text, usr_id.Text, role_id.Text), nil)
							if err != nil {
								dialog.ShowError(err, windows.Program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, windows.Program)
							} else if res.StatusCode != 204 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, windows.Program)
							} else {
								dialog.ShowInformation("Success", "The role has been successfully added to the provided member!", windows.Program)
							}
						}
						break
					}
				case "Remove a role from a member":
					{
						windows.Program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(guild_id, actions, usr_id, role_id, confirm_action)))
						windows.Program.Resize(fyne.NewSize(400, 240))
						confirm_action.SetText("Remove")
						confirm_action.OnTapped = func() {
							req, err := http.NewRequest("DELETE", fmt.Sprintf("https://discord.com/api/v10/guilds/%s/members/%s/roles/%s", guild_id.Text, usr_id.Text, role_id.Text), nil)
							if err != nil {
								dialog.ShowError(err, windows.Program)
							}
							req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
							c := &http.Client{}
							res, err := c.Do(req)
							if err != nil {
								dialog.ShowError(err, windows.Program)
							} else if res.StatusCode != 204 {
								var body struct {
									Message string
								}
								bytes, _ := io.ReadAll(res.Body)
								j.Unmarshal(bytes, &body)
								dialog.ShowInformation("Error", body.Message, windows.Program)
							} else {
								dialog.ShowInformation("Success", "The role has been successfully removed from the provided member!", windows.Program)
							}
						}
						break
					}
				}
			}
			actions.SetSelected("Write a message")
			windows.Program.SetContent(container.NewBorder(navbar_edit, nil, nil, nil, container.NewVBox(chn_id, actions, msg_box, confirm_action)))
			windows.Program.Show()
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
			act_description.SetText("http_requests")
		}
	}
	windows.Login.SetContent(container.NewVBox(container.NewCenter(custom_activity), tkn, validate, activity_box))
	windows.Login.Show()
	windows.App.Run()
}
