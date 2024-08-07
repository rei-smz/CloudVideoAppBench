package HTTPController

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

type HTTPController interface {
	Post(url string, body map[string]string) map[string]any
}

type httpController struct {
	client *http.Client
}

func NewHTTPController() HTTPController {
	return &httpController{
		client: &http.Client{},
	}
}

func (c *httpController) Post(url string, body map[string]string) map[string]any {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		log.Println("Error converting map to JSON:", err)
		return nil
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Println("Error creating request:", err)
		return nil
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		log.Println("Error sending request:", err)
		return nil
	}
	defer resp.Body.Close()

	var respBody []byte
	resp.Body.Read(respBody)
	ret := make(map[string]any)
	json.Unmarshal(respBody, &ret)
	ret["status_code"] = resp.StatusCode

	return ret
}
