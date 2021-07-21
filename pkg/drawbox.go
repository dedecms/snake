package pkg

// Box9Slice is used by DrawBox functions to draw frames around text content by
// defining the corner and edge characters. See DefaultBox9Slice for an
// example
type Box9Slice struct {
	Top         string
	TopRight    string
	Right       string
	BottomRight string
	Bottom      string
	BottomLeft  string
	Left        string
	TopLeft     string
}

var defaultBox9Slice = Box9Slice{
	Top:         "─",
	TopRight:    "┐",
	Right:       "│",
	BottomRight: "┘",
	Bottom:      "─",
	BottomLeft:  "└",
	Left:        "│",
	TopLeft:     "┌",
}

// DefaultBox9Slice defines the character object to use with "CustomBox".
// It is used as Box9Slice object in "DrawBox" function.
//
// Usage:
// DrawCustomBox("Hello World", 20, AligntTypeCenter, DefaultBox9Slice())
//
// Outputs:
// <code>
//   ┌──────────────────┐
//   │   Hello World    │
//   └──────────────────┘
// </code>
func DefaultBox9Slice() Box9Slice {
	return defaultBox9Slice
}

var simpleBox9Slice = Box9Slice{
	Top:         "-",
	TopRight:    "+",
	Right:       "|",
	BottomRight: "+",
	Bottom:      "-",
	BottomLeft:  "+",
	Left:        "|",
	TopLeft:     "+",
}

// SimpleBox9Slice defines a character set to use with DrawCustomBox. It uses
// only simple ASCII characters
//
// Usage:
//   DrawCustomBox("Hello World", 20, Center, SimpleBox9Slice(), "\n")
//
// Outputs:
//   +------------------+
//   |   Hello World    |
//   +------------------+
func SimpleBox9Slice() Box9Slice {
	return simpleBox9Slice
}

// // DrawBox creates a frame with "content" in it. DefaultBox9Slice object is used to
// // define characters in the frame. "align" sets the alignment of the content.
// // It must be one of the strutil.AlignType constants.
// //
// // Usage:
// //   DrawBox("Hello World", 20, Center)
// //
// // Outputs:
// //   ┌──────────────────┐
// //   │   Hello World    │
// //   └──────────────────┘
// func DrawBox(content string, width int, align AlignType) (string, error) {
// 	return DrawCustomBox(content, width, align, defaultBox9Slice, "\n")
// }
