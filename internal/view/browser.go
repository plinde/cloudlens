package view

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/one2nc/cloudlens/internal/render"
	"github.com/one2nc/cloudlens/internal/ui"
)

// Browser represents a generic resource browser.
type Browser struct {
	*Table
	contextFn  ContextFunc
	cancelFn   context.CancelFunc
	refreshCtx context.Context
	mx         sync.RWMutex ``
}

// NewBrowser returns a new browser.
func NewBrowser(resource string) ResourceViewer {
	return &Browser{
		Table: NewTable(resource),
	}
}

// Init watches all running pods in given namespace.
func (b *Browser) Init(ctx context.Context) error {
	if err := b.Table.Init(ctx); err != nil {
		return err
	}
	b.SetContextFn(func(c context.Context) context.Context {
		return ctx
	})
	b.bindKeys(b.Actions())
	for _, f := range b.bindKeysFn {
		f(b.Actions())
	}

	row, _ := b.GetSelection()
	if row == 0 && b.GetRowCount() > 0 {
		b.Select(1, 0)
	}
	b.GetModel().SetRefreshRate(DefaultRefreshRate)
	return nil
}

func (b *Browser) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyR:       ui.NewSharedKeyAction("Refresh", b.resetCmd, false),
		ui.KeyQ:       ui.NewSharedKeyAction("Quit", b.backOrQuitCmd, false),
		tcell.KeyHelp: ui.NewSharedKeyAction("Help", b.helpCmd, false),
	})
}

// Start initializes browser updates.
func (b *Browser) Start() {
	b.Stop()
	b.GetModel().AddListener(b)
	b.Table.Start()
	b.refreshCtx = b.prepareContext()
	b.Table.GetModel().Refresh(b.refreshCtx)
	b.Refresh()
	if err := b.GetModel().Watch(b.refreshCtx); err != nil {
		b.App().Flash().Err(fmt.Errorf("Watcher failed for %s -- %w", b.Resource(), err))
	}
}

// Stop terminates browser updates.
func (b *Browser) Stop() {
	b.mx.Lock()
	{
		if b.cancelFn != nil {
			b.cancelFn()
			b.cancelFn = nil
		}
	}
	b.mx.Unlock()
	b.GetModel().RemoveListener(b)
	b.Table.Stop()
}

func (b *Browser) prepareContext() context.Context {
	ctx := context.Background()
	ctx, b.cancelFn = context.WithCancel(ctx)
	if b.contextFn != nil {
		ctx = b.contextFn(ctx)
	}
	return ctx
}

// Name returns the component name.
func (b *Browser) Name() string { return b.Table.Resource() }

// SetContextFn populates a custom context.
func (b *Browser) SetContextFn(f ContextFunc) { b.contextFn = f }

// GetTable returns the underlying table.
func (b *Browser) GetTable() *Table { return b.Table }

func (b *Browser) helpCmd(evt *tcell.EventKey) *tcell.EventKey {

	return evt
}

func (b *Browser) resetCmd(evt *tcell.EventKey) *tcell.EventKey {
	if b.refreshCtx != nil {
		b.GetModel().Refresh(b.refreshCtx)
		b.GetModel().ResetTimer()
	}
	return evt
}

func (b *Browser) backOrQuitCmd(evt *tcell.EventKey) *tcell.EventKey {
	if b.App().Content.IsLast() {
		b.App().BailOut()
	} else {
		b.App().PrevCmd(evt)
	}
	return nil
}

// TableDataChanged notifies the model data changed.
func (b *Browser) TableDataChanged(data *render.TableData) {
	b.App().Footer().SetLastSync(time.Now())
	b.App().QueueUpdateDraw(func() {
		b.Table.Update(data)
	})
}
