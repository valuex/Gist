package handler_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"gist/backend/internal/handler"
	"gist/backend/internal/model"
	"gist/backend/internal/service"
	"gist/backend/internal/service/mock"
)

func TestSubscriptionHandler_Add_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeedSvc := mock.NewMockFeedService(ctrl)
	mockFolderSvc := mock.NewMockFolderService(ctrl)
	h := handler.NewSubscriptionHandlerHelper(mockFeedSvc, mockFolderSvc)

	e := newTestEcho()
	reqBody := map[string]interface{}{
		"url":   "https://example.com/feed.xml",
		"title": "Example Feed",
	}
	req := newJSONRequest(http.MethodPost, "/subscriptions", reqBody)
	c, rec := newTestContext(e, req)

	expectedFeed := model.Feed{
		ID:    42,
		Title: "Example Feed",
		URL:   "https://example.com/feed.xml",
	}

	mockFeedSvc.EXPECT().
		Add(gomock.Any(), "https://example.com/feed.xml", (*int64)(nil), "Example Feed", "article").
		Return(expectedFeed, nil)

	err := h.Add(c)
	require.NoError(t, err)

	var resp handler.SubscriptionResponse
	assertJSONResponse(t, rec, http.StatusCreated, &resp)
	require.Equal(t, "42", resp.ID)
	require.Equal(t, "Example Feed", resp.Title)
	require.True(t, resp.IsNew)
}

func TestSubscriptionHandler_Add_AlreadyExists_Returns200(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeedSvc := mock.NewMockFeedService(ctrl)
	mockFolderSvc := mock.NewMockFolderService(ctrl)
	h := handler.NewSubscriptionHandlerHelper(mockFeedSvc, mockFolderSvc)

	e := newTestEcho()
	reqBody := map[string]interface{}{
		"url": "https://example.com/feed.xml",
	}
	req := newJSONRequest(http.MethodPost, "/subscriptions", reqBody)
	c, rec := newTestContext(e, req)

	existingFeed := model.Feed{
		ID:    99,
		Title: "Existing Feed",
		URL:   "https://example.com/feed.xml",
	}
	conflictErr := &service.FeedConflictError{ExistingFeed: existingFeed}

	mockFeedSvc.EXPECT().
		Add(gomock.Any(), "https://example.com/feed.xml", (*int64)(nil), "", "article").
		Return(model.Feed{}, conflictErr)

	err := h.Add(c)
	require.NoError(t, err)

	var resp handler.SubscriptionResponse
	assertJSONResponse(t, rec, http.StatusOK, &resp)
	require.Equal(t, "99", resp.ID)
	require.Equal(t, "Existing Feed", resp.Title)
	require.False(t, resp.IsNew)
}

func TestSubscriptionHandler_Add_InvalidURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeedSvc := mock.NewMockFeedService(ctrl)
	mockFolderSvc := mock.NewMockFolderService(ctrl)
	h := handler.NewSubscriptionHandlerHelper(mockFeedSvc, mockFolderSvc)

	e := newTestEcho()
	reqBody := map[string]interface{}{
		"url": "not-a-valid-url",
	}
	req := newJSONRequest(http.MethodPost, "/subscriptions", reqBody)
	c, rec := newTestContext(e, req)

	mockFeedSvc.EXPECT().
		Add(gomock.Any(), "not-a-valid-url", (*int64)(nil), "", "article").
		Return(model.Feed{}, service.ErrInvalid)

	err := h.Add(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSubscriptionHandler_Add_MissingURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeedSvc := mock.NewMockFeedService(ctrl)
	mockFolderSvc := mock.NewMockFolderService(ctrl)
	h := handler.NewSubscriptionHandlerHelper(mockFeedSvc, mockFolderSvc)

	e := newTestEcho()
	req := newJSONRequest(http.MethodPost, "/subscriptions", map[string]interface{}{})
	c, rec := newTestContext(e, req)

	err := h.Add(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSubscriptionHandler_Add_WithCategory_ExistingFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeedSvc := mock.NewMockFeedService(ctrl)
	mockFolderSvc := mock.NewMockFolderService(ctrl)
	h := handler.NewSubscriptionHandlerHelper(mockFeedSvc, mockFolderSvc)

	e := newTestEcho()
	reqBody := map[string]interface{}{
		"url":      "https://tech.example.com/feed",
		"category": "Tech",
	}
	req := newJSONRequest(http.MethodPost, "/subscriptions", reqBody)
	c, rec := newTestContext(e, req)

	folderID := int64(7)
	folders := []model.Folder{
		{ID: folderID, Name: "Tech", Type: "article"},
	}
	mockFolderSvc.EXPECT().List(gomock.Any()).Return(folders, nil)

	expectedFeed := model.Feed{
		ID:       55,
		Title:    "Tech News",
		URL:      "https://tech.example.com/feed",
		FolderID: &folderID,
	}
	mockFeedSvc.EXPECT().
		Add(gomock.Any(), "https://tech.example.com/feed", &folderID, "", "article").
		Return(expectedFeed, nil)

	err := h.Add(c)
	require.NoError(t, err)

	var resp handler.SubscriptionResponse
	assertJSONResponse(t, rec, http.StatusCreated, &resp)
	require.Equal(t, "55", resp.ID)
	require.NotNil(t, resp.FolderID)
	require.Equal(t, "7", *resp.FolderID)
	require.True(t, resp.IsNew)
}

func TestSubscriptionHandler_Add_WithCategory_NewFolder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeedSvc := mock.NewMockFeedService(ctrl)
	mockFolderSvc := mock.NewMockFolderService(ctrl)
	h := handler.NewSubscriptionHandlerHelper(mockFeedSvc, mockFolderSvc)

	e := newTestEcho()
	reqBody := map[string]interface{}{
		"url":      "https://science.example.com/feed",
		"category": "Science",
	}
	req := newJSONRequest(http.MethodPost, "/subscriptions", reqBody)
	c, rec := newTestContext(e, req)

	folderID := int64(20)
	mockFolderSvc.EXPECT().List(gomock.Any()).Return([]model.Folder{}, nil)
	mockFolderSvc.EXPECT().
		Create(gomock.Any(), "Science", (*int64)(nil), "article").
		Return(model.Folder{ID: folderID, Name: "Science", Type: "article"}, nil)

	expectedFeed := model.Feed{
		ID:       77,
		Title:    "Science Feed",
		URL:      "https://science.example.com/feed",
		FolderID: &folderID,
	}
	mockFeedSvc.EXPECT().
		Add(gomock.Any(), "https://science.example.com/feed", &folderID, "", "article").
		Return(expectedFeed, nil)

	err := h.Add(c)
	require.NoError(t, err)

	var resp handler.SubscriptionResponse
	assertJSONResponse(t, rec, http.StatusCreated, &resp)
	require.Equal(t, "77", resp.ID)
	require.True(t, resp.IsNew)
}

func TestSubscriptionHandler_AddBatch_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeedSvc := mock.NewMockFeedService(ctrl)
	mockFolderSvc := mock.NewMockFolderService(ctrl)
	h := handler.NewSubscriptionHandlerHelper(mockFeedSvc, mockFolderSvc)

	e := newTestEcho()
	reqBody := map[string]interface{}{
		"items": []map[string]interface{}{
			{"url": "https://a.example.com/rss"},
			{"url": "https://b.example.com/rss"},
		},
	}
	req := newJSONRequest(http.MethodPost, "/subscriptions/batch", reqBody)
	c, rec := newTestContext(e, req)

	now := time.Now()
	feedA := model.Feed{ID: 1, Title: "Feed A", URL: "https://a.example.com/rss", CreatedAt: now, UpdatedAt: now}
	feedB := model.Feed{ID: 2, Title: "Feed B", URL: "https://b.example.com/rss", CreatedAt: now, UpdatedAt: now}

	mockFeedSvc.EXPECT().
		AddWithoutFetch(gomock.Any(), "https://a.example.com/rss", (*int64)(nil), "", "article").
		Return(feedA, true, nil)
	mockFeedSvc.EXPECT().
		AddWithoutFetch(gomock.Any(), "https://b.example.com/rss", (*int64)(nil), "", "article").
		Return(feedB, false, nil) // already exists

	err := h.AddBatch(c)
	require.NoError(t, err)

	var resp handler.BatchSubscriptionResponse
	assertJSONResponse(t, rec, http.StatusOK, &resp)
	require.Equal(t, 1, resp.Created)
	require.Equal(t, 1, resp.Skipped)
	require.Equal(t, 0, resp.Errors)
	require.Len(t, resp.Results, 2)
	require.Equal(t, "created", resp.Results[0].Status)
	require.Equal(t, "exists", resp.Results[1].Status)
}

func TestSubscriptionHandler_AddBatch_EmptyItems(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeedSvc := mock.NewMockFeedService(ctrl)
	mockFolderSvc := mock.NewMockFolderService(ctrl)
	h := handler.NewSubscriptionHandlerHelper(mockFeedSvc, mockFolderSvc)

	e := newTestEcho()
	reqBody := map[string]interface{}{
		"items": []interface{}{},
	}
	req := newJSONRequest(http.MethodPost, "/subscriptions/batch", reqBody)
	c, rec := newTestContext(e, req)

	err := h.AddBatch(c)
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestSubscriptionHandler_AddBatch_InvalidURL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockFeedSvc := mock.NewMockFeedService(ctrl)
	mockFolderSvc := mock.NewMockFolderService(ctrl)
	h := handler.NewSubscriptionHandlerHelper(mockFeedSvc, mockFolderSvc)

	e := newTestEcho()
	reqBody := map[string]interface{}{
		"items": []map[string]interface{}{
			{"url": "bad-url"},
		},
	}
	req := newJSONRequest(http.MethodPost, "/subscriptions/batch", reqBody)
	c, rec := newTestContext(e, req)

	mockFeedSvc.EXPECT().
		AddWithoutFetch(gomock.Any(), "bad-url", (*int64)(nil), "", "article").
		Return(model.Feed{}, false, service.ErrInvalid)

	err := h.AddBatch(c)
	require.NoError(t, err)

	var resp handler.BatchSubscriptionResponse
	assertJSONResponse(t, rec, http.StatusOK, &resp)
	require.Equal(t, 0, resp.Created)
	require.Equal(t, 1, resp.Errors)
	require.Len(t, resp.Results, 1)
	require.Equal(t, "error", resp.Results[0].Status)
}
