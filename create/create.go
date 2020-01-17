package create

import (
	"crypto/tls"
	"fmt"
	"log"
	"io/ioutil"
    "time"

    "github.com/terraform-providers/terraform-provider-nutanix/utils"
	"github.com/go-resty/resty/v2"
	v3 "github.com/terraform-providers/terraform-provider-nutanix/client/v3"
)

// NutanixClient represents Nutanix API client
type NutanixClient struct {
	*resty.Client
}

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

// Create the image record in Nutanix
func (c *NutanixClient) CreateImage(name string) (interface{}, error) {
	currentTime := time.Now()
	response, err := c.R().
	    SetHeader("Content-Type", "application/json").
	    SetBody(`{"spec": {"name": "` + name + `", "description": "Built with Packer on ` + currentTime.Format("2006-01-02") +
	        `", "resources": {"image_type": "DISK_IMAGE"} }, "metadata": {"kind": "image"} }`).
	    SetResult(&v3.ImageIntentResponse{}).
	    Post("/images")
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode() == 202 {
	    result := response.Result().(*v3.ImageIntentResponse)
	    imageUUID := utils.StringValue(result.Metadata.UUID)
	    return imageUUID, nil
	}
	return nil, fmt.Errorf("Unexpected response from Nutanix API: %s", response.Status())
}

// Upload an image to Nutanix
func (c *NutanixClient) UploadImage(uuid string, source string) error {
	fileBytes, err := ioutil.ReadFile(source)
	if err != nil {
		return fmt.Errorf("error: Cannot read file %s", err)
	}

	response, err := c.R().
	    SetHeader("Content-Type", "application/octet-stream").
	    SetBody(fileBytes).
	    Put(fmt.Sprintf("/images/%s/file", uuid))
	if err != nil {
		log.Fatal(err)
	}

	if response.StatusCode() == 200 {
		return nil
	}

	return fmt.Errorf("Unexpected response from Nutanix API: %s", response.Status())
}
