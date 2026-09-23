package utils

import (
	"context"
	"encoding/json"
	"fmt"
)

// EmbedCreatedBy extracts the Cognito sub from the request context
// and injects it into the provided JSON payload as "created_by", returning the new payload.
func EmbedCreatedBy(ctx context.Context, originalPayload []byte) ([]byte, error) {
	userSub, ok := ctx.Value("user_sub").(string)
	if !ok || userSub == "" {
		return nil, fmt.Errorf("unauthorized: missing user identity in context")
	}

	// Embed the Cognito sub into the JSON payload
	var payloadMap map[string]interface{}
	if len(originalPayload) > 0 {
		if err := json.Unmarshal(originalPayload, &payloadMap); err != nil {
			return nil, fmt.Errorf("invalid original JSON payload")
		}
	} else {
		payloadMap = make(map[string]interface{})
	}

	payloadMap["created_by"] = userSub

	enriched, err := json.Marshal(payloadMap)
	return enriched, err
}
