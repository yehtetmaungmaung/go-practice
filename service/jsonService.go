package service

import "strings"

// Nodes is a list of nodes that are displayed in fibermaps.
// NextLink takes you to the next page in fibermaps UI.
//
//	{
//		Nodes: [
//			{
//				"Name": "CA2-B21XXZZINC-B03",
//				"NodeAttr": {
//					"Type": "SPLITTER",
//					"MaxPorts": 16,
//					"UsedPorts": 15,
//					"Tags":
//					"ConnectionType":
//				}
//				"Township":
//				"Latitude":
//				"Longitude":
//			},
//			{
//				"Name":"CA2-A4XXXZZINC-C05",
//				"NodeAttr": {
//					"Type":"SPLITTER",
//					"MaxPorts": 16
//					"UsedPorts": 14
//					"Tags":
//					"ConnectionType":
//				}
//				"Township":
//				"Latitude":
//				"Longitude":
//			}
//		]
//		"NextLink": "http://fiber-maps-api-qa.svc.cluster.local/pets/api/nodes/external/node/check_port?page=3&per_page=2&device_type=CA2"
//		"TotalCount": 24962
//	}
type Nodes struct {
	Nodes      []NodeJson `json:"nodes" bson:"nodes"`
	NextLink   string     `json:"next_link" bson:"next_link"`
	TotalCount int        `json:"total_count" bson:"total_count"`
}

type NodeJson struct {
	Name      string        `json:"name" bson:"name"`
	NodeAttr  NodeAttribute `json:"node_attr" bson:"node_attr"`
	Township  string        `json:"township" bson:"township"`
	Latitude  string        `json:"latitude" bson:"latitude"`
	Longitude string        `json:"longitude" bson:"longitude"`
}

type NodeAttribute struct {
	Type           string      `json:"device_type" bson:"device_type"`
	MaxPorts       float32     `json:"max_ports" bson:"max_ports"`
	UsedPorts      float32     `json:"used_ports" bson:"used_ports"`
	Tags           interface{} `json:"tags" bson:"tags"`
	ConnectionType string      `json:"connection_type" bson:"connection_type"`
}

func (node *NodeJson) SetNodeName() string {
	var name string
	if node.NodeAttr.Type == "CPE" {
		if len(node.Name) > 4 {
			nameList := strings.Split(node.Name, "-")
			name = nameList[len(nameList)-1]
		}
	} else {
		name = node.Name
	}
	return name
}
