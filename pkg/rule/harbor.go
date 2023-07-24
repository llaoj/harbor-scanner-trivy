package rule

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const LabelIDLegal = 1
const LabelIDIlegal = 2

type HarborClient struct {
	baseUrl string
	auth    string
	client  *http.Client
}

type HistoryRecord struct {
	Created   string `json:"created"`
	CreatedBy string `json:"created_by"`
}

func NewHarborClient(baseUrl, auth string) *HarborClient {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	return &HarborClient{
		baseUrl: baseUrl,
		auth:    auth,
		client: &http.Client{
			Transport: tr,
			Timeout:   30 * time.Second},
	}
}

func (c *HarborClient) Get(url string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", c.auth)
	request.Header.Add("Accept", "application/json")
	log.WithFields(log.Fields{
		"url":  url,
		"auth": c.auth,
	}).Trace("Harbor Post API")
	return c.client.Do(request)
}

func (c *HarborClient) Post(url string, body io.Reader) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", c.auth)
	request.Header.Add("Accept", "application/json")
	request.Header.Add("Content-Type", "application/json")
	log.WithFields(log.Fields{
		"url":  url,
		"auth": c.auth,
	}).Trace("Harbor Post API")
	return c.client.Do(request)
}

func (c *HarborClient) Delete(url string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("Authorization", c.auth)
	request.Header.Add("Accept", "application/json")
	log.WithFields(log.Fields{
		"url":  url,
		"auth": c.auth,
	}).Trace("Harbor Delete API")
	return c.client.Do(request)
}

func parseRepository(repository string) (project, repo string, ok bool) {
	s := strings.Split(repository, "/")
	if len(s) == 2 {
		ok = true
		project = s[0]
		repo = s[1]
	}
	ok = false
	return
}

//	curl -X 'POST' \
//	  'https://registry3-qingdao.cosmoplat.com/api/v2.0/projects/61_dddyyy/repositories/memcached/artifacts/sha256%3A8d4a31b8f8e9d7764de35f8dbdd078c1372399c77c3530d71a16154041c52eeb/labels' \
//	  -H 'accept: application/json' \
//	  -H 'Content-Type: application/json' \
//	  -H 'X-Harbor-CSRF-Token: sgqQg1OmdmBlL8wmt6kF2uiHr/nb000wA7qQIB+6knT6BWrKCVHBcCf8/ZbglLED0lmCJ6w++wJmE3m2BZg4Kw==' \
//	  -d '{"id":2}'
//
// Response 200
// Response 409
//
//	{
//		"errors": [
//		  {
//			"code": "CONFLICT",
//			"message": "label 2 is already added to the artifact 64"
//		  }
//		]
//	}
func (c *HarborClient) AddLabel(project, repo, artifactRef string, lableID int) error {
	url := fmt.Sprintf("%s/api/v2.0/projects/%s/repositories/%s/artifacts/%s/labels",
		c.baseUrl,
		project,
		repo,
		artifactRef)
	body := []byte(`{"id":` + strconv.Itoa(lableID) + `}`)
	response, err := c.Post(url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	rawBody, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil
	}
	if response.StatusCode == http.StatusOK || response.StatusCode == http.StatusConflict {
		return nil
	}
	return errors.New(string(rawBody))
}

//	curl -X 'DELETE' \
//	  'https://registry3-qingdao.cosmoplat.com/api/v2.0/projects/61_dddyyy/repositories/memcached/artifacts/sha256%3A8d4a31b8f8e9d7764de35f8dbdd078c1372399c77c3530d71a16154041c52eeb/labels/1' \
//	  -H 'accept: application/json' \
//	  -H 'X-Harbor-CSRF-Token: ex68Iq6jTDTQG0F06dGQBc2Ne7TBIgGGhlAXYguitogzEUZr9FT7JJLIcMS+7CTc91NWarbPt7Tj+f70EYAc1w=='
//
// Response 200
// Response 404
//
//	{
//	    "errors": [
//	      {
//	        "code": "NOT_FOUND",
//	        "message": "reference with label 1 and artifact 64 not found"
//	      }
//	    ]
//	}
func (c *HarborClient) DeleteLabel(project, repo, artifactRef string, lableID int) error {
	url := fmt.Sprintf("%s/api/v2.0/projects/%s/repositories/%s/artifacts/%s/labels/%d",
		c.baseUrl,
		project,
		repo,
		artifactRef,
		lableID)

	response, err := c.Delete(url)
	if err != nil {
		return err
	}
	rawBody, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil
	}
	if response.StatusCode == http.StatusOK || response.StatusCode == http.StatusNotFound {
		return nil
	}
	return errors.New(string(rawBody))
}

func (c *HarborClient) AddLegalLabel(project, repo, artifactRef string) error {
	err := c.DeleteLabel(project, repo, artifactRef, LabelIDIlegal)
	if err != nil {
		return err
	}
	return c.AddLabel(project, repo, artifactRef, LabelIDLegal)
}

func (c *HarborClient) AddIlegalLabel(project, repo, artifactRef string) error {
	err := c.DeleteLabel(project, repo, artifactRef, LabelIDLegal)
	if err != nil {
		return err
	}
	return c.AddLabel(project, repo, artifactRef, LabelIDIlegal)
}

func (c *HarborClient) GetBuildHistory(project, repo, artifactRef string) ([]HistoryRecord, error) {
	url := fmt.Sprintf("%s/api/v2.0/projects/%s/repositories/%s/artifacts/%s/additions/build_history",
		c.baseUrl,
		project,
		repo,
		artifactRef)
	response, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	rawBody, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}
	if response.StatusCode == http.StatusOK {
		history := make([]HistoryRecord, 2)
		if err = json.Unmarshal(rawBody, &history); err != nil {
			return nil, err
		}
		return history, nil
	}
	return nil, errors.New(string(rawBody))
}
