# elenchus

## Developing

1. Clone the repository onto your machine:

    ```bash
    git clone https://github.com/tony-montemuro/elenchus.git

    # or, if you are using ssh:
    git clone git@github.com:tony-montemuro/elenchus.git
    ```

2. Run the project initialization script:

   ⚠️ Before running this script, please ensure the following conditions are met:

   - MariaDB is installed and running.
   - The root password must be provided to setup the database.

   ```bash
   ./init.sh
   ```

    This script will:

    - Install project dependencies
    - Create the database
    - Create a database user for the application, as well as running migrations
    - Run database migrations

    Additionally, if you specify that you are in a local (non-production) environment...

    - Create a test database
    - Create a test database user for testing
    - Create a self-signed TLS certificate for HTTPS locally

## Running

To run the application:

```bash
go run ./cmd/web
```

Note that this application has many configurable options, which can be declared using [command-line flags](https://pkg.go.dev/flag). To view all configurable options, run:

```bash
go run ./cmd/web --help
```

In most cases, however, the default configuration options should suffice.
