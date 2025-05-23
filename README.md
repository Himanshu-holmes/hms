# Hospital Management System 
# Docker Setup Documentation

## Overview

This project uses Docker Compose to orchestrate a multi-container application consisting of:
- **App Service**: Main application running on port 3000
- **Database Service**: PostgreSQL database running on port 5432

## Prerequisites

- Docker installed on your system
- Docker Compose installed on your system
- `wait-for-it.sh` script in your project root (for database connection waiting)

## Project Structure

```
project-root/
├── docker-compose.yml
├── db.env
├── Dockerfile
├── wait-for-it.sh
└── hms (executable)
```

## Configuration Files

### docker-compose.yml

The main orchestration file that defines two services:

- **app**: Your main application service
  - Built from local Dockerfile
  - Exposed on port 3000
  - Waits for database to be ready before starting
  - Connects to PostgreSQL database

- **db**: PostgreSQL database service
  - Uses official PostgreSQL image
  - Exposed on port 5432
  - Configuration loaded from `db.env`

### db.env

Database environment configuration file:
```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=12345
POSTGRES_DB=hms
```

## Setup Instructions

### 1. Ensure Required Files Exist

Make sure you have the following files in your project root:
- `Dockerfile` - Contains instructions to build your app image
- `wait-for-it.sh` - Script to wait for database availability
- `hms` - Your application executable
- `db.env` - Database environment variables
- `app.env` - Application environment variables

### 2. Configure Docker Compose

 if you want to use a different configuration, Edit the `docker-compose.yml`  file to match your project's configuration else leave it as is

### 3. Build and Start Services

```bash
# Build and start all services
docker-compose up --build

# Or run in detached mode
docker-compose up --build -d
```

### 4. Verify Services

Check that both services are running:
```bash
docker-compose ps
```

## Usage Commands

### Start Services
```bash
docker-compose up
```

### Start Services in Background
```bash
docker-compose up -d
```

### Stop Services
```bash
docker-compose down
```

### View Logs
```bash
# All services
docker-compose logs

# Specific service
docker-compose logs app
docker-compose logs db
```

### Rebuild Services
```bash
docker-compose up --build
```

### Remove Everything (including volumes)
```bash
docker-compose down -v
```

## Application Access

- **Application**: http://localhost:3000
- **API Documentation**: http://localhost:3000/swagger/index.html#/
- **Database**: localhost:5432
  - Username: `postgres`
  - Password: `12345`
  - Database: `hms`

## Environment Variables

The application uses the following environment variables:

| Variable | Value | Description |
|----------|-------|-------------|
| `PORT` | 3000 | Application port |
| `DB_URL` | postgresql://postgres:12345@db:5432/hms | Database connection string |

## Database Connection

The application connects to PostgreSQL using:
- Host: `db` (Docker service name)
- Port: `5432`
- Username: `postgres`
- Password: `12345`
- Database: `hms`

## Troubleshooting

### Common Issues

1. **Port conflicts**: If ports 3000 or 5432 are already in use, modify the port mappings in `docker-compose.yml`

2. **Database connection issues**: The `wait-for-it.sh` script ensures the database is ready before starting the app. If you encounter connection issues, verify:
   - Database credentials match between `db.env` and `DB_URL`
   - Database service is healthy: `docker-compose logs db`

3. **Build failures**: Ensure your `Dockerfile` is properly configured and all required files are present

### Useful Debug Commands

```bash
# Check container status
docker-compose ps

# Access database directly
docker-compose exec db psql -U postgres -d hms

# Access app container shell
docker-compose exec app sh

# View real-time logs
docker-compose logs -f
```

## Development Workflow

1. Make code changes
2. Rebuild and restart: `docker-compose up --build`
3. Test your changes at http://localhost:3000
4. Check logs if needed: `docker-compose logs app`

## Production Considerations

For production deployment, consider:
- Using environment-specific configuration files
- Implementing proper secrets management
- Setting up SSL/TLS termination
- Configuring proper logging and monitoring
- Using production-ready PostgreSQL configuration