package google

import (
	"context"
	_ "embed"
	"fmt"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"strings"
	"time"
)

type DummyGoogleAPIClient struct{}

func NewDummy() AndroidPublisherClient {
	return &DummyGoogleAPIClient{}
}

func (d DummyGoogleAPIClient) GetProductPurchase(ctx context.Context, packageName string, productId string, token string) (*ProductPurchase, error) {
	if strings.HasPrefix(token, "PURCHASED:") {
		return &ProductPurchase{
			Kind:                        "",
			PurchaseTimeMillis:          fmt.Sprintf("%d", time.Date(2000, 1, 1, 12, 0, 0, 0, time.UTC).UnixMilli()),
			PurchaseState:               langext.Ptr(PurchaseStatePurchased),
			ConsumptionState:            ConsumptionStateConsumed,
			DeveloperPayload:            "{}",
			OrderId:                     "000",
			PurchaseType:                nil,
			AcknowledgementState:        AcknowledgementStateAcknowledged,
			PurchaseToken:               nil,
			ProductId:                   langext.Ptr("1234-5678"),
			Quantity:                    nil,
			ObfuscatedExternalAccountId: "000",
			ObfuscatedExternalProfileId: "000",
			RegionCode:                  "DE",
		}, nil
	}
	return nil, nil // = purchase not found
}
