package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"math/rand"
	"os"
	"time"
)

var (
	url         string
	connections int
	auto        bool
	verbose     bool
	summary     bool
	standard    bool
)

func init() {
	// Random seed, we are using for generating random temp directories
	rand.Seed(time.Now().UnixNano())

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output during the download.")
	rootCmd.PersistentFlags().BoolVarP(&summary, "info", "i", false, "Info (summary) in the process of download.")
	rootCmd.PersistentFlags().BoolVarP(&standard, "standard", "s", false, "Should we use the standard download (much slower)?")
	rootCmd.PersistentFlags().IntVarP(&connections, "connections", "c", 50, "How many connections to fire to a single URL?")
	rootCmd.PersistentFlags().BoolVarP(&auto, "auto", "a", true, "Do you want to calculate automatically the number connections?")
	rootCmd.PersistentFlags().StringVarP(&filename, "filename", "f", "", "Different name to the target file?")
}

var rootCmd = &cobra.Command{
	Use:   "gotake",
	Short: "Fast and reliable file downloads",
	Long:  longDescription,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			printMessage("Please provide a valid URL. Use GoTake like: gotake http://demo.com/test.png OR gotake -h", "warning")
			return
		}

		// Catch the URL from arguments
		url = args[0]

		// Validate the URL
		if validateUrl(url) == false {
			message := fmt.Sprintf("%s doesn't have protocol, trying with http:// prefix", url)
			printMessage(message, "info")

			// Append the default protocol on the address
			url = fmt.Sprintf("%s%s", "http://", url)
		}

		// Chose which method to use for the downloading
		if standard {
			downloadStandard(url)
		} else {
			downloadRanges(url)
		}

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
