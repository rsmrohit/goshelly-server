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
		SSL_EMAIL, _  := cmd.Flags().GetString("EMAIL")
		s.StartServer(PORT, SSL_EMAIL) ///note the order of parameters matters and the size can only be 2. This is a variadic argument 
	},
}

func init() {
	rootCmd.AddCommand(demoCmd)
	demoCmd.PersistentFlags().String("PORT", "443", "PORT to listen for incoming connections.")
	demoCmd.PersistentFlags().String("EMAIL", "goshellydemo@araalinetworks.com", "Email address to generate SSL certificate.")

	
	


}
