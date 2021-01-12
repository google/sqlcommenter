## Postgres sample with OpenTelemetry
To run the sample:
```bash
# Create virtualenv and install dependencies
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt

# Run with Postgres DSN environment variable
export POSTGRES_DSN="postgresql://<user>:<pass>@<dbname>?host=<host>"
python psycopgy_opentelemetry_sample.py 
```
