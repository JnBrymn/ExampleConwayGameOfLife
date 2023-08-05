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
	finishedWithAdditionalReferenceRetrieval bool
  finishedWithTurn bool
}

func (t *Turn) preamble() *prompt.Builder {
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
}

// NewTurn creates a new explain turn.
func NewTurn(bt *base.Turn) (*Turn, error) {
	
	t := &Turn{Turn: bt}
  
  pb = t.preamble()
 
  relatedRefs := turn.PrevTurns().References()
  relatedRefs = append(relatedRefs,turn.NearbyContext())
  
  pb = pb.AddMessage(
    systemMessage.Role,
    "The user's specific request is \"" + 
    turn.UserMessage().Content + "\"\n\n" +
    "But before we do, we must pick out any other relevant references to review. Which of these references are likely to be very relevant to the user's request?\n\n" +
    "Options:" +
    relatedRefs.OptionMenu() + "\n\n" +
    "Select the references that are likely to be useful."
  )

  t.pb = pb
  return t, nil
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
	if !t.finishedWithAdditionalReferenceRetrieval {
		additionalReferences := t.parseAdditionalReferences(modelResponse)
		additionalReferences.hydrate(t.skills)
		
		pb = t.preamble() // This recreated the top of the prompt just like before
		
		pb = pb.AddMessage(systemMessage.Role,
			"The " + refs.ShortDescriptor() + " below will also likely be helpful." +
			additionalReferences.IdentifiersAndText()
		)	
		pb = pb.AddMessage(userMessage.Role,
			turn.UserMessage().Content
		)
		t.pb = pb
		
		t.finishedWithAdditionalReferenceRetrieval = true
	}
		t.AddMessage(conversation.NewAssistantMessage(modelResponse))
		t.finishedWithTurn = true
	}
}

func (t *Turn) parseAdditionalReferences(modelResponse string) reference.Group {
  ... TODO ...
}

// Finished returns whether a Turn is finished.
func (t *Turn) Finished() bool {
	return t.finishedWithTurn
}
