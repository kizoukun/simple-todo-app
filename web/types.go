package web

type Todo struct {
	Task      string `json:"task"`
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
}

type TodoGetResponse struct {
	Todos []Todo `json:"todos"`
}

type ToggleTodoRequest struct {
	Completed bool `json:"completed"`
}
