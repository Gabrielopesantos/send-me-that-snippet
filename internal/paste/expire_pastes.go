package paste

import (
	"bytes"
	"encoding/json"
	"github.com/gabrielopesantos/smts/config"
	"github.com/gabrielopesantos/smts/internal/model"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func StartExpirePastesProcess(srvCfg config.ServerConfig) {
	// Need a way to gracefully stop this
	t := time.NewTicker(srvCfg.DeletePastesIntervalMins * time.Minute)
	client := http.Client{Timeout: 10 * time.Second}
	defer client.CloseIdleConnections()
	for {
		select {
		case <-t.C:
			expirePastes(&client, srvCfg)
		default:
		}
	}
}

func expirePastes(client *http.Client, srvCfg config.ServerConfig) {
	// getAllNonExpiredPastes
	pastes, _ := getNonExpiredPastes(client, srvCfg)
	log.Printf("Len pastes: %d", len(pastes))
	// Check which notes should be expired / Check if pastes is empty
	pastesToExpire := getPastesToExpire(pastes)
	// Expire them
	if len(pastesToExpire) > 0 {
		updatePastes(client, pastesToExpire, srvCfg)
	}
}

func getNonExpiredPastes(client *http.Client, srvCfg config.ServerConfig) ([]model.Paste, error) {
	u, _ := url.Parse("http://localhost:5000/api/v1/pastes/filter?expired=false")
	req, err := http.NewRequest("GET", u.String(), nil)
	req.SetBasicAuth(srvCfg.BasicAuthUser, srvCfg.BasicAuthPassword)
	if err != nil { // Cannot fail here
		log.Fatal(err)
		return nil, err
	}
	resp, err := client.Do(req)
	body, _ := io.ReadAll(resp.Body)
	var pastes []model.Paste
	err = json.Unmarshal(body, &pastes)
	if err != nil {
		return nil, err
	}

	return pastes, nil
}

func getPastesToExpire(pastes []model.Paste) []model.Paste {
	var pastesToExpire []model.Paste
	now := time.Now()
	for _, p := range pastes {
		if p.CreatedAt.Add(p.ExpiresIn * time.Minute).Before(now) {
			pastesToExpire = append(pastesToExpire, p)
		}
	}
	return pastesToExpire
}

func updatePastes(client *http.Client, pastes []model.Paste, srvCfg config.ServerConfig) {
	u, _ := url.Parse("http://localhost:5000/api/v1/pastes")
	body := model.Paste{Expired: true}
	jsonBody, _ := json.Marshal(body)
	for _, p := range pastes {
		req, _ := http.NewRequest("PUT", u.String()+"/"+p.Id, bytes.NewReader(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		req.SetBasicAuth(srvCfg.BasicAuthUser, srvCfg.BasicAuthPassword)
		resp, err := client.Do(req)
		if err != nil {
			log.Print(err)
		}
		log.Println(resp)
	}
}
