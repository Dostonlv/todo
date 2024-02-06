package structs

type User struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Password  string `json:"password"`
	CreatedAt string `json:"createdAt"`
}

type UserList struct {
	Users []User `json:"users"`
	Count int    `json:"count"`
}

type Filter struct {
	Search string `json:"search"`
	Active string `json:"active"`
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
}
