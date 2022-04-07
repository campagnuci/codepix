/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/campagnuci/codepix/go/application/grpc"
	"github.com/campagnuci/codepix/go/application/kafka"
	"github.com/campagnuci/codepix/go/infrastructure/db"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/spf13/cobra"
)

var gRPCPortNumber int

// allCmd represents the all command
var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Run gRPC and Apache Kafka Consumer",
	Run: func(cmd *cobra.Command, args []string) {
		database := db.ConnectDB(os.Getenv("env"))
		go grpc.StartGrpcServer(database, portNumber)

		deliveryChan := make(chan ckafka.Event)
		producer := kafka.NewKafkaProducer()
		go kafka.DeliveryReport(deliveryChan)
		kafkaProcessor := kafka.NewKafkaProcessor(database, producer, deliveryChan)
		kafkaProcessor.Consume()
	},
}

func init() {
	rootCmd.AddCommand(allCmd)
	allCmd.Flags().IntVarP(&gRPCPortNumber, "grpc-port", "p", 50051, "gRPC Server Port")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// allCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// allCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
