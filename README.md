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

### 3. Setup Webhook Testing

1. Go to [webhook.site](https://webhook.site)
2. Click **Edit** on your unique webhook URL
3. Set the following configurations:
   - **Content Type:** `application/json`
   - **Response Body:**

```json
{
  "message": "Webhook received and processed successfully.",
  "messageId": "73c68e1a-4d22-4f81-9b6d-62f92f15a9e3"
}
```

4. Save the configuration and copy your webhook URL to use in the config file

### 4. Run Application

Start the application using Docker Compose:

```bash
docker compose up -d
```

### 5. Seed Database

To populate the database with initial data:

```bash
make db-seed
```

### 6. Access Swagger Documentation

Once the application is running, access the API documentation:

```
http://localhost:8080/swagger/index.html
```

### 7. View Logs

To access application logs:

```bash
docker logs <container_name>
```

To view logs in real-time:

```bash
docker logs -f <container_name>
```

---