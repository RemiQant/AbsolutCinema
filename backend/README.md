# Project AbsolutCinema

One Paragraph of project description goes here

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development and testing purposes. See deployment for notes on how to deploy the project on a live system.

## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```

Live reload the application:
```bash
make watch
```

### Docker Commands (Development)

Start all services (app + database):
```bash
make docker-run
```

Start only the database:
```bash
make docker-run-db
```

Shutdown all containers:
```bash
make docker-down
```

### Docker Commands (Production)

Start all services in production mode:
```bash
make docker-run-prod
```

Shutdown production containers:
```bash
make docker-down-prod
```

### Integration Tests

Run database integration tests:
```bash
make itest
```
