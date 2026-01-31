package services

import (
	"lottery-system/constants"
	"lottery-system/models"
	"lottery-system/repositories"
	"lottery-system/utils"
	"lottery-system/validators"
)

// CompanyService handles company business logic
type CompanyService struct {
	companyRepo *repositories.CompanyRepository
}

// NewCompanyService creates a new company service
func NewCompanyService() *CompanyService {
	return &CompanyService{
		companyRepo: repositories.NewCompanyRepository(),
	}
}

// GetCompanyByCode gets a company by code (active only)
func (s *CompanyService) GetCompanyByCode(code string) (*models.Company, error) {
	if code == "" {
		code = constants.DefaultCompanyCode
	}

	return s.companyRepo.FindByCodeWithActive(code)
}

// GetCompanyByID gets a company by ID
func (s *CompanyService) GetCompanyByID(id int) (*models.Company, error) {
	return s.companyRepo.FindByID(id)
}

// CreateCompany creates a new company
func (s *CompanyService) CreateCompany(company *models.Company) error {
	// Validate code
	if err := validators.ValidateCompanyCode(company.Code); err != nil {
		return err
	}

	// Check if code already exists
	exists, err := s.companyRepo.ExistsByCode(company.Code)
	if err != nil {
		return err
	}
	if exists {
		return utils.NewBusinessLogicError(constants.ErrCompanyCodeExists)
	}

	return s.companyRepo.Create(company)
}

// UpdateCompany updates a company
func (s *CompanyService) UpdateCompany(company *models.Company) error {
	// Validate code if changed
	if company.Code != "" {
		if err := validators.ValidateCompanyCode(company.Code); err != nil {
			return err
		}

		// Check if code already exists (excluding current company)
		exists, err := s.companyRepo.ExistsByCodeAndID(company.Code, company.ID)
		if err != nil {
			return err
		}
		if exists {
			return utils.NewBusinessLogicError(constants.ErrCompanyCodeExists)
		}
	}

	return s.companyRepo.Update(company)
}

// DeleteCompany deletes a company
func (s *CompanyService) DeleteCompany(id int) error {
	// Check if company exists
	_, err := s.companyRepo.FindByID(id)
	if err != nil {
		return utils.NewNotFoundError("公司")
	}

	return s.companyRepo.Delete(id)
}

// GetAllCompanies gets all companies
func (s *CompanyService) GetAllCompanies() ([]models.Company, error) {
	return s.companyRepo.FindAll()
}

// GetCompanyStats gets company statistics
func (s *CompanyService) GetCompanyStats() (map[string]interface{}, error) {
	// Count total companies
	total, err := s.companyRepo.Count()
	if err != nil {
		return nil, err
	}

	// Count active companies
	active, err := s.companyRepo.CountActive()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"total_companies":    total,
		"active_companies":   active,
		"inactive_companies": total - active,
	}, nil
}
