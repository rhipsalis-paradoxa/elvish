// Command widget allows manually testing a single widget.
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/elves/elvish/cli"
	"github.com/elves/elvish/cli/el"
	"github.com/elves/elvish/cli/el/codearea"
	"github.com/elves/elvish/cli/el/combobox"
	"github.com/elves/elvish/cli/el/listbox"
	"github.com/elves/elvish/cli/term"
	"github.com/elves/elvish/ui"
)

var (
	maxHeight  = flag.Int("max-height", 10, "maximum height")
	horizontal = flag.Bool("horizontal", false, "use horizontal listbox layout")
)

func makeWidget() el.Widget {
	items := listbox.TestItems{Prefix: "list item "}
	w := combobox.New(combobox.Spec{
		CodeArea: codearea.Spec{
			Prompt: codearea.ConstPrompt(
				ui.T(" NUMBER ", ui.Bold, ui.BgMagenta).
					ConcatText(ui.T(" "))),
		},
		ListBox: listbox.Spec{
			State:       listbox.State{Items: &items},
			Placeholder: ui.T("(no items)"),
			Horizontal:  *horizontal,
		},
		OnFilter: func(w combobox.Widget, filter string) {
			if n, err := strconv.Atoi(filter); err == nil {
				items.NItems = n
			}
		},
	})
	return w
}

func main() {
	flag.Parse()
	widget := makeWidget()

	tty := cli.NewTTY(os.Stdin, os.Stderr)
	restore, err := tty.Setup()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer restore()
	events := tty.StartInput()
	defer tty.StopInput()
	for {
		h, w := tty.Size()
		if h > *maxHeight {
			h = *maxHeight
		}
		tty.UpdateBuffer(nil, widget.Render(w, h), false)
		event := <-events
		handled := widget.Handle(event)
		if !handled && event == term.K('D', ui.Ctrl) {
			tty.UpdateBuffer(nil, term.NewBufferBuilder(w).Buffer(), true)
			break
		}
	}
}