package view

import (
	"fmt"
	"strconv"

	"github.com/derailed/tview"
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal/dao"
	"github.com/one2nc/cloudlens/internal/ui"
	"github.com/one2nc/cloudlens/internal/ui/dialog"
)

type ASG struct {
	ResourceViewer
}

func NewASG(resource string) ResourceViewer {
	var a ASG
	a.ResourceViewer = NewBrowser(resource)
	a.AddBindKeysFn(a.bindKeys)
	return &a
}

func (a *ASG) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftN:    ui.NewKeyAction("Sort Name", a.GetTable().SortColCmd("Name", true), true),
		ui.KeyShiftM:    ui.NewKeyAction("Sort Min-Size", a.GetTable().SortColCmd("Min-Size", true), true),
		ui.KeyShiftX:    ui.NewKeyAction("Sort Max-Size", a.GetTable().SortColCmd("Max-Size", true), true),
		ui.KeyShiftD:    ui.NewKeyAction("Sort Desired", a.GetTable().SortColCmd("Desired", true), true),
		ui.KeyShiftS:    ui.NewKeyAction("Sort Status", a.GetTable().SortColCmd("Status", true), true),
		ui.KeyE:         ui.NewKeyAction("Edit", a.editCmd, true),
		tcell.KeyEscape: ui.NewKeyAction("Back", a.App().PrevCmd, false),
		tcell.KeyEnter:  ui.NewKeyAction("View", a.enterCmd, false),
	})
}

func (a *ASG) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	asgName := a.GetTable().GetSelectedItem()
	if asgName != "" {
		f := describeResource
		if a.GetTable().enterFn != nil {
			f = a.GetTable().enterFn
		}
		f(a.App(), a.GetTable().GetModel(), a.Resource(), asgName)
		a.App().Flash().Info("ASG: " + asgName)
	}
	return nil
}

func (a *ASG) editCmd(evt *tcell.EventKey) *tcell.EventKey {
	asgName := a.GetTable().GetSelectedItem()
	if asgName == "" {
		return nil
	}
	row := a.GetTable().GetSelectedRowIndex()
	if row == 0 {
		return nil
	}
	currentMin := a.GetTable().GetSelectedCell(1)
	currentMax := a.GetTable().GetSelectedCell(2)
	currentDesired := a.GetTable().GetSelectedCell(3)

	form := tview.NewForm()
	form.SetItemPadding(0)
	form.SetButtonsAlign(tview.AlignCenter).
		SetButtonBackgroundColor(tcell.ColorDarkSlateBlue).
		SetButtonTextColor(tcell.ColorBlack.TrueColor()).
		SetLabelColor(tcell.ColorWhite.TrueColor()).
		SetFieldTextColor(tcell.ColorIndianRed)

	form.AddInputField("Min Size", currentMin, 10, nil, nil)
	form.AddInputField("Max Size", currentMax, 10, nil, nil)
	form.AddInputField("Desired Capacity", currentDesired, 10, nil, nil)

	form.AddButton("Apply", func() {
		a.handleApply(form, asgName)
	})
	form.AddButton("Cancel", func() {
		a.App().Content.Pages.RemovePage("asg-edit")
	})

	form.SetFocus(0)
	if b := form.GetButton(0); b != nil {
		b.SetBackgroundColorActivated(tcell.ColorDodgerBlue)
		b.SetLabelColorActivated(tcell.ColorBlack.TrueColor())
	}

	modal := tview.NewModalForm(fmt.Sprintf("<ASG Edit> %s", asgName), form)
	modal.SetTextColor(tcell.ColorOrangeRed)
	modal.SetBackgroundColor(tcell.ColorBlack.TrueColor())
	modal.SetBorderColor(tcell.ColorBlue)
	modal.SetDoneFunc(func(int, string) {
		a.App().Content.Pages.RemovePage("asg-edit")
	})

	a.App().Content.Pages.AddPage("asg-edit", modal, false, true)
	return nil
}

func (a *ASG) handleApply(form *tview.Form, asgName string) {
	minField := form.GetFormItem(0).(*tview.InputField)
	maxField := form.GetFormItem(1).(*tview.InputField)
	desiredField := form.GetFormItem(2).(*tview.InputField)

	minSize, err := strconv.Atoi(minField.GetText())
	if err != nil {
		dialog.ShowError(a.App().Content.Pages, "Min Size must be a valid integer")
		return
	}
	maxSize, err := strconv.Atoi(maxField.GetText())
	if err != nil {
		dialog.ShowError(a.App().Content.Pages, "Max Size must be a valid integer")
		return
	}
	desiredCapacity, err := strconv.Atoi(desiredField.GetText())
	if err != nil {
		dialog.ShowError(a.App().Content.Pages, "Desired Capacity must be a valid integer")
		return
	}
	if minSize > maxSize {
		dialog.ShowError(a.App().Content.Pages, "Min Size cannot be larger than Max Size")
		return
	}
	if desiredCapacity < minSize || desiredCapacity > maxSize {
		dialog.ShowError(a.App().Content.Pages, fmt.Sprintf("Desired Capacity must be between %d and %d", minSize, maxSize))
		return
	}

	asgDAO := &dao.ASG{}
	asgDAO.Init(a.App().GetContext())
	err = asgDAO.UpdateSize(a.App().GetContext(), asgName, int32(minSize), int32(maxSize), int32(desiredCapacity))
	if err != nil {
		dialog.ShowError(a.App().Content.Pages, "Update failed: "+err.Error())
		return
	}

	a.App().Content.Pages.RemovePage("asg-edit")
	a.App().Flash().Info(fmt.Sprintf("ASG %s updated (min=%d, max=%d, desired=%d)", asgName, minSize, maxSize, desiredCapacity))
	a.GetTable().GetModel().Refresh(a.App().GetContext())
}
