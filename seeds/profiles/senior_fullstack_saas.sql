-- Fictional CareerMesh sample profile; all employers, projects, and identifiers are illustrative.
-- Run from the repository root only against a dedicated development/test database:
-- psql "$TEST_DATABASE_URL" -v ON_ERROR_STOP=1 -f seeds/profiles/senior_fullstack_saas.sql
-- The seeded account has an intentionally invalid password hash and cannot be used to log in.

BEGIN;

INSERT INTO users (id, email, password_hash)
VALUES (
    '11000000-0000-7000-8000-000000000001',
    'elena.park@example.com',
    '!disabled-seed-account-no-password!'
)
ON CONFLICT (email) DO UPDATE SET
    password_hash = EXCLUDED.password_hash
WHERE users.id = EXCLUDED.id;

-- SQL seeds bypass user.Service.CreateWithHashedPassword, so assign only the default role.
INSERT INTO user_roles (user_id, role_id)
VALUES (
    '11000000-0000-7000-8000-000000000001',
    (SELECT id FROM roles WHERE name = 'user')
)
ON CONFLICT (user_id) DO UPDATE SET
    role_id = EXCLUDED.role_id;

INSERT INTO profiles (id, user_id, name, headline, location, phone, website, github, linkedin)
VALUES (
    '21000000-0000-7000-8000-000000000001',
    '11000000-0000-7000-8000-000000000001',
    'Elena Park',
    'Senior Full-Stack SaaS Engineer | TypeScript, React, Node.js',
    'Seattle, Washington',
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
        '31000000-0000-7000-8000-000000000001',
        '11000000-0000-7000-8000-000000000001',
        'Juniper Grove Learning', '', 'Software Engineer', 'full_time', 'Seattle, Washington',
        DATE '2017-06-12', DATE '2019-08-30',
        'Built web features for a subscription learning platform, working across React interfaces and Node.js APIs with designers, educators, and customer support.'
    ),
    (
        '31000000-0000-7000-8000-000000000002',
        '11000000-0000-7000-8000-000000000001',
        'Harborwell Commerce', '', 'Full-Stack Engineer', 'full_time', 'Seattle, Washington',
        DATE '2019-09-03', DATE '2022-10-14',
        'Delivered customer and operations workflows for a multi-tenant commerce SaaS product, integrating React and TypeScript clients with REST APIs, PostgreSQL, and automated release pipelines.'
    ),
    (
        '31000000-0000-7000-8000-000000000003',
        '11000000-0000-7000-8000-000000000001',
        'Northwind Orchard Software', '', 'Senior Full-Stack SaaS Engineer', 'full_time', 'Remote',
        DATE '2022-10-17', NULL,
        'Leads delivery of product capabilities across a SaaS application stack, partnering with product and design on usability, API contracts, data models, testing strategy, and incremental releases.'
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
    '31000000-0000-7000-8000-000000000001',
    '31000000-0000-7000-8000-000000000002',
    '31000000-0000-7000-8000-000000000003'
);

INSERT INTO experience_bullets (id, experience_id, content, position)
VALUES
    ('41000000-0000-7000-8000-000000000001', '31000000-0000-7000-8000-000000000001', 'Built responsive course and account screens in React, translating design prototypes into accessible flows used across desktop and mobile.', 1),
    ('41000000-0000-7000-8000-000000000002', '31000000-0000-7000-8000-000000000001', 'Extended Node.js REST endpoints and PostgreSQL models for enrollment and progress tracking, keeping client and server validation consistent.', 2),
    ('41000000-0000-7000-8000-000000000003', '31000000-0000-7000-8000-000000000001', 'Added component and API tests to the continuous integration workflow, catching regressions before weekly product releases.', 3),
    ('41000000-0000-7000-8000-000000000004', '31000000-0000-7000-8000-000000000002', 'Delivered a shared TypeScript component library with product design, reducing duplicated form and table behavior across three customer workflows.', 1),
    ('41000000-0000-7000-8000-000000000005', '31000000-0000-7000-8000-000000000002', 'Designed REST API changes and PostgreSQL migrations for subscription changes and invoice history, coordinating backward-compatible rollout with support and QA.', 2),
    ('41000000-0000-7000-8000-000000000006', '31000000-0000-7000-8000-000000000002', 'Containerized local development services with Docker and improved automated test coverage for checkout and account-management paths.', 3),
    ('41000000-0000-7000-8000-000000000007', '31000000-0000-7000-8000-000000000003', 'Shipped a Next.js customer workspace that consolidates usage, invoices, and team settings, informed by usability sessions with customer success.', 1),
    ('41000000-0000-7000-8000-000000000008', '31000000-0000-7000-8000-000000000003', 'Defined API contracts with backend teammates and delivered Node.js endpoints for metered usage summaries with paginated PostgreSQL queries.', 2),
    ('41000000-0000-7000-8000-000000000009', '31000000-0000-7000-8000-000000000003', 'Introduced preview deployments and required test checks in CI/CD, giving product and design partners earlier feedback on proposed changes.', 3)
ON CONFLICT (id) DO UPDATE SET
    experience_id = EXCLUDED.experience_id,
    content = EXCLUDED.content,
    position = EXCLUDED.position;

INSERT INTO skills (id, user_id, name, category)
VALUES
    ('51000000-0000-7000-8000-000000000001', '11000000-0000-7000-8000-000000000001', 'TypeScript', 'Programming Languages'),
    ('51000000-0000-7000-8000-000000000002', '11000000-0000-7000-8000-000000000001', 'JavaScript', 'Programming Languages'),
    ('51000000-0000-7000-8000-000000000003', '11000000-0000-7000-8000-000000000001', 'React', 'Frontend'),
    ('51000000-0000-7000-8000-000000000004', '11000000-0000-7000-8000-000000000001', 'Next.js', 'Frontend'),
    ('51000000-0000-7000-8000-000000000005', '11000000-0000-7000-8000-000000000001', 'Node.js', 'Backend'),
    ('51000000-0000-7000-8000-000000000006', '11000000-0000-7000-8000-000000000001', 'REST APIs', 'API Design'),
    ('51000000-0000-7000-8000-000000000007', '11000000-0000-7000-8000-000000000001', 'PostgreSQL', 'Databases'),
    ('51000000-0000-7000-8000-000000000008', '11000000-0000-7000-8000-000000000001', 'Automated Testing', 'Quality Engineering'),
    ('51000000-0000-7000-8000-000000000009', '11000000-0000-7000-8000-000000000001', 'Docker', 'Containers'),
    ('51000000-0000-7000-8000-000000000010', '11000000-0000-7000-8000-000000000001', 'CI/CD', 'Delivery')
ON CONFLICT (id) DO UPDATE SET
    user_id = EXCLUDED.user_id,
    name = EXCLUDED.name,
    category = EXCLUDED.category;

INSERT INTO projects (id, user_id, name, description, url, repository_url, start_date, end_date)
VALUES
    (
        '61000000-0000-7000-8000-000000000001',
        '11000000-0000-7000-8000-000000000001',
        'Harborwell Subscription Console',
        'A multi-tenant account console for subscription changes, invoice history, and team access. The project joined React and TypeScript workflows to versioned Node.js REST endpoints and PostgreSQL records, with automated coverage for billing state transitions.',
        '', '', DATE '2020-02-01', DATE '2021-07-30'
    ),
    (
        '61000000-0000-7000-8000-000000000002',
        '11000000-0000-7000-8000-000000000001',
        'Northwind Usage Workspace',
        'A SaaS workspace for reviewing product usage, managing seats, and exporting account activity. Built with Next.js, React, TypeScript, and Node.js APIs, supported by PostgreSQL query pagination, Docker-based local services, and preview deployments in CI/CD.',
        '', '', DATE '2023-03-01', DATE '2024-11-29'
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
    '61000000-0000-7000-8000-000000000001',
    '61000000-0000-7000-8000-000000000002'
);

INSERT INTO project_skills (project_id, skill_id)
VALUES
    ('61000000-0000-7000-8000-000000000001', '51000000-0000-7000-8000-000000000001'),
    ('61000000-0000-7000-8000-000000000001', '51000000-0000-7000-8000-000000000003'),
    ('61000000-0000-7000-8000-000000000001', '51000000-0000-7000-8000-000000000005'),
    ('61000000-0000-7000-8000-000000000001', '51000000-0000-7000-8000-000000000006'),
    ('61000000-0000-7000-8000-000000000001', '51000000-0000-7000-8000-000000000007'),
    ('61000000-0000-7000-8000-000000000001', '51000000-0000-7000-8000-000000000008'),
    ('61000000-0000-7000-8000-000000000002', '51000000-0000-7000-8000-000000000001'),
    ('61000000-0000-7000-8000-000000000002', '51000000-0000-7000-8000-000000000003'),
    ('61000000-0000-7000-8000-000000000002', '51000000-0000-7000-8000-000000000004'),
    ('61000000-0000-7000-8000-000000000002', '51000000-0000-7000-8000-000000000005'),
    ('61000000-0000-7000-8000-000000000002', '51000000-0000-7000-8000-000000000007'),
    ('61000000-0000-7000-8000-000000000002', '51000000-0000-7000-8000-000000000008'),
    ('61000000-0000-7000-8000-000000000002', '51000000-0000-7000-8000-000000000009'),
    ('61000000-0000-7000-8000-000000000002', '51000000-0000-7000-8000-000000000010')
ON CONFLICT (project_id, skill_id) DO NOTHING;

DELETE FROM experience_skills
WHERE experience_id IN (
    '31000000-0000-7000-8000-000000000001',
    '31000000-0000-7000-8000-000000000002',
    '31000000-0000-7000-8000-000000000003'
);

INSERT INTO experience_skills (experience_id, skill_id)
VALUES
    ('31000000-0000-7000-8000-000000000001', '51000000-0000-7000-8000-000000000001'),
    ('31000000-0000-7000-8000-000000000001', '51000000-0000-7000-8000-000000000002'),
    ('31000000-0000-7000-8000-000000000001', '51000000-0000-7000-8000-000000000003'),
    ('31000000-0000-7000-8000-000000000001', '51000000-0000-7000-8000-000000000005'),
    ('31000000-0000-7000-8000-000000000001', '51000000-0000-7000-8000-000000000006'),
    ('31000000-0000-7000-8000-000000000001', '51000000-0000-7000-8000-000000000007'),
    ('31000000-0000-7000-8000-000000000002', '51000000-0000-7000-8000-000000000001'),
    ('31000000-0000-7000-8000-000000000002', '51000000-0000-7000-8000-000000000003'),
    ('31000000-0000-7000-8000-000000000002', '51000000-0000-7000-8000-000000000005'),
    ('31000000-0000-7000-8000-000000000002', '51000000-0000-7000-8000-000000000006'),
    ('31000000-0000-7000-8000-000000000002', '51000000-0000-7000-8000-000000000007'),
    ('31000000-0000-7000-8000-000000000002', '51000000-0000-7000-8000-000000000008'),
    ('31000000-0000-7000-8000-000000000002', '51000000-0000-7000-8000-000000000009'),
    ('31000000-0000-7000-8000-000000000002', '51000000-0000-7000-8000-000000000010'),
    ('31000000-0000-7000-8000-000000000003', '51000000-0000-7000-8000-000000000001'),
    ('31000000-0000-7000-8000-000000000003', '51000000-0000-7000-8000-000000000003'),
    ('31000000-0000-7000-8000-000000000003', '51000000-0000-7000-8000-000000000004'),
    ('31000000-0000-7000-8000-000000000003', '51000000-0000-7000-8000-000000000005'),
    ('31000000-0000-7000-8000-000000000003', '51000000-0000-7000-8000-000000000006'),
    ('31000000-0000-7000-8000-000000000003', '51000000-0000-7000-8000-000000000007'),
    ('31000000-0000-7000-8000-000000000003', '51000000-0000-7000-8000-000000000008'),
    ('31000000-0000-7000-8000-000000000003', '51000000-0000-7000-8000-000000000010')
ON CONFLICT (experience_id, skill_id) DO NOTHING;

INSERT INTO education (id, user_id, institution, degree, field, start_date, end_date, description)
VALUES (
    '71000000-0000-7000-8000-000000000001',
    '11000000-0000-7000-8000-000000000001',
    'Cascadia State University',
    'Bachelor of Science',
    'Information Systems',
    DATE '2013-09-23',
    DATE '2017-06-09',
    'Coursework included web application development, database design, human-computer interaction, software project management, and information security.'
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
        '81000000-0000-7000-8000-000000000001',
        '11000000-0000-7000-8000-000000000001',
        'AWS Certified Cloud Practitioner',
        'Amazon Web Services',
        DATE '2023-05-12',
        DATE '2026-05-12',
        'SAMPLE-AWS-CCP-101',
        ''
    ),
    (
        '81000000-0000-7000-8000-000000000002',
        '11000000-0000-7000-8000-000000000001',
        'Professional Scrum Developer I',
        'Scrum.org',
        DATE '2021-09-18',
        NULL,
        'SAMPLE-PSD1-102',
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
    ('91000000-0000-7000-8000-000000000001', '11000000-0000-7000-8000-000000000001', 'English', 'native'),
    ('91000000-0000-7000-8000-000000000002', '11000000-0000-7000-8000-000000000001', 'Spanish', 'upper_intermediate')
ON CONFLICT (id) DO UPDATE SET
    user_id = EXCLUDED.user_id,
    name = EXCLUDED.name,
    proficiency = EXCLUDED.proficiency;

COMMIT;