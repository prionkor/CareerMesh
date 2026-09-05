# CareerMesh AI

## Fundamental rule

The AI system must not become the source of truth.

The career database is authoritative.

## Allowed AI operations

AI may:

- parse job descriptions
- extract requirements
- classify skills
- identify relevant career evidence
- rank evidence
- generate resume content
- generate cover letters
- provide recommendations

## Evidence

Generated career claims should originate from stored career data.

Whenever practical, generated content should retain references to the
source records used to produce it.

## Hallucination prevention

The system should prefer omission over fabrication.

If the career database does not contain evidence for a requirement,
the AI should not claim that the user has that experience.

## Resume generation

A generated resume may reorganize, summarize, and rewrite existing
career information.

It must not create new factual claims.

## Job matching

Matching should distinguish between:

- Explicitly demonstrated skills
- Related experience
- Potentially transferable experience
- Missing requirements
