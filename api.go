package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type payload struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	N      int    `json:"n"`
	Size   string `json:"size"`
}

func generateImage(c *gin.Context) {
	data := payload{
		Model:  "dall-e-3",
		Prompt: "Create a random abstract art image using dark colors and geometric shapes",
		N:      1,
		Size:   "1024x1024",
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal payload"})
		return
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/images/generations", body)
	if err != nil {
		// handle err
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", os.ExpandEnv("Bearer $OPENAI_API_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send request"})
		return
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle err
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response body"})
		return
	}

	c.Data(http.StatusCreated, "application/json", respBody)
}
