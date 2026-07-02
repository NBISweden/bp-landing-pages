# Bigpicture-landing-pages-service

Bigpicture project dataset landing page generator

## Development environment

To start S3 minio service, navigate to dev_utils folder and run

```bash
docker compose up
```

Set config.yaml in environment file

```bash
export CONFIGFILE="dev_utils/config.yaml"
```

Start the application by running

```bash
go run .
```
