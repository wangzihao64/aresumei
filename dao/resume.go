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

func (u *ResumeDao) GetResumeByIdAndUserId(resumeId, userId uint) (resume *model.Resume, err error) {
	err = u.Model(&model.Resume{}).Where("id = ? AND user_id = ?", resumeId, userId).First(&resume).Error
	return
}

func (u *ResumeDao) MarkReportProcessing(resumeId uint) error {
	return u.Model(&model.Resume{}).Where("id = ?", resumeId).Updates(map[string]interface{}{
		"status":        model.ResumeStatusProcessing,
		"report":        "",
		"error_message": "",
	}).Error
}

func (u *ResumeDao) MarkReportCompleted(resumeId uint, report string) error {
	return u.Model(&model.Resume{}).Where("id = ?", resumeId).Updates(map[string]interface{}{
		"status":        model.ResumeStatusCompleted,
		"report":        report,
		"error_message": "",
	}).Error
}

func (u *ResumeDao) MarkReportFailed(resumeId uint, message string) error {
	return u.Model(&model.Resume{}).Where("id = ?", resumeId).Updates(map[string]interface{}{
		"status":        model.ResumeStatusFailed,
		"error_message": message,
	}).Error
}
