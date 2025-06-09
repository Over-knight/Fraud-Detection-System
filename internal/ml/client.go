package ml

import (
    "bytes"
    "encoding/json"
    "net/http"
    "time"
)

// Client to call the ML service

type Client struct {
    url string
    http *http.Client
}

func NewClient(url string) *Client {
    return &Client{
        url: url,
        http: &http.Client{Timeout: 5 * time.Second},
    }
}

func (c *Client) Predict(tx Transaction) (string, error) {
    payload, _ := json.Marshal(tx)
    resp, err := c.http.Post(c.url+"/predict", "application/json", bytes.NewBuffer(payload))
    if err != nil {
        return "error", err
    }
    defer resp.Body.Close()
    var result struct{Verdict string `json:"verdict"`}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return "error", err
    }
    return result.Verdict, nil
}