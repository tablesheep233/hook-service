package server

type ReleaseArg struct {
	Action string `json:"action"`
	Release Release `json:"release"`
	Repo Repository `json:"repository"`
}

type Release struct {
	Tag_name string `json:"tag_name"`
}

type Repository struct {
	Name string `json:"name"`
}

type ResponseMsg struct {
	Msg string `json:"msg"`
}
