# PostgreSQL + pgAdmin + Docker Setup

## 1. Define Environment Variables

**Purpose:** Store database credentials and configuration outside source code.

`.env`

```env
POSTGRES_DB=
POSTGRES_USER=
POSTGRES_PASSWORD=

PGADMIN_DEFAULT_EMAIL=
PGADMIN_DEFAULT_PASSWORD=
```

### Definitions

- `POSTGRES_DB` → Database to create on first startup.
- `POSTGRES_USER` → Database user to create.
- `POSTGRES_PASSWORD` → Password for the database user.
- `PGADMIN_DEFAULT_EMAIL` → Login email for pgAdmin.
- `PGADMIN_DEFAULT_PASSWORD` → Login password for pgAdmin.

---

## 2. Create Docker Compose

**Purpose:** Run PostgreSQL and pgAdmin as containers.

`docker-compose.yml`

```yaml
services:
  postgres:
    image: postgres:17
    container_name: postgres-db
    env_file:
      - .env
    ports:
      - '5432:5432'
    volumes:
      - postgres_data:/var/lib/postgresql/data

  pgadmin:
    image: dpage/pgadmin4
    container_name: pgadmin
    env_file:
      - .env
    ports:
      - '5050:80'
    depends_on:
      - postgres

volumes:
  postgres_data:
```

### Definitions

- `image` → Docker image to run.
- `container_name` → Name of the running container.
- `ports` → Expose container ports to host machine.
- `volumes` → Persist database data across restarts.
- `depends_on` → Start PostgreSQL before pgAdmin.

---

## 3. Start Containers

**Purpose:** Create and run PostgreSQL and pgAdmin.

```bash
docker compose up -d
```

Verify:

```bash
docker ps
```

Expected:

```text
postgres-db
pgadmin
```

---

## 4. Access pgAdmin

**Purpose:** GUI for managing PostgreSQL.

Open:

```text
http://localhost:5050
```

Login:

```text
Email = PGADMIN_DEFAULT_EMAIL
Password = PGADMIN_DEFAULT_PASSWORD
```

---

## 5. Register PostgreSQL Server

**Purpose:** Connect pgAdmin to the PostgreSQL container.

### General Tab

```text
Name: POSTGRES_DB
```

### Connection Tab

```text
Host name/address: postgres (dont put localhost)
Port: 5432
Maintenance database: goauth
Username:POSTGRES_USER
Password:POSTGRES_PASSWORD
```

### Important

Use:

```text
postgres
```

NOT:

```text
localhost
```

because pgAdmin and PostgreSQL communicate over Docker's internal network.

---

## 6. Verify Database Access

Connect directly to PostgreSQL:

```bash
docker exec -it postgres-db psql -U srijan -d goauth
```

Expected:

```sql
goauth=>
```

Useful commands:

```sql
\l
```

List databases.

```sql
\du
```

List users/roles.

Exit:

```sql
\q
```

---

## 7. Connect From Go

Connection string:

```env
DATABASE_URL=postgres://srijan:srijanpassword@localhost:5432/goauth?sslmode=disable
```

Explanation:

- `srijan` → DB username
- `srijanpassword` → DB password
- `localhost:5432` → PostgreSQL exposed by Docker
- `goauth` → Database name

---

## 8. Recreate Database After Changing Credentials

PostgreSQL only uses `POSTGRES_*` variables during first initialization.

If credentials change:

```bash
docker compose down -v
docker compose up -d
```

### Why?

`-v` removes the database volume so PostgreSQL initializes again using the new `.env`.

Without deleting the volume, old credentials remain active.

---

## Workflow Summary

```text
.env
  ↓
docker compose up -d
  ↓
PostgreSQL container starts
  ↓
Creates DB + User from POSTGRES_* vars
  ↓
pgAdmin starts
  ↓
Register Server (host = postgres)
  ↓
Connect and manage database
  ↓
Go app connects using DATABASE_URL
```
