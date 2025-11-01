package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"
)

func main() {
	// Setup routes
	http.HandleFunc("/fshare-links", func(w http.ResponseWriter, r *http.Request) {
		links, err := linksHandler(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		writeJSON(w, http.StatusOK, links)
	})

	// Start server
	port := ":8080"
	fmt.Printf("Server starting on port %s\n", port)
	fmt.Println("Available endpoints:")
	fmt.Println("  POST    /fshare-links        - Get fshare links")
	fmt.Println("  Example payload:")
	fmt.Println("  curl --location 'http://localhost:8080/fshare-links' \\")
	fmt.Println("  --header 'Content-Type: application/json' \\")
	fmt.Println("  --data '{")
	fmt.Println("      \"link\": \"https://www.fshare.vn/folder/IYZU2DHZ929T\",")
	fmt.Println("      \"cookie\": \"__cfduid=d4caae757dd449d8a3fdfaabddad0e48a1512302783; user=vietnth0602%40gmail.com; pass=78321e89c3e254e911a18c4de61837b1; __zlcmid=k0gyT6eoKkbPCg; PHPSESSID=9h5a793piavt55u53bu23os012; __atuvc=63%7C2%2C6%7C3\"")
	fmt.Println("  }'")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

// linksHandler
func linksHandler(r *http.Request) (FinalResponse, error) {
	if r.Method != http.MethodPost {
		return FinalResponse{}, fmt.Errorf("method not allowed")
	}

	// Parse JSON payload
	var payload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return FinalResponse{}, fmt.Errorf("failed to parse JSON payload: %w", err)
	}

	// Validate required fields
	if payload.Link == "" {
		return FinalResponse{}, fmt.Errorf("link is required")
	}

	if payload.Cookie == "" {
		return FinalResponse{}, fmt.Errorf("cookie is required")
	}

	links, err := fetchFshareLinks(
		"https://linksvip.net/GetLinkFs",
		fmt.Sprintf("%s&pass=undefined&hash=&captcha=undefined", payload.Link),
		payload.Cookie,
	)
	if err != nil {
		return FinalResponse{}, err
	}
	return links, nil
}

// fileGetContentsCurl is the Go equivalent of the PHP file_get_contents_curl function
// It makes HTTP requests with retry logic, similar to PHP's cURL implementation
func fileGetContentsCurl(url string, retries int, post string, cookie string) (string, error) {
	// Create HTTP client
	client := &http.Client{
		Timeout: 30 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Handle redirects if needed
			return nil
		},
	}

	// Create request
	var req *http.Request
	var err error

	if post != "" {
		// POST request
		req, err = http.NewRequest("POST", url, strings.NewReader(post))
		if err != nil {
			return "", fmt.Errorf("failed to create POST request: %w", err)
		}
	} else {
		// GET request
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			return "", fmt.Errorf("failed to create GET request: %w", err)
		}
	}

	// Set headers (equivalent to CURLOPT_HTTPHEADER)
	req.Header.Set("Host", "linksvip.net")
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Origin", "https://linksvip.net")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("DNT", "1")
	req.Header.Set("Referer", "https://linksvip.net/")
	req.Header.Set("Cookie", cookie)

	// Make the request
	resp, err := client.Do(req)

	if err != nil {
		// Retry logic if request fails
		if retries > 0 {
			time.Sleep(1 * time.Second)
			return fileGetContentsCurl(url, retries-1, post, cookie)
		}
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check for HTTP errors (equivalent to CURLOPT_FAILONERROR)
	if resp.StatusCode >= 400 {
		// Retry logic on HTTP error
		if retries > 0 {
			time.Sleep(1 * time.Second)
			return fileGetContentsCurl(url, retries-1, post, cookie)
		}
		return "", fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	result := string(body)

	// Check if result is empty and retry
	if result == "" && retries > 0 {
		time.Sleep(1 * time.Second)
		return fileGetContentsCurl(url, retries-1, post, cookie)
	}

	return result, nil
}

// cleanString removes non-printable characters except \r and \n (equivalent to preg_replace('/[^[:print:]\r\n]/', ”, ...))
func cleanString(s string) string {
	var result strings.Builder
	for _, r := range s {
		// Keep printable characters, \r, and \n
		if unicode.IsPrint(r) || r == '\r' || r == '\n' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// Response structure for the API response
type APIResponse struct {
	Trangthai string `json:"trangthai"`
	Filename  string `json:"filename"`
	Linkvip   string `json:"linkvip"`
	Loi       string `json:"loi"`
}

// Message structure
type Message struct {
	Text string `json:"text"`
}

// FinalResponse structure
type FinalResponse struct {
	Messages []Message `json:"messages"`
}

// RequestPayload structure for incoming JSON payload
type RequestPayload struct {
	Link   string `json:"link"`
	Cookie string `json:"cookie"`
}

func fetchFshareLinks(url string, link string, cookie string) (FinalResponse, error) {
	// Call fileGetContentsCurl (equivalent to PHP: file_get_contents_curl($URL, 5, $link))
	result, err := fileGetContentsCurl(url, 5, link, cookie)
	if err != nil {
		return FinalResponse{}, err
	}

	// Clean result - remove non-printable characters except \r and \n
	// (equivalent to preg_replace('/[^[:print:]\r\n]/', '', ...))
	cleanedResult := cleanString(result)

	// Parse JSON (equivalent to json_decode($result, TRUE))
	var apiResp APIResponse
	if err := json.Unmarshal([]byte(cleanedResult), &apiResp); err != nil {
		return FinalResponse{}, fmt.Errorf("failed to parse JSON: %w", err)
	}

	var finalResp FinalResponse

	// Check trangthai (equivalent to $echo['trangthai'] == 1)
	switch apiResp.Trangthai {
	case "1":
		// Build success response
		finalResp = FinalResponse{
			Messages: []Message{
				{Text: "Link của bạn đã sẵn sàng để download. <3"},
				{Text: "*Tên file:* " + apiResp.Filename},
				{Text: "*Download:* " + apiResp.Linkvip},
			},
		}
	default:
		finalResp = FinalResponse{
			Messages: []Message{
				{Text: apiResp.Loi},
			},
		}
	}

	return finalResp, nil
}

// writeJSON is a helper function to write JSON responses
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON: %v", err)
	}
}
