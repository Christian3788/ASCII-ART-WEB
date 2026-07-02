# Docker Quick Reference Guide

## 🚀 Quick Commands

### Build & Run

```bash
# 1. Compile Go binary
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o ascii-art-web .

# 2. Build Docker image
docker build --tag ascii-art-web:1.0 --tag ascii-art-web:latest -f Dockerfile .

# 3. Run container
docker run -d --name ascii-art-web-container -p 8080:8080 ascii-art-web:latest

# 4. Test application
curl http://localhost:8080
```

## 📊 Information Commands

```bash
# View image
docker images | grep ascii-art-web
docker inspect ascii-art-web:1.0

# View container
docker ps | grep ascii-art-web-container
docker inspect ascii-art-web-container

# View logs
docker logs ascii-art-web-container
docker logs -f ascii-art-web-container  # Follow logs

# View metadata
docker inspect ascii-art-web:1.0 --format='{{json .Config.Labels}}' | jq .
docker inspect ascii-art-web-container --format='{{json .Config.Labels}}' | jq .

# Disk usage
docker system df
docker system df -v
```

## 🧹 Garbage Collection Commands

```bash
# Safe cleanup (dangling objects only)
./docker-gc.sh cleanup-safe
# OR
docker system prune -f

# Full cleanup (all unused)
./docker-gc.sh cleanup-full
# OR
docker system prune -a -f

# Daily maintenance
./docker-gc.sh maintenance

# Remove specific objects
docker stop ascii-art-web-container
docker rm ascii-art-web-container
docker rmi ascii-art-web:1.0
docker rmi -f ascii-art-web:1.0  # Force remove
```

## 🐳 Docker Compose Commands

```bash
# Start service
docker-compose up -d

# View status
docker-compose ps

# View logs
docker-compose logs -f

# Stop service
docker-compose down

# Remove containers and volumes
docker-compose down -v
```

## 🔍 Troubleshooting

```bash
# Check if port is in use
lsof -i :8080
docker port ascii-art-web-container

# Get container IP
docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' ascii-art-web-container

# Execute command in container
docker exec ascii-art-web-container ls -la

# Check container resource usage
docker stats ascii-art-web-container

# Inspect network
docker network ls
docker network inspect bridge
```

## 📋 Files Overview

| File | Purpose |
|------|---------|
| `Dockerfile` | Container image definition |
| `.dockerignore` | Build context optimization |
| `docker-compose.yml` | Orchestration configuration |
| `docker-gc.sh` | Garbage collection script |
| `DOCKER.md` | Complete documentation |
| `DOCKER-DELIVERABLES.md` | Project summary |

## ✅ Checklist

- [x] Dockerfile created
- [x] Docker image built (`ascii-art-web:1.0`)
- [x] Docker container running (`ascii-art-web-container`)
- [x] Metadata labels applied
- [x] Garbage collection procedures documented
- [x] Go best practices followed
- [x] Docker best practices implemented
- [x] Application tested and working
- [x] Documentation complete

## 📞 Support

For detailed information, see:
- [DOCKER.md](DOCKER.md) - Complete documentation
- [DOCKER-DELIVERABLES.md](DOCKER-DELIVERABLES.md) - Project summary
- [docker-gc.sh](docker-gc.sh) - Run `./docker-gc.sh help`
