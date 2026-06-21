package service

import (
	"aresumei/dao"
	"aresumei/internal/agent/resume"
	"aresumei/pkg/e"
	"aresumei/pkg/util"
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
	resume_, err := resumeDao.GetResumeById(r.ResumeId)
	if err != nil {
		code = e.ErrorNotExistResume
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	resumePath := resume_.FilePath
	pdfstring, _ := util.ReadPdfToString(resumePath)
	resume.Execute(ctx, pdfstring)
	return serizlizer.Response{
		Status: code,
		Data:   e.GetMsg(code),
	}
}
