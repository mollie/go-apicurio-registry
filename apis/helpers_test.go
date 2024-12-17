package apis_test

import (
	"context"
	"fmt"
	"github.com/mollie/go-apicurio-registry/apis"
	"github.com/mollie/go-apicurio-registry/models"
	"time"
)

func generateArtifactForTest(ctx context.Context, artifactsAPI *apis.ArtifactsAPI) (string, error) {
	// Helper to generate unique artifact IDs
	generateArtifactID := func(prefix string) string {
		return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
	}

	newArtifactID := generateArtifactID("test-artifact")

	artifact := models.CreateArtifactRequest{
		ArtifactID:   newArtifactID,
		ArtifactType: models.Json,
		Name:         newArtifactID,
		FirstVersion: models.CreateVersionRequest{
			Version: "1.0.0",
			Content: models.CreateContentRequest{
				Content:     stubArtifactContent,
				ContentType: "application/json",
			},
			IsDraft: true,
		},
	}
	createParams := &models.CreateArtifactParams{
		IfExists: models.IfExistsFail,
	}
	_, err := artifactsAPI.CreateArtifact(ctx, stubGroupId, artifact, createParams)
	if err != nil {
		return "", err
	}
	return newArtifactID, nil
}

func generateGroupForTest(ctx context.Context, groupAPI *apis.GroupAPI) (*models.GroupInfo, string, error) {
	newGroupID := generateRandomName("test-group")
	resp, err := groupAPI.CreateGroup(ctx, newGroupID, stubDescription, stubLabels)
	return resp, newGroupID, err
}

func deleteGroupAfterTest(ctx context.Context, groupAPI *apis.GroupAPI, groupID string) error {
	return groupAPI.DeleteGroup(ctx, groupID)
}

func generateRandomName(prefix string) string {
	return fmt.Sprintf("%s-%d", prefix, time.Now().UnixNano())
}
