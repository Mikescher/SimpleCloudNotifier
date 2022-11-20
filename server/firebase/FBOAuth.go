package firebase

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gogs.mikescher.com/BlackForestBytes/goext/langext"
	"gogs.mikescher.com/BlackForestBytes/goext/timeext"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type FBOAuth2 struct {
	client *http.Client

	scopes       []string
	tokenURL     string
	privateKeyID string
	clientMail   string

	currToken   *string
	tokenExpiry *time.Time
	privateKey  *rsa.PrivateKey
}

func NewAuth(tokenURL string, privKeyID string, cmail string, pemstr string) (*FBOAuth2, error) {

	pkey, err := decodePemKey(pemstr)
	if err != nil {
		return nil, err
	}

	return &FBOAuth2{
		client:       &http.Client{Timeout: 3 * time.Second},
		tokenURL:     tokenURL,
		privateKey:   pkey,
		privateKeyID: privKeyID,
		clientMail:   cmail,
		scopes: []string{
			"https://www.googleapis.com/auth/cloud-platform",
			"https://www.googleapis.com/auth/datastore",
			"https://www.googleapis.com/auth/devstorage.full_control",
			"https://www.googleapis.com/auth/firebase",
			"https://www.googleapis.com/auth/identitytoolkit",
			"https://www.googleapis.com/auth/userinfo.email",
		},
	}, nil
}

func decodePemKey(pemstr string) (*rsa.PrivateKey, error) {
	var raw []byte

	block, _ := pem.Decode([]byte(pemstr))

	if block != nil {
		raw = block.Bytes
	} else {
		raw = []byte(pemstr)
	}

	pkey8, err1 := x509.ParsePKCS8PrivateKey(raw)
	if err1 == nil {
		privkey, ok := pkey8.(*rsa.PrivateKey)
		if !ok {
			return nil, errors.New("private key is invalid")
		}
		return privkey, nil
	}

	pkey1, err2 := x509.ParsePKCS1PrivateKey(raw)
	if err2 == nil {
		return pkey1, nil
	}

	return nil, errors.New(fmt.Sprintf("failed to parse private-key: [ %v | %v ]", err1, err2))
}

func (a *FBOAuth2) Token(ctx context.Context) (string, error) {
	if a.currToken == nil || a.tokenExpiry == nil || a.tokenExpiry.Before(time.Now()) {
		err := a.Refresh(ctx)
		if err != nil {
			return "", err
		}
	}

	return *a.currToken, nil
}

func (a *FBOAuth2) Refresh(ctx context.Context) error {

	assertion, err := a.encodeAssertion(a.privateKey)
	if err != nil {
		return err
	}

	body := url.Values{
		"assertion":  []string{assertion},
		"grant_type": []string{"urn:ietf:params:oauth:grant-type:jwt-bearer"},
	}.Encode()

	req, err := http.NewRequestWithContext(ctx, "POST", a.tokenURL, strings.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	reqNow := time.Now()

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if bstr, err := io.ReadAll(resp.Body); err == nil {
			return errors.New(fmt.Sprintf("Auth-Request returned %d: %s", resp.StatusCode, string(bstr)))
		} else {
			return errors.New(fmt.Sprintf("Auth-Request returned %d", resp.StatusCode))
		}
	}

	respBodyBin, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var respBody struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(respBodyBin, &respBody); err != nil {
		return err
	}

	a.currToken = langext.Ptr(respBody.AccessToken)
	a.tokenExpiry = langext.Ptr(reqNow.Add(timeext.FromSeconds(respBody.ExpiresIn)))

	return nil
}

func (a *FBOAuth2) encodeAssertion(key *rsa.PrivateKey) (string, error) {
	headBin, err := json.Marshal(gin.H{"alg": "RS256", "typ": "JWT", "kid": a.privateKeyID})
	if err != nil {
		return "", err
	}
	head := base64.RawURLEncoding.EncodeToString(headBin)

	now := time.Now().Add(-10 * time.Second) // jwt hack against unsynced clocks

	claimBin, err := json.Marshal(gin.H{"iss": a.clientMail, "scope": strings.Join(a.scopes, " "), "aud": a.tokenURL, "exp": now.Add(time.Hour).Unix(), "iat": now.Unix()})
	if err != nil {
		return "", err
	}
	claim := base64.RawURLEncoding.EncodeToString(claimBin)

	checksum := sha256.New()
	checksum.Write([]byte(head + "." + claim))
	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, checksum.Sum(nil))
	if err != nil {
		return "", err
	}

	return head + "." + claim + "." + base64.RawURLEncoding.EncodeToString(sig), nil
}
