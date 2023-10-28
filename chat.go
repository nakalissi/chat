package entity

import (
	"errors"

	"github.com/google/uuid"
)

func main() {
	print("Welcome to chat")
}

type ChatConfig struct {
	Model            *Model
	Temperature      float32
	TopP             float32
	N                int
	Stop             []string
	MaxTokens        int
	PresencePenalty  float32
	FrequencyPenalty float32
}

type Chat struct {
	ID                   string
	UserID               string
	InitialSystemMessage *Message
	Messages             []*Message
	ErasedMessages       []*Message
	Status               string
	TokenUsage           int
	Config               *ChatConfig
}

func NewChat(UserID string, InitialSystemMessage *Message, chatConfig *ChatConfig) (*Chat, error) {
	chat := &Chat{
		ID:                   uuid.New().String(),
		UserID:               UserID,
		InitialSystemMessage: InitialSystemMessage,
		Status:               "active",
		Config:               chatConfig,
		TokenUsage:           0,
	}
	chat.AddMessage(InitialSystemMessage)

	if err := chat.Validate(); err != nil {
		return nil, err
	}
	return chat, nil
}

func (c *Chat) Validate() error {
	if c.UserID == "" {
		return errors.New("User id is empty")
	}
	if c.Status != "active" && c.Status != "ended" {
		return errors.New("Invalid status")
	}
	if c.Config.Temperature < 0 || c.Config.Temperature > 2 {
		return errors.New("Invalid temperature")
	}
	return nil
}

func (c *Chat) AddMessage(m *Message) error {
	if c.Status == "ended" {
		return errors.New("chat is ended. no more messages allowed")
	}

	for {
		if c.Config.Model.GetMaxTokens() >= m.GetQtdTokens()+c.TokenUsage {
			c.Messages = append(c.Messages, m)
		}
		c.ErasedMessages = append(c.ErasedMessages, c.Messages[0])
		c.Messages = c.Messages[1:]
		c.RefreshTokensUsage()
	}
}

func (c *Chat) GetMessages() []*Message {
	return c.Messages
}

func (c *Chat) CountMessages() int {
	return len(c.Messages)
}

func (c *Chat) End() {
	c.Status = "ended"
}

func (c *Chat) RefreshTokensUsage() {
	c.TokenUsage = 0
	for m := range c.Messages {
		c.TokenUsage += c.Messages[m].GetQtdTokens()
	}
}
