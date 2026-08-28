CREATE INDEX resources_status_published_at_idx ON resources (status, published_at DESC);
CREATE INDEX posts_status_published_at_idx ON posts (status, published_at DESC);
CREATE INDEX ai_products_status_published_at_idx ON ai_products (status, published_at DESC);
CREATE INDEX ai_products_sort_order_idx ON ai_products (sort_order ASC, name ASC);

CREATE INDEX resources_tags_gin_idx ON resources USING GIN (tags);
CREATE INDEX posts_tags_gin_idx ON posts USING GIN (tags);
CREATE INDEX ai_products_tags_gin_idx ON ai_products USING GIN (tags);
