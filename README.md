# Get Folder Links VIP Net

A simple Go REST API service to fetch and process Fshare links from linksvip.net. This application provides an endpoint to retrieve download links for Fshare files and folders.

## Features

- 🚀 Simple REST API built with Go standard library
- 🔄 Automatic retry logic for failed requests
- 📦 Docker support with multi-stage builds
- 🔒 Secure HTTPS requests with proper headers
- 📝 Clean JSON request/response format
- ⚡ Lightweight and fast

## Prerequisites

- Go 1.21 or higher
- Docker (optional, for containerized deployment)

## Installation

### Using Go

1. Clone the repository:
```bash
git clone <repository-url>
cd get-folder-linksvip-net
```

2. Build the application:
```bash
go build -o app main.go
```

3. Run the application:
```bash
./app
```

The server will start on port `8080`.

### Using Docker

1. Build the Docker image:
```bash
docker build -t get-folder-linksvip-net .
```

2. Run the container:
```bash
docker run -p 8080:8080 get-folder-linksvip-net
```

For detached mode:
```bash
docker run -d -p 8080:8080 --name fshare-api get-folder-linksvip-net
```

## API Documentation

### Endpoint

**POST** `/fshare-links`

Get Fshare download links from linksvip.net.

### Request

**Content-Type:** `application/json`

**Request Body:**
```json
{
  "link": "https://www.fshare.vn/folder/IYZU2DHZ929T",
  "cookie": "__cfduid=...; user=...; pass=...; PHPSESSID=...; ..."
}
```

**Parameters:**
- `link` (string, required): The Fshare folder or file link
- `cookie` (string, required): Your linksvip.net session cookie string

### Response

**Success Response (200 OK):**
```json
{
  "messages": [
    {
      "text": "Link của bạn đã sẵn sàng để download. <3"
    },
    {
      "text": "*Tên file:* filename.ext"
    },
    {
      "text": "*Download:* https://download-link..."
    }
  ]
}
```

**Error Response (400/500):**
```json
{
  "messages": [
    {
      "text": "Error message here"
    }
  ]
}
```

### Example Usage

**Using curl:**
```bash
curl --location 'http://localhost:8080/fshare-links' \
  --header 'Content-Type: application/json' \
  --data '{
    "link": "https://www.fshare.vn/folder/IYZU2DHZ929T",
    "cookie": "__cfduid=d4caae757dd449d8a3fdfaabddad0e48a1512302783; user=vietnth0602%40gmail.com; pass=78321e89c3e254e911a18c4de61837b1; __zlcmid=k0gyT6eoKkbPCg; PHPSESSID=9h5a793piavt55u53bu23os012; __atuvc=63%7C2%2C6%7C3"
  }'
```

**Using HTTPie:**
```bash
http POST http://localhost:8080/fshare-links \
  link="https://www.fshare.vn/folder/IYZU2DHZ929T" \
  cookie="__cfduid=...; user=...; pass=..."
```

**Using JavaScript (fetch):**
```javascript
fetch('http://localhost:8080/fshare-links', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify({
    link: 'https://www.fshare.vn/folder/IYZU2DHZ929T',
    cookie: '__cfduid=...; user=...; pass=...'
  })
})
.then(response => response.json())
.then(data => console.log(data))
.catch(error => console.error('Error:', error));
```

## Project Structure

```
get-folder-linksvip-net/
├── main.go           # Main application file
├── go.mod            # Go module file
├── Dockerfile        # Docker build configuration
├── .dockerignore     # Docker ignore file
└── README.md        # This file
```

## How It Works

1. The API receives a POST request with a Fshare link and cookie
2. Makes an HTTP POST request to linksvip.net's API
3. Cleans the response (removes non-printable characters)
4. Parses the JSON response
5. Formats the response based on the status (`trangthai`)
6. Returns formatted messages with download links

## Error Handling

The application includes:
- Automatic retry logic (up to 5 retries)
- HTTP error handling
- JSON parsing error handling
- Input validation

## Development

### Running in Development Mode

```bash
go run main.go
```

### Building for Production

```bash
go build -ldflags="-s -w" -o app main.go
```

### Testing

Test the API endpoint:
```bash
curl -X POST http://localhost:8080/fshare-links \
  -H "Content-Type: application/json" \
  -d '{"link":"your-link","cookie":"your-cookie"}'
```

## Docker Commands

**View logs:**
```bash
docker logs fshare-api
```

**Stop the container:**
```bash
docker stop fshare-api
```

**Remove the container:**
```bash
docker rm fshare-api
```

**Rebuild and restart:**
```bash
docker build -t get-folder-linksvip-net .
docker stop fshare-api && docker rm fshare-api
docker run -d -p 8080:8080 --name fshare-api get-folder-linksvip-net
```

## License

This project is provided as-is for educational purposes.

## Contributing

Feel free to submit issues and enhancement requests!
