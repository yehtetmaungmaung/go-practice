package fttxHelper

import (
	"cli-test/service"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func FiberMapsPortSync() {
	startTime := time.Now()
	log.Println("Port synchronization starting.........")
	syncPortProcess("OTB")
	log.Println("Port synchronization !!!", time.Since(startTime))
}

func syncPortProcess(deviceType string) {
	var nodes service.Nodes
	var nodeArray []service.NodeJson

	var filter bson.M
	var updatedFields bson.M

	nodes.GetFiberMapsNodesByDeviceType(1, deviceType)
	nodeArray = append(nodeArray, nodes.Nodes...)
	if len(nodes.Nodes) > 0 {
		getNodesByDeviceType(&nodeArray, getPageNo(nodes.NextLink), deviceType)
	}
	log.Println("Total"+deviceType+" node count from fiberMaps", len(nodeArray))

	for _, node := range nodeArray {
		filter = bson.M{"properties.id": node.SetNodeName()}
		fields := bson.M{
			"properties.additional_info.max_ports":         fmt.Sprintf("%v", node.NodeAttr.MaxPorts),
			"properties.additional_info.port_availability": fmt.Sprintf("%v", node.NodeAttr.MaxPorts-node.NodeAttr.UsedPorts),
			"properties.updated_on":                        time.Now(),
		}
		updatedFields = bson.M{"$set": fields}
	}
}

func getNodesByDeviceType(nodeArray *[]service.NodeJson, pageNo int, deviceType string) service.Nodes {
	defer func() {
		if x := recover(); x != nil {
			fmt.Printf("Panic: %v\n", x)
		}
	}()

	var nodes service.Nodes
	nodes.GetFiberMapsNodesByDeviceType(pageNo, deviceType)
	nodeLength := len(nodes.Nodes)
	if nodeLength > 0 {
		*nodeArray = append(*nodeArray, nodes.Nodes...)
		getNodesByDeviceType(nodeArray, getPageNo(nodes.NextLink), deviceType)
	}
	return nodes
}

func getPageNo(nextLink string) int {
	pageNoIndex := strings.Index(nextLink, "page=")
	perPageIndex := strings.Index(nextLink, "&per_page")
	pageNo := string(nextLink[pageNoIndex+5 : perPageIndex])
	pNumber, err := strconv.Atoi(pageNo)
	if err != nil {
		fmt.Println("Error converting page number")
		return 0
	}
	log.Println(pNumber)
	return pNumber
}
