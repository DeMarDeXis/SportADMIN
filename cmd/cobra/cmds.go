package cobra

import (
	"github.com/spf13/cobra"
)

var nhlCmd = &cobra.Command{
	Use:   "nhl",
	Short: "NHL",
	Run:   parseNHL,
}

var nflCmd = &cobra.Command{
	Use:   "nfl",
	Short: "NFL",
	Run:   parseNFL,
}

func init() {
	rootCmd.AddCommand(nhlCmd)
	rootCmd.AddCommand(nflCmd)

	nhlCmd.Flags().StringP("method", "m", "", "What are we parsing")

	if err := nhlCmd.MarkFlagRequired("method"); err != nil {
		panic(err)
	}
}
