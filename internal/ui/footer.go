package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/derailed/tview"
	"github.com/one2nc/cloudlens/internal/model"
)

// Footer represents the bottom status bar with hotkeys and last-sync timestamp.
type Footer struct {
	*tview.TextView

	lastSync time.Time
	hints    model.MenuHints
}

// NewFooter returns a new footer view.
func NewFooter() *Footer {
	f := Footer{
		TextView: tview.NewTextView(),
	}
	f.SetTextAlign(tview.AlignLeft)
	f.SetDynamicColors(true)
	f.SetText("[white::-]Ready")
	return &f
}

// SetLastSync records the last refresh time and updates the display.
func (f *Footer) SetLastSync(t time.Time) {
	f.lastSync = t
	f.refresh()
}

// UpdateHints sets the hotkey hints and updates the display.
func (f *Footer) UpdateHints(hh model.MenuHints) {
	f.hints = hh
	f.refresh()
}

// StackPushed notifies a component was added.
func (f *Footer) StackPushed(c model.Component) {
	f.UpdateHints(c.Hints())
}

// StackPopped notifies a component was removed.
func (f *Footer) StackPopped(_, top model.Component) {
	if top != nil {
		f.UpdateHints(top.Hints())
	} else {
		f.hints = nil
		f.refresh()
	}
}

// StackTop notifies the top component.
func (f *Footer) StackTop(t model.Component) {
	f.UpdateHints(t.Hints())
}

// Name returns the component name.
func (f *Footer) Name() string { return "footer" }

// Hints returns menu hints (satisfies Hinter interface).
func (f *Footer) Hints() model.MenuHints { return nil }

func (f *Footer) refresh() {
	var parts []string

	// Hotkeys section — format visible hints in e2c style
	if len(f.hints) > 0 {
		var hotkeys []string
		for _, h := range f.hints {
			if !h.Visible || h.Mnemonic == "" || h.Description == "" {
				continue
			}
			mnemonic := "<" + strings.ToLower(h.Mnemonic) + ">"
			hotkeys = append(hotkeys, fmt.Sprintf("[yellow::b]%s[white::-] %s", mnemonic, h.Description))
		}
		if len(hotkeys) > 0 {
			parts = append(parts, strings.Join(hotkeys, "  "))
		}
	}

	// Last-sync section
	if !f.lastSync.IsZero() {
		syncPart := fmt.Sprintf("[silver::]Last sync: [white::]%s", f.lastSync.Format("15:04:05"))
		parts = append(parts, syncPart)
	}

	if len(parts) == 0 {
		f.SetText("[white::-]Ready")
		return
	}

	f.SetText(" " + strings.Join(parts, " [silver::]│ ") + " ")
}
