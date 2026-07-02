package dao

import (
	"aresumei/model"
	"context"

	"gorm.io/gorm"
)

type CompanyDao struct {
	*gorm.DB
}

func NewCompanyDao(ctx context.Context) *CompanyDao {
	return &CompanyDao{NewDBClient(ctx)}
}

func (u *CompanyDao) CreateCompanyInfo(company *model.CompanyInfo) error {
	return u.Model(&model.CompanyInfo{}).Create(company).Error
}

func (u *CompanyDao) GetCompanyInfoByIdAndUserId(companyId, userId uint) (company *model.CompanyInfo, err error) {
	err = u.Model(&model.CompanyInfo{}).Where("id = ? AND user_id = ?", companyId, userId).First(&company).Error
	return
}

func (u *CompanyDao) GetCompanyReportByIdAndUserId(companyId, userId uint) (company *model.CompanyInfo, err error) {
	err = u.Model(&model.CompanyInfo{}).
		Select("id", "user_id", "status", "report", "error_message").
		Where("id = ? AND user_id = ?", companyId, userId).
		First(&company).Error
	return
}

func (u *CompanyDao) UpdateCompanyInfoById(companyId uint, values map[string]interface{}) error {
	return u.Model(&model.CompanyInfo{}).Where("id = ?", companyId).Updates(values).Error
}

func (u *CompanyDao) UpdateCompanyLLMFilePath(companyId uint, filePath string) error {
	return u.UpdateCompanyInfoById(companyId, map[string]interface{}{
		"llm_file_path": filePath,
	})
}

func (u *CompanyDao) MarkCompanyReportProcessing(companyId uint) error {
	return u.Model(&model.CompanyInfo{}).Where("id = ?", companyId).Updates(map[string]interface{}{
		"status":        model.CompanyStatusProcessing,
		"error_message": "",
	}).Error
}

func (u *CompanyDao) MarkCompanyReportCompleted(companyId uint, report string) error {
	return u.Model(&model.CompanyInfo{}).Where("id = ?", companyId).Updates(map[string]interface{}{
		"status":        model.CompanyStatusCompleted,
		"error_message": "",
	}).Error
}

func (u *CompanyDao) MarkCompanyReportFailed(companyId uint, message string) error {
	return u.Model(&model.CompanyInfo{}).Where("id = ?", companyId).Updates(map[string]interface{}{
		"status":        model.CompanyStatusFailed,
		"error_message": message,
	}).Error
}
