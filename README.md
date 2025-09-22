
# Project 1: JWKS Server

Repo path: **github.com/MohamedBabker/project1-jwks**

Features:
- RSA key generation with kid + expiry
- Active + expired key
- `/.well-known/jwks.json` → only unexpired keys
- `/auth` → signed JWT with active key
- `/auth?expired=1` → signed JWT with expired key + expired exp

## Run
```bash
make tidy
make run
```

## Endpoints
- GET /.well-known/jwks.json
- POST /auth
- POST /auth?expired=1
- GET /healthz

## Tests
```bash
make cover
```

## Deliverables
- Push to GitHub at https://github.com/MohamedBabker/project1-jwks
- Screenshot: test client working
- Screenshot: coverage output >80%
