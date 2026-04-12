package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"wishlist/internal/application"
	"wishlist/internal/domain"
	"wishlist/internal/infrastructure/in/dto"
	"wishlist/internal/infrastructure/in/middleware"
	"wishlist/internal/infrastructure/in/response"
)

type WishlistHandler struct {
	wishlistService *application.WishlistService
	itemService     *application.ItemService
	validate        *validator.Validate
}

func NewWishlistHandler(wishlistService *application.WishlistService, itemService *application.ItemService) *WishlistHandler {
	return &WishlistHandler{
		wishlistService: wishlistService,
		itemService:     itemService,
		validate:        validator.New(),
	}
}

func toWishlistResponse(w domain.Wishlist) dto.WishlistResponse {
	return dto.WishlistResponse{
		ID:          w.ID,
		Title:       w.Title,
		Description: w.Description,
		EventDate:   w.EventDate,
		PublicToken: w.PublicToken,
	}
}

func (h *WishlistHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	var req dto.WishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	eventDate, err := time.Parse(time.DateOnly, req.EventDate)
	if err != nil {
		response.BadRequest(w, "invalid event_date format, use YYYY-MM-DD")
		return
	}

	created, err := h.wishlistService.Create(r.Context(), domain.Wishlist{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		EventDate:   eventDate,
	})
	if err != nil {
		response.FromDomainError(w, err)
		return
	}

	response.Created(w, toWishlistResponse(created))
}

func (h *WishlistHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	wishlists, err := h.wishlistService.GetAllByUser(r.Context(), userID)
	if err != nil {
		response.FromDomainError(w, err)
		return
	}

	result := make([]dto.WishlistResponse, 0, len(wishlists))
	for _, wl := range wishlists {
		result = append(result, toWishlistResponse(wl))
	}

	response.OK(w, result)
}

func (h *WishlistHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	wishlistID, err := parseID(r, "id")
	if err != nil {
		response.BadRequest(w, "invalid wishlist id")
		return
	}

	wl, err := h.wishlistService.GetByID(r.Context(), userID, wishlistID)
	if err != nil {
		response.FromDomainError(w, err)
		return
	}

	response.OK(w, toWishlistResponse(wl))
}

func (h *WishlistHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	wishlistID, err := parseID(r, "id")
	if err != nil {
		response.BadRequest(w, "invalid wishlist id")
		return
	}

	var req dto.WishlistRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	eventDate, err := time.Parse(time.DateOnly, req.EventDate)
	if err != nil {
		response.BadRequest(w, "invalid event_date format, use YYYY-MM-DD")
		return
	}

	updated, err := h.wishlistService.Update(r.Context(), domain.Wishlist{
		UserID:      userID,
		ID:          wishlistID,
		Title:       req.Title,
		Description: req.Description,
		EventDate:   eventDate,
	})
	if err != nil {
		response.FromDomainError(w, err)
		return
	}

	response.OK(w, toWishlistResponse(updated))
}

func (h *WishlistHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	wishlistID, err := parseID(r, "id")
	if err != nil {
		response.BadRequest(w, "invalid wishlist id")
		return
	}

	if err := h.wishlistService.Delete(r.Context(), userID, wishlistID); err != nil {
		response.FromDomainError(w, err)
		return
	}

	response.NoContent(w)
}

func (h *WishlistHandler) GetPublic(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	wl, err := h.wishlistService.GetByPublicToken(r.Context(), token)
	if err != nil {
		response.FromDomainError(w, err)
		return
	}

	items, err := h.itemService.GetAllByWishlist(r.Context(), wl.UserID, wl.ID)
	if err != nil {
		response.FromDomainError(w, err)
		return
	}

	itemsResp := make([]dto.ItemResponse, 0, len(items))
	for _, item := range items {
		itemsResp = append(itemsResp, toItemResponse(item))
	}

	response.OK(w, map[string]any{
		"wishlist": toWishlistResponse(wl),
		"items":    itemsResp,
	})
}

func parseID(r *http.Request, param string) (int64, error) {
	return strconv.ParseInt(chi.URLParam(r, param), 10, 64)
}
