package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var (
	broker   = "broker:9092"
	topic    = "marbles"
	outTopic = "out"
	group    = "marbles-grp"
)

var rootCmd = &cobra.Command{
	Use:   "strmproc",
	Short: "Stream processor",
	Run: func(cmd *cobra.Command, args []string) {
		proc, err := newProcessor([]string{broker}, group, topic, outTopic)

		if err != nil {
			log.Fatalln(err)
		}

		ctx := context.Background()

		go func() {
			if err := proc.Run(ctx); err != nil {
				log.Fatalln(err)
			}
		}()
		wait := make(chan os.Signal, 1)
		signal.Notify(wait, syscall.SIGINT, syscall.SIGTERM)
		<-wait // wait for SIGINT/SIGTERM
		ctx, cancel := context.WithTimeout(ctx, time.Second*10)
		defer cancel()
		proc.Close()
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	// default value of the flag
	rootCmd.PersistentFlags().StringVar(&broker, "broker", "broker:9092", "default topic is 'broker:9092'")
	rootCmd.PersistentFlags().StringVar(&topic, "topic", "marbles", "default topic is 'marbles'")
	rootCmd.PersistentFlags().StringVar(&outTopic, "outTopic", "out", "default out topic prefix is 'out'")
	rootCmd.PersistentFlags().StringVar(&group, "group", "marbles-grp", "default is 'marbles-grp'")

}

func initConfig() {

}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
