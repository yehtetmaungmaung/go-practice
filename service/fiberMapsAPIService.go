package service

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// call fibermaps api and retrieve nodes. The nodes are then stored in &nodes.
func (nodes *Nodes) GetFiberMapsNodesByDeviceType(pageNo int, deviceType string) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := http.Client{
		Timeout:   60 * time.Second,
		Transport: tr,
	}
	url := os.Getenv("FIBERMAPS_URL") + "pets/api/nodes/external/node/check_port"
	query := "?page=" + strconv.Itoa(pageNo) + "&per_page=" + os.Getenv("PER_PAGE") + "&device_type=" + deviceType
	req, _ := http.NewRequest("GET", url+query, nil)
	req.Header.Set("Pre-Shared-Key", os.Getenv("FIBERMAPS_URL_PRESHAREDKEY"))
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		panic(nil)
	}
	defer resp.Body.Close()
	defer client.CloseIdleConnections()
	if resp.StatusCode == 200 {
		responseData, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		json.Unmarshal(responseData, &nodes)
	}
}
