-- Fictional CareerMesh sample profile; all employers, projects, and identifiers are illustrative.
-- Run from the repository root only against a dedicated development/test database:
-- psql "$TEST_DATABASE_URL" -v ON_ERROR_STOP=1 -f seeds/profiles/senior_backend_platform.sql
-- The seeded account has an intentionally invalid password hash and cannot be used to log in.

BEGIN;

INSERT INTO users (id, email, password_hash)
VALUES (
    '10000000-0000-7000-8000-000000000001',
    'avery.lin@example.com',
    '!disabled-seed-account-no-password!'
)
ON CONFLICT (email) DO UPDATE SET
    password_hash = EXCLUDED.password_hash
WHERE users.id = EXCLUDED.id;

-- SQL seeds bypass user.Service.CreateWithHashedPassword, so assign only the default role.
INSERT INTO user_roles (user_id, role_id)
VALUES (
    '10000000-0000-7000-8000-000000000001',
    (SELECT id FROM roles WHERE name = 'user')
)
ON CONFLICT (user_id) DO UPDATE SET
    role_id = EXCLUDED.role_id;

INSERT INTO profiles (id, user_id, name, headline, location, phone, website, github, linkedin)
VALUES (
    '20000000-0000-7000-8000-000000000001',
    '10000000-0000-7000-8000-000000000001',
    'Avery Lin',
    'Senior Backend / Platform Engineer | Go, Distributed Systems, Cloud Infrastructure',
    'Portland, Oregon',
    '',
    '',
    '',
    ''
)
ON CONFLICT (user_id) DO UPDATE SET
    id = EXCLUDED.id,
    name = EXCLUDED.name,
    headline = EXCLUDED.headline,
    location = EXCLUDED.location,
    phone = EXCLUDED.phone,
    website = EXCLUDED.website,
    github = EXCLUDED.github,
    linkedin = EXCLUDED.linkedin;

INSERT INTO experiences (
    id, user_id, company, website, title, employment_type, location,
    start_date, end_date, description
)
VALUES
    (
        '30000000-0000-7000-8000-000000000001',
        '10000000-0000-7000-8000-000000000001',
        'Ashline Software', '', 'Software Engineer', 'full_time', 'Portland, Oregon',
        DATE '2016-01-11', DATE '2018-06-30',
        'Built backend services and internal tooling for a growing subscription software platform, developing a strong foundation in API design, relational data modeling, and production support.'
    ),
    (
        '30000000-0000-7000-8000-000000000002',
        '10000000-0000-7000-8000-000000000001',
        'Copperleaf Digital', '', 'Backend Engineer', 'full_time', 'Portland, Oregon',
        DATE '2018-07-01', DATE '2021-02-28',
        'Owned Go services for customer-facing workflows and asynchronous processing. Improved PostgreSQL query performance and introduced Redis-backed caching and Kafka event consumers for high-volume integrations.'
    ),
    (
        '30000000-0000-7000-8000-000000000003',
        '10000000-0000-7000-8000-000000000001',
        'RelayForge', '', 'Senior Backend Engineer', 'full_time', 'Remote',
        DATE '2021-03-01', DATE '2023-12-31',
        'Led backend design for a multi-tenant event processing platform, setting service reliability practices and partnering with infrastructure engineers on containerized deployments and delivery automation.'
    ),
    (
        '30000000-0000-7000-8000-000000000004',
        '10000000-0000-7000-8000-000000000001',
        'CloudKite Systems', '', 'Senior Backend / Platform Engineer', 'full_time', 'Portland, Oregon',
        DATE '2024-01-01', NULL,
        'Designs shared platform services and paved-road tooling for product teams, with a focus on Go, Kubernetes, distributed systems, observability, and dependable release workflows.'
    )
ON CONFLICT (id) DO UPDATE SET
    user_id = EXCLUDED.user_id,
    company = EXCLUDED.company,
    website = EXCLUDED.website,
    title = EXCLUDED.title,
    employment_type = EXCLUDED.employment_type,
    location = EXCLUDED.location,
    start_date = EXCLUDED.start_date,
    end_date = EXCLUDED.end_date,
    description = EXCLUDED.description;

DELETE FROM experience_bullets
WHERE experience_id IN (
    '30000000-0000-7000-8000-000000000001',
    '30000000-0000-7000-8000-000000000002',
    '30000000-0000-7000-8000-000000000003',
    '30000000-0000-7000-8000-000000000004'
);

INSERT INTO experience_bullets (id, experience_id, content, position)
VALUES
    ('40000000-0000-7000-8000-000000000001', '30000000-0000-7000-8000-000000000001', 'Implemented versioned REST endpoints and request validation for account and billing workflows.', 1),
    ('40000000-0000-7000-8000-000000000002', '30000000-0000-7000-8000-000000000001', 'Added PostgreSQL indexes and rewrote slow reporting queries after reviewing production query plans.', 2),
    ('40000000-0000-7000-8000-000000000003', '30000000-0000-7000-8000-000000000001', 'Automated test and deployment steps in CI/CD pipelines, reducing manual release work for the team.', 3),
    ('40000000-0000-7000-8000-000000000004', '30000000-0000-7000-8000-000000000002', 'Developed Go workers that consumed Kafka events and safely retried transient integration failures.', 1),
    ('40000000-0000-7000-8000-000000000005', '30000000-0000-7000-8000-000000000002', 'Introduced Redis caching for frequently requested reference data with explicit expiry and invalidation behavior.', 2),
    ('40000000-0000-7000-8000-000000000006', '30000000-0000-7000-8000-000000000002', 'Worked with product and support teams to trace incidents across API, database, and background-job boundaries.', 3),
    ('40000000-0000-7000-8000-000000000007', '30000000-0000-7000-8000-000000000003', 'Designed tenant-aware event ingestion services with idempotent consumers and bounded retry policies.', 1),
    ('40000000-0000-7000-8000-000000000008', '30000000-0000-7000-8000-000000000003', 'Standardized Docker images and Kubernetes deployment templates shared by six service teams.', 2),
    ('40000000-0000-7000-8000-000000000009', '30000000-0000-7000-8000-000000000003', 'Established service-level dashboards and actionable alerts using Prometheus metrics and structured logs.', 3),
    ('40000000-0000-7000-8000-000000000010', '30000000-0000-7000-8000-000000000004', 'Building reusable Go libraries for service configuration, health checks, and graceful shutdown.', 1),
    ('40000000-0000-7000-8000-000000000011', '30000000-0000-7000-8000-000000000004', 'Improving platform visibility with OpenTelemetry traces connected to service and infrastructure metrics.', 2),
    ('40000000-0000-7000-8000-000000000012', '30000000-0000-7000-8000-000000000004', 'Coaching engineers on API contracts, production readiness, and operational ownership.', 3)
ON CONFLICT (id) DO UPDATE SET
    experience_id = EXCLUDED.experience_id,
    content = EXCLUDED.content,
    position = EXCLUDED.position;

INSERT INTO skills (id, user_id, name, category)
VALUES
    ('50000000-0000-7000-8000-000000000001', '10000000-0000-7000-8000-000000000001', 'Go', 'Programming Languages'),
    ('50000000-0000-7000-8000-000000000002', '10000000-0000-7000-8000-000000000001', 'PostgreSQL', 'Databases'),
    ('50000000-0000-7000-8000-000000000003', '10000000-0000-7000-8000-000000000001', 'Redis', 'Databases'),
    ('50000000-0000-7000-8000-000000000004', '10000000-0000-7000-8000-000000000001', 'Kafka', 'Messaging'),
    ('50000000-0000-7000-8000-000000000005', '10000000-0000-7000-8000-000000000001', 'Docker', 'Containers'),
    ('50000000-0000-7000-8000-000000000006', '10000000-0000-7000-8000-000000000001', 'Kubernetes', 'Platform Engineering'),
    ('50000000-0000-7000-8000-000000000007', '10000000-0000-7000-8000-000000000001', 'REST APIs', 'API Design'),
    ('50000000-0000-7000-8000-000000000008', '10000000-0000-7000-8000-000000000001', 'CI/CD', 'Delivery'),
    ('50000000-0000-7000-8000-000000000009', '10000000-0000-7000-8000-000000000001', 'Prometheus', 'Observability'),
    ('50000000-0000-7000-8000-000000000010', '10000000-0000-7000-8000-000000000001', 'OpenTelemetry', 'Observability')
ON CONFLICT (id) DO UPDATE SET
    user_id = EXCLUDED.user_id,
    name = EXCLUDED.name,
    category = EXCLUDED.category;

INSERT INTO projects (id, user_id, name, description, url, repository_url, start_date, end_date)
VALUES
    (
        '60000000-0000-7000-8000-000000000001',
        '10000000-0000-7000-8000-000000000001',
        'Streamline Event Gateway',
        'A reference event-ingestion service that validates partner webhooks, writes durable Kafka messages, and exposes replay and delivery-status APIs. Includes idempotency keys, bounded retries, and operational dashboards.',
        '', '', DATE '2021-06-01', DATE '2022-03-31'
    ),
    (
        '60000000-0000-7000-8000-000000000002',
        '10000000-0000-7000-8000-000000000001',
        'Atlas Deployment Control Plane',
        'A Kubernetes-oriented deployment service for coordinating staged releases across multiple environments. Provides a Go REST API, PostgreSQL-backed release state, policy checks, and traceable rollout history.',
        '', '', DATE '2023-02-01', DATE '2024-01-31'
    )
ON CONFLICT (id) DO UPDATE SET
    user_id = EXCLUDED.user_id,
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    url = EXCLUDED.url,
    repository_url = EXCLUDED.repository_url,
    start_date = EXCLUDED.start_date,
    end_date = EXCLUDED.end_date;

DELETE FROM project_skills
WHERE project_id IN (
    '60000000-0000-7000-8000-000000000001',
    '60000000-0000-7000-8000-000000000002'
);

INSERT INTO project_skills (project_id, skill_id)
VALUES
    ('60000000-0000-7000-8000-000000000001', '50000000-0000-7000-8000-000000000001'),
    ('60000000-0000-7000-8000-000000000001', '50000000-0000-7000-8000-000000000002'),
    ('60000000-0000-7000-8000-000000000001', '50000000-0000-7000-8000-000000000003'),
    ('60000000-0000-7000-8000-000000000001', '50000000-0000-7000-8000-000000000004'),
    ('60000000-0000-7000-8000-000000000001', '50000000-0000-7000-8000-000000000007'),
    ('60000000-0000-7000-8000-000000000002', '50000000-0000-7000-8000-000000000001'),
    ('60000000-0000-7000-8000-000000000002', '50000000-0000-7000-8000-000000000002'),
    ('60000000-0000-7000-8000-000000000002', '50000000-0000-7000-8000-000000000005'),
    ('60000000-0000-7000-8000-000000000002', '50000000-0000-7000-8000-000000000006'),
    ('60000000-0000-7000-8000-000000000002', '50000000-0000-7000-8000-000000000008'),
    ('60000000-0000-7000-8000-000000000002', '50000000-0000-7000-8000-000000000010')
ON CONFLICT (project_id, skill_id) DO NOTHING;

DELETE FROM experience_skills
WHERE experience_id IN (
    '30000000-0000-7000-8000-000000000001',
    '30000000-0000-7000-8000-000000000002',
    '30000000-0000-7000-8000-000000000003',
    '30000000-0000-7000-8000-000000000004'
);

INSERT INTO experience_skills (experience_id, skill_id)
VALUES
    ('30000000-0000-7000-8000-000000000001', '50000000-0000-7000-8000-000000000001'),
    ('30000000-0000-7000-8000-000000000001', '50000000-0000-7000-8000-000000000002'),
    ('30000000-0000-7000-8000-000000000001', '50000000-0000-7000-8000-000000000007'),
    ('30000000-0000-7000-8000-000000000001', '50000000-0000-7000-8000-000000000008'),
    ('30000000-0000-7000-8000-000000000002', '50000000-0000-7000-8000-000000000001'),
    ('30000000-0000-7000-8000-000000000002', '50000000-0000-7000-8000-000000000002'),
    ('30000000-0000-7000-8000-000000000002', '50000000-0000-7000-8000-000000000003'),
    ('30000000-0000-7000-8000-000000000002', '50000000-0000-7000-8000-000000000004'),
    ('30000000-0000-7000-8000-000000000003', '50000000-0000-7000-8000-000000000001'),
    ('30000000-0000-7000-8000-000000000003', '50000000-0000-7000-8000-000000000004'),
    ('30000000-0000-7000-8000-000000000003', '50000000-0000-7000-8000-000000000005'),
    ('30000000-0000-7000-8000-000000000003', '50000000-0000-7000-8000-000000000006'),
    ('30000000-0000-7000-8000-000000000003', '50000000-0000-7000-8000-000000000009'),
    ('30000000-0000-7000-8000-000000000004', '50000000-0000-7000-8000-000000000001'),
    ('30000000-0000-7000-8000-000000000004', '50000000-0000-7000-8000-000000000002'),
    ('30000000-0000-7000-8000-000000000004', '50000000-0000-7000-8000-000000000006'),
    ('30000000-0000-7000-8000-000000000004', '50000000-0000-7000-8000-000000000008'),
    ('30000000-0000-7000-8000-000000000004', '50000000-0000-7000-8000-000000000010')
ON CONFLICT (experience_id, skill_id) DO NOTHING;

INSERT INTO education (id, user_id, institution, degree, field, start_date, end_date, description)
VALUES (
    '70000000-0000-7000-8000-000000000001',
    '10000000-0000-7000-8000-000000000001',
    'Pacific Northwest Institute of Technology',
    'Bachelor of Science',
    'Computer Science',
    DATE '2011-09-01',
    DATE '2015-06-15',
    'Coursework included operating systems, database systems, computer networks, algorithms, and distributed computing.'
)
ON CONFLICT (id) DO UPDATE SET
    user_id = EXCLUDED.user_id,
    institution = EXCLUDED.institution,
    degree = EXCLUDED.degree,
    field = EXCLUDED.field,
    start_date = EXCLUDED.start_date,
    end_date = EXCLUDED.end_date,
    description = EXCLUDED.description;

INSERT INTO certifications (id, user_id, name, issuer, issue_date, expiry_date, credential_id, url)
VALUES
    (
        '80000000-0000-7000-8000-000000000001',
        '10000000-0000-7000-8000-000000000001',
        'AWS Certified Developer - Associate',
        'Amazon Web Services',
        DATE '2022-04-15',
        DATE '2025-04-15',
        'SAMPLE-AWS-DEV-001',
        ''
    ),
    (
        '80000000-0000-7000-8000-000000000002',
        '10000000-0000-7000-8000-000000000001',
        'Certified Kubernetes Administrator (CKA)',
        'The Linux Foundation',
        DATE '2025-03-10',
        DATE '2028-03-10',
        'SAMPLE-CKA-002',
        ''
    )
ON CONFLICT (id) DO UPDATE SET
    user_id = EXCLUDED.user_id,
    name = EXCLUDED.name,
    issuer = EXCLUDED.issuer,
    issue_date = EXCLUDED.issue_date,
    expiry_date = EXCLUDED.expiry_date,
    credential_id = EXCLUDED.credential_id,
    url = EXCLUDED.url;

INSERT INTO languages (id, user_id, name, proficiency)
VALUES
    ('90000000-0000-7000-8000-000000000001', '10000000-0000-7000-8000-000000000001', 'English', 'native'),
    ('90000000-0000-7000-8000-000000000002', '10000000-0000-7000-8000-000000000001', 'Mandarin Chinese', 'advanced')
ON CONFLICT (id) DO UPDATE SET
    user_id = EXCLUDED.user_id,
    name = EXCLUDED.name,
    proficiency = EXCLUDED.proficiency;

COMMIT;