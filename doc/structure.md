# Go-HAR 数据结构详解

本文档详细说明了 HTTP Archive (HAR) 格式的数据结构以及在 Go-HAR 中的实现。

## HAR 文件概述

HAR (HTTP Archive) 文件是 HTTP 交互的存档格式，以 JSON 格式存储 HTTP 请求和响应数据。它由浏览器开发工具导出，用于分析网页性能、调试问题和记录网络交互。

## 主要结构

### Har 结构

`Har` 是整个 HAR 文件的根结构，包含唯一的 `Log` 字段：

```go
type Har struct {
    Log Log `json:"log"`
}
```

### Log 结构

`Log` 包含所有 HAR 数据的主容器：

```go
type Log struct {
    Version string   `json:"version"`     // HAR 格式版本，通常是 "1.2"
    Creator Creator  `json:"creator"`     // 创建 HAR 文件的工具信息
    Browser Browser  `json:"browser,omitempty"` // 用于捕获的浏览器信息（可选）
    Pages   []Pages  `json:"pages,omitempty"`   // 捕获的页面信息（可选）
    Entries []Entries `json:"entries"`     // 所有 HTTP 请求和响应条目
    Comment string   `json:"comment,omitempty"` // 用户提供的注释（可选）
}
```

各字段详解：
- `Version`: HAR 格式的版本，规范版本为 "1.2"
- `Creator`: 创建 HAR 文件的应用信息
- `Browser`: 生成这些请求的浏览器信息
- `Pages`: 记录的页面信息
- `Entries`: 所有HTTP请求和响应的条目
- `Comment`: 用户注释（可选）

### Creator 和 Browser 结构

描述创建工具和浏览器的信息：

```go
type Creator struct {
    Name    string `json:"name"`    // 应用名称
    Version string `json:"version"` // 应用版本
    Comment string `json:"comment,omitempty"` // 注释（可选）
}

type Browser struct {
    Name    string `json:"name"`    // 浏览器名称
    Version string `json:"version"` // 浏览器版本
    Comment string `json:"comment,omitempty"` // 注释（可选）
}
```

### Pages 结构

描述捕获的页面信息：

```go
type Pages struct {
    StartedDateTime time.Time   `json:"startedDateTime"` // 页面加载开始时间（ISO 8601）
    ID              string      `json:"id"`              // 页面唯一标识符
    Title           string      `json:"title"`           // 页面标题
    PageTimings     PageTimings `json:"pageTimings"`     // 页面加载时间信息
    Comment         string      `json:"comment,omitempty"` // 注释（可选）
}

type PageTimings struct {
    OnContentLoad float64 `json:"onContentLoad,omitempty"` // DOMContentLoaded 事件，毫秒，-1表示不可用
    OnLoad        float64 `json:"onLoad,omitempty"`        // load 事件，毫秒，-1表示不可用
    Comment       string  `json:"comment,omitempty"`       // 注释（可选）
}
```

### Entries 结构

`Entries` 是 HAR 文件中最重要的部分，包含完整的 HTTP 请求/响应信息：

```go
type Entries struct {
    Pageref         string     `json:"pageref,omitempty"` // 引用的页面ID
    StartedDateTime time.Time  `json:"startedDateTime"`   // 请求开始时间（ISO 8601）
    Time            float64    `json:"time"`              // 请求总耗时（毫秒）
    Request         Request    `json:"request"`           // 请求信息
    Response        Response   `json:"response"`          // 响应信息
    Cache           Cache      `json:"cache"`             // 缓存信息
    Timings         Timings    `json:"timings"`           // 请求各阶段时间
    ServerIPAddress string     `json:"serverIPAddress,omitempty"` // 服务器IP地址
    Connection      string     `json:"connection,omitempty"`      // 连接信息
    Comment         string     `json:"comment,omitempty"`         // 注释
}
```

各字段详解：
- `Pageref`: 对应的页面ID，可用于关联请求所属的页面
- `StartedDateTime`: 请求开始的精确时间，ISO 8601格式
- `Time`: 整个请求耗时（毫秒），包括从请求发起到接收完响应的全部时间
- `Request`: 包含完整的HTTP请求信息
- `Response`: 包含完整的HTTP响应信息
- `Cache`: 浏览器缓存状态信息
- `Timings`: 请求各个阶段的详细耗时
- `ServerIPAddress`: 服务器IP地址（可选）
- `Connection`: 连接标识符，如"52492"（可选）

### Request 结构

HTTP 请求的详细信息：

```go
type Request struct {
    Method      string     `json:"method"`      // HTTP方法（GET、POST等）
    URL         string     `json:"url"`         // 完整URL
    HTTPVersion string     `json:"httpVersion"` // HTTP版本
    Cookies     []Cookie   `json:"cookies"`     // Cookie信息
    Headers     []Headers  `json:"headers"`     // 请求头信息
    QueryString []QueryString `json:"queryString"` // URL查询参数
    PostData    PostData   `json:"postData,omitempty"` // POST数据（可选）
    HeadersSize int        `json:"headersSize"` // 请求头大小（字节）
    BodySize    int        `json:"bodySize"`    // 请求体大小（字节）
    Comment     string     `json:"comment,omitempty"`  // 注释（可选）
}
```

### Response 结构

HTTP 响应的详细信息：

```go
type Response struct {
    Status       int       `json:"status"`       // HTTP状态码
    StatusText   string    `json:"statusText"`   // 状态文本
    HTTPVersion  string    `json:"httpVersion"`  // HTTP版本
    Cookies      []Cookie  `json:"cookies"`      // Cookie信息
    Headers      []Headers `json:"headers"`      // 响应头信息
    Content      Content   `json:"content"`      // 响应内容
    RedirectURL  string    `json:"redirectURL"`  // 重定向URL
    HeadersSize  int       `json:"headersSize"`  // 响应头大小（字节）
    BodySize     int       `json:"bodySize"`     // 响应体大小（字节）
    TransferSize int       `json:"_transferSize,omitempty"` // 传输大小（字节）
    Error        string    `json:"_error,omitempty"`        // 错误信息
    Comment      string    `json:"comment,omitempty"`       // 注释（可选）
}
```

### Headers 和 Cookie 结构

HTTP 头和 Cookie 信息：

```go
type Headers struct {
    Name    string `json:"name"`    // 头部名称
    Value   string `json:"value"`   // 头部值
    Comment string `json:"comment,omitempty"` // 注释（可选）
}

type Cookie struct {
    Name     string    `json:"name"`     // Cookie名称
    Value    string    `json:"value"`    // Cookie值
    Path     string    `json:"path,omitempty"`     // 路径
    Domain   string    `json:"domain,omitempty"`   // 域名
    Expires  time.Time `json:"expires,omitempty"`  // 过期时间
    HTTPOnly bool      `json:"httpOnly,omitempty"` // 是否HTTPOnly
    Secure   bool      `json:"secure,omitempty"`   // 是否安全
    SameSite string    `json:"sameSite,omitempty"` // SameSite属性
    Comment  string    `json:"comment,omitempty"`  // 注释（可选）
}
```

### QueryString 和 PostData 结构

URL 查询参数和 POST 数据：

```go
type QueryString struct {
    Name    string `json:"name"`    // 参数名
    Value   string `json:"value"`   // 参数值
    Comment string `json:"comment,omitempty"` // 注释（可选）
}

type PostData struct {
    MimeType string    `json:"mimeType"` // MIME类型
    Params   []Params  `json:"params"`   // 参数（适用于表单）
    Text     string    `json:"text"`     // 文本内容
    Comment  string    `json:"comment,omitempty"` // 注释（可选）
}

type Params struct {
    Name        string  `json:"name"`        // 参数名
    Value       string  `json:"value,omitempty"` // 参数值
    FileName    string  `json:"fileName,omitempty"` // 文件名（用于文件上传）
    ContentType string  `json:"contentType,omitempty"` // 内容类型
    Comment     string  `json:"comment,omitempty"` // 注释（可选）
}
```

### Content 结构

HTTP 响应内容：

```go
type Content struct {
    Size        int    `json:"size"`        // 内容大小（字节）
    MimeType    string `json:"mimeType"`    // MIME类型
    Compression int    `json:"compression,omitempty"` // 压缩大小（字节）
    Text        string `json:"text,omitempty"`        // 实际内容（可选）
    Encoding    string `json:"encoding,omitempty"`    // 编码方式（如base64）
    Comment     string `json:"comment,omitempty"`     // 注释（可选）
}
```

### Cache 结构

浏览器缓存信息：

```go
type Cache struct {
    BeforeRequest *BeforeRequest `json:"beforeRequest,omitempty"` // 请求前缓存状态
    AfterRequest  *AfterRequest  `json:"afterRequest,omitempty"`  // 请求后缓存状态
    Comment       string         `json:"comment,omitempty"`       // 注释（可选）
}

type BeforeRequest struct {
    Expires    time.Time `json:"expires,omitempty"`    // 过期时间
    LastAccess time.Time `json:"lastAccess"`           // 最后访问时间
    ETag       string    `json:"eTag"`                 // ETag
    HitCount   int       `json:"hitCount"`             // 命中次数
    Comment    string    `json:"comment,omitempty"`    // 注释（可选）
}

type AfterRequest struct {
    Expires    time.Time `json:"expires,omitempty"`    // 过期时间
    LastAccess time.Time `json:"lastAccess"`           // 最后访问时间
    ETag       string    `json:"eTag"`                 // ETag
    HitCount   int       `json:"hitCount"`             // 命中次数
    Comment    string    `json:"comment,omitempty"`    // 注释（可选）
}
```

### Timings 结构

请求各阶段耗时信息：

```go
type Timings struct {
    Blocked float64 `json:"blocked"`             // 阻塞时间
    DNS     float64 `json:"dns"`                 // DNS解析时间
    Connect float64 `json:"connect"`             // TCP连接时间
    Send    float64 `json:"send"`                // 发送请求时间
    Wait    float64 `json:"wait"`                // 等待响应时间
    Receive float64 `json:"receive"`             // 接收响应时间
    Ssl     float64 `json:"ssl,omitempty"`       // SSL/TLS协商时间
    Comment string  `json:"comment,omitempty"`   // 注释（可选）
}
```

各字段详解：
- `Blocked`: 请求被阻塞的时间（如等待可用TCP连接）
- `DNS`: DNS解析耗时，如不适用或值为-1表示DNS已缓存
- `Connect`: 建立TCP连接耗时，如不适用或值为-1表示重用连接
- `Send`: 发送HTTP请求到服务器的耗时
- `Wait`: 等待服务器首字节响应的耗时
- `Receive`: 接收响应数据的耗时
- `Ssl`: SSL/TLS协商耗时，仅适用于HTTPS请求

## Go-HAR的优化结构

除了标准HAR结构外，Go-HAR还提供了几种优化结构，用于不同场景：

### OptimizedHar

内存优化版本，适用于处理大型HAR文件：

```go
type OptimizedHar struct {
    Log struct {
        Version string
        Creator Creator
        Pages   []Pages
        Entries []OptimizedEntries
    }
}

type OptimizedEntries struct {
    StartedDateTime time.Time
    Time            float64
    PageRef         *string              // 使用指针允许nil值
    Request         OptimizedRequest
    Response        OptimizedResponse
    Timings         OptimizedTimings
}

type OptimizedRequest struct {
    Method      HTTPMethod        // 使用枚举而不是字符串
    URL         string
    HTTPVersion string
    Cookies     []Cookie
    Headers     map[string]string // 使用map而不是数组，优化查找
    QueryString map[string]string // 使用map而不是数组
    HeadersSize *int              // 使用指针允许nil值
    BodySize    *int              // 使用指针允许nil值
}
```

主要优化点：
- 使用枚举代替字符串
- 使用map代替数组，优化查找性能
- 使用指针代表可选值，避免不必要的内存占用

### LazyHar

懒加载版本，适用于需要延迟加载大型内容的场景：

```go
type LazyHar struct {
    Log struct {
        Version string
        Creator Creator
        Pages   []Pages
        Entries []LazyEntries
    }
}

type LazyEntries struct {
    StartedDateTime time.Time
    Time            float64
    Pageref         string
    Request         Request
    Response        LazyResponse
    Timings         Timings
}

type LazyResponse struct {
    // 基本字段直接加载
    Status       int
    StatusText   string
    HTTPVersion  string
    // 内容懒加载
    Content      *LazyContent
    // 其他字段...
}

type LazyContent struct {
    // 基本信息直接加载
    Size        int
    MimeType    string
    // 大型内容懒加载
    Text        *string
    Encoding    *string
    // 内部状态
    loaded      bool
    dataSource  ContentDataSource
}
```

主要特点：
- 基本信息（大小、类型等）直接加载
- 大型内容（如响应体）延迟到需要时才加载
- 提供统一的接口，使用方式与标准结构相同

## 接口

Go-HAR 定义了一系列接口，确保不同实现之间的互操作性：

```go
// HARProvider 所有HAR实现的统一接口
type HARProvider interface {
    GetVersion() string
    GetCreator() Creator
    GetEntries() []EntryProvider
    GetPages() []PageProvider
    ToStandard() *Har
}

// EntryProvider 单个条目的接口
type EntryProvider interface {
    GetStartedDateTime() time.Time
    GetTime() float64
    GetRequest() RequestProvider
    GetResponse() ResponseProvider
    GetTimings() TimingsProvider
    GetPageref() string
    ToStandard() Entries
}

// 其他接口: RequestProvider, ResponseProvider, ContentProvider等
```

这些接口确保不同的HAR实现（标准、优化、懒加载、流式）可以统一使用，大大提高了代码的可复用性和灵活性。

## 最佳实践

### 选择合适的解析模式

- **小型HAR文件**: 使用标准解析即可
  ```go
  har, err := har.ParseFile("small.har")
  ```

- **大型HAR文件**: 使用内存优化模式减少内存占用
  ```go
  har, err := har.ParseFile("large.har", har.WithMemoryOptimized())
  ```

- **包含大型响应的HAR文件**: 使用懒加载延迟加载大型内容
  ```go
  har, err := har.ParseFile("large_content.har", har.WithLazyLoading())
  ```

- **超大HAR文件**: 使用流式解析逐条处理
  ```go
  iterator, err := har.NewStreamingParserFromFile("huge.har")
  ```

### 处理错误

始终检查解析错误，并考虑使用增强的错误处理：

```go
result, err := har.ParseHarFileWithWarnings("problematic.har")
if err != nil {
    log.Fatalf("解析完全失败: %v", err)
} else if len(result.Warnings) > 0 {
    log.Printf("解析成功，但有 %d 个警告", len(result.Warnings))
    for _, w := range result.Warnings {
        log.Printf("警告: %s", w.Error())
    }
}
```

### 使用接口进行编程

编写的函数应接受接口类型而非具体实现：

```go
// 好的做法 - 接受任何实现了HARProvider的类型
func ProcessHAR(harFile har.HARProvider) {
    // 处理逻辑
}

// 不推荐 - 只接受具体类型
func ProcessHAR(harFile *har.Har) {
    // 处理逻辑
}
```

## 性能考虑

- **内存使用**: 对于GB级HAR文件，标准解析可能会耗尽内存，考虑使用流式解析
- **解析速度**: 跳过验证可以提高解析速度，但可能错过格式问题
- **懒加载权衡**: 懒加载可以减少初始内存使用，但随后访问内容时可能有性能损失 