package serviceaccountjwt

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/mongodb/atlas-sdk-go/auth"

	"github.com/mongodb/terraform-provider-mongodbatlas/internal/config"
)

func (r *ES) terraformVersion() string {
	if r.EphemeralResourceData != nil {
		return r.EphemeralResourceData.TerraformVersion
	}
	return ""
}

// withUserAgentClient injects an HTTP client into the context via
// auth.HTTPClient so the atlas-sdk-go token exchange picks it up.
//
// Transport chain:
//
//	baseUserAgentTransport  (sets "terraform-provider-mongodbatlas/<ver> Terraform/<ver>")
//	  → UserAgentTransport  (appends "Name/service_account_jwt Operation/..." from context)
//	    → NetworkLoggingTransport  (logs method/URL/timing/status)
//	      → http.DefaultTransport
func (r *ES) withUserAgentClient(ctx context.Context) context.Context {
	networkLog := config.NewTransportWithNetworkLogging(http.DefaultTransport, logging.IsDebugOrHigher())
	uaTransport := &config.UserAgentTransport{Transport: networkLog, Enabled: true}
	client := &http.Client{
		Transport: &baseUserAgentTransport{
			base:      uaTransport,
			userAgent: config.UserAgent(r.terraformVersion()),
		},
	}
	return context.WithValue(ctx, auth.HTTPClient, client)
}

// baseUserAgentTransport sets the base User-Agent header on every outgoing
// request so that UserAgentTransport can append resource/operation extras.
type baseUserAgentTransport struct {
	base      http.RoundTripper
	userAgent string
}

func (t *baseUserAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", t.userAgent)
	return t.base.RoundTrip(req)
}
