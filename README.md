# PromeConfig

A modern web application for managing Prometheus configurations with AI-powered alert rule generation.

## Features

- **Target Management**: Add, edit, and delete monitoring targets
- **Alert Rule Management**: Create and manage Prometheus alert rules
- **AI Integration**: Generate alert rules using AI (OpenAI, Claude, etc.)
- **Configuration Preview**: Preview generated Prometheus configuration
- **User Authentication**: Secure login and user management
- **Modern UI**: Clean and responsive interface built with React and Tailwind CSS
- **Production Ready**: Complete deployment solution with Docker, CI/CD, and monitoring

## Tech Stack

### Frontend
- React 18 with TypeScript
- Vite for build tooling
- Tailwind CSS for styling
- Lucide React for icons

### Backend
- Go with Gin framework
- GORM for database operations
- PostgreSQL database
- JWT authentication
- Redis for caching

### Infrastructure
- Docker & Docker Compose
- Nginx reverse proxy
- Prometheus monitoring
- Grafana dashboards
- GitHub Actions CI/CD

## Quick Start

### Development Environment

```bash
# Clone the repository
git clone <repository-url>
cd PromeConfig

# Start development environment
make dev

# In separate terminals:
make dev-backend  # Start backend
make dev-frontend # Start frontend
```

### Production Deployment

```bash
# One-click production deployment
make prod

# Or using deployment scripts
./deploy.sh prod     # Linux/Mac
deploy.bat prod      # Windows
```

## Installation & Setup

### Prerequisites

- Docker 20.10+
- Docker Compose 2.0+
- Node.js 18+ (for development)
- Go 1.21+ (for development)
- PostgreSQL 15+ (for development)

### Development Setup

1. **Clone and install dependencies:**
```bash
git clone <repository-url>
cd PromeConfig
make install
```

2. **Setup environment:**
```bash
# Copy environment file
cp backend/.env.example backend/.env

# Start database
make dev

# Initialize database
make db-init
```

3. **Start development servers:**
```bash
# Backend (Terminal 1)
make dev-backend

# Frontend (Terminal 2)
make dev-frontend
```

4. **Access the application:**
- Frontend: http://localhost:5173
- Backend API: http://localhost:8080
- Database: localhost:5432

### Production Deployment

#### Quick Deployment

```bash
# Setup secrets and SSL certificates
chmod +x scripts/setup-secrets.sh
./scripts/setup-secrets.sh setup

# Deploy to production
make prod
```

#### Manual Deployment

1. **Setup secrets:**
```bash
./scripts/setup-secrets.sh setup
```

2. **Configure environment:**
```bash
# Edit production configuration
vim backend/.env.production
```

3. **Deploy services:**
```bash
# Build and start all services
docker-compose -f docker-compose.prod.yml up -d --build
```

4. **Verify deployment:**
```bash
make health
make status
```

## Available Commands

### Make Commands

```bash
make help           # Show all available commands
make dev            # Start development environment
make prod           # Start production environment
make build          # Build application images
make test           # Run all tests
make lint           # Run code linting
make clean          # Clean containers and images
make logs           # View service logs
make status         # Check service status
make health         # Run health checks
make backup         # Backup database
make deploy         # Full deployment pipeline
```

### Deployment Scripts

```bash
# Linux/Mac
./deploy.sh [dev|prod|stop|logs|clean]

# Windows
deploy.bat [dev|prod|stop|logs|clean]

# Examples
./deploy.sh prod    # Start production
./deploy.sh logs    # View logs
./deploy.sh clean   # Clean environment
```

### Secret Management

```bash
./scripts/setup-secrets.sh setup   # Generate all secrets
./scripts/setup-secrets.sh backup  # Backup secrets
./scripts/setup-secrets.sh info    # Show secret info
./scripts/setup-secrets.sh clean   # Clean secrets
```

## Configuration

### Environment Variables

**Development (.env):**
```env
ENV=development
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres123
DB_NAME=promeconfig
JWT_SECRET=your-dev-secret
PROMETHEUS_URL=http://localhost:9090
```

**Production (.env.production):**
```env
ENV=production
DB_HOST=postgres
DB_PASSWORD_FILE=/run/secrets/db_password
JWT_SECRET_FILE=/run/secrets/jwt_secret
REDIS_URL=redis://redis:6379
LOG_LEVEL=info
```

### Service Configuration

- **Nginx**: `deploy/nginx.prod.conf`
- **Redis**: `deploy/redis.conf`
- **Prometheus**: `deploy/prometheus.yml`
- **Docker Compose**: `docker-compose.yml` (dev), `docker-compose.prod.yml` (prod)

## API Documentation

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - User login
- `POST /api/auth/logout` - User logout

### Targets
- `GET /api/targets` - Get all targets
- `POST /api/targets` - Create new target
- `PUT /api/targets/:id` - Update target
- `DELETE /api/targets/:id` - Delete target

### Alert Rules
- `GET /api/alertrules` - Get all alert rules
- `POST /api/alertrules` - Create new alert rule
- `PUT /api/alertrules/:id` - Update alert rule
- `DELETE /api/alertrules/:id` - Delete alert rule

### AI Settings
- `GET /api/aisettings` - Get AI settings
- `POST /api/aisettings` - Create AI settings
- `PUT /api/aisettings/:id` - Update AI settings
- `DELETE /api/aisettings/:id` - Delete AI settings

### System
- `GET /api/health` - Health check
- `GET /api/metrics` - Prometheus metrics

## Monitoring & Observability

### Built-in Monitoring

- **Prometheus**: Metrics collection (http://localhost:9090)
- **Grafana**: Dashboards and visualization (http://localhost:3000)
- **Health Checks**: Automated service health monitoring
- **Logging**: Structured JSON logging with log rotation

### Accessing Monitoring

```bash
# Start with monitoring stack
docker-compose -f docker-compose.prod.yml --profile monitoring up -d

# View metrics
curl http://localhost:8080/api/metrics

# Check health
curl http://localhost:8080/api/health
```

## Security

### Production Security Features

- **HTTPS/TLS**: SSL termination at Nginx
- **Secrets Management**: Docker secrets for sensitive data
- **Non-root Containers**: All services run as non-root users
- **Network Isolation**: Docker networks with restricted access
- **Rate Limiting**: API rate limiting and DDoS protection
- **Security Headers**: HSTS, CSP, and other security headers
- **Input Validation**: Comprehensive input validation and sanitization

### Security Best Practices

1. **Change default passwords**: Update all default credentials
2. **Use strong secrets**: Generate cryptographically secure secrets
3. **Enable HTTPS**: Use valid SSL certificates in production
4. **Regular updates**: Keep dependencies and base images updated
5. **Monitor logs**: Set up log monitoring and alerting
6. **Backup strategy**: Implement regular database backups

## Troubleshooting

### Common Issues

1. **Port conflicts:**
```bash
# Check port usage
netstat -tulpn | grep :80
netstat -tulpn | grep :8080
```

2. **Container startup issues:**
```bash
# View detailed logs
make logs
docker-compose logs app
```

3. **Database connection issues:**
```bash
# Check database status
make shell-db
psql -U postgres -c "SELECT version();"
```

4. **Permission issues:**
```bash
# Fix script permissions
chmod +x deploy.sh scripts/setup-secrets.sh
```

### Debug Mode

```bash
# Start in debug mode
docker-compose up --build

# Access container shells
make shell-backend
make shell-db

# View real-time logs
make logs
```

## CI/CD Pipeline

### GitHub Actions

The project includes a complete CI/CD pipeline:

- **Code Quality**: Linting, testing, security scanning
- **Build & Push**: Multi-platform Docker image builds
- **Deployment**: Automated deployment to dev/prod environments
- **Monitoring**: Health checks and notifications

### Pipeline Stages

1. **Lint & Test**: Code quality checks and unit tests
2. **Build & Push**: Docker image build and registry push
3. **Security Scan**: Vulnerability scanning with Trivy
4. **Deploy**: Environment-specific deployments
5. **Monitor**: Health checks and notifications

## Performance

### Optimization Features

- **Multi-stage builds**: Optimized Docker images
- **Static asset caching**: Long-term browser caching
- **Gzip compression**: Response compression
- **Connection pooling**: Database connection optimization
- **Redis caching**: Application-level caching
- **Load balancing**: Multi-instance deployment support

### Performance Monitoring

```bash
# View performance metrics
curl http://localhost:8080/api/metrics

# Monitor resource usage
docker stats

# Check response times
curl -w "@curl-format.txt" -o /dev/null -s http://localhost/api/health
```

## Backup & Recovery

### Database Backup

```bash
# Create backup
make backup

# Restore from backup
make restore FILE=backups/backup_20240101_120000.sql

# Automated backups (add to crontab)
0 2 * * * cd /path/to/promeconfig && make backup
```

### Full System Backup

```bash
# Backup secrets
./scripts/setup-secrets.sh backup

# Backup data volumes
docker run --rm -v promeconfig_postgres_data:/data -v $(pwd)/backups:/backup alpine tar czf /backup/postgres_data.tar.gz -C /data .
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Run tests (`make test`)
4. Run linting (`make lint`)
5. Commit your changes (`git commit -m 'Add some amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Development Guidelines

- Follow existing code style and conventions
- Add tests for new features
- Update documentation as needed
- Ensure all CI checks pass

## Documentation

- [Deployment Guide](DEPLOYMENT.md) - Detailed deployment instructions
- [Backend README](backend/README.md) - Backend-specific documentation
- [API Documentation](docs/api.md) - Complete API reference
- [Architecture Guide](docs/architecture.md) - System architecture overview

## Support

If you encounter any issues:

1. Check the [troubleshooting section](#troubleshooting)
2. Review the logs: `make logs`
3. Check service status: `make status`
4. Run health checks: `make health`
5. Open an issue on GitHub

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**PromeConfig** - Making Prometheus configuration management simple and intelligent.