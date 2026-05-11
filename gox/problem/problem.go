// Package problem implements RFC 9457 "Problem Details for HTTP APIs".
// https://www.rfc-editor.org/rfc/rfc9457.html
package problem

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"

	"github.com/victormf2/gox/internal/errorsx"
)

const (
	// ContentTypeJSON is the media type for JSON problem details (RFC 9457 §3).
	ContentTypeJSON = "application/problem+json"
	// ContentTypeXML is the media type for XML problem details (RFC 9457 Appendix B).
	ContentTypeXML = "application/problem+xml"

	// TypeAboutBlank is the default type URI when no specific type is given (RFC 9457 §3.1).
	TypeAboutBlank = "about:blank"
)

// Problem represents an RFC 9457 Problem Details object.
//
// All members are optional; omitted members use their default values as
// described in the spec. Consumers MUST ignore unrecognised extension members.
//
// RFC 9457 §3.1 — Members of a Problem Details Object
type Problem struct {
	// Type is a URI reference that identifies the problem type.
	// When dereferenced, it SHOULD provide human-readable documentation.
	// Defaults to "about:blank" when absent.
	Type string `json:"type,omitempty" xml:"type,omitempty"`

	// Title is a short, human-readable summary of the problem type.
	// It SHOULD NOT change between occurrences of the same type.
	Title string `json:"title,omitempty" xml:"title,omitempty"`

	// Status is the HTTP status code for this occurrence of the problem.
	// It mirrors the status sent in the HTTP response.
	Status int `json:"status,omitempty" xml:"status,omitempty"`

	// Detail is a human-readable explanation specific to this occurrence.
	Detail string `json:"detail,omitempty" xml:"detail,omitempty"`

	// Instance is a URI reference identifying this specific occurrence.
	// It MAY be dereferenceable (e.g. linking to a support ticket).
	Instance string `json:"instance,omitempty" xml:"instance,omitempty"`

	// Errors is a list of ValidationErrors, according to
	// RFC 9457 §4.2 example
	Errors []ValidationError `json:"errors,omitempty" xml:"errors,omitempty"`

	// Extensions holds any additional members defined by the problem type.
	// These are serialised at the top level of the JSON/XML object.
	// RFC 9457 §3.2 — Extension Members
	Extensions map[string]any `json:"-" xml:"-"`
}

// Keys defined by RFC 9457
var RFC9457Keys = map[string]bool{
	"type":     true,
	"title":    true,
	"status":   true,
	"detail":   true,
	"instance": true,
}

// MarshalJSON serialises Problem to JSON, merging Extensions into the top-level
// object as required by RFC 9457 §3.2.
func (p Problem) MarshalJSON() ([]byte, error) {
	problemJSON := make(map[string]any, 5+len(p.Extensions))

	if p.Type == "" {
		problemJSON["type"] = TypeAboutBlank
	}

	if p.Title != "" {
		problemJSON["title"] = p.Title
	}
	if p.Status != 0 {
		problemJSON["status"] = p.Status
	}
	if p.Detail != "" {
		problemJSON["detail"] = p.Detail
	}
	if p.Instance != "" {
		problemJSON["instance"] = p.Instance
	}

	if len(p.Errors) > 0 {
		problemJSON["errors"] = p.Instance
	}

	for key, value := range p.Extensions {
		if RFC9457Keys[key] {
			continue
		}
		problemJSON[key] = value
	}

	return json.Marshal(problemJSON)
}

var ErrProblemUnmarshalJSON = errors.New("unmarshalling Problem JSON")

// UnmarshalJSON deserialises JSON into a Problem, placing unknown members into
// Extensions so no information is lost.
func (problem *Problem) UnmarshalJSON(data []byte) error {
	var jsonMap map[string]json.RawMessage
	err := json.Unmarshal(data, &jsonMap)
	if err != nil {
		return errorsx.JoinErrors(ErrProblemUnmarshalJSON, "unmarshalling full JSON", err)
	}

	problemType, hasProblemType := jsonMap["type"]
	if hasProblemType {
		err = json.Unmarshal(problemType, &problem.Type)
		if err != nil {
			return errorsx.JoinErrors(ErrProblemUnmarshalJSON, "unmarshalling type", err)
		}
	}

	problemTitle, hasProblemTitle := jsonMap["title"]
	if hasProblemTitle {
		err = json.Unmarshal(problemTitle, &problem.Title)
		if err != nil {
			return errorsx.JoinErrors(ErrProblemUnmarshalJSON, "unmarshalling tile", err)
		}
	}
	problemStatus, hasProblemStatus := jsonMap["status"]
	if hasProblemStatus {
		err = json.Unmarshal(problemStatus, &problem.Status)
		if err != nil {
			return errorsx.JoinErrors(ErrProblemUnmarshalJSON, "unmarshalling status", err)
		}
	}
	problemDetail, hasProblemDetail := jsonMap["detail"]
	if hasProblemDetail {
		err = json.Unmarshal(problemDetail, &problem.Detail)
		if err != nil {
			return errorsx.JoinErrors(ErrProblemUnmarshalJSON, "unmarshalling detail", err)
		}
	}
	problemInstance, hasProblemInstance := jsonMap["instance"]
	if hasProblemInstance {
		err = json.Unmarshal(problemInstance, &problem.Instance)
		if err != nil {
			return errorsx.JoinErrors(ErrProblemUnmarshalJSON, "unmarshalling instance", err)
		}
	}

	problemErrors, hasProblemErrors := jsonMap["errors"]
	if hasProblemErrors {
		// ignore errors, will be added to the extensions
		_ = json.Unmarshal(problemErrors, &problem.Errors)
	}

	for key, jsonValue := range jsonMap {
		if RFC9457Keys[key] {
			continue
		}
		if key == "errors" && len(problem.Errors) > 0 {
			continue
		}
		if problem.Extensions == nil {
			problem.Extensions = make(map[string]any)
		}
		var val any
		err = json.Unmarshal(jsonValue, &val)
		if err != nil {
			message := fmt.Sprintf("unmarshalling %s", key)
			return errorsx.JoinErrors(ErrProblemUnmarshalJSON, message, err)
		}
		problem.Extensions[key] = val
	}

	return nil
}

// --- RFC 9457 §3.1 default handling ----------------------------------------

// EffectiveType returns the problem type URI, defaulting to "about:blank" when
// the Type field is empty, as required by RFC 9457 §3.1.
func (p *Problem) EffectiveType() string {
	if p.Type == "" {
		return TypeAboutBlank
	}
	return p.Type
}

// --- Validation errors (RFC 9457 §4.2 example) ------------------------------

// ValidationError describes a single field-level validation failure, used
// inside the "errors" extension of a validation problem.
type ValidationError struct {
	// Detail is a human-readable description of the validation failure.
	Detail string `json:"detail" xml:"detail"`

	// Pointer is a JSON Pointer (RFC 6901) to the field that caused the error.
	// e.g. "/name" or "/address/postcode"
	Pointer string `json:"pointer,omitempty" xml:"pointer,omitempty"`
}

// NewValidationProblem returns a 422 problem whose "errors" extension
// carries per-field validation errors, as shown in RFC 9457 §4.2.
func NewValidationProblem(typeURI, detail string, errs []ValidationError) *Problem {
	return &Problem{
		Type:   typeURI,
		Title:  "Validation Error",
		Status: http.StatusUnprocessableEntity,
		Detail: detail,
		Extensions: map[string]any{
			"errors": errs,
		},
	}
}

// ---------------------------------------------------------------------------
// Field interfaces
//
// Each interface corresponds to one RFC 9457 field. Custom error types may
// implement any combination of them; FromError probes for each one and copies
// whatever it finds into a Problem.
//
// Implementing all five plus Extender gives you full spec coverage. You are
// free to implement only the interfaces that make sense for a given error type.
// ---------------------------------------------------------------------------

// IProblemType provides the "type" URI that identifies the problem type (RFC 9457 §3.1).
// The returned string SHOULD be an absolute URI. When it is a locator (http/https)
// it SHOULD resolve to human-readable documentation.
// If not implemented, "about:blank" is used as the default.
type IProblemType interface {
	ProblemType() string
}

// IProblemTitle provides the "title" field: a short, human-readable summary of the
// problem type that SHOULD NOT change between occurrences (RFC 9457 §3.1).
// If not implemented, the title is derived from the HTTP status code when
// possible, or omitted.
type IProblemTitle interface {
	ProblemTitle() string
}

// IProblemStatus provides the "status" field: the HTTP status code generated for this
// occurrence (RFC 9457 §3.1). The value SHOULD mirror the actual response status.
// If not implemented, the status is omitted (treated as 0 / unknown).
type IProblemStatus interface {
	ProblemStatus() int
}

// IProblemDetail provides the "detail" field: a human-readable explanation specific
// to this occurrence of the problem (RFC 9457 §3.1).
// If not implemented, detail is omitted from the output.
type IProblemDetail interface {
	ProblemDetail() string
}

// IProblemInstance provides the "instance" field: a URI reference that identifies this
// specific occurrence of the problem (RFC 9457 §3.1). It MAY be dereferenceable,
// e.g. a link to a support ticket or log entry.
// If not implemented, instance is omitted from the output.
type IProblemInstance interface {
	ProblemInstance() string
}

// IProblemExtensions provides extension members that are merged into the top-level
// problem object alongside the standard fields (RFC 9457 §3.2).
// Keys MUST NOT shadow the five standard field names.
// Consumers that do not recognise an extension key MUST ignore it.
type IProblemExtensions interface {
	ProblemExtensions() map[string]any
}

// IProblemErrors provides a list of ValidationErrors according to
// RFC 9457 §4.2 example
type IProblemErrors interface {
	ProblemErrors() []ValidationError
}

// ---------------------------------------------------------------------------
// FromError converts any error into a *Problem by probing for the field
// interfaces above. Fields present on the error take precedence over defaults.
//
// Conversion order:
//  1. If err already is a *Problem, return a shallow copy.
//  2. Otherwise build a new Problem by checking each interface in turn.
//  3. Detailer falls back to err.Error() so the message is never lost.
//  4. Typer falls back to "about:blank" per RFC 9457 §3.1.
//  5. When Titler is absent and a status code is available, Title is set from
//     http.StatusText so clients always get a human-readable summary.
// ---------------------------------------------------------------------------

// FromError converts any error into a *Problem, probing for the field
// interfaces (Typer, Titler, Statuser, Detailer, Instancer, Extender).
// It never returns nil; unrecognised errors become a 500 with the error
// message in Detail.
func FromError(err error) *Problem {
	if err == nil {
		return nil
	}

	// Fast path: already a *Problem.
	problem, isProblem := err.(*Problem)
	if isProblem {
		return problem
	}

	problem = &Problem{}

	// type
	problemType, hasType := err.(IProblemType)
	if hasType {
		problem.Type = problemType.ProblemType()
	} else {
		problem.Type = TypeAboutBlank
	}

	// title
	problemTitle, hasTitle := err.(IProblemTitle)
	if hasTitle {
		problem.Title = problemTitle.ProblemTitle()
	}

	// status
	problemStatus, hasStatus := err.(IProblemStatus)
	if hasStatus {
		problem.Status = problemStatus.ProblemStatus()
	} else {
		problem.Status = http.StatusInternalServerError
	}

	// Derive title from status when the error didn't supply one.
	if problem.Title == "" && problem.Status != 0 {
		problem.Title = http.StatusText(problem.Status)
	}

	// detail — falls back to "Some instability happening" so we don't expose the error to the client for security reasons.
	problemDetail, hasProblemDetails := err.(IProblemDetail)
	if hasProblemDetails {
		problem.Detail = problemDetail.ProblemDetail()
	} else {
		problem.Detail = "Some instability happening"
	}

	// instance
	problemInstance, hasInstance := err.(IProblemInstance)
	if hasInstance {
		problem.Instance = problemInstance.ProblemInstance()
	}

	problemErrors, hasErrors := err.(IProblemErrors)
	if hasErrors {
		problem.Errors = problemErrors.ProblemErrors()
	}

	// extensions
	problemExtensions, hasExtensions := err.(IProblemExtensions)
	if hasExtensions {
		extensions := problemExtensions.ProblemExtensions()
		if len(extensions) > 0 {
			problem.Extensions = extensions
		}
	}

	return problem
}

// WriteErrorJSON is a convenience helper that calls FromError and writes the
// result as an application/problem+json response.
func WriteErrorJSON(w http.ResponseWriter, err error) error {
	return FromError(err).WriteJSON(w)
}

// WriteErrorXML is a convenience helper that calls FromError and writes the
// result as an application/problem+xml response.
func WriteErrorXML(w http.ResponseWriter, err error) error {
	return FromError(err).WriteXML(w)
}

// --- Constructor helpers ----------------------------------------------------

// New creates a Problem with the given type URI and HTTP status code.
// Use the With* option functions to attach additional fields.
func New(typeURI string, status int, opts ...Option) *Problem {
	p := &Problem{
		Type:   typeURI,
		Status: status,
	}
	for _, o := range opts {
		o(p)
	}
	return p
}

// FromStatus creates a minimal Problem from an HTTP status code.
// The type defaults to "about:blank" and the title is set from the status text.
func FromStatus(status int, opts ...Option) *Problem {
	p := &Problem{
		Type:   TypeAboutBlank,
		Status: status,
		Title:  http.StatusText(status),
	}
	for _, o := range opts {
		o(p)
	}
	return p
}

// Option is a functional option for building a Problem.
type Option func(*Problem)

// WithTitle sets the Title field.
func WithTitle(title string) Option {
	return func(p *Problem) { p.Title = title }
}

// WithDetail sets the Detail field.
func WithDetail(detail string) Option {
	return func(p *Problem) { p.Detail = detail }
}

// WithInstance sets the Instance field.
func WithInstance(instance string) Option {
	return func(p *Problem) { p.Instance = instance }
}

// WithErrors sets the Errors field.
func WithErrors(errors []ValidationError) Option {
	return func(p *Problem) { p.Errors = errors }
}

// WithExtension adds a single extension member (RFC 9457 §3.2).
func WithExtension(key string, value any) Option {
	return func(p *Problem) {
		if p.Extensions == nil {
			p.Extensions = make(map[string]any)
		}
		p.Extensions[key] = value
	}
}

// --- HTTP response helpers --------------------------------------------------

// WriteJSON writes the problem as an RFC 9457 JSON response.
// It sets Content-Type to application/problem+json and the given status code.
func (p *Problem) WriteJSON(w http.ResponseWriter) error {
	status := p.Status
	if status == 0 {
		status = http.StatusInternalServerError
	}
	w.Header().Set("Content-Type", ContentTypeJSON)
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(p)
}

// WriteXML writes the problem as an RFC 9457 XML response (Appendix B).
func (p *Problem) WriteXML(w http.ResponseWriter) error {
	status := p.Status
	if status == 0 {
		status = http.StatusInternalServerError
	}
	w.Header().Set("Content-Type", ContentTypeXML)
	w.WriteHeader(status)
	_, err := w.Write([]byte(xml.Header))
	if err != nil {
		return err
	}
	return xml.NewEncoder(w).Encode(p)
}

// --- error interface --------------------------------------------------------

// Error implements the error interface so Problem values can be returned as
// errors and unwrapped with errors.As.
func (p *Problem) Error() string {
	if p.Detail != "" {
		return p.Detail
	}
	if p.Title != "" {
		return p.Title
	}
	return p.EffectiveType()
}

// AsProblem is a convenience wrapper around errors.As for Problem pointers.
func AsProblem(err error) (*Problem, bool) {
	var p *Problem
	ok := errors.As(err, &p)
	return p, ok
}

// --- Pre-built common problems ----------------------------------------------

// Common pre-built problem constructors for frequently used HTTP status codes.

func BadRequest(title string, opts ...Option) *Problem {
	optsWithTitle := append([]Option{WithTitle(title)}, opts...)
	return FromStatus(http.StatusBadRequest, optsWithTitle...)
}

func Unauthorized(title string, opts ...Option) *Problem {
	optsWithTitle := append([]Option{WithTitle(title)}, opts...)
	return FromStatus(http.StatusUnauthorized, optsWithTitle...)
}

func Forbidden(title string, opts ...Option) *Problem {
	optsWithTitle := append([]Option{WithTitle(title)}, opts...)
	return FromStatus(http.StatusForbidden, optsWithTitle...)
}

func NotFound(title string, opts ...Option) *Problem {
	optsWithTitle := append([]Option{WithTitle(title)}, opts...)
	return FromStatus(http.StatusNotFound, optsWithTitle...)
}

func Conflict(title string, opts ...Option) *Problem {
	optsWithTitle := append([]Option{WithTitle(title)}, opts...)
	return FromStatus(http.StatusConflict, optsWithTitle...)
}

func UnprocessableEntity(title string, opts ...Option) *Problem {
	optsWithTitle := append([]Option{WithTitle(title)}, opts...)
	return FromStatus(http.StatusUnprocessableEntity, optsWithTitle...)
}

func InternalServerError(title string, opts ...Option) *Problem {
	optsWithTitle := append([]Option{WithTitle(title)}, opts...)
	return FromStatus(http.StatusInternalServerError, optsWithTitle...)
}
