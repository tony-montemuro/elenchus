# elenchus

elenchus is an AI-powered education platform. The following demos showcase some of the functionality:

https://github.com/user-attachments/assets/b6aceea2-a9e9-4719-a1d0-a520f34f49d2

https://github.com/user-attachments/assets/6e43b954-b30f-4873-b9dd-40aec7bd292f

https://github.com/user-attachments/assets/f5869a68-1355-40c8-af69-20c4ece2065d

https://github.com/user-attachments/assets/82fe0091-de2a-4aae-b3cf-8edad6969989

## Goals

- Often we hear about the negative effects AI has on cognitive abilities. This application is to designed as an educational AI application that can actually enhance the learning process.
- Advance my understanding of the Go programming language.
- Advance my understanding of backend concepts: 
   - Databases 
   - Session management 
   - Middleware 
   - API design
   - Remote API calls
- Advance my understanding of frontend fundamentals: vanilla HTML, CSS, and JavaScript.
- Introduction to formal testing.

## Project Status

This project is still a **work in progress**. Many of the core features are complete, but there is much work to be done before this application is production ready. Some examples of missing features include:

- [ ] Quiz search / filtration
- [ ] Quiz deletion
- [ ] More advanced error handling
- [ ] More advanced quiz edits (reordering, adding, deleting, etc.)
- [ ] Manual quiz creation
- [ ] More advanced authentication (2FA / OAuth)
- [ ] Stripe intergartion

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

3. Export your OpenAI API key:

   ```bash
   export OPENAI_API_KEY="<YOUR_API_KEY>"
   ```

   This key is **required** to enable quiz creation. If you don't have one yet, you can create an API key at [https://platform.openai.com/api-keys](https://platform.openai.com/api-keys).

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
