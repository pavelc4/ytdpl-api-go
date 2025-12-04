# ytdpl-api-go

A high-performance, secure, and robust Go API wrapper for `yt-dlp`. This API allows you to extract video information, download URLs, and formats from various platforms supported by `yt-dlp` without downloading the actual files.

## 🚀 Features

- **Cloudflare R2 Integration**:
  - **Direct Upload**: Downloads videos/audio and uploads them directly to Cloudflare R2.
  - **Auto-Cleanup**: Automatically deletes files from R2 older than 7 days.
  - **Storage Separation**: Organizes files into `vidioe/` and `audio/` folders.
- **High Performance**:
  - **Smart Caching**: Caches R2 upload results for 1 hour to prevent redundant processing.
  - **In-Memory Caching**: Caches metadata results for 15 minutes.
  - **Concurrency Control**: Limits concurrent `yt-dlp` processes (default: 10).
  - **Response Compression**: Uses Gzip/Brotli compression.
- **Security**:
  - **Rate Limiting**: 
    - Global: 20 requests/minute per IP.
    - Upload/Merge: 5 requests/minute per IP.
  - **Input Validation**: Validates URLs before processing.
  - **Context Cancellation**: Automatically kills `yt-dlp` processes if the client disconnects.
- **Containerless Ready**: Supports loading YouTube cookies directly from R2 storage.
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
    Copy `.env.example` to `.env` and configure the following:
    ```bash
    cp .env.example .env
    ```
    
    **Required for R2 Uploads:**
    - `R2_ACCOUNT_ID`
    - `R2_ACCESS_KEY_ID`
    - `R2_SECRET_ACCESS_KEY`
    - `R2_BUCKET_NAME`
    - `R2_ENDPOINT`
    - `R2_PUBLIC_URL`

    **Optional:**
    - `R2_COOKIE_KEY`: Path to cookie file in R2 bucket (e.g., `cookies/youtube.txt`) for containerless deployments.

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

### 4. Merge & Upload (R2)
Downloads the video/audio, processes it, and uploads it to Cloudflare R2.

- **URL**: `/api/v1/merge`
- **Method**: `GET`
- **Query Params**:
  - `url` (required)
  - `quality` (optional): `best` (default), `1080p`, `720p`.
  - `type` (optional): `video` (default), `audio`.
- **Examples**:

  **Best Quality Video (Default):**
  ```bash
  curl "http://localhost:3000/api/v1/merge?url=https://youtu.be/..."
  ```

  **Audio Only (MP3):**
  ```bash
  curl "http://localhost:3000/api/v1/merge?url=https://youtu.be/...&type=audio"
  ```

  **Specific Resolution (1080p):**
  ```bash
  curl "http://localhost:3000/api/v1/merge?url=https://youtu.be/...&quality=1080p"
  ```

### 5. Health Check
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
