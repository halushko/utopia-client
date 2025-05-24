// auth.go
package page_navigator

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Client struct {
	loginClient  *http.Client
	scrapeClient *http.Client
}

var c *Client
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
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	scrape := &http.Client{Jar: jar}
	return &Client{loginClient: login, scrapeClient: scrape}, nil
}

func Login() error {
	client, err := newClient()
	if err != nil {
		log.Fatalf("Init error: %v", err)
		return err
	}
	c = client

	res, err := c.loginClient.Get(loginURL)
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
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT; Win64; x64)")
	req.Header.Set("Referer", loginURL)

	resp, err := c.loginClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed: status %s", resp.Status)
	}

	u, _ := url.Parse(base)
	for _, ck := range c.loginClient.Jar.Cookies(u) {
		if ck.Name == "laravel_session" {
			return nil
		}
	}
	return errors.New("login failed: session cookie not found")
}

func HTTP() *http.Client {
	if c == nil {
		if err := Login(); err != nil {
			log.Fatalf("Init error: %v", err)
		}
	}
	return c.scrapeClient
}
