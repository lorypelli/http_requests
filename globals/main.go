package windows

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
)

var App fyne.App

var Login, Program, Bot, Msg_list fyne.Window

func Setup() {
	App = app.New()
	Login = App.NewWindow("Login")
	Program = App.NewWindow("http_requests")
	Bot = App.NewWindow("bot_info")
	Msg_list = App.NewWindow("msg_list")
	req, _ := http.NewRequest("GET", "https://raw.githubusercontent.com/lorypelli/http_requests/main/icon.png", nil)
	res, _ := http.DefaultClient.Do(req)
	_, ok := App.(desktop.App)
	if ok {
		homeDirectory, _ := os.UserHomeDir()
		path := fmt.Sprintf("%s/http_requests/", homeDirectory)
		os.MkdirAll(path, os.ModePerm)
		filePath := filepath.Join(path, "icon.png")
		file, _ := os.Create(filePath)
		io.Copy(file, res.Body)
		file, _ = os.Open(filePath)
		stats, _ := os.Stat(filePath)
		size := stats.Size()
		fileBytes := make([]byte, size)
		fileSlice := fileBytes[:]
		file.Read(fileSlice)
		icon := fyne.NewStaticResource("icon.png", fileBytes)
		Login.SetIcon(icon)
		Program.SetIcon(icon)
		Bot.SetIcon(icon)
		Msg_list.SetIcon(icon)
	}
	Login.Resize(fyne.NewSize(670, 165))
	Program.Resize(fyne.NewSize(400, 240))
	Bot.Resize(fyne.NewSize(400, 200))
	Msg_list.Resize(fyne.NewSize(1280, 720))
	Login.SetFixedSize(true)
	Program.SetFixedSize(true)
	Bot.SetFixedSize(true)
	Msg_list.SetFixedSize(true)
	Login.CenterOnScreen()
	Program.CenterOnScreen()
	Bot.CenterOnScreen()
	Msg_list.CenterOnScreen()
	Login.SetCloseIntercept(func() {
		App.Quit()
	})
	Bot.SetCloseIntercept(func() {
		Bot.Hide()
	})
	Msg_list.SetCloseIntercept(func() {
		Msg_list.Hide()
	})
}