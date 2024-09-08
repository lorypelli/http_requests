package actions

import (
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

func DeleteChannel(navbar *fyne.Container, tkn, chn_id *widget.Entry, actions *widget.Select, confirm_action *widget.Button) {
	windows.Program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(chn_id, actions, confirm_action)))
	windows.Program.Resize(fyne.NewSize(400, 150))
	confirm_action.SetText("Delete")
	confirm_action.OnTapped = func() {
		internalDeleteChannel(tkn, chn_id)
	}
}

func internalDeleteChannel(tkn, chn_id *widget.Entry) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://discord.com/api/v10/channels/%s", chn_id.Text), nil)
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
		dialog.ShowInformation("Success", "The channel has been successfully deleted!", windows.Program)
	}
}
