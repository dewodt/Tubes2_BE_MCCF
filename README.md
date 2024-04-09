# Wikirace

## How To Run

1. Make sure Docker Desktop is running

2. For development (with hot reload), run

```bash
docker compose -f "docker-compose.development.yml" up --build
```

3. For production, run

```bash
docker compose -f "docker-compose.production.yml" up --build
```

4. Visit [http://localhost:8080](http://localhost:8080)
5. To stop the process, find the container id

```bash
docker ps
```

6. Stop the process

```bash
docker stop <container_id>
```
