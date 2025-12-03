FROM golang:1.23-alpine

# Install Python, pip, and ffmpeg
RUN apk add --no-cache python3 py3-pip ffmpeg

# Install yt-dlp
RUN pip3 install yt-dlp --break-system-packages

WORKDIR /app

# Copy Go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go application
RUN go build -o main .

# Create downloads directory
RUN mkdir -p downloads

# Expose port
EXPOSE 3000

# Run the application
CMD ["./main"]
