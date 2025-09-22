
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
- <img width="554" height="198" alt="image" src="https://github.com/user-attachments/assets/88d56a65-f0f5-42ae-9c5f-d844b0fa72c7" />

- Screenshot: coverage output >80%
- <img width="914" height="276" alt="image" src="https://github.com/user-attachments/assets/b2146e4b-51e2-4a2f-8140-842c27418868" />

