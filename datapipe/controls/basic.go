package controls

import (
	"github.com/plusev-terminal/go-plugin-common/datapipe/types"
)

type Control struct {
	Label string
	Name  string
	Type  types.GuiControlType
}

func NewControl(label, name string, controlType types.GuiControlType, options map[string]any) *types.GuiControl {
	return &types.GuiControl{
		Label:   label,
		Name:    name,
		Type:    controlType,
		Options: options,
	}
}
