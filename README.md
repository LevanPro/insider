# Insider Application

## Getting Started

### 1. Clone Repository

```bash
git clone https://github.com/LevanPro/insider.git
cd insider
```

### 2. Configure Application

Modify the configuration file to correspond with your requirements:

```bash
/config/dev.yml
```

**Important:** Add your Webhook URL and Webhook API Key to the configuration file.

### 3. Run Application

Start the application using Docker Compose:

```bash
docker compose up -d
```

### 4. Seed Database

To populate the database with initial data:

```bash
make db-seed
```

### 5. Access Swagger Documentation

Once the application is running, access the API documentation:

```
http://localhost:8080/swagger/index.html
```

### 6. View Logs

To access application logs:

```bash
docker logs <container_name>
```

To view logs in real-time:

```bash
docker logs -f <container_name>
```

---