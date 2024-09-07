package actions

import (
	j "encoding/json"
	"fmt"
	"http_requests/globals"
	"io"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func BotInfo(tkn *widget.Entry) {
	req, err := http.NewRequest("GET", "https://discord.com/api/v10/users/@me", nil)
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
		var bots struct {
			Id       string
			Username string
			Avatar   string
		}
		bytes, _ := io.ReadAll(res.Body)
		j.Unmarshal(bytes, &bots)
		req, err = http.NewRequest("GET", "https://discord.com/api/v10/users/@me/guilds", nil)
		if err != nil {
			dialog.ShowError(err, windows.Program)
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn.Text))
		res, err = http.DefaultClient.Do(req)
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
			var guilds []struct{}
			bytes, _ = io.ReadAll(res.Body)
			j.Unmarshal(bytes, &guilds)
			var img fyne.Resource
			if bots.Avatar == "" {
				img, err = fyne.LoadResourceFromURLString("https://cdn.discordapp.com/embed/avatars/0.png")
			} else {
				img, err = fyne.LoadResourceFromURLString(fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png", bots.Id, bots.Avatar))
			}
			if err != nil {
				dialog.ShowError(err, windows.Program)
			} else {
				img_box := canvas.NewImageFromResource(img)
				img_box.FillMode = canvas.ImageFillContain
				img_box.SetMinSize(fyne.NewSquareSize(32))
				windows.Bot.SetContent(container.NewCenter(container.NewVBox(img_box, widget.NewLabel(fmt.Sprintf("Username: %s", bots.Username)), widget.NewLabel(fmt.Sprintf("ID: %s", bots.Id)), widget.NewLabel(fmt.Sprintf("Server Count: %d", len(guilds))))))
				windows.Bot.Show()
			}
		}
	}
}