build-amd64:
	GOOS=linux GOARCH=amd64 go build -o bin/cloudflare_update_dns main.go