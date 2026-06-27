package schema

type CompanySkillGroup struct {
	CompanyName      string                    `json:"company_name,omitempty"`
	PositionTitle    string                    `json:"position_title,omitempty"`      //岗位名
	JobLevel         string                    `json:"job_level,omitempty"`           //实习、初级、中级、高级、专家等
	JobSummary       string                    `json:"job_summary,omitempty"`         //岗位整体描述
	Responsibilities []string                  `json:"responsibilities,omitempty"`    //工作职责
	MustHaveSkills   []companySkillRequirement `json:"must_have_skills,omitempty"`    //硬性技能要求
	NiceToHaveSkills []companySkillRequirement `json:"nice_to_have_skills,omitempty"` //加分项
	Technologies     []string                  `json:"technologies,omitempty"`        //技术关键词扁平列表
	Domains          []string                  `json:"domains,omitempty"`             //业务领域，比如 SaaS、金融、电商、AI、招聘系统
	Experience       []string                  `json:"experience,omitempty"`          //年限
	Education        []string                  `json:"education,omitempty"`
	SoftSkills       []string                  `json:"soft_skills,omitempty"`  //沟通、协作、owner 意识等软技能
	Keywords         []string                  `json:"keywords,omitempty"`     //简历优化时应该优先出现的关键词
	ResumeFocus      []string                  `json:"resume_focus,omitempty"` //优化简历时应该重点突出的方向
}
type companySkillRequirement struct {
	Name       string `json:"name,omitempty"`
	Category   string `json:"category,omitempty"`
	Importance string `json:"importance,omitempty"`
	Evidence   string `json:"evidence,omitempty"`
}
