package model

import "gorm.io/gorm"

const (
	CompanyStatusUploaded   = "uploaded"
	CompanyStatusProcessing = "processing"
	CompanyStatusCompleted  = "completed"
	CompanyStatusFailed     = "failed"
)

type CompanyInfo struct {
	*gorm.Model
	UserID       uint   `gorm:"not null;index"`
	User         User   `gorm:"foreignkey:UserID"`
	Name         string `gorm:"size:128"`                  // 公司名或岗位名
	Content      string `gorm:"type:longtext"`             // 公司介绍、JD、岗位要求等
	FilePath     string `gorm:"size:512"`                  // 原始公司信息文本保存路径
	Status       string `gorm:"not null;default:uploaded"` //uploaded/processing/completed/failed
	ErrorMessage string `gorm:"type:text"`                 //生成失败原因
	LLMFilePath  string `gorm:"not null"`                  //LLM生成报告的位置
}
