package prompt

import (
	"github.com/d3code/xlog"
	"github.com/manifoldco/promptui"
)

func p() {

	y := promptui.Prompt{
		Label:       nil,
		Default:     "",
		AllowEdit:   false,
		Validate:    nil,
		Mask:        0,
		HideEntered: false,
		Templates:   nil,
		IsConfirm:   false,
		IsVimMode:   false,
		Pointer:     nil,
		Stdin:       nil,
		Stdout:      nil,
	}
	_, err := y.Run()
	if err != nil {
		return
	}

	x := promptui.Select{
		Label:             nil,
		Items:             []string{"1", "2"},
		Size:              0,
		CursorPos:         0,
		IsVimMode:         false,
		HideHelp:          false,
		HideSelected:      false,
		Templates:         nil,
		Keys:              nil,
		Searcher:          nil,
		StartInSearchMode: false,
		Pointer:           nil,
		Stdin:             nil,
		Stdout:            nil,
	}

	_, i, err := x.Run()
	if err != nil {
		xlog.Error(err.Error())
		return
	}

	xlog.Error(i)
}
