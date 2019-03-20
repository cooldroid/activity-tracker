package service

import (
	"log"
	"time"

	"github.com/prashantgupta24/activity-tracker/internal/pkg/mouse"
	"github.com/prashantgupta24/activity-tracker/pkg/activity"
)

type MouseCursorHandler struct {
	tickerCh chan struct{}
}

type cursorInfo struct {
	didCursorMove   bool
	currentMousePos *mouse.Position
}

func (m *MouseCursorHandler) Start(activityCh chan *activity.Type) {

	m.tickerCh = make(chan struct{})

	go func() {
		lastMousePos := mouse.GetPosition()
		for range m.tickerCh {
			log.Printf("mouse cursor checked at : %v\n", time.Now())
			commCh := make(chan *cursorInfo)
			go checkCursorChange(commCh, lastMousePos)
			select {
			case cursorInfo := <-commCh:
				if cursorInfo.didCursorMove {
					activityCh <- &activity.Type{
						ActivityType: activity.MOUSE_CURSOR_MOVEMENT,
					}
					lastMousePos = cursorInfo.currentMousePos
				}
			case <-time.After(timeout * time.Millisecond):
				//timeout, do nothing
				log.Printf("timeout happened after %vms while checking mouse cursor handler", timeout)
			}
		}
		log.Printf("stopping cursor handler")
		return
	}()
}

func (m *MouseCursorHandler) Trigger() {
	//doing it the non-blocking sender way
	select {
	case m.tickerCh <- struct{}{}:
	default:
		//service is blocked, handle it somehow?
	}
}
func (m *MouseCursorHandler) Close() {
	close(m.tickerCh)
}

func checkCursorChange(commCh chan *cursorInfo, lastMousePos *mouse.Position) {
	currentMousePos := mouse.GetPosition()
	//log.Printf("current mouse position: %v\n", currentMousePos)
	//log.Printf("last mouse position: %v\n", lastMousePos)
	if currentMousePos.MouseX == lastMousePos.MouseX &&
		currentMousePos.MouseY == lastMousePos.MouseY {
		commCh <- &cursorInfo{
			didCursorMove:   false,
			currentMousePos: nil,
		}
	} else {
		commCh <- &cursorInfo{
			didCursorMove:   true,
			currentMousePos: currentMousePos,
		}
	}
}
