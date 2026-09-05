# CareerMesh Architecture

## Current architecture

CareerMesh is a modular monolith.

There is one Go application process.

The application contains separate internal modules for different
domains.

## Current stack

- Go 1.27
- PostgreSQL
- pgx
- Goose
- HTML/CSS/JavaScript

## Database

PostgreSQL is the primary persistence layer.

## AI

AI functionality lives inside the application initially.

Do not create a separate AI service.

## Semantic search

Semantic search may later use PostgreSQL + pgvector.

Do not introduce a separate vector database unless there is a demonstrated
need.

## Deployment

The initial deployment should remain simple.

Avoid distributed infrastructure until actual requirements justify it.
