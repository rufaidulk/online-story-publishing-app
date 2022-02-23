package msgbroker

import (
	"log"
	"storyservice/adapters"
	"storyservice/msgbroker/events"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ChapterReadLogEventDispatch(userUuid string, storyId primitive.ObjectID) {
	log.Println("chapter read log event dispatching...")
	ch, err := adapters.GetRabbitmqConn().Channel()
	if err != nil {
		log.Println("Failed to open a channel")
		log.Fatal(err)
	}
	defer ch.Close()

	event := events.NewChapterReadLogEvent(ch, "chapter_read_log", userUuid, storyId)
	event.CreateExchange().PublishMessage()
	log.Println("chapter read log event dispatched")
}
