#!/bin/bash

# Docker Garbage Collection and Maintenance Script
# Handles cleanup of unused Docker objects while preserving active resources

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to display Docker resource usage
show_usage() {
    log_info "Current Docker resource usage:"
    docker system df
}

# Function to list running containers with labels
show_containers() {
    log_info "Running containers:"
    docker ps --format "table {{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}"
    
    log_info "\nAll containers (including stopped):"
    docker ps -a --format "table {{.ID}}\t{{.Names}}\t{{.Image}}\t{{.Status}}"
}

# Function to list images with labels
show_images() {
    log_info "Docker images:"
    docker images --format "table {{.Repository}}\t{{.Tag}}\t{{.Size}}\t{{.Created}}"
}

# Function to remove dangling images (orphaned by new builds)
cleanup_dangling_images() {
    log_info "Removing dangling images..."
    
    local count=$(docker images -f "dangling=true" -q | wc -l)
    
    if [ "$count" -gt 0 ]; then
        docker image prune -f
        log_success "Removed $count dangling image(s)"
    else
        log_info "No dangling images found"
    fi
}

# Function to remove stopped containers
cleanup_containers() {
    log_info "Removing stopped containers..."
    
    local count=$(docker ps -a -f "status=exited" -q | wc -l)
    
    if [ "$count" -gt 0 ]; then
        docker container prune -f
        log_success "Removed $count stopped container(s)"
    else
        log_info "No stopped containers found"
    fi
}

# Function to remove unused volumes
cleanup_volumes() {
    log_info "Removing unused volumes..."
    
    local count=$(docker volume ls -f "dangling=true" -q | wc -l)
    
    if [ "$count" -gt 0 ]; then
        docker volume prune -f
        log_success "Removed $count unused volume(s)"
    else
        log_info "No unused volumes found"
    fi
}

# Function to remove build cache
cleanup_build_cache() {
    log_info "Cleaning build cache..."
    
    docker builder prune -f
    log_success "Build cache cleaned"
}

# Function to perform comprehensive cleanup
full_cleanup() {
    log_warning "Starting comprehensive Docker cleanup..."
    
    log_info "Step 1/4: Removing stopped containers..."
    cleanup_containers
    
    log_info "Step 2/4: Removing dangling images..."
    cleanup_dangling_images
    
    log_info "Step 3/4: Removing unused volumes..."
    cleanup_volumes
    
    log_info "Step 4/4: Cleaning build cache..."
    cleanup_build_cache
    
    log_success "Comprehensive cleanup completed"
}

# Function to perform safe cleanup (only dangling objects)
safe_cleanup() {
    log_info "Starting safe cleanup (dangling objects only)..."
    
    docker system prune -f
    log_success "Safe cleanup completed"
}

# Function to stop and remove specific container
stop_container() {
    local container_name=$1
    
    log_info "Stopping container: $container_name"
    
    if docker ps -a --format "{{.Names}}" | grep -q "^${container_name}$"; then
        docker stop "$container_name" 2>/dev/null || log_warning "Container was not running"
        docker rm "$container_name"
        log_success "Container $container_name removed"
    else
        log_warning "Container $container_name not found"
    fi
}

# Function to remove image
remove_image() {
    local image_name=$1
    local force=${2:-false}
    
    log_info "Removing image: $image_name"
    
    if [ "$force" = true ]; then
        docker rmi -f "$image_name"
    else
        docker rmi "$image_name"
    fi
    
    log_success "Image $image_name removed"
}

# Function to inspect container metadata
inspect_container() {
    local container_name=$1
    
    log_info "Container metadata for: $container_name"
    docker inspect "$container_name" --format='{{json .Config.Labels}}' | python3 -m json.tool
}

# Function to inspect image metadata
inspect_image() {
    local image_name=$1
    
    log_info "Image metadata for: $image_name"
    docker inspect "$image_name" --format='{{json .Config.Labels}}' | python3 -m json.tool
}

# Function to display help
show_help() {
    cat << EOF
Docker Garbage Collection and Maintenance Script

USAGE:
    $0 [COMMAND] [OPTIONS]

COMMANDS:
    help                    Show this help message
    status                  Show Docker resource usage
    containers              List all containers
    images                  List all images
    cleanup-dangling        Remove dangling images only
    cleanup-containers      Remove stopped containers
    cleanup-volumes         Remove unused volumes
    cleanup-cache           Remove build cache
    cleanup-safe            Safe cleanup (dangling objects only)
    cleanup-full            Comprehensive cleanup (all unused objects)
    stop <container>        Stop and remove specific container
    remove-image <image>    Remove specific image
    remove-image-f <image>  Force remove specific image
    inspect-container <name>    Show container metadata
    inspect-image <name>        Show image metadata
    maintenance             Run daily maintenance (safe cleanup + status)

EXAMPLES:
    $0 status                              # Show current usage
    $0 cleanup-safe                        # Safe cleanup
    $0 cleanup-full                        # Full cleanup
    $0 stop ascii-art-web-container       # Stop specific container
    $0 maintenance                         # Daily maintenance

EOF
}

# Main script logic
main() {
    if [ $# -eq 0 ]; then
        show_help
        exit 0
    fi
    
    case "${1:-}" in
        help)
            show_help
            ;;
        status)
            show_usage
            ;;
        containers)
            show_containers
            ;;
        images)
            show_images
            ;;
        cleanup-dangling)
            cleanup_dangling_images
            ;;
        cleanup-containers)
            cleanup_containers
            ;;
        cleanup-volumes)
            cleanup_volumes
            ;;
        cleanup-cache)
            cleanup_build_cache
            ;;
        cleanup-safe)
            safe_cleanup
            show_usage
            ;;
        cleanup-full)
            full_cleanup
            show_usage
            ;;
        stop)
            if [ $# -lt 2 ]; then
                log_error "Container name required"
                exit 1
            fi
            stop_container "$2"
            ;;
        remove-image)
            if [ $# -lt 2 ]; then
                log_error "Image name required"
                exit 1
            fi
            remove_image "$2"
            ;;
        remove-image-f)
            if [ $# -lt 2 ]; then
                log_error "Image name required"
                exit 1
            fi
            remove_image "$2" true
            ;;
        inspect-container)
            if [ $# -lt 2 ]; then
                log_error "Container name required"
                exit 1
            fi
            inspect_container "$2"
            ;;
        inspect-image)
            if [ $# -lt 2 ]; then
                log_error "Image name required"
                exit 1
            fi
            inspect_image "$2"
            ;;
        maintenance)
            log_info "Running daily maintenance..."
            safe_cleanup
            log_info ""
            show_usage
            log_success "Daily maintenance completed"
            ;;
        *)
            log_error "Unknown command: $1"
            show_help
            exit 1
            ;;
    esac
}

# Run main function
main "$@"
