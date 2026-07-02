# Minimal Dockerfile using scratch for production
# This uses a pre-compiled static Go binary
# Build the binary locally: CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ascii-art-web .

FROM scratch

# Add comprehensive metadata labels for container identification
LABEL maintainer="christianotieno" \
      version="1.0" \
      description="ASCII Art Web Application in Docker" \
      org.opencontainers.image.title="ASCII Art Web" \
      org.opencontainers.image.documentation="Generate ASCII art from text input" \
      org.opencontainers.image.authors="Christian Otieno" \
      org.opencontainers.image.licenses="MIT"

# Set working directory
WORKDIR /app

# Copy pre-compiled binary
COPY ascii-art-web .

# Copy static assets (templates and banners)
COPY templates/ ./templates/
COPY banners/ ./banners/

# Expose port for HTTP traffic
EXPOSE 8080

# Set entrypoint to run the application
ENTRYPOINT ["./ascii-art-web"]

