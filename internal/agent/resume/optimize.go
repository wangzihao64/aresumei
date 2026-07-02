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
	ParsedResumeToOptimize    schema.Resume            `json:"parsed_resume_to_optimize"`
	TargetCompanyRequirements schema.CompanySkillGroup `json:"target_company_requirements"`
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
		ParsedResumeToOptimize:    parsedResume,
		TargetCompanyRequirements: parsedCompany,
	}
	payload, err := json.Marshal(input)
	if err != nil {
		return schema.Resume{}, err
	}
	systemPrompt := "你是资深简历优化专家。\n你的任务是根据 target_company_requirements 优化 parsed_resume_to_optimize，但只能输出优化后的简历 JSON。\n硬性规则：\n1. 只输出 JSON，不要输出 markdown、解释、注释或代码块。\n2. 输出必须只是一份优化后的简历，不要输出公司要求，不要把公司要求合并进简历。\n3. 使用和 parsed_resume_to_optimize 完全一致的 JSON 字段结构。\n4. 不得编造 parsed_resume_to_optimize 中不存在的工作、项目、学校、证书、工具、指标、技能或经历。\n5. target_company_requirements 只能作为优化方向，用来决定突出哪些原有经历、技能和项目。\n6. 严禁翻译 parsed_resume_to_optimize 中已有的任何文本。中文保持中文，英文保持英文，中英混合保持原样。\n7. 严禁把中文概念替换成英文术语。例如：链表必须仍然是链表，不能变成 linked list；数组不能变成 array；队列不能变成 queue。\n8. 公司要求里的英文关键词只能用于匹配，不得写进简历，除非该英文词已经原样出现在 parsed_resume_to_optimize 中。\n9. 如果无法在不翻译的情况下优化某句话，就保留原句。\n10. 保持联系方式、教育、证书、语言、日期、公司名、项目名和链接等事实字段不变。"
	userPrompt := fmt.Sprintf(`请根据 target_company_requirements 优化 parsed_resume_to_optimize。

两份输入的角色：
1. parsed_resume_to_optimize：这是用户原本的结构化简历，也是唯一可以被改写的简历内容。
2. target_company_requirements：这是目标公司的岗位要求，只能作为优化依据，不能直接合并到简历里。

优化目标：
1. 让 parsed_resume_to_optimize 更匹配 target_company_requirements。
2. 优化 summary、skills、work_experiences、project_experiences 的表达和排序。
3. 优先突出原简历中已经存在、且和目标岗位匹配的技术、项目和成果。
4. 所有内容必须来自 parsed_resume_to_optimize，不得新增事实。

语言保持规则，非常重要：
1. 不要翻译用户原简历里的任何内容。
2. 原简历字段是中文，优化后仍然必须是中文。
3. 原简历字段是英文，优化后仍然必须是英文。
4. 原简历字段是中英混合，优化后保持同样的中英混合风格。
5. 原简历里已有的技术名、缩写、中文概念必须原样保留。
6. 如果原简历写“链表”，输出仍然写“链表”，不能写“linked list”。
7. 如果原简历写“SpringBoot”，可以保留“SpringBoot”；但不能把原简历里的中文表述整体改成英文句子。
8. target_company_requirements 里的英文词只用于理解岗位偏好，不能导致原简历中文被翻译成英文。

输出要求：
1. 输出一个合法的简历 JSON 对象。
2. 保持原 JSON 字段名。
3. 不要输出 target_company_requirements。
4. 不要输出优化说明。
5. 如果某个字段没有可靠优化空间，直接保留 parsed_resume_to_optimize 中的原值。

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
