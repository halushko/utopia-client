package page_navigator

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Client struct {
	loginClient    *http.Client
	navigateClient *http.Client
}

var c *Client
var clientMu sync.Mutex
var base = "https://utp.to"

var user = os.Getenv("UTP_USER")
var pass = os.Getenv("UTP_PASSWORD")

var loginURL = base + "/login"

func newClient() (*Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	login := &http.Client{
		Timeout: 30 * time.Second,
		Jar:     jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	scrape := &http.Client{Jar: jar, Timeout: 30 * time.Second}
	return &Client{loginClient: login, navigateClient: scrape}, nil
}

func setHeaders(req *http.Request) {
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT; Win64; x64)")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Referer", "https://utp.to/login")
}

func GET(pageURL string) ([]byte, error) {
	return get(pageURL, 2)
}

func get(pageURL string, tryCount int) ([]byte, error) {
	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		return nil, err
	}
	setHeaders(req)

	client, err := getClient()
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		if tryCount == 0 {
			return nil, fmt.Errorf("не вдалося залогінитися на Utopia")
		}
		time.Sleep(4 * time.Second)
		return get(pageURL, tryCount-1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch page%s: status %s", pageURL, resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	doc, parseErr := goquery.NewDocumentFromReader(strings.NewReader(string(body)))
	if parseErr == nil && doc.Find(`input[type="password"], form[action$="/login"]`).Length() > 0 {
		clientMu.Lock()
		c = nil
		clientMu.Unlock()
		return nil, errors.New("UTOPIA authentication required")
	}
	return body, nil
}

func login() error {
	client, err := newClient()
	if err != nil {
		return err
	}

	res, err := client.loginClient.Get(loginURL)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}
	token, ok := doc.Find(`input[name="_token"]`).Attr("value")
	if !ok {
		token, ok = doc.Find(`meta[name="csrf-token"]`).Attr("content")
	}
	if !ok {
		return errors.New("csrf token not found")
	}

	form := url.Values{
		"_token":   {token},
		"username": {user},
		"password": {pass},
		"remember": {"on"},
	}
	req, err := http.NewRequest("POST", loginURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	setHeaders(req)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.loginClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed: status %s", resp.Status)
	}

	u, _ := url.Parse(base)
	for _, ck := range client.loginClient.Jar.Cookies(u) {
		if ck.Name == "laravel_session" {
			c = client
			return nil
		}
	}
	return errors.New("login failed: session cookie not found")
}

func getClient() (*http.Client, error) {
	clientMu.Lock()
	defer clientMu.Unlock()
	if c == nil {
		if err := login(); err != nil {
			return nil, fmt.Errorf("UTOPIA login: %w", err)
		}
	}
	return c.navigateClient, nil
}
