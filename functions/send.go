package actions

import (
	b "bytes"
	j "encoding/json"
	"fmt"
	"io"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func Send(msg, chn_id, tkn *widget.Entry, program fyne.Window) {
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
	res, err := http.DefaultClient.Do(req)
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
