package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"gist/backend/internal/model"
	"gist/backend/internal/service"
	"gist/backend/pkg/logger"
	"gist/backend/pkg/network"
)

// SubscriptionHandler handles RSS subscription management via external API.
type SubscriptionHandler struct {
	feedService   service.FeedService
	folderService service.FolderService
}

// NewSubscriptionHandler creates a new SubscriptionHandler.
func NewSubscriptionHandler(feedService service.FeedService, folderService service.FolderService) *SubscriptionHandler {
	return &SubscriptionHandler{feedService: feedService, folderService: folderService}
}

type addSubscriptionRequest struct {
	URL      string `json:"url"`
	Title    string `json:"title"`
	Category string `json:"category"`
}

type subscriptionResponse struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	URL       string  `json:"url"`
	FolderID  *string `json:"folderId,omitempty"`
	IsNew     bool    `json:"isNew"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

type batchAddSubscriptionRequest struct {
	Items []addSubscriptionRequest `json:"items"`
}

type batchSubscriptionResult struct {
	URL    string                `json:"url"`
	Status string                `json:"status"` // "created", "exists", "error"
	Error  string                `json:"error,omitempty"`
	Feed   *subscriptionResponse `json:"feed,omitempty"`
}

type batchSubscriptionResponse struct {
	Created int                       `json:"created"`
	Skipped int                       `json:"skipped"`
	Errors  int                       `json:"errors"`
	Results []batchSubscriptionResult `json:"results"`
}

// RegisterRoutes registers subscription routes on the given group (requires auth).
func (h *SubscriptionHandler) RegisterRoutes(g *echo.Group) {
	g.POST("/subscriptions", h.Add)
	g.POST("/subscriptions/batch", h.AddBatch)
}

// Add subscribes to a feed idempotently.
// @Summary Subscribe to a feed
// @Description Subscribe to an RSS/Atom feed. Returns 201 if created, 200 if already subscribed.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param subscription body addSubscriptionRequest true "Subscription request"
// @Success 201 {object} subscriptionResponse "Feed subscription created"
// @Success 200 {object} subscriptionResponse "Already subscribed (idempotent)"
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Add(c echo.Context) error {
	var req addSubscriptionRequest
	if err := c.Bind(&req); err != nil {
		logger.Debug("subscription add invalid request", "module", "handler", "action", "create", "resource", "subscription", "result", "failed", "error", err)
		return c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid request"})
	}

	feedURL := strings.TrimSpace(req.URL)
	if feedURL == "" {
		return c.JSON(http.StatusBadRequest, errorResponse{Error: "url is required"})
	}

	folderID, httpErr := h.resolveCategoryToFolder(c, strings.TrimSpace(req.Category))
	if httpErr != nil {
		return httpErr
	}

	feed, err := h.feedService.Add(c.Request().Context(), feedURL, folderID, strings.TrimSpace(req.Title), "article")
	if err != nil {
		var conflictErr *service.FeedConflictError
		if errors.As(err, &conflictErr) {
			logger.Info("subscription already exists", "module", "handler", "action", "create", "resource", "subscription", "result", "ok", "host", network.ExtractHost(feedURL), "feed_id", conflictErr.ExistingFeed.ID)
			return c.JSON(http.StatusOK, toSubscriptionResponse(conflictErr.ExistingFeed, false))
		}
		logger.Error("subscription add failed", "module", "handler", "action", "create", "resource", "subscription", "result", "failed", "host", network.ExtractHost(feedURL), "error", err)
		return writeServiceError(c, err)
	}

	logger.Info("subscription created", "module", "handler", "action", "create", "resource", "subscription", "result", "ok", "feed_id", feed.ID, "feed_title", feed.Title, "host", network.ExtractHost(feed.URL))
	return c.JSON(http.StatusCreated, toSubscriptionResponse(feed, true))
}

// AddBatch subscribes to multiple feeds at once.
// @Summary Batch subscribe to feeds
// @Description Subscribe to multiple RSS/Atom feeds in a single request. Each item is processed independently.
// @Tags subscriptions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body batchAddSubscriptionRequest true "Batch subscription request"
// @Success 200 {object} batchSubscriptionResponse
// @Failure 400 {object} errorResponse
// @Failure 401 {object} errorResponse
// @Router /subscriptions/batch [post]
func (h *SubscriptionHandler) AddBatch(c echo.Context) error {
	var req batchAddSubscriptionRequest
	if err := c.Bind(&req); err != nil {
		logger.Debug("subscription batch add invalid request", "module", "handler", "action", "create", "resource", "subscription", "result", "failed", "error", err)
		return c.JSON(http.StatusBadRequest, errorResponse{Error: "invalid request"})
	}
	if len(req.Items) == 0 {
		return c.JSON(http.StatusBadRequest, errorResponse{Error: "items must not be empty"})
	}

	resp := batchSubscriptionResponse{
		Results: make([]batchSubscriptionResult, 0, len(req.Items)),
	}

	for _, item := range req.Items {
		feedURL := strings.TrimSpace(item.URL)
		if feedURL == "" {
			resp.Errors++
			resp.Results = append(resp.Results, batchSubscriptionResult{
				URL:    feedURL,
				Status: "error",
				Error:  "url is required",
			})
			continue
		}

		folderID, err := h.resolveCategory(c.Request().Context(), strings.TrimSpace(item.Category))
		if err != nil {
			logger.Warn("subscription batch category resolve failed", "module", "handler", "action", "create", "resource", "subscription", "result", "failed", "host", network.ExtractHost(feedURL), "error", err)
			resp.Errors++
			resp.Results = append(resp.Results, batchSubscriptionResult{
				URL:    feedURL,
				Status: "error",
				Error:  "failed to resolve category",
			})
			continue
		}

		feed, isNew, err := h.feedService.AddWithoutFetch(c.Request().Context(), feedURL, folderID, strings.TrimSpace(item.Title), "article")
		if err != nil {
			logger.Warn("subscription batch item failed", "module", "handler", "action", "create", "resource", "subscription", "result", "failed", "host", network.ExtractHost(feedURL), "error", err)
			resp.Errors++
			resp.Results = append(resp.Results, batchSubscriptionResult{
				URL:    feedURL,
				Status: "error",
				Error:  err.Error(),
			})
			continue
		}

		feedResp := toSubscriptionResponse(feed, isNew)
		if isNew {
			resp.Created++
			resp.Results = append(resp.Results, batchSubscriptionResult{URL: feedURL, Status: "created", Feed: &feedResp})
		} else {
			resp.Skipped++
			resp.Results = append(resp.Results, batchSubscriptionResult{URL: feedURL, Status: "exists", Feed: &feedResp})
		}
	}

	logger.Info("subscription batch add", "module", "handler", "action", "create", "resource", "subscription", "result", "ok", "created", resp.Created, "skipped", resp.Skipped, "errors", resp.Errors)
	return c.JSON(http.StatusOK, resp)
}

// resolveCategoryToFolder looks up or creates a folder for the given category name.
// Returns nil folderID when category is empty.
// Returns (nil, httpErr) when an HTTP error response has already been written.
func (h *SubscriptionHandler) resolveCategoryToFolder(c echo.Context, category string) (*int64, error) {
	if category == "" {
		return nil, nil
	}
	folderID, err := h.resolveCategory(c.Request().Context(), category)
	if err != nil {
		logger.Error("subscription resolve category failed", "module", "handler", "action", "create", "resource", "subscription", "result", "failed", "error", err)
		return nil, c.JSON(http.StatusInternalServerError, errorResponse{Error: "internal error"})
	}
	return folderID, nil
}

// resolveCategory returns the folder ID for the given category name,
// creating the folder if it does not already exist.
func (h *SubscriptionHandler) resolveCategory(ctx context.Context, category string) (*int64, error) {
	if category == "" {
		return nil, nil
	}

	folders, err := h.folderService.List(ctx)
	if err != nil {
		return nil, err
	}

	for _, f := range folders {
		if strings.EqualFold(f.Name, category) {
			id := f.ID
			return &id, nil
		}
	}

	// Folder not found — create it
	folder, err := h.folderService.Create(ctx, category, nil, "article")
	if err != nil {
		if errors.Is(err, service.ErrConflict) {
			// Race condition: another request created it; re-list and find
			folders, listErr := h.folderService.List(ctx)
			if listErr == nil {
				for _, f := range folders {
					if strings.EqualFold(f.Name, category) {
						id := f.ID
						return &id, nil
					}
				}
			}
		}
		return nil, err
	}

	return &folder.ID, nil
}

func toSubscriptionResponse(feed model.Feed, isNew bool) subscriptionResponse {
	return subscriptionResponse{
		ID:        idToString(feed.ID),
		Title:     feed.Title,
		URL:       feed.URL,
		FolderID:  idPtrToString(feed.FolderID),
		IsNew:     isNew,
		CreatedAt: feed.CreatedAt.UTC().Format(time.RFC3339),
		UpdatedAt: feed.UpdatedAt.UTC().Format(time.RFC3339),
	}
}
