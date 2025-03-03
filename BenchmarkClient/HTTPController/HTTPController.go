package HTTPController

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

type HTTPController interface {
	Post(url string, body map[string]any) map[string]any
}

type httpController struct {
	client *http.Client
}

func NewHTTPController() HTTPController {
	return &httpController{
		client: &http.Client{},
	}
}

func (c *httpController) Post(url string, body map[string]any) map[string]any {
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

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return nil
	}
	//log.Println(string(respBody))
	ret := make(map[string]any)
	json.Unmarshal(respBody, &ret)
	ret["status_code"] = resp.StatusCode

	return ret
}
