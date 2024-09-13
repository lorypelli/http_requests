package actions

import (
	b "bytes"
	j "encoding/json"
	"fmt"
	"http_requests/windows"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func CreateChannel(navbar *fyne.Container, chn_name, tkn, guild_id *widget.Entry, actions, chn_type *widget.Select, confirm_action *widget.Button) {
	windows.Program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(guild_id, actions, chn_type, chn_name, confirm_action)))
	windows.Program.Resize(fyne.NewSize(400, 240))
	confirm_action.SetText("Create")
	var choice uint8
	confirm_action.OnTapped = func() {
		internalCreateChannel(chn_name, tkn, guild_id, choice)
	}
	chn_type.OnChanged = func(s string) {
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
			internalCreateChannel(chn_name, tkn, guild_id, choice)
		}
	}
}

func internalCreateChannel(chn_name, tkn, guild_id *widget.Entry, choice uint8) {
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
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		dialog.ShowError(err, windows.Program)
	} else if res.StatusCode != 201 {
		ShowError(res.Body)
	} else {
		dialog.ShowInformation("Success", "The channel has been successfully created!", windows.Program)
	}
}
