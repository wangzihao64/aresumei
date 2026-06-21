package dao

import (
	"aresumei/model"
	"context"

	"gorm.io/gorm"
)

type ResumeDao struct {
	*gorm.DB
}

func NewResumeDao(ctx context.Context) *ResumeDao {
	return &ResumeDao{NewDBClient(ctx)}
}
func (u *ResumeDao) CreateResume(resume *model.Resume) error {
	return u.Model(&model.Resume{}).Create(resume).Error
}
func (u *ResumeDao) GetResumeById(resumeId uint) (resume *model.Resume, err error) {
	err = u.Model(&model.Resume{}).Where("id = ?", resumeId).First(&resume).Error
	return
}
