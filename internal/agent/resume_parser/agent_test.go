package resume_parser

import (
	"context"
	"errors"
	"testing"

	"github.com/cloudwego/eino/components/model"
	einoschema "github.com/cloudwego/eino/schema"
)

type fakeChatModel struct {
	response string
}

func (m fakeChatModel) Generate(
	_ context.Context,
	_ []*einoschema.Message,
	_ ...model.Option,
) (*einoschema.Message, error) {
	return &einoschema.Message{
		Role:    einoschema.Assistant,
		Content: m.response,
	}, nil
}

func (m fakeChatModel) Stream(
	_ context.Context,
	_ []*einoschema.Message,
	_ ...model.Option,
) (*einoschema.StreamReader[*einoschema.Message], error) {
	return nil, errors.New("stream is not implemented in fakeChatModel")
}

func TestAgentParse(t *testing.T) {
	agent, err := NewAgent(Config{
		ChatModel: fakeChatModel{
			response: `{
  "basics": {
    "full_name": "张三",
    "email": "zhangsan@example.com",
    "phone": "13800000000",
    "location": "上海",
    "github": "https://github.com/zhangsan",
    "linkedin": "",
    "website": ""
  },
  "summary": "5 年 Golang 后端开发经验，熟悉微服务、MySQL、Redis 和 Kubernetes。",
  "skills": [
    {
      "category": "Backend",
      "skills": ["Go", "Gin", "gRPC", "MySQL", "Redis", "Kafka", "Kubernetes"]
    }
  ],
  "work_experiences": [
    {
      "company": "某科技公司",
      "position": "Golang 后端工程师",
      "location": "上海",
      "start_date": "2021-03",
      "end_date": "",
      "is_current": true,
      "description": "负责核心交易系统后端开发。",
      "highlights": [
        "优化订单查询链路，将接口 P95 延迟从 800ms 降低到 220ms",
        "设计 Redis 缓存方案，降低数据库读压力"
      ],
      "technologies": ["Go", "MySQL", "Redis", "Kafka"]
    }
  ],
  "project_experiences": [
    {
      "name": "订单履约平台",
      "role": "后端负责人",
      "start_date": "2022-01",
      "end_date": "2023-06",
      "description": "建设订单履约链路，支持订单状态流转和异步消息处理。",
      "highlights": [
        "使用 Kafka 解耦订单状态变更事件",
        "引入 gRPC 拆分核心服务接口"
      ],
      "technologies": ["Go", "gRPC", "Kafka", "MySQL"],
      "links": []
    }
  ],
  "education": [
    {
      "school": "某大学",
      "degree": "本科",
      "major": "计算机科学与技术",
      "start_date": "2016-09",
      "end_date": "2020-06",
      "location": "杭州"
    }
  ],
  "certifications": [],
  "languages": [
    {
      "name": "英语",
      "proficiency": "CET-6"
    }
  ]
}`,
		},
	})
	if err != nil {
		t.Fatalf("NewAgent() error = %v", err)
	}

	result, err := agent.Parse(context.Background(), ParseInput{
		RawText: `张三
邮箱：zhangsan@example.com
电话：13800000000
地点：上海
GitHub：https://github.com/zhangsan

5 年 Golang 后端开发经验，熟悉微服务、MySQL、Redis 和 Kubernetes。

工作经历：
某科技公司 Golang 后端工程师 2021.03 - 至今
- 负责核心交易系统后端开发
- 优化订单查询链路，将接口 P95 延迟从 800ms 降低到 220ms
- 设计 Redis 缓存方案，降低数据库读压力

项目经历：
订单履约平台，后端负责人，2022.01 - 2023.06
- 建设订单履约链路，支持订单状态流转和异步消息处理
- 使用 Kafka 解耦订单状态变更事件
- 引入 gRPC 拆分核心服务接口

教育经历：
某大学 计算机科学与技术 本科 2016.09 - 2020.06

语言：
英语 CET-6`,
	})
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if result.Resume.Basics.FullName != "张三" {
		t.Fatalf("full_name = %q, want %q", result.Resume.Basics.FullName, "张三")
	}

	if len(result.Resume.WorkExperiences) != 1 {
		t.Fatalf("work_experiences len = %d, want 1", len(result.Resume.WorkExperiences))
	}

	if result.Resume.WorkExperiences[0].Company != "某科技公司" {
		t.Fatalf("company = %q, want %q", result.Resume.WorkExperiences[0].Company, "某科技公司")
	}

	if result.Resume.WorkExperiences[0].Position != "Golang 后端工程师" {
		t.Fatalf("position = %q, want %q", result.Resume.WorkExperiences[0].Position, "Golang 后端工程师")
	}

	if len(result.Resume.ProjectExperiences) != 1 {
		t.Fatalf("project_experiences len = %d, want 1", len(result.Resume.ProjectExperiences))
	}

	if result.RawJSON == "" {
		t.Fatal("RawJSON is empty")
	}
}
