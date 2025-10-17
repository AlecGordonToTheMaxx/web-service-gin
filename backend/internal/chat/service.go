package chat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"web-service-gin/backend/internal/album"

	openai "github.com/sashabaranov/go-openai"
)

// Service handles chat operations with OpenAI
type Service struct {
	client    *openai.Client
	albumRepo album.Repository
}

// NewService creates a new chat service
func NewService(albumRepo album.Repository) *Service {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		panic("OPENAI_API_KEY environment variable is required")
	}

	client := openai.NewClient(apiKey)

	return &Service{
		client:    client,
		albumRepo: albumRepo,
	}
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatRequest represents the incoming chat request
type ChatRequest struct {
	Messages []Message `json:"messages" binding:"required"`
}

// ChatResponse represents the chat response
type ChatResponse struct {
	Message     string              `json:"message"`
	ToolCalls   []openai.ToolCall   `json:"tool_calls,omitempty"`
	ToolResults []ToolResult        `json:"tool_results,omitempty"`
}

// ToolResult represents the result of a tool execution
type ToolResult struct {
	ToolCallID string `json:"tool_call_id"`
	Output     string `json:"output"`
}

// GetToolDefinitions returns the tool definitions for OpenAI
func (s *Service) GetToolDefinitions() []openai.Tool {
	return []openai.Tool{
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "get_albums",
				Description: "Get all albums, optionally filtered by search term",
				Parameters: json.RawMessage(`{
					"type": "object",
					"properties": {
						"search": {
							"type": "string",
							"description": "Optional search term to filter albums by title, artist"
						}
					}
				}`),
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "get_album_by_id",
				Description: "Get a specific album by its ID",
				Parameters: json.RawMessage(`{
					"type": "object",
					"properties": {
						"id": {
							"type": "number",
							"description": "The ID of the album to retrieve"
						}
					},
					"required": ["id"]
				}`),
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "create_album",
				Description: "Create a new album",
				Parameters: json.RawMessage(`{
					"type": "object",
					"properties": {
						"title": {
							"type": "string",
							"description": "Title of the album"
						},
						"artist": {
							"type": "string",
							"description": "Artist name"
						},
						"price": {
							"type": "number",
							"description": "Price of the album"
						}
					},
					"required": ["title", "artist", "price"]
				}`),
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "update_album",
				Description: "Update an existing album by ID",
				Parameters: json.RawMessage(`{
					"type": "object",
					"properties": {
						"id": {
							"type": "number",
							"description": "ID of the album to update"
						},
						"title": {
							"type": "string",
							"description": "Title of the album"
						},
						"artist": {
							"type": "string",
							"description": "Artist name"
						},
						"price": {
							"type": "number",
							"description": "Price of the album"
						}
					},
					"required": ["id"]
				}`),
			},
		},
		{
			Type: openai.ToolTypeFunction,
			Function: &openai.FunctionDefinition{
				Name:        "delete_album",
				Description: "Delete an album by ID",
				Parameters: json.RawMessage(`{
					"type": "object",
					"properties": {
						"id": {
							"type": "number",
							"description": "ID of the album to delete"
						}
					},
					"required": ["id"]
				}`),
			},
		},
	}
}

// ExecuteTool executes a tool function based on the tool name
func (s *Service) ExecuteTool(ctx context.Context, toolName string, argsJSON string) (string, error) {
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}

	switch toolName {
	case "get_albums":
		return s.getAlbums(ctx, args)
	case "get_album_by_id":
		return s.getAlbumByID(ctx, args)
	case "create_album":
		return s.createAlbum(ctx, args)
	case "update_album":
		return s.updateAlbum(ctx, args)
	case "delete_album":
		return s.deleteAlbum(ctx, args)
	default:
		return "", fmt.Errorf("unknown tool: %s", toolName)
	}
}

// Tool implementation functions
func (s *Service) getAlbums(ctx context.Context, args map[string]interface{}) (string, error) {
	albums, err := s.albumRepo.FindAll(ctx)
	if err != nil {
		return "", err
	}

	result := map[string]interface{}{
		"albums": albums,
		"count":  len(albums),
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func (s *Service) getAlbumByID(ctx context.Context, args map[string]interface{}) (string, error) {
	id, ok := args["id"].(float64)
	if !ok {
		return "", errors.New("id must be a number")
	}

	albumResult, err := s.albumRepo.FindByID(ctx, int(id))
	if err != nil {
		if errors.Is(err, album.ErrNotFound) {
			return fmt.Sprintf(`{"error": "Album with ID %d not found"}`, int(id)), nil
		}
		return "", err
	}

	jsonData, err := json.Marshal(map[string]interface{}{
		"album": albumResult,
	})
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func (s *Service) createAlbum(ctx context.Context, args map[string]interface{}) (string, error) {
	title, ok := args["title"].(string)
	if !ok {
		return "", errors.New("title must be a string")
	}

	artist, ok := args["artist"].(string)
	if !ok {
		return "", errors.New("artist must be a string")
	}

	price, ok := args["price"].(float64)
	if !ok {
		return "", errors.New("price must be a number")
	}

	newAlbum := &album.Album{
		Title:  title,
		Artist: artist,
		Price:  price,
	}

	if err := s.albumRepo.Create(ctx, newAlbum); err != nil {
		return "", err
	}

	result := map[string]interface{}{
		"album":   newAlbum,
		"message": "Album created successfully",
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func (s *Service) updateAlbum(ctx context.Context, args map[string]interface{}) (string, error) {
	id, ok := args["id"].(float64)
	if !ok {
		return "", errors.New("id must be a number")
	}

	existingAlbum, err := s.albumRepo.FindByID(ctx, int(id))
	if err != nil {
		if errors.Is(err, album.ErrNotFound) {
			return fmt.Sprintf(`{"error": "Album with ID %d not found"}`, int(id)), nil
		}
		return "", err
	}

	// Update fields if provided
	if title, ok := args["title"].(string); ok {
		existingAlbum.Title = title
	}
	if artist, ok := args["artist"].(string); ok {
		existingAlbum.Artist = artist
	}
	if price, ok := args["price"].(float64); ok {
		existingAlbum.Price = price
	}

	if err := s.albumRepo.Update(ctx, existingAlbum); err != nil {
		return "", err
	}

	result := map[string]interface{}{
		"album":   existingAlbum,
		"message": "Album updated successfully",
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func (s *Service) deleteAlbum(ctx context.Context, args map[string]interface{}) (string, error) {
	id, ok := args["id"].(float64)
	if !ok {
		return "", errors.New("id must be a number")
	}

	if err := s.albumRepo.Delete(ctx, int(id)); err != nil {
		if errors.Is(err, album.ErrNotFound) {
			return fmt.Sprintf(`{"error": "Album with ID %d not found"}`, int(id)), nil
		}
		return "", err
	}

	result := map[string]interface{}{
		"message": fmt.Sprintf("Album %d deleted successfully", int(id)),
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

// Chat handles the main chat interaction
func (s *Service) Chat(ctx context.Context, messages []Message) (*ChatResponse, error) {
	// Convert messages to OpenAI format
	chatMessages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: `You are an intelligent album management assistant. You can help users manage their album collection by:
- Viewing and searching albums
- Creating new albums
- Updating existing albums
- Deleting albums

Always be helpful and provide clear explanations of what actions you're taking. When presenting data, format it in a user-friendly way. If asked to create or update albums, ask for clarification on any required fields that are missing (title, artist, price are required).`,
		},
	}

	// Add user messages
	for _, msg := range messages {
		chatMessages = append(chatMessages, openai.ChatCompletionMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	// Make initial API call
	resp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    openai.GPT4oMini,
		Messages: chatMessages,
		Tools:    s.GetToolDefinitions(),
	})

	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	message := resp.Choices[0].Message

	// Handle tool calls if present
	if len(message.ToolCalls) > 0 {
		toolResults := make([]ToolResult, 0, len(message.ToolCalls))

		// Execute each tool call
		for _, toolCall := range message.ToolCalls {
			result, err := s.ExecuteTool(ctx, toolCall.Function.Name, toolCall.Function.Arguments)
			if err != nil {
				result = fmt.Sprintf(`{"error": "%s"}`, err.Error())
			}

			toolResults = append(toolResults, ToolResult{
				ToolCallID: toolCall.ID,
				Output:     result,
			})
		}

		// Add assistant's message with tool calls
		chatMessages = append(chatMessages, message)

		// Add tool results as messages
		for _, toolResult := range toolResults {
			chatMessages = append(chatMessages, openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				Content:    toolResult.Output,
				ToolCallID: toolResult.ToolCallID,
			})
		}

		// Make second API call with tool results
		finalResp, err := s.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
			Model:    openai.GPT4oMini,
			Messages: chatMessages,
		})

		if err != nil {
			return nil, fmt.Errorf("OpenAI API error on second call: %w", err)
		}

		return &ChatResponse{
			Message:     finalResp.Choices[0].Message.Content,
			ToolCalls:   message.ToolCalls,
			ToolResults: toolResults,
		}, nil
	}

	// No tool calls, return direct response
	return &ChatResponse{
		Message:     message.Content,
		ToolCalls:   nil,
		ToolResults: nil,
	}, nil
}
