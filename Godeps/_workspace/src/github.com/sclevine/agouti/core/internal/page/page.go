package page

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sclevine/agouti/core/internal/types"
)

type Page struct {
	Client types.Client
}

func (p *Page) Destroy() error {
	if err := p.Client.DeleteSession(); err != nil {
		return fmt.Errorf("failed to destroy session: %s", err)
	}
	return nil
}

func (p *Page) Navigate(url string) error {
	if err := p.Client.SetURL(url); err != nil {
		return fmt.Errorf("failed to navigate: %s", err)
	}
	return nil
}

func (p *Page) SetCookie(cookie interface{}) error {
	if err := p.Client.SetCookie(cookie); err != nil {
		return fmt.Errorf("failed to set cookie: %s", err)
	}
	return nil
}

func (p *Page) DeleteCookie(name string) error {
	if err := p.Client.DeleteCookie(name); err != nil {
		return fmt.Errorf("failed to delete cookie %s: %s", name, err)
	}
	return nil
}

func (p *Page) ClearCookies() error {
	if err := p.Client.DeleteCookies(); err != nil {
		return fmt.Errorf("failed to clear cookies: %s", err)
	}
	return nil
}

func (p *Page) URL() (string, error) {
	url, err := p.Client.GetURL()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve URL: %s", err)
	}
	return url, nil
}

func (p *Page) Size(width, height int) error {
	window, err := p.Client.GetWindow()
	if err != nil {
		return fmt.Errorf("failed to retrieve window: %s", err)
	}

	if err := window.SetSize(width, height); err != nil {
		return fmt.Errorf("failed to set window size: %s", err)
	}

	return nil
}

func (p *Page) Screenshot(filename string) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0750); err != nil {
		return fmt.Errorf("failed to create directory for screenshot: %s", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file for screenshot: %s", err)
	}
	defer file.Close()

	screenshot, err := p.Client.GetScreenshot()
	if err != nil {
		os.Remove(filename)
		return fmt.Errorf("failed to retrieve screenshot: %s", err)
	}

	if _, err := file.Write(screenshot); err != nil {
		return fmt.Errorf("failed to write file for screenshot: %s", err)
	}

	return nil
}

func (p *Page) Title() (string, error) {
	title, err := p.Client.GetTitle()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve page title: %s", err)
	}
	return title, nil
}

func (p *Page) HTML() (string, error) {
	html, err := p.Client.GetSource()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve page HTML: %s", err)
	}
	return html, nil
}

func (p *Page) RunScript(body string, arguments map[string]interface{}, result interface{}) error {
	var (
		keys   []string
		values []interface{}
	)

	for key, value := range arguments {
		keys = append(keys, key)
		values = append(values, value)
	}

	argumentList := strings.Join(keys, ", ")
	cleanBody := fmt.Sprintf("return (function(%s) { %s; }).apply(this, arguments);", argumentList, body)

	if err := p.Client.Execute(cleanBody, values, result); err != nil {
		return fmt.Errorf("failed to run script: %s", err)
	}

	return nil
}

func (p *Page) PopupText() (string, error) {
	text, err := p.Client.GetAlertText()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve popup text: %s", err)
	}
	return text, nil
}

func (p *Page) EnterPopupText(text string) error {
	if err := p.Client.SetAlertText(text); err != nil {
		return fmt.Errorf("failed to enter popup text: %s", err)
	}
	return nil
}

func (p *Page) ConfirmPopup() error {
	if err := p.Client.SetAlertText("\u000d"); err != nil {
		return fmt.Errorf("failed to confirm popup: %s", err)
	}
	return nil
}

func (p *Page) CancelPopup() error {
	if err := p.Client.SetAlertText("\u001b"); err != nil {
		return fmt.Errorf("failed to cancel popup: %s", err)
	}
	return nil
}

func (p *Page) Forward() error {
	if err := p.Client.Forward(); err != nil {
		return fmt.Errorf("failed to navigate forward in history: %s", err)
	}
	return nil
}

func (p *Page) Back() error {
	if err := p.Client.Back(); err != nil {
		return fmt.Errorf("failed to navigate backwards in history: %s", err)
	}
	return nil
}

func (p *Page) Refresh() error {
	if err := p.Client.Refresh(); err != nil {
		return fmt.Errorf("failed to refresh page: %s", err)
	}
	return nil
}

func (p *Page) SwitchToParentFrame() error {
	if err := p.Client.FrameParent(); err != nil {
		return fmt.Errorf("failed to switch to parent frame: %s", err)
	}
	return nil
}

func (p *Page) SwitchToRootFrame() error {
	if err := p.Client.Frame(nil); err != nil {
		return fmt.Errorf("failed to switch to original page frame: %s", err)
	}
	return nil
}