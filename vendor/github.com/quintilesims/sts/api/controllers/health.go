package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/quintilesims/sts/models"
	"github.com/zpatrick/fireball"
	"log"
	"os"
	"strings"
	"time"
)

const TIME_FORMAT = ""

type HealthController struct {
	created time.Time
	mode    string
}

func NewHealthController() *HealthController {
	return &HealthController{
		created: time.Now(),
		mode:    "normal",
	}
}

func (h *HealthController) Routes() []*fireball.Route {
	return []*fireball.Route{
		{
			Path: "/health",
			Handlers: fireball.Handlers{
				"GET":  h.getHealth,
				"POST": h.setHealth,
			},
		},
	}
}

func (h *HealthController) getHealth(c *fireball.Context) (fireball.Response, error) {
	if h.mode == "slow" {
		log.Println("[Mode: Slow] waiting 20 seconds before proceeding")
		time.Sleep(time.Second * 20)
	}

	health := models.Health{
		TimeCreated: h.created.Format(TIME_FORMAT),
		Mode:        h.mode,
	}

	return fireball.NewJSONResponse(200, health)
}

func (h *HealthController) setHealth(c *fireball.Context) (fireball.Response, error) {
	var req models.SetHealthRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&req); err != nil {
		return nil, err
	}

	switch mode := strings.ToLower(req.Mode); mode {
	case "normal", "slow":
		log.Printf("Running in %s mode\n", mode)
		h.mode = mode
	case "die":
		h.mode = mode
		log.Println("[Mode: Die] exiting program in 5 seconds")
		go func() {
			time.Sleep(time.Second * 5)
			os.Exit(1)
		}()
	default:
		return nil, fmt.Errorf("Unknown mode '%s'", req.Mode)
	}

	health := models.Health{
		TimeCreated: h.created.Format(TIME_FORMAT),
		Mode:        h.mode,
	}

	return fireball.NewJSONResponse(200, health)
}
