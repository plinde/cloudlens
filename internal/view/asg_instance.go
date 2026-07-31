package view

import (
	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal/ui"
)

type ASGInstanceView struct {
	ResourceViewer
}

func NewASGInstance(resource string) ResourceViewer {
	var a ASGInstanceView
	a.ResourceViewer = NewBrowser(resource)
	a.AddBindKeysFn(a.bindKeys)
	return &a
}

func (a *ASGInstanceView) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftI:    ui.NewKeyAction("Sort Instance-Id", a.GetTable().SortColCmd("Instance-Id", true), true),
		ui.KeyShiftS:    ui.NewKeyAction("Sort Lifecycle-State", a.GetTable().SortColCmd("Lifecycle-State", true), true),
		tcell.KeyEscape: ui.NewKeyAction("Back", a.App().PrevCmd, false),
		tcell.KeyEnter:  ui.NewKeyAction("View", a.enterCmd, false),
	})
}

func (a *ASGInstanceView) enterCmd(evt *tcell.EventKey) *tcell.EventKey {
	instanceId := a.GetTable().GetSelectedItem()
	if instanceId != "" {
		f := describeResource
		if a.GetTable().enterFn != nil {
			f = a.GetTable().enterFn
		}
		f(a.App(), a.GetTable().GetModel(), a.Resource(), instanceId)
		a.App().Flash().Info("Instance: " + instanceId)
	}
	return nil
}
