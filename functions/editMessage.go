package actions

import (
	b "bytes"
	j "encoding/json"
	"fmt"
	"http_requests/windows"
	"io"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func EditMessage(navbar_edit, msg_box *fyne.Container, tkn, msg, chn_id, msg_id *widget.Entry, actions *widget.Select, confirm_action *widget.Button) {
	windows.Program.Resize(fyne.NewSize(400, 270))
	windows.Program.SetContent(container.NewBorder(navbar_edit, nil, nil, nil, container.NewVBox(chn_id, actions, msg_id, msg_box, confirm_action)))
	confirm_action.SetText("Edit")
	confirm_action.OnTapped = func() {
		internalEditMessage(tkn, msg, chn_id, msg_id)
	}
}

func internalEditMessage(tkn, msg, chn_id, msg_id *widget.Entry) {
	body := map[string]interface{}{
		"content": msg.Text,
	}
	json, _ := j.Marshal(body)
	req, err := http.NewRequest("PATCH", fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages/%s", chn_id.Text, msg_id.Text), b.NewBuffer(json))
	if err != nil {
		dialog.ShowError(err, windows.Program)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
	req.Header.Add("Content-Type", "application/json")
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
		dialog.ShowInformation("Success", "The message has been successfully edited!", windows.Program)
	}
}