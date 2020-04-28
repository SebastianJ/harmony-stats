package logger

import (
	"fmt"
	"strings"
	"time"

	"github.com/SebastianJ/harmony-stats/config"
	"github.com/gookit/color"
)

var (
	timeFormat = "2006-01-02 15:04:05"
)

// Title - outputs the application title
func Title() {
	fmt.Println()
	config.Configuration.Styling.Header.Println(
		fmt.Sprintf("\tHarmony Stats - Network: %s (%s mode)%s",
			strings.Title(config.Configuration.Network.Name),
			strings.ToUpper(config.Configuration.Network.Mode),
			strings.Repeat("\t", 15),
		),
	)
	fmt.Println()
}

// Log - logs default testing messages
func Log(message string) {
	OutputLog(message, "default")
}

// AccountLog - logs account related messages
func AccountLog(message string) {
	OutputLog(message, "account")
}

// FundingLog - logs funding related messages
func FundingLog(message string) {
	OutputLog(message, "funding")
}

// BalanceLog - logs balance related messages
func BalanceLog(message string) {
	OutputLog(message, "balance")
}

// TransactionLog - logs transaction related messages
func TransactionLog(message string) {
	OutputLog(message, "transaction")
}

// StakingLog - logs staking related messages
func StakingLog(message string) {
	OutputLog(message, "staking")
}

// TeardownLog - logs teardown related messages
func TeardownLog(message string) {
	OutputLog(message, "teardown")
}

// InfoLog - logs success related messages
func InfoLog(message string) {
	OutputLog(message, "info")
}

// SuccessLog - logs success related messages
func SuccessLog(message string) {
	OutputLog(message, "success")
}

// WarningLog - logs error related messages
func WarningLog(message string) {
	OutputLog(message, "warning")
}

// ErrorLog - logs error related messages
func ErrorLog(message string) {
	OutputLog(message, "error")
}

// ResultLog - logs result related messages - will switch between green (successful) and red (failed) depending on the passed boolean
func ResultLog(result bool, expected bool) {
	if config.Configuration.Verbose {
		var formattedCategory string
		message := fmt.Sprintf("Test successful: %t, Expected: %t", result, expected)
		formattedMessage := ResultColoring(result, expected).Render(message)

		if result == expected {
			formattedCategory = color.Style{color.FgGreen, color.OpBold}.Render("RESULT")
		} else {
			formattedCategory = color.Style{color.FgRed, color.OpBold}.Render("RESULT")
		}

		fmt.Println(fmt.Sprintf("\n[%s] %s - %s", time.Now().Format(timeFormat), formattedCategory, formattedMessage))
	}
}

// ResultColoring - generate a green or red color setup depending on if the result matches the expected result
func ResultColoring(result bool, expected bool) *color.Style {
	if result == expected {
		return config.Configuration.Styling.Success
	}

	return config.Configuration.Styling.Error
}

// OutputLog - time stamped logging messages for test cases
func OutputLog(message string, category string) {
	var c *color.Style

	switch category {
	case "default":
		c = config.Configuration.Styling.Default
	case "account":
		c = config.Configuration.Styling.Account
	case "funding":
		c = config.Configuration.Styling.Funding
	case "balance":
		c = config.Configuration.Styling.Balance
	case "transaction":
		c = config.Configuration.Styling.Transaction
	case "staking":
		c = config.Configuration.Styling.Staking
	case "teardown":
		c = config.Configuration.Styling.Teardown
	case "info":
		c = config.Configuration.Styling.Info
	case "success":
		c = config.Configuration.Styling.Success
	case "warning":
		c = config.Configuration.Styling.Warning
	case "error":
		c = config.Configuration.Styling.Error
	default:
		c = config.Configuration.Styling.Default
	}

	formattedCategory := c.Render(strings.ToUpper(category))
	title := fmt.Sprintf("[%s] %s - ", time.Now().Format(timeFormat), formattedCategory)
	fmt.Println(fmt.Sprintf("%s%s", title, message))
}
