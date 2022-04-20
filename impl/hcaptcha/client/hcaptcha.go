package client

import (
	"context"
	"encoding/json"
	"github.com/reaganiwadha/arah/domain"
	lr "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	siteVerifyEndpoint = "https://hcaptcha.com/siteverify"
)

type hCaptchaVerifyResponse struct {
	Success     bool      `json:"success"`
	ChallengeTs time.Time `json:"challenge_ts"`
}

type hCaptchaClient struct {
	serverSecret string
}

func NewHCaptchaClient(serverSecret string) domain.HCaptchaClient {
	return &hCaptchaClient{serverSecret: serverSecret}
}

func (h hCaptchaClient) Verify(ctx context.Context, clientResponse string) (success bool, err error) {
	payload := url.Values{}
	payload.Set("response", clientResponse)
	payload.Set("secret", h.serverSecret)

	encoded := payload.Encode()

	r, err := http.NewRequestWithContext(ctx, http.MethodPost, siteVerifyEndpoint, strings.NewReader(encoded))
	if err != nil {
		lr.Fatal(err)
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(encoded)))

	res, err := http.DefaultClient.Do(r)

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return
	}

	var responseJson hCaptchaVerifyResponse

	err = json.Unmarshal(body, &responseJson)

	success = responseJson.Success

	return
}
