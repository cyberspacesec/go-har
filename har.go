// Package har provides functionality for parsing and manipulating HAR (HTTP Archive) files.
// This is a compatibility wrapper that forwards to the implementation in pkg/har.
package har

import (
	"github.com/cyberspacesec/go-har/pkg/har"
)

// Har represents a HAR file
type Har = har.Har

// Log represents the log section of a HAR file
type Log = har.Log

// Creator represents the creator information
type Creator = har.Creator

// Browser represents browser information
type Browser = har.Browser

// PageTimings represents page timing information
type PageTimings = har.PageTimings

// Pages represents a page in the HAR file
type Pages = har.Pages

// Headers represents HTTP headers
type Headers = har.Headers

// Request represents an HTTP request
type Request = har.Request

// Cookie represents an HTTP cookie
type Cookie = har.Cookie

// Content represents response content
type Content = har.Content

// Response represents an HTTP response
type Response = har.Response

// BeforeRequest represents cache state before request
type BeforeRequest = har.BeforeRequest

// AfterRequest represents cache state after request
type AfterRequest = har.AfterRequest

// Cache represents cache information
type Cache = har.Cache

// Timings represents timing information for a request
type Timings = har.Timings

// Entries represents an entry in the HAR file
type Entries = har.Entries

// Initiator represents the initiator of a request
type Initiator = har.Initiator

// Stack represents a call stack
type Stack = har.Stack

// Parent represents parent information in a call stack
type Parent = har.Parent

// ParentID represents a parent ID in a call stack
type ParentID = har.ParentID

// CallFrame represents a frame in a call stack
type CallFrame = har.CallFrame

// HTTPMethod enum type for HTTP methods
type HTTPMethod = har.HTTPMethod

// HTTP Method constants
const (
	MethodUnknown = har.MethodUnknown
	MethodGET     = har.MethodGET
	MethodPOST    = har.MethodPOST
	MethodPUT     = har.MethodPUT
	MethodDELETE  = har.MethodDELETE
	MethodHEAD    = har.MethodHEAD
	MethodOPTIONS = har.MethodOPTIONS
	MethodPATCH   = har.MethodPATCH
	MethodCONNECT = har.MethodCONNECT
	MethodTRACE   = har.MethodTRACE
)

// ConvertFormat for conversion formats
type ConvertFormat = har.ConvertFormat

// Format constants
const (
	FormatCSV      = har.FormatCSV
	FormatMarkdown = har.FormatMarkdown
	FormatHTML     = har.FormatHTML
	FormatText     = har.FormatText
)

// Error types
type (
	ErrorCode              = har.ErrorCode
	HarError               = har.HarError
	ParseOptions           = har.ParseOptions
	FilterOptions          = har.FilterOptions
	FilterResult           = har.FilterResult
	Result                 = har.Result
	ConvertOptions         = har.ConvertOptions
	OptimizedHar           = har.OptimizedHar
	OptimizedEntries       = har.OptimizedEntries
	OptimizedRequest       = har.OptimizedRequest
	OptimizedResponse      = har.OptimizedResponse
	OptimizedContent       = har.OptimizedContent
	OptimizedTimings       = har.OptimizedTimings
	StreamingHar           = har.StreamingHar
	EntryIterator          = har.EntryIterator
	StreamingEntryIterator = har.StreamingEntryIterator
	LazyHar                = har.LazyHar
	LazyContent            = har.LazyContent
	LazyResponse           = har.LazyResponse
	LazyEntries            = har.LazyEntries

	// 接口类型
	HARProvider         = har.HARProvider
	EntryProvider       = har.EntryProvider
	RequestProvider     = har.RequestProvider
	ResponseProvider    = har.ResponseProvider
	HeaderProvider      = har.HeaderProvider
	CookieProvider      = har.CookieProvider
	ContentProvider     = har.ContentProvider
	TimingsProvider     = har.TimingsProvider
	PageProvider        = har.PageProvider
	PageTimingsProvider = har.PageTimingsProvider

	// 选项类型
	Option = har.Option
)

// Error code constants
const (
	ErrCodeUnknown       = har.ErrCodeUnknown
	ErrCodeFileSystem    = har.ErrCodeFileSystem
	ErrCodeJSONParse     = har.ErrCodeJSONParse
	ErrCodeInvalidFormat = har.ErrCodeInvalidFormat
	ErrCodeValidation    = har.ErrCodeValidation
	ErrCodeMissingField  = har.ErrCodeMissingField
	ErrCodeInvalidValue  = har.ErrCodeInvalidValue
	ErrCodeUnsupported   = har.ErrCodeUnsupported
)

// Forward all functions
var (
	// Basic operations
	ParseHarFile = har.ParseHarFile
	ParseHar     = har.ParseHar
	NewHar       = har.NewHar

	// Optimized parsing
	ParseHarFileOptimized = har.ParseHarFileOptimized
	ParseHarOptimized     = har.ParseHarOptimized
	ToOptimizedHar        = har.ToOptimizedHar

	// Lazy loading
	ParseHarWithLazyLoading     = har.ParseHarWithLazyLoading
	ParseHarFileWithLazyLoading = har.ParseHarFileWithLazyLoading

	// Streaming
	NewStreamingHarFromFile = har.NewStreamingHarFromFile

	// Enhanced parsing
	ParseHarWithOptions      = har.ParseHarWithOptions
	ParseHarFileWithOptions  = har.ParseHarFileWithOptions
	ParseHarEnhanced         = har.ParseHarEnhanced
	ParseHarFileEnhanced     = har.ParseHarFileEnhanced
	ParseHarLenient          = har.ParseHarLenient
	ParseHarFileLenient      = har.ParseHarFileLenient
	ParseHarWithWarnings     = har.ParseHarWithWarnings
	ParseHarFileWithWarnings = har.ParseHarFileWithWarnings
	DefaultParseOptions      = har.DefaultParseOptions

	// Error utilities
	NewHarError            = har.NewHarError
	NewFileSystemError     = har.NewFileSystemError
	NewJSONParseError      = har.NewJSONParseError
	WrapJSONUnmarshalError = har.WrapJSONUnmarshalError
	NewValidationError     = har.NewValidationError
	NewInvalidFormatError  = har.NewInvalidFormatError
	NewMissingFieldError   = har.NewMissingFieldError
	NewInvalidValueError   = har.NewInvalidValueError
	NewUnsupportedError    = har.NewUnsupportedError

	// Utilities
	ParseMethod           = har.ParseMethod
	DefaultConvertOptions = har.DefaultConvertOptions

	// 新的函数选项模式API
	Parse                      = har.Parse
	ParseFile                  = har.ParseFile
	NewStreamingParser         = har.NewStreamingParser
	NewStreamingParserFromFile = har.NewStreamingParserFromFile

	// 选项函数
	WithLenient         = har.WithLenient
	WithSkipValidation  = har.WithSkipValidation
	WithCollectWarnings = har.WithCollectWarnings
	WithMaxWarnings     = har.WithMaxWarnings
	WithMemoryOptimized = har.WithMemoryOptimized
	WithLazyLoading     = har.WithLazyLoading
	WithStreaming       = har.WithStreaming

	// 预定义选项组
	OptMemoryEfficient = har.OptMemoryEfficient
	OptFast            = har.OptFast
	OptLenient         = har.OptLenient
	OptPerformance     = har.OptPerformance
)
