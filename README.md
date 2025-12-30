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
api_key: "your-api-key"
api_url: "https://api.openai.com/v1"
model: "gpt-4o"
target_lang: "Chinese"
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

### Flags

| Flag | Description |
|------|-------------|
| `-i` | Input file path (required) |
| `-o` | Output file path (required) |
| `-target-lang` | Target language (optional, overrides config) |
| `-c` | Config file path (optional) |

## License

Apache-2.0 License
