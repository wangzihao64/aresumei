package model

import "gorm.io/gorm"

const (
	ResumeStatusUploaded   = "uploaded"
	ResumeStatusProcessing = "processing"
	ResumeStatusCompleted  = "completed"
	ResumeStatusFailed     = "failed"
)

type Resume struct {
	*gorm.Model
	UserID       uint   `gorm:"not null;index"`
	User         User   `gorm:"foreignkey:UserID"`
	OriginalName string `gorm:"not null"`                  //用户上传的原始的简历名
	FilePath     string `gorm:"not null"`                  //简历保存的路径
	Status       string `gorm:"not null;default:uploaded"` //uploaded/processing/completed/failed
	Report       string `gorm:"type:longtext"`             //面试报告
	ErrorMessage string `gorm:"type:text"`                 //生成失败原因
}
