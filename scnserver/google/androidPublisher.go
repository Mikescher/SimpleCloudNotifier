package google

import (
	scn "blackforestbytes.com/simplecloudnotifier"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strings"
	"time"
)

// https://developers.google.com/android-publisher/api-ref/rest/v3/purchases.products/get
// https://developers.google.com/android-publisher/api-ref/rest/v3/purchases.products#ProductPurchase

type AndroidPublisher struct {
	client  http.Client
	auth    *GoogleOAuth2
	baseURL string
}

func NewAndroidPublisherAPI(conf scn.Config) (AndroidPublisherClient, error) {

	pkey := strings.ReplaceAll(conf.GoogleAPIPrivateKey, "\\n", "\n")

	googauth, err := NewAuth(conf.GoogleAPITokenURI, conf.GoogleAPIPrivKeyID, conf.GoogleAPIClientMail, pkey)
	if err != nil {
		return nil, err
	}

	return &AndroidPublisher{
		client:  http.Client{Timeout: 5 * time.Second},
		auth:    googauth,
		baseURL: "https://androidpublisher.googleapis.com/androidpublisher",
	}, nil
}

type PurchaseType int //@enum:type

const (
	PurchaseTypeTest     PurchaseType = 0 // i.e. purchased from a license testing account
	PurchaseTypePromo    PurchaseType = 1 // i.e. purchased using a promo code
	PurchaseTypeRewarded PurchaseType = 2 // i.e. from watching a video ad instead of paying
)

type ConsumptionState int //@enum:type

const (
	ConsumptionStateYetToBeConsumed ConsumptionState = 0
	ConsumptionStateConsumed        ConsumptionState = 1
)

type PurchaseState int //@enum:type

const (
	PurchaseStatePurchased PurchaseState = 0
	PurchaseStateCanceled  PurchaseState = 1
	PurchaseStatePending   PurchaseState = 2
)

type AcknowledgementState int //@enum:type

const (
	AcknowledgementStateYetToBeAcknowledged AcknowledgementState = 0
	AcknowledgementStateAcknowledged        AcknowledgementState = 1
)

type ProductPurchase struct {
	Kind                        string               `json:"kind"`
	PurchaseTimeMillis          string               `json:"purchaseTimeMillis"`
	PurchaseState               *PurchaseState       `json:"purchaseState"`
	ConsumptionState            ConsumptionState     `json:"consumptionState"`
	DeveloperPayload            string               `json:"developerPayload"`
	OrderId                     string               `json:"orderId"`
	PurchaseType                *PurchaseType        `json:"purchaseType"`
	AcknowledgementState        AcknowledgementState `json:"acknowledgementState"`
	PurchaseToken               *string              `json:"purchaseToken"`
	ProductId                   *string              `json:"productId"`
	Quantity                    *int                 `json:"quantity"`
	ObfuscatedExternalAccountId string               `json:"obfuscatedExternalAccountId"`
	ObfuscatedExternalProfileId string               `json:"obfuscatedExternalProfileId"`
	RegionCode                  string               `json:"regionCode"`
}

type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (ap AndroidPublisher) GetProductPurchase(ctx context.Context, packageName string, productId string, token string) (*ProductPurchase, error) {

	uri := fmt.Sprintf("%s/v3/applications/%s/purchases/products/%s/tokens/%s", ap.baseURL, packageName, productId, token)

	request, err := http.NewRequestWithContext(ctx, "GET", uri, nil)
	if err != nil {
		return nil, err
	}

	tok, err := ap.auth.Token(ctx)
	if err != nil {
		log.Err(err).Msg("Refreshing FB token failed")
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+tok)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Accept", "application/json")

	response, err := ap.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() { _ = response.Body.Close() }()

	respBodyBin, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode == 400 {

		var errBody struct {
			Error apiError `json:"error"`
		}
		if err := json.Unmarshal(respBodyBin, &errBody); err != nil {
			return nil, err
		}
		if errBody.Error.Code == 400 {
			return nil, nil // probably token not found
		}
	}

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		if bstr, err := io.ReadAll(response.Body); err == nil {
			return nil, errors.New(fmt.Sprintf("GetProducts-Request returned %d: %s", response.StatusCode, string(bstr)))
		} else {
			return nil, errors.New(fmt.Sprintf("GetProducts-Request returned %d", response.StatusCode))
		}
	}

	var respBody ProductPurchase
	if err := json.Unmarshal(respBodyBin, &respBody); err != nil {
		return nil, err
	}

	if respBody.Kind != "androidpublisher#productPurchase" {
		return nil, errors.New(fmt.Sprintf("Invalid ProductPurchase.kind: '%s'", respBody.Kind))
	}

	return &respBody, nil
}
