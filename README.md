# Laconic & LLMHub CLI Tools

这个项目包含了四个独立的命令行工具，展示了如何使用 [llmhub](https://github.com/smhanov/llmhub) 和 [laconic](https://github.com/smhanov/laconic) 构建灵活且功能强大的 LLM 应用。

这四个工具分别是：
- `text-cli`: 标准文本生成。
- `stream-cli`: 流式文本生成。
- `vision-cli`: 多模态视觉图像理解。
- `search-cli`: 基于 Agent 的自动联网搜索（使用 DuckDuckGo 和 Laconic）。

## 环境配置

所有工具都共享一套一致的配置体系。你可以通过 **命令行参数** 或 **环境变量** 来配置：

- `--key` 或 `OPENAI_API_KEY`: 你的 API Key。
- `--base-url` 或 `OPENAI_BASE_URL`: (可选) 自定义模型服务的 Base URL (例如本地的 Ollama, vLLM 等)。
- `--model` 或 `OPENAI_MODEL`: (可选) 模型名称 (例如 `gpt-4o`, `qwen2.5` 等)。

## 编译与运行方式

在项目根目录下编译所有工具：

```bash
go build -o text-cli ./cmd/text
go build -o stream-cli ./cmd/stream
go build -o vision-cli ./cmd/vision
go build -o search-cli ./cmd/search
```

### 1. Text CLI (标准文本生成)

用于发送一段 Prompt，并在大模型生成完毕后一次性输出。

```bash
# 使用环境变量
export OPENAI_API_KEY="sk-..."
./text-cli "Hello, what is Go?"

# 使用命令行参数
./text-cli --key "sk-..." --model "gpt-4o" "请解释一下什么是协程？"
```

### 2. Stream CLI (流式文本生成)

用于流式输出大模型的结果，提供打字机体验，非常适合终端交互。

```bash
export OPENAI_API_KEY="sk-..."
./stream-cli "写一首关于秋天的现代诗。"
```

### 3. Vision CLI (视觉图像分析)

用于处理图像（支持本地文件路径、网址或 Base64）。你需要使用 `--image` 参数指定图像。

```bash
export OPENAI_API_KEY="sk-..."
# 解析本地图片
./vision-cli --image "cat.jpg" "这张图片里有什么？"

# 解析网络图片
./vision-cli --image "https://example.com/logo.png" "描述这个 Logo 的颜色"
```

### 4. Search CLI (Agentic 联网搜索)

基于 `laconic` 的 Agent 框架实现。当模型不确定答案时，会自动使用 DuckDuckGo 搜索网络，总结知识，最终回答你的问题。由于会进行多次思考与搜索，运行可能需要几秒到十几秒的时间。

```bash
export OPENAI_API_KEY="sk-..."
# 问一个实时性强的问题
./search-cli "今天苹果公司的股票价格是多少？"
```

---

## 架构设计与开发思路

### 1. 架构概览

本项目采用 **单体仓库、多入口 (Monorepo with Multiple Commands)** 的模式。
每一个具体的 CLI 都作为一个独立的包存放在 `cmd/` 目录下（如 `cmd/text`, `cmd/stream` 等）。
共享逻辑存放在 `internal/` 目录下。

```text
├── cmd/
│   ├── search/     # 联网搜索 Agent
│   ├── stream/     # 流式输出工具
│   ├── text/       # 标准文本输出
│   └── vision/     # 图像视觉分析
├── internal/
│   └── config/     # 共享的命令行参数解析及环境配置逻辑
├── go.mod
└── README.md
```

### 2. 开发思路

*   **配置复用 (DRY 原则):**
    四个命令行工具在初始化大模型时，都需要 API Key, Base URL 和 Model 等参数。为了避免重复代码，在 `internal/config/config.go` 中封装了统一的 `ParseFlags()` 函数。该函数首先解析命令行 Flag，如果 Flag 为空，则自动回退(Fallback)去读取环境变量。

*   **无缝对接任何提供商 (LLMHub 的核心优势):**
    所有工具均依赖 `llmhub` 进行初始化。由于 `llmhub` 抹平了提供商（如 OpenAI, Anthropic, Gemini, 本地 Ollama）的差异，代码中统一使用 `llmhub.New("openai", key, opts...)`。
    如果用户想对接本地的 Ollama 服务，完全**不需要修改代码**，只需要传入 `--base-url http://localhost:11434/v1` 和对应的 `--model` 即可。

*   **多模态适配 (Vision CLI):**
    在处理图像时，考虑到用户可能传入 HTTP URL，也可能传入本地文件。我们在代码中实现了一个 `resolveImage` 函数：如果识别到是 URL 或 Base64 就直接放行，如果是本地路径，就读取文件并转换为 Data URI，这样就能无缝通过 `llmhub` 的 `Image()` 接口传给大模型。

*   **框架适配桥接 (Search CLI):**
    `laconic` 核心需要一个实现 `laconic.LLMProvider` 接口的对象，而 `llmhub.Client` 默认的签名与之略有不同。我们在 `cmd/search/adapter.go` 中采用了**适配器模式 (Adapter Pattern)**，将 `llmhub` 的接口适配为 `laconic` 所需的简单文本接口，完美整合了这两个出色的库。
