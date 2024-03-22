package service

type Sessions struct {
	r SessionsRepo
}

func NewSessions(r SessionsRepo) *Sessions {
	return &Sessions{r}
}
