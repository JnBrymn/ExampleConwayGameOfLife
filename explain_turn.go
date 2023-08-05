// Package explain contains a conversation.Turn implementation for the explain intent.
package explain

import (
	"context"
	"errors"

	"github.com/github/copilot-api/pkg/chat/conversation"
	"github.com/github/copilot-api/pkg/chat/tokens"
	"github.com/github/copilot-api/pkg/chat/turn/base"
	"github.com/github/copilot-api/pkg/chat/turn/prompt"
	"github.com/github/copilot-api/pkg/chat/turn/reference"
)

// Turn represents an explain turn in a chat thread.
type Turn struct {
	*base.Turn

	pb *prompt.Builder
	finished bool
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

	refs := bt.References()
	if len(refs) == 1 {
		ref := refs[0]
		pb = pb.AddMessage(systemMessage.Role,
			// ShortDescriptor() returns things like "snippet" or "symbol"
			"Help the user understand the " + ref.ShortDescriptor() + " below.\n" +
			ref.Identifier() +  // things like https://github.com/JnBrymn/ExampleConwayGameOfLife/blob/main/GameOfLife.java#L86-L87
			"\n" + ref.Text() + "\n" 
			"\nThis was found in a greater context of:\n" +
			reference.GreaterContext(ref) // see the last example here https://github.com/github/copilot-api/issues/565#issuecomment-1666237978
		)	
	} else {
		pb = pb.AddMessage(systemMessage.Role,
			// ShortDescriptor() here returns things like "code" or "snippets" or "symbols"
			"Help the user understand the " + refs.ShortDescriptor() + " below.\n" +
			refs.IdentifiersAndText() // see comments above
			// ignoring greater context because we don't know what that reasonably looks like for multiple references
		)	
	}
	pb = pb.AddMessage(userMessage.Role, bt.UserMessage().Content)

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

// ProcessModelResponse processes the response from the model and adds it to the turn.
func (t *Turn) ProcessModelResponse(modelResponse string) {
	// Nothing much to do here
	t.AddMessage(conversation.NewAssistantMessage(modelResponse))
	t.finished = true
}

// Finished returns whether a Turn is finished.
func (t *Turn) Finished() bool {
	return t.finished
}
