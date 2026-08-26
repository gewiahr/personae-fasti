# Storyshard production deployment

This deployment keeps the existing VPS Nginx as the public TLS router. Docker
runs the `api` and `web` services in one Compose project named `storyshard`.
Neither PostgreSQL nor Garage needs a public port.

## One-time server preparation

Create the shared network before connecting the existing services:

```sh
docker network inspect storyshard >/dev/null 2>&1 || docker network create storyshard
docker network connect storyshard postgres
```

The API reaches PostgreSQL as `postgres` through Docker DNS. Garage runs on the
host and is reached through the private Nginx bridge described below.

Create and populate the configuration volume:

```sh
docker volume create storyshard
docker run --rm \
  -v storyshard:/target \
  -v "$PWD":/source:ro \
  alpine:3.23 sh -c 'cp /source/storyshard.json /target/storyshard.json && chmod 0444 /target/storyshard.json'
```

Copy `docker-compose.prod.yml` to `/opt/storyshard/docker-compose.yml`. The
server user used by GitHub Actions must be able to run Docker and write to
`/opt/storyshard`.

## Database rehearsal and release

Do not develop migrations against the only live database. First dump and
restore a disposable rehearsal copy, apply migrations, inspect the data, then
repeat from a fresh restore until the result is reliable.

For the final rollout:

1. Stop writes to the old application.
2. Take a final custom-format dump.
3. Restore it into a new empty database on the new server.
4. As a PostgreSQL administrator, run `CREATE EXTENSION IF NOT EXISTS pgcrypto`.
5. Run the application migration container.
6. Start the API and web containers and verify health before switching Nginx.

```sh
cd /opt/storyshard
docker compose --profile tools run --rm migrate
docker compose up -d api web
docker compose ps
curl --fail http://127.0.0.1:17001/readyz
curl --fail http://127.0.0.1:17003/healthz
```

Production configuration sets `migrateOnStart` to `false`; a normal API restart
will not silently mutate the database. Development continues to migrate on
startup unless explicitly configured otherwise.

## VPS Nginx

The host Nginx proxies:

- `app.storyshard.ru` to `http://127.0.0.1:17003`
- `api.storyshard.ru` to `http://127.0.0.1:17001`
- `img.storyshard.ru` to the public `storyshard` Garage bucket

The API listens on `17001` both inside the container and on the host. The web
container listens on its standard internal port `80`, published as host port
`17003`.

The API virtual host needs `client_max_body_size 32m` and proxy read/send
timeouts of at least 130 seconds. `storyshard.nginx.conf` is the complete file
for the existing host Nginx; it is not an additional container.

Garage's public web endpoint remains on `127.0.0.1:3902`. The API container
uses `http://host.docker.internal:17002` for authenticated S3 operations; that
private Nginx listener proxies to Garage on `127.0.0.1:3900`. Compose maps
`host.docker.internal` to Docker's host gateway.

## Tagged deployments

Each repository workflow builds an immutable image tagged with its commit SHA,
also updates the local `latest` tag, and recreates only its own service. The API
workflow tests and applies migrations before replacing the API container. The
web workflow replaces the whole static-file container, so there is no interval
where the live `dist` directory is empty.
