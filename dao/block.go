package dao

type Block struct {
	ID          int64  `json:"id"`
	BlockNumber string `json:"block_number"`
	BlockHash   string `json:"block_hash"`
	ParentHash  string `json:"parent_hash"`
	CreateTime  int64  `json:"create_time"`
	Fork        bool   `json:"fork"`
}
