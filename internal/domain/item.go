package domain

type Item struct {
	ID          int64
	WishlistID  int64
	Title       string
	Description string
	URL         string
	Priority    int
	IsReserved  bool
}
