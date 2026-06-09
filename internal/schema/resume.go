package schema

type Resume struct {
	Basics             ResumeBasic         `json:"basics"`
	Summary            string              `json:"summary,omitempty"`
	Skills             []ResumeSkillGroup  `json:"skills,omitempty"`
	WorkExperiences    []WorkExperience    `json:"work_experiences,omitempty"`
	ProjectExperiences []ProjectExperience `json:"project_experiences"`
	Educations         []Education         `json:"education"`
	Certifications     []Certification     `json:"certifications,omitempty"`
	Languages          []Language          `json:"languages,omitempty"`
}
type ResumeBasic struct {
	FullName string `json:"full_name,omitempty"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Location string `json:"location,omitempty"`
	Github   string `json:"github,omitempty"`
	LinkedIn string `json:"linkedin,omitempty"`
	Website  string `json:"website,omitempty"`
}
type ResumeSkillGroup struct {
	Category string   `json:"category"`
	Skills   []string `json:"skills"`
}
type WorkExperience struct {
	Company      string   `json:"company"`
	Position     string   `json:"position,omitempty"`
	Location     string   `json:"location,omitempty"`
	StartDate    string   `json:"start_date,omitempty"`
	EndDate      string   `json:"end_date,omitempty"`
	IsCurrent    bool     `json:"is_current,omitempty"`
	Description  string   `json:"description,omitempty"`
	Highlights   []string `json:"highlights,omitempty"`
	Technologies []string `json:"technologies,omitempty"`
}
type ProjectExperience struct {
	Name         string   `json:"name"`
	Role         string   `json:"role"`
	StartDate    string   `json:"start_date"`
	EndDate      string   `json:"end_date"`
	Description  string   `json:"description"`
	Highlights   []string `json:"highlights,omitempty"`
	Technologies []string `json:"technologies,omitempty"`
	Links        []string `json:"links,omitempty"`
}
type Education struct {
	School    string `json:"school"`
	Degree    string `json:"degree,omitempty"`
	Major     string `json:"major,omitempty"`
	StartDate string `json:"start_date,omitempty"`
	EndDate   string `json:"end_date,omitempty"`
	Location  string `json:"location,omitempty"`
}
type Certification struct {
	Name     string `json:"name"`
	Issuer   string `json:"issuer,omitempty"`
	Date     string `json:"date,omitempty"`
	ExpireAt string `json:"expires_at,omitempty"`
}
type Language struct {
	Name        string `json:"name"`
	Proficiency string `json:"proficiency"`
}
