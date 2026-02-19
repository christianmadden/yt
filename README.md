# yt

A friendly [yt-dlp](https://github.com/yt-dlp/yt-dlp) wrapper — download videos without memorizing flags.

## Features

- Conversational prompts: quality selection, filename, save location
- Smart defaults via `~/.config/yt/config.yaml` (auto-created on first run)
- Quality aliases: `4k`, `hd`, `1080p`, `720p`, etc.
- Live yt-dlp progress bar
- `-y` flag for fully non-interactive scripting

## Requirements

- [yt-dlp](https://github.com/yt-dlp/yt-dlp): `brew install yt-dlp`
- [ffmpeg](https://ffmpeg.org/) (for merging video+audio): `brew install ffmpeg`

## Install

```bash
brew install go          # one-time
cd ~/dev/yt
make install             # builds and copies to /usr/local/bin/yt
```

## Usage

```
yt <url> [flags]

Flags:
  -q, --quality string    Target quality height (360/480/720/1080/1440/2160)
  -o, --output string     Output directory
  -n, --name string       Output filename (extension auto-added)
  -f, --format string     Container format: mp4, mkv, webm
  -y, --yes               Non-interactive: accept all defaults, no prompts
      --audio-only        Download best audio only (m4a)
      --list-formats      Print available formats and exit
      --config string     Override config file path
  -v, --verbose           Pass verbose flag through to yt-dlp
```

## Examples

```bash
# Interactive (recommended for everyday use)
yt "https://www.youtube.com/watch?v=dQw4w9WgXcQ"

# Non-interactive download at 720p
yt "https://..." -y -q 720 -n "my-video" -o ~/Videos

# Audio only
yt "https://..." --audio-only

# See available resolutions
yt "https://..." --list-formats

# Use 4K alias
yt "https://..." -q 4k
```

## Config

`~/.config/yt/config.yaml` is created automatically on first run:

```yaml
preferred_quality: 1080   # target height in pixels
min_quality: 720          # warn if best available is below this
default_format: mp4       # output container
default_dir: ""           # empty = current working directory
```

## Module path

Module path: `github.com/christianmadden/yt`
