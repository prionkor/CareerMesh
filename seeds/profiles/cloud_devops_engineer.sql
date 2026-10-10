-- Fictional CareerMesh sample profile; all employers, projects, and identifiers are illustrative.
-- Run from the repository root only against a dedicated development/test database:
-- psql "$TEST_DATABASE_URL" -v ON_ERROR_STOP=1 -f seeds/profiles/cloud_devops_engineer.sql
-- The seeded account has an intentionally invalid password hash and cannot be used to log in.

BEGIN;

INSERT INTO users (id, email, password_hash)
VALUES (
    '12000000-0000-7000-8000-000000000001',
    'daniel.alvarez@example.com',
    '!disabled-seed-account-no-password!'
)
ON CONFLICT (email) DO UPDATE SET
    password_hash = EXCLUDED.password_hash
WHERE users.id = EXCLUDED.id;

-- SQL seeds bypass user.Service.CreateWithHashedPassword, so assign only the default role.
INSERT INTO user_roles (user_id, role_id)
VALUES (
    '12000000-0000-7000-8000-000000000001',
    (SELECT id FROM roles WHERE name = 'user')
)
ON CONFLICT (user_id) DO UPDATE SET
    role_id = EXCLUDED.role_id;

INSERT INTO profiles (id, user_id, name, headline, location, phone, website, github, linkedin)
VALUES (
    '22000000-0000-7000-8000-000000000001',
    '12000000-0000-7000-8000-000000000001',
    'Daniel Alvarez',
    'Cloud / DevOps Engineer | AWS, Terraform, Kubernetes, SRE',
    'Denver, Colorado',
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
        '32000000-0000-7000-8000-000000000001',
        '12000000-0000-7000-8000-000000000001',
        'Cedarline Cloud Services', '', 'Systems Engineer', 'full_time', 'Denver, Colorado',
        DATE '2017-06-05', DATE '2019-07-31',
        'Operated Linux-based customer environments and supported application deployments for a regional cloud hosting provider, building practical experience in networking, scripting, monitoring, and incident response.'
    ),
    (
        '32000000-0000-7000-8000-000000000002',
        '12000000-0000-7000-8000-000000000001',
        'Ternary Harbor Logistics', '', 'Cloud Infrastructure Engineer', 'full_time', 'Denver, Colorado',
        DATE '2019-08-05', DATE '2022-12-30',
        'Automated AWS and Azure cloud foundations and delivery workflows for logistics and analytics services, working with application teams on infrastructure changes, release safety, capacity planning, and production support.'
    ),
    (
        '32000000-0000-7000-8000-000000000003',
        '12000000-0000-7000-8000-000000000001',
        'Blue Mesa Health Systems', '', 'Senior Cloud / DevOps Engineer', 'full_time', 'Remote',
        DATE '2023-01-03', NULL,
        'Leads cloud platform reliability and deployment improvements for healthcare software teams, emphasizing infrastructure as code, Kubernetes operations, observability, incident learning, and pragmatic SRE practices.'
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
    '32000000-0000-7000-8000-000000000001',
    '32000000-0000-7000-8000-000000000002',
    '32000000-0000-7000-8000-000000000003'
);

INSERT INTO experience_bullets (id, experience_id, content, position)
VALUES
    ('42000000-0000-7000-8000-000000000001', '32000000-0000-7000-8000-000000000001', 'Maintained Linux hosts, deployment scripts, and service checks for customer applications, resolving availability issues with clear handoffs to engineering teams.', 1),
    ('42000000-0000-7000-8000-000000000002', '32000000-0000-7000-8000-000000000001', 'Added host and service monitoring with documented alert thresholds, helping the on-call rotation distinguish application faults from capacity problems.', 2),
    ('42000000-0000-7000-8000-000000000003', '32000000-0000-7000-8000-000000000001', 'Standardized routine environment setup and backup checks with shell automation, reducing repeated manual steps during customer onboarding.', 3),
    ('42000000-0000-7000-8000-000000000004', '32000000-0000-7000-8000-000000000002', 'Created Terraform modules for AWS networking, IAM, and service environments, reducing a typical new environment setup from several days to a few hours.', 1),
    ('42000000-0000-7000-8000-000000000005', '32000000-0000-7000-8000-000000000002', 'Moved containerized services to managed Kubernetes with documented resource requests, health checks, and rollback procedures.', 2),
    ('42000000-0000-7000-8000-000000000006', '32000000-0000-7000-8000-000000000002', 'Introduced CI/CD checks for infrastructure plans and application images, giving reviewers repeatable previews before production changes.', 3),
    ('42000000-0000-7000-8000-000000000007', '32000000-0000-7000-8000-000000000003', 'Built GitOps deployment workflows for Kubernetes services, separating reviewed configuration changes from runtime promotion across environments.', 1),
    ('42000000-0000-7000-8000-000000000008', '32000000-0000-7000-8000-000000000003', 'Improved service dashboards and actionable alerts using infrastructure and application telemetry, then used incident reviews to prioritize recurring failure modes.', 2),
    ('42000000-0000-7000-8000-000000000009', '32000000-0000-7000-8000-000000000003', 'Coordinated incident response with product engineering and security, documenting recovery steps and validating follow-up reliability work.', 3)
ON CONFLICT (id) DO UPDATE SET
    experience_id = EXCLUDED.experience_id,
    content = EXCLUDED.content,
    position = EXCLUDED.position;

INSERT INTO skills (id, user_id, name, category)
VALUES
    ('52000000-0000-7000-8000-000000000001', '12000000-0000-7000-8000-000000000001', 'AWS', 'Cloud Platforms'),
    ('52000000-0000-7000-8000-000000000002', '12000000-0000-7000-8000-000000000001', 'Azure', 'Cloud Platforms'),
    ('52000000-0000-7000-8000-000000000003', '12000000-0000-7000-8000-000000000001', 'Terraform', 'Infrastructure as Code'),
    ('52000000-0000-7000-8000-000000000004', '12000000-0000-7000-8000-000000000001', 'Kubernetes', 'Platform Engineering'),
    ('52000000-0000-7000-8000-000000000005', '12000000-0000-7000-8000-000000000001', 'Docker', 'Containers'),
    ('52000000-0000-7000-8000-000000000006', '12000000-0000-7000-8000-000000000001', 'GitOps', 'Delivery'),
    ('52000000-0000-7000-8000-000000000007', '12000000-0000-7000-8000-000000000001', 'CI/CD', 'Delivery'),
    ('52000000-0000-7000-8000-000000000008', '12000000-0000-7000-8000-000000000001', 'Linux', 'Operating Systems'),
    ('52000000-0000-7000-8000-000000000009', '12000000-0000-7000-8000-000000000001', 'Monitoring', 'Observability'),
    ('52000000-0000-7000-8000-000000000010', '12000000-0000-7000-8000-000000000001', 'Observability', 'Reliability Engineering'),
    ('52000000-0000-7000-8000-000000000011', '12000000-0000-7000-8000-000000000001', 'Incident Response', 'Operations'),
    ('52000000-0000-7000-8000-000000000012', '12000000-0000-7000-8000-000000000001', 'SRE Practices', 'Reliability Engineering')
ON CONFLICT (id) DO UPDATE SET
    user_id = EXCLUDED.user_id,
    name = EXCLUDED.name,
    category = EXCLUDED.category;

INSERT INTO projects (id, user_id, name, description, url, repository_url, start_date, end_date)
VALUES
    (
        '62000000-0000-7000-8000-000000000001',
        '12000000-0000-7000-8000-000000000001',
        'Ternary Harbor Cloud Foundations',
        'A reusable cloud foundation for logistics services, provisioning AWS network boundaries and IAM roles alongside Azure virtual networks and role assignments with Terraform. Change plans ran in CI/CD and included reviewable outputs and documented recovery steps.',
        '', '', DATE '2020-01-06', DATE '2021-10-29'
    ),
    (
        '62000000-0000-7000-8000-000000000002',
        '12000000-0000-7000-8000-000000000001',
        'Blue Mesa Kubernetes Delivery Platform',
        'A managed Kubernetes delivery path for product teams, combining Docker image checks, GitOps configuration promotion, deployment health signals, and rollback guidance. The platform paired service-level dashboards with incident review actions to improve operational consistency.',
        '', '', DATE '2023-04-03', DATE '2025-02-28'
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
    '62000000-0000-7000-8000-000000000001',
    '62000000-0000-7000-8000-000000000002'
);

INSERT INTO project_skills (project_id, skill_id)
VALUES
    ('62000000-0000-7000-8000-000000000001', '52000000-0000-7000-8000-000000000001'),
    ('62000000-0000-7000-8000-000000000001', '52000000-0000-7000-8000-000000000002'),
    ('62000000-0000-7000-8000-000000000001', '52000000-0000-7000-8000-000000000003'),
    ('62000000-0000-7000-8000-000000000001', '52000000-0000-7000-8000-000000000007'),
    ('62000000-0000-7000-8000-000000000001', '52000000-0000-7000-8000-000000000008'),
    ('62000000-0000-7000-8000-000000000001', '52000000-0000-7000-8000-000000000009'),
    ('62000000-0000-7000-8000-000000000002', '52000000-0000-7000-8000-000000000003'),
    ('62000000-0000-7000-8000-000000000002', '52000000-0000-7000-8000-000000000004'),
    ('62000000-0000-7000-8000-000000000002', '52000000-0000-7000-8000-000000000005'),
    ('62000000-0000-7000-8000-000000000002', '52000000-0000-7000-8000-000000000006'),
    ('62000000-0000-7000-8000-000000000002', '52000000-0000-7000-8000-000000000007'),
    ('62000000-0000-7000-8000-000000000002', '52000000-0000-7000-8000-000000000009'),
    ('62000000-0000-7000-8000-000000000002', '52000000-0000-7000-8000-000000000010'),
    ('62000000-0000-7000-8000-000000000002', '52000000-0000-7000-8000-000000000011'),
    ('62000000-0000-7000-8000-000000000002', '52000000-0000-7000-8000-000000000012')
ON CONFLICT (project_id, skill_id) DO NOTHING;

DELETE FROM experience_skills
WHERE experience_id IN (
    '32000000-0000-7000-8000-000000000001',
    '32000000-0000-7000-8000-000000000002',
    '32000000-0000-7000-8000-000000000003'
);

INSERT INTO experience_skills (experience_id, skill_id)
VALUES
    ('32000000-0000-7000-8000-000000000001', '52000000-0000-7000-8000-000000000008'),
    ('32000000-0000-7000-8000-000000000001', '52000000-0000-7000-8000-000000000009'),
    ('32000000-0000-7000-8000-000000000001', '52000000-0000-7000-8000-000000000011'),
    ('32000000-0000-7000-8000-000000000002', '52000000-0000-7000-8000-000000000001'),
    ('32000000-0000-7000-8000-000000000002', '52000000-0000-7000-8000-000000000002'),
    ('32000000-0000-7000-8000-000000000002', '52000000-0000-7000-8000-000000000003'),
    ('32000000-0000-7000-8000-000000000002', '52000000-0000-7000-8000-000000000004'),
    ('32000000-0000-7000-8000-000000000002', '52000000-0000-7000-8000-000000000005'),
    ('32000000-0000-7000-8000-000000000002', '52000000-0000-7000-8000-000000000007'),
    ('32000000-0000-7000-8000-000000000002', '52000000-0000-7000-8000-000000000009'),
    ('32000000-0000-7000-8000-000000000003', '52000000-0000-7000-8000-000000000001'),
    ('32000000-0000-7000-8000-000000000003', '52000000-0000-7000-8000-000000000002'),
    ('32000000-0000-7000-8000-000000000003', '52000000-0000-7000-8000-000000000003'),
    ('32000000-0000-7000-8000-000000000003', '52000000-0000-7000-8000-000000000004'),
    ('32000000-0000-7000-8000-000000000003', '52000000-0000-7000-8000-000000000006'),
    ('32000000-0000-7000-8000-000000000003', '52000000-0000-7000-8000-000000000007'),
    ('32000000-0000-7000-8000-000000000003', '52000000-0000-7000-8000-000000000009'),
    ('32000000-0000-7000-8000-000000000003', '52000000-0000-7000-8000-000000000010'),
    ('32000000-0000-7000-8000-000000000003', '52000000-0000-7000-8000-000000000011'),
    ('32000000-0000-7000-8000-000000000003', '52000000-0000-7000-8000-000000000012')
ON CONFLICT (experience_id, skill_id) DO NOTHING;

INSERT INTO education (id, user_id, institution, degree, field, start_date, end_date, description)
VALUES (
    '72000000-0000-7000-8000-000000000001',
    '12000000-0000-7000-8000-000000000001',
    'Front Range Polytechnic',
    'Bachelor of Science',
    'Information Technology',
    DATE '2013-08-26',
    DATE '2017-05-12',
    'Coursework included operating systems, computer networking, cloud systems, scripting, and information security.'
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
        '82000000-0000-7000-8000-000000000001',
        '12000000-0000-7000-8000-000000000001',
        'AWS Certified SysOps Administrator - Associate',
        'Amazon Web Services',
        DATE '2024-02-20',
        DATE '2027-02-20',
        'SAMPLE-AWS-SOA-201',
        ''
    ),
    (
        '82000000-0000-7000-8000-000000000002',
        '12000000-0000-7000-8000-000000000001',
        'Certified Kubernetes Administrator (CKA)',
        'The Linux Foundation',
        DATE '2025-06-16',
        DATE '2028-06-16',
        'SAMPLE-CKA-202',
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
    ('92000000-0000-7000-8000-000000000001', '12000000-0000-7000-8000-000000000001', 'English', 'native'),
    ('92000000-0000-7000-8000-000000000002', '12000000-0000-7000-8000-000000000001', 'Spanish', 'advanced')
ON CONFLICT (id) DO UPDATE SET
    user_id = EXCLUDED.user_id,
    name = EXCLUDED.name,
    proficiency = EXCLUDED.proficiency;

COMMIT;