package main

import (
	"fmt"
	"time"

	"github.com/Deichindianer/poe-item-bot/internal/ladderpoller"
)

func main() {
	ladderPoller := ladderpoller.NewLadderPoller("SSF Ritual HC")
	ladderPoller.Poll()

	i := 0
	for i < 5 {
		fmt.Printf("Ladder cached since: %+v\n", ladderPoller.Ladder.CachedSince)
		fmt.Printf("Ladder updated? %+v\n", ladderPoller.Updated)
		time.Sleep(60 * time.Second)
	}
	ladderPoller.StopPoll()
}
