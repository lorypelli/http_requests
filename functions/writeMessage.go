package actions

import (
	b "bytes"
	j "encoding/json"
	"fmt"
	"io"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func WriteMessage(navbar_edit, msg_box *fyne.Container, msg, chn_id, tkn *widget.Entry, actions *widget.Select, confirm_action *widget.Button, program fyne.Window) {
	program.SetContent(container.NewBorder(navbar_edit, nil, nil, nil, container.NewVBox(chn_id, actions, msg_box, confirm_action)))
	program.Resize(fyne.NewSize(400, 240))
	confirm_action.SetText("Send")
	confirm_action.OnTapped = func() {
		body := map[string]interface{}{
			"content": msg.Text,
		}
		json, _ := j.Marshal(body)
		req, err := http.NewRequest("POST", fmt.Sprintf("https://discord.com/api/v10/channels/%s/messages", chn_id.Text), b.NewBuffer(json))
		if err != nil {
			dialog.ShowError(err, program)
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
		req.Header.Add("Content-Type", "application/json")
		c := &http.Client{}
		res, err := c.Do(req)
		if err != nil {
			dialog.ShowError(err, program)
		} else if res.StatusCode != 200 {
			var body struct {
				Message string
			}
			bytes, _ := io.ReadAll(res.Body)
			j.Unmarshal(bytes, &body)
			dialog.ShowInformation("Error", body.Message, program)
		} else {
			dialog.ShowInformation("Success", "The message has been successfully sent!", program)
		}
	}
}
