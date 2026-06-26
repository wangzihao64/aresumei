package service

import (
	"aresumei/dao"
	"aresumei/pkg/e"
	"aresumei/serizlizer"
	"context"
)

type ResumeService struct {
	ResumeId uint `json:"resume_id" form:"resume_id" binding:"required"`
}

func (r *ResumeService) Resume(ctx context.Context, id uint) serizlizer.Response {
	code := e.Success
	if r.ResumeId == 0 {
		code = e.ErrorNotExistResume
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  "resume_id 不能为空",
		}
	}
	resumeDao := dao.NewResumeDao(ctx)
	resume_, err := resumeDao.GetResumeByIdAndUserId(r.ResumeId, id)
	if err != nil {
		code = e.ErrorNotExistResume
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serizlizer.Response{
		Status: code,
		Data: map[string]interface{}{
			"resume_id": resume_.ID,
			"status":    resume_.Status,
			"report":    resume_.Report,
			"error":     resume_.ErrorMessage,
		},
		Msg: e.GetMsg(code),
	}
}
