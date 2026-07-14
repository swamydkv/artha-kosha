#!/bin/bash

# ArthaKosha Authentication MVP Deployment Script for Docker Desktop
# This script helps deploy, manage, and test the authentication system

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
COMPOSE_FILE="infra/docker-compose.yml"
PROJECT_NAME="artha-kosha"
API_PORT=8080
WEB_PORT=3000
DB_PORT=5432

# Functions
print_header() {
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

print_error() {
    echo -e "${RED}✗ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed. Please install Docker Desktop first."
        exit 1
    fi

    # Check if Docker is running
    if ! docker info &> /dev/null; then
        print_error "Docker is not running. Please start Docker Desktop."
        exit 1
    fi

    print_success "Docker is running"
}

check_prerequisites() {
    print_header "Checking Prerequisites"
    
    check_docker
    
    # Check if required files exist
    if [ ! -f "$COMPOSE_FILE" ]; then
        print_error "Docker Compose file not found: $COMPOSE_FILE"
        exit 1
    fi
    print_success "Docker Compose file found"
    
    if [ ! -d "apps/finance-api" ]; then
        print_error "Backend directory not found: apps/finance-api"
        exit 1
    fi
    print_success "Backend directory found"
    
    if [ ! -d "apps/web" ]; then
        print_error "Frontend directory not found: apps/web"
        exit 1
    fi
    print_success "Frontend directory found"
    
    # Check if go.mod exists
    if [ ! -f "apps/finance-api/go.mod" ]; then
        print_warning "go.mod not found - Go dependencies may not be initialized"
    else
        print_success "Go module found"
    fi
    
    # Check if package.json exists
    if [ ! -f "apps/web/package.json" ]; then
        print_warning "package.json not found - Node dependencies may not be initialized"
    else
        print_success "Node package found"
    fi
}

start_services() {
    print_header "Starting Services in Docker Desktop"
    
    # Build and start services
    print_info "Building and starting PostgreSQL, API, and Web services..."
    docker compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" up -d --build
    
    print_success "Services started in Docker Desktop"
    
    # Wait for services to be healthy
    print_header "Waiting for Services to be Ready"
    print_info "Waiting for PostgreSQL to initialize..."
    sleep 15
    
    # Check service status
    print_header "Service Status"
    docker compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" ps
    
    # Wait additional time for API to compile and start
    print_info "Waiting for API to compile and start..."
    sleep 20
    
    # Check PostgreSQL
    if docker compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" ps postgres | grep -q "Up"; then
        print_success "PostgreSQL is running"
    else
        print_error "PostgreSQL failed to start"
        docker compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" logs postgres
        exit 1
    fi
    
    # Check API health
    print_info "Checking API health..."
    MAX_RETRIES=10
    RETRY_COUNT=0
    
    while [ $RETRY_COUNT -lt $MAX_RETRIES ]; do
        if curl -s "http://localhost:$API_PORT/health" > /dev/null 2>&1; then
            print_success "API is healthy and responding"
            break
        else
            RETRY_COUNT=$((RETRY_COUNT + 1))
            print_info "API not ready yet... retrying ($RETRY_COUNT/$MAX_RETRIES)"
            sleep 5
        fi
    done
    
    if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
        print_warning "API health check failed - checking logs"
        docker compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" logs finance-api
    fi
    
    # Check Web
    print_info "Checking Web accessibility..."
    if curl -s "http://localhost:$WEB_PORT" > /dev/null 2>&1; then
        print_success "Web is accessible"
    else
        print_warning "Web may still be starting - check logs if needed"
    fi
    
    print_header "Service Access Information"
    echo -e "🌐 Web Application: ${GREEN}http://localhost:$WEB_PORT${NC}"
    echo -e "🔧 API Endpoints: ${GREEN}http://localhost:$API_PORT${NC}"
    echo -e "🗄️  Database: ${GREEN}localhost:$DB_PORT${NC}"
    echo -e ""
    echo -e "API Endpoints:"
    echo -e "  - POST ${GREEN}http://localhost:$API_PORT/register${NC}"
    echo -e "  - POST ${GREEN}http://localhost:$API_PORT/login${NC}"
    echo -e "  - POST ${GREEN}http://localhost:$API_PORT/logout${NC}"
    echo -e "  - GET  ${GREEN}http://localhost:$API_PORT/health${NC}"
    echo -e ""
    echo -e "Database Connection:"
    echo -e "  - Host: ${GREEN}localhost${NC}"
    echo -e "  - Port: ${GREEN}$DB_PORT${NC}"
    echo -e "  - Database: ${GREEN}artha_kosha${NC}"
    echo -e "  - User: ${GREEN}postgres${NC}"
    echo -e "  - Password: ${GREEN}postgres${NC}"
}

stop_services() {
    print_header "Stopping Services"
    
    docker compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" down
    
    print_success "Services stopped"
}

clean_data() {
    print_header "Cleaning All Data"
    
    print_warning "This will stop all services and remove all data including database volumes"
    read -p "Are you sure? (y/N): " -n 1 -r
    echo
    
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        # Stop services and remove volumes
        docker compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" down -v
        
        # Remove any orphaned containers
        docker compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" down --remove-orphans
        
        print_success "All data cleaned - volumes removed"
    else
        print_info "Clean data operation cancelled"
    fi
}

restart_services() {
    print_header "Restarting Services"
    stop_services
    sleep 3
    start_services
}

show_logs() {
    local service=$1
    
    if [ -z "$service" ]; then
        print_header "Showing All Service Logs"
        docker compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" logs -f
    else
        print_header "Showing Logs for $service"
        docker compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" logs -f "$service"
    fi
}

show_status() {
    print_header "Service Status"
    
    docker compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" ps
    
    echo ""
    print_header "Container Details"
    docker ps --filter "name=$PROJECT_NAME" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
}

run_api_tests() {
    print_header "Running API Integration Tests"
    
    # Check if API is running
    if ! curl -s "http://localhost:$API_PORT/health" > /dev/null; then
        print_error "API is not running. Start services first with: ./deploy.sh start"
        return 1
    fi
    
    print_success "API is running - starting tests"
    
    # Test health endpoint
    echo ""
    print_info "Test 1: Health Endpoint"
    HEALTH_RESPONSE=$(curl -s "http://localhost:$API_PORT/health")
    if [ "$HEALTH_RESPONSE" == "ok" ]; then
        print_success "Health endpoint working: $HEALTH_RESPONSE"
    else
        print_error "Health endpoint failed: $HEALTH_RESPONSE"
    fi
    
    # Test registration
    echo ""
    print_info "Test 2: Registration Endpoint"
    REGISTER_RESPONSE=$(curl -s -X POST "http://localhost:$API_PORT/register" \
        -H "Content-Type: application/json" \
        -d '{
            "full_name": "Test User",
            "date_of_birth": "1990-01-01",
            "mobile_number": "+1234567890",
            "email": "test@example.com",
            "username": "testuser_'$(date +%s)'",
            "password": "TestPass123!",
            "confirm_password": "TestPass123!"
        }')
    
    if echo "$REGISTER_RESPONSE" | grep -q "user_id"; then
        print_success "Registration endpoint working"
        USER_ID=$(echo "$REGISTER_RESPONSE" | grep -o '"user_id":"[^"]*"' | cut -d'"' -f4)
        USERNAME=$(echo "$REGISTER_RESPONSE" | grep -o '"username":"[^"]*"' | cut -d'"' -f4)
        echo "  Created user - ID: $USER_ID, Username: $USERNAME"
    else
        print_error "Registration endpoint failed"
        echo "  Response: $REGISTER_RESPONSE"
    fi
    
    # Test duplicate registration
    echo ""
    print_info "Test 3: Duplicate Registration (should fail)"
    DUPLICATE_RESPONSE=$(curl -s -X POST "http://localhost:$API_PORT/register" \
        -H "Content-Type: application/json" \
        -d '{
            "full_name": "Test User",
            "date_of_birth": "1990-01-01",
            "mobile_number": "+1234567890",
            "email": "test@example.com",
            "username": "testuser_'$(date +%s)'",
            "password": "TestPass123!",
            "confirm_password": "TestPass123!"
        }')
    
    if echo "$DUPLICATE_RESPONSE" | grep -q "error\|duplicate"; then
        print_success "Duplicate registration properly rejected"
    else
        print_warning "Duplicate registration response unexpected: $DUPLICATE_RESPONSE"
    fi
    
    # Test login with the created user
    echo ""
    print_info "Test 4: Login Endpoint"
    LOGIN_RESPONSE=$(curl -s -X POST "http://localhost:$API_PORT/login" \
        -H "Content-Type: application/json" \
        -d "{
            \"username\": \"$USERNAME\",
            \"password\": \"TestPass123!\"
        }")
    
    if echo "$LOGIN_RESPONSE" | grep -q "session_id"; then
        print_success "Login endpoint working"
        SESSION_ID=$(echo "$LOGIN_RESPONSE" | grep -o '"session_id":"[^"]*"' | cut -d'"' -f4)
        WELCOME_MSG=$(echo "$LOGIN_RESPONSE" | grep -o '"welcome_message":"[^"]*"' | cut -d'"' -f4)
        echo "  Created session - ID: $SESSION_ID"
        echo "  Welcome message: $WELCOME_MSG"
    else
        print_error "Login endpoint failed"
        echo "  Response: $LOGIN_RESPONSE"
    fi
    
    # Test logout
    echo ""
    print_info "Test 5: Logout Endpoint"
    LOGOUT_RESPONSE=$(curl -s -X POST "http://localhost:$API_PORT/logout" \
        -H "Content-Type: application/json" \
        -H "X-Session-ID: $SESSION_ID")
    
    if echo "$LOGOUT_RESPONSE" | grep -q "logged out\|success"; then
        print_success "Logout endpoint working"
    else
        print_error "Logout endpoint failed"
        echo "  Response: $LOGOUT_RESPONSE"
    fi
    
    # Test invalid login
    echo ""
    print_info "Test 6: Invalid Login (should fail)"
    INVALID_LOGIN_RESPONSE=$(curl -s -X POST "http://localhost:$API_PORT/login" \
        -H "Content-Type: application/json" \
        -d '{
            "username": "wronguser",
            "password": "wrongpass"
        }')
    
    if echo "$INVALID_LOGIN_RESPONSE" | grep -q "error\|invalid"; then
        print_success "Invalid login properly rejected"
    else
        print_warning "Invalid login response unexpected: $INVALID_LOGIN_RESPONSE"
    fi
    
    print_header "API Tests Completed"
}

run_backend_tests() {
    print_header "Running Backend Unit Tests"
    
    if [ ! -d "apps/finance-api" ]; then
        print_error "Backend directory not found"
        return 1
    fi
    
    cd apps/finance-api
    
    if [ ! -f "go.mod" ]; then
        print_error "Go module not found - cannot run tests"
        cd ../..
        return 1
    fi
    
    print_info "Running Go tests..."
    if go test -v ./...; then
        print_success "Backend tests passed"
    else
        print_error "Backend tests failed"
        cd ../..
        return 1
    fi
    
    cd ../..
}

setup_database() {
    print_header "Database Setup"
    
    # Check if PostgreSQL is running
    if ! docker compose -f "$COMPOSE_FILE" -p "$PROJECT_NAME" ps postgres | grep -q "Up"; then
        print_error "PostgreSQL is not running. Start services first with: ./deploy.sh start"
        return 1
    fi
    
    print_success "PostgreSQL is running"
    
    # Wait for PostgreSQL to be fully ready
    print_info "Waiting for PostgreSQL to accept connections..."
    sleep 5
    
    # Check if migration files exist
    if [ -d "apps/finance-api/migrations" ]; then
        print_header "Available Migration Files"
        ls -la apps/finance-api/migrations/
        
        print_info "To apply migrations manually, you can:"
        echo "  docker exec -i ${PROJECT_NAME}-postgres-1 psql -U postgres -d artha_kosha < apps/finance-api/migrations/001_create_users_table.sql"
        echo "  docker exec -i ${PROJECT_NAME}-postgres-1 psql -U postgres -d artha_kosha < apps/finance-api/migrations/002_create_sessions_table.sql"
        
        # Ask if user wants to apply migrations
        read -p "Do you want to apply the migrations now? (y/N): " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            print_info "Applying migrations..."
            
            for migration_file in apps/finance-api/migrations/*.sql; do
                if [ -f "$migration_file" ]; then
                    print_info "Applying $migration_file"
                    docker exec -i ${PROJECT_NAME}-postgres-1 psql -U postgres -d artha_kosha < "$migration_file"
                    if [ $? -eq 0 ]; then
                        print_success "Migration applied: $(basename $migration_file)"
                    else
                        print_error "Migration failed: $(basename $migration_file)"
                    fi
                fi
            done
            
            print_success "Database setup completed"
        else
            print_info "Migrations not applied"
        fi
    else
        print_warning "No migration files found in apps/finance-api/migrations/"
    fi
}

open_browser() {
    print_header "Opening Application in Browser"
    
    if command -v open &> /dev/null; then
        # macOS
        open "http://localhost:$WEB_PORT"
    elif command -v xdg-open &> /dev/null; then
        # Linux
        xdg-open "http://localhost:$WEB_PORT"
    elif command -v start &> /dev/null; then
        # Windows
        start "http://localhost:$WEB_PORT"
    else
        print_info "Please open http://localhost:$WEB_PORT in your browser"
    fi
}

show_help() {
    cat << EOF
ArthaKosha Authentication MVP Deployment Script for Docker Desktop

Usage: ./deploy.sh [COMMAND]

Commands:
    start           Start all services in Docker Desktop (PostgreSQL, API, Web)
    stop            Stop all services
    restart         Restart all services
    status          Show detailed status of all services
    logs            Show logs for all services (follow mode)
    logs [service]  Show logs for specific service (postgres, finance-api, web)
    clean           Stop services and remove all data including volumes
    test-api        Run API integration tests
    test-backend    Run backend unit tests
    setup-db        Setup database (apply migrations)
    open            Open web application in browser
    help            Show this help message

Examples:
    ./deploy.sh start              # Start all services
    ./deploy.sh logs finance-api    # Show API logs
    ./deploy.sh test-api            # Test API endpoints
    ./deploy.sh clean               # Remove all data
    ./deploy.sh open                # Open application in browser

Testing Workflow:
1. Start services:     ./deploy.sh start
2. Open application:   ./deploy.sh open
3. Test registration manually in browser
4. Test login manually in browser
5. Run API tests:      ./deploy.sh test-api
6. View logs:         ./deploy.sh logs
7. Clean data:         ./deploy.sh clean

Service Access:
- Web: http://localhost:3000
- API: http://localhost:8080
- Database: localhost:5432

Database Connection Info:
- Host: localhost
- Port: 5432
- Database: artha_kosha
- User: postgres
- Password: postgres

EOF
}

# Main script logic
case "${1:-help}" in
    start)
        check_prerequisites
        start_services
        ;;
    stop)
        stop_services
        ;;
    restart)
        check_prerequisites
        restart_services
        ;;
    status)
        show_status
        ;;
    logs)
        show_logs "$2"
        ;;
    clean)
        clean_data
        ;;
    test-api)
        run_api_tests
        ;;
    test-backend)
        run_backend_tests
        ;;
    setup-db)
        setup_database
        ;;
    open)
        open_browser
        ;;
    help|--help|-h)
        show_help
        ;;
    *)
        print_error "Unknown command: $1"
        show_help
        exit 1
        ;;
esac