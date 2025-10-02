#!/bin/bash
set -euo pipefail

: "${POSTGRES_HOST:=postgres}"
: "${POSTGRES_PORT:=5432}"
: "${POSTGRES_USER:?POSTGRES_USER must be set}"
: "${POSTGRES_PASSWORD:?POSTGRES_PASSWORD must be set}"
: "${POSTGRES_DB:?POSTGRES_DB must be set}"
MIGRATIONS_DIR=${MIGRATIONS_DIR:-/migrations}

if [[ ! -d "$MIGRATIONS_DIR" ]]; then
  log "Migrations directory $MIGRATIONS_DIR not found" >&2
  exit 1
fi

export PGPASSWORD="$POSTGRES_PASSWORD"

log() {
  echo "[postgres-bootstrap] $*"
}

wait_for_db() {
  for attempt in $(seq 1 30); do
    if pg_isready -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" >/dev/null 2>&1; then
      return 0
    fi
    log "Waiting for Postgres to become ready (attempt $attempt/30)"
    sleep 2
  done
  log "Postgres did not become ready in time" >&2
  return 1
}

run_psql() {
  psql -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" "$@"
}

run_sql() {
  local file="$1"
  if [[ ! -f "$file" ]]; then
    log "Migration file $file not found" >&2
    exit 1
  fi
  log "Running migration $(basename "$file")"
  run_psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1 -f "$file"
}

ensure_database() {
  local exists
  exists=$(run_psql -U "$POSTGRES_USER" -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname = '$POSTGRES_DB'")
  if [[ "$exists" != "1" ]]; then
    log "Creating database '$POSTGRES_DB'"
    createdb -h "$POSTGRES_HOST" -p "$POSTGRES_PORT" -U "$POSTGRES_USER" "$POSTGRES_DB"
  else
    log "Database '$POSTGRES_DB' already exists"
  fi
}

apply_migrations_if_needed() {
  local has_users_table
  has_users_table=$(run_psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -tAc "SELECT 1 FROM pg_tables WHERE schemaname = 'public' AND tablename = 'users'") || true
  if [[ "$has_users_table" == "1" ]]; then
    log "Schema already present in '$POSTGRES_DB'; skipping migrations"
    return 0
  fi

  log "Applying baseline migrations to '$POSTGRES_DB'"
  shopt -s nullglob
  readarray -t migrations < <(find "$MIGRATIONS_DIR" -maxdepth 1 -type f -name '*.sql' -print | sort)

  if [[ ${#migrations[@]} -eq 0 ]]; then
    log "No migration files found in $MIGRATIONS_DIR" >&2
    exit 1
  fi

  for migration in "${migrations[@]}"; do
    run_sql "$migration"
  done

  shopt -u nullglob
  log "Migrations applied successfully"
}

main() {
  wait_for_db
  ensure_database
  apply_migrations_if_needed
}

main "$@"
