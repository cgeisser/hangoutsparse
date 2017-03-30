package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type ChatLog struct {
	Continuation_end_timestamp string
	Conversation_state         []ConversationState
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
	for _, e := range csi.Event {
		output += csi.EventString(&e) + "\n"
	}
	return output
}

func (csi *ConversationStateInner) EventString(e *Event) string {
	return fmt.Sprintf("[%v] %v",
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

type ChatMessage struct {
	Message_content MessageContent
}

type MessageContent struct {
	Segment []ChatSegment
}

type ChatSegment struct {
	Text string
}

func main() {
	dec := json.NewDecoder(os.Stdin)
	var chatlog ChatLog

	err := dec.Decode(&chatlog)
	fmt.Println(err)
	fmt.Printf("%v\n", chatlog)
}
