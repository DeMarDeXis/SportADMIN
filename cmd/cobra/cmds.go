package cobra

import (
	"github.com/spf13/cobra"
)

var nhlCmd = &cobra.Command{
	Use:   "nhl",
	Short: "NHL",
	Long:  "Main NHL command which contains subcommands",
	//Run:   parseNHL, //TODO: check and delete(16.01.25)
}

var nhlParseCmd = &cobra.Command{
	Use:   "nhl-prs",
	Short: "NHL parser",
	Long:  "It is NHL parser",
	Run:   parseNHL,
}

var nhlLoadToDBCmd = &cobra.Command{
	Use:   "nhl-db",
	Short: "NHL loader",
	Long:  "It is NHL loader to DB",
	Run:   loadNHLToDB, //TODO: make loader to DB(16.01.25)
}

var nflCmd = &cobra.Command{
	Use:   "nfl",
	Short: "NFL",
	Run:   parseNFL,
}

func init() {
	//ROOT CMD
	rootCmd.AddCommand(nhlCmd)
	rootCmd.AddCommand(nflCmd)

	//NHL CMD
	nhlCmd.AddCommand(nhlParseCmd)
	nhlCmd.AddCommand(nhlLoadToDBCmd)

	nhlParseCmd.Flags().StringP("method", "m", "", "What are we parsing")

	if err := nhlParseCmd.MarkFlagRequired("method"); err != nil {
		panic(err)
	}

	nhlLoadToDBCmd.Flags().StringP("method", "m", "", "What are we loading to DB")
	if err := nhlLoadToDBCmd.MarkFlagRequired("method"); err != nil {
		panic(err)
	}
}
