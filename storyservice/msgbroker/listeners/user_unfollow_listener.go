package listeners

import (
	"log"
	"storyservice/collections"
)

type UserUnfollowEventData struct {
	Data map[string]interface{}
}

func (u *UserUnfollowEventData) Handle() {
	log.Println("Executing the user unfollow event listener...")
	storyFeed := collections.NewStoryFeed()
	storyFeed.LoadByUser(u.Data["user_uuid"].(string))
	if err := storyFeed.RemoveFollowingAuthorFromDocument(u.Data["follower_uuid"].(string)); err != nil {
		log.Fatal(err)
	}
	log.Println("User unfollow event listener completed")
}
