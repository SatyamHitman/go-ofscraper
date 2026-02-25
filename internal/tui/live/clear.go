// =============================================================================
// FILE: internal/tui/live/clear.go
// PURPOSE: Terminal clear operations. Provides functions to clear the terminal
//          screen or specific lines for live display updates.
//          Ports Python utils/live/clear.py.
// =============================================================================

package live

import (
	"fmt"
	"os"
	"strings"
)

// ---------------------------------------------------------------------------
// ANSI escape sequences
// ---------------------------------------------------------------------------

const (
	ansiClearScreen   = "\033[2J"
	ansiMoveCursorTop = "\033[H"
	ansiClearLine     = "\033[2K"
	ansiMoveUp        = "\033[%dA"
	ansiMoveToCol0    = "\033[0G"
)

// ---------------------------------------------------------------------------
// ClearScreen
// ---------------------------------------------------------------------------

// ClearScreen provides terminal clear operations.
type ClearScreen struct {
	lastLineCount int
}

// NewClearScreen creates a new ClearScreen.
func NewClearScreen() *ClearScreen {
	return &ClearScreen{}
}

// ClearAll clears the entire terminal screen and moves the cursor to the top.
func (cs *ClearScreen) ClearAll() {
	fmt.Fprint(os.Stdout, ansiClearScreen+ansiMoveCursorTop)
	cs.lastLineCount = 0
}

// ClearLines clears the last N lines of output by moving up and clearing.
func (cs *ClearScreen) ClearLines(n int) {
	if n <= 0 {
		return
	}
	var sb strings.Builder
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf(ansiMoveUp, 1))
		sb.WriteString(ansiMoveToCol0)
		sb.WriteString(ansiClearLine)
	}
	fmt.Fprint(os.Stdout, sb.String())
}

// ClearPrevious clears the lines from the last render pass.
func (cs *ClearScreen) ClearPrevious() {
	cs.ClearLines(cs.lastLineCount)
}

// SetLastLineCount records how many lines were rendered, so they can be
// cleared on the next update.
func (cs *ClearScreen) SetLastLineCount(n int) {
	cs.lastLineCount = n
}

// LastLineCount returns the recorded line count.
func (cs *ClearScreen) LastLineCount() int {
	return cs.lastLineCount
}

// CountLines counts the number of newline-delimited lines in a string.
func CountLines(s string) int {
	if s == "" {
		return 0
	}
	n := 1
	for _, c := range s {
		if c == '\n' {
			n++
		}
	}
	return n
}
