package listeners

import (
	"log"
	"storyservice/collections"
)

type UserFollowEventListener struct {
	Data map[string]interface{}
}

func (u *UserFollowEventListener) Handle() {
	log.Println("Executing the user follow event listener...")
	storyFeed := collections.NewStoryFeed()
	storyFeed.LoadByUser(u.Data["user_uuid"].(string))
	storyFeed.FollowingAuthors = append(storyFeed.FollowingAuthors, u.Data["follower_uuid"].(string))
	if err := storyFeed.AddFollowingAuthorToDocument(); err != nil {
		log.Fatal(err)
	}
	log.Println("User follow event listener completed")
}
