package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"UAS/app/models"
	"UAS/app/repository"
	"UAS/utils"
)

// AchievementService defines all achievement-related operations
type AchievementService interface {
	CreateAchievement(c *fiber.Ctx) error
	UpdateAchievement(c *fiber.Ctx) error
	DeleteAchievement(c *fiber.Ctx) error
	SubmitAchievement(c *fiber.Ctx) error
	ListAchievements(c *fiber.Ctx) error
	GetAchievementDetail(c *fiber.Ctx) error
	GetAchievementHistory(c *fiber.Ctx) error
	GetStatistics(c *fiber.Ctx) error
	GetStudentReport(c *fiber.Ctx) error
	VerifyAchievement(c *fiber.Ctx) error
	RejectAchievement(c *fiber.Ctx) error
	UploadAttachment(c *fiber.Ctx) error
}

type achievementServiceImpl struct {
	pgRepo       *repository.AchievementRepository
	mongoRepo    *repository.MongoAchievementRepository
	studentRepo  *repository.StudentRepository
	userRepo     *repository.UserRepository
	lecturerRepo *repository.LecturerRepository
}

func NewAchievementService() AchievementService {
	return &achievementServiceImpl{
		pgRepo:       repository.NewAchievementRepository(),
		mongoRepo:    repository.NewMongoAchievementRepository(),
		studentRepo:  repository.NewStudentRepository(),
		userRepo:     repository.NewUserRepository(),
		lecturerRepo: repository.NewLecturerRepository(),
	}
}

// CreateAchievement handles achievement creation
// @Summary Create new achievement
// @Description Create a new achievement for the logged-in student
// @Tags Achievements
// @Accept json
// @Produce json
// @Param body body models.CreateAchievementRequest true "Achievement data"
// @Success 201 {object} models.AchievementReference
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /achievements [post]
// @Security Bearer
func (s *achievementServiceImpl) CreateAchievement(c *fiber.Ctx) error {
	var req models.CreateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	if req.Title == "" || req.AchievementType == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "title and achievement_type are required")
	}

	if c.Locals("role") != "Mahasiswa" {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "only mahasiswa can create achievements")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoAch, err := s.mongoRepo.Create(ctx, &models.MongoAchievement{
		StudentID:       c.Locals("userID").(string),
		Title:           req.Title,
		Description:     req.Description,
		AchievementType: req.AchievementType,
		Details:         req.Details,
		Tags:            req.Tags,
		Points:          req.Points,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	})
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to save achievement")
	}

	pgAch := &models.AchievementReference{
		ID:                 uuid.New().String(),
		StudentID:          c.Locals("userID").(string),
		MongoAchievementID: mongoAch.ID.Hex(),
		Status:             "draft",
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	if err := s.pgRepo.Create(pgAch); err != nil {
		s.mongoRepo.SoftDelete(ctx, mongoAch.ID.Hex())
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to save achievement reference")
	}

	return utils.CreatedResponse(c, "Prestasi berhasil dibuat", pgAch)
}

// ListAchievements handles listing achievements
// @Summary List achievements
// @Description Get list of achievements based on user role
// @Tags Achievements
// @Produce json
// @Success 200 {array} models.AchievementReference
// @Failure 500 {object} map[string]interface{}
// @Router /achievements [get]
// @Security Bearer
func (s *achievementServiceImpl) ListAchievements(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	role := c.Locals("role").(string)

	// Get query parameters for filtering, sorting, and pagination
	status := c.Query("status", "")            // draft, submitted, verified, rejected
	achievementType := c.Query("type", "")     // competition, publication, organization, certification
	sortBy := c.Query("sort_by", "created_at") // created_at, updated_at, title
	sortOrder := c.Query("sort_order", "desc") // asc, desc
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("page_size", 10)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var achievements []models.AchievementReference
	var err error

	switch role {
	case "Admin":
		// Admin can see all achievements
		achievements, err = s.pgRepo.FindAll()

	case "Mahasiswa":
		// Student can only see their own achievements
		achievements, err = s.pgRepo.FindByStudentID(userID)

	case "Dosen", "Dosen Wali":
		// Lecturer can only see achievements from their advisees (anak wali)
		// First, get all students where AdvisorID = current user
		students, err := s.studentRepo.FindByAdvisorID(userID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to retrieve advisees")
		}

		// Then get achievements for all advisees
		achievements = []models.AchievementReference{}
		for _, student := range students {
			studentAchievements, err := s.pgRepo.FindByStudentID(student.UserID)
			if err != nil {
				continue // Skip if error
			}
			achievements = append(achievements, studentAchievements...)
		}

	default:
		return utils.ErrorResponse(c, fiber.StatusForbidden, "insufficient permissions to view achievements")
	}

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to retrieve achievements: "+err.Error())
	}

	// Debug: log jumlah achievements
	// fmt.Printf("DEBUG - Role: %s, Total achievements found: %d\n", role, len(achievements))

	// If no achievements found, return empty response with message
	if len(achievements) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  true,
			"message": "no achievements found for user (role: " + role + ")",
			"data":    []fiber.Map{},
			"pagination": fiber.Map{
				"page":        1,
				"page_size":   pageSize,
				"total":       0,
				"total_pages": 0,
			},
		})
	}

	// Apply filters
	var filteredAchievements []models.AchievementReference
	for _, ach := range achievements {
		// Filter by status
		if status != "" && ach.Status != status {
			continue
		}
		// Filter by achievement type (need to fetch from MongoDB)
		if achievementType != "" {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			mongoAch, err := s.mongoRepo.FindByID(ctx, ach.MongoAchievementID)
			cancel()
			if err != nil || mongoAch.AchievementType != achievementType {
				continue
			}
		}
		filteredAchievements = append(filteredAchievements, ach)
	}

	// Apply sorting
	for i := 0; i < len(filteredAchievements)-1; i++ {
		for j := 0; j < len(filteredAchievements)-1-i; j++ {
			shouldSwap := false
			if sortOrder == "asc" {
				if sortBy == "title" {
					shouldSwap = filteredAchievements[j].Status > filteredAchievements[j+1].Status
				} else if sortBy == "updated_at" {
					shouldSwap = filteredAchievements[j].UpdatedAt.After(filteredAchievements[j+1].UpdatedAt)
				} else { // default: created_at
					shouldSwap = filteredAchievements[j].CreatedAt.After(filteredAchievements[j+1].CreatedAt)
				}
			} else { // desc
				if sortBy == "title" {
					shouldSwap = filteredAchievements[j].Status < filteredAchievements[j+1].Status
				} else if sortBy == "updated_at" {
					shouldSwap = filteredAchievements[j].UpdatedAt.Before(filteredAchievements[j+1].UpdatedAt)
				} else { // default: created_at
					shouldSwap = filteredAchievements[j].CreatedAt.Before(filteredAchievements[j+1].CreatedAt)
				}
			}
			if shouldSwap {
				filteredAchievements[j], filteredAchievements[j+1] = filteredAchievements[j+1], filteredAchievements[j]
			}
		}
	}

	// Apply pagination
	totalRecords := len(filteredAchievements)
	totalPages := (totalRecords + pageSize - 1) / pageSize
	startIdx := (page - 1) * pageSize
	endIdx := startIdx + pageSize

	if startIdx >= totalRecords {
		startIdx = totalRecords
	}
	if endIdx > totalRecords {
		endIdx = totalRecords
	}

	paginatedAchievements := filteredAchievements[startIdx:endIdx]

	// Fetch MongoDB details for each achievement
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	response := fiber.Map{
		"data": []fiber.Map{},
		"pagination": fiber.Map{
			"page":        page,
			"page_size":   pageSize,
			"total":       totalRecords,
			"total_pages": totalPages,
		},
	}

	responseData := make([]fiber.Map, len(paginatedAchievements))
	for i, ach := range paginatedAchievements {
		mongoAch, err := s.mongoRepo.FindByID(ctx, ach.MongoAchievementID)
		if err != nil {
			mongoAch = nil // If not found, just continue
		}

		responseData[i] = fiber.Map{
			"id":              ach.ID,
			"student_id":      ach.StudentID,
			"mongo_id":        ach.MongoAchievementID,
			"status":          ach.Status,
			"created_at":      ach.CreatedAt,
			"updated_at":      ach.UpdatedAt,
			"verified_at":     ach.VerifiedAt,
			"mongodb_details": mongoAch,
		}
	}
	response["data"] = responseData

	return c.Status(fiber.StatusOK).JSON(response)
}

// GetAchievementDetail handles getting achievement detail
// @Summary Get achievement detail
// @Description Retrieve detailed information of a specific achievement
// @Tags Achievements
// @Produce json
// @Param id path string true "Achievement ID"
// @Success 200 {object} models.AchievementReference
// @Failure 404 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Router /achievements/{id} [get]
// @Security Bearer
func (s *achievementServiceImpl) GetAchievementDetail(c *fiber.Ctx) error {
	achievement, err := s.pgRepo.FindByID(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "achievement not found")
	}

	role := c.Locals("role").(string)
	userID := c.Locals("userID").(string)

	// Mahasiswa can only view their own achievements
	if role == "Mahasiswa" && achievement.StudentID != userID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "you can only view your own achievements")
	}

	// Dosen Wali can only view achievements of their advisees
	if role == "Dosen Wali" {
		// Get student info
		student, err := s.studentRepo.FindByUserID(achievement.StudentID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to get student info")
		}

		// Get lecturer info
		lecturer, err := s.lecturerRepo.FindByUserID(userID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "lecturer profile not found")
		}

		// Check if this lecturer is the advisor
		if student.AdvisorID != lecturer.ID {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "you can only view achievements of your advisees")
		}
	}

	return utils.SuccessResponse(c, "achievement detail retrieved", achievement)
}

// UpdateAchievement handles updating achievement
// @Summary Update achievement
// @Description Update an existing achievement (only draft status)
// @Tags Achievements
// @Accept json
// @Produce json
// @Param id path string true "Achievement ID"
// @Param body body models.UpdateAchievementRequest true "Updated achievement data"
// @Success 200 {object} models.AchievementReference
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /achievements/{id} [put]
// @Security Bearer
func (s *achievementServiceImpl) UpdateAchievement(c *fiber.Ctx) error {
	var req models.UpdateAchievementRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	achievement, err := s.pgRepo.FindByID(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "achievement not found")
	}

	if c.Locals("role") == "Mahasiswa" && achievement.StudentID != c.Locals("userID").(string) {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "you can only update your own achievements")
	}

	if achievement.Status != "draft" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "only draft achievements can be updated")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if _, err := s.mongoRepo.Update(ctx, achievement.MongoAchievementID, &models.MongoAchievement{
		StudentID:       c.Locals("userID").(string),
		Title:           req.Title,
		Description:     req.Description,
		AchievementType: req.AchievementType,
		Details:         req.Details,
		Tags:            req.Tags,
		Points:          req.Points,
	}); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to update achievement")
	}

	achievement.UpdatedAt = time.Now()
	if err := s.pgRepo.Update(c.Params("id"), achievement); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to update achievement")
	}

	return utils.SuccessResponse(c, "Prestasi berhasil diperbarui", achievement)
}

// DeleteAchievement handles deleting achievement
// @Summary Delete achievement
// @Description Delete an achievement (only draft status)
// @Tags Achievements
// @Produce json
// @Param id path string true "Achievement ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /achievements/{id} [delete]
// @Security Bearer
func (s *achievementServiceImpl) DeleteAchievement(c *fiber.Ctx) error {
	achievement, err := s.pgRepo.FindByID(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "achievement not found")
	}

	if c.Locals("role") == "Mahasiswa" && achievement.StudentID != c.Locals("userID").(string) {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "you can only delete your own achievements")
	}

	if achievement.Status != "draft" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "only draft achievements can be deleted")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.mongoRepo.SoftDelete(ctx, achievement.MongoAchievementID); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to delete achievement")
	}

	if err := s.pgRepo.Delete(c.Params("id")); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to delete achievement")
	}

	return utils.DeletedResponse(c, "Prestasi berhasil dihapus")
}

// SubmitAchievement handles submitting achievement for verification
// @Summary Submit achievement
// @Description Submit an achievement for verification (changes status from draft to submitted)
// @Tags Achievements
// @Produce json
// @Param id path string true "Achievement ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /achievements/{id}/submit [post]
// @Security Bearer
func (s *achievementServiceImpl) SubmitAchievement(c *fiber.Ctx) error {
	achievement, err := s.pgRepo.FindByID(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "achievement not found")
	}

	if c.Locals("role") == "Mahasiswa" && achievement.StudentID != c.Locals("userID").(string) {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "you can only submit your own achievements")
	}

	if achievement.Status != "draft" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "only draft achievements can be submitted")
	}

	if err := s.pgRepo.UpdateStatus(c.Params("id"), "submitted"); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to submit achievement")
	}

	return utils.SuccessResponse(c, "Prestasi berhasil disubmit untuk verifikasi", fiber.Map{"id": c.Params("id"), "status": "submitted"})
}

// GetAchievementHistory handles getting achievement history
// @Summary Get achievement history
// @Description Get the timeline/history of an achievement's status changes
// @Tags Achievements
// @Produce json
// @Param id path string true "Achievement ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /achievements/{id}/history [get]
// @Security Bearer
func (s *achievementServiceImpl) GetAchievementHistory(c *fiber.Ctx) error {
	achievement, err := s.pgRepo.FindByID(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "achievement not found")
	}

	if c.Locals("role") == "Mahasiswa" && achievement.StudentID != c.Locals("userID").(string) {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "you can only view your own achievement history")
	}

	timeline := []map[string]interface{}{
		{
			"status":    "draft",
			"timestamp": achievement.CreatedAt,
			"message":   "Achievement created",
		},
	}

	if !achievement.SubmittedAt.IsZero() {
		timeline = append(timeline, map[string]interface{}{
			"status":    "submitted",
			"timestamp": achievement.SubmittedAt,
			"message":   "Achievement submitted for verification",
		})
	}

	if !achievement.VerifiedAt.IsZero() {
		timeline = append(timeline, map[string]interface{}{
			"status":      "verified",
			"timestamp":   achievement.VerifiedAt,
			"verified_by": achievement.VerifiedBy,
			"message":     "Achievement verified",
		})
	}

	if achievement.Status == "rejected" && achievement.RejectionNote != "" {
		timeline = append(timeline, map[string]interface{}{
			"status":  "rejected",
			"message": achievement.RejectionNote,
		})
	}

	return utils.SuccessResponse(c, "achievement history retrieved", map[string]interface{}{
		"id":       c.Params("id"),
		"status":   achievement.Status,
		"timeline": timeline,
	})
}

// GetStatistics handles getting achievement statistics
// @Summary Get achievement statistics
// @Description Get comprehensive statistics of achievements based on user role
// @Tags Reports
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /reports/statistics [get]
// @Security Bearer
func (s *achievementServiceImpl) GetStatistics(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	role := c.Locals("role").(string)

	var achievements []models.AchievementReference
	var err error

	switch role {
	case "Mahasiswa":
		// Student sees only their own
		achievements, err = s.pgRepo.FindByStudentID(userID)
	case "Dosen", "Dosen Wali":
		// Lecturer sees their advisees
		students, err := s.studentRepo.FindByAdvisorID(userID)
		if err != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to retrieve advisees")
		}
		achievements = []models.AchievementReference{}
		for _, student := range students {
			studentAchievements, err := s.pgRepo.FindByStudentID(student.UserID)
			if err != nil {
				continue
			}
			achievements = append(achievements, studentAchievements...)
		}
	default:
		// Admin sees all
		achievements, err = s.pgRepo.FindAll()
	}

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to retrieve achievements")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	stats := s.buildStatistics(ctx, achievements)
	return utils.SuccessResponse(c, "statistics retrieved successfully", stats)
}

// buildStatistics builds comprehensive statistics from achievements
func (s *achievementServiceImpl) buildStatistics(ctx context.Context, achievements []models.AchievementReference) fiber.Map {
	// 1. Status distribution
	statusCount := map[string]int64{
		"draft":     0,
		"submitted": 0,
		"verified":  0,
		"rejected":  0,
	}
	for _, ach := range achievements {
		statusCount[ach.Status]++
	}

	// 2. Achievement type distribution
	typeCount := map[string]int64{
		"competition":   0,
		"publication":   0,
		"organization":  0,
		"certification": 0,
	}

	// 3. Competition level distribution
	levelCount := map[string]int64{
		"school":        0,
		"city":          0,
		"provincial":    0,
		"national":      0,
		"international": 0,
	}

	// 4. Period distribution (by year)
	periodCount := make(map[string]int64)

	// 5. Student with most achievements
	studentAchievementCount := make(map[string]int64)

	for _, ach := range achievements {
		// Get MongoDB details for type and level
		mongoAch, err := s.mongoRepo.FindByID(ctx, ach.MongoAchievementID)
		if err == nil && mongoAch != nil {
			// Count by type
			typeCount[mongoAch.AchievementType]++

			// Count by period (year)
			year := mongoAch.CreatedAt.Year()
			yearStr := fmt.Sprintf("%d", year)
			periodCount[yearStr]++

			// Count competition levels if it's competition type
			if mongoAch.AchievementType == "competition" {
				if mongoAch.Details != nil {
					if level, exists := mongoAch.Details["competition_level"]; exists {
						if levelStr, ok := level.(string); ok {
							levelCount[levelStr]++
						}
					}
				}
			}
		}

		// Count by student
		studentAchievementCount[ach.StudentID]++
	}

	// Find top 5 students
	type studentStat struct {
		StudentID        string
		AchievementCount int64
	}
	var topStudents []studentStat
	for studentID, count := range studentAchievementCount {
		topStudents = append(topStudents, studentStat{studentID, count})
	}
	// Simple sorting (bubble sort)
	for i := 0; i < len(topStudents)-1; i++ {
		for j := 0; j < len(topStudents)-1-i; j++ {
			if topStudents[j].AchievementCount < topStudents[j+1].AchievementCount {
				topStudents[j], topStudents[j+1] = topStudents[j+1], topStudents[j]
			}
		}
	}
	if len(topStudents) > 5 {
		topStudents = topStudents[:5]
	}

	total := int64(len(achievements))
	verificationRate := 0.0
	if total > 0 {
		verificationRate = float64(statusCount["verified"]) / float64(total) * 100
	}

	return fiber.Map{
		"summary": fiber.Map{
			"total":             total,
			"verification_rate": verificationRate,
		},
		"by_status":            statusCount,
		"by_type":              typeCount,
		"by_competition_level": levelCount,
		"by_period":            periodCount,
		"top_students":         topStudents,
	}
}

// VerifyAchievement handles achievement verification by lecturer
// @Summary Verify achievement
// @Description Verify an achievement submission (Dosen Wali only)
// @Tags Achievements
// @Accept json
// @Produce json
// @Param id path string true "Achievement ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /achievements/{id}/verify [post]
// @Security Bearer
func (s *achievementServiceImpl) VerifyAchievement(c *fiber.Ctx) error {
	achievementID := c.Params("id")
	dosenID := c.Locals("userID").(string)

	// Get achievement first to update it properly
	achievement, err := s.pgRepo.FindByID(achievementID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "achievement not found")
	}

	// Check if status is submitted
	if achievement.Status != "submitted" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "only submitted achievements can be verified")
	}

	// Verify that the dosen is the advisor of the student
	student, err := s.studentRepo.FindByUserID(achievement.StudentID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "student not found")
	}

	if student.AdvisorID != dosenID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "only the student's advisor can verify achievements")
	}

	// Update achievement status, verified_at, and verified_by
	achievement.Status = "verified"
	achievement.VerifiedAt = time.Now()
	achievement.VerifiedBy = dosenID
	achievement.UpdatedAt = time.Now()

	if err := s.pgRepo.Update(achievementID, achievement); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to verify achievement")
	}

	return utils.SuccessResponse(c, "achievement verified successfully", nil)
}

// RejectAchievement handles achievement rejection by lecturer
// @Summary Reject achievement
// @Description Reject an achievement submission with notes (Dosen Wali only)
// @Tags Achievements
// @Accept json
// @Produce json
// @Param id path string true "Achievement ID"
// @Param body body map[string]interface{} true "Rejection data"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /achievements/{id}/reject [post]
// @Security Bearer
func (s *achievementServiceImpl) RejectAchievement(c *fiber.Ctx) error {
	achievementID := c.Params("id")
	dosenID := c.Locals("userID").(string)

	var req struct {
		RejectionNote string `json:"rejection_note" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "invalid request body")
	}

	// Get achievement first to update it properly
	achievement, err := s.pgRepo.FindByID(achievementID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "achievement not found")
	}

	// Check if status is submitted
	if achievement.Status != "submitted" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "only submitted achievements can be rejected")
	}

	// Verify that the dosen is the advisor of the student
	student, err := s.studentRepo.FindByUserID(achievement.StudentID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "student not found")
	}

	if student.AdvisorID != dosenID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "only the student's advisor can reject achievements")
	}

	// Update achievement status, rejection note, verified_at, and verified_by
	achievement.Status = "rejected"
	achievement.RejectionNote = req.RejectionNote
	achievement.VerifiedAt = time.Now()
	achievement.VerifiedBy = dosenID
	achievement.UpdatedAt = time.Now()

	if err := s.pgRepo.Update(achievementID, achievement); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to reject achievement")
	}

	return utils.SuccessResponse(c, "achievement rejected successfully", nil)
}

// UploadAttachment handles file attachment upload for achievements
// @Summary Upload achievement attachment
// @Description Upload proof files for an achievement
// @Tags Achievements
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Achievement ID"
// @Param file formData file true "File to upload"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /achievements/{id}/attachments [post]
// @Security Bearer
func (s *achievementServiceImpl) UploadAttachment(c *fiber.Ctx) error {
	achievementID := c.Params("id")

	// Get file from request
	file, err := c.FormFile("file")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "file is required")
	}

	// Validate file size (max 10MB)
	if file.Size > 10*1024*1024 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "file size exceeds 10MB limit")
	}

	// Verify achievement exists
	achievement, err := s.pgRepo.FindByID(achievementID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "achievement not found")
	}

	// Check if achievement is still in draft status
	// Can only upload attachments when status is "draft"
	if achievement.Status != "draft" {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "attachments can only be uploaded for draft achievements. Current status: "+achievement.Status)
	}

	// Verify ownership
	if c.Locals("role") == "Mahasiswa" && achievement.StudentID != c.Locals("userID").(string) {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "you can only upload attachments to your own achievements")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get MongoDB achievement to add attachment
	mongoAch, err := s.mongoRepo.FindByID(ctx, achievement.MongoAchievementID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "achievement not found in database")
	}

	// Save file to local storage or cloud storage
	// For now, we'll store the file in a local directory
	uploadDir := "uploads/achievements"

	// Create directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to create upload directory")
	}

	filename := uuid.New().String() + "_" + file.Filename
	filepath := uploadDir + "/" + filename

	// Save the file
	if err := c.SaveFile(file, filepath); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to save file: "+err.Error())
	}

	// Create attachment object
	attachment := models.Attachment{
		FileName:   file.Filename,
		FileURL:    "/uploads/achievements/" + filename, // URL path for serving file
		FileType:   file.Header.Get("Content-Type"),
		UploadedAt: time.Now(),
	}

	// Add attachment to MongoDB achievement
	if mongoAch.Attachments == nil {
		mongoAch.Attachments = []models.Attachment{}
	}
	mongoAch.Attachments = append(mongoAch.Attachments, attachment)
	mongoAch.UpdatedAt = time.Now()

	// Update MongoDB document
	if _, err := s.mongoRepo.Update(ctx, achievement.MongoAchievementID, mongoAch); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to save attachment record")
	}

	// Update PostgreSQL timestamp
	achievement.UpdatedAt = time.Now()
	if err := s.pgRepo.Update(achievementID, achievement); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to update achievement")
	}

	return utils.SuccessResponse(c, "file uploaded successfully", fiber.Map{
		"achievement_id": achievementID,
		"file_name":      file.Filename,
		"file_url":       "/uploads/achievements/" + filename,
		"file_type":      file.Header.Get("Content-Type"),
		"file_size":      file.Size,
		"uploaded_at":    time.Now(),
	})
}

// GetStudentReport handles getting detailed report of specific student's achievements
// @Summary Get student achievement report
// @Description Get comprehensive achievement report for a specific student
// @Tags Reports
// @Produce json
// @Param id path string true "Student User ID (UUID)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /reports/student/{id} [get]
// @Security Bearer
func (s *achievementServiceImpl) GetStudentReport(c *fiber.Ctx) error {
	userID := c.Locals("userID").(string)
	role := c.Locals("role").(string)
	studentUserID := c.Params("id") // This is the User ID of the student

	if studentUserID == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "student id is required")
	}

	// Authorization check
	if role == "Mahasiswa" && userID != studentUserID {
		return utils.ErrorResponse(c, fiber.StatusForbidden, "you can only view your own report")
	}

	if role == "Dosen" || role == "Dosen Wali" {
		// Check if student is their advisee
		student, err := s.studentRepo.FindByUserID(studentUserID)
		if err != nil || student == nil {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "student not found")
		}
		if student.AdvisorID != userID {
			return utils.ErrorResponse(c, fiber.StatusForbidden, "you can only view your advisees' reports")
		}
	}

	// Get student info
	student, err := s.studentRepo.FindByUserID(studentUserID)
	if err != nil || student == nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "student not found")
	}

	// Get user info for name
	user, err := repository.NewUserRepository().FindByID(studentUserID)
	if err != nil || user == nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "user not found")
	}

	// Get all achievements for student
	achievements, err := s.pgRepo.FindByStudentID(studentUserID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to retrieve achievements")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build detailed report
	report := s.buildStudentReport(ctx, student, user, achievements)
	return utils.SuccessResponse(c, "student report retrieved successfully", report)
}

// buildStudentReport builds comprehensive student achievement report
func (s *achievementServiceImpl) buildStudentReport(ctx context.Context, student *models.Student, user *models.User, achievements []models.AchievementReference) fiber.Map {
	// Basic student info
	report := fiber.Map{
		"student_id":    student.ID,
		"user_id":       student.UserID,
		"nim":           student.StudentID,
		"name":          user.FullName,
		"email":         user.Email,
		"program_study": student.ProgramStudy,
		"academic_year": student.AcademicYear,
	}

	// Count achievements by status
	statusCount := map[string]int64{
		"draft":     0,
		"submitted": 0,
		"verified":  0,
		"rejected":  0,
	}

	// Count by type
	typeCount := map[string]int64{
		"competition":   0,
		"publication":   0,
		"organization":  0,
		"certification": 0,
	}

	// Competition levels
	levelCount := map[string]int64{
		"school":        0,
		"city":          0,
		"provincial":    0,
		"national":      0,
		"international": 0,
	}

	// Detailed achievements with mongo details
	var detailedAchievements []fiber.Map

	for _, ach := range achievements {
		statusCount[ach.Status]++

		mongoAch, err := s.mongoRepo.FindByID(ctx, ach.MongoAchievementID)
		if err == nil && mongoAch != nil {
			typeCount[mongoAch.AchievementType]++

			// Count competition levels
			if mongoAch.AchievementType == "competition" && mongoAch.Details != nil {
				if level, exists := mongoAch.Details["competition_level"]; exists {
					if levelStr, ok := level.(string); ok {
						levelCount[levelStr]++
					}
				}
			}

			// Add to detailed list
			detailedAchievements = append(detailedAchievements, fiber.Map{
				"achievement_id": ach.ID,
				"title":          mongoAch.Title,
				"type":           mongoAch.AchievementType,
				"description":    mongoAch.Description,
				"status":         ach.Status,
				"points":         mongoAch.Points,
				"details":        mongoAch.Details,
				"created_at":     mongoAch.CreatedAt,
				"updated_at":     mongoAch.UpdatedAt,
			})
		}
	}

	// Calculate statistics
	total := int64(len(achievements))
	verificationRate := 0.0
	if total > 0 {
		verificationRate = float64(statusCount["verified"]) / float64(total) * 100
	}

	// Build final report
	report["statistics"] = fiber.Map{
		"total":             total,
		"verified":          statusCount["verified"],
		"submitted":         statusCount["submitted"],
		"draft":             statusCount["draft"],
		"rejected":          statusCount["rejected"],
		"verification_rate": verificationRate,
	}
	report["by_type"] = typeCount
	report["by_competition_level"] = levelCount
	report["achievements"] = detailedAchievements

	return report
}
