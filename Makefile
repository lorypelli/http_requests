windows:
	@GOOS=windows go build -o bin/http_requests_win32.exe http_requests.go
linux:
	@GOOS=linux go build -o bin/http_requests_linux http_requests.go
darwin:
	@GOOS=darwin go build -o bin/http_requests_darwin http_requests.go
all: windows linux darwin