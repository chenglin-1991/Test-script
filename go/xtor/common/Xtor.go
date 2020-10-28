package common

const (
	XTOR_INITIAL_ADMIN_PASSWORD string = "x70r@x7a0"
)

type XTXtorAPIRequest struct {
	VolName string `json:"VolName"`
}

type XTXtorAPIReply struct {
	Status int `json:"Status"`
	Errmsg string `json:"Errmsg"`
	Result string `json:"Result"`
}

/*
 * Du API message
 */
type XTXtorDuAPIRequest struct {
	VolName string `json:"VolName"`
	Path string `json:"Path"`
	Obj bool `json:"Obj"`
}

type XTXtorDuReply struct {
	Size string `json:"Size"`
	Dirs string `json:"Dirs"`
	Files string `json:"Files"`
}

type XTXtorAPIDuResult struct {
	Status string `json:"Status"`
	Errmsg string `json:"Errmsg"`
	Result XTXtorDuReply `json:"Result"`
}

/*
 * Fsstat API message
 */
type XTXtorBrickStat struct {
	Device string `json:"Device"`
	Dir string `json:"Dir"`
	Total string `json:"Total"`
	Free  string `json:"Free"`
}

type XTXtorFsstatReply struct {
	Bricks []XTXtorBrickStat `json:"Bricks"`
	Total float64 `json:"Total"`
        Free float64 `json:"Free"`
}

type XTXtorAPIFsstatResult struct {
	Status string `json:"Status"`
	Errmsg string `json:"Errmsg"`
	Result XTXtorFsstatReply `json:"Result"`	
}

/*
 * API request
 */
type XTXtorQuotaListAPIRequest struct {
	VolName string `json:"VolName"`
	Path string `json:"Path"`
	Obj bool `json:"Obj"`
}

type XTXtorQuotaSetAPIRequest struct {
	VolName string `json:"VolName"`
	Path string `json:"Path"`
	Obj bool `json:"Obj"`
	Limit string `json:"Limit"`
}
