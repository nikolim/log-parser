# Log Analyser

Package to pull logs from Github Actions, analyse them and write them into a database.

### Run
```bash
go build
./log-analyser --config config.toml
``` 

### Token
In order to analyse logs of private repositories, you need to provide a Github-API token.
These values can be saved in a file called ```.env``` with or as an environment variables.
```
GITHUB_API_TOKEN="your-api-token"
```

### Database
Make sure to have a Postgres database running. For example:
```
docker run --name postgres -p 5432:5432 -e POSTGRES_PASSWORD=test -e POSTGRES_USER=test postgres:alpine
```
To create the required tables, you can use the ```log_analyser/schema.sql``` file.

### Configuration
The sample configuration file is located at ```log_analyser/config.toml```.
Make sure to set the correct values for the database connection.

