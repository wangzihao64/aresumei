package resume

import (
	"aresumei/internal/infrastructure/llm"
	"aresumei/internal/schema"
	"aresumei/pkg/util"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/openai/openai-go/v3"
)

const defaultBaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
const defaultModel = "deepseek-v3.2"

type interviewReportInput struct {
	Skills             []schema.ResumeSkillGroup    `json:"skills,omitempty"`
	WorkExperiences    []interviewWorkExperience    `json:"work_experiences,omitempty"`
	ProjectExperiences []interviewProjectExperience `json:"project_experiences,omitempty"`
}

type interviewWorkExperience struct {
	Company      string   `json:"company,omitempty"`
	Position     string   `json:"position,omitempty"`
	Highlights   []string `json:"highlights,omitempty"`
	Technologies []string `json:"technologies,omitempty"`
}

type interviewProjectExperience struct {
	Name         string   `json:"name,omitempty"`
	Role         string   `json:"role,omitempty"`
	Description  string   `json:"description,omitempty"`
	Highlights   []string `json:"highlights,omitempty"`
	Technologies []string `json:"technologies,omitempty"`
}

func Execute(ctx context.Context, pdf string) (string, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("DASHSCOPE_API_KEY")
	}
	if apiKey == "" {
		return "", errors.New("missing OPENAI_API_KEY or DASHSCOPE_API_KEY")
	}
	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	ai := llm.NewOpenAIClient(llm.Config{
		APIKey:  apiKey,
		BaseURL: baseURL,
	})
	parsedResume, err := parseResume(ctx, ai, pdf)
	if err != nil {
		return "", err
	}
	return generateInterviewReport(ctx, ai, parsedResume)
}

func parseResume(ctx context.Context, ai *llm.OpenAiClient, pdf string) (schema.Resume, error) {
	systemPrompt := "You are a resume parsing agent.\nYour task is to convert raw resume text into a strict JSON object.\nRules:\n1. Output JSON only. Do not output markdown, comments, or explanations.\n2. Do not invent facts that are not present in the resume.\n3. If a field is unknown, use an empty string or an empty array.\n4. Preserve the candidate's real experience and achievements.\n5. Use the exact JSON field names required by the schema.\n6. Dates should use YYYY-MM or YYYY-MM-DD when possible.\n7. Highlights should be concrete resume bullet points.\n8. Technologies should contain explicit tools, languages, frameworks, databases, cloud services, or platforms mentioned in the resume."
	userPromot := fmt.Sprintf(`Convert the following resume text into this JSON shape:
{
  "basics": {
    "full_name": "",
    "email": "",
    "phone": "",
    "location": "",
    "github": "",
    "linkedin": "",
    "website": ""
  },
  "summary": "",
  "skills": [
    {
      "category": "",
      "skills": []
    }
  ],
  "work_experiences": [
    {
      "company": "",
      "position": "",
      "location": "",
      "start_date": "",
      "end_date": "",
      "is_current": false,
      "description": "",
      "highlights": [],
      "technologies": []
    }
  ],
  "project_experiences": [
    {
      "name": "",
      "role": "",
      "start_date": "",
      "end_date": "",
      "description": "",
      "highlights": [],
      "technologies": [],
      "links": []
    }
  ],
  "education": [
    {
      "school": "",
      "degree": "",
      "major": "",
      "start_date": "",
      "end_date": "",
      "location": ""
    }
  ],
  "certifications": [
    {
      "name": "",
      "issuer": "",
      "date": "",
      "expires_at": ""
    }
  ],
  "languages": [
    {
      "name": "",
      "proficiency": ""
    }
  ]
}
Resume text:
%s`, pdf)
	param := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userPromot),
		},
		Model: defaultModel,
	}
	completion, err := ai.Client.Chat.Completions.New(ctx, param)
	if err != nil {
		return schema.Resume{}, err
	}
	if len(completion.Choices) == 0 {
		return schema.Resume{}, errors.New("empty chat completion choices")
	}
	var resume schema.Resume
	rawcontent := completion.Choices[0].Message.Content
	content := util.CleanJsonContent(rawcontent)
	if err := json.Unmarshal([]byte(content), &resume); err != nil {
		return schema.Resume{}, err
	}
	return resume, nil
}

func generateInterviewReport(ctx context.Context, ai *llm.OpenAiClient, parsedResume schema.Resume) (string, error) {
	reportInput := buildInterviewReportInput(parsedResume)
	payload, err := json.Marshal(reportInput)
	if err != nil {
		return "", err
	}
	systemPrompt := "You are a senior technical interviewer.\nGenerate a polished plain-text interview report for the user based only on the provided structured resume JSON.\nFocus on work_experiences.highlights, work_experiences.technologies, skills, project_experiences.description, project_experiences.highlights, and project_experiences.technologies.\nDo not invent facts or evaluate information that is not present in the JSON.\nDo not use Markdown syntax, including # headings, *, -, code fences, tables, or block quotes."
	userPrompt := fmt.Sprintf(`请基于以下结构化简历信息生成一份中文面试报告。

报告要求：
1. 输出适合直接展示给用户的纯文本，不要输出 JSON，不要使用 Markdown。
2. 不要使用 #、*、-、>、表格、代码块等格式符号。
3. 用清晰的小标题和编号组织内容，小标题格式为“1. 技术栈匹配度”。
4. 每个小标题下用自然中文短句分段说明，不要使用项目符号列表。
5. 重点覆盖技术栈匹配度、工作经历亮点、项目经历亮点、可追问方向、候选人风险点、建议面试问题。
6. 所有判断必须能从输入 JSON 中找到依据；信息不足时明确写“简历未体现”。
7. 面试问题要结合候选人的技术、工作亮点和项目描述，不要生成泛泛的问题。

结构化简历信息：
%s`, string(payload))
	param := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userPrompt),
		},
		Model: defaultModel,
	}
	completion, err := ai.Client.Chat.Completions.New(ctx, param)
	if err != nil {
		return "", err
	}
	if len(completion.Choices) == 0 {
		return "", errors.New("empty interview report completion choices")
	}
	return cleanInterviewReport(completion.Choices[0].Message.Content), nil
}

var (
	markdownHeadingRE     = regexp.MustCompile(`(?m)^[ \t]{0,3}#{1,6}[ \t]*`)
	markdownListMarkerRE  = regexp.MustCompile(`(?m)^[ \t]*[*+-][ \t]+`)
	markdownQuoteRE       = regexp.MustCompile(`(?m)^[ \t]*>[ \t]?`)
	markdownBoldItalicRE  = regexp.MustCompile(`[*_]{1,3}([^*_]+)[*_]{1,3}`)
	markdownInlineCodeRE  = regexp.MustCompile("`([^`]+)`")
	markdownFenceMarkerRE = regexp.MustCompile("(?m)^```[\\w-]*\\s*$")
	blankLinesRE         = regexp.MustCompile(`\n{3,}`)
)

func cleanInterviewReport(report string) string {
	report = strings.ReplaceAll(report, "\r\n", "\n")
	report = markdownFenceMarkerRE.ReplaceAllString(report, "")
	report = markdownHeadingRE.ReplaceAllString(report, "")
	report = markdownListMarkerRE.ReplaceAllString(report, "")
	report = markdownQuoteRE.ReplaceAllString(report, "")
	report = markdownBoldItalicRE.ReplaceAllString(report, "$1")
	report = markdownInlineCodeRE.ReplaceAllString(report, "$1")
	report = strings.ReplaceAll(report, "|", " ")
	report = blankLinesRE.ReplaceAllString(report, "\n\n")
	return strings.TrimSpace(report)
}

func buildInterviewReportInput(parsedResume schema.Resume) interviewReportInput {
	input := interviewReportInput{
		Skills:             make([]schema.ResumeSkillGroup, 0, len(parsedResume.Skills)),
		WorkExperiences:    make([]interviewWorkExperience, 0, len(parsedResume.WorkExperiences)),
		ProjectExperiences: make([]interviewProjectExperience, 0, len(parsedResume.ProjectExperiences)),
	}
	input.Skills = append(input.Skills, parsedResume.Skills...)
	for _, work := range parsedResume.WorkExperiences {
		input.WorkExperiences = append(input.WorkExperiences, interviewWorkExperience{
			Company:      work.Company,
			Position:     work.Position,
			Highlights:   append([]string(nil), work.Highlights...),
			Technologies: append([]string(nil), work.Technologies...),
		})
	}
	for _, project := range parsedResume.ProjectExperiences {
		input.ProjectExperiences = append(input.ProjectExperiences, interviewProjectExperience{
			Name:         project.Name,
			Role:         project.Role,
			Description:  project.Description,
			Highlights:   append([]string(nil), project.Highlights...),
			Technologies: append([]string(nil), project.Technologies...),
		})
	}
	return input
}
