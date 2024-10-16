# Use Alpine Linux as the base image
FROM alpine:latest

# Install SQLite
RUN apk add --no-cache sqlite

# Set the working directory to /app
WORKDIR /app

# Command to start msqlite with the mounted database on container start
CMD ["sqlite3", "/sample.db"]