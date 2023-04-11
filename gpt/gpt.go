package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var completionsParam struct {
	Prompt string `json:"prompt"`
}

var chatParam struct {
	Msg string `json:"msg"`
}

func main() {
	fmt.Println("start")

	http.HandleFunc("/chat", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "text/event-stream")
		rw.Header().Set("Cache-Control", "no-cache")
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE")
		rw.Header().Set("Access-Control-Allow-Headers", "*")
		rw.Header().Set("Connection", "keep-alive")
		rw.Header().Set("Keep-Alive", "timeout=10")

		if r.Method == "OPTIONS" {
			return
		}
		// fmt.Println("=============================================")
		var url string = "https://api.openai.com/v1/completions"

		json.NewDecoder(r.Body).Decode(&chatParam)

		body := fmt.Sprintf(`{
			"model": "gpt-3.5-turbo",
			"stream":true,
			"messages":[{"role": "user", "content": "%s"}]
		}`, chatParam.Msg)
		var jsonStr = []byte(body)

		client := &http.Client{}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		if err != nil {
			fmt.Println("req ===", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer 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")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("resp ===", err)
			return
		}
		defer resp.Body.Close()
		// data, err := ioutil.ReadAll(resp.Body)
		// fmt.Println("abc", string(data))

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			eventData := scanner.Text()
			if eventData == "" {
				continue
			}
			// fmt.Fprintf(rw, "%s", eventData)
			// eventData = `data: {"id":"chatcmpl-744OIi2u1F8f9K7Xm1S5C1ferP6XC","object":"chat.completion.chunk","created":1681204546,"model":"gpt-3.5-turbo-0301","choices":[{"delta":{"role":"assistant"},"index":0,"finish_reason":null}]}`
			fmt.Fprintf(rw, "%s\n\n", eventData)
			// fmt.Println(eventData)
			flusher, ok := rw.(http.Flusher)
			if ok {
				flusher.Flush()
			} else {
				log.Println("Flushing not supported")
			}
		}
	})

	http.HandleFunc("/completions", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "text/event-stream")
		rw.Header().Set("Cache-Control", "no-cache")
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Methods", "POST,GET,OPTIONS,PUT,DELETE")
		rw.Header().Set("Access-Control-Allow-Headers", "*")
		rw.Header().Set("Connection", "keep-alive")
		rw.Header().Set("Keep-Alive", "timeout=10")

		if r.Method == "OPTIONS" {
			return
		}
		// fmt.Println("=============================================")
		var url string = "https://api.openai.com/v1/completions"

		json.NewDecoder(r.Body).Decode(&completionsParam)

		body := fmt.Sprintf(`{
			"model": "text-davinci-003",
			"stream":true,
			"max_tokens":4000,
			"prompt": "%s"
		}`, completionsParam.Prompt)
		var jsonStr = []byte(body)

		client := &http.Client{}
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
		if err != nil {
			fmt.Println("req ===", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer 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")
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("resp ===", err)
			return
		}
		defer resp.Body.Close()
		// data, err := ioutil.ReadAll(resp.Body)
		// fmt.Println("abc", string(data))

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			eventData := scanner.Text()
			if eventData == "" {
				continue
			}
			// fmt.Fprintf(rw, "%s", eventData)
			// eventData = `data: {"id":"chatcmpl-744OIi2u1F8f9K7Xm1S5C1ferP6XC","object":"chat.completion.chunk","created":1681204546,"model":"gpt-3.5-turbo-0301","choices":[{"delta":{"role":"assistant"},"index":0,"finish_reason":null}]}`
			fmt.Fprintf(rw, "%s\n\n", eventData)
			// fmt.Println(eventData)
			flusher, ok := rw.(http.Flusher)
			if ok {
				flusher.Flush()
			} else {
				log.Println("Flushing not supported")
			}
		}

		// fmt.Fprintf(rw, "娃哈哈")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
