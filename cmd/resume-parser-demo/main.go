package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"aresumei/internal/agent/resume_parser"
	"aresumei/internal/infrastructure/llm"
)

const (
	dashScopeAPIKeyEnv = "DASHSCOPE_API_KEY"
	dashScopeModel     = "qwen-plus"
	dashScopeBaseURL   = "https://dashscope.aliyuncs.com/compatible-mode/v1"
)

func main() {
	ctx := context.Background()

	chatModel, err := llm.NewOpenAIChatModel(ctx, llm.OpenAIConfig{
		APIKey:  os.Getenv(dashScopeAPIKeyEnv),
		Model:   dashScopeModel,
		BaseURL: dashScopeBaseURL,
	})
	if err != nil {
		log.Fatalf("new openai chat model: %v", err)
	}

	agent, err := resume_parser.NewAgent(resume_parser.Config{
		ChatModel: chatModel,
	})
	if err != nil {
		log.Fatalf("new resume parser agent: %v", err)
	}

	result, err := agent.Parse(ctx, resume_parser.ParseInput{
		RawText: sampleResume(),
	})
	if err != nil {
		log.Fatalf("parse resume: %v", err)
	}

	fmt.Println(result.RawJSON)
}

func sampleResume() string {
	return `张三
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
英语 CET-6`
}
