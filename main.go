package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"text/tabwriter"

	"aurora_version_check/logger_util"
	"aurora_version_check/rds_util"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

func readInput(prompt string, validInputs []string) string {
	for {
		fmt.Printf("%s: \n", prompt)
		var input string
		_, err := fmt.Scan(&input)
		if err != nil {
			fmt.Printf("Input Data Error: %v\n", err)
			continue
		}
		for _, validInput := range validInputs {
			if input == validInput {
				return input
			}
		}
		fmt.Printf("Input correct %s(%s)\n", prompt, input)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide a number as argument.")
		return
	}

	num, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		fmt.Printf("Error parsing number: %v", err)
		return
	}

	if num < 1001 {
		fmt.Println("Input more than 1000 milliseconds")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 0, ' ', tabwriter.AlignRight)

	fmt.Printf("================================\n")
	fmt.Printf("Aurora Version & Parameter Check\n")
	fmt.Printf("================================\n")

	region := readInput("Region(kr/jp/ca/uk)", []string{"kr", "jp", "ca", "uk"})
	worktype := readInput("WorkType(cluster/instance)", []string{"cluster", "instance"})

	regionMap := map[string]string{
		"kr": "ap-northeast-2",
		"jp": "ap-northeast-1",
		"ca": "ca-central-1",
		"uk": "eu-west-2",
	}

	realRegion, ok := regionMap[region]
	if !ok {
		fmt.Println("Invalid region provided")
		return
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
		Config: aws.Config{
			Region: aws.String(realRegion),
		},
	})

	if err != nil {
		log.Fatal("Failed to create session:", err)
	}

	var describable rds_util.Describable
	var logger *log.Logger
	var cleanup func()

	if worktype == "cluster" {
		describable = rds_util.Cluster{}
		logger, cleanup, err = logger_util.NewLogger("aurora_cluster_version_check.log")
		if err != nil {
			log.Fatalf("Failed to open cluster log file: %v", err)
		}
	} else {
		describable = rds_util.Instance{}
		logger, cleanup, err = logger_util.NewLogger("aurora_instance_version_check.log")
		if err != nil {
			log.Fatalf("Failed to open instance log file: %v", err)
		}
	}
	defer cleanup()

	createFile := readInput("Do you want to create the list file? (yes/no)", []string{"yes", "no"})
	if createFile == "yes" {
		fileName := fmt.Sprintf("./db_%s_list_%s.txt", worktype, region)
		err := describable.SaveListToFile(rds_util.NewRDS(sess), fileName)
		if err != nil {
			log.Fatalf("Failed to save list to file: %v", err)
		}
		fmt.Printf("List saved to file: %s\n", fileName)
		return
	}

	svc, filteredLines := rds_util.PrepareCheck(region, worktype, sess)
	rds_util.AuroraVersionParam(svc, filteredLines, w, num, describable, logger)
}
