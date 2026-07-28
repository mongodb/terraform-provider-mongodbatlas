package organization2

import (
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TFModel struct {
	Name                 types.String `tfsdk:"name"`
	OrgID                types.String `tfsdk:"org_id"`
	ClientID             types.String `tfsdk:"client_id"`
	ClientSecret         types.String `tfsdk:"client_secret"`
	ClientSecretRotation types.Object `tfsdk:"client_secret_rotation"`
}

type TFClientSecretRotationModel struct {
	Interval        types.String `tfsdk:"interval"`
	NextRenewal     types.String `tfsdk:"next_renewal"`
	ExpiresAt       types.String `tfsdk:"expires_at"`
	CurrentSecretID types.String `tfsdk:"current_secret_id"`
	OldSecretID     types.String `tfsdk:"old_secret_id"`
	SecretVersion   types.Int64  `tfsdk:"secret_version"`
}

type orgState struct {
	nextRenewal      time.Time
	expiresAt        time.Time
	secretCreatedAt  time.Time
	name             string
	orgID            string
	clientID         string
	clientSecret     string
	interval         string
	currentSecretID  string
	oldSecretID      string
	secretVersion    int64
	hasRotationBlock bool
}

func RotationDue(now, nextRenewal, expiresAt time.Time) bool {
	return !now.Before(nextRenewal) || !now.Before(expiresAt)
}

func computeRenewalTimes(secretCreatedAt time.Time, interval time.Duration) (nextRenewal, expiresAt time.Time) {
	nextRenewal = secretCreatedAt.Add(interval)
	expiresAt = secretCreatedAt.Add(2 * interval)
	return nextRenewal, expiresAt
}

func parseInterval(interval string) (time.Duration, error) {
	return time.ParseDuration(interval)
}

func formatRFC3339(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}
