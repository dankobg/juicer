package dbtype

type GameStats struct {
	Win   int32 `json:"win"`
	Loss  int32 `json:"loss"`
	Draw  int32 `json:"draw"`
	Total int32 `json:"total"`
}
