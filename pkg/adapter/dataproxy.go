package adapter

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"regexp"
)

const (
	MaxClinicalPayloadBytes = 64 * 1024
	AuthorizedJurisdiction  = "AR-SOUTHAMERICA-WEST1"
)

var (
	ErrInvalidHMACSecret     = errors.New("DATAPROXY_VIOLATION: HMAC secret key must be at least 32 bytes")
	ErrTrailingOrInvalidJSON = errors.New("DATAPROXY_VIOLATION: malformed JSON, unknown fields, or trailing smuggling data")
	ErrJurisdictionFence     = errors.New("SOVEREIGNTY_VIOLATION: cross-border data transfer prohibited under Ley 25.326")
	reDNI   = regexp.MustCompile(`\b\d{1,2}\.?\d{3}\.?\d{3}\b`)
	reCUIL  = regexp.MustCompile(`\b(20|23|24|27|30|33|34)-?\d{8}-?\d\b`)
	reEmail = regexp.MustCompile(`\b[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}\b`)
	rePhone = regexp.MustCompile(`\b(\+?54\s?9?\s?)?(\d{2,4})[\s\-]?\d{4}[\s\-]?\d{4}\b`)
)

type DataProxy struct { hmacKey []byte; region string }

func NewDataProxy(hmacKey []byte, region string) (*DataProxy, error) {
	if len(hmacKey) < 32 { return nil, ErrInvalidHMACSecret }
	if region != AuthorizedJurisdiction { return nil, ErrJurisdictionFence }
	keyCopy := make([]byte, len(hmacKey)); copy(keyCopy, hmacKey)
	return &DataProxy{hmacKey: keyCopy, region: region}, nil
}

func (dp *DataProxy) PseudonymizeIdentifier(rawID string) string {
	mac := hmac.New(sha256.New, dp.hmacKey); mac.Write([]byte(rawID))
	return "PSEUDO-AR-" + hex.EncodeToString(mac.Sum(nil))[:24]
}

func (dp *DataProxy) SanitizeClinicalText(input string) string {
	out := reCUIL.ReplaceAllString(input, "[REDACTED-CUIL]")
	out = reDNI.ReplaceAllString(out, "[REDACTED-DNI]")
	out = reEmail.ReplaceAllString(out, "[REDACTED-EMAIL]")
	return rePhone.ReplaceAllString(out, "[REDACTED-PHONE]")
}

func StrictDecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, MaxClinicalPayloadBytes)
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body); dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil { return ErrTrailingOrInvalidJSON }
	if dec.More() || dec.Decode(&struct{}{}) != io.EOF { return ErrTrailingOrInvalidJSON }
	return nil
}
