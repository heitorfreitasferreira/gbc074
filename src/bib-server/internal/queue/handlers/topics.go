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

type BookTopic string

const (
	BookCreateTopic BookTopic = "book/create"
	BookRemoveTopic BookTopic = "book/remove"
	BookUpdateTopic BookTopic = "book/update"
)

type UserTopic string

const (
	UserCreateTopic UserTopic = "user/create"
	UserUpdateTopic UserTopic = "user/update"
	UserRemoveTopic UserTopic = "user/remove"
)
