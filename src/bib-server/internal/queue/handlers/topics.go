package handlers

type UserBookTopic string

const (
	BookLoanTopic         UserBookTopic = "book/loan"
	BookReturnTopic       UserBookTopic = "book/return"
	BookListBorrowedTopic UserBookTopic = "book/list-borrowed"
	BookListMissingTopic  UserBookTopic = "book/list-missing"
	BookSearchTopic       UserBookTopic = "book/search"
	UserBlockTopic        UserBookTopic = "user/block"
	UserFreeTopic         UserBookTopic = "user/free"
	UserListBlockedTopic  UserBookTopic = "user/list-blocked"
)
