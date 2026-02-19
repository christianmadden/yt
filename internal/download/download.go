package download

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// downloadFilter passes only lines starting with "[download]" to the
// underlying writer. It handles both \r (progress bar) and \n terminators.
type downloadFilter struct {
	buf []byte
	out io.Writer
}

func (f *downloadFilter) Write(p []byte) (n int, err error) {
	f.buf = append(f.buf, p...)
	for {
		i := bytes.IndexAny(f.buf, "\r\n")
		if i < 0 {
			break
		}
		line := f.buf[:i]
		term := f.buf[i]
		f.buf = f.buf[i+1:]
		if bytes.HasPrefix(line, []byte("[download]")) {
			f.out.Write(line)
			f.out.Write([]byte{term})
		}
	}
	return len(p), nil
}

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

	args = append(args, opts.URL)

	cmd := exec.Command("yt-dlp", args...)
	if opts.Verbose {
		cmd.Stdout = os.Stdout
	} else {
		cmd.Stdout = &downloadFilter{out: os.Stdout}
	}
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
