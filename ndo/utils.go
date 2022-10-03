package ndo

import (
	"context"
	"fmt"
	"log"

	"steampipe-plugin-ndo/client"

	"github.com/turbot/steampipe-plugin-sdk/plugin"
)

func connect(ctx context.Context, d *plugin.QueryData) (*client.Client, error) {
	log.Printf("[DEBUG] Getting connection to MSO")
	ndoConfig := GetConfig(d.Connection)

	log.Printf("[TRACE] connection found: %v", ndoConfig)

	// Initial values. Env vars will be overridden by configuration if values are set in there
	//vsphereServer := os.Getenv("VSPHERE_SERVER")
	//user := os.Getenv("VSPHERE_USER")
	//password := os.Getenv("VSPHERE_PASSWORD")
	allowUnverifiedSSL := false
	clusterURI := ""
	user := ""
	password := ""
	loginDomain := "DefaultAuth"
	platform := "nd"

	// Override potential env values with config values
	if ndoConfig.AllowUnverifiedSSL != nil {
		allowUnverifiedSSL = *ndoConfig.AllowUnverifiedSSL
	}

	if ndoConfig.ClusterURI != nil {
		clusterURI = *ndoConfig.ClusterURI
	}

	if ndoConfig.User != nil {
		user = *ndoConfig.User
	}

	if ndoConfig.Password != nil {
		password = *ndoConfig.Password
	}

	if ndoConfig.LoginDomain != nil {
		loginDomain = *ndoConfig.LoginDomain
	}

	if ndoConfig.Platform != nil {
		platform = *ndoConfig.Platform
	}
	// Make sure we have all required arguments set via either env or config
	if clusterURI == "" || user == "" || password == "" || platform == "" {
		errorMsg := ""
		if clusterURI == "" {
			errorMsg += "Missing clusterURI from config'\n"
		}
		if user == "" {
			errorMsg += "Missing user from config'\n"
		}
		if password == "" {
			errorMsg += "Missing password from config'\n"
		}
		return nil, fmt.Errorf("Error in configuraiton: %s", errorMsg)
	}

	log.Printf("[TRACE] Connection config:\n[TRACE] URI: %s\n[TRACE] User: %s\n[TRACE] Domain: %s\n[TRACE] Platform: %s", clusterURI, user, loginDomain, platform)

	if platform == "nd" {
		ndoClient := client.GetClient(clusterURI, user, client.Password(password), client.Insecure(allowUnverifiedSSL), client.Domain(loginDomain), client.Platform("nd"))
		log.Printf("[DEBUG] Got ND client")
		log.Printf("[TRACE] client: %v", ndoClient)
		return ndoClient, nil
	} else {
		ndoClient := client.GetClient(clusterURI, user, client.Password(password), client.Insecure(allowUnverifiedSSL), client.Domain(loginDomain))
		log.Printf("[DEBUG] Got MSO client")
		log.Printf("[TRACE] client: %v", ndoClient)
		return ndoClient, nil
	}
}
