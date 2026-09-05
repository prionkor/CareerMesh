# CareerMesh Agent Instructions

## Project

CareerMesh is an AI-powered career intelligence platform that connects
experience, skills, projects, resumes, and job opportunities to create
targeted job applications.

The project is written in Go and currently uses a monolithic architecture.

## Core principles

- Keep the architecture simple.
- Prefer a modular monolith over microservices.
- Do not introduce services, queues, event buses, or distributed systems
  unless there is a clear requirement.
- Keep the career database as the source of truth.
- AI must not invent career experience, skills, achievements, or education.
- Generated content must be traceable to stored career data.
- Prefer explicit, readable code over abstractions that provide little value.
- Avoid premature optimization.
- Avoid adding dependencies unless they provide meaningful value.

## Project Documentation

Read this file before making changes.

Before working on a feature, read the relevant documentation:

- `docs/product.md` - product concepts and domain rules
- `docs/architecture.md` - system architecture and architectural decisions
- `docs/data-model.md` - database entities and relationships
- `docs/ai.md` - AI architecture, evidence, and hallucination rules

Documentation should reflect decisions that have actually been made.
Do not change documentation speculatively or introduce architectural
decisions without a corresponding requirement.

## Technology

- Go 1.27
- PostgreSQL
- pgx/v5
- Goose migrations
- HTML/CSS/JavaScript for the initial frontend
- PostgreSQL pgvector may be introduced later for semantic search
- LLM APIs will be integrated later

## Go conventions

- Follow standard Go conventions.
- Run gofmt on changed Go files.
- Prefer the standard library where practical.
- Return errors instead of hiding them.
- Wrap errors with useful context.
- Keep functions reasonably small.
- Avoid unnecessary interfaces.
- Do not create abstractions before they are needed.

## Database

- PostgreSQL is the primary database.
- Use UUID primary keys.
- Use TIMESTAMPTZ for timestamps.
- Use DATE for career dates where time-of-day is irrelevant.
- Use foreign keys.
- Use ON DELETE CASCADE where ownership clearly implies cascading deletion.
- Use PostgreSQL ENUMs for stable controlled vocabularies.
- Database changes must be implemented through Goose migrations.
- Never modify an existing migration that may already have been applied.
- Create a new migration for schema changes.

## Data ownership

The users table represents an account.

Profiles represent career information belonging to a user.

Top-level career entities contain `user_id`:

- profiles
- experiences
- projects
- skills
- education
- certifications
- languages

Join tables should not duplicate `user_id` when ownership can be derived
from their parent records.

## AI rules

The AI layer is not the source of truth.

AI may:

- analyze job descriptions
- identify relevant skills
- rank career evidence
- suggest relevant experience
- generate resume content
- generate cover letters

AI must not:

- invent employment
- invent projects
- invent technologies
- invent achievements
- invent education
- invent certifications
- invent dates
- claim experience that does not exist in the career database

Whenever practical, generated content should retain references to the
career data used to produce it.

## Product architecture

Career data is the canonical source.

Resumes are representations of career data.

Job descriptions are external inputs.

Applications record what was actually submitted.

Do not treat generated resumes or cover letters as the canonical
representation of the user's career.

## Before changing architecture

Before introducing a new dependency, service, database, queue, framework,
or architectural pattern:

1. Check whether the existing architecture can solve the problem.
2. Prefer the simplest solution.
3. Consider whether the change creates unnecessary operational complexity.
4. Preserve the existing project conventions.

## Working on tasks

Before making significant changes:

1. Inspect the relevant existing code.
2. Understand the current architecture.
3. Reuse existing patterns.
4. Make the smallest coherent change.
5. Run formatting and tests.
6. Report any assumptions or unresolved issues.

Do not rewrite unrelated code.

## Git

Keep commits focused.

Do not modify generated files manually unless they are intended to be
tracked source files.

Never commit secrets, credentials, `.env` files, or personal career data.
