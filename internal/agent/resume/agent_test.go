package resume

import (
	"aresumei/internal/schema"
	"reflect"
	"testing"
)

func TestBuildInterviewReportInputUsesRequiredResumeFields(t *testing.T) {
	parsedResume := schema.Resume{
		Summary: "should not be included",
		Skills: []schema.ResumeSkillGroup{
			{Category: "Backend", Skills: []string{"Go", "MySQL"}},
		},
		WorkExperiences: []schema.WorkExperience{
			{
				Company:      "AResumeI",
				Position:     "Backend Engineer",
				Description:  "should not be included",
				Highlights:   []string{"Built resume parsing pipeline"},
				Technologies: []string{"Go", "OpenAI"},
			},
		},
		ProjectExperiences: []schema.ProjectExperience{
			{
				Name:         "Interview Report",
				Role:         "Developer",
				Description:  "Generate reports from resume facts",
				Highlights:   []string{"Generated targeted interview questions"},
				Technologies: []string{"Go", "LLM"},
				Links:        []string{"https://example.com"},
			},
		},
	}

	got := buildInterviewReportInput(parsedResume)
	want := interviewReportInput{
		Skills: []schema.ResumeSkillGroup{
			{Category: "Backend", Skills: []string{"Go", "MySQL"}},
		},
		WorkExperiences: []interviewWorkExperience{
			{
				Company:      "AResumeI",
				Position:     "Backend Engineer",
				Highlights:   []string{"Built resume parsing pipeline"},
				Technologies: []string{"Go", "OpenAI"},
			},
		},
		ProjectExperiences: []interviewProjectExperience{
			{
				Name:         "Interview Report",
				Role:         "Developer",
				Description:  "Generate reports from resume facts",
				Highlights:   []string{"Generated targeted interview questions"},
				Technologies: []string{"Go", "LLM"},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("buildInterviewReportInput() = %#v, want %#v", got, want)
	}
}
