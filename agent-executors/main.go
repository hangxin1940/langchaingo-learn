package main

import (
	"context"
	"fmt"
	"github.com/tmc/langchaingo/agents"
	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/ollama"
	"github.com/tmc/langchaingo/tools"
	"github.com/tmc/langchaingo/tools/serpapi"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	llm, err := ollama.New(
		// 配置ollama地址
		ollama.WithServerURL("http://127.0.0.1:1234"),
		// 指定模型
		ollama.WithModel("llama3:8b"),
	)
	if err != nil {
		return err
	}

	// serpapi工具
	search, err := serpapi.New()
	if err != nil {
		return err
	}

	//  设置工具
	agentTools := []tools.Tool{
		// 计算器工具
		tools.Calculator{},
		search,
	}

	// 创建代理
	agent := agents.NewOneShotAgent(llm,
		agentTools,
		// 最大迭代3次
		agents.WithMaxIterations(3),
		agents.WithReturnIntermediateSteps(),
	)
	executor := agents.NewExecutor(agent)

	question := "Who is Olivia Wilde's boyfriend? What is his current age raised to the 0.23 power?"
	answer, err := chains.Run(
		context.Background(),
		executor,
		question,
		chains.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Print(string(chunk))
			return nil
		}),
	)
	fmt.Println("\n-------------\n")
	fmt.Println(answer)
	return err
}
