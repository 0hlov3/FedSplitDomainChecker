package cmd

import (
	"github.com/0hlov3/FedSplitDomainChecker/internal/logger"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "FedSplitDomainChecker",
	Short: "Fediverse split-domain deployment checker",
	Long: `A CLI tool to check and validate Fediverse split-domain deployments,
including account and host domains.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		logger.Init(Verbose)
	},
}

var (
	Verbose       bool
	Debug         bool
	AccountDomain string
	HostDomain    string
	Account       string
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&AccountDomain, "accountDomain", "gotosocial.org", "The account domain to check (default: gotosocial.org)")
	rootCmd.PersistentFlags().StringVar(&HostDomain, "hostDomain", "gts.gotosocial.org", "The host domain to check (default: gts.gotosocial.org)")
	rootCmd.PersistentFlags().StringVar(&Account, "account", "admin@gotosocial.org", "The account to check (default: admin@gotosocial.org)")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Enable verbose logging")
	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Enable debug logging")

	viper.BindPFlag("accountDomain", rootCmd.PersistentFlags().Lookup("accountDomain"))
	viper.BindPFlag("hostDomain", rootCmd.PersistentFlags().Lookup("hostDomain"))
	viper.BindPFlag("account", rootCmd.PersistentFlags().Lookup("account"))
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}
