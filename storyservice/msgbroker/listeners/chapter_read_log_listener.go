package listeners

import (
	"log"
	"storyservice/collections"
)

type ChapterReadLogListener struct {
	Data map[string]interface{}
}

func (u *ChapterReadLogListener) Handle() {
	log.Println("Executing the chapter read log event listener...")
	story := collections.NewStory()
	story.LoadById(u.Data["story_id"].(string))
	storyFeed := collections.NewStoryFeed()
	storyFeed.LoadByUser(u.Data["user_uuid"].(string))
	if err := storyFeed.AddCategoriesBasedOnReadingHistoryToDocument(story.Categories); err != nil {
		log.Fatal(err)
	}
	log.Println("chapter read log event listener completed")
}
