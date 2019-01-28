package inputhandler

import (
	"fmt"
	"strings"

	"github.com/hajimehoshi/ebiten/inpututil"

	"github.com/hajimehoshi/ebiten"
)

func KeyPressed(name string) bool {
	key, ok := KeyNameToKey(name)
	if !ok {
		fmt.Println("no such key", name)
		return false
	}
	return ebiten.IsKeyPressed(key)
}

func KeyJustPressed(name string) bool {
	key, ok := KeyNameToKey(name)
	if !ok {
		fmt.Println("no such key", name)
		return false
	}
	return inpututil.IsKeyJustPressed(key)
}

func KeyNameToKey(name string) (ebiten.Key, bool) {
	switch strings.ToLower(name) {
	case "0":
		return ebiten.Key0, true
	case "1":
		return ebiten.Key1, true
	case "2":
		return ebiten.Key2, true
	case "3":
		return ebiten.Key3, true
	case "4":
		return ebiten.Key4, true
	case "5":
		return ebiten.Key5, true
	case "6":
		return ebiten.Key6, true
	case "7":
		return ebiten.Key7, true
	case "8":
		return ebiten.Key8, true
	case "9":
		return ebiten.Key9, true
	case "a":
		return ebiten.KeyA, true
	case "b":
		return ebiten.KeyB, true
	case "c":
		return ebiten.KeyC, true
	case "d":
		return ebiten.KeyD, true
	case "e":
		return ebiten.KeyE, true
	case "f":
		return ebiten.KeyF, true
	case "g":
		return ebiten.KeyG, true
	case "h":
		return ebiten.KeyH, true
	case "i":
		return ebiten.KeyI, true
	case "j":
		return ebiten.KeyJ, true
	case "k":
		return ebiten.KeyK, true
	case "l":
		return ebiten.KeyL, true
	case "m":
		return ebiten.KeyM, true
	case "n":
		return ebiten.KeyN, true
	case "o":
		return ebiten.KeyO, true
	case "p":
		return ebiten.KeyP, true
	case "q":
		return ebiten.KeyQ, true
	case "r":
		return ebiten.KeyR, true
	case "s":
		return ebiten.KeyS, true
	case "t":
		return ebiten.KeyT, true
	case "u":
		return ebiten.KeyU, true
	case "v":
		return ebiten.KeyV, true
	case "w":
		return ebiten.KeyW, true
	case "x":
		return ebiten.KeyX, true
	case "y":
		return ebiten.KeyY, true
	case "z":
		return ebiten.KeyZ, true
	case "alt":
		return ebiten.KeyAlt, true
	case "apostrophe":
		return ebiten.KeyApostrophe, true
	case "backslash":
		return ebiten.KeyBackslash, true
	case "backspace":
		return ebiten.KeyBackspace, true
	case "capslock":
		return ebiten.KeyCapsLock, true
	case "comma":
		return ebiten.KeyComma, true
	case "control":
		return ebiten.KeyControl, true
	case "delete":
		return ebiten.KeyDelete, true
	case "down":
		return ebiten.KeyDown, true
	case "end":
		return ebiten.KeyEnd, true
	case "enter":
		return ebiten.KeyEnter, true
	case "equal":
		return ebiten.KeyEqual, true
	case "escape":
		return ebiten.KeyEscape, true
	case "f1":
		return ebiten.KeyF1, true
	case "f2":
		return ebiten.KeyF2, true
	case "f3":
		return ebiten.KeyF3, true
	case "f4":
		return ebiten.KeyF4, true
	case "f5":
		return ebiten.KeyF5, true
	case "f6":
		return ebiten.KeyF6, true
	case "f7":
		return ebiten.KeyF7, true
	case "f8":
		return ebiten.KeyF8, true
	case "f9":
		return ebiten.KeyF9, true
	case "f10":
		return ebiten.KeyF10, true
	case "f11":
		return ebiten.KeyF11, true
	case "f12":
		return ebiten.KeyF12, true
	case "graveaccent":
		return ebiten.KeyGraveAccent, true
	case "home":
		return ebiten.KeyHome, true
	case "insert":
		return ebiten.KeyInsert, true
	case "kp0":
		return ebiten.KeyKP0, true
	case "kp1":
		return ebiten.KeyKP1, true
	case "kp2":
		return ebiten.KeyKP2, true
	case "kp3":
		return ebiten.KeyKP3, true
	case "kp4":
		return ebiten.KeyKP4, true
	case "kp5":
		return ebiten.KeyKP5, true
	case "kp6":
		return ebiten.KeyKP6, true
	case "kp7":
		return ebiten.KeyKP7, true
	case "kp8":
		return ebiten.KeyKP8, true
	case "kp9":
		return ebiten.KeyKP9, true
	case "kpadd":
		return ebiten.KeyKPAdd, true
	case "kpdecimal":
		return ebiten.KeyKPDecimal, true
	case "kpdivide":
		return ebiten.KeyKPDivide, true
	case "kpenter":
		return ebiten.KeyKPEnter, true
	case "kpequal":
		return ebiten.KeyKPEqual, true
	case "kpmultiply":
		return ebiten.KeyKPMultiply, true
	case "kpsubtract":
		return ebiten.KeyKPSubtract, true
	case "left":
		return ebiten.KeyLeft, true
	case "leftbracket":
		return ebiten.KeyLeftBracket, true
	case "menu":
		return ebiten.KeyMenu, true
	case "minus":
		return ebiten.KeyMinus, true
	case "numlock":
		return ebiten.KeyNumLock, true
	case "pagedown":
		return ebiten.KeyPageDown, true
	case "pageup":
		return ebiten.KeyPageUp, true
	case "pause":
		return ebiten.KeyPause, true
	case "period":
		return ebiten.KeyPeriod, true
	case "printscreen":
		return ebiten.KeyPrintScreen, true
	case "right":
		return ebiten.KeyRight, true
	case "rightbracket":
		return ebiten.KeyRightBracket, true
	case "scrolllock":
		return ebiten.KeyScrollLock, true
	case "semicolon":
		return ebiten.KeySemicolon, true
	case "shift":
		return ebiten.KeyShift, true
	case "slash":
		return ebiten.KeySlash, true
	case "space":
		return ebiten.KeySpace, true
	case "tab":
		return ebiten.KeyTab, true
	case "up":
		return ebiten.KeyUp, true
	}
	return 0, false
}
