package turingpi

type SDCard struct {
	Total int64 `json:"total"`
	Free  int64 `json:"free"`
	Use   int64 `json:"use"`
}

type SDCardResponse struct {
	Response []SDCard `json:"response"`
}

type NodeInfo struct {
	Node1 string `json:"node1"`
	Node2 string `json:"node2"`
	Node3 string `json:"node3"`
	Node4 string `json:"node4"`
}

type NodeInfoResponse struct {
	Response []NodeInfo `json:"response"`
}
