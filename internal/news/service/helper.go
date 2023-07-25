package service

type Sort string

var (
	SortDESC Sort = "DESC"
	SortASC  Sort = "ASC"
)

func (s Sort) Valid() bool {
	switch s {
	case SortDESC:
		return true
	case SortASC:
		return true
	}
	return false
}
