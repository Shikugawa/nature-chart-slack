package nature

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	host = "https://api.nature.global"
)

type Response struct {
	Name            string `json:"name"`
	Id              string `json:"id"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	MacAddress      string `json:"mac_address"`
	SerialNumber    string `json:"serial_number"`
	FirmwareVersion string `json:"firmware_version"`
	TematureOffset  int    `json:"temature_offset"`
	HumidityOffset  int    `json:"humidity_offset"`
	Users           []struct {
		Id        string `json:"id"`
		Nickname  string `json:"nickname"`
		Superuser bool   `json:"superuser"`
	} `json:"users"`
	NewestEvents struct {
		Te struct {
			Value     float64 `json:"val"`
			CreatedAt string  `json:"created_at"`
		} `json:"te"`
	} `json:"newest_events"`
}

type Client struct {
	atoken string
}

func NewClient(token string) *Client {
	return &Client{
		atoken: token,
	}
}

func (c *Client) Request() ([]Response, error) {
	req, err := http.NewRequest(http.MethodGet, host+"/1/devices", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.atoken)

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode > http.StatusUnauthorized {
		return nil, fmt.Errorf("got status: %v", res.StatusCode)
	}

	var resb []Response
	b, _ := ioutil.ReadAll(res.Body)

	if err := json.Unmarshal(b, &resb); err != nil {
		return nil, err
	}

	return resb, nil
}
