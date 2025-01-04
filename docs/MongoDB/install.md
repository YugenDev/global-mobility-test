# MongoDB Installation

## Prerequisites
- Docker installed
- Docker Compose installed
- Git (optional)

## Quick Start
Start MongoDB container:
```bash
docker-compose up mongodb -d
```

## Configuration
MongoDB is configured with:
- Port: 27017
- Username: admin
- Password: globalMobility
- Database: ecommerce
- Volumes:
  - mongodb_data:/data/db
  - mongodb_config:/data/configdb

## Connection String
```
mongodb://admin:globalMobility@localhost:27017
```

## Verify Installation
1. Check container status:
```bash
docker ps | grep mongodb
```

2. Access MongoDB shell:
```bash
docker exec -it mongodb mongosh -u admin -p globalMobility
```

## Common Issues
1. Port already in use:
   - Stop local MongoDB service
   - Change port in 

docker-compose.yml



2. Authentication failed:
   - Verify environment variables
   - Check credentials in 

docker-compose.yml



## Data Persistence
Data is persisted through Docker volumes:
- mongodb_data: Stores database files
- mongodb_config: Stores configuration files

## Security Notes
- Change default credentials in production
- Use environment variables for sensitive data
- Configure network security appropriately

## Cleanup
To remove containers and volumes:
```bash
docker-compose down -v
```