package render

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/one2nc/cloudlens/internal/aws"
)

type ASG struct {
}

func (a ASG) Header() Header {
	return Header{
		HeaderColumn{Name: "Name", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Min-Size", SortIndicatorIdx: -1, Align: tview.AlignCenter, Hide: false, Wide: false, MX: true, Time: false},
		HeaderColumn{Name: "Max-Size", SortIndicatorIdx: -1, Align: tview.AlignCenter, Hide: false, Wide: false, MX: true, Time: false},
		HeaderColumn{Name: "Desired", SortIndicatorIdx: -1, Align: tview.AlignCenter, Hide: false, Wide: false, MX: true, Time: false},
		HeaderColumn{Name: "Status", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Availability-Zones", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: true, MX: false, Time: false},
		HeaderColumn{Name: "Health-Check-Type", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: true, MX: false, Time: false},
		HeaderColumn{Name: "Default-Cooldown", SortIndicatorIdx: -1, Align: tview.AlignCenter, Hide: false, Wide: true, MX: false, Time: false},
		HeaderColumn{Name: "Launch-Config-Name", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: true, MX: false, Time: false},
		HeaderColumn{Name: "Termination-Policies", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: true, MX: false, Time: false},
	}
}

func (a ASG) Render(o interface{}, ns string, row *Row) error {
	asgResp, ok := o.(aws.ASGResp)

	if !ok {
		return fmt.Errorf("Expected ASGResp, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		asgResp.Name,
		asgResp.MinSize,
		asgResp.MaxSize,
		asgResp.DesiredCapacity,
		asgResp.Status,
		asgResp.AvailabilityZones,
		asgResp.HealthCheckType,
		asgResp.DefaultCooldown,
		asgResp.LaunchConfigurationName,
		asgResp.TerminationPolicies,
	}

	return nil
}
