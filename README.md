# CloudPGDump

## Usage

```
s3pgdump v1.0.0

env:
  - PGPASSWORD=<password>
  - DROPBOX_TOKEN=<token>
  - AWS_REGION=<region>
  - AWS_<creds>
USAGE: s3pgdump <pg_dump args> --s3bucket=bucket --s3prefix=/prefix --dropboxpath=/path
```

## Unitests

Make sure you have the following env vars for integration tests:
- `DROPBOX_TOKEN`: The token for your test Dropbox app
- `S3_TEST_BUCKET`: Test S3 bucket. It doesn't need to be public.
- `S3_TEST_REGION`: Test S3 region.
- `AWS_ACCESS_KEY_ID`: Key for S3 tests. Also, any other security token or secret.