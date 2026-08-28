INSERT INTO users (id, email, display_name, role)
VALUES
    ('00000000-0000-4000-8000-000000000001', 'member@example.com', 'Demo Member', 'member'),
    ('00000000-0000-4000-8000-000000000002', 'admin@example.com', 'Demo Admin', 'admin')
ON CONFLICT (email) DO UPDATE SET
    display_name = EXCLUDED.display_name,
    role = EXCLUDED.role,
    updated_at = now();

INSERT INTO resources (id, title, type, description, file_path, tags, author_id, status, published_at)
VALUES (
    '10000000-0000-4000-8000-000000000001',
    'Organization Workbench Guide',
    'document',
    'A short guide to the organization workbench and its shared resources.',
    'README.md',
    ARRAY['onboarding', 'guide'],
    '00000000-0000-4000-8000-000000000002',
    'published',
    now()
)
ON CONFLICT (id) DO UPDATE SET
    title = EXCLUDED.title,
    type = EXCLUDED.type,
    description = EXCLUDED.description,
    file_path = EXCLUDED.file_path,
    tags = EXCLUDED.tags,
    author_id = EXCLUDED.author_id,
    status = EXCLUDED.status,
    published_at = EXCLUDED.published_at,
    updated_at = now();

INSERT INTO posts (id, title, summary, body, tags, author_id, status, published_at)
VALUES (
    '20000000-0000-4000-8000-000000000001',
    'Welcome to the Organization Workbench',
    'The first development update for our shared workbench.',
    'This post introduces the shared resources, development updates, and AI experiments available to members.',
    ARRAY['announcement', 'workbench'],
    '00000000-0000-4000-8000-000000000002',
    'published',
    now()
)
ON CONFLICT (id) DO UPDATE SET
    title = EXCLUDED.title,
    summary = EXCLUDED.summary,
    body = EXCLUDED.body,
    tags = EXCLUDED.tags,
    author_id = EXCLUDED.author_id,
    status = EXCLUDED.status,
    published_at = EXCLUDED.published_at,
    updated_at = now();

INSERT INTO ai_products (id, name, summary, url, tags, status, sort_order, published_at)
VALUES (
    '30000000-0000-4000-8000-000000000001',
    'Dify Assistant',
    'An internal Dify-powered assistant for experimenting with organization knowledge.',
    'https://dify.ai',
    ARRAY['ai', 'dify'],
    'published',
    10,
    now()
)
ON CONFLICT (id) DO UPDATE SET
    name = EXCLUDED.name,
    summary = EXCLUDED.summary,
    url = EXCLUDED.url,
    tags = EXCLUDED.tags,
    status = EXCLUDED.status,
    sort_order = EXCLUDED.sort_order,
    published_at = EXCLUDED.published_at,
    updated_at = now();

INSERT INTO user_layouts (user_id, layout)
VALUES (
    '00000000-0000-4000-8000-000000000001',
    '{"modules":["overview","resources","posts","ai-lab"],"hidden":[],"density":"comfortable","theme":"light"}'::jsonb
)
ON CONFLICT (user_id) DO UPDATE SET
    layout = EXCLUDED.layout,
    updated_at = now();
