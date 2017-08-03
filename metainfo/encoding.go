package metainfo

import "github.com/axgle/mahonia"

// ToUtf8 converts the content to the encoding of UTF8 from the encoding of from.
//
// The charset of content is 'from', and the result is UTF-8.
func ToUtf8(from, content string) string {
	dec := mahonia.NewDecoder(from)
	if dec == nil {
		return ""
	}
	return dec.ConvertString(content)
}
