package service

import (
	"aresumei/cache"
	"aresumei/dao"
	agentcompany "aresumei/internal/agent/company"
	agentresume "aresumei/internal/agent/resume"
	"aresumei/model"
	"aresumei/pkg/e"
	"aresumei/pkg/util"
	"aresumei/serizlizer"
	"context"
	"mime/multipart"
)

const resumesPath = "uploads/resumes"
const companyPath = "uploads/company"

type UserService struct {
	Nickname  string `json:"nick_name" form:"nick_name"`
	Username  string `json:"user_name" form:"user_name"`
	Password  string `json:"password" form:"password"`
	Key       string `json:"key" form:"key"`
	Email     string `json:"email" form:"email"`
	EmailCode string `json:"code" form:"code"`
}
type TextService struct {
	Name string `json:"name" form:"name"`
	Text string `json:"text" form:"text" binding:"required"`
}

func (service *UserService) Login(ctx context.Context) serizlizer.Response {
	code := e.Success
	userDao := dao.NewUserDao(ctx)
	user, exist, err := userDao.ExistOrNotByUserName(service.Username)
	if err != nil {
		code = e.Error
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if !exist {
		code = e.ErrorNotExistUser
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  e.GetMsg(code),
		}
	}
	if !user.CheckPassword(service.Password) {
		code = e.ErrorPassword
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  e.GetMsg(code),
		}
	}
	token, err := util.GenerateToken(user.ID, user.Username, 0)
	if err != nil {
		code = e.ErrorGenerateToken
		return serizlizer.Response{
			Status: code,
		}
	}
	return serizlizer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data: serizlizer.TokenData{
			Token: token,
			User:  serizlizer.BuildUser(user),
		},
	}
}
func (service *UserService) Register(ctx context.Context) serizlizer.Response {
	var user model.User
	code := e.Success
	if service.Username == "" {
		code = e.ErrorUsernameIsEmpty
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  e.GetMsg(code),
		}
	}
	if service.Password == "" {
		code = e.ErrorPasswordIsEmpty
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  e.GetMsg(code),
		}
	}
	if service.Email == "" {
		code = e.ErrorEmailIsEmpty
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  e.GetMsg(code),
		}
	}
	rdb := cache.NewRDB(ctx)
	rdbCode, err := rdb.GetVerificationCode(service.Email)
	if err != nil {
		code = e.ErrorRedisGetkey
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if rdbCode != service.EmailCode {
		code = e.ErrorRedisGetkey
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	//todo 目前密钥没有用到
	if service.Key == "" || len(service.Key) != 16 {
		code = e.Error
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  "密钥长度不足",
		}
	}
	userDao := dao.NewUserDao(ctx)
	_, exist, err := userDao.ExistOrNotByUserName(service.Username)
	if err != nil {
		code = e.Error
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if exist {
		code = e.ErrorExistUser
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  "该用户已经存在",
		}
	}
	user = model.User{
		Nickname:      service.Nickname,
		Username:      service.Username,
		Email:         service.Email,
		EmailVerified: true,
	}
	//密码加密
	if err := user.SetPassword(service.Password); err != nil {
		code = e.ErrorFailEncryption
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if err := userDao.CreateUser(&user); err != nil {
		code = e.Error
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serizlizer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}
func (service *UserService) UpLoadResume(ctx context.Context, file *multipart.FileHeader, id uint) serizlizer.Response {
	code := e.Success
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserById(id)
	if err != nil {
		code = e.ErrorNotExistUser
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	//todo 校验文件后缀(目前只判断是否是pdf文件)
	if err := util.IsPdf(file); err != nil {
		code = e.ErrorInvalidFileFormat
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	//创建保存目录
	filePath, err := util.SaveFiles(resumesPath, file, user.Username)
	if err != nil {
		code = e.ErrorSaveFile
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	resumeDao := dao.NewResumeDao(ctx)
	resume := model.Resume{
		UserID:       user.ID,
		OriginalName: file.Filename,
		FilePath:     filePath,
		Status:       model.ResumeStatusProcessing,
	}
	if err := resumeDao.CreateResume(&resume); err != nil {
		code = e.Error
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	go generateResumeReport(resume.ID, resume.FilePath)
	return serizlizer.Response{
		Status: code,
		Data: map[string]interface{}{
			"resume_id": resume.ID,
			"status":    resume.Status,
		},
		Msg: e.GetMsg(code),
	}
}

func generateResumeReport(resumeId uint, filePath string) {
	ctx := context.Background()
	resumeDao := dao.NewResumeDao(ctx)
	if err := resumeDao.MarkReportProcessing(resumeId); err != nil {
		return
	}
	pdfstring, err := util.ReadPdfToString(filePath)
	if err != nil {
		_ = resumeDao.MarkReportFailed(resumeId, err.Error())
		return
	}
	report, err := agentresume.Execute(ctx, pdfstring)
	if err != nil {
		_ = resumeDao.MarkReportFailed(resumeId, err.Error())
		return
	}
	_ = resumeDao.MarkReportCompleted(resumeId, report)
}
func generateCompanyReport(companyId uint, filePath string) {
	ctx := context.Background()
	companyDao := dao.NewCompanyDao(ctx)
	if err := companyDao.MarkCompanyReportProcessing(companyId); err != nil {
		return
	}
	text, err := util.ReadText(filePath)
	if err != nil {
		_ = companyDao.MarkCompanyReportFailed(companyId, err.Error())
		return
	}
	report, err := agentcompany.Execute(ctx, text)
	if err != nil {
		_ = companyDao.MarkCompanyReportFailed(companyId, err.Error())
		return
	}
	_ = companyDao.MarkCompanyReportCompleted(companyId, report)
}

// todo 需要后期解耦合
func (service *TextService) UploadCompany(ctx context.Context, id uint) serizlizer.Response {
	code := e.Success
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserById(id)
	if err != nil {
		code = e.ErrorNotExistUser
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	//保存text到文件里
	filePath, err := util.SaveText(companyPath, service.Text, user.Username)
	if err != nil {
		code = e.ErrorSaveFile
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	company := model.CompanyInfo{
		UserID:   user.ID,
		Name:     service.Name,
		Content:  service.Text,
		FilePath: filePath,
	}
	companyDao := dao.NewCompanyDao(ctx)
	if err := companyDao.CreateCompanyInfo(&company); err != nil {
		code = e.Error
		return serizlizer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	go generateCompanyReport(company.ID, company.FilePath)
	return serizlizer.Response{
		Status: code,
		Data: map[string]interface{}{
			"company_id": company.ID,
			"status":     company.Status,
		},
		Msg: e.GetMsg(code),
	}
}
