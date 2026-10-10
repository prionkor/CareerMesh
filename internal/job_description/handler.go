package job_description

import (
	"errors"
	"log"
	"strconv"

	"uuid"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
	"github.com/prionkor/careermesh/internal/authorization"
	"github.com/prionkor/careermesh/internal/httpapi"
	"github.com/prionkor/careermesh/internal/middleware"
	"github.com/prionkor/careermesh/internal/utility"
	"github.com/prionkor/careermesh/models"
)

const (
	defaultPageSize = 20
	maximumPageSize = 100
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type JobDescriptionRequest struct {
	Title          string                 `json:"title" validate:"notblank"`
	Company        string                 `json:"company" validate:"notblank"`
	CompanyWebsite string                 `json:"company_website" validate:"omitempty,absolute_http_url"`
	SourceURL      string                 `json:"source_url" validate:"omitempty,absolute_http_url"`
	Location       string                 `json:"location"`
	EmploymentType *models.EmploymentType `json:"employment_type" validate:"omitempty,model_valid"`
	OriginalText   string                 `json:"original_text" validate:"notblank"`
}

type ListJobDescriptionsResponse struct {
	Items  []*models.JobDescription `json:"items"`
	Limit  int                      `json:"limit"`
	Offset int                      `json:"offset"`
}

func (h *Handler) Create(c fiber.Ctx) error {
	var req JobDescriptionRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if err := utility.ValidateStruct(req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}

	jobDescription := jobDescriptionFromRequest(req)
	jobDescription.UserID = middleware.UserID(c)
	created, err := h.service.Create(c.Context(), jobDescription)
	if err != nil {
		return handleServiceError(c, "create job description", err)
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

func (h *Handler) GetByID(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid job description id")
	}

	jobDescription, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, "get job description by id", err)
	}
	if !authorization.CanAccessResource(
		middleware.UserID(c), jobDescription.UserID, middleware.Permissions(c),
		authorization.PermJobDescriptionsReadOwn, authorization.PermJobDescriptionsReadAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	return c.Status(fiber.StatusOK).JSON(jobDescription)
}

func (h *Handler) List(c fiber.Ctx) error {
	limit, offset, err := parsePagination(c)
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, err.Error())
	}

	var ownerID *uuid.UUID
	if !authorization.HasPermission(middleware.Permissions(c), authorization.PermJobDescriptionsReadAll) {
		userID := middleware.UserID(c)
		ownerID = &userID
	}

	jobDescriptions, err := h.service.List(c.Context(), ownerID, limit, offset)
	if err != nil {
		return handleServiceError(c, "list job descriptions", err)
	}

	return c.Status(fiber.StatusOK).JSON(ListJobDescriptionsResponse{
		Items: jobDescriptions, Limit: limit, Offset: offset,
	})
}

func (h *Handler) Update(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid job description id")
	}

	var req JobDescriptionRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if err := utility.ValidateStruct(req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}

	saved, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, "get job description for update", err)
	}
	if !authorization.CanAccessResource(
		middleware.UserID(c), saved.UserID, middleware.Permissions(c),
		authorization.PermJobDescriptionsUpdateOwn, authorization.PermJobDescriptionsUpdateAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	jobDescription := jobDescriptionFromRequest(req)
	jobDescription.ID = id
	jobDescription.UserID = saved.UserID
	updated, err := h.service.Update(c.Context(), jobDescription)
	if err != nil {
		return handleServiceError(c, "update job description", err)
	}

	return c.Status(fiber.StatusOK).JSON(updated)
}

func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid job description id")
	}

	saved, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, "get job description for delete", err)
	}
	if !authorization.CanAccessResource(
		middleware.UserID(c), saved.UserID, middleware.Permissions(c),
		authorization.PermJobDescriptionsDeleteOwn, authorization.PermJobDescriptionsDeleteAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		return handleServiceError(c, "delete job description", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func jobDescriptionFromRequest(req JobDescriptionRequest) *models.JobDescription {
	return &models.JobDescription{
		Title: req.Title, Company: req.Company, CompanyWebsite: req.CompanyWebsite,
		SourceURL: req.SourceURL, Location: req.Location, EmploymentType: req.EmploymentType,
		OriginalText: req.OriginalText,
	}
}

func parsePagination(c fiber.Ctx) (int, int, error) {
	limit := defaultPageSize
	if rawLimit := c.Query("limit"); rawLimit != "" {
		parsedLimit, err := strconv.Atoi(rawLimit)
		if err != nil || parsedLimit < 1 {
			return 0, 0, errors.New("limit must be a positive integer")
		}
		limit = parsedLimit
		if limit > maximumPageSize {
			limit = maximumPageSize
		}
	}

	offset := 0
	if rawOffset := c.Query("offset"); rawOffset != "" {
		parsedOffset, err := strconv.Atoi(rawOffset)
		if err != nil || parsedOffset < 0 {
			return 0, 0, errors.New("offset must be a non-negative integer")
		}
		offset = parsedOffset
	}
	return limit, offset, nil
}

func handleServiceError(c fiber.Ctx, action string, err error) error {
	if errors.Is(err, authorization.ErrForbidden) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}
	if errors.Is(err, pgx.ErrNoRows) {
		return httpapi.WriteError(c, fiber.StatusNotFound, "job description not found")
	}

	log.Printf("%s: %v", action, err)
	return httpapi.WriteError(c, fiber.StatusInternalServerError, "internal server error")
}
