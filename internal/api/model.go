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

type Power struct {
	Node1 int64 `json:"node1"`
	Node2 int64 `json:"node2"`
	Node3 int64 `json:"node3"`
	Node4 int64 `json:"node4"`
}

type PowerResponse struct {
	Response []Power `json:"response"`
}

type Usb struct {
	Mode int64 `json:"mode"`
	Node int64 `json:"node"`
}

type UsbResponse struct {
	Response []Usb `json:"response"`
}
