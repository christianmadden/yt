package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/christianmadden/yt/internal/config"
	"github.com/christianmadden/yt/internal/download"
	"github.com/christianmadden/yt/internal/fetch"
	"github.com/christianmadden/yt/internal/ui"
)

var (
	flagQuality     string
	flagOutput      string
	flagName        string
	flagFormat      string
	flagYes         bool
	flagAudioOnly   bool
	flagListFormats bool
	flagConfig      string
	flagVerbose     bool
)

func main() {
	root := &cobra.Command{
		Use:   "yt <url>",
		Short: "A friendly yt-dlp wrapper",
		Long: `yt — download videos without memorizing flags.

Wraps yt-dlp with smart defaults, conversational prompts, and a
clean config file. Pass -y to skip all prompts for scripting.`,
		Args:          cobra.ExactArgs(1),
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE:          run,
	}

	root.Flags().StringVarP(&flagQuality, "quality", "q", "", "Target quality height (360/480/720/1080/1440/2160)")
	root.Flags().StringVarP(&flagOutput, "output", "o", "", "Output directory")
	root.Flags().StringVarP(&flagName, "name", "n", "", "Output filename (extension auto-added)")
	root.Flags().StringVarP(&flagFormat, "format", "f", "", "Container format: mp4, mkv, webm")
	root.Flags().BoolVarP(&flagYes, "yes", "y", false, "Non-interactive: accept all defaults, no prompts")
	root.Flags().BoolVar(&flagAudioOnly, "audio-only", false, "Download best audio only (m4a)")
	root.Flags().BoolVar(&flagListFormats, "list-formats", false, "Print available formats and exit")
	root.Flags().StringVar(&flagConfig, "config", "", "Override config file path")
	root.Flags().BoolVarP(&flagVerbose, "verbose", "v", false, "Pass verbose flag through to yt-dlp")

	if err := root.Execute(); err != nil {
		ui.Error(err.Error())
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	url := args[0]

	cfg, err := config.Load(flagConfig)
	if err != nil {
		return fmt.Errorf("config error: %w", err)
	}

	// --list-formats: probe and print, then exit
	if flagListFormats {
		info, err := fetch.Probe(url)
		if err != nil {
			return err
		}
		ui.PrintTitle(info.Title)
		ui.PrintFormats(info.Heights)
		return nil
	}

	// Resolve format
	outputFormat := cfg.DefaultFormat
	if flagFormat != "" {
		outputFormat = flagFormat
	}

	// Audio-only shortcut — skip quality negotiation
	if flagAudioOnly {
		outputDir := flagOutput
		if outputDir == "" {
			outputDir = cfg.DefaultDir
		}
		name := flagName

		if !flagYes {
			if name == "" {
				info, err := fetch.Probe(url)
				if err != nil {
					return err
				}
				ui.PrintTitle(info.Title)
				name = ui.FilenameDialogue(info.Title)
			}
			outputDir = ui.DirectoryDialogue(outputDir)
		}

		return download.Run(download.Options{
			URL:       url,
			OutputDir: outputDir,
			Name:      name,
			Format:    outputFormat,
			AudioOnly: true,
			Verbose:   flagVerbose,
		})
	}

	// Probe video info
	ui.Info("Fetching video info…")
	info, err := fetch.Probe(url)
	if err != nil {
		return err
	}
	ui.PrintTitle(info.Title)

	// Resolve quality
	var chosenQuality int
	if flagQuality != "" {
		chosenQuality, err = ui.NormalizeQuality(flagQuality)
		if err != nil {
			return err
		}
	} else if flagYes {
		chosenQuality = cfg.PreferredQuality
	} else {
		chosenQuality, err = ui.QualityDialogue(cfg.PreferredQuality, cfg.MinQuality, info.Heights)
		if err != nil {
			return err
		}
	}

	// Resolve filename
	name := flagName
	if name == "" && !flagYes {
		name = ui.FilenameDialogue(info.Title)
	}

	// Resolve output dir
	outputDir := flagOutput
	if outputDir == "" {
		outputDir = cfg.DefaultDir
	}
	if !flagYes && flagOutput == "" {
		outputDir = ui.DirectoryDialogue(outputDir)
	}

	fmt.Println()
	return download.Run(download.Options{
		URL:       url,
		Quality:   chosenQuality,
		OutputDir: outputDir,
		Name:      name,
		Format:    outputFormat,
		AudioOnly: false,
		Verbose:   flagVerbose,
	})
}
