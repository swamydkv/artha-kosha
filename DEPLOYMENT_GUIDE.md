# ArthaKosha Authentication MVP - Deployment Guide

## Quick Start

### Prerequisites
- Docker Desktop installed and running
- Git (for cloning the repository)
- Bash shell (macOS/Linux) or Git Bash (Windows)

### Deployment Steps

1. **Start Docker Desktop**
   - Make sure Docker Desktop is running before executing the deployment script

2. **Navigate to Project Directory**
   ```bash
   cd /Users/swamydkv/Desktop/myProjects/artha-kosha
   ```

3. **Start All Services**
   ```bash
   ./deploy.sh start
   ```

   This will:
   - Build and start PostgreSQL database
   - Build and start the Go API server
   - Build and start the Next.js web application
   - Wait for services to be healthy
   - Display service access information

4. **Open the Application**
   ```bash
   ./deploy.sh open
   ```
   Or manually open: http://localhost:3000

5. **Test the Application**
   - Click "Create Account" to register a new user
   - Fill in the registration form with valid data
   - After successful registration, you'll be redirected to login
   - Sign in with your credentials
   - You should see the authenticated home page with your first name
   - Click "Logout" to end the session

## Deployment Script Commands

### Service Management
```bash
./deploy.sh start       # Start all services
./deploy.sh stop        # Stop all services
./deploy.sh restart     # Restart all services
./deploy.sh status      # Show service status
./deploy.sh logs        # Show all logs
./deploy.sh logs finance-api  # Show specific service logs
```

### Testing
```bash
./deploy.sh test-api    # Run API integration tests
./deploy.sh test-backend    # Run backend unit tests
```

### Database Management
```bash
./deploy.sh setup-db    # Apply database migrations
```

### Data Management
```bash
./deploy.sh clean       # Remove all data (requires confirmation)
```

### Utilities
```bash
./deploy.sh open        # Open web app in browser
./deploy.sh help        # Show help message
```

## Service Access Information

### Web Application
- **URL**: http://localhost:3000
- **Features**: Landing page, registration, login, authenticated home

### API Endpoints
- **Base URL**: http://localhost:8080
- **Endpoints**:
  - `POST /register` - User registration
  - `POST /login` - User authentication
  - `POST /logout` - Session termination
  - `GET /health` - Health check

### Database
- **Host**: localhost
- **Port**: 5432
- **Database**: artha_kosha
- **User**: postgres
- **Password**: postgres

## Manual Testing

### Registration Flow
1. Navigate to http://localhost:3000
2. Click "Create Account"
3. Fill in the form:
   - Full Name: John Doe
   - Date of Birth: 1990-01-15
   - Mobile Number: +1234567890
   - Email: john@example.com
   - Username: johndoe
   - Password: SecurePass123!
   - Confirm Password: SecurePass123!
4. Click "Create Account"
5. You should be redirected to the login page

### Login Flow
1. Navigate to http://localhost:3000/login
2. Enter your username and password
3. Click "Sign In"
4. You should see the authenticated home page with welcome message

### Logout Flow
1. From the authenticated home page
2. Click "Logout"
3. You should be redirected to the login page

## API Testing with Curl

### Health Check
```bash
curl http://localhost:8080/health
```

### Register User
```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Test User",
    "date_of_birth": "1990-01-01",
    "mobile_number": "+1234567890",
    "email": "test@example.com",
    "username": "testuser",
    "password": "TestPass123!",
    "confirm_password": "TestPass123!"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "TestPass123!"
  }'
```

### Logout
```bash
curl -X POST http://localhost:8080/logout \
  -H "Content-Type: application/json" \
  -H "X-Session-ID: YOUR_SESSION_ID"
```

## Troubleshooting

### Services Won't Start
1. Check Docker Desktop is running
2. Check port conflicts (8080, 3000, 5432)
3. View logs: `./deploy.sh logs finance-api`
4. Try clean restart: `./deploy.sh clean && ./deploy.sh start`

### API Not Responding
1. Check API logs: `./deploy.sh logs finance-api`
2. Verify API is built correctly
3. Check if port 8080 is available
4. Restart services: `./deploy.sh restart`

### Database Connection Issues
1. Check PostgreSQL logs: `./deploy.sh logs postgres`
2. Verify database is created: `./deploy.sh setup-db`
3. Check database credentials in docker-compose.yml

### Frontend Not Loading
1. Check web logs: `./deploy.sh logs web`
2. Verify Node.js dependencies are installed
3. Check if port 3000 is available
4. Try accessing API directly to isolate the issue

### Clean Reset
If you encounter persistent issues:
```bash
./deploy.sh clean
./deploy.sh start
```

This will remove all data and start fresh.

## Development Workflow

### Making Changes
1. Make changes to the code
2. Restart services: `./deploy.sh restart`
3. Test your changes
4. Run tests: `./deploy.sh test-api`

### Viewing Logs
```bash
# All logs
./deploy.sh logs

# Specific service
./deploy.sh logs finance-api
./deploy.sh logs postgres
./deploy.sh logs web
```

### Database Access
```bash
# Connect to PostgreSQL
docker exec -it artha-kosha-postgres-1 psql -U postgres -d artha_kosha

# View tables
\dt

# View users
SELECT * FROM users;

# Exit
\q
```

## Performance Testing

The system includes performance benchmarks that can be run:

```bash
# Run Go benchmarks
cd apps/finance-api
go test -bench=. -benchmem ./internal/auth/integration/
```

Expected performance targets:
- Registration: < 2 seconds
- Login: < 1 second
- Logout: < 500ms

## Stopping Services

When you're done:
```bash
./deploy.sh stop
```

To completely remove all data:
```bash
./deploy.sh clean
```

## Architecture Overview

The deployment consists of three main services:

1. **PostgreSQL Database**
   - Stores user accounts and sessions
   - Persistent data storage
   - Accessible on port 5432

2. **Go API Server**
   - REST API for authentication
   - Business logic and validation
   - Accessible on port 8080

3. **Next.js Web Application**
   - User interface for registration/login
   - Client-side validation
   - Accessible on port 3000

All services run in Docker Desktop containers and are orchestrated via Docker Compose.

## Security Notes

This is a development environment with the following security considerations:

- Passwords are hashed (currently SHA-256, to be upgraded to Argon2id)
- Session management uses localStorage (for MVP)
- No HTTPS in development
- Default database credentials (postgres/postgres)

For production deployment, additional security measures would be required:
- Argon2id password hashing
- Secure session management (HttpOnly cookies)
- HTTPS/TLS encryption
- Environment-based configuration
- Secrets management (HashiCorp Vault)