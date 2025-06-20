package cobra

func init() {
	//ROOT CMD
	rootCmd.AddCommand(nhlCmd)
	rootCmd.AddCommand(nflCmd)
	rootCmd.AddCommand(nbaCmd)
	rootCmd.AddCommand(mlbCmd)

	//NHL CMD
	nhlCmd.AddCommand(nhlParseCmd)
	nhlCmd.AddCommand(nhlLoadToDBCmd)

	//NBA CMD
	nbaCmd.AddCommand(nbaParseCmd)
	nbaCmd.AddCommand(nbaLoadToDBCmd)

	//NFL CMD
	nflCmd.AddCommand(nflParseCmd)
	nflCmd.AddCommand(nflLoadToDBCmd)

	//MLB CMD
	mlbCmd.AddCommand(mlbParseCmd)
	mlbCmd.AddCommand(mlbLoadToDBCmd)

	//FLAGS for cmds
	//NHL
	nhlParseCmd.Flags().StringP("method", "m", "", "What are we parsing")
	if err := nhlParseCmd.MarkFlagRequired("method"); err != nil {
		panic(err)
	}
	nhlLoadToDBCmd.Flags().StringP("method", "m", "", "What are we loading to DB")
	if err := nhlLoadToDBCmd.MarkFlagRequired("method"); err != nil {
		panic(err)
	}

	//NBA
	nbaParseCmd.Flags().StringP("method", "m", "", "What are we parsing")
	if err := nbaParseCmd.MarkFlagRequired("method"); err != nil {
		panic(err)
	}
	nbaLoadToDBCmd.Flags().StringP("method", "m", "", "What are we loading to DB")
	if err := nbaLoadToDBCmd.MarkFlagRequired("method"); err != nil {
		panic(err)
	}

	//NFL
	nflParseCmd.Flags().StringP("method", "m", "", "What are we parsing")
	if err := nflParseCmd.MarkFlagRequired("method"); err != nil {
		panic(err)
	}
	nflLoadToDBCmd.Flags().StringP("method", "m", "", "What are we loading to DB")
	if err := nflLoadToDBCmd.MarkFlagRequired("method"); err != nil {
		panic(err)
	}

	//MLB
	mlbParseCmd.Flags().StringP("method", "m", "", "What are we parsing")
	if err := mlbParseCmd.MarkFlagRequired("method"); err != nil {
		panic(err)
	}
	mlbLoadToDBCmd.Flags().StringP("method", "m", "", "What are we loading to DB")
	if err := mlbLoadToDBCmd.MarkFlagRequired("method"); err != nil {
		panic(err)
	}

}
