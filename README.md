# War Events Clustering Backend

This project provides a backend system for clustering war events using spatial-temporal data. It consists of a Go backend API and a Python clustering microservice.

## Project Structure

```
.
├── go-backend/      # Go backend API (Fiber, PostgreSQL)
├── pyservice/       # Python clustering service (FastAPI)
├── dendogram.py     # Script for dendrogram visualization using ete3
├── .env             # Environment variables
├── run_servers.txt  # Commands for running servers
```

## Setup

### 1. Python Dependencies

```sh
cd pyservice
pip install -r requirements.txt
```

### 2. Go Dependencies

```sh
cd go-backend
go mod tidy
```

## Running the Servers

See `run_servers.txt` for commands:

- **Go Backend:**  
  Uses `air` for hot reload  
  ```sh
  cd go-backend
  air
  ```

- **Python Service:**  
  ```sh
  cd pyservice
  python -m uvicorn app:app --reload
  ```

## How It Works

- The Go backend connects to a PostgreSQL database (Supabase) and calls the Python service for clustering.
- The Python service receives event data and performs clustering using KD-Tree and hierarchical algorithms.
- Results are saved to the database and can be retrieved via API.

## Dendrogram Visualization

Use `dendogram.py` to visualize clusters from a CSV file.

```sh
python dendogram.py
```
**Note:**  
- Edit `<csv-file>` in `dendogram.py` to your CSV path.
- Requires `ete3`, `matplotlib`, and `pandas`.

## Environment Variables

See `.env` for examples:

- `DATABASE_URL` for database connection
- `PY_PORT` for Python service (default: 8000)
- `GO_PORT` for Go backend (default: 5000)

## API Endpoints (Examples)

- `POST /api/process` : Send event data for clustering
- `POST /api/events-lat-lon-date` : Retrieve events for clustering and save clusters
- `POST /api/clusters/hierarchical` : Get hierarchical cluster data