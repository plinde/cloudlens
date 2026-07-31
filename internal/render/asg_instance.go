package render

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/one2nc/cloudlens/internal/aws"
)

type ASGInstance struct {
}

func (a ASGInstance) Header() Header {
	return Header{
		HeaderColumn{Name: "Instance-Id", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Type", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Availability-Zone", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Lifecycle-State", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Health-Status", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Protected", SortIndicatorIdx: -1, Align: tview.AlignCenter, Hide: false, Wide: false, MX: false, Time: false},
	}
}

func (a ASGInstance) Render(o interface{}, ns string, row *Row) error {
	resp, ok := o.(aws.ASGInstanceResp)
	if !ok {
		return fmt.Errorf("Expected ASGInstanceResp, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		resp.InstanceId,
		resp.InstanceType,
		resp.AvailabilityZone,
		resp.LifecycleState,
		resp.HealthStatus,
		resp.Protected,
	}

	return nil
}
