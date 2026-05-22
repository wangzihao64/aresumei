package schema

import "errors"

type ResumeSchema struct {
	Basics             ResumeBasics        `json:"basics"`
	Summary            string              `json:"summary"`
	Skills             []ResumeSkillGroup  `json:"skills"`
	WorkExperiences    []WorkExperience    `json:"work_experiences"`
	ProjectExperiences []ProjectExperience `json:"project_experiences"`
	Education          []Education         `json:"education"`
	Certifications     []Certification     `json:"certifications,omitempty"`
	Languages          []Language          `json:"languages,omitempty"`
	RawText            string              `json:"raw_text,omitempty"`
}

type ResumeBasics struct {
	FullName string `json:"full_name"`
	Email    string `json:"email,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Location string `json:"location,omitempty"`
	GitHub   string `json:"github,omitempty"`
	LinkedIn string `json:"linkedin,omitempty"`
	Website  string `json:"website,omitempty"`
}

type ResumeSkillGroup struct {
	Category string   `json:"category"`
	Skills   []string `json:"skills"`
}

type WorkExperience struct {
	Company      string   `json:"company"`
	Position     string   `json:"position"`
	Location     string   `json:"location,omitempty"`
	StartDate    string   `json:"start_date"`
	EndDate      string   `json:"end_date,omitempty"`
	IsCurrent    bool     `json:"is_current"`
	Description  string   `json:"description,omitempty"`
	Highlights   []string `json:"highlights"`
	Technologies []string `json:"technologies,omitempty"`
}

type ProjectExperience struct {
	Name         string   `json:"name"`
	Role         string   `json:"role,omitempty"`
	StartDate    string   `json:"start_date,omitempty"`
	EndDate      string   `json:"end_date,omitempty"`
	Description  string   `json:"description"`
	Highlights   []string `json:"highlights"`
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
	Name      string `json:"name"`
	Issuer    string `json:"issuer,omitempty"`
	Date      string `json:"date,omitempty"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

type Language struct {
	Name        string `json:"name"`
	Proficiency string `json:"proficiency,omitempty"`
}

func (r ResumeSchema) Validate() error {
	if r.Basics.FullName == "" {
		return errors.New("resume basics.full_name is required")
	}

	if len(r.WorkExperiences) == 0 && len(r.ProjectExperiences) == 0 {
		return errors.New("resume must contain at least one work or project experience")
	}

	for _, exp := range r.WorkExperiences {
		if exp.Company == "" {
			return errors.New("work experience company is required")
		}
		if exp.Position == "" {
			return errors.New("work experience position is required")
		}
	}

	for _, project := range r.ProjectExperiences {
		if project.Name == "" {
			return errors.New("project experience name is required")
		}
	}

	return nil
}
