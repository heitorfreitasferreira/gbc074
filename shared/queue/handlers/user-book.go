package handlers

import "log"

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

func LogTopicMessage(topic UserBookTopic, msg []byte) {
	log.Printf("Mensagem recebida no tópico %s: %s\n", topic, msg)
}
