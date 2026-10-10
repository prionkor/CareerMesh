# CareerMesh Data Model

## User

Represents an account.

One user can have one profile.

## Profile

Represents the user's core career identity.

A profile belongs to exactly one user.

## Experience

Represents employment or professional experience.

An experience belongs to one user.

Experience bullets contain individual pieces of career evidence.

## Job Description

Represents an external job posting saved by a user. The original submitted
description text is retained as provided and is not replaced by later
derived data.

A job description belongs to exactly one user.

## Job Application

Represents a user's application to a saved job description. It records the
application status, an optional applied date, and notes. A job description can
have many applications, and an application references exactly one job
description.

A job application belongs to exactly one user, who must also own the referenced
job description. A job description that has applications cannot be deleted.

`applied_at` is serialized as `YYYY-MM-DD` in JSON through `models.Date`.
Other resources' date fields keep their existing `time.Time` (RFC3339) format.

## Project

Represents a project associated with a user.

Projects can be associated with skills.

## Skill

Represents a skill known by the user.

Skills can be associated with experiences and projects.

## Education

Represents educational history.

## Certification

Represents a professional certification.

## Language

Represents a spoken/written language and its proficiency.

## Relationships

User
├── Profile
├── Experiences
│ └── Experience Bullets
├── Job Descriptions
│ └── Job Applications
├── Projects
├── Skills
├── Education
├── Certifications
└── Languages

Experience ──< ExperienceSkill >── Skill

Project ──< ProjectSkill >── Skill
