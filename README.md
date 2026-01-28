# Go-HAR

[![Go Reference](https://pkg.go.dev/badge/github.com/cyberspacesec/go-har.svg)](https://pkg.go.dev/github.com/cyberspacesec/go-har)
[![Go Report Card](https://goreportcard.com/badge/github.com/cyberspacesec/go-har)](https://goreportcard.com/report/github.com/cyberspacesec/go-har)
[![License](https://img.shields.io/github/license/cyberspacesec/go-har)](https://github.com/cyberspacesec/go-har/blob/main/LICENSE)

Go-HAR 是一个高性能、灵活的 HTTP Archive (HAR) 解析和处理库，用 Go 语言实现。它为处理 HAR 文件提供了多种策略，从简单的小型文件到需要优化内存使用的大型文件都能高效处理。

## 特性

- **多种解析策略**
  - 标准解析：适用于常规场景
  - 内存优化：减少大型 HAR 文件的内存占用
  - 懒加载：延迟加载大型内容字段
  - 流式处理：逐条处理超大 HAR 文件

- **灵活的接口设计**
  - 基于接口的设计，支持不同实现之间的互操作
  - 统一的 API，无论使用哪种解析策略

- **增强的错误处理**
  - 详细的错误信息和上下文
  - 部分解析能力和警告收集

- **高级功能**
  - 高效过滤和搜索
  - 丰富的统计分析
  - 可视化和报告生成

## 安装

```bash
go get github.com/cyberspacesec/go-har
```

## 快速开始

### 基本用法

```go
package main

import (
	"fmt"
	"log"
	
	"github.com/cyberspacesec/go-har"
)

func main() {
	// 解析 HAR 文件
	harData, err := har.ParseHarFile("example.har")
	if err != nil {
		log.Fatalf("无法解析 HAR 文件: %v", err)
	}
	
	// 访问 HAR 数据
	fmt.Printf("HAR 版本: %s\n", harData.Log.Version)
	fmt.Printf("条目数量: %d\n", len(harData.Log.Entries))
	
	// 遍历所有请求
	for i, entry := range harData.Log.Entries {
		fmt.Printf("请求 #%d: %s %s\n", i+1, entry.Request.Method, entry.Request.URL)
	}
}
```

### 内存优化模式

```go
// 使用内存优化模式处理大型文件
harData, err := har.ParseHarFile("large.har", har.WithMemoryOptimized())
if err != nil {
	log.Fatalf("无法解析 HAR 文件: %v", err)
}

// 接口保持一致，使用方式相同
for _, entry := range harData.GetEntries() {
	fmt.Printf("URL: %s\n", entry.GetRequest().GetURL())
}
```

### 懒加载模式

```go
// 使用懒加载模式延迟加载大型内容
harData, err := har.ParseHarFile("large_content.har", har.WithLazyLoading())
if err != nil {
	log.Fatalf("无法解析 HAR 文件: %v", err)
}

// 基本信息直接可用，大型内容仅在需要时加载
for _, entry := range harData.GetEntries() {
	resp := entry.GetResponse()
	fmt.Printf("状态码: %d, 内容大小: %d\n", 
		resp.GetStatus(), 
		resp.GetContent().GetSize())
	
	// 内容仅在需要时加载
	if resp.GetStatus() == 200 {
		content := resp.GetContent()
		text := content.GetText() // 此时才加载内容
		fmt.Printf("内容长度: %d\n", len(text))
	}
}
```

## 高级用法

查看 [详细文档](./doc/usage.md) 了解更多高级功能，包括：

- 流式解析超大 HAR 文件
- 增强的错误处理
- 过滤和搜索功能
- 统计分析和可视化
- 命令行工具

## 项目结构

- `pkg/har/` - 核心 HAR 解析和处理代码
- `examples/` - 示例代码和实用工具
  - `examples/statistics/` - 统计分析示例
  - `examples/visualization/` - 可视化示例
  - `examples/cli-tool/` - 命令行工具示例
- `doc/` - 详细文档
  - `doc/usage.md` - 使用文档
  - `doc/structure.md` - HAR 结构文档

## 贡献

欢迎贡献！请查看 [贡献指南](CONTRIBUTING.md) 了解如何参与项目开发。

## 许可证

本项目使用 [MIT 许可证](LICENSE)。

