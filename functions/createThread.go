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

func CreateThread(navbar *fyne.Container, chn_id, tkn, msg_id, thread_name *widget.Entry, actions *widget.Select, confirm_action *widget.Button) {
	windows.Program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(chn_id, actions, msg_id, thread_name, confirm_action)))
	windows.Program.Resize(fyne.NewSize(400, 220))
	confirm_action.SetText("Create")
	confirm_action.OnTapped = func() {
		internalCreateThread(chn_id, tkn, msg_id, thread_name)
	}
}

func internalCreateThread(chn_id, tkn, msg_id, thread_name *widget.Entry) {
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
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		dialog.ShowError(err, windows.Program)
	} else if res.StatusCode != 201 {
		ShowError(res.Body)
	} else {
		dialog.ShowInformation("Success", "The thread has been successfully created!", windows.Program)
	}
}
