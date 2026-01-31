package services

import (
	"lottery-system/config"
	"lottery-system/constants"
	"lottery-system/models"
	"lottery-system/repositories"
	"lottery-system/utils"

	"gorm.io/gorm"
)

// DrawService handles lottery draw operations
type DrawService struct {
	drawRepo    *repositories.DrawRepository
	prizeRepo   *repositories.PrizeRepository
	userRepo    *repositories.UserRepository
	companyRepo *repositories.CompanyRepository
}

// NewDrawService creates a new draw service
func NewDrawService() *DrawService {
	return &DrawService{
		drawRepo:    repositories.NewDrawRepository(),
		prizeRepo:   repositories.NewPrizeRepository(),
		userRepo:    repositories.NewUserRepository(),
		companyRepo: repositories.NewCompanyRepository(),
	}
}

// DrawPrize executes a lottery draw for a specified user
// Parameters:
//   - userID: The user ID
//   - companyID: The company ID
//   - levelID: The prize level ID (0 for random draw)
//   - ip: The client IP address
//
// Returns: Draw record and error
func (s *DrawService) DrawPrize(userID, companyID, levelID int, ip string) (*models.DrawRecord, error) {
	// Get user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, utils.NewNotFoundError("用户")
	}

	// Check if user can draw
	if err := s.CheckUserCanDraw(user); err != nil {
		return nil, err
	}

	// Check if user belongs to company
	if user.CompanyID != companyID {
		return nil, utils.NewAuthorizationError(constants.ErrPermissionDenied)
	}

	// Execute draw
	if levelID == 0 {
		return s.DrawRandom(user, ip, companyID)
	}
	return s.DrawWithLevel(user, ip, companyID, levelID)
}

// DrawRandom executes a random draw from all available prize levels
func (s *DrawService) DrawRandom(user *models.User, ip string, companyID int) (*models.DrawRecord, error) {
	// Check if user has drawn
	if user.HasDrawn {
		return nil, utils.NewBusinessLogicError(constants.ErrUserAlreadyDrawn)
	}

	// Begin transaction
	tx := config.DB.Begin()

	// Get all active prize levels with stock
	levels, err := s.prizeRepo.FindActiveLevelsByCompany(companyID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Filter levels with available stock
	availableLevels := make([]models.PrizeLevel, 0)
	totalProbability := 0.0
	for _, level := range levels {
		if level.UsedStock < level.TotalStock {
			availableLevels = append(availableLevels, level)
			totalProbability += level.Probability
		}
	}

	// Check if any prizes available
	if len(availableLevels) == 0 {
		tx.Rollback()
		return nil, utils.NewBusinessLogicError(constants.ErrNoPrizesAvailable)
	}

	// Weighted random selection
	selectedLevel := s.selectLevelByProbability(availableLevels, totalProbability)
	if selectedLevel == nil {
		tx.Rollback()
		return nil, utils.NewBusinessLogicError(constants.ErrOperationFailed)
	}

	// Execute draw with selected level
	record, err := s.executeDraw(tx, user, selectedLevel.ID, ip, companyID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Reload with associations
	return s.drawRepo.FindByIDWithPreload(record.ID)
}

// DrawWithLevel executes a draw for a specific prize level
func (s *DrawService) DrawWithLevel(user *models.User, ip string, companyID, levelID int) (*models.DrawRecord, error) {
	// Check if user has drawn
	if user.HasDrawn {
		return nil, utils.NewBusinessLogicError(constants.ErrUserAlreadyDrawn)
	}

	// Get prize level
	level, err := s.prizeRepo.FindActiveLevelByID(levelID, companyID)
	if err != nil {
		return nil, utils.NewNotFoundError("奖项等级")
	}

	// Check stock
	if level.UsedStock >= level.TotalStock {
		return nil, utils.NewBusinessLogicError(constants.ErrPrizeOutOfStock)
	}

	// Begin transaction
	tx := config.DB.Begin()

	// Execute draw
	record, err := s.executeDraw(tx, user, levelID, ip, companyID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	// Reload with associations
	return s.drawRepo.FindByIDWithPreload(record.ID)
}

// DrawMultiple executes draws for multiple users
func (s *DrawService) DrawMultiple(userIDs []int, companyID, levelID, count int, ip string) ([]models.DrawRecord, error) {
	if len(userIDs) == 0 {
		return nil, utils.NewValidationError(constants.ErrNoUsersAvailable)
	}

	// Limit count to number of users
	if count > len(userIDs) {
		count = len(userIDs)
	}

	// Select random users
	selectedIndices := utils.RandomIndices(len(userIDs), count)
	var records []models.DrawRecord

	for _, index := range selectedIndices {
		userID := userIDs[index]

		// Get user
		user, err := s.userRepo.FindByID(userID)
		if err != nil {
			continue // Skip invalid users
		}

		// Execute draw
		var record *models.DrawRecord
		if levelID == 0 {
			record, err = s.DrawRandom(user, ip, companyID)
		} else {
			record, err = s.DrawWithLevel(user, ip, companyID, levelID)
		}

		if err != nil {
			continue // Skip failed draws
		}

		records = append(records, *record)
	}

	if len(records) == 0 {
		return nil, utils.NewBusinessLogicError(constants.ErrOperationFailed)
	}

	return records, nil
}

// executeDraw executes the actual draw operation within a transaction
func (s *DrawService) executeDraw(tx *gorm.DB, user *models.User, levelID int, ip string, companyID int) (*models.DrawRecord, error) {
	// Get prize level
	var level models.PrizeLevel
	if err := tx.Where("id = ? AND company_id = ?", levelID, companyID).First(&level).Error; err != nil {
		return nil, err
	}

	// Check stock again (transaction-level check)
	if level.UsedStock >= level.TotalStock {
		return nil, utils.NewBusinessLogicError(constants.ErrPrizeOutOfStock)
	}

	// Get any prize for the level
	prize, err := s.prizeRepo.FindAnyPrizeByLevel(levelID)
	if err != nil {
		return nil, err
	}

	// Create draw record
	record := &models.DrawRecord{
		CompanyID: companyID,
		UserID:    user.ID,
		LevelID:   levelID,
		PrizeID:   prize.ID,
		IP:        ip,
	}

	if err := tx.Create(record).Error; err != nil {
		return nil, err
	}

	// Update prize level stock
	if err := tx.Model(&level).Update("used_stock", level.UsedStock+1).Error; err != nil {
		return nil, err
	}

	// Update user status
	if err := tx.Model(user).Update("has_drawn", true).Error; err != nil {
		return nil, err
	}

	return record, nil
}

// selectLevelByProbability selects a prize level using weighted random selection
func (s *DrawService) selectLevelByProbability(levels []models.PrizeLevel, totalProbability float64) *models.PrizeLevel {
	if len(levels) == 0 {
		return nil
	}

	// Generate random value
	randomValue := utils.RandomFloat() * totalProbability

	// Select level
	cumulativeProbability := 0.0
	for i := range levels {
		cumulativeProbability += levels[i].Probability
		if randomValue <= cumulativeProbability {
			return &levels[i]
		}
	}

	// Fallback to last level
	return &levels[len(levels)-1]
}

// CheckUserCanDraw checks if a user can draw
func (s *DrawService) CheckUserCanDraw(user *models.User) error {
	if user.HasDrawn {
		return utils.NewBusinessLogicError(constants.ErrUserAlreadyDrawn)
	}
	return nil
}

// GetDrawRecords gets draw records for a company
func (s *DrawService) GetDrawRecords(companyID int, limit int) ([]models.DrawRecord, error) {
	return s.drawRepo.FindByCompany(companyID, limit)
}

// GetMyPrize gets the prize for a specific user
func (s *DrawService) GetMyPrize(userID, companyID int) (*models.DrawRecord, error) {
	// Get user
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, utils.NewNotFoundError("用户")
	}

	if !user.HasDrawn {
		return nil, utils.NewBusinessLogicError("您还没有参与抽奖")
	}

	// Get draw record
	record, err := s.drawRepo.FindByCompanyAndUser(companyID, userID)
	if err != nil {
		return nil, utils.NewNotFoundError("抽奖记录")
	}

	return record, nil
}

// GetUserStats gets user statistics for a company
func (s *DrawService) GetUserStats(companyID int) (map[string]interface{}, error) {
	// Get available users
	availableUsers, err := s.userRepo.FindAvailableUsers(companyID)
	if err != nil {
		return nil, err
	}

	// Count total users
	total, err := s.userRepo.CountByCompany(companyID)
	if err != nil {
		return nil, err
	}

	// Count drawn users
	drawn, err := s.userRepo.CountByCompanyAndStatus(companyID, true)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_users":     total,
		"available_users": len(availableUsers),
		"drawn_users":     drawn,
		"undrawn_users":   total - drawn,
	}, nil
}
