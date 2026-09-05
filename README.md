# CareerMesh

**CareerMesh is an AI powered career intelligence system that connects your experience, skills, projects, resumes, and job opportunities to help you create better, more targeted job applications.**

> 🚧 CareerMesh is currently under active development.

## The idea

Most people maintain their career information across multiple resumes, LinkedIn, documents, project notes, and old job applications.

CareerMesh aims to bring that information into one structured career knowledge base.

Instead of treating a resume as the source of truth, CareerMesh treats your **career history** as the source of truth and generates different representations of it when needed.

For example:

```text
                    CareerMesh
                        │
             ┌──────────┼──────────┐
             │          │          │
        Experience   Projects    Skills
             │          │          │
             └──────────┼──────────┘
                        │
                Career Knowledge
                        │
              ┌─────────┴─────────┐
              │                   │
         Job Description      Existing CVs
              │                   │
              └─────────┬─────────┘
                        │
                 Job Matching
                        │
              ┌─────────┴─────────┐
              │                   │
         Tailored Resume     Cover Letter
              │                   │
              └─────────┬─────────┘
                        │
                   Application
```

## Planned capabilities

- Store structured career information
  - Experience
  - Projects
  - Skills
  - Education
  - Achievements

- Import and manage existing resumes
- Analyze job descriptions
- Match jobs against career experience
- Select or generate an appropriate resume
- Generate tailored cover letters
- Keep a history of job applications
- Track which resume and cover letter were used for each application
- Semantic search across career information
- Generate public, tailored career/resume profiles
- Eventually provide an AI-assisted career intelligence layer

## Design principles

### Career data is the source of truth

AI should help organize and express career information, not invent it.

Generated resumes and cover letters should be based on verifiable information stored in CareerMesh.

### One career, many representations

A person should not need to maintain six completely independent resumes.

CareerMesh should maintain one underlying career profile and generate different representations depending on the target role.

### Evidence over invention

When possible, generated claims should be traceable back to the underlying career data.

### Simple first

CareerMesh starts as a **Go monolith**.

The initial goal is to build a useful product before introducing unnecessary infrastructure or architectural complexity.

## Technology

The initial stack is intentionally small:

- **Go**
- **PostgreSQL**
- **pgvector** for semantic search as the AI layer develops
- HTML/CSS/JavaScript for the web interface
- LLM and embedding APIs

The technology stack may evolve as the project develops.

## Project status

CareerMesh is currently in the early development stage.

The first milestone is to build the underlying career knowledge base and application workflow before adding more advanced AI capabilities.

## Source availability

CareerMesh is **source-available**.

The source code is publicly available so developers, engineers, and recruiters can inspect the implementation and understand how the project is built.

The project is not currently intended to be a conventional open-source project with unrestricted commercial use.

See [`LICENSE`](LICENSE) for the applicable terms.

## Why this project?

CareerMesh is being developed as a real product and as an experiment in building an AI-assisted application around structured personal data.

The project is also an opportunity to explore:

- Go application development
- AI application architecture
- Retrieval-augmented generation (RAG)
- Vector search
- Document processing
- Structured LLM output
- Resume generation
- Product-oriented software design

## Roadmap

The roadmap will evolve as the project develops.

### Phase 1 - Career Foundation

- [ ] Project structure
- [ ] PostgreSQL setup
- [ ] Career profile
- [ ] Experience
- [ ] Projects
- [ ] Skills
- [ ] Education
- [ ] Resume/document storage

### Phase 2 - Job Management

- [ ] Job description storage
- [ ] Job parsing
- [ ] Job requirements
- [ ] Career/job matching

### Phase 3 - AI

- [ ] Embeddings
- [ ] Semantic career search
- [ ] RAG pipeline
- [ ] Resume tailoring
- [ ] Cover letter generation

### Phase 4 - Applications

- [ ] Application tracking
- [ ] Resume version history
- [ ] Cover letter history
- [ ] Application analytics

### Phase 5 - Public Career Profiles

- [ ] Dynamic resume generation
- [ ] Printable resumes
- [ ] Shareable career profiles
- [ ] Signed resume URLs

## License

### CareerMesh Source Available License

Copyright (c) 2026 Sisir Kanti Adhikari. All rights reserved.

The source code of CareerMesh is made publicly available for inspection and educational purposes.

Permission is granted to view and study the source code through the public repository.

No permission is granted to:

- use the software for personal, commercial, or production purposes;
- deploy or operate the software;
- copy, modify, or create derivative works for use or distribution;
- redistribute or sublicense the software;
- sell, rent, license, host, or provide the software as a service;
- incorporate the software or substantial portions of it into another product or service.

Any rights not expressly granted above are reserved by the copyright holder.

For permissions beyond those granted here, please contact the copyright holder.
