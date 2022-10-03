package client

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"steampipe-plugin-ndo/container"
)

const msoAuthPayload = `{
	"username": "%s",
	"password": "%s"
}`

const ndAuthPayload = `{
	"userName": "%s",
	"userPasswd": "%s"
}`

// Default timeout for NGINX in ACI is 90 Seconds.
// Allow the client to set a shorter or longer time depending on their
// environment
const DefaultReqTimeoutVal int = 100
const DefaultBackoffMinDelay int = 4
const DefaultBackoffMaxDelay int = 60
const DefaultBackoffDelayFactor float64 = 3
const DefaultAPIURL = "/mso/api/v1/"

// Client is the main entry point
type Client struct {
	BaseURL            *url.URL
	APIURL             string
	httpClient         *http.Client
	AuthToken          *Auth
	username           string
	password           string
	insecure           bool
	proxyUrl           string
	domain             string
	platform           string
	reqTimeoutSet      bool
	reqTimeoutVal      uint32
	skipLoggingPayload bool
	*ServiceManager
}

// singleton implementation of a client
var ndoClientImpl *Client

type Option func(*Client)

func Insecure(insecure bool) Option {
	return func(client *Client) {
		client.insecure = insecure
	}
}

func Password(password string) Option {
	return func(client *Client) {
		client.password = password
	}
}

func ProxyUrl(pUrl string) Option {
	return func(client *Client) {
		client.proxyUrl = pUrl
	}
}

func Domain(domain string) Option {
	return func(client *Client) {
		client.domain = domain
	}
}

func Platform(platform string) Option {
	return func(client *Client) {
		client.platform = platform
	}
}

func initClient(clientUrl, username string, options ...Option) *Client {
	var transport *http.Transport
	bUrl, err := url.Parse(clientUrl)
	if err != nil {
		// cannot move forward if url is undefined
		log.Fatal("Client URL not defined: %s", err)
	}
	client := &Client{
		BaseURL:  bUrl,
		username: username,
		APIURL:   DefaultAPIURL,
	}

	for _, option := range options {
		option(client)
	}

	if client.httpClient == nil {
		transport = client.useInsecureHTTPClient(client.insecure)
		if client.proxyUrl != "" {
			transport = client.configProxy(transport)
		}
		client.httpClient = &http.Client{
			Transport: transport,
		}
	}

	var timeout time.Duration
	if client.reqTimeoutSet {
		timeout = time.Second * time.Duration(client.reqTimeoutVal)
	} else {
		timeout = time.Second * time.Duration(DefaultReqTimeoutVal)
	}

	client.httpClient.Timeout = timeout
	client.ServiceManager = NewServiceManager(client.APIURL, client)
	return client
}

// GetClient returns a singleton
func GetClient(clientUrl, username string, options ...Option) *Client {
	if ndoClientImpl == nil {
		log.Printf("[DEBUG] Initializing new client...")
		ndoClientImpl = initClient(clientUrl, username, options...)
	} else {
		log.Printf("[DEBUG] Verifying cached client...")
		// making sure it is the same client
		bUrl, err := url.Parse(clientUrl)
		if err != nil {
			// cannot move forward if url is undefined
			log.Fatal("URL undefined: %s", err)
		}
		if bUrl != ndoClientImpl.BaseURL {
			log.Printf("[DEBUG] Cached client miss...")
			ndoClientImpl = initClient(clientUrl, username, options...)
		}
	}
	return ndoClientImpl
}

func (c *Client) configProxy(transport *http.Transport) *http.Transport {
	pUrl, err := url.Parse(c.proxyUrl)
	if err != nil {
		log.Fatal(err)
	}
	transport.Proxy = http.ProxyURL(pUrl)
	return transport
}

func (c *Client) useInsecureHTTPClient(insecure bool) *http.Transport {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			CipherSuites: []uint16{
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			},
			PreferServerCipherSuites: true,
			InsecureSkipVerify:       insecure,
			MinVersion:               tls.VersionTLS11,
			MaxVersion:               tls.VersionTLS12,
		},
	}

	return transport
}

func (c *Client) MakeRestRequest(method string, path string, body *container.Container, authenticated bool) (*http.Request, error) {
	if c.platform == "nd" && path != "/login" {
		if strings.HasPrefix(path, "/") {
			path = path[1:]
		}
		path = fmt.Sprintf("mso/%v", path)
	}
	url, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	fURL := c.BaseURL.ResolveReference(url)
	var req *http.Request
	if method == "GET" || method == "DELETE" {
		req, err = http.NewRequest(method, fURL.String(), nil)
	} else {
		req, err = http.NewRequest(method, fURL.String(), bytes.NewBuffer((body.Bytes())))
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	log.Printf("HTTP request %s %s", method, path)

	if authenticated {
		req, err = c.InjectAuthenticationHeader(req, path)
		if err != nil {
			return req, err
		}
	}
	log.Printf("HTTP request after injection %s %s", method, path)

	return req, nil
}

// Authenticate is used to
func (c *Client) Authenticate() error {
	method := "POST"
	path := "/api/v1/auth/login"
	var authPayload string

	if c.platform == "nd" {
		authPayload = ndAuthPayload
		if c.domain == "" {
			c.domain = "DefaultAuth"
		}
		path = "/login"
	} else {
		authPayload = msoAuthPayload
	}
	body, err := container.ParseJSON([]byte(fmt.Sprintf(authPayload, c.username, c.password)))
	if err != nil {
		return err
	}

	if c.domain != "" {
		if c.platform == "nd" {
			body.Set(c.domain, "domain")
		} else {
			domainId, err := c.GetDomainId(c.domain)
			if err != nil {
				return err
			}
			body.Set(domainId, "domainId")
		}
	}

	req, err := c.MakeRestRequest(method, path, body, false)
	if err != nil {
		return err
	}

	obj, _, err := c.Do(req)
	if err != nil {
		return err
	}

	if obj == nil {
		return errors.New("Empty response")
	}
	req.Header.Set("Content-Type", "application/json")

	token := StripQuotes(obj.S("token").String())

	if token == "" || token == "{}" {
		return errors.New("Invalid Username or Password")
	}

	if c.AuthToken == nil {
		c.AuthToken = &Auth{}
	}
	c.AuthToken.Token = stripQuotes(token)
	c.AuthToken.CalculateExpiry(1200) //refreshTime=1200 Sec

	return nil
}

func (c *Client) GetDomainId(domain string) (string, error) {
	req, err := c.MakeRestRequest("GET", "/api/v1/auth/login-domains", nil, false)
	if err != nil {
		return "", err
	}

	obj, _, err := c.Do(req)

	if err != nil {
		return "", err
	}
	err = CheckForErrors(obj, "GET", false)
	if err != nil {
		return "", err
	}
	count, err := obj.ArrayCount("domains")
	if err != nil {
		return "", err
	}

	for i := 0; i < count; i++ {
		domainCont, err := obj.ArrayElement(i, "domains")
		if err != nil {
			return "", err
		}
		domainName := StripQuotes(domainCont.S("name").String())

		if domainName == domain {
			return StripQuotes(domainCont.S("id").String()), nil
		}
	}
	return "", fmt.Errorf("Unable to find domain id for domain %s", domain)
}

func StrtoInt(s string, startIndex int, bitSize int) (int64, error) {
	return strconv.ParseInt(s, startIndex, bitSize)
}

func (c *Client) Do(req *http.Request) (*container.Container, *http.Response, error) {
	log.Printf("[DEBUG] Begining DO method %s", req.URL.String())
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	log.Printf("\nHTTP Request: %s %s", req.Method, req.URL.String())
	log.Printf("\nHTTP Response: %d %s %v", resp.StatusCode, resp.Status, resp)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	bodyStr := string(bodyBytes)
	resp.Body.Close()
	log.Printf("\nHTTP response unique string %s %s %s", req.Method, req.URL.String(), bodyStr)
	if req.Method != "DELETE" && resp.StatusCode != 204 {
		obj, err := container.ParseJSON(bodyBytes)

		if err != nil {
			log.Printf("Error occured while json parsing %+v", err)
			return nil, resp, err
		}
		log.Printf("[DEBUG] Exit from do method")
		return obj, resp, err
	} else if resp.StatusCode == 204 {
		return nil, nil, nil
	} else {
		return nil, resp, err
	}
}

func stripQuotes(word string) string {
	if strings.HasPrefix(word, "\"") && strings.HasSuffix(word, "\"") {
		return strings.TrimSuffix(strings.TrimPrefix(word, "\""), "\"")
	}
	return word
}
