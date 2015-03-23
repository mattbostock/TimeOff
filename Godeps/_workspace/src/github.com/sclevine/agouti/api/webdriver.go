package api

import (
	"fmt"
	"time"

	"github.com/sclevine/agouti/api/internal/bus"
	"github.com/sclevine/agouti/api/internal/service"
)

type WebDriver struct {
	Timeout  time.Duration
	service  driverService
	sessions []*Session
}

type driverService interface {
	URL() (string, error)
	Start() error
	Stop() error
	WaitForBoot(timeout time.Duration) error
}

func NewWebDriver(url string, command []string) *WebDriver {
	driverService := &service.Service{
		URLTemplate: url,
		CmdTemplate: command,
	}

	return &WebDriver{
		Timeout: 5 * time.Second,
		service: driverService,
	}
}

func (w *WebDriver) Open(desiredCapabilites map[string]interface{}) (*Session, error) {
	url, err := w.service.URL()
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve URL: %s", err)
	}

	busClient, err := bus.Connect(url, desiredCapabilites)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %s", err)
	}

	session := &Session{busClient}
	w.sessions = append(w.sessions, session)
	return session, nil
}

func (w *WebDriver) Start() error {
	if err := w.service.Start(); err != nil {
		return fmt.Errorf("failed to start service: %s", err)
	}

	if err := w.service.WaitForBoot(w.Timeout); err != nil {
		w.service.Stop()
		return err
	}

	return nil
}

func (w *WebDriver) Stop() error {
	for _, session := range w.sessions {
		session.Delete()
	}

	if err := w.service.Stop(); err != nil {
		return fmt.Errorf("failed to stop service: %s", err)
	}

	return nil
}
