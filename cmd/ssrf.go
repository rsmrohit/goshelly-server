/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ssrfCmd represents the ssrf command
var ssrfCmd = &cobra.Command{
	Use:   "ssrf",
	Short: "Simulate SSRF attack.",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ssrf called")
	},
}

func init() {
	rootCmd.AddCommand(ssrfCmd)
}
