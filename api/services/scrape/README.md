# Scraper Service

A microservice for scraping chess tournament data from chess-results.com.

## Overview

The scraper is configured to monitor tournaments in the **Barcelona metropolitan region** of Catalonia, Spain. It automatically scrapes chess-results.com every 6 hours to discover and track tournaments in the area.

## Features

- **Automatic Periodic Scraping**: Scheduled scraping of configured federations
- **Location Filtering**: Filters tournaments by geographic location (e.g., "Barcelona")
- **Manual Triggering**: API endpoint to trigger scrapes on-demand
- **Job Management**: Track scraping jobs with status and statistics
- **Tournament Storage**: Stores scraped tournament data in the database
- **Event Publishing**: Publishes events when tournaments are discovered or updated
- **Full Observability**: Metrics, tracing, and structured logging

## Architecture

```
┌─────────────────────────────────────────┐
│          HTTP API (main.go)             │
│  - Job endpoints                        │
│  - Tournament endpoints                 │
│  - Manual trigger                       │
│  - Statistics                           │
└─────────────────────────────────────────┘
                  │
┌─────────────────────────────────────────┐
│       Worker Components                 │
├─────────────────────────────────────────┤
│  Scheduler → creates periodic jobs      │
│       ↓                                 │
│  Orchestrator → polls pending jobs      │
│       ↓                                 │
│  Worker → scrapes & processes           │
│       ↓                                 │
│  Client → fetches HTML                  │
│  Parser → extracts data                 │
│  Publisher → emits events               │
└─────────────────────────────────────────┘
```

## API Endpoints

### Jobs

- `POST /api/v1/scraper/jobs` - Create a scrape job
- `GET /api/v1/scraper/jobs` - List jobs (filterable by federation/status)
- `GET /api/v1/scraper/jobs/{id}` - Get specific job

### Tournaments

- `GET /api/v1/scraper/tournaments` - List scraped tournaments (filterable by federation)

### Actions

- `POST /api/v1/scraper/trigger` - Manually trigger a scrape
- `GET /api/v1/scraper/statistics` - Get aggregated statistics

## Configuration

Environment variables (prefix: `SCRAPER_`):

| Variable | Default | Description |
|----------|---------|-------------|
| `SCRAPER_FEDERATIONS` | `CAT` | Comma-separated list of federations to scrape (CAT=Catalonia) |
| `SCRAPER_LOCATION_FILTER` | `Barcelona` | Filter tournaments by location (case-insensitive substring match) |
| `SCRAPER_DISCOVERY_INTERVAL` | `6h` | Interval between automatic scrapes |
| `SCRAPER_REQUEST_DELAY` | `2s` | Delay between HTTP requests (politeness) |
| `SCRAPER_USER_AGENT` | `Mozilla/5.0...` | User agent for HTTP requests |
| `SCRAPER_PUBLISHER_TYPE` | `log` | Publisher type (log, webhook, nats) |
| `SCRAPER_API_HOST` | `0.0.0.0:8080` | API server address |
| `SCRAPER_DEBUG_HOST` | `0.0.0.0:8081` | Debug server address |
| `DATABASE_URL` | - | PostgreSQL connection string (required) |

## Barcelona Region Configuration

The scraper is pre-configured for the Barcelona metropolitan region:

- **Federation**: CAT (Federació Catalana d'Escacs - Catalan Chess Federation)
- **Location Filter**: "Barcelona" (matches Barcelona, L'Hospitalet de Llobregat, Badalona, etc.)
- **How it works**:
  1. Scrapes all tournaments from the Catalan federation
  2. Filters for tournaments with "Barcelona" in the location or name
  3. Stores only tournaments matching the filter

## Usage Examples

### Create a Scrape Job

```bash
curl -X POST http://localhost:8080/api/v1/scraper/jobs \
  -H "Content-Type: application/json" \
  -d '{"federation":"ESP"}'
```

### List Jobs

```bash
# All jobs
curl http://localhost:8080/api/v1/scraper/jobs

# Filter by status
curl "http://localhost:8080/api/v1/scraper/jobs?status=completed"

# Filter by federation
curl "http://localhost:8080/api/v1/scraper/jobs?federation=ESP"
```

### Get Job Details

```bash
curl http://localhost:8080/api/v1/scraper/jobs/{job-id}
```

### List Scraped Tournaments

```bash
# All tournaments
curl http://localhost:8080/api/v1/scraper/tournaments

# Filter by federation
curl "http://localhost:8080/api/v1/scraper/tournaments?federation=ESP"
```

### Manually Trigger a Scrape

```bash
# Trigger Barcelona tournament scrape
curl -X POST http://localhost:8080/api/v1/scraper/trigger \
  -H "Content-Type: application/json" \
  -d '{"federation":"CAT"}'
```

Response:
```json
{
  "job_id": "123e4567-e89b-12d3-a456-426614174000",
  "federation": "CAT",
  "message": "scrape job triggered"
}
```

**Note**: The scraper will automatically filter for Barcelona tournaments based on the `SCRAPER_LOCATION_FILTER` configuration.

### Get Statistics

```bash
curl http://localhost:8080/api/v1/scraper/statistics
```

Response:
```json
{
  "total_jobs": 45,
  "by_status": {
    "completed": 40,
    "running": 2,
    "pending": 1,
    "failed": 2
  },
  "by_federation": {
    "ESP": 30,
    "CAT": 15
  },
  "total_found": 1250,
  "total_created": 980,
  "total_updated": 270
}
```

## Development

### Local Setup

```bash
# Build the service
make scraper

# Run all services (including scraper)
make dev

# Clean and restart
make dev-clean

# Stop all services
make down-dev
```

### Monitoring

- **Metrics**: http://localhost:8081/debug/vars
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3100 (admin/admin)
- **Tempo**: http://localhost:3200

### Logs

View logs with docker compose:
```bash
docker compose -f zoltan/compose/docker-compose.yml \
  -f zoltan/compose/docker-compose.dev.yml \
  logs -f scraper
```

## Database Schema

### scrape_jobs

Tracks scraping operations:

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| federation | VARCHAR(10) | Federation code (ESP, CAT, etc.) |
| status | VARCHAR(20) | pending, running, completed, failed |
| started_at | TIMESTAMP | When the job started |
| completed_at | TIMESTAMP | When the job completed (nullable) |
| error | TEXT | Error message if failed (nullable) |
| found | INT | Number of tournaments found |
| created | INT | Number of tournaments created |
| updated | INT | Number of tournaments updated |
| date_created | TIMESTAMP | Record creation time |
| date_updated | TIMESTAMP | Record last update time |

### scraped_tournaments

Stores tournament data:

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| external_id | VARCHAR(50) | Unique chess-results.com ID |
| name | VARCHAR(255) | Tournament name |
| url | VARCHAR(500) | Tournament URL |
| federation | VARCHAR(10) | Federation code |
| organizer | VARCHAR(255) | Organizer name (nullable) |
| arbiter | VARCHAR(255) | Arbiter name (nullable) |
| location | VARCHAR(255) | Tournament location (nullable) |
| start_date | VARCHAR(50) | Start date (nullable) |
| players | INT | Number of players |
| rounds | INT | Number of rounds |
| time_control | VARCHAR(255) | Time control (nullable) |
| last_scraped_at | TIMESTAMP | Last time scraped |
| date_created | TIMESTAMP | Record creation time |
| date_updated | TIMESTAMP | Record last update time |

## How It Works

### Automatic Scraping

1. **Scheduler** creates jobs every 6 hours for configured federations (CAT)
2. **Orchestrator** polls for pending jobs every 10 seconds
3. **Worker** processes each job:
   - Fetches federation page (chess-results.com/fed.aspx?fed=CAT)
   - Parses tournament listings from the page
   - **Filters tournaments** by location (keeps only "Barcelona" matches)
   - For each filtered tournament:
     - Fetches tournament details
     - Parses tournament info
     - Upserts to database
     - Publishes event
4. Job status updated (pending → running → completed/failed)

### Location Filtering

The worker applies location filtering in two ways:
1. **Location field**: Checks if tournament location contains "Barcelona"
2. **Tournament name**: Checks if tournament name contains "Barcelona"

This ensures we capture tournaments like:
- "Open de Barcelona 2024" in "Barcelona"
- "Torneig Rapid" in "Barcelona, Spain"
- "Catalunya Championship" in "L'Hospitalet de Llobregat, Barcelona"

### Manual Scraping

Use the `/api/v1/scraper/trigger` endpoint to immediately create and process a job for any federation.

## Events

The scraper publishes these events:

- `tournament.discovered` - New tournament found
- `tournament.updated` - Existing tournament updated
- `scrape.error` - Error during scraping

Events include:
- Event type
- Timestamp
- Tournament ID
- Federation
- Payload (tournament data or error details)

## Error Handling

- Failed jobs are marked with status `failed` and error message stored
- Scraper continues processing other tournaments if one fails
- HTTP requests retry up to 3 times with exponential backoff
- Rate limiting respects configured delay between requests

## Performance

- Configurable request delay prevents overwhelming chess-results.com
- Parallel processing of pending jobs
- Database connection pooling
- Efficient upsert operations (create or update in single query)
