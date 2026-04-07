# WebGoDB

A simple Go web server that logs visitor IPs to PostgreSQL and displays the visitor count.  
Designed for Kubernetes deployment with GitOps (Argo CD) and CI/CD (GitHub Actions).

---

## Features

- Displays visitor count and pod name on a simple HTML page
- Logs visitor IP and timestamp to PostgreSQL
- `/healthz` endpoint for load balancer health checks
- Multi-stage Docker build (final image ~15MB)
- CI/CD via GitHub Actions (triggered by git tag)
- GitOps deployment via Argo CD

---

## Project Structure

```
WebGoDummy/
├── main.go                        # Go web server
├── index.html                     # HTML template (embedded in binary)
├── go.mod                         # Go module dependencies
├── go.sum                         # Dependency checksums
├── Dockerfile                     # Multi-stage Docker build
└── .github/workflows/ci.yaml      # GitHub Actions CI/CD pipeline
```

---

## Prerequisites

- [Go 1.20+](https://golang.org/dl/)
- [Docker](https://docs.docker.com/get-docker/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- PostgreSQL (or use the Kubernetes manifest from the GitOps repo)

---

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/juijeong8324/WebGoDummy.git
cd WebGoDummy
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Run locally (requires PostgreSQL running)

Set environment variables and run:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=yourpassword
export DB_NAME=visitorsdb

go run main.go
```

Visit `http://localhost:8080`

---

## Docker

### Build

```bash
docker build -t web-go-db:v1 .
```

### Run

```bash
docker run -p 8080:8080 \
  -e DB_HOST=<your-db-host> \
  -e DB_PORT=5432 \
  -e DB_USER=postgres \
  -e DB_PASSWORD=yourpassword \
  -e DB_NAME=visitorsdb \
  web-go-db:v1
```

---

## CI/CD Pipeline

Pushing a git tag triggers the following pipeline:

**CI (GitHub Actions)**
1. Build Docker image
2. Push image to Docker Hub
3. Update image tag in the GitOps repository (`04-deploy.yaml`)

**CD (Argo CD)**

4. Detects the change in GitOps repository
5. Automatically syncs to the cluster

```bash
git add .
git commit -m "feat: something"
git push

# Deploy
git tag v1
git push origin v1
```

### Required GitHub Secrets

| Secret | Description |
|--------|-------------|
| `DOCKERHUB_USERNAME` | Docker Hub username |
| `DOCKERHUB_TOKEN` | Docker Hub access token |
| `OPS_GITHUB_TOKEN` | GitHub personal access token (repo scope) |

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DB_HOST` | - | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | - | PostgreSQL username |
| `DB_PASSWORD` | - | PostgreSQL password |
| `DB_NAME` | - | PostgreSQL database name |

---

## Kubernetes (GitOps)

Kubernetes manifests are managed separately in the [GitOps repository](https://github.com/juijeong8324/GitOps).

| File | Description |
|------|-------------|
| `03-postgres.yaml` | PostgreSQL StatefulSet |
| `04-deploy.yaml` | Go app Deployment (replicas: 2) |
| `05-hpa.yaml` | Horizontal Pod Autoscaler |
| `06-ingress.yaml` | Ingress |
| `07-networkpolicy.yaml` | NetworkPolicy |
| `08-service.yaml` | Service |
