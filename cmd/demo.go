/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	s "cobrashelly/basic"
	"github.com/spf13/cobra"

)

// demoCmd represents the demo command
var demoCmd = &cobra.Command{
	Use:   "demo",
	Short: "Run a set of cli commands.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) { 

		fmt.Println("Running GoShelly Server - Demo.")
		PORT, _ := cmd.Flags().GetString("PORT")
		SSL_EMAIL, _  := cmd.Flags().GetString("SSLEMAIL")
		NOT_EMAIL, _  := cmd.Flags().GetString("NOTEMAIL")
		NOT_SLACK, _  := cmd.Flags().GetString("NOTSLACK")
		SLACK_CHN, _  := cmd.Flags().GetString("SLACKCHN")
		EMAIL_EN, _ := cmd.Flags().GetBool("EMAIL_EN")
		SLACK_EN, _ := cmd.Flags().GetBool("SLACKEN")
		s.StartServer(PORT, SSL_EMAIL, NOT_EMAIL, NOT_SLACK, SLACK_CHN, EMAIL_EN, SLACK_EN) ///note the order of parameters matters and the size can only be 2. This is a variadic argument 
	},
}

func init() {
	rootCmd.AddCommand(demoCmd)
	demoCmd.PersistentFlags().String("PORT", "443", "PORT to listen for incoming connections.")
	demoCmd.PersistentFlags().String("SSLEMAIL", "goshellydemo@araalinetworks.com", "Email address to generate SSL certificate.")
	demoCmd.PersistentFlags().String("NOTEMAIL", "all@araalinetworks.com", "Email to be notified after a client is connected.")
	demoCmd.PersistentFlags().String("NOTSLACK", "", "SLACK HOOK")
	demoCmd.PersistentFlags().String("SLACKCHN", "", "Slack channel to send the message to.")
	demoCmd.PersistentFlags().Bool("SLACKEN", false, "Enable/Disable email notifications")
	demoCmd.PersistentFlags().Bool("EMAILEN", false, "Enable/Disable email notifications")
}
