/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/dhquero/stress-test-go/internal/usecase"
	"github.com/dhquero/stress-test-go/pkg/web"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "stresstest",
	Short: "Stress test API client",
	Long: `This is a client to stress test an API.
	You configure the URL and the number of requests and workers.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		url, err := cmd.Flags().GetString("url")
		if err != nil {
			return err
		}

		if url == "" {
			return errors.New("you must provide a URL")
		}

		if !web.IsValidURL(url) {
			return errors.New("invalid URL")
		}

		_, err = cmd.Flags().GetUint("requests")
		if err != nil {
			return err
		}

		_, err = cmd.Flags().GetUint("concurrency")
		if err != nil {
			return err
		}

		_, err = cmd.Flags().GetUint("timeout")
		if err != nil {
			return err
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		url, _ := cmd.Flags().GetString("url")
		requests, _ := cmd.Flags().GetUint("requests")
		concurrency, _ := cmd.Flags().GetUint("concurrency")
		timeout, _ := cmd.Flags().GetUint("timeout")

		stressTest := usecase.NewStressTestUseCase(url, requests, concurrency, timeout)

		report, err := stressTest.Execute()

		if err != nil {
			return err
		}

		println("REPORT - STRESS TEST")
		println("======================================================================")
		println(fmt.Sprintf("URL: %s", report.URL))
		println(fmt.Sprintf("Requests: %d", report.Requests))
		println(fmt.Sprintf("Workers: %d", report.Concurrency))
		println(fmt.Sprintf("Timeout: %.0f seconds", report.Timeout.Seconds()))
		println(fmt.Sprintf("Total time: %s", report.TotalTime))
		println("")
		println(fmt.Sprintf("Total code 200: %d", report.StatusCode[http.StatusOK]))
		if len(report.StatusCode) > 0 {
			delete(report.StatusCode, http.StatusOK)
			println("Other status codes")
			for k, v := range report.StatusCode {
				println(fmt.Sprintf("Code %d: %d", k, v))
			}
		}
		println("======================================================================")

		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().String("url", "", "Target URL")
	rootCmd.Flags().Uint("requests", 1, "Number of requests")
	rootCmd.Flags().Uint("concurrency", 1, "Number of workers")
	rootCmd.Flags().Uint("timeout", 1, "Timeout (seconds)")
	rootCmd.MarkFlagRequired("url")
}
