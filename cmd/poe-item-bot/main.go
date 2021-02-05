package main

import (
	"fmt"
	"time"

	"github.com/Deichindianer/poe-item-bot/internal/characterpoller"
	"github.com/Deichindianer/poe-item-bot/internal/ladderpoller"
)

func main() {
	fmt.Println("Starting ladder poll")
	lp := ladderpoller.NewLadderPoller("SSF Ritual HC")
	lp.Poll(time.Minute)
	// wait for the ladderPoller to do some work for testing
	time.Sleep(5)
	var pollList []characterpoller.PollCharacter
	for _, entry := range lp.Ladder.Entries {
		pollList = append(
			pollList,
			characterpoller.PollCharacter{
				AccountName:   entry.Account.Name,
				CharacterName: entry.Character.Name,
			},
		)
	}
	fmt.Printf("Length of pollList: %d", len(pollList))
	fmt.Println("Finished ladder poll")
	characterPoller := characterpoller.NewCharacterPoller(pollList)
	characterPoller.Poll(time.Minute)
	lp.StopPoll()
	i := 0
	for i < 5 {
		fmt.Println("Checking for characters...")
		for _, char := range characterPoller.Characters {
			fmt.Printf("Current character: %s\n", char.Character.Name)
		}
		time.Sleep(60 * time.Second)
	}
	characterPoller.StopPoll()
}
