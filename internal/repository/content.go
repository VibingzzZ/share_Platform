package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"share-platform/internal/model"
)

// ContentRepository provides published content reads and administrator writes.
type ContentRepository struct {
	pool *pgxpool.Pool
}

func NewContent(pool *pgxpool.Pool) *ContentRepository {
	return &ContentRepository{pool: pool}
}

func (r *ContentRepository) ListResources(ctx context.Context, filter model.ListFilter) (model.Page[model.Resource], error) {
	filter = filter.Normalized()
	const where = "status = 'published' AND ($1 = '' OR type = $1) AND ($2 = '' OR $2 = ANY(tags))"
	var total int
	if err := r.pool.QueryRow(ctx, "SELECT count(*) FROM resources WHERE "+where, filter.Type, filter.Tag).Scan(&total); err != nil {
		return model.Page[model.Resource]{}, fmt.Errorf("count resources: %w", err)
	}
	rows, err := r.pool.Query(ctx, `SELECT id, title, type, description, url, file_path, tags, author_id, status, published_at, created_at, updated_at
		FROM resources WHERE `+where+` ORDER BY published_at DESC, created_at DESC LIMIT $3 OFFSET $4`, filter.Type, filter.Tag, filter.PageSize, (filter.Page-1)*filter.PageSize)
	if err != nil {
		return model.Page[model.Resource]{}, fmt.Errorf("list resources: %w", err)
	}
	defer rows.Close()
	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Resource])
	if err != nil {
		return model.Page[model.Resource]{}, fmt.Errorf("scan resources: %w", err)
	}
	return model.Page[model.Resource]{Items: items, Page: filter.Page, PageSize: filter.PageSize, Total: total}, nil
}

func (r *ContentRepository) ListPosts(ctx context.Context, filter model.ListFilter) (model.Page[model.Post], error) {
	filter = filter.Normalized()
	const where = "status = 'published' AND ($1 = '' OR $1 = ANY(tags))"
	var total int
	if err := r.pool.QueryRow(ctx, "SELECT count(*) FROM posts WHERE "+where, filter.Tag).Scan(&total); err != nil {
		return model.Page[model.Post]{}, fmt.Errorf("count posts: %w", err)
	}
	rows, err := r.pool.Query(ctx, `SELECT id, title, summary, body, cover_image, tags, author_id, status, published_at, created_at, updated_at
		FROM posts WHERE `+where+` ORDER BY published_at DESC, created_at DESC LIMIT $2 OFFSET $3`, filter.Tag, filter.PageSize, (filter.Page-1)*filter.PageSize)
	if err != nil {
		return model.Page[model.Post]{}, fmt.Errorf("list posts: %w", err)
	}
	defer rows.Close()
	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Post])
	if err != nil {
		return model.Page[model.Post]{}, fmt.Errorf("scan posts: %w", err)
	}
	return model.Page[model.Post]{Items: items, Page: filter.Page, PageSize: filter.PageSize, Total: total}, nil
}

func (r *ContentRepository) ListAIProducts(ctx context.Context, filter model.ListFilter) (model.Page[model.AIProduct], error) {
	filter = filter.Normalized()
	const where = "status = 'published' AND ($1 = '' OR $1 = ANY(tags))"
	var total int
	if err := r.pool.QueryRow(ctx, "SELECT count(*) FROM ai_products WHERE "+where, filter.Tag).Scan(&total); err != nil {
		return model.Page[model.AIProduct]{}, fmt.Errorf("count AI products: %w", err)
	}
	rows, err := r.pool.Query(ctx, `SELECT id, name, summary, url, cover_image, tags, status, sort_order, published_at, created_at, updated_at
		FROM ai_products WHERE `+where+` ORDER BY sort_order ASC, name ASC LIMIT $2 OFFSET $3`, filter.Tag, filter.PageSize, (filter.Page-1)*filter.PageSize)
	if err != nil {
		return model.Page[model.AIProduct]{}, fmt.Errorf("list AI products: %w", err)
	}
	defer rows.Close()
	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.AIProduct])
	if err != nil {
		return model.Page[model.AIProduct]{}, fmt.Errorf("scan AI products: %w", err)
	}
	return model.Page[model.AIProduct]{Items: items, Page: filter.Page, PageSize: filter.PageSize, Total: total}, nil
}

func (r *ContentRepository) Dashboard(ctx context.Context) (model.Dashboard, error) {
	resources, err := r.ListResources(ctx, model.ListFilter{Page: 1, PageSize: 5})
	if err != nil {
		return model.Dashboard{}, err
	}
	posts, err := r.ListPosts(ctx, model.ListFilter{Page: 1, PageSize: 5})
	if err != nil {
		return model.Dashboard{}, err
	}
	products, err := r.ListAIProducts(ctx, model.ListFilter{Page: 1, PageSize: 5})
	if err != nil {
		return model.Dashboard{}, err
	}
	return model.Dashboard{Resources: resources.Items, Posts: posts.Items, AIProducts: products.Items}, nil
}

func (r *ContentRepository) CreateResource(ctx context.Context, resource model.Resource) (model.Resource, error) {
	return r.resource(ctx, `INSERT INTO resources (title, type, description, url, file_path, tags, author_id, status, published_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, CASE WHEN $8 = 'published' THEN now() END)
		RETURNING id, title, type, description, url, file_path, tags, author_id, status, published_at, created_at, updated_at`, resource)
}

func (r *ContentRepository) UpdateResource(ctx context.Context, id string, resource model.Resource) (model.Resource, error) {
	return r.resource(ctx, `UPDATE resources SET title = $1, type = $2, description = $3, url = $4, file_path = $5, tags = $6,
		status = $8, published_at = CASE WHEN $8 = 'published' THEN COALESCE(published_at, now()) ELSE NULL END, updated_at = now()
		WHERE id = $9 RETURNING id, title, type, description, url, file_path, tags, author_id, status, published_at, created_at, updated_at`, resource, id)
}

func (r *ContentRepository) resource(ctx context.Context, query string, resource model.Resource, extra ...string) (model.Resource, error) {
	args := []any{resource.Title, resource.Type, resource.Description, resource.URL, resource.FilePath, resource.Tags, resource.AuthorID, resource.Status}
	for _, value := range extra {
		args = append(args, value)
	}
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return model.Resource{}, fmt.Errorf("write resource: %w", err)
	}
	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Resource])
	if err != nil {
		return model.Resource{}, fmt.Errorf("write resource: %w", err)
	}
	return item, nil
}

func (r *ContentRepository) DeleteResource(ctx context.Context, id string) error {
	return r.delete(ctx, "resources", id)
}

func (r *ContentRepository) CreatePost(ctx context.Context, post model.Post) (model.Post, error) {
	return r.post(ctx, `INSERT INTO posts (title, summary, body, cover_image, tags, author_id, status, published_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, CASE WHEN $7 = 'published' THEN now() END)
		RETURNING id, title, summary, body, cover_image, tags, author_id, status, published_at, created_at, updated_at`, post)
}

func (r *ContentRepository) UpdatePost(ctx context.Context, id string, post model.Post) (model.Post, error) {
	return r.post(ctx, `UPDATE posts SET title = $1, summary = $2, body = $3, cover_image = $4, tags = $5, status = $7,
		published_at = CASE WHEN $7 = 'published' THEN COALESCE(published_at, now()) ELSE NULL END, updated_at = now()
		WHERE id = $8 RETURNING id, title, summary, body, cover_image, tags, author_id, status, published_at, created_at, updated_at`, post, id)
}

func (r *ContentRepository) post(ctx context.Context, query string, post model.Post, extra ...string) (model.Post, error) {
	args := []any{post.Title, post.Summary, post.Body, post.CoverImage, post.Tags, post.AuthorID, post.Status}
	for _, value := range extra {
		args = append(args, value)
	}
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return model.Post{}, fmt.Errorf("write post: %w", err)
	}
	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Post])
	if err != nil {
		return model.Post{}, fmt.Errorf("write post: %w", err)
	}
	return item, nil
}

func (r *ContentRepository) DeletePost(ctx context.Context, id string) error {
	return r.delete(ctx, "posts", id)
}

func (r *ContentRepository) CreateAIProduct(ctx context.Context, product model.AIProduct) (model.AIProduct, error) {
	return r.aiProduct(ctx, `INSERT INTO ai_products (name, summary, url, cover_image, tags, status, sort_order, published_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, CASE WHEN $6 = 'published' THEN now() END)
		RETURNING id, name, summary, url, cover_image, tags, status, sort_order, published_at, created_at, updated_at`, product)
}

func (r *ContentRepository) UpdateAIProduct(ctx context.Context, id string, product model.AIProduct) (model.AIProduct, error) {
	return r.aiProduct(ctx, `UPDATE ai_products SET name = $1, summary = $2, url = $3, cover_image = $4, tags = $5, status = $6, sort_order = $7,
		published_at = CASE WHEN $6 = 'published' THEN COALESCE(published_at, now()) ELSE NULL END, updated_at = now()
		WHERE id = $8 RETURNING id, name, summary, url, cover_image, tags, status, sort_order, published_at, created_at, updated_at`, product, id)
}

func (r *ContentRepository) aiProduct(ctx context.Context, query string, product model.AIProduct, extra ...string) (model.AIProduct, error) {
	args := []any{product.Name, product.Summary, product.URL, product.CoverImage, product.Tags, product.Status, product.SortOrder}
	for _, value := range extra {
		args = append(args, value)
	}
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return model.AIProduct{}, fmt.Errorf("write AI product: %w", err)
	}
	item, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.AIProduct])
	if err != nil {
		return model.AIProduct{}, fmt.Errorf("write AI product: %w", err)
	}
	return item, nil
}

func (r *ContentRepository) DeleteAIProduct(ctx context.Context, id string) error {
	return r.delete(ctx, "ai_products", id)
}

func (r *ContentRepository) delete(ctx context.Context, table, id string) error {
	command, err := r.pool.Exec(ctx, "DELETE FROM "+table+" WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("delete %s: %w", table, err)
	}
	if command.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}
