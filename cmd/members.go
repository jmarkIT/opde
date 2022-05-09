/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"opde/logic"

	"github.com/spf13/cobra"
)

// membersCmd represents the members command
var membersCmd = &cobra.Command{
	Use:   "members",
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
				if managersFlag {
					managers := group.GetManagers(account)
					logic.PrintGroupMembers(group, managers, csv)
				} else {
					members := group.GetMembers(account)
					logic.PrintGroupMembers(group, members, csv)
				}
			}
		}
	},
}

func init() {
	groupsCmd.AddCommand(membersCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// membersCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// membersCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	membersCmd.Flags().BoolP("managers", "m", false, "Only print group managers")
}
