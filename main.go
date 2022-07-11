package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type dnsList struct {
	Success  bool          `json:"success"`
	Errors   []interface{} `json:"errors"`
	Messages []interface{} `json:"messages"`
	Result   []struct {
		Id         string    `json:"id"`
		Type       string    `json:"type"`
		Name       string    `json:"name"`
		Content    string    `json:"content"`
		Proxiable  bool      `json:"proxiable"`
		Proxied    bool      `json:"proxied"`
		Ttl        int       `json:"ttl"`
		Locked     bool      `json:"locked"`
		ZoneId     string    `json:"zone_id"`
		ZoneName   string    `json:"zone_name"`
		CreatedOn  time.Time `json:"created_on"`
		ModifiedOn time.Time `json:"modified_on"`
		Data       struct {
		} `json:"data"`
		Meta struct {
			AutoAdded bool   `json:"auto_added"`
			Source    string `json:"source"`
		} `json:"meta"`
	} `json:"result"`
}

type updateDNS struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Content string `json:"content"`
	Ttl     int    `json:"ttl"`
	Proxied bool   `json:"proxied"`
}

type updateDNSResponse struct {
	Result struct {
		Id        string `json:"id"`
		ZoneId    string `json:"zone_id"`
		ZoneName  string `json:"zone_name"`
		Name      string `json:"name"`
		Type      string `json:"type"`
		Content   string `json:"content"`
		Proxiable bool   `json:"proxiable"`
		Proxied   bool   `json:"proxied"`
		Ttl       int    `json:"ttl"`
		Locked    bool   `json:"locked"`
		Meta      struct {
			AutoAdded           bool   `json:"auto_added"`
			ManagedByApps       bool   `json:"managed_by_apps"`
			ManagedByArgoTunnel bool   `json:"managed_by_argo_tunnel"`
			Source              string `json:"source"`
		} `json:"meta"`
		CreatedOn  time.Time `json:"created_on"`
		ModifiedOn time.Time `json:"modified_on"`
	} `json:"result"`
	Success  bool          `json:"success"`
	Errors   []interface{} `json:"errors"`
	Messages []interface{} `json:"messages"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	baseURL := "https://api.cloudflare.com/client/v4"
	zoneID := os.Getenv("ZONE_ID")
	request, err := http.NewRequest("GET", fmt.Sprintf("%s/zones/%s/dns_records", baseURL, zoneID), nil)
	if err != nil {
		log.Panic(err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("TOKEN")))

	var client = &http.Client{}
	var data dnsList
	response, err := client.Do(request)
	if err != nil {
		log.Panic(err)
	}
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		log.Panic(err)
	}
	publicIP := getPublicIP()
	log.Println("Public IP:", publicIP)
	if data.Success {
		for _, res := range data.Result {
			var dns updateDNS
			if res.Type == "A" {
				log.Println(fmt.Sprintf("Zone %s", res.Name))
				dns.Type = "A"
				dns.Name = res.Name
				dns.Ttl = 600
				dns.Proxied = true
				dns.Content = publicIP
			}
			putUrl := fmt.Sprintf("%s/zones/%s/dns_records/%s", baseURL, zoneID, res.Id)
			body, err := json.Marshal(dns)
			postBody := bytes.NewBuffer(body)
			request, err = http.NewRequest("PUT", putUrl, postBody)
			if err != nil {
				log.Panic(err)
			}
			request.Header.Set("Content-Type", "application/json")
			request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("TOKEN")))
			var client = &http.Client{}
			response, err = client.Do(request)
			if err != nil {
				log.Panic(err)
			}
			var updateDNSResponse updateDNSResponse
			defer response.Body.Close()
			err = json.NewDecoder(response.Body).Decode(&updateDNSResponse)
			if err != nil {
				log.Println(err)
			}
			if updateDNSResponse.Success {
				log.Println(fmt.Sprintf("Success update Zone %s", res.Name))
			} else {
				log.Println(fmt.Sprintf("Error update Zone %s", res.Name))
				log.Println(updateDNSResponse.Errors)
			}
			log.Println("================")
		}
	}
}

func getPublicIP() string {
	url := "https://checkip.amazonaws.com"
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("getPublicIP", "err", err)
	}
	var client = &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Panic(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	bodyString := string(body)
	return bodyString
}
