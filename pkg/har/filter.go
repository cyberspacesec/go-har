package har

import (
	"regexp"
	"strings"
	"time"
)

// FilterResult 过滤结果
type FilterResult struct {
	Entries []Entries
}

// FilterOptions 过滤选项
type FilterOptions struct {
	URL             string    // URL包含的字符串或正则表达式
	Method          string    // 请求方法
	StatusCode      int       // 响应状态码
	StatusCodeMin   int       // 最小状态码
	StatusCodeMax   int       // 最大状态码
	ContentType     string    // 内容类型
	StartTime       time.Time // 开始时间
	EndTime         time.Time // 结束时间
	MinDuration     float64   // 最小请求持续时间(ms)
	MaxDuration     float64   // 最大请求持续时间(ms)
	ResourceType    string    // 资源类型
	HasError        bool      // 是否有错误
	HeaderName      string    // 请求头名
	HeaderValue     string    // 请求头值
	RespHeaderName  string    // 响应头名
	RespHeaderValue string    // 响应头值
	UseRegex        bool      // 使用正则表达式匹配
}

// Filter 按条件过滤条目
func (h *Har) Filter(options FilterOptions) *FilterResult {
	var result []Entries

	for _, entry := range h.Log.Entries {
		if matchesFilter(entry, options) {
			result = append(result, entry)
		}
	}

	return &FilterResult{
		Entries: result,
	}
}

// 检查条目是否符合过滤条件
func matchesFilter(entry Entries, options FilterOptions) bool {
	// URL过滤
	if options.URL != "" {
		if options.UseRegex {
			re, err := regexp.Compile(options.URL)
			if err == nil && !re.MatchString(entry.Request.URL) {
				return false
			}
		} else if !strings.Contains(entry.Request.URL, options.URL) {
			return false
		}
	}

	// 请求方法过滤
	if options.Method != "" && entry.Request.Method != options.Method {
		return false
	}

	// 状态码过滤
	if options.StatusCode > 0 && entry.Response.Status != options.StatusCode {
		return false
	}

	// 状态码范围过滤
	if options.StatusCodeMin > 0 && entry.Response.Status < options.StatusCodeMin {
		return false
	}
	if options.StatusCodeMax > 0 && entry.Response.Status > options.StatusCodeMax {
		return false
	}

	// 内容类型过滤
	if options.ContentType != "" {
		matched := false
		for _, header := range entry.Response.Headers {
			if strings.EqualFold(header.Name, "Content-Type") && strings.Contains(header.Value, options.ContentType) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// 时间范围过滤
	if !options.StartTime.IsZero() && entry.StartedDateTime.Before(options.StartTime) {
		return false
	}
	if !options.EndTime.IsZero() && entry.StartedDateTime.After(options.EndTime) {
		return false
	}

	// 持续时间过滤
	if options.MinDuration > 0 && entry.Time < options.MinDuration {
		return false
	}
	if options.MaxDuration > 0 && entry.Time > options.MaxDuration {
		return false
	}

	// 资源类型过滤
	if options.ResourceType != "" && entry.ResourceType != options.ResourceType {
		return false
	}

	// 错误过滤
	if options.HasError && (entry.Response.Status < 400 || entry.Response.Status >= 600) {
		return false
	}

	// 请求头过滤
	if options.HeaderName != "" {
		matched := false
		for _, header := range entry.Request.Headers {
			if strings.EqualFold(header.Name, options.HeaderName) {
				if options.HeaderValue == "" || strings.Contains(header.Value, options.HeaderValue) {
					matched = true
					break
				}
			}
		}
		if !matched {
			return false
		}
	}

	// 响应头过滤
	if options.RespHeaderName != "" {
		matched := false
		for _, header := range entry.Response.Headers {
			if strings.EqualFold(header.Name, options.RespHeaderName) {
				if options.RespHeaderValue == "" || strings.Contains(header.Value, options.RespHeaderValue) {
					matched = true
					break
				}
			}
		}
		if !matched {
			return false
		}
	}

	return true
}

// 快捷过滤方法

// FindByURL 按URL查找
func (h *Har) FindByURL(urlStr string, useRegex bool) *FilterResult {
	return h.Filter(FilterOptions{
		URL:      urlStr,
		UseRegex: useRegex,
	})
}

// FindByMethod 按HTTP方法查找
func (h *Har) FindByMethod(method string) *FilterResult {
	return h.Filter(FilterOptions{
		Method: method,
	})
}

// FindByStatusCode 按状态码查找
func (h *Har) FindByStatusCode(statusCode int) *FilterResult {
	return h.Filter(FilterOptions{
		StatusCode: statusCode,
	})
}

// FindErrors 查找所有错误请求
func (h *Har) FindErrors() *FilterResult {
	return h.Filter(FilterOptions{
		HasError: true,
	})
}

// FindByTimeRange 按时间范围查找
func (h *Har) FindByTimeRange(start, end time.Time) *FilterResult {
	return h.Filter(FilterOptions{
		StartTime: start,
		EndTime:   end,
	})
}

// FindByContentType 按内容类型查找
func (h *Har) FindByContentType(contentType string) *FilterResult {
	return h.Filter(FilterOptions{
		ContentType: contentType,
	})
}

// FindSlowRequests 查找慢请求
func (h *Har) FindSlowRequests(minDuration float64) *FilterResult {
	return h.Filter(FilterOptions{
		MinDuration: minDuration,
	})
}

// Count 获取过滤结果数量
func (fr *FilterResult) Count() int {
	return len(fr.Entries)
}

// First 获取第一个结果
func (fr *FilterResult) First() *Entries {
	if len(fr.Entries) > 0 {
		return &fr.Entries[0]
	}
	return nil
}

// ToHar 将过滤结果转换为新的Har对象
func (fr *FilterResult) ToHar() *Har {
	har := NewHar()
	har.Log.Entries = fr.Entries
	return har
}
