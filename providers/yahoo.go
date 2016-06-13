package providers

type JsonResp struct {
	List JsonRespList
}

type JsonRespMeta struct {
	Count int
}

type JsonRespList struct {
	Resources []JsonRespResourceCont
	Meta      JsonRespMeta
}

type JsonRespResourceCont struct {
	Resource JsonRespResource
}

type JsonRespResource struct {
	Fields JsonRespFields
}

type JsonRespFields struct {
	Price       float64 `json:",string"`
	Chg_percent float64 `json:",string"`
}
