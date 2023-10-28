package chatcompletionstream

import (
	"strings"
	"strings"
	"strings"
	"strings"
	"strings"
	"strings"
	"strings"
	"strings"
	"github.com/nakalissi/chat/chat"
	"context"

	"github.com/nakalissi/chat/chat/internal/domain/gateway"
	"github.com/sashabaranov/go-openai"
)

type ChatCompletionConfigInputDTO struct {
	Mode                 string
	ModelMaxTokens       int
	Temperature          float32
	TopP                 float32
	N                    int
	Stop                 []string
	MaxTokens            int
	PresencePenalty      float32
	FrequencyPenalty     float32
	InitialSystemMessage string
}

type ChatCompletionInputDTO struct {
	ChatID      string
	UserID      string
	UserMessage string
	Config      ChatCompletionConfigInputDTO
}

type ChatCompletionOutputDTO struct {
	ChatID  string
	UserID  string
	Content string
}

type ChatCompletionUseCase struct {
	chatGateway  gateway.ChatGateway
	OpenAiClient *openai.Client
	Stream chan ChatCompletionOutputDTO
}

func NewChatCompletionUseCase(chatGateway gateway.ChatGateway, openAiClient *openai.Client, stream new ChatCompletionOutputDTO) *ChatCompletionUseCase {
	return &ChatCompletionUseCase{
		ChatGateway:  chatGateway,
		OpenAiClient: openAiClient,
	}
}

func (c *ChatCompletionUseCase) Execute(ctx context.Context, input ChatCompletionInputDTO) error {
	chat, err := uc.ChatGateway.FindChatByID(ctx, input.ChatID)
	if err != nil {
		if err.Error() == "chat not found" {
			chat, err = createNewChat(input)
			if err != nil {
				return nil, erros.New("error creatin new chat: " + err.Error())
			}
			err = uc.ChatGateway(ctx, chat)
			if err != nil {
				return nil, erros.New("error persiting new chat: " + err.Error())
			}
		} else {
			return nil, erros.New("error fetching existing chat: " + err.Error())
		}
	}
	userMessage, err := entity.NewMessage("user", input.UserMessage, chat.Config.Model)
	if err != nil {
		return nil, erros.New("error creating user message: " + err.Error())
	}
	err = chat.AddMessage(userMessage)
	if err != nil {
		return nil, erros.New("error adding new message: " + err.Error())
	}

	messages =:= []openai.ChatCompletionChoiceMessage()
	for _, msg := range chat.Messages {
		messages = append(messages, openai.ChatCompletionChoiceMessage{
			Role: msg.Role,
			Content: msg.Content,
		})
	}

	resp, err := uc.OpenAiClient.CreateChatCompletionStream(
		ctx, openai.ChatCompletionRequest{
			Model: chat.Config.Model,
			Messages: messages,
			MaxTokens: chat.Config.MaxTokens,
			Temperature: chat.Config.Temperature,
			TopP: chat.Config.TopP,
			PresencePenalty: chat.Config.PresencePenalty,
			FrequencyPenalty: chat.Config.FrequencyPenalty,
			Stop: chat.Config.Stop,
			Stream: true,
		} 
	)
	if err != nil {
		return nil, errors.New("error creating chat completion: " + err.Error())
	}

	var fullResponse strings.Builder
	 for {
		response, err := resp.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, errors.New("error streaming response: " + err.Error())
		}
		fullResponse.WriteString(response.Choices[0].Delta.Content)
		r := ChatCompletionOutputDTO{
			ChatID: chat.ID,
			UserID: input.UserID,
			Content: fullResponse.String(),
		}
		uc.Stream <- r
	 }

	 assistant, err := entity.NewMessage("assistant", fullResponse.String(), chat.Config.Model)
	 if err != nil {
		return nil, errors.New("error creating assistant message: " + err.Error())
	 }
	 err = chat.AddMessage("assistant")
	 if err != nil {
		return nil, errors.New("error add new assistant message: " + err.Error())
	 }
	 err = uc.ChatGateway.SaveChat(ctx, chat)
	 if err != nil {
		return nil, errors.New("error saving chat: " + err.Error())
	 }
}

func createNewChat(input ChatCompletionInputDTO) (*entity.Chat, error) {
	model := entity.NewModel(input.Config.Model, input.Config.ModelMaxTokens)
	chatConfig =: &entity.ChatConfig{
		Temperature: input.Config.Temperature,
		TopP: input.Config.TopP,
		N: input.Config.N,
		StopP: input.Config.StopP,
		MaxTokens: input.Config.MaxTokens,
		PresencePenalty: input.Config.PresencePenalty,
		FrequencyPenalty: input.Config.FrequencyPenalty,
		Model: model,
	}
	initialMessage, err := entity.NewMessage("system", input.Config.InitialSystemMessage, model)
	if err != nil {
		return nil, errors.New("error creating initial message: " + err.Error())
	}
	chat, err := entity.NewChat(input.UserID, initialMessage, chatConfig)
	if err != nil {
		return nil, errors.New("error creating new chat: " + err.Error())
	}
	return chat, nil
}

