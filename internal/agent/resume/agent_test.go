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

func TestCleanInterviewReportRemovesMarkdownSyntax(t *testing.T) {
	raw := "# 面试报告\n\n## 技术栈匹配度\n* **Go** 与 `MySQL` 经验明确\n- 简历未体现云服务经验\n> 可继续追问项目落地细节\n\n```text\n建议问题\n```\n"

	got := cleanInterviewReport(raw)
	want := "面试报告\n\n技术栈匹配度\nGo 与 MySQL 经验明确\n简历未体现云服务经验\n可继续追问项目落地细节\n\n建议问题"

	if got != want {
		t.Fatalf("cleanInterviewReport() = %q, want %q", got, want)
	}
}
