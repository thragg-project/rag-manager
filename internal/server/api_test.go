package server

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/thragg-project/rag-manager/internal/core"
	apperrors "github.com/thragg-project/rag-manager/internal/errors"
	"github.com/thragg-project/rag-manager/internal/server/handlers"
)

type fakeDataSourceService struct{}

var _ handlers.DataSourceService = (*fakeDataSourceService)(nil)

func (s *fakeDataSourceService) Get(ctx context.Context, id int64) (core.DataSource, error) {
	return core.DataSource{}, apperrors.ErrNotFound
}

func (s *fakeDataSourceService) Create(
	ctx context.Context,
	dataSource core.DataSource,
) (core.DataSource, error) {
	id := int64(1)
	dataSource.Id = &id
	return dataSource, nil
}

func (s *fakeDataSourceService) Search(ctx context.Context) ([]core.DataSource, error) {
	return nil, nil
}

func TestAPIServerValidatesRequests(t *testing.T) {
	apiServer := NewAPIServer("127.0.0.1", 0, testLogger())
	err := apiServer.RegisterHandlers(
		handlers.NewDataSourceHandler(&fakeDataSourceService{}),
		nil,
		nil,
	)
	if err != nil {
		t.Fatalf("register handlers: %v", err)
	}

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		"/api/v1/datasources",
		strings.NewReader(
			`{"provider":"s3","credentials":{},"dataType":"html","extra":"rejected"}`,
		),
	)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	apiServer.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d: %s",
			http.StatusBadRequest,
			rec.Code,
			rec.Body.String(),
		)
	}
}

func TestAPIServerMapsNotFoundError(t *testing.T) {
	apiServer := NewAPIServer("127.0.0.1", 0, testLogger())
	err := apiServer.RegisterHandlers(
		handlers.NewDataSourceHandler(&fakeDataSourceService{}),
		nil,
		nil,
	)
	if err != nil {
		t.Fatalf("register handlers: %v", err)
	}

	req := httptest.NewRequestWithContext(
		context.Background(),
		http.MethodGet,
		"/api/v1/datasources/1",
		nil,
	)
	rec := httptest.NewRecorder()

	apiServer.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d: %s", http.StatusNotFound, rec.Code, rec.Body.String())
	}
}

func testLogger() *slog.Logger {
	return slog.New(slog.DiscardHandler)
}
