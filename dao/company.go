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
