package fetch

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sort"
)

type VideoInfo struct {
	Title   string
	Heights []int
}

type ytdlpInfo struct {
	Title   string   `json:"title"`
	Formats []format `json:"formats"`
}

type format struct {
	Height int    `json:"height"`
	VCodec string `json:"vcodec"`
	ACodec string `json:"acodec"`
}

func Probe(url string) (*VideoInfo, error) {
	cmd := exec.Command("yt-dlp", "-J", "--no-warnings", url)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("yt-dlp error: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("failed to run yt-dlp: %w", err)
	}

	var info ytdlpInfo
	if err := json.Unmarshal(out, &info); err != nil {
		return nil, fmt.Errorf("failed to parse yt-dlp output: %w", err)
	}

	heightSet := make(map[int]bool)
	for _, f := range info.Formats {
		if f.Height > 0 && f.VCodec != "none" && f.VCodec != "" {
			heightSet[f.Height] = true
		}
	}

	heights := make([]int, 0, len(heightSet))
	for h := range heightSet {
		heights = append(heights, h)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(heights)))

	return &VideoInfo{
		Title:   info.Title,
		Heights: heights,
	}, nil
}
