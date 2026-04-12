package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"wishlist/internal/application"
	"wishlist/internal/domain"
	"wishlist/internal/infrastructure/in/dto"
	"wishlist/internal/infrastructure/in/middleware"
	"wishlist/internal/infrastructure/in/response"
)

type ItemHandler struct {
	itemService     *application.ItemService
	wishlistService *application.WishlistService
	validate        *validator.Validate
}

func NewItemHandler(itemService *application.ItemService, wishlistService *application.WishlistService) *ItemHandler {
	return &ItemHandler{
		itemService:     itemService,
		wishlistService: wishlistService,
		validate:        validator.New(),
	}
}

func toItemResponse(i domain.Item) dto.ItemResponse {
	return dto.ItemResponse{
		ID:          i.ID,
		WishlistID:  i.WishlistID,
		Title:       i.Title,
		Description: i.Description,
		URL:         i.URL,
		Priority:    i.Priority,
		IsReserved:  i.IsReserved,
	}
}

func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	wishlistID, err := parseID(r, "wishlist_id")
	if err != nil {
		response.BadRequest(w, "invalid wishlist id")
		return
	}

	var req dto.ItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	created, err := h.itemService.Create(r.Context(), userID, domain.Item{
		WishlistID:  wishlistID,
		Title:       req.Title,
		Description: req.Description,
		URL:         req.URL,
		Priority:    req.Priority,
	})
	if err != nil {
		response.FromDomainError(w, err)
		return
	}

	response.Created(w, toItemResponse(created))
}

func (h *ItemHandler) GetAllByWishlist(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	wishlistID, err := parseID(r, "wishlist_id")
	if err != nil {
		response.BadRequest(w, "invalid wishlist id")
		return
	}

	items, err := h.itemService.GetAllByWishlist(r.Context(), userID, wishlistID)
	if err != nil {
		response.FromDomainError(w, err)
		return
	}

	result := make([]dto.ItemResponse, 0, len(items))
	for _, i := range items {
		result = append(result, toItemResponse(i))
	}

	response.OK(w, result)
}

func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	itemID, err := parseID(r, "id")
	if err != nil {
		response.BadRequest(w, "invalid item id")
		return
	}

	var req dto.ItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.BadRequest(w, err.Error())
		return
	}

	updated, err := h.itemService.Update(r.Context(), userID, domain.Item{
		ID:          itemID,
		Title:       req.Title,
		Description: req.Description,
		URL:         req.URL,
		Priority:    req.Priority,
	})
	if err != nil {
		response.FromDomainError(w, err)
		return
	}

	response.OK(w, toItemResponse(updated))
}

func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	itemID, err := parseID(r, "id")
	if err != nil {
		response.BadRequest(w, "invalid item id")
		return
	}

	if err := h.itemService.Delete(r.Context(), userID, itemID); err != nil {
		response.FromDomainError(w, err)
		return
	}

	response.NoContent(w)
}

func (h *ItemHandler) Reserve(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	itemID, err := parseID(r, "id")
	if err != nil {
		response.BadRequest(w, "invalid item id")
		return
	}

	if err := h.itemService.Reserve(r.Context(), token, itemID); err != nil {
		response.FromDomainError(w, err)
		return
	}

	response.NoContent(w)
}
