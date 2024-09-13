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

func EditChannel(navbar *fyne.Container, chn_name, tkn, chn_id *widget.Entry, actions *widget.Select, confirm_action *widget.Button) {
	windows.Program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(chn_id, actions, chn_name, confirm_action)))
	windows.Program.Resize(fyne.NewSize(400, 200))
	confirm_action.SetText("Edit")
	confirm_action.OnTapped = func() {
		internalEditChannel(chn_name, tkn, chn_id)
	}
}

func internalEditChannel(chn_name, tkn, chn_id *widget.Entry) {
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
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		dialog.ShowError(err, windows.Program)
	} else if res.StatusCode != 200 {
		ShowError(res.Body)
	} else {
		dialog.ShowInformation("Success", "The channel has been successfully edited!", windows.Program)
	}
}
