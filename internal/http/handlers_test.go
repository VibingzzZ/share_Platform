package http

import (
	"context"
	"encoding/json"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"

	"share-platform/internal/config"
	"share-platform/internal/model"
)

func TestPublicResourcesReturnsPublishedPaginatedEnvelope(t *testing.T) {
	store := contentStoreStub{
		resources: model.Page[model.Resource]{
			Items:    []model.Resource{{ID: "resource-1", Title: "Published guide", Status: model.StatusPublished}},
			Page:     2,
			PageSize: 5,
			Total:    6,
		},
	}
	req := httptest.NewRequest(stdhttp.MethodGet, "/api/v1/resources?page=2&pageSize=5&type=document&tag=guide", nil)
	res := httptest.NewRecorder()

	NewRouter(config.Config{}, Dependencies{Content: &store}).ServeHTTP(res, req)

	if res.Code != stdhttp.StatusOK {
		t.Fatalf("expected status 200, got %d", res.Code)
	}
	var body struct {
		Items    []model.Resource `json:"items"`
		Page     int              `json:"page"`
		PageSize int              `json:"pageSize"`
		Total    int              `json:"total"`
	}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Page != 2 || body.PageSize != 5 || body.Total != 6 || len(body.Items) != 1 {
		t.Fatalf("unexpected page response: %#v", body)
	}
	if body.Items[0].Status != model.StatusPublished {
		t.Fatalf("public read returned non-published item: %#v", body.Items[0])
	}
	if store.resourceFilter != (model.ListFilter{Page: 2, PageSize: 5, Type: "document", Tag: "guide"}) {
		t.Fatalf("unexpected resource filter: %#v", store.resourceFilter)
	}
}

func TestAdminResourceCreateRequiresAdminRole(t *testing.T) {
	req := httptest.NewRequest(stdhttp.MethodPost, "/api/v1/admin/resources", nil)
	req.Header.Set("X-User-ID", "00000000-0000-4000-8000-000000000001")
	req.Header.Set("X-User-Role", "member")
	res := httptest.NewRecorder()

	NewRouter(config.Config{}, Dependencies{}).ServeHTTP(res, req)

	if res.Code != stdhttp.StatusForbidden {
		t.Fatalf("expected status 403, got %d", res.Code)
	}
	var body errorResponse
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Code != "forbidden" {
		t.Fatalf("expected forbidden error code, got %q", body.Code)
	}
}

type contentStoreStub struct {
	resources      model.Page[model.Resource]
	resourceFilter model.ListFilter
}

func (s *contentStoreStub) ListResources(_ context.Context, filter model.ListFilter) (model.Page[model.Resource], error) {
	s.resourceFilter = filter
	return s.resources, nil
}

func (s *contentStoreStub) ListPosts(context.Context, model.ListFilter) (model.Page[model.Post], error) {
	return model.Page[model.Post]{}, nil
}

func (s *contentStoreStub) ListAIProducts(context.Context, model.ListFilter) (model.Page[model.AIProduct], error) {
	return model.Page[model.AIProduct]{}, nil
}

func (s *contentStoreStub) Dashboard(context.Context) (model.Dashboard, error) {
	return model.Dashboard{}, nil
}

func (s *contentStoreStub) CreateResource(context.Context, model.Resource) (model.Resource, error) {
	return model.Resource{}, nil
}

func (s *contentStoreStub) UpdateResource(context.Context, string, model.Resource) (model.Resource, error) {
	return model.Resource{}, nil
}

func (s *contentStoreStub) DeleteResource(context.Context, string) error { return nil }

func (s *contentStoreStub) CreatePost(context.Context, model.Post) (model.Post, error) {
	return model.Post{}, nil
}

func (s *contentStoreStub) UpdatePost(context.Context, string, model.Post) (model.Post, error) {
	return model.Post{}, nil
}

func (s *contentStoreStub) DeletePost(context.Context, string) error { return nil }

func (s *contentStoreStub) CreateAIProduct(context.Context, model.AIProduct) (model.AIProduct, error) {
	return model.AIProduct{}, nil
}

func (s *contentStoreStub) UpdateAIProduct(context.Context, string, model.AIProduct) (model.AIProduct, error) {
	return model.AIProduct{}, nil
}

func (s *contentStoreStub) DeleteAIProduct(context.Context, string) error { return nil }
