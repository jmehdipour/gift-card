package cmd

import (
	"github.com/jmehdipour/gift-card/internal/interface/http"
	"github.com/spf13/cobra"
)

// startCMD represents the start command of the application.
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "serve API",
	Run:   startFunc,
}

func startFunc(_ *cobra.Command, _ []string) {
	http.NewServer().Serve()
}
