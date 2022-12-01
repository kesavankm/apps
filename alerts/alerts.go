package alerts

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
)

type AlertConfig struct {
	ClientType string
	ApiToken   string
}

// Alert - supported operations
type AlertIntf interface {
	PostMessage(channelID, message string) (err error)
}

func NewAlertClient(ctx context.Context, ac AlertConfig) (intf AlertIntf, err error) {
	switch {
	case strings.EqualFold(ac.ClientType, "slack"):
		return NewSlackClient(ctx, ac.ApiToken)
	default:
		errMsg := fmt.Sprintf("Unexpected clientType %s\n", ac.ClientType)
		log.Printf("err: %s", errMsg)
		return nil, errors.New(errMsg)
	}
	return nil, nil
}
