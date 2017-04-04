package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"
)

type ChatLog struct {
	Conversation_state []ConversationState
}

type ConversationState struct {
	Conversation_state ConversationStateInner
}

type ConversationStateInner struct {
	Event        []Event
	Conversation Conversation
}

func (csi ConversationStateInner) String() string {
	output := ""
	sort.Sort(EventByTime(csi.Event))
	for _, e := range csi.Event {
		output += csi.EventString(&e) + "\n"
	}
	return output
}

func (csi *ConversationStateInner) EventString(e *Event) string {
	return fmt.Sprintf("%v [%v] %v",
		e.Time().Format("2006-01-02 15:04:05"),
		csi.Conversation.GetNameForId(e.Sender_id.Gaia_id),
		e.Chat_message)
}

type Conversation struct {
	Participant_data []Participant

	participantNames map[string]string
}

func (c *Conversation) GetNameForId(id string) string {
	if len(c.participantNames) == 0 {
		//initialize the map
		c.participantNames = make(map[string]string)
		for _, i := range c.Participant_data {
			c.participantNames[i.Id.Gaia_id] = i.Fallback_name
		}
	}
	return c.participantNames[id]
}

type Participant struct {
	Id            SenderId
	Fallback_name string
}

type SenderId struct {
	Gaia_id string
}

type Event struct {
	Sender_id    SenderId
	Timestamp    string
	Chat_message ChatMessage
}

func (e *Event) Time() time.Time {
	timestamp, _ := strconv.ParseInt(e.Timestamp, 10, 64)
	return time.Unix(timestamp/1000000, timestamp%1000000)
}

type EventByTime []Event

func (a EventByTime) Len() int           { return len(a) }
func (a EventByTime) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a EventByTime) Less(i, j int) bool {
	return a[i].Time().Before(a[j].Time())
}

type ChatMessage struct {
	Message_content MessageContent
}

func (cm ChatMessage) String() string {
	return cm.Message_content.String()
}

type MessageContent struct {
	Segment []ChatSegment
}

func (mc MessageContent) String() string {
	output := ""
	for _, segment := range mc.Segment {
		output += segment.String()
	}
	return output
}

type ChatSegment struct {
	Text string
}

func (cs ChatSegment) String() string {
	return cs.Text
}

func main() {
	dec := json.NewDecoder(os.Stdin)
	var chatlog ChatLog

	dec.Decode(&chatlog)
	fmt.Printf("%v\n", chatlog)
}
