# kouji-app-backend2

## MySQL environment

1. Ensure Docker Desktop is running.
2. Copy the example environment file and adjust credentials if needed:

   ```bash
   cp .env.example .env
   # edit .env to match your local MySQL settings
   ```

3. Start the bundled database:

   ```bash
   docker compose up -d mysql
   ```

   - Default connection values match `docker-compose.yml` (`MYSQL_USER=user`, `MYSQL_PASSWORD=pass1234`, `MYSQL_DATABASE=kouji_app`, port `3306`).
   - The Go commands described below read their settings from `.env` (or a custom `MYSQL_DSN`).

## Running migrations with GORM AutoMigrate

All schema changes are performed through GORM, not raw SQL files.

```bash
go run ./cmd/migrate
```

The migrator reads the credentials from `.env` (or a full `MYSQL_DSN`) and runs `AutoMigrate` for every registered model (currently `models.User`). Running the command is idempotent, so it is safe to execute multiple times.

## Seeding data

Use the faker-backed seeder to populate the `users` table:

```bash
# Insert 25 fake users directly into the database using the .env credentials
go run ./cmd/seed-users -n 25 -mode db
```

- The tool automatically runs `AutoMigrate` before inserting, so you can skip `cmd/migrate` if preferred.
- After inserting it prints a few sample credential pairs (`email / plain password`) to help with manual testing.
- Switch to JSON output instead of inserting with `-mode json`:

  ```bash
  go run ./cmd/seed-users -n 5 -mode json
  ```

## API

Start the API server (it uses the same `.env` for DB credentials):

```bash
go run .
```

### `GET /users`

- Returns a paginated list of users ordered by newest first.
- Optional query params:
  - `limit` (default 50, max 100)
  - `page` (default 1)
- Response body:

  ```json
  {
    "data": [
      {
        "id": 1,
        "uuid": "...",
        "name": "Demo User",
        "email": "demo@example.com",
        "avatar_url": null,
        "status": "active",
        "last_login_at": "2024-05-01T12:30:00Z",
        "created_at": "2024-04-30T09:15:00Z",
        "updated_at": "2024-05-01T12:30:00Z"
      }
    ],
    "meta": {
      "page": 1,
      "limit": 50,
      "count": 1
    }
  }
  ```

### `GET /users/:id`

- Retrieves a single user by numeric `id` or by `uuid` (e.g., `/users/42` or `/users/9b8d...`).
- Returns the same shape as an individual element in the list.

## Notes

- No SQL files are executed by the container anymore; remove the `mysql_data` volume if you need a clean slate.
- Customize DSNs using `MYSQL_DSN` for advanced options (e.g., Cloud SQL). When `MYSQL_DSN` is set it takes precedence over the individual env vars.
- The generated `PlainPassword` is **only** returned by the seeder and never saved to the DB; redact logs if sharing externally.
