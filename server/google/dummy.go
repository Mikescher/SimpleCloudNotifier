package google

import (
	"context"
	_ "embed"
)

type DummyGoogleAPIClient struct{}

func NewDummy() AndroidPublisherClient {
	return &DummyGoogleAPIClient{}
}

func (d DummyGoogleAPIClient) GetProductPurchase(ctx context.Context, packageName string, productId string, token string) (*ProductPurchase, error) {
	return nil, nil // = purchase not found
}
