package job_application

import (
	"errors"
	"log"
	"strconv"

	"uuid"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/prionkor/careermesh/internal/authorization"
	"github.com/prionkor/careermesh/internal/httpapi"
	"github.com/prionkor/careermesh/internal/job_description"
	"github.com/prionkor/careermesh/internal/middleware"
	"github.com/prionkor/careermesh/internal/utility"
	"github.com/prionkor/careermesh/models"
)

const (
	defaultPageSize = 20
	maximumPageSize = 100

	// SQLSTATE foreign_key_violation and the job_description_id constraint name.
	foreignKeyViolation = "23503"
	jobDescriptionFK    = "job_applications_job_description_id_fkey"
)

type Handler struct {
	service               *Service
	jobDescriptionService *job_description.Service
}

func NewHandler(service *Service, jobDescriptionService *job_description.Service) *Handler {
	return &Handler{service: service, jobDescriptionService: jobDescriptionService}
}

type CreateJobApplicationRequest struct {
	JobDescriptionID uuid.UUID                    `json:"job_description_id" validate:"required"`
	Status           *models.JobApplicationStatus `json:"status" validate:"omitempty,model_valid"`
	AppliedAt        *models.Date                 `json:"applied_at"`
	Notes            string                       `json:"notes"`
}

type UpdateJobApplicationRequest struct {
	Status    models.JobApplicationStatus `json:"status" validate:"model_valid"`
	AppliedAt *models.Date                `json:"applied_at"`
	Notes     string                      `json:"notes"`
}

type ListJobApplicationsResponse struct {
	Items  []*models.JobApplication `json:"items"`
	Limit  int                      `json:"limit"`
	Offset int                      `json:"offset"`
}

func (h *Handler) Create(c fiber.Ctx) error {
	var req CreateJobApplicationRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if err := utility.ValidateStruct(req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}

	// Applications are always owned by the caller, who must also own the job description.
	userID := middleware.UserID(c)
	jobDescription, err := h.jobDescriptionService.GetByID(c.Context(), req.JobDescriptionID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return httpapi.WriteError(c, fiber.StatusNotFound, "job description not found")
		}
		return handleServiceError(c, "get job description for application", err)
	}
	if jobDescription.UserID != userID {
		return httpapi.WriteError(c, fiber.StatusNotFound, "job description not found")
	}

	status := models.JobApplicationStatusDraft
	if req.Status != nil {
		status = *req.Status
	}
	created, err := h.service.Create(c.Context(), &models.JobApplication{
		UserID: userID, JobDescriptionID: req.JobDescriptionID,
		Status: status, AppliedAt: req.AppliedAt, Notes: req.Notes,
	})
	if err != nil {
		if isMissingJobDescription(err) {
			return httpapi.WriteError(c, fiber.StatusNotFound, "job description not found")
		}
		return handleServiceError(c, "create job application", err)
	}

	return c.Status(fiber.StatusCreated).JSON(created)
}

func (h *Handler) GetByID(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid job application id")
	}

	application, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, "get job application by id", err)
	}
	if !authorization.CanAccessResource(
		middleware.UserID(c), application.UserID, middleware.Permissions(c),
		authorization.PermJobApplicationsReadOwn, authorization.PermJobApplicationsReadAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	return c.Status(fiber.StatusOK).JSON(application)
}

func (h *Handler) List(c fiber.Ctx) error {
	limit, offset, err := parsePagination(c)
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, err.Error())
	}

	var jobDescriptionID *uuid.UUID
	if raw := c.Query("job_description_id"); raw != "" {
		parsed, err := uuid.Parse(raw)
		if err != nil {
			return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid job description id")
		}
		jobDescriptionID = &parsed
	}

	var ownerID *uuid.UUID
	if !authorization.HasPermission(middleware.Permissions(c), authorization.PermJobApplicationsReadAll) {
		userID := middleware.UserID(c)
		ownerID = &userID
	}

	applications, err := h.service.List(c.Context(), ownerID, jobDescriptionID, limit, offset)
	if err != nil {
		return handleServiceError(c, "list job applications", err)
	}

	return c.Status(fiber.StatusOK).JSON(ListJobApplicationsResponse{
		Items: applications, Limit: limit, Offset: offset,
	})
}

func (h *Handler) Update(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid job application id")
	}

	var req UpdateJobApplicationRequest
	if err := c.Bind().Body(&req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}
	if err := utility.ValidateStruct(req); err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid request body")
	}

	saved, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, "get job application for update", err)
	}
	if !authorization.CanAccessResource(
		middleware.UserID(c), saved.UserID, middleware.Permissions(c),
		authorization.PermJobApplicationsUpdateOwn, authorization.PermJobApplicationsUpdateAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	updated, err := h.service.Update(c.Context(), &models.JobApplication{
		ID: id, Status: req.Status, AppliedAt: req.AppliedAt, Notes: req.Notes,
	})
	if err != nil {
		return handleServiceError(c, "update job application", err)
	}

	return c.Status(fiber.StatusOK).JSON(updated)
}

func (h *Handler) Delete(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return httpapi.WriteError(c, fiber.StatusBadRequest, "invalid job application id")
	}

	saved, err := h.service.GetByID(c.Context(), id)
	if err != nil {
		return handleServiceError(c, "get job application for delete", err)
	}
	if !authorization.CanAccessResource(
		middleware.UserID(c), saved.UserID, middleware.Permissions(c),
		authorization.PermJobApplicationsDeleteOwn, authorization.PermJobApplicationsDeleteAll,
	) {
		return httpapi.WriteError(c, fiber.StatusForbidden, "insufficient permissions")
	}

	if err := h.service.Delete(c.Context(), id); err != nil {
		return handleServiceError(c, "delete job application", err)
	}
	return c.SendStatus(fiber.StatusNoContent)
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
		return httpapi.WriteError(c, fiber.StatusNotFound, "job application not found")
	}

	log.Printf("%s: %v", action, err)
	return httpapi.WriteError(c, fiber.StatusInternalServerError, "internal server error")
}

// isMissingJobDescription reports a job_description_id foreign-key violation, e.g. the job description was deleted after the ownership check.
func isMissingJobDescription(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) &&
		pgErr.Code == foreignKeyViolation && pgErr.ConstraintName == jobDescriptionFK
}
