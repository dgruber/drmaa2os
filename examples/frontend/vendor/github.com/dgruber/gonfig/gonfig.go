package gonfig

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/cloudfoundry-community/go-cfenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"io/ioutil"
	"net/http"
	"os"
)

type ConfigServerResponse struct {
	Name     string   `json:"name"`
	Profiles []string `json:"profiles"`
	Label    string   `json:"label"`
	Version  string   `json:"version"`
	// State           interface{} `json:"state"`
	PropertySources []struct {
		// Name is the path to the github repository for example
		Name string `json:"name"`
		// Source contains the actual configuration (float64 for ints)
		Source map[string]interface{} `json:"source"` // Rest of the fields should go here.
	} `json:"propertySources"`
}

type VCAPServices []struct {
	Service []struct {
		Name  string      `json:"name"`
		Label string      `json:"label"`
		Tags  []string    `json:"tags"`
		Plan  string      `json:"plan"`
		Cred  Credentials `json:"credentials"`
	} `json:"-"`
}

type Credentials struct {
	AccessTokenURI string
	ClientID       string
	ClientSecret   string
	// URI is the URI of the config service to make the request
	URI string
}

func (c *Credentials) GetConfigurationFromServer() (map[string]interface{}, error) {
	conf := &clientcredentials.Config{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		TokenURL:     c.AccessTokenURI,
	}

	var client *http.Client
	if os.Getenv("gonfig_testing") == "1" {
		client = &http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}
	} else {
		client = oauth2.NewClient(oauth2.NoContext, conf.TokenSource(oauth2.NoContext))
	}

	resp, errGet := client.Get(c.URI)
	if errGet != nil {
		return nil, fmt.Errorf("Error during getting configuration from URL %s: %s", c.URI, errGet.Error())
	}
	defer resp.Body.Close()

	body, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil {
		return nil, fmt.Errorf("Error during reading response body from config server: %s", errRead.Error())
	}

	var response ConfigServerResponse
	errDecode := json.Unmarshal(body, &response)
	if errDecode != nil {
		return nil, fmt.Errorf("Error during deconding response from config server: %s", errDecode.Error())
	} else if len(response.PropertySources) == 0 {
		return nil, fmt.Errorf("Response from config server has zero length: %s", body)
	}

	fmt.Printf("Configuration: \n %v\n", response)
	return response.PropertySources[0].Source, nil
}

func getConfigServerCredentials() (*Credentials, error) {
	env, err := cfenv.Current()
	if err != nil {
		return nil, fmt.Errorf("Error during fetching current configuration: %s", err.Error())
	}

	services, err := env.Services.WithLabel("p-config-server")
	if err != nil {
		return nil, err
	}

	var credentials Credentials

	credentials.AccessTokenURI, _ = services[0].CredentialString("access_token_uri")
	if credentials.AccessTokenURI == "" {
		return nil, fmt.Errorf("Hostname not found in credentials")
	}
	credentials.ClientID, _ = services[0].CredentialString("client_id")
	if credentials.ClientID == "" {
		return nil, fmt.Errorf("Username not found in credentials")
	}
	credentials.ClientSecret, _ = services[0].CredentialString("client_secret")
	if credentials.ClientSecret == "" {
		return nil, fmt.Errorf("ClientSecret not found in credentials")
	}
	uri, _ := services[0].CredentialString("uri")
	if uri == "" {
		return nil, fmt.Errorf("URI not found in credentials")
	}
	// URI + App name (name of the *.xml) + App profile (space name) + label (default master / git branch)
	credentials.URI = fmt.Sprintf("%s/%s/%s/%s", uri, env.Name, env.SpaceName, "master")
	if os.Getenv("gonfig_testing") == "1" {
		credentials.URI = uri
	}
	return &credentials, nil
}

// FetchConfig returns the configuration given by the PCF Config Server which is bound
// to the app.
func FetchConfig() (map[string]interface{}, error) {
	credentials, err := getConfigServerCredentials()
	if err != nil {
		return nil, err
	}
	return credentials.GetConfigurationFromServer()
}
