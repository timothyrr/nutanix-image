package get

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
)

// NutanixClient represents Nutanix API client
type NutanixClient struct {
	*resty.Client
}

// NutanixImage represents an Image returned from API
type NutanixImage struct {
	APIVersion string `json:"api_version"`
	Entities   []struct {
		Metadata struct {
			Categories  struct{} `json:"categories"`
			Kind        string   `json:"kind"`
			SpecVersion int64    `json:"spec_version"`
			UUID        string   `json:"uuid"`
		} `json:"metadata"`
		Spec struct {
			Name      string `json:"name"`
			Resources struct {
				ImageType string `json:"image_type"`
				SourceURI string `json:"source_uri"`
				ArchType  string `json:"architecture"`
			} `json:"resources"`
		} `json:"spec"`
		Status struct {
			Name      string `json:"name"`
			Resources struct {
				RetrievalUriList   []interface{} `json:"retrieval_uri_list"`
				SizeBytes          int64         `json:"size_bytes"`
				ImageType          string        `json:"image_type"`
				SourceURI          string        `json:"source_uri"`
				ArchType           string        `json:"architecture"`
			} `json:"resources"`
			State string `json:"state"`
		} `json:"status"`
	} `json:"entities"`
	Metadata struct {
		Filter       string `json:"filter"`
		Kind         string `json:"kind"`
		Length       int64  `json:"length"`
		Offset       int64  `json:"offset"`
		TotalMatches int64  `json:"total_matches"`
	} `json:"metadata"`
}

// NewNutanixClient instantiates REST client for Nutanix
func NewNutanixClient(url string, username string, password string, insecure bool, debug bool) NutanixClient {
	client := resty.New()
	if debug {
		client.SetDebug(true)
	}
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: insecure})
	client.SetBasicAuth(username, password)
	client.SetHostURL("https://" + url + ":9440" + "/api/nutanix/v3")
	return NutanixClient{client}
}

// GetImage retrieves Images by name from Nutanix
func (c *NutanixClient) GetImage(name string) (*NutanixImage, error) {
	response, err := c.R().
		SetBody(map[string]interface{}{"filter": "name==" + name}).
		SetResult(&NutanixImage{}).
		Post("/images/list")
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode() == 200 {
		image := response.Result().(*NutanixImage)
		if len(image.Entities) > 0 {
			return image, nil
		}
		return new(NutanixImage), nil
	}
	return nil, fmt.Errorf("Unexpected response from Nutanix API: %s", response.Status())
}

// DownloadImage downloads the given Nutanix image by UUID
func (c *NutanixClient) DownloadImage(out_file string, uuid string, name string) error {
    out, err := os.Create(out_file)
    if err != nil {
        return err
    }
    defer out.Close()

	response, err := c.R().SetOutput(out_file).Get(fmt.Sprintf("/images/%s/file", uuid))
	if err != nil {
		log.Fatal(err)
	}

	if response.StatusCode() == 202 {
		return nil
	}

	return fmt.Errorf("Unexpected response from Nutanix API: %s", response.Status())
}
