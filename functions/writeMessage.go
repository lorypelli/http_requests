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

func WriteMessage(navbar_edit, msg_box *fyne.Container, tkn, msg, chn_id *widget.Entry, actions *widget.Select, confirm_action *widget.Button) {
	windows.Program.SetContent(container.NewBorder(navbar_edit, nil, nil, nil, container.NewVBox(chn_id, actions, msg_box, confirm_action)))
	windows.Program.Resize(fyne.NewSize(400, 240))
	confirm_action.SetText("Send")
	confirm_action.OnTapped = func() {
		internalWriteMessage(tkn, msg, chn_id)
	}
}

func internalWriteMessage(tkn, msg, chn_id *widget.Entry) {
	body := map[string]interface{}{
		"content": msg.Text,
	}
	json, _ := j.Marshal(body)
	req, err := http.NewRequest("POST", fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages", chn_id.Text), b.NewBuffer(json))
	if err != nil {
		dialog.ShowError(err, windows.Program)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
	req.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		dialog.ShowError(err, windows.Program)
	} else if res.StatusCode != 200 {
		ShowError(res.Body)
	} else {
		dialog.ShowInformation("Success", "The message has been successfully sent!", windows.Program)
	}
}
