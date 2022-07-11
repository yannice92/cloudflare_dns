update your public IP to A record on DNS cloudflare

**Warning: it will update all your A records to existing public IP**

## Requirement
* go 1.18

## How To Build

* type on your terminal `GOOS=linux GOARCH=amd64 go build -o bin/cloudflare_update_dns main.go`
* rename .env.example to .env
* input your token and zone ID from cloudflare dashboar

## How To Use

Just create scheduler and add cloudflare_update_dns binary
