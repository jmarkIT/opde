/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"opde/logic"

	"github.com/spf13/cobra"
)

// vaultCmd represents the vault command
var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 1 {
			fmt.Println("Too many arguments")
		} else if len(args) < 1 {
			fmt.Println("Please provide a vault to list")
		} else {
			vault := args[0]
			managersFlag, _ := cmd.Flags().GetBool("managers")
			account, _ := cmd.Flags().GetString("account")
			csv, _ := cmd.Flags().GetBool("csv")

			groups := logic.GetVaultGroups(vault, account)
			for _, group := range groups {
				group.SetMembers(account)
				managers := group.GetManagers(account)
				if managersFlag {
					logic.PrintOutput(group, managers, csv)
				} else {
					logic.PrintOutput(group, managers, csv)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(vaultCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// vaultCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// vaultCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	vaultCmd.Flags().BoolP("managers", "m", false, "Only print managers.")
}
