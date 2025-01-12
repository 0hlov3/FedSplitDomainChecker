package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"

	"github.com/0hlov3/FedSplitDomainChecker/internal/logger"
	"github.com/spf13/cobra"
)

type WebfingerResponse struct {
	Subject string   `json:"subject"`
	Aliases []string `json:"aliases"`
	Links   []struct {
		Rel  string `json:"rel"`
		Type string `json:"type"`
		Href string `json:"href"`
	} `json:"links"`
}

var checkSplitDomainCmd = &cobra.Command{
	Use:   "checkSplitDomain",
	Short: "Check a split-domain Fediverse deployment",
	Long:  `Validates the split-domain deployment by verifying webfinger configurations and responses.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("❓ Starting split-domain check")
		logger.ZapLog.Debug("Starting split-domain check")

		if err := validateEndpoints(); err != nil {
			printFailureStep("Split-domain check", err)
			logger.ZapLog.Error("Split-domain check failed", zap.Error(err))
			return
		}

		printSuccessStep("Split-domain check")
		logger.ZapLog.Debug("Split-domain check completed successfully")
	},
}

func init() {
	rootCmd.AddCommand(checkSplitDomainCmd)
}

func validateEndpoints() error {
	endpoints := []struct {
		name string
		url  string
	}{
		{"Host-Meta", fmt.Sprintf("https://%s/.well-known/host-meta", AccountDomain)},
		{"Nodeinfo", fmt.Sprintf("https://%s/.well-known/nodeinfo", AccountDomain)},
		{"Webfinger", fmt.Sprintf("https://%s/.well-known/webfinger?resource=acct:%s", AccountDomain, Account)},
	}

	for _, endpoint := range endpoints {
		printSuccessStep(fmt.Sprintf("Checking %s endpoint", endpoint.name))
		logger.ZapLog.Debug("Checking endpoint", zap.String("name", endpoint.name), zap.String("url", endpoint.url))

		resp, err := makeRequest(endpoint.url)
		if err != nil {
			return fmt.Errorf("%s validation failed: %w", endpoint.name, err)
		}
		defer resp.Body.Close()
		resp, err = validateRedirectLocation(resp, HostDomain, endpoint.name)
		if err != nil {
			return fmt.Errorf("%s redirect validation failed: %w", endpoint.name, err)
		}

		if endpoint.name == "Webfinger" {
			if err := validateWebfinger(resp, HostDomain, Account); err != nil {
				return fmt.Errorf("Webfinger validation failed: %w", err)
			}
		}

		printSuccessStep(fmt.Sprintf("%s endpoint validation", endpoint.name))
	}

	return nil
}

func makeRequest(url string) (*http.Response, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to query URL: %w", err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("404 when query URL: %s", url)
	}

	return resp, nil
}

func validateWebfinger(resp *http.Response, hostDomain, account string) error {
	if !strings.Contains(resp.Header.Get("Content-Type"), "application/jrd+json") {
		return fmt.Errorf("unexpected content type: %s", resp.Header.Get("Content-Type"))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	var webfingerResp WebfingerResponse
	if err := json.Unmarshal(body, &webfingerResp); err != nil {
		return fmt.Errorf("failed to parse webfinger response JSON: %w", err)
	}

	expectedSubject := fmt.Sprintf("acct:%s", account)
	if webfingerResp.Subject != expectedSubject {
		return fmt.Errorf("subject mismatch: got %s, expected %s", webfingerResp.Subject, expectedSubject)
	}

	fmt.Printf("✅ Webfinger Subject: %s\n", webfingerResp.Subject)
	return validateSelfLink(webfingerResp, hostDomain, extractUsername(account))
}

func validateRedirectLocation(resp *http.Response, hostDomain, endpointName string) (*http.Response, error) {
	if resp.StatusCode == http.StatusMovedPermanently {
		redirectURL := resp.Header.Get("Location")
		if redirectURL == "" {
			return nil, fmt.Errorf("%s endpoint: redirect received, but no location header provided", endpointName)
		}
		if !strings.HasPrefix(redirectURL, fmt.Sprintf("https://%s", hostDomain)) {
			return nil, fmt.Errorf("%s endpoint: redirect location mismatch: got %s, expected to start with https://%s", endpointName, redirectURL, hostDomain)
		}
		fmt.Printf("✅ Redirect recived: %s\n", redirectURL)

		logger.ZapLog.Debug("Redirect location validated", zap.String("endpoint", endpointName), zap.String("location", redirectURL))
		resp.Body.Close()
		return http.Get(redirectURL)
	}

	return nil, errors.New("validateRedirectLocation error")
}

func validateSelfLink(resp WebfingerResponse, hostDomain, username string) error {
	expectedSelfLink := fmt.Sprintf("https://%s/users/%s", hostDomain, username)
	for _, link := range resp.Links {
		if link.Rel == "self" && link.Type == "application/activity+json" {
			if link.Href == expectedSelfLink {
				logger.ZapLog.Debug("Self-link validation passed", zap.String("selfLink", link.Href))
				return nil
			}
			return fmt.Errorf("unexpected self-link: got %s, expected %s", link.Href, expectedSelfLink)
		}
	}

	return errors.New("self-link not found in webfinger response")
}

func extractUsername(account string) string {
	parts := strings.Split(account, "@")
	if len(parts) == 2 {
		return parts[0]
	}
	return ""
}

func printSuccessStep(step string) {
	fmt.Printf("✅ %s: success\n", step)
}

func printFailureStep(step string, err error) {
	fmt.Printf("❌ %s: failed (%v)\n", step, err)
}
