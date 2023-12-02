package cli

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	currency string
)

var addCmd = &cobra.Command{
	Use:   "add [type]",
	Short: "add a new resource",
}

var addCompanyCmd = &cobra.Command{
	Use:   "company [name] [price]",
	Short: "add a company",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if getCompanyIfExists(args[0]) == nil {
			log.Fatalf("company %s already exists", args[0])
		}
		name := args[0]
		price, err := strconv.ParseInt(args[1], 10, 32)
		if err != nil || price <= 0 {
			log.Fatalf("invalid price %s", args[1])
		}
		timeTracker.AddCompany(name, int32(price))
	},
}

func init() {
	defaultCurrency := viper.GetString("settings.default_currency")
	if defaultCurrency == "" {
		defaultCurrency = "USD"
	}
	rootCmd.AddCommand(addCmd)
	addCmd.AddCommand(addCompanyCmd)
	addCompanyCmd.Flags().StringVar(&currency, "currency", defaultCurrency,
		"currency in which the company pays (defaults to USD)")
}
