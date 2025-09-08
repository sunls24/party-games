package wzq

type Conf struct {
}

type State struct {
	Black   int8     `json:"black"`
	Current int8     `json:"current"`
	Win     int8     `json:"win"`
	LastIdx uint8    `json:"lastIdx"`
	LastCt  int8     `json:"lastCt"`
	Regret  int8     `json:"regret"`
	Board   [][]int8 `json:"board"`
}
