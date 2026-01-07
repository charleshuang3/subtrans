# subtrans

A subtitle translation tool powered by LLM. Translate subtitle files (SRT, VTT, ASS, etc.) using OpenAI-compatible APIs.

## Features

- Support for multiple subtitle formats (SRT, VTT, ASS, SSA, etc.)
- Batch translation with configurable batch size (50 items per batch)
- Progress logging for long-running translations
- Configurable API endpoint and model
- Configuration file support with sensible defaults

## Installation

```bash
go install github.com/charleshuang3/subtrans@latest
```

Or build from source:

```bash
git clone https://github.com/charleshuang3/subtrans.git
cd subtrans
go build -o subtrans .
```

## Configuration

Create a configuration file at one of the following locations (in order of priority):

1. Path specified via `-c` flag
2. `.env.yaml` in the current directory
3. `~/.config/subtrans/config.yaml`

Configuration file format:

```yaml
default_llm: "openai"  # name of the default LLM provider
target_lang: "Chinese"
llms:  # map of LLM provider configurations
  openai:
    api: "openai"  # or "gemini"
    api_key: "your-openai-api-key"  # required
    api_url: "https://api.openai.com/v1"  # optional, defaults to API provider's URL
    model: "gpt-4o"  # required
    max_tokens: 128000  # optional, defaults to 128000
    structure_output: "json_schema"  # optional for OpenAI, "json_object" or "json_schema"
  gemini:
    api: "gemini"
    api_key: "your-gemini-api-key"  # required
    model: "gemini-1.5-pro"  # required
    max_tokens: 128000  # optional, defaults to 128000
prompts:  # optional, custom prompts for different translation contexts
  default: "Translate the following subtitle text to {target_lang}, preserving timing and formatting:"
  formal: "Translate the following subtitle text to {target_lang} using formal language:"
  casual: "Translate the following subtitle text to {target_lang} using casual, conversational language:"
```

## Usage

Basic usage:

```bash
subtrans -i input.srt -o output.srt
```

Specify target language:

```bash
subtrans -i input.srt -o output.srt -target-lang "Japanese"
```

Use custom config file:

```bash
subtrans -i input.srt -o output.srt -c /path/to/config.yaml
```

Use specific LLM provider:

```bash
subtrans -i input.srt -o output.srt -llm "gemini"
```

Resume translation from specific index:

```bash
subtrans -i input.srt -o output.srt -from "0,5,2"
```

Use custom prompt:

```bash
subtrans -i input.srt -o output.srt -prompt "formal"
```

Dry run (no API calls, returns empty translations):

```bash
subtrans -i input.srt -o output.srt --dry-run
```

### Flags

| Flag | Description |
|------|-------------|
| `-i` | Input file path (required) |
| `-o` | Output file path (required) |
| `-target-lang` | Target language (optional, overrides config) |
| `-c` | Config file path (optional) |
| `-prompt` | Prompt key from config (optional, defaults to "default") |
| `-llm` | LLM provider to use (optional, defaults to "default") |
| `-from` | Resume from index (item,line,seg) (optional) |
| `--dry-run` | Dry run without making API calls (optional) |

## Tests

There are some integration tests in this project requires keys for LLM (Gemini, Deepseek and XAI).
You need to copy `.env.example` to `.env`, and fill in your keys. Then run `just alltest`.

## License

Apache-2.0 License
