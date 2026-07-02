package service

import (
	"aresumei/dao"
	agentoptimize "aresumei/internal/agent/resume"
	"aresumei/pkg/e"
	"aresumei/pkg/util"
	"aresumei/serizlizer"
	"context"
)

type ResumeService struct {
	ResumeId uint `json:"resume_id" form:"resume_id" binding:"required"`
}
type OptimizeResume struct {
	ResumeId  uint `json:"resume_id" form:"resume_id" binding:"required"`
	CompanyId uint `json:"company_id" form:"company_id" binding:"required"`
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
	report, _ := util.ReadText(resume_.LLMFilePath)
	return serizlizer.Response{
		Status: code,
		Data: map[string]interface{}{
			"resume_id": resume_.ID,
			"status":    resume_.Status,
			"report":    report,
			"error":     resume_.ErrorMessage,
		},
		Msg: e.GetMsg(code),
	}
}

func (o *OptimizeResume) Optimize(ctx context.Context, id uint) serizlizer.Response {
	code := e.Success
	if o.ResumeId == 0 {
		code = e.ErrorNotExistResume
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  "resume_id 不能为空",
		}
	}
	if o.CompanyId == 0 {
		code = e.ErrorNotExistCompany
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  "company_id 不能为空",
		}
	}
	resumeDao := dao.NewResumeDao(ctx)
	resume_, err := resumeDao.GetResumeByIdAndUserId(o.ResumeId, id)
	if err != nil {
		code = e.ErrorNotExistResume
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	companyDao := dao.NewCompanyDao(ctx)
	company_, err := companyDao.GetCompanyInfoByIdAndUserId(o.CompanyId, id)
	if err != nil {
		code = e.ErrorNotExistCompany
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  "company_id 不能为空",
		}
	}
	companyText, _ := util.ReadText(company_.LLMFilePath)
	resumeText, _ := util.ReadText(resume_.ParsedResumeFilePath)
	optimizeText, _ := agentoptimize.OptimizeResumeJSON(ctx, resumeText, companyText)
	return serizlizer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data: map[string]interface{}{
			"optimize_text": optimizeText,
		},
	}
}
