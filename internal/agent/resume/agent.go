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

const defaultBaseURL = "https://dashscope.aliyuncs.com/compatible-mode/v1"

func Execute(ctx context.Context, pdf string) error {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("DASHSCOPE_API_KEY")
	}
	if apiKey == "" {
		return errors.New("missing OPENAI_API_KEY or DASHSCOPE_API_KEY")
	}
	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	ai := llm.NewOpenAIClient(llm.Config{
		APIKey:  apiKey,
		BaseURL: baseURL,
	})
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
		Model: "deepseek-v3.2",
	}
	completion, err := ai.Client.Chat.Completions.New(ctx, param)
	if err != nil {
		return err
	}
	if len(completion.Choices) == 0 {
		return errors.New("empty chat completion choices")
	}
	var resume schema.Resume
	rawcontent := completion.Choices[0].Message.Content
	content := util.CleanJsonContent(rawcontent)
	fmt.Println(content)
	if err := json.Unmarshal([]byte(content), &resume); err != nil {
		return err
	}
	fmt.Println(resume)
	return nil
}
