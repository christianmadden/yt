package ui

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var (
	green  = color.New(color.FgGreen).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	bold   = color.New(color.Bold).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
)

var reader = bufio.NewReader(os.Stdin)

// Ask prompts the user with a default value. Returns defaultVal on bare Enter.
func Ask(prompt, defaultVal string) string {
	if defaultVal != "" {
		fmt.Printf("%s [%s]: ", prompt, cyan(defaultVal))
	} else {
		fmt.Printf("%s: ", prompt)
	}
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "" {
		return defaultVal
	}
	return input
}

// NormalizeQuality converts quality aliases to integer pixel height.
func NormalizeQuality(input string) (int, error) {
	s := strings.TrimSpace(input)
	lower := strings.ToLower(s)

	switch lower {
	case "4k", "uhd":
		return 2160, nil
	case "hd", "1080p", "1080":
		return 1080, nil
	}

	// Strip trailing 'p' (e.g. "720p" → "720")
	s = strings.TrimSuffix(s, "p")
	s = strings.TrimSuffix(s, "P")

	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("unrecognized quality %q — try 360, 480, 720, 1080, 1440, or 2160", input)
	}
	return n, nil
}

// SanitizeFilename removes characters unsafe for filenames.
func SanitizeFilename(title string) string {
	re := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1f]`)
	s := re.ReplaceAllString(title, "_")
	s = strings.TrimSpace(s)
	if len(s) > 200 {
		s = s[:200]
	}
	return s
}

// PrintFormats prints available heights in a friendly list.
func PrintFormats(heights []int) {
	fmt.Println(bold("Available video heights:"))
	for _, h := range heights {
		label := heightLabel(h)
		fmt.Printf("  %s %s\n", green(fmt.Sprintf("%dp", h)), label)
	}
}

func heightLabel(h int) string {
	switch {
	case h >= 2160:
		return yellow("(4K / UHD)")
	case h >= 1440:
		return "(2K / QHD)"
	case h >= 1080:
		return "(Full HD)"
	case h >= 720:
		return "(HD)"
	case h >= 480:
		return "(SD)"
	default:
		return red("(low quality)")
	}
}

// QualityDialogue runs the conversational quality selection.
// Returns the chosen quality height.
func QualityDialogue(preferred, minimum int, available []int) (int, error) {
	if len(available) == 0 {
		return 0, fmt.Errorf("no video streams found")
	}

	best := available[0] // already sorted descending

	// Check if preferred is available
	preferredAvailable := false
	for _, h := range available {
		if h == preferred {
			preferredAvailable = true
			break
		}
	}

	// Find closest at-or-below preferred
	closestBelow := 0
	for _, h := range available {
		if h <= preferred {
			closestBelow = h
			break
		}
	}

	for {
		if preferredAvailable {
			if best > preferred {
				fmt.Printf("\n%s Your preferred %s is available. %s is also an option if you want it.\n",
					green("✓"),
					bold(fmt.Sprintf("%dp", preferred)),
					yellow(fmt.Sprintf("%dp", best)),
				)
				fmt.Printf("  Press Enter for %s, or type a height (or %s for all options): ",
					cyan(fmt.Sprintf("%dp", preferred)),
					cyan("options"),
				)
			} else {
				fmt.Printf("\n%s Your preferred %s is available. Grabbing that.\n",
					green("✓"),
					bold(fmt.Sprintf("%dp", preferred)),
				)
				return preferred, nil
			}
		} else if closestBelow >= minimum {
			fmt.Printf("\n%s Best available is %s (your preferred is %dp). Still good?\n",
				yellow("!"),
				bold(fmt.Sprintf("%dp", closestBelow)),
				preferred,
			)
			fmt.Printf("  Press Enter for %s, or type a height (or %s for all options): ",
				cyan(fmt.Sprintf("%dp", closestBelow)),
				cyan("options"),
			)
		} else if closestBelow > 0 {
			fmt.Printf("\n%s Heads up: only %s available — that's below your minimum of %dp. Download anyway?\n",
				red("⚠"),
				bold(fmt.Sprintf("%dp", closestBelow)),
				minimum,
			)
			fmt.Printf("  Press Enter for %s, type a height, %s or %s to cancel (or %s for all options): ",
				cyan(fmt.Sprintf("%dp", closestBelow)),
				cyan("yes"),
				cyan("no"),
				cyan("options"),
			)
		} else {
			fmt.Printf("\n%s No suitable video stream found. Available: ", red("⚠"))
			for i, h := range available {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%dp", h)
			}
			fmt.Println()
			fmt.Printf("  Type a height to use, or %s to cancel: ", cyan("no"))
		}

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		lower := strings.ToLower(input)

		if lower == "no" || lower == "n" || lower == "cancel" {
			return 0, fmt.Errorf("download cancelled")
		}

		if lower == "" || lower == "yes" || lower == "y" {
			if preferredAvailable {
				return preferred, nil
			}
			if closestBelow > 0 {
				return closestBelow, nil
			}
		}

		if lower == "more" || lower == "options" || lower == "list" {
			fmt.Println()
			PrintFormats(available)
			fmt.Println()
			continue
		}

		if input != "" {
			q, err := NormalizeQuality(input)
			if err != nil {
				fmt.Printf("  %s\n", red(err.Error()))
				continue
			}
			// Find best available at or below requested
			for _, h := range available {
				if h <= q {
					fmt.Printf("  %s Using %dp\n", green("✓"), h)
					return h, nil
				}
			}
			fmt.Printf("  %s No stream available at or below %dp. Try: ", red("✗"), q)
			for i, h := range available {
				if i > 0 {
					fmt.Print(", ")
				}
				fmt.Printf("%dp", h)
			}
			fmt.Println()
		}
	}
}

// FilenameDialogue asks user to confirm or override the suggested filename.
func FilenameDialogue(suggested string) string {
	sanitized := SanitizeFilename(suggested)
	fmt.Printf("\n%s Save as: ", bold("Filename"))
	result := Ask("", sanitized)
	return SanitizeFilename(result)
}

// DirectoryDialogue asks user to confirm or override the save directory.
func DirectoryDialogue(defaultDir string) string {
	if defaultDir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			cwd = "."
		}
		defaultDir = cwd
	}
	fmt.Printf("\n%s Save to: ", bold("Directory"))
	return Ask("", defaultDir)
}

// PrintTitle prints the video title in a friendly way.
func PrintTitle(title string) {
	fmt.Printf("\n%s %s\n", bold("Video:"), cyan(title))
}

// Info prints a plain informational message.
func Info(msg string) {
	fmt.Println(msg)
}

// Success prints a green success message.
func Success(msg string) {
	fmt.Println(green("✓ " + msg))
}

// Warn prints a yellow warning message.
func Warn(msg string) {
	fmt.Println(yellow("! " + msg))
}

// Error prints a red error message.
func Error(msg string) {
	fmt.Println(red("✗ " + msg))
}
