package actions

import (
	"encoding/json"
	"http_requests/windows"
	"io"

	"fyne.io/fyne/v2/dialog"
)

func ShowError(b io.ReadCloser) {
	var body struct {
		Message string
	}
	bytes, _ := io.ReadAll(b)
	json.Unmarshal(bytes, &body)
	dialog.ShowInformation("Error", body.Message, windows.Program)
}
