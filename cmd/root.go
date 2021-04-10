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
	verbose     bool
	summary     bool
)

func init() {
	// Random seed, we are using for generating random temp directories
	rand.Seed(time.Now().UnixNano())

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output during the download")
	rootCmd.PersistentFlags().BoolVarP(&summary, "summary", "s", false, "Summary after the download")
	rootCmd.PersistentFlags().IntVarP(&connections, "connections", "c", 50, "How many connections to fire to a single file?")
}

var rootCmd = &cobra.Command{
	Use:   "gotake",
	Short: "Fast HTTP file downloads",
	Long:  longDescription,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			printMessage("Please provide a valid URL. Use GoTake like: gotake http://demo.com/test.png OR gotake -h", "warning")
			return
		}

		url = args[0]

		downloadRanges(url)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
