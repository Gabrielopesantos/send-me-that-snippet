package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gabrielopesantos/smts/config"
	"github.com/gabrielopesantos/smts/internal/middleware"
	"github.com/gabrielopesantos/smts/internal/model"
	"github.com/gabrielopesantos/smts/internal/paste"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

type Server struct {
	app    *fiber.App
	dbConn *gorm.DB
	cfg    *config.Config
}

func New(dbConn *gorm.DB, cfg *config.Config) *Server {
	return &Server{
		app:    fiber.New(),
		dbConn: dbConn,
		cfg:    cfg,
	}
}

func (s *Server) Start() {
	s.mapRoutes()

	addr := fmt.Sprintf("%s:%s", s.cfg.ServerConfig.Host, s.cfg.ServerConfig.Port)
	err := s.app.Listen(addr)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Server) mapRoutes() {
	// This should be split
	mm := middleware.NewMiddlewareManager(s.cfg)
	if s.cfg.ServerConfig.CreateDashboardEndpoint {
		s.app.Get("/dashboard", mm.BasicAuthMiddleware(monitor.New()))
	}

	// V1
	v1Group := s.app.Group("/api/v1")

	// Paste
	go func() {
		startExpirePastesProcess(s.cfg.ServerConfig)
	}()
	pasteGroup := v1Group.Group("/pastes")
	pasteHandlers := paste.NewHandlers(s.dbConn, s.cfg)
	paste.MapRoutes(pasteGroup, pasteHandlers, mm)
}

func startExpirePastesProcess(srvCfg config.ServerConfig) {
	t := time.NewTicker(10 * time.Minute) // Time in minutes defined in config
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
	// Check which notes should be expired / Check if pastes is empty
	pastesToExpire := getPastesToExpire(pastes)
	// Expire them
	updatePastes(client, pastesToExpire, srvCfg)
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
