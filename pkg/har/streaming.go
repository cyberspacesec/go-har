package har

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"
)

// Browser 表示HAR文件中的浏览器信息
type Browser struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Comment string `json:"comment,omitempty"`
}

// EntryIterator 提供流式迭代HAR文件中的条目的接口
type EntryIterator interface {
	// Next 移动到下一个条目，如果没有更多条目则返回false
	Next() bool
	// Entry 返回当前条目
	Entry() *Entries
	// Err 返回迭代过程中出现的错误
	Err() error
	// Close 关闭迭代器和相关资源
	Close() error
}

// StreamingHar 表示一个流式处理的HAR文件
// 它不会一次性加载整个HAR文件到内存中
type StreamingHar struct {
	file       *os.File
	fileOffset int64
	mutex      sync.Mutex
	creator    Creator
	pages      []Pages
	version    string
	data       []byte // 新增：保存原始字节数据
}

// StreamingEntryIterator 是HAR条目的迭代器
type StreamingEntryIterator struct {
	har        *StreamingHar
	decoder    *json.Decoder
	err        error
	file       *os.File      // 可能为nil，如果是基于内存的处理
	reader     *bytes.Reader // 新增：内存读取器
	currentPos int
	entry      Entries
	closed     bool
}

// NewStreamingHarFromFile 从文件路径创建一个流式HAR对象
func NewStreamingHarFromFile(filePath string) (*StreamingHar, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open HAR file: %w", err)
	}

	// 创建解码器
	decoder := json.NewDecoder(file)

	// 解析基本HAR信息，但不加载entries
	har := &StreamingHar{
		file: file,
	}

	// 查找HAR对象开始
	if err := findHarObjectStart(decoder); err != nil {
		file.Close()
		return nil, err
	}

	// 解析HAR基本信息
	if err := parseHarBasicInfo(decoder, har); err != nil {
		file.Close()
		return nil, err
	}

	// 记录当前文件位置
	offset, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to get file position: %w", err)
	}
	har.fileOffset = offset

	return har, nil
}

// NewStreamingHarFromBytes 从字节数据创建一个流式HAR对象
func NewStreamingHarFromBytes(data []byte) (*StreamingHar, error) {
	// 首先，为了提取基本信息，将整个HAR解析到一个临时对象中
	tempHar := &Har{}
	err := json.Unmarshal(data, tempHar)
	if err != nil {
		return nil, fmt.Errorf("无法解析HAR数据: %w", err)
	}

	// 创建HAR对象并填充基本信息
	har := &StreamingHar{
		data:    data,
		version: tempHar.Log.Version,
		creator: tempHar.Log.Creator,
		pages:   tempHar.Log.Pages,
	}

	return har, nil
}

// 查找HAR对象开始
func findHarObjectStart(decoder *json.Decoder) error {
	// 查找"{"字符
	token, err := decoder.Token()
	if err != nil {
		return fmt.Errorf("failed to read first token: %w", err)
	}
	if delim, ok := token.(json.Delim); !ok || delim != '{' {
		return errors.New("expected { at the start of HAR file")
	}

	// 查找"log"字段
	for {
		token, err := decoder.Token()
		if err != nil {
			return fmt.Errorf("failed to find log field: %w", err)
		}

		if str, ok := token.(string); ok && str == "log" {
			break
		}
	}

	// 检查log后面的是对象开始符号
	token, err = decoder.Token()
	if err != nil {
		return fmt.Errorf("failed to read token after log: %w", err)
	}
	if delim, ok := token.(json.Delim); !ok || delim != '{' {
		return errors.New("expected { after log field")
	}

	return nil
}

// 解析HAR基本信息
func parseHarBasicInfo(decoder *json.Decoder, har *StreamingHar) error {
	for {
		// 获取下一个字段名
		token, err := decoder.Token()
		if err != nil {
			return fmt.Errorf("failed to read field name: %w", err)
		}

		// 检查是否到达对象结束
		if delim, ok := token.(json.Delim); ok && delim == '}' {
			break
		}

		// 处理字段
		fieldName, ok := token.(string)
		if !ok {
			return fmt.Errorf("expected string field name, got %T", token)
		}

		switch fieldName {
		case "version":
			if err := decoder.Decode(&har.version); err != nil {
				return fmt.Errorf("failed to decode version: %w", err)
			}
		case "creator":
			if err := decoder.Decode(&har.creator); err != nil {
				return fmt.Errorf("failed to decode creator: %w", err)
			}
		case "browser":
			// 跳过browser字段，因为结构体中已移除
			var dummy interface{}
			if err := decoder.Decode(&dummy); err != nil {
				return fmt.Errorf("failed to skip browser field: %w", err)
			}
		case "pages":
			if err := decoder.Decode(&har.pages); err != nil {
				return fmt.Errorf("failed to decode pages: %w", err)
			}
		case "entries":
			// 只找到entries数组的开始，不解析内容
			token, err := decoder.Token()
			if err != nil {
				return fmt.Errorf("failed to find entries array start: %w", err)
			}
			if delim, ok := token.(json.Delim); !ok || delim != '[' {
				return errors.New("expected [ at the start of entries")
			}

			// 记录当前位置并退出
			return nil
		default:
			// 跳过其他字段
			var dummy interface{}
			if err := decoder.Decode(&dummy); err != nil {
				return fmt.Errorf("failed to skip field %s: %w", fieldName, err)
			}
		}
	}

	return nil
}

// Close 关闭StreamingHar并释放资源
func (h *StreamingHar) Close() error {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if h.file != nil {
		err := h.file.Close()
		h.file = nil
		return err
	}
	return nil
}

// GetVersion 返回HAR版本
func (h *StreamingHar) GetVersion() string {
	return h.version
}

// GetCreator 返回HAR创建者信息
func (h *StreamingHar) GetCreator() Creator {
	return h.creator
}

// GetPages 返回页面信息
func (h *StreamingHar) GetPages() []Pages {
	return h.pages
}

// Entries 返回一个条目迭代器
func (h *StreamingHar) Entries() *StreamingEntryIterator {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	// 通过直接解码JSON创建迭代器
	decoder := json.NewDecoder(bytes.NewReader(h.data))

	return &StreamingEntryIterator{
		har:     h,
		decoder: decoder,
		entry:   Entries{},
	}
}

// Next 获取下一个条目
func (it *StreamingEntryIterator) Next() bool {
	if it.closed || it.err != nil {
		return false
	}

	// 如果是第一次调用Next，需要找到entries数组位置
	if it.currentPos == 0 {
		// 简化的方法：直接搜索直到找到"entries"和"["
		// 使用更直接的方式查找entries数组
		found := false
		for !found {
			token, err := it.decoder.Token()
			if err != nil {
				it.err = err
				return false
			}

			// 检查是否是"entries"字段
			if str, ok := token.(string); ok && str == "entries" {
				// 下一个token应该是"["
				token, err = it.decoder.Token()
				if err != nil {
					it.err = err
					return false
				}

				if delim, ok := token.(json.Delim); ok && delim == '[' {
					found = true // 找到了entries数组
				} else {
					it.err = fmt.Errorf("预期在'entries'字段后找到'['，但实际为: %v", token)
					return false
				}
			}
		}
	}

	// 检查是否有更多元素
	if !it.decoder.More() {
		return false
	}

	// 解析下一个条目
	var entry Entries
	if err := it.decoder.Decode(&entry); err != nil {
		it.err = err
		return false
	}

	it.entry = entry
	it.currentPos++
	return true
}

// Entry 返回当前条目
func (it *StreamingEntryIterator) Entry() *Entries {
	return &it.entry
}

// Position 返回当前位置
func (it *StreamingEntryIterator) Position() int {
	return it.currentPos
}

// Err 返回迭代过程中的错误
func (it *StreamingEntryIterator) Err() error {
	if it.err == io.EOF {
		return nil
	}
	return it.err
}

// Close 关闭迭代器和相关资源
func (it *StreamingEntryIterator) Close() error {
	if it.closed {
		return nil
	}

	it.closed = true

	// 关闭文件（如果有）
	if it.file != nil {
		return it.file.Close()
	}

	// 对于内存读取器，不需要特殊关闭操作
	return nil
}

// GetAllEntries 获取所有条目（便捷方法，但会加载所有内容到内存）
func (sh *StreamingHar) GetAllEntries() ([]Entries, error) {
	var entries []Entries

	it := sh.Entries()
	for it.Next() {
		entries = append(entries, *it.Entry())
	}

	if err := it.Err(); err != nil {
		return entries, err
	}

	return entries, nil
}

// 使用示例：
//
// har, err := har.NewStreamingHarFromFile("large.har")
// if err != nil {
//     panic(err)
// }
//
// iterator := har.Entries()
// defer iterator.Close()
//
// for iterator.Next() {
//     entry := iterator.Entry()
//     fmt.Println(entry.Request.URL)
// }
//
// if err := iterator.Err(); err != nil {
//     panic(err)
// }
