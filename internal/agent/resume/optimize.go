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

	"github.com/openai/openai-go/v3"
)

type optimizeResumeInput struct {
	OriginalResume schema.Resume            `json:"original_resume"`
	Company        schema.CompanySkillGroup `json:"company"`
}

func OptimizeResume(ctx context.Context, parsedResume schema.Resume, parsedCompany schema.CompanySkillGroup) (schema.Resume, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("DASHSCOPE_API_KEY")
	}
	if apiKey == "" {
		return schema.Resume{}, errors.New("missing OPENAI_API_KEY or DASHSCOPE_API_KEY")
	}
	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	ai := llm.NewOpenAIClient(llm.Config{
		APIKey:  apiKey,
		BaseURL: baseURL,
	})
	return optimizeResume(ctx, ai, parsedResume, parsedCompany)
}

func OptimizeResumeJSON(ctx context.Context, resumeJSON, companyJSON string) (string, error) {
	var parsedResume schema.Resume
	if err := json.Unmarshal([]byte(util.CleanJsonContent(resumeJSON)), &parsedResume); err != nil {
		return "", err
	}
	var parsedCompany schema.CompanySkillGroup
	if err := json.Unmarshal([]byte(util.CleanJsonContent(companyJSON)), &parsedCompany); err != nil {
		return "", err
	}
	optimized, err := OptimizeResume(ctx, parsedResume, parsedCompany)
	if err != nil {
		return "", err
	}
	content, err := json.Marshal(optimized)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func optimizeResume(ctx context.Context, ai *llm.OpenAiClient, parsedResume schema.Resume, parsedCompany schema.CompanySkillGroup) (schema.Resume, error) {
	input := optimizeResumeInput{
		OriginalResume: parsedResume,
		Company:        parsedCompany,
	}
	payload, err := json.Marshal(input)
	if err != nil {
		return schema.Resume{}, err
	}
	systemPrompt := "You are a senior resume optimization expert.\nYour task is to rewrite an existing candidate resume so it better matches a target company or job requirement.\nRules:\n1. Output JSON only. Do not output markdown, comments, explanations, or code fences.\n2. Use the exact same JSON shape as the input resume schema.\n3. Do not invent jobs, projects, schools, certifications, tools, or metrics that are not present in the original resume.\n4. You may rewrite wording, reorder skills, strengthen bullet points, and emphasize relevant experience from the original resume.\n5. Use the target company requirements to prioritize keywords, technologies, and resume focus.\n6. Keep factual fields stable. Preserve contact information, education, certifications, and languages.\n7. If information is missing, keep it empty instead of fabricating it.\n8. Make the resume concise, concrete, and ATS-friendly."
	userPrompt := fmt.Sprintf(`Please optimize the candidate resume based on the following structured inputs.

Optimization goals:
1. Make the resume more aligned with the target company requirements.
2. Strengthen the summary, skill ordering, work experience highlights, and project experience highlights.
3. Keep all facts truthful and limited to the original resume content.
4. Reorder or rephrase content to surface the most relevant keywords from the target company.
5. Preserve the same JSON field names and output a valid resume JSON object.

Output requirements:
1. Keep basics, education, certifications, and languages factual and unchanged.
2. Rewrite summary, skills, work_experiences, and project_experiences to better match the target role.
3. Do not add new roles, projects, or technologies that are not in the original resume.
4. If the original resume does not contain a field, leave it empty.

Input JSON:
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
		return schema.Resume{}, err
	}
	if len(completion.Choices) == 0 {
		return schema.Resume{}, errors.New("empty optimize resume completion choices")
	}
	content := util.CleanJsonContent(completion.Choices[0].Message.Content)
	var optimized schema.Resume
	if err := json.Unmarshal([]byte(content), &optimized); err != nil {
		return schema.Resume{}, err
	}
	optimized = mergeOptimizedResume(parsedResume, optimized)
	return optimized, nil
}

func mergeOptimizedResume(original, optimized schema.Resume) schema.Resume {
	merged := optimized
	merged.Basics = original.Basics
	merged.Educations = original.Educations
	merged.Certifications = original.Certifications
	merged.Languages = original.Languages

	if merged.Summary == "" {
		merged.Summary = original.Summary
	}
	if len(merged.Skills) == 0 {
		merged.Skills = original.Skills
	}
	if len(merged.WorkExperiences) == 0 {
		merged.WorkExperiences = original.WorkExperiences
	}
	if len(merged.ProjectExperiences) == 0 {
		merged.ProjectExperiences = original.ProjectExperiences
	}
	return merged
}
