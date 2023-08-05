// Package explain contains a conversation.Turn implementation for the explain intent.
package explain

import (
	"context"
	"errors"

	"github.com/github/copilot-api/pkg/chat/conversation"
	"github.com/github/copilot-api/pkg/chat/tokens"
	"github.com/github/copilot-api/pkg/chat/turn/base"
	"github.com/github/copilot-api/pkg/chat/turn/prompt"
)

const (
	header = `Help the user understand the code below.
Format the response as markdown. Surround code with backticks ` + "`like_this()`" + `.
Add markdown links to major external sites if they are certain to exist. Do not guess links.
Only consider the content below. Ignore any previous knowledge of this file.
Be concise — the explanation should be a single paragraph. The user will ask additional follow-up questions.`

	footer = "Now, provide a brief explanation for this code as requested."
)

// Turn represents an explain turn in a chat thread.
type Turn struct {
	*base.Turn

	pb *prompt.Builder
}

// NewTurn creates a new explain turn.
func NewTurn(bt *base.Turn) (*Turn, error) {
	inputRefs := bt.InputReferences()
	if len(inputRefs) == 0 {
		return nil, errors.New("no input references")
	}

	tc := tokens.NewCounter(bt.Tokenizer(), base.ReservedResponseTokens, bt.MaxTokens)

	pb := prompt.
		NewBuilder(tc).
		SafetyMessage().
		PreviousTurns(bt.PrevTurns())

	params := prompt.Message{
		Header:     header,
		References: inputRefs,
		Footer:     footer,
		Role:       conversation.RoleSystem,
	}
	pb.AddReferenceMessage(params)

	userMessage := bt.UserMessage()
	pb.AddMessage(userMessage.Role, userMessage.Content)

	return &Turn{
		Turn: bt,
		pb:   pb,
	}, nil
}

// Prompt generates a prompt for a Turn.
func (t *Turn) Prompt(ctx context.Context) ([]*conversation.Message, error) {
	m, err := t.pb.Build()
	if err != nil {
		return nil, err
	}

	return m, nil
}
