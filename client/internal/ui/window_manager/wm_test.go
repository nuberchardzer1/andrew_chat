package wm

import (
	"andrew_chat/client/internal/ui/types"
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

type testWindow struct {
	width  int
	height int
}

func (d *testWindow) Init() tea.Cmd {
	return nil
}

func (d *testWindow) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if sizeMsg, ok := msg.(tea.WindowSizeMsg); ok {
		d.width = sizeMsg.Width
		d.height = sizeMsg.Height
	}
	return d, nil
}

func (d *testWindow) View() string {
	return ""
}

func TestAddWindow(t *testing.T) {
	wm := NewWM()
	tw := &testWindow{}

	wm.AddWindow(types.PositionTopLeft, tw, true)

	ow := wm.GetWindow(types.PositionTopLeft)
	if ow == nil {
		t.Fatalf("expected window at position %v, got nil", types.PositionTopLeft)
	}

	otw, ok := ow.(*testWindow)
	if !ok {
		t.Fatalf("expected *testWindow, got %T", ow)
	}

	if tw != otw {
		t.Fatalf("expected same instance, got different pointers")
	}
}


func TestCloseWindow(t *testing.T) {
	wm := NewWM()
	tw := &testWindow{}

	wm.AddWindow(types.PositionTopLeft, tw, true)
	wm.CloseWindow(tw)

	ow := wm.GetWindow(types.PositionTopLeft)
	if ow != nil {
		t.Fatalf("expected nil after close, got %T", ow)
	}
}

func TestGetWindowWidth(t *testing.T) {
	wm := NewWM()

	height := 500
	width := 950

	wm.Update(tea.WindowSizeMsg{
		Height: height,
		Width:  width,
	})

	tests := []struct {
		name     string
		position types.Position
		expected int
	}{
		{
			name:     "left panel",
			position: types.PositionTopLeft,
			expected: int(float32(width)*leftPanelWidthRatio) - types.BorderSize,
		},
		{
			name:     "right panel",
			position: types.PositionTopRight,
			expected: int(float32(width)*(1-leftPanelWidthRatio)) - types.BorderSize,
		},
	}

	for _, tt := range tests {
		got := wm.getWindowWidthWithoutBorder(tt.position)

		if got != tt.expected {
			t.Fatalf(
				"unexpected width: got=%d want=%d",
				got,
				tt.expected,
			)
		}
	}
}

func TestGetWindowHeight(t *testing.T) {

	wm := NewWM()

	height := 500
	width := 950

	wm.Update(tea.WindowSizeMsg{
		Height: height,
		Width:  width,
	})

	expected := height - types.BorderSize

	positionNames := map[types.Position]string{
		types.PositionTopLeft:  "TopLeft",
		types.PositionTopRight: "TopRight",
	}

	positions := []types.Position{
		types.PositionTopLeft,
		types.PositionTopRight,
	}

	for _, pos := range positions {
		name := positionNames[pos]
		got := wm.getWindowHeightWithoutBorder(pos)

		if got != expected {
			t.Fatalf(
				"unexpected height for %s: got=%d want=%d",
				name,
				got,
				expected,
			)
		}
	}
}

func runCmd(wm *WindowManager, cmd tea.Cmd) {
	for cmd != nil {
		msg := cmd()

		switch m := msg.(type) {

		case tea.BatchMsg:
			for _, c := range m {
				runCmd(wm, c)
			}
			return

		default:
			_, cmd = wm.Update(msg)
		}
	}
}

func TestAddWindowLayout(t *testing.T) {
	for i, test := range addLayoutTests {
		t.Run(
			fmt.Sprintf("case_%d_pos_%v", i, test.targetPos),
			func(t *testing.T) {

				wm := NewWM()

				_, cmd := wm.Update(tea.WindowSizeMsg{
					Height: test.h,
					Width:  test.w,
				})
				runCmd(wm, cmd)

				for _, pos := range test.existing {
					var tw testWindow
					cmd := wm.AddWindow(pos, &tw, false)
					runCmd(wm, cmd)
				}

				var target testWindow
				cmd = wm.AddWindow(test.targetPos, &target, true)
				runCmd(wm, cmd)

				outmodel := wm.GetWindow(test.targetPos)
				if outmodel == nil {
					t.Fatalf("target window not found at %v", test.targetPos)
				}

				otw, ok := outmodel.(*testWindow)
				if !ok {
					t.Fatalf("expected *testWindow, got %T", outmodel)
				}

				if otw != &target {
					t.Fatalf("expected same instance, got different pointer")
				}

				if target.height != test.expectH {
					t.Fatalf(
						"unexpected height: got=%d, want=%d (window=%v, totalH=%d)",
						target.height,
						test.expectH,
						test.targetPos,
						test.h,
					)
				}

				if target.width != test.expectW {
					t.Fatalf(
						"unexpected width: got=%d, want=%d (window=%v, totalW=%d)",
						target.width,
						test.expectW,
						test.targetPos,
						test.w,
					)
				}
			},
		)
	}
}


var addLayoutTests = []struct {
		name string

		h int
		w int

		existing []types.Position
		targetPos   types.Position

		expectH int
		expectW int
	}{
		{
			name: "single left window full height",

			h: 1000,
			w: 800,

			existing: []types.Position{
				types.PositionTopLeft,
			},

			targetPos: types.PositionTopLeft,

			expectH: 1000 - types.BorderSize,
			expectW: int(float32(800)*leftPanelWidthRatio) - types.BorderSize,
		},

		{
			name: "left column split",

			h: 1000,
			w: 800,

			existing: []types.Position{
				types.PositionTopLeft,
				types.PositionBotLeft,
			},

			targetPos: types.PositionTopLeft,

			expectH: (1000 / 2) - types.BorderSize,
			expectW: int(float32(800)*leftPanelWidthRatio) - types.BorderSize,
		},

		{
			name: "right column split",

			h: 1000,
			w: 800,

			existing: []types.Position{
				types.PositionTopRight,
				types.PositionBotRight,
			},

			targetPos: types.PositionBotRight,

			expectH: (1000 / 2) - types.BorderSize,
			expectW: int(float32(800)*(1-leftPanelWidthRatio)) - types.BorderSize,
		},

		{
			name: "four quadrants",

			h: 1000,
			w: 800,

			existing: []types.Position{
				types.PositionTopLeft,
				types.PositionBotLeft,
				types.PositionTopRight,
				types.PositionBotRight,
			},

			targetPos: types.PositionBotLeft,

			expectH: (1000 / 2) - types.BorderSize,
			expectW: int(float32(800)*leftPanelWidthRatio) - types.BorderSize,
		},

		{
			name: "right split does not affect left",

			h: 1000,
			w: 800,

			existing: []types.Position{
				types.PositionTopLeft,
				types.PositionTopRight,
				types.PositionBotRight,
			},

			targetPos: types.PositionTopLeft,

			expectH: 1000 - types.BorderSize,
			expectW: int(float32(800)*leftPanelWidthRatio) - types.BorderSize,
		},
	}