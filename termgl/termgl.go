package termgl

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/muesli/ansi"
	"golang.org/x/image/font"
	"image"
	"image/color"
	"io"
	"strings"
	"sync"
)

var mappedKeys = map[ebiten.Key]tea.KeyType{
	ebiten.KeyEnter:      tea.KeyEnter,
	ebiten.KeyTab:        tea.KeyTab,
	ebiten.KeySpace:      tea.KeySpace,
	ebiten.KeyBackspace:  tea.KeyBackspace,
	ebiten.KeyDelete:     tea.KeyDelete,
	ebiten.KeyHome:       tea.KeyHome,
	ebiten.KeyEnd:        tea.KeyEnd,
	ebiten.KeyPageUp:     tea.KeyPgUp,
	ebiten.KeyArrowUp:    tea.KeyUp,
	ebiten.KeyArrowDown:  tea.KeyDown,
	ebiten.KeyArrowLeft:  tea.KeyLeft,
	ebiten.KeyArrowRight: tea.KeyRight,
	ebiten.KeyEscape:     tea.KeyEscape,
}

var mappedRunes = map[ebiten.Key][]rune{
	ebiten.Key1: {'1'},
	ebiten.Key2: {'2'},
	ebiten.Key3: {'3'},
	ebiten.Key4: {'4'},
	ebiten.Key5: {'5'},
	ebiten.Key6: {'6'},
	ebiten.Key7: {'7'},
	ebiten.Key8: {'8'},
	ebiten.Key9: {'9'},
	ebiten.Key0: {'0'},
}

type FontWeight byte

const (
	FontWeightNormal FontWeight = iota
	FontWeightBold
	FontWeightItalic
)

type GridCell struct {
	Char   rune
	Fg     color.Color
	Bg     color.Color
	Weight FontWeight
}

type Game struct {
	sync.Mutex

	cellsWidth  int
	cellsHeight int
	cellWidth   int
	cellHeight  int
	cellOffsetY int

	prog *tea.Program

	faceNormal font.Face
	faceBold   font.Face
	faceItalic font.Face

	tty      io.Reader
	grid     [][]GridCell
	bgColors *image.RGBA

	cursorX int
	cursorY int

	mouseCellX int
	mouseCellY int

	defaultBg color.Color
	curFg     color.Color
	curBg     color.Color
	curWeight FontWeight

	routine sync.Once
}

// NewGame creates a new terminal game with the given dimensions and font faces.
func NewGame(width int, height int, fontNormal font.Face, fontBold font.Face, fontItalic font.Face, tty io.Reader, prog *tea.Program, defaultBg color.Color) *Game {
	if defaultBg == nil {
		defaultBg = color.Black
	}

	bounds, _, _ := fontNormal.GlyphBounds([]rune("█")[0])
	size := bounds.Max.Sub(bounds.Min)

	cellWidth := size.X.Round()
	cellHeight := size.Y.Round()
	cellOffsetY := -bounds.Min.Y.Round()

	grid := make([][]GridCell, height)
	for y := 0; y < height; y++ {
		grid[y] = make([]GridCell, width)
		for x := 0; x < width; x++ {
			grid[y][x] = GridCell{
				Char:   ' ',
				Fg:     color.White,
				Bg:     defaultBg,
				Weight: FontWeightNormal,
			}
		}
	}

	game := &Game{
		cellsWidth:  width,
		cellsHeight: height,
		cellWidth:   cellWidth,
		cellHeight:  cellHeight,
		cellOffsetY: cellOffsetY,
		prog:        prog,
		faceNormal:  fontNormal,
		faceBold:    fontBold,
		faceItalic:  fontItalic,
		defaultBg:   defaultBg,
		grid:        grid,
		tty:         tty,
		bgColors:    image.NewRGBA(image.Rect(0, 0, width*cellWidth, height*cellHeight)),
	}

	game.ResetSGR()

	return game
}

// ResetSGR resets the SGR attributes to their default values.
func (g *Game) ResetSGR() {
	g.curFg = color.White
	g.curBg = g.defaultBg
	g.curWeight = FontWeightNormal
}

// SetBgPixels sets a chunk of background pixels in the size of the cell.
func (g *Game) SetBgPixels(x, y int, c color.Color) {
	for i := 0; i < g.cellWidth; i++ {
		for j := 0; j < g.cellHeight; j++ {
			g.bgColors.Set(x*g.cellWidth+i, y*g.cellHeight+j, c)
		}
	}
}

func (g *Game) handleCSI(csi any) {
	switch seq := csi.(type) {
	case CursorUpSeq:
		g.cursorY -= seq.Count
		if g.cursorY < 0 {
			g.cursorY = 0
		}
	case CursorDownSeq:
		g.cursorY += seq.Count
		if g.cursorY >= g.cellsHeight {
			g.cursorY = g.cellsHeight - 1
		}
	case CursorForwardSeq:
		g.cursorX += seq.Count
		if g.cursorX >= g.cellsWidth {
			g.cursorX = g.cellsWidth - 1
		}
	case CursorBackSeq:
		g.cursorX -= seq.Count
		if g.cursorX < 0 {
			g.cursorX = 0
		}
	case CursorNextLineSeq:
		g.cursorY += seq.Count
		if g.cursorY >= g.cellsHeight {
			g.cursorY = g.cellsHeight - 1
		}
		g.cursorX = 0
	case CursorPreviousLineSeq:
		g.cursorY -= seq.Count
		if g.cursorY < 0 {
			g.cursorY = 0
		}
		g.cursorX = 0
	case CursorHorizontalSeq:
		g.cursorX = seq.Count - 1
	case CursorPositionSeq:
		g.cursorX = seq.Col - 1
		g.cursorY = seq.Row - 1

		if g.cursorX < 0 {
			g.cursorX = 0
		} else if g.cursorX >= g.cellsWidth {
			g.cursorX = g.cellsWidth - 1
		}

		if g.cursorY < 0 {
			g.cursorY = 0
		} else if g.cursorY >= g.cellsHeight {
			g.cursorY = g.cellsHeight - 1
		}
	case EraseDisplaySeq:
		if seq.Type != 2 {
			return // only support 2 (erase entire display)
		}

		for i := 0; i < g.cellsWidth; i++ {
			for j := 0; j < g.cellsHeight; j++ {
				g.grid[j][i].Char = ' '
				g.grid[j][i].Fg = color.White
				g.grid[j][i].Bg = g.defaultBg
			}
		}
	case EraseLineSeq:
		switch seq.Type {
		case 0: // erase from cursor to end of line
			for i := g.cursorX; i < g.cellsWidth-g.cursorX; i++ {
				g.grid[g.cursorY][g.cursorX+i].Char = ' '
				g.grid[g.cursorY][g.cursorX+i].Fg = color.White
				g.grid[g.cursorY][g.cursorX+i].Bg = g.defaultBg
				g.SetBgPixels(g.cursorX+i, g.cursorY, g.defaultBg)
			}
		case 1: // erase from start of line to cursor
			for i := 0; i < g.cursorX; i++ {
				g.grid[g.cursorY][i].Char = ' '
				g.grid[g.cursorY][i].Fg = color.White
				g.grid[g.cursorY][i].Bg = g.defaultBg
				g.SetBgPixels(i, g.cursorY, g.defaultBg)
			}
		case 2: // erase entire line
			for i := 0; i < g.cellsWidth; i++ {
				g.grid[g.cursorY][i].Char = ' '
				g.grid[g.cursorY][i].Fg = color.White
				g.grid[g.cursorY][i].Bg = g.defaultBg
				g.SetBgPixels(i, g.cursorY, g.defaultBg)
			}
		}
	case ScrollUpSeq:
		fmt.Println("UNSUPPORTED: ScrollUpSeq", seq.Count)
	case ScrollDownSeq:
		fmt.Println("UNSUPPORTED: ScrollDownSeq", seq.Count)
	case SaveCursorPositionSeq:
		fmt.Println("UNSUPPORTED: SaveCursorPositionSeq")
	case RestoreCursorPositionSeq:
		fmt.Println("UNSUPPORTED: RestoreCursorPositionSeq")
	case ChangeScrollingRegionSeq:
		fmt.Println("UNSUPPORTED: ChangeScrollingRegionSeq")
	case InsertLineSeq:
		fmt.Println("UNSUPPORTED: InsertLineSeq")
	case DeleteLineSeq:
		fmt.Println("UNSUPPORTED: DeleteLineSeq")
	}
}

func (g *Game) handleSGR(sgr any) {
	switch seq := sgr.(type) {
	case SGRReset:
		g.ResetSGR()
	case SGRBold:
		g.curWeight = FontWeightBold
	case SGRItalic:
		g.curWeight = FontWeightItalic
	case SGRUnsetBold:
		g.curWeight = FontWeightNormal
	case SGRUnsetItalic:
		g.curWeight = FontWeightNormal
	case SGRFgTrueColor:
		g.curFg = color.RGBA{seq.R, seq.G, seq.B, 255}
	case SGRBgTrueColor:
		g.curBg = color.RGBA{seq.R, seq.G, seq.B, 255}
	}
}

func (g *Game) parseSequences(str string, printExtra bool) int {
	runes := []rune(str)

	lastFound := 0
	for i := 0; i < len(runes); i++ {
		if sgr, ok := extractSGR(string(runes[i:])); ok {
			i += len(sgr) - 1

			if sgr, ok := parseSGR(sgr); ok {
				lastFound = i
				for i := range sgr {
					g.handleSGR(sgr[i])
				}
			}
		} else if csi, ok := extractCSI(string(runes[i:])); ok {
			i += len(csi) - 1

			if csi, ok := parseCSI(csi); ok {
				lastFound = i
				g.handleCSI(csi)
			}
		} else if printExtra {
			g.PrintChar(runes[i], g.curFg, g.curBg, g.curWeight)
		}
	}
	return lastFound
}

func (g *Game) RecalculateBackgrounds() {
	for i := 0; i < g.cellsWidth; i++ {
		for j := 0; j < g.cellsHeight; j++ {
			g.SetBgPixels(i, j, g.grid[j][i].Bg)
		}
	}
}

// PrintChar prints a character to the screen.
func (g *Game) PrintChar(r rune, fg, bg color.Color, weight FontWeight) {
	if r == '\n' {
		g.cursorX = 0
		g.cursorY++
		return
	}

	if ansi.PrintableRuneWidth(string(r)) == 0 {
		return
	}

	// Wrap around if we're at the end of the line.
	if g.cursorX >= g.cellsWidth {
		g.cursorX = 0
		g.cursorY++
	}

	// Scroll down if we're at the bottom and add a new line.
	if g.cursorY >= g.cellsHeight {
		diff := g.cursorY - g.cellsHeight + 1
		g.grid = g.grid[diff:]
		for i := 0; i < diff; i++ {
			g.grid = append(g.grid, make([]GridCell, g.cellsWidth))
			for i := 0; i < g.cellsWidth; i++ {
				g.grid[len(g.grid)-1][i].Char = ' '
				g.grid[len(g.grid)-1][i].Fg = color.White
				g.grid[len(g.grid)-1][i].Bg = g.defaultBg
			}
		}
		g.cursorY = g.cellsHeight - 1
		g.RecalculateBackgrounds()
	}

	// Set the cell.
	g.grid[g.cursorY][g.cursorX].Char = r
	g.grid[g.cursorY][g.cursorX].Fg = fg
	g.grid[g.cursorY][g.cursorX].Bg = bg
	g.grid[g.cursorY][g.cursorX].Weight = weight

	// Set the pixels.
	g.SetBgPixels(g.cursorX, g.cursorY, g.grid[g.cursorY][g.cursorX].Bg)

	// Move the cursor.
	g.cursorX++
}

func (g *Game) Update() error {
	g.routine.Do(func() {
		go func() {
			buf := make([]byte, 1024)
			for {
				n, err := g.tty.Read(buf)
				if err != nil {
					fmt.Println("ERROR: ", err)
					continue
				}

				if n == 0 {
					continue
				}

				g.Lock()
				{
					line := string(buf[:n])
					g.parseSequences(line, true)
				}
				g.Unlock()
			}
		}()
	})

	mx, my := ebiten.CursorPosition()
	mcx, mcy := mx/g.cellWidth, my/g.cellHeight

	if mcx != g.mouseCellX || mcy != g.mouseCellY {
		g.mouseCellX = mcx
		g.mouseCellY = mcy

		g.prog.Send(tea.MouseMsg{
			X:      g.mouseCellX,
			Y:      g.mouseCellY,
			Shift:  false,
			Alt:    false,
			Ctrl:   false,
			Action: tea.MouseActionMotion,
			Button: 0,
			Type:   tea.MouseMotion,
		})
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.prog.Send(tea.MouseMsg{
			X:      g.mouseCellX,
			Y:      g.mouseCellY,
			Shift:  false,
			Alt:    false,
			Ctrl:   false,
			Action: tea.MouseActionRelease,
			Button: tea.MouseButtonLeft,
			Type:   tea.MouseLeft,
		})
	}

	var keys []ebiten.Key
	keys = inpututil.AppendJustReleasedKeys(keys)

	for _, k := range keys {
		runes := []rune(strings.ToLower(k.String()))
		if val, ok := mappedRunes[k]; ok {
			runes = val
		}
		g.prog.Send(tea.KeyMsg{
			Type:  tea.KeyRunes,
			Runes: runes,
			Alt:   ebiten.IsKeyPressed(ebiten.KeyAlt),
		})
	}

	for k, v := range mappedKeys {
		if inpututil.IsKeyJustReleased(k) {
			runes := []rune(strings.ToLower(k.String()))
			g.prog.Send(tea.KeyMsg{
				Type:  v,
				Runes: runes,
				Alt:   ebiten.IsKeyPressed(ebiten.KeyAlt),
			})
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Lock()
	defer g.Unlock()

	// Draw background
	screen.WritePixels(g.bgColors.Pix)

	// Draw text
	for y := 0; y < g.cellsHeight; y++ {
		for x := 0; x < g.cellsWidth; x++ {
			if g.grid[y][x].Char == ' ' {
				continue
			}

			switch g.grid[y][x].Weight {
			case FontWeightNormal:
				text.Draw(screen, string(g.grid[y][x].Char), g.faceNormal, x*g.cellWidth, y*g.cellHeight+g.cellOffsetY, g.grid[y][x].Fg)
			case FontWeightBold:
				text.Draw(screen, string(g.grid[y][x].Char), g.faceBold, x*g.cellWidth, y*g.cellHeight+g.cellOffsetY, g.grid[y][x].Fg)
			case FontWeightItalic:
				text.Draw(screen, string(g.grid[y][x].Char), g.faceItalic, x*g.cellWidth, y*g.cellHeight+g.cellOffsetY, g.grid[y][x].Fg)
			}
		}
	}

	/*
		screen.Set(g.cursorX*cellWidth, g.cursorY*cellHeight, color.White)
		screen.Set(g.cursorX*cellWidth, g.cursorY*cellHeight+cellHeight, color.White)
		screen.Set(g.cursorX*cellWidth+cellWidth, g.cursorY*cellHeight, color.White)
		screen.Set(g.cursorX*cellWidth+cellWidth, g.cursorY*cellHeight+cellHeight, color.White)
	*/

	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.cellsWidth * g.cellWidth, g.cellsHeight * g.cellHeight
}
