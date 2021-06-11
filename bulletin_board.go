package main

import (
	"sync"
	"time"

	"github.com/diamondburned/arikawa/v2/discord"
)

// BulletinBoard maps an address -> messageID
type BulletinBoard struct {
	sync.Mutex
	ChannelID           uint64
	AddressMessageMap   map[string]discord.MessageID
	MessageSchedulerMap map[discord.MessageID]*Scheduler
}

func NewBulletinBoard(channelID uint64) *BulletinBoard {
	return &BulletinBoard{
		ChannelID:           channelID,
		AddressMessageMap:   make(map[string]discord.MessageID, 16),
		MessageSchedulerMap: make(map[discord.MessageID]*Scheduler, 16),
	}
}

func (bb *BulletinBoard) Register(address string, msgID discord.MessageID, interval time.Duration, f func() error) {
	bb.Lock()
	defer bb.Unlock()

	bb.AddressMessageMap[address] = msgID
	bb.MessageSchedulerMap[msgID] = NewScheduler()

	bb.MessageSchedulerMap[msgID].Start(interval, f)
}

func (bb *BulletinBoard) UpdateMsgID(address string, new discord.MessageID) {
	bb.Lock()
	defer bb.Unlock()

	// old values
	oldId := bb.AddressMessageMap[address]
	scheduler := bb.MessageSchedulerMap[oldId]

	// update values
	bb.AddressMessageMap[address] = new
	bb.MessageSchedulerMap[new] = scheduler

	// delete old id reference
	delete(bb.MessageSchedulerMap, oldId)
}
