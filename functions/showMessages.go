package actions

import (
	j "encoding/json"
	"fmt"
	"http_requests/windows"
	"io"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func ShowMessages(tkn, chn_id *widget.Entry) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages?limit=100", chn_id.Text), nil)
	if err != nil {
		dialog.ShowError(err, windows.Program)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
	res, err := http.DefaultClient.Do(req)
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
}