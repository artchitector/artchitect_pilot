package model

import "time"

type CreationState struct {
	PreviousCardID       uint
	Version              string
	Seed                 uint
	TagsCount            uint
	Tags                 []string
	LastCardPaintTime    uint // seconds
	CurrentCardPaintTime uint // seconds
	CardID               uint
	EnjoyTime            uint
	CurrentEnjoyTime     uint
}

type LotteryState struct {
	Lottery          Lottery
	EnjoyTotalTime   uint
	EnjoyCurrentTime uint
}

type PrayState struct {
	Queue   uint
	Started bool
}

type HeartState struct {
	Rnd []uint // Some random images for heart-entertainment. Usually 4 images
}

type EntropyState struct {
	Timestamp time.Time
	Phase     string
	Image     string // base64 jpeg-encoded
}

const (
	UnityStateCollectingChildren = "collecting_children"
	UnityStateUnifyChildren      = "unify_children"
	UnityStatePromoteLeads       = "promote_leads"
	UnityStatePrepareThumb       = "prepare_thumb"
	UnityStateFinish             = "finished"
)

type UnityState struct {
	Unifications []*UnityStateUnification
}

func (us *UnityState) Last() *UnityStateUnification {
	return us.Unifications[len(us.Unifications)-1]
}

func (us *UnityState) Add(unity Unity) {
	us.Unifications = append(us.Unifications, &UnityStateUnification{
		Unity: unity,
	})
}

func (us *UnityState) Remove() {
	us.Unifications = us.Unifications[:len(us.Unifications)-1]
}

func (us *UnityState) AddChild(unity Unity) {
	last := us.Last()
	last.AddChild(unity)
}

func (us *UnityState) AddLead(cardID uint) {
	last := us.Last()
	last.AddLead(cardID)
}

func (us *UnityState) SetThumb(thumb string) {
	last := us.Last()
	last.SetThumb(thumb)
}

func (us *UnityState) SetChildState(mask string, state string) {
	last := us.Last()
	last.SetChildState(mask, state)
}

func (us *UnityState) SetLeads(leads []uint) {
	last := us.Last()
	last.SetLeads(leads)
}

func (us *UnityState) SetUnity(unity Unity) {
	last := us.Last()
	last.SetUnity(unity)
}

func (us *UnityState) SetState(state string, currentProgress int, totalProgress int) {
	us.Last().SetState(state, currentProgress, totalProgress)
}

type UnityStateUnification struct {
	Unity           Unity
	Rank            uint
	State           string
	CurrentProgress int
	TotalProgress   int
	Children        []Unity
	Leads           []uint
	Thumb           string
}

func (uss *UnityStateUnification) SetUnity(unity Unity) {
	uss.Unity = unity
	uss.Rank = unity.Rank
}

func (uss *UnityStateUnification) SetState(state string, currentProgress int, totalProgress int) {
	uss.State = state
	uss.CurrentProgress = currentProgress
	uss.TotalProgress = totalProgress
}

func (uss *UnityStateUnification) AddChild(unity Unity) {
	uss.Children = append(uss.Children, unity)
}

func (uss *UnityStateUnification) AddLead(cardID uint) {
	uss.Leads = append(uss.Leads, cardID)
}

func (uss *UnityStateUnification) SetLeads(leads []uint) {
	uss.Leads = leads
}

func (uss *UnityStateUnification) SetThumb(thumb string) {
	uss.Thumb = thumb
}

func (uss *UnityStateUnification) SetChildState(mask string, state string) {
	for idx, child := range uss.Children {
		if child.Mask == mask {
			uss.Children[idx].State = state
			break
		}
	}
}
