package client

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"steampipe-plugin-ndo/container"
)

type ServiceManager struct {
	APIURL string
	client *Client
}

func NewServiceManager(moURL string, client *Client) *ServiceManager {

	sm := &ServiceManager{
		APIURL: moURL,
		client: client,
	}
	return sm
}

func (sm *ServiceManager) Get(dn string) (*container.Container, error) {
	finalURL := fmt.Sprintf("%s/%s.json", sm.APIURL, dn)
	req, err := sm.client.MakeRestRequest("GET", finalURL, nil, true)

	if err != nil {
		return nil, err
	}

	obj, _, err := sm.client.Do(req)
	if err != nil {
		return nil, err
	}

	if obj == nil {
		return nil, errors.New("Empty response body")
	}
	log.Printf("[DEBUG] Exit from GET %s", finalURL)
	return obj, CheckForErrors(obj, "GET", sm.client.skipLoggingPayload)
}

func createJsonPayload(payload map[string]string) (*container.Container, error) {
	containerJSON := []byte(fmt.Sprintf(`{
		"%s": {
			"attributes": {
			}
		}
	}`, payload["classname"]))

	return container.ParseJSON(containerJSON)
}

func StripQuotes(word string) string {
	if strings.HasPrefix(word, "\"") && strings.HasSuffix(word, "\"") {
		return strings.TrimSuffix(strings.TrimPrefix(word, "\""), "\"")
	}
	return word
}

func G(cont *container.Container, key string) string {
	return StripQuotes(cont.S(key).String())
}

// CheckForErrors parses the response and checks of there is an error attribute in the response
func CheckForErrors(cont *container.Container, method string, skipLoggingPayload bool) error {
	return nil
}

func (sm *ServiceManager) GetViaURL(url string) (*container.Container, error) {
	req, err := sm.client.MakeRestRequest("GET", url, nil, true)

	if err != nil {
		return nil, err
	}

	obj, _, err := sm.client.Do(req)
	if !sm.client.skipLoggingPayload {
		log.Printf("Getvia url %+v", obj)
	}
	if err != nil {
		return nil, err
	}

	if obj == nil {
		return nil, errors.New("Empty response body")
	}
	return obj, CheckForErrors(obj, "GET", sm.client.skipLoggingPayload)

}
