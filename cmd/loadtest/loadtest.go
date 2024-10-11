package loadtest

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	// Cmd is the load test command.
	Cmd = &cobra.Command{
		Use:     "loadtest",
		Short:   "Perform load testing by sending FIX orders",
		Long:    "Load testing tool for sending multiple FIX orders at a configurable rate.",
		Example: "qf loadtest --orders 100 --rate 10",
		RunE:    execute,
	}

	orderCount int
	rate        int
)

func init() {
	Cmd.Flags().IntVarP(&orderCount, "orders", "o", 100, "Number of orders to send")
	Cmd.Flags().IntVarP(&rate, "rate", "r", 10, "Orders per second")
}

func execute(cmd *cobra.Command, args []string) error {
	fmt.Printf("Starting load test: %d orders at %d orders/second\n", orderCount, rate)

	for i := 0; i < orderCount; i++ {
		sendOrder(i) // Replace this with actual order sending logic
		time.Sleep(time.Second / time.Duration(rate)) // Control the rate
	}

	return nil
}

// sendOrder is a stub for sending an order. Replace with actual implementation.
func sendOrder(orderID int) {
	fmt.Printf("Order %d sent\n", orderID)
}
