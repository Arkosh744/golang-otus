package service

type Repository interface {
}

type calendarService struct {
	repo Repository
}

func New(repo Repository) *calendarService {
	return &calendarService{
		repo: repo,
	}
}
