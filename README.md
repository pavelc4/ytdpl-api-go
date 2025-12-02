# ytdpl-api-go

A high-performance, secure, and robust Go API wrapper for `yt-dlp`. This API allows you to extract video information, download URLs, and formats from various platforms supported by `yt-dlp` without downloading the actual files.

## 🚀 Features

- **Link Extraction Only**: Strictly configured to extract links and metadata only (`--no-playlist`), ensuring no large files are downloaded to the server.
- **High Performance**:
  - **In-Memory Caching**: Caches results for 15 minutes to provide instant responses for repeated requests.
  - **Concurrency Control**: Limits concurrent `yt-dlp` processes (default: 10) to prevent server overload.
  - **Response Compression**: Uses Gzip/Brotli compression for faster data transfer.
- **Security**:
  - **Rate Limiting**: Protects the API from abuse (default: 20 requests/minute per IP).
  - **Input Validation**: Validates URLs before processing.
  - **Context Cancellation**: Automatically kills `yt-dlp` processes if the client disconnects.
- **Clean API**: Returns a standardized, RESTful JSON response structure.

## 🛠️ Prerequisites

- **Go**: Version 1.20 or higher.
- **yt-dlp**: Must be installed and available in the system PATH.
  ```bash
  # Linux / macOS
  sudo curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
  sudo chmod a+rx /usr/local/bin/yt-dlp
  ```

##  Installation

1.  **Clone the repository**
    ```bash
    git clone https://github.com/pavelc4/ytdpl-api-go.git
    cd ytdpl-api-go
    ```

2.  **Install dependencies**
    ```bash
    go mod tidy
    ```

3.  **Configuration**
    Copy `.env.example` to `.env` and adjust settings if needed.
    ```bash
    cp .env.example .env
    ```

4.  **Run the server**
    ```bash
    go run cmd/server/main.go
    ```

##  API Endpoints

### 1. Get Download URLs
Extracts direct video and audio download links.

- **URL**: `/api/v1/dl`
- **Method**: `GET`
- **Query Params**: `url` (required)
- **Example**:
  ```bash
  curl "http://localhost:3000/api/v1/dl?url=https://www.youtube.com/watch?v=dQw4w9WgXcQ"
  ```

### 2. Get Video Info
Retrieves detailed metadata about the video.

- **URL**: `/api/v1/info`
- **Method**: `GET`
- **Query Params**: `url` (required)
- **Example**:
  ```bash
  curl "http://localhost:3000/api/v1/info?url=https://www.youtube.com/watch?v=dQw4w9WgXcQ"
  ```

### 3. Get Formats
Lists all available video and audio formats.

- **URL**: `/api/v1/formats`
- **Method**: `GET`
- **Query Params**: `url` (required)
- **Example**:
  ```bash
  curl "http://localhost:3000/api/v1/formats?url=https://www.youtube.com/watch?v=dQw4w9WgXcQ"
  ```

### 4. Health Check
Checks the API status and `yt-dlp` availability.

- **URL**: `/health`
- **Method**: `GET`

##  Response Structure

The API uses a standardized JSON envelope:

**Success Response:**
```json
{
  "data": {
    "video_url": "https://...",
    "audio_url": "https://..."
  },
  "meta": {
    "request_id": "...",
    "timestamp": 1701540000,
    "version": "1.0"
  }
}
```

**Error Response:**
```json
{
  "error": {
    "code": "EXTRACTION_FAILED",
    "message": "Failed to extract download URLs",
    "details": "..."
  }
}
```

##  Project Structure

```
.
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── handler/         # HTTP request handlers
│   ├── models/          # Data structures
│   ├── routes/          # Route definitions & middleware
│   └── services/        # Business logic (yt-dlp wrapper)
├── config/              # Configuration loader
├── .env                 # Environment variables
└── README.md            # Project documentation
```
