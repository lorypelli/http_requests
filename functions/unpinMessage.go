package actions

import (
	"fmt"
	"http_requests/windows"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func UnpinMessage(navbar *fyne.Container, msg_id, tkn, chn_id *widget.Entry, actions *widget.Select, confirm_action *widget.Button) {
	windows.Program.SetContent(container.NewBorder(navbar, nil, nil, nil, container.NewVBox(chn_id, actions, msg_id, confirm_action)))
	windows.Program.Resize(fyne.NewSize(400, 200))
	confirm_action.SetText("Unpin")
	confirm_action.OnTapped = func() {
		internalUnpinMessage(msg_id, tkn, chn_id)
	}
}

func internalUnpinMessage(msg_id, tkn, chn_id *widget.Entry) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("https://discord.com/api/v10/channels/%s/pins/%s", chn_id.Text, msg_id.Text), nil)
	if err != nil {
		dialog.ShowError(err, windows.Program)
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		dialog.ShowError(err, windows.Program)
	} else if res.StatusCode != 204 {
		ShowError(res.Body)
	} else {
		dialog.ShowInformation("Success", "The message has been successfully unpinned!", windows.Program)
	}
}
