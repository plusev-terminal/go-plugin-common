package tax

import (
	"fmt"
	"time"

	"github.com/plusev-terminal/go-plugin-common/utils"
)

// StartImportParams contains parameters for the startImport command.
type StartImportParams struct {
	Account       string     `json:"account,omitempty" mapstructure:"account"`
	From          *time.Time `json:"from,omitempty" mapstructure:"from"`
	To            *time.Time `json:"to,omitempty" mapstructure:"to"`
	LastFetchedTs int64      `json:"lastFetchedTs,omitempty" mapstructure:"lastFetchedTs"` // Unix millis
}

func (p StartImportParams) Validate() error { return nil }

// ResolveFrom returns the effective start time: From takes precedence, then LastFetchedTs.
func (p StartImportParams) ResolveFrom() time.Time {
	if p.From != nil {
		return *p.From
	}
	if p.LastFetchedTs > 0 {
		return time.UnixMilli(p.LastFetchedTs)
	}
	return time.Time{}
}

// ResolveTo returns the effective end time, or zero time if not set.
func (p StartImportParams) ResolveTo() time.Time {
	if p.To != nil {
		return *p.To
	}
	return time.Time{}
}

// StartImportParamsFromMap extracts StartImportParams from a validated map.
func StartImportParamsFromMap(data map[string]any) StartImportParams {
	return StartImportParams{
		Account:       utils.GetValue[string]("account", data),
		From:          utils.ExtractTime("from", data),
		To:            utils.ExtractTime("to", data),
		LastFetchedTs: int64(utils.GetValue[float64]("lastFetchedTs", data)),
	}
}

// GenerateReportParams contains parameters for the generateReport command.
type GenerateReportParams struct {
	TaxYear int `json:"taxYear" mapstructure:"taxYear"`
}

func (p GenerateReportParams) Validate() error {
	if p.TaxYear < 2009 || p.TaxYear > 2100 {
		return fmt.Errorf("taxYear must be between 2009 and 2100, got %d", p.TaxYear)
	}
	return nil
}

// GenerateReportParamsFromMap extracts GenerateReportParams from a validated map.
func GenerateReportParamsFromMap(data map[string]any) GenerateReportParams {
	return GenerateReportParams{
		TaxYear: utils.ExtractInt("taxYear", data),
	}
}
