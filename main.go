package main

import (
	"cli-test/helper/fttxHelper"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	if len(os.Args) >= 2 {
		switch os.Args[1] {
		case "customer":
			customerSubCommands()
		case "fttx":
			fttxSubCommands()
		default:
			showHelper()
		}
	} else {
		showHelper()
	}
}

// showHelper prints command arguments
func showHelper() {
	fmt.Println("Usage ./mvp-cli command [OPTION] inputfile")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("customer                              choose customer or fttx to go to next subcommands")
	fmt.Println("fttx                                  choose customer or fttx to go to next subcommands")
	fmt.Println()
	fmt.Println("Options for customer command")
	fmt.Println("customer [OPTION]")
	fmt.Println("--subscription                        add subscription data (Usage : ./mvp-cli customer --subscription inputfile data_source_name) ")
	fmt.Println("--profile                             add profile data (Usage : ./mvp-cli  customer --profile inputfile data_source_name)")
	fmt.Println("--cloudShareRT                        generate the latest csv and sync data from cloudshare's rtCrm Folder (Usage : ./mvp-cli customer --cloudShareRT)")
	fmt.Println("--cloudShareFirstBytes                generate the latest csv and sync data from cloudshare's FirstByte Folder (Usage : ./mvp-cli customer --cloudShareFirstBytes)")
	fmt.Println("--checkMissedFiles                    generate the missing files from cloudshare's rtCrm Folder (Usage : ./mvp-cli customer --checkMissedFiles)")
	fmt.Println("--deleteLogFiles                      delete daily sync's log files from cloudshare's Logs Folder (Usage : ./mvp-cli customer --deleteLogFiles)")
	fmt.Println()
	fmt.Println("Options for fttx command")
	fmt.Println("fttx [OPTION]")
	fmt.Println("--fiberMapsDataOneTimeSync            Sync FTTx Data For the one time (Usage : ./mvp-cli fttx --fiberMapsDataOneTimeSync [node.json file path] [edge.json file path])")
	fmt.Println("--syncCPEStatus                       Sync Last CPE status from CPEMS , not the real time (Usage : ./mvp-cli fttx --syncCPEStatus)")
	fmt.Println("--syncSplitterStatus                  Sync status of CA2 and OTB Status based on CPE status (Usage : ./mvp-cli fttx --syncSplitterStatus)")
	fmt.Println("--syncOLTMSStatus                     Sync status of OLT and ONU (Usage : ./mvp-cli fttx --syncOLTMSStatus)")
	fmt.Println("--syncFiberMaps                       Sync with Fibermaps through RabbitMQ (Usage : ./mvp-cli fttx --syncFiberMaps)")
	fmt.Println("--syncDLFiberMaps                     Sync with Fibermaps through RabbitMQ (Usage : ./mvp-cli fttx --syncDLFiberMaps)")
	fmt.Println("--oltBasicInfoSync                    Sync OLT device info from OLTMS through OLTMS API (Usage : ./mvp-cli fttx --oltBasicInfoSync)")
	fmt.Println("--syncCPELatLngWithSubscription       Sync OLT device info from OLTMS through OLTMS API (Usage : ./mvp-cli fttx --syncCPELatLngWithSubscription)")
	fmt.Println("--portUsageOneTimeSync                Sync FTTx Port Usages For the one time (Usage : ./mvp-cli fttx --portUsageOneTimeSync [node.json file path])")
	fmt.Println("--fiberMapsTagsOneTimeSync            Sync FTTx Tags For the one time (Usage : ./mvp-cli fttx --fiberMapsTagsOneTimeSync [node.json file path])")
	fmt.Println("--updateFTSBSType                     Modify FT-SBS type to CPE For one time (Usage : ./mvp-cli fttx --updateFTSBSType")
	fmt.Println("--fiberMapsConneTypeOneTimeSync       Sync FTTx Connection Type For the one time (Usage : ./mvp-cli fttx --fiberMapsConnecTypeOneTimeSync [node.json file path])")
	fmt.Println("--portUsageRegularSync                Sync FTTx Port Usages For daily (Usage : ./mvp-cli fttx --portUsageRegularSync)")
	fmt.Println("--updateDuplicatePoints               Fix FTTx duplicate and update points (Usage : ./mvp-cli fttx --updateDuplicatePoints)")
	fmt.Println()
	fmt.Println("--help                                For more information on a command(Usage : ./mvp-cli --help or ./mvp-cli command --help")
}

func customerSubCommands() {
	fmt.Println("--customer")
}

func fttxSubCommands() {
	if len(os.Args) == 2 {
		showHelper()
	} else {
		switch os.Args[2] {
		case "--portUsageRegularSync":
			portUsageRegularSync()
		default:
			showHelper()
		}
	}
}

func portUsageRegularSync() {
	fttxHelper.FiberMapsPortSync()
}
