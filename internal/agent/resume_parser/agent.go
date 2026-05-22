package resume_parser

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"aresumei/internal/pkg/textutil"
	appschema "aresumei/internal/schema"

	"github.com/cloudwego/eino/components/model"
	einoschema "github.com/cloudwego/eino/schema"
)

var (
	ErrEmptyResumeText = errors.New("resume text is empty")
	ErrEmptyLLMOutput  = errors.New("llm output is empty")
)

type Agent struct {
	chatModel     model.BaseChatModel
	maxInputChars int
}

type Config struct {
	ChatModel     model.BaseChatModel
	MaxInputChars int
}

type ParseInput struct {
	RawText string
}

type ParseResult struct {
	Resume  appschema.ResumeSchema
	RawJSON string
}

func NewAgent(cfg Config) (*Agent, error) {
	if cfg.ChatModel == nil {
		return nil, errors.New("chat model is required")
	}

	if cfg.MaxInputChars <= 0 {
		cfg.MaxInputChars = 60000
	}

	return &Agent{
		chatModel:     cfg.ChatModel,
		maxInputChars: cfg.MaxInputChars,
	}, nil
}

func (a *Agent) Parse(ctx context.Context, input ParseInput) (*ParseResult, error) {
	rawText := strings.TrimSpace(input.RawText)
	if rawText == "" {
		return nil, ErrEmptyResumeText
	}

	rawText = textutil.TruncateByRune(rawText, a.maxInputChars)

	messages := []*einoschema.Message{
		{
			Role:    einoschema.System,
			Content: systemPrompt(),
		},
		{
			Role:    einoschema.User,
			Content: buildUserPrompt(rawText),
		},
	}

	resp, err := a.chatModel.Generate(
		ctx,
		messages,
		model.WithTemperature(0),
	)
	fmt.Println(resp.Content)
	if err != nil {
		return nil, fmt.Errorf("generate resume schema: %w", err)
	}

	content := strings.TrimSpace(resp.Content)
	if content == "" {
		return nil, ErrEmptyLLMOutput
	}

	jsonText, err := textutil.ExtractJSONObject(content)
	if err != nil {
		return nil, fmt.Errorf("extract resume json: %w", err)
	}

	var resume appschema.ResumeSchema
	if err := json.Unmarshal([]byte(jsonText), &resume); err != nil {
		return nil, fmt.Errorf("unmarshal resume json: %w", err)
	}

	resume.RawText = rawText

	if err := resume.Validate(); err != nil {
		return nil, fmt.Errorf("validate resume schema: %w", err)
	}

	return &ParseResult{
		Resume:  resume,
		RawJSON: jsonText,
	}, nil
}

func systemPrompt() string {
	return `You are a resume parsing agent.

Your task is to convert raw resume text into a strict JSON object.

Rules:
1. Output JSON only. Do not output markdown, comments, or explanations.
2. Do not invent facts that are not present in the resume.
3. If a field is unknown, use an empty string or an empty array.
4. Preserve the candidate's real experience and achievements.
5. Use the exact JSON field names required by the schema.
6. Dates should use YYYY-MM or YYYY-MM-DD when possible.
7. Highlights should be concrete resume bullet points.
8. Technologies should contain explicit tools, languages, frameworks, databases, cloud services, or platforms mentioned in the resume.`
}

func buildUserPrompt(rawText string) string {
	return fmt.Sprintf(`Convert the following resume text into this JSON shape:

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

%s`, rawText)
}
