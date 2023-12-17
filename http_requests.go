package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	login := a.NewWindow("Login")
	program := a.NewWindow("http_requests")
	fileName := "app_icon.png"
	file, err := os.Open(fileName)
	if err == nil {
		stats, _ := os.Stat(fileName)
		size := stats.Size()
		fileByte := make([]byte, size)
		fileSlice := fileByte[:]
		file.Read(fileSlice)
		login.SetIcon(fyne.NewStaticResource(fileName, fileByte))
		program.SetIcon(fyne.NewStaticResource(fileName, fileByte))
	}
	login.Resize(fyne.NewSize(840, 170))
	program.Resize(fyne.NewSize(500, 250))
	login.SetFixedSize(true)
	program.SetFixedSize(true)
	login.CenterOnScreen()
	program.CenterOnScreen()
	tkn_textbox := widget.NewEntry()
	tkn_textbox.SetPlaceHolder("Insert bot token")
	login.SetContent(container.NewVBox(tkn_textbox, widget.NewButton("Validate", func() {
		req, err := http.NewRequest("POST", "https://discord.com/api/v10/auth/login", nil)
		if err != nil {
			dialog.ShowError(err, login)
		}
		req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn_textbox.Text))
		c := &http.Client{}
		res, err := c.Do(req)
		if err != nil {
			dialog.ShowError(err, login)
		} else if res.StatusCode != 200 {
			var body struct {
				Message string
			}
			bytes, _ := io.ReadAll(res.Body)
			json.Unmarshal(bytes, &body)
			dialog.ShowInformation("Error", body.Message, login)
		} else {
			req, err := http.NewRequest("GET", "https://discord.com/api/v10/users/@me", nil)
			if err != nil {
				dialog.ShowError(err, login)
			}
			req.Header.Add("Authorization", fmt.Sprintf("Bot %s", tkn_textbox.Text))
			c := &http.Client{}
			res, err := c.Do(req)
			if err != nil {
				dialog.ShowError(err, login)
			}
			var body struct {
				Id       string
				Username string
			}
			bytes, _ := io.ReadAll(res.Body)
			json.Unmarshal(bytes, &body)
			botId := body.Id
			botUsername := body.Username
			login.Hide()
			program.SetContent(container.NewBorder(container.NewHBox(widget.NewLabel(botId), layout.NewSpacer(), widget.NewButton("Logout", func() {
				dialog.ShowConfirm("Logout", "Are you sure you want to logout?", func(b bool) {
					if b {
						login.Show()
						program.Hide()
					}
				}, program)
			}), layout.NewSpacer(), widget.NewLabel(botUsername)), nil, nil, nil))
			program.Show()
		}
	})))
	login.Show()
	a.Run()
}
