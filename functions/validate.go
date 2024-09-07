package actions

import (
	j "encoding/json"
	"fmt"
	"http_requests/globals"
	"io"
	"net/http"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func Validate(tkn, act_description *widget.Entry, act_type, act_status *widget.Select, show bool) {
	req, err := http.NewRequest("POST", "https://discord.com/api/v10/auth/login", nil)
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
				go Connect(tkn.Text, act_description.Text, act_type.Selected, act_status.Selected, timer, stop)
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
		go Connect(tkn.Text, act_description.Text, act_type.Selected, act_status.Selected, timer, nil)
		bot_info := widget.NewButton("Bot Info", func() {
			BotInfo(tkn)
		})
		logout_btn := widget.NewButton("Logout", func() {
			dialog.ShowConfirm("Logout", "Are you sure you want to logout?", logout, windows.Program)
			show = true
		})
		//navbar := container.NewHBox(bot_info, layout.NewSpacer(), logout_btn)
		chn_id := widget.NewEntry()
		chn_id.SetPlaceHolder("Insert channel ID")
		show_msg_list := widget.NewButton("Show Message List", func() {
			showMessages(tkn, chn_id)
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
			Send(msg, chn_id, tkn, windows.Program)
		})
		actions := widget.NewSelect([]string{"Write a message", "Edit a message", "Pin a message", "Create a channel", "Edit a channel", "Create a thread", "Delete a channel", "Delete a message", "Unpin a message", "Kick a user", "Ban a user", "Unban a user", "Create a role", "Edit a role", "Delete a role", "Add a role to a member", "Remove a role from a member"}, nil)
		msg.OnChanged = func(s string) {
			count.SetText(fmt.Sprint(len(s)))
			if len(msg.Text) > 4096 {
				confirm_action.Disable()
			}
		}
		actions.SetSelected("Write a message")
		windows.Program.SetContent(container.NewBorder(navbar_edit, nil, nil, nil, container.NewVBox(chn_id, actions, msg_box, confirm_action)))
		windows.Program.Show()
	}
}