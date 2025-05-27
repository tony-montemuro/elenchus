# elenchus

## Developing

1. Clone the repository onto your machine:

    ```bash
    git clone https://github.com/tony-montemuro/elenchus.git

    # or, if you are using ssh:
    git clone git@github.com:tony-montemuro/elenchus.git
    ```

2. Run the project initialization script:

   ```bash
   ./init.sh
   ```

    This script will:

    - Install project dependencies
    - Create a self-signed TLS certificate for HTTPS

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
