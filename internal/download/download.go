package download

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type Options struct {
	URL       string
	Quality   int
	OutputDir string
	Name      string
	Format    string
	AudioOnly bool
	Verbose   bool
}

func Run(opts Options) error {
	args := []string{}

	if opts.AudioOnly {
		args = append(args,
			"-f", "bestaudio[ext=m4a]/bestaudio",
			"--extract-audio",
			"--audio-format", "m4a",
		)
	} else {
		fmtStr := fmt.Sprintf(
			"bestvideo[height<=%d][ext=mp4]+bestaudio[ext=m4a]/bestvideo[height<=%d]+bestaudio/best[height<=%d]",
			opts.Quality, opts.Quality, opts.Quality,
		)
		args = append(args,
			"-f", fmtStr,
			"--merge-output-format", opts.Format,
		)
	}

	outDir := opts.OutputDir
	if outDir == "" {
		var err error
		outDir, err = os.Getwd()
		if err != nil {
			outDir = "."
		}
	}

	var outputTemplate string
	if opts.Name != "" {
		outputTemplate = filepath.Join(outDir, opts.Name+".%(ext)s")
	} else {
		outputTemplate = filepath.Join(outDir, "%(title)s.%(ext)s")
	}
	args = append(args, "-o", outputTemplate)

	if opts.Verbose {
		args = append(args, "--verbose")
	}

	args = append(args, opts.URL)

	cmd := exec.Command("yt-dlp", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
