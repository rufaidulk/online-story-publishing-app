package listeners

import "log"

type UserFollowEventData struct {
	UserUuid     string
	FollowerUuid string
}

func (u *UserFollowEventData) Handle() {
	log.Println("Executing the user follow event listener")
}
