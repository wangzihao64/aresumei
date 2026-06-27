package company

import (
	"aresumei/internal/infrastructure/llm"
	"aresumei/internal/schema"
	"aresumei/pkg/util"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/openai/openai-go/v3"
)

const defaultBaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"
const defaultModel = "deepseek-v3.2"

func Execute(ctx context.Context, text string) (string, error) {
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
	parseCompany, err := parseCompany(ctx, ai, text)
	if err != nil {
		return "", err
	}
	content, err := json.Marshal(parseCompany)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func parseCompany(ctx context.Context, ai *llm.OpenAiClient, text string) (schema.CompanySkillGroup, error) {
	systemPrompt := "You are a job requirement parsing agent.\nYour task is to convert raw company, job description, or recruitment requirement text into a strict JSON object.\nRules:\n1. Output JSON only. Do not output markdown, comments, or explanations.\n2. Do not invent requirements that are not present in the source text.\n3. If a field is unknown, use an empty string or an empty array.\n4. Preserve the original meaning of the company's requirements.\n5. Separate must-have skills from nice-to-have skills based on wording such as required, must, familiar with, preferred, plus, bonus, 加分, 优先, 必须, 熟悉, 掌握.\n6. Technologies should contain explicit tools, languages, frameworks, databases, cloud services, platforms, and engineering practices mentioned in the source text.\n7. Evidence must quote or briefly paraphrase the source phrase that supports the extracted requirement.\n8. Keywords and resume_focus should help rewrite a candidate resume toward this role without fabricating experience."
	userPrompt := fmt.Sprintf(`Convert the following company or job requirement text into this JSON shape:
{
  "company_name": "",
  "position_title": "",
  "job_level": "",
  "job_summary": "",
  "responsibilities": [],
  "must_have_skills": [
    {
      "name": "",
      "category": "",
      "importance": "",
      "evidence": ""
    }
  ],
  "nice_to_have_skills": [
    {
      "name": "",
      "category": "",
      "importance": "",
      "evidence": ""
    }
  ],
  "technologies": [],
  "domains": [],
  "experience": [],
  "education": [],
  "soft_skills": [],
  "keywords": [],
  "resume_focus": []
}

Field guidance:
- company_name: company name if present.
- position_title: role title if present.
- job_level: internship, junior, mid-level, senior, expert, or the original level wording if present.
- job_summary: one concise sentence summarizing the role.
- responsibilities: concrete job duties.
- must_have_skills: hard requirements required for the role.
- nice_to_have_skills: preferred or bonus requirements.
- category: use values like programming_language, framework, database, cache, message_queue, cloud, devops, architecture, testing, ai, domain, soft_skill, other.
- importance: use must_have or nice_to_have.
- technologies: flat list of technical keywords explicitly mentioned.
- domains: business or product domains such as SaaS, finance, e-commerce, AI, recruitment, enterprise software.
- experience: years of experience, project experience, team experience, or industry experience requirements.
- education: degree, major, school, or certification requirements.
- soft_skills: communication, collaboration, ownership, learning ability, documentation, leadership, and similar requirements.
- keywords: important words that should be considered when optimizing a resume for this role.
- resume_focus: directions a resume should emphasize for this company, based only on the source text.

Company or job requirement text:
%s`, text)

	param := openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userPrompt),
		},
		Model: defaultModel,
	}
	completion, err := ai.Client.Chat.Completions.New(ctx, param)
	if err != nil {
		return schema.CompanySkillGroup{}, err
	}
	if len(completion.Choices) == 0 {
		return schema.CompanySkillGroup{}, errors.New("empty chat completion choices")
	}
	content := util.CleanJsonContent(completion.Choices[0].Message.Content)
	var company schema.CompanySkillGroup
	if err := json.Unmarshal([]byte(content), &company); err != nil {
		return schema.CompanySkillGroup{}, err
	}
	return company, nil
}
