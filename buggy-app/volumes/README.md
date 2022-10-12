Directories within `volumes` will be attached to the running containers by `docker compose`.

- `data` will contain data from the Postgres database
- `init` contains initialisation scripts for Postgres. **Note:** these will only be run if the data directory is empty, so if you want to run these again you need to do `make volumes-reset`
- `secrets` will contain files containing passwords or other information that should not be committed
