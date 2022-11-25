package google

import (
	"context"
)

type AndroidPublisherClient interface {
	GetProductPurchase(ctx context.Context, packageName string, productId string, token string) (*ProductPurchase, error)
}
