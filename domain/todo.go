package domain

type Todo struct {
	Id      string `json:id`
	User    string `json:user`
	Subject string `json:subject`
	Status  string `json:status`
}
