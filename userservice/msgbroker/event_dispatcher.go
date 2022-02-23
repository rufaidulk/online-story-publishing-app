package msgbroker

import (
	"log"
	"userservice/adapters"
	"userservice/msgbroker/events"
)

func UserFollowEventDispatch(userUuid string, followerUuid string) {
	log.Println("user follow event dispatching...")
	ch, err := adapters.GetRabbitmqConn().Channel()
	if err != nil {
		log.Println("Failed to open a channel")
		log.Fatal(err)
	}
	defer ch.Close()

	storyFeedEvent := events.NewUserFollowEvent(ch, "follow", userUuid, followerUuid)
	storyFeedEvent.CreateExchange().PublishMessage()
	log.Println("user follow event dispatched")
}

func UserUnFollowEventDispatch(userUuid string, followerUuid string) {
	log.Println("user unfollow event dispatching...")
	ch, err := adapters.GetRabbitmqConn().Channel()
	if err != nil {
		log.Println("Failed to open a channel")
		log.Fatal(err)
	}
	defer ch.Close()

	storyFeedEvent := events.NewUserFollowEvent(ch, "unfollow", userUuid, followerUuid)
	storyFeedEvent.CreateExchange().PublishMessage()
	log.Println("user unfollow event dispatched")
}
