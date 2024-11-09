package healthcheck

import (
	"context"
	"fmt"
	"time"
)

func (h *healthcheck) asyncChecker() {
	var (
		failure int
		success int
	)
	ticker := time.NewTicker(h.opt.Readiness.PeriodSec)
	initialDelay := h.opt.Readiness.InitDelaySec - h.opt.Readiness.PeriodSec
	if initialDelay < 0 {
		initialDelay = 0
	}

	// init delay sec
	time.Sleep(initialDelay)

	h.efLogger.Info(`[OK]`, `Health: `, fmt.Sprintf("Readiness starts after %s delay", initialDelay))

	for {
		select {
		case <-ticker.C:
			ctx, cancel := context.WithTimeout(context.Background(), h.opt.Readiness.CheckTimeout)
			err := h.opt.Readiness.CheckF(ctx, cancel)
			if err != nil {
				failure++
				if failure > h.opt.Readiness.FailureThreshold {
					h.setReadinessStatus(false)
					failure = 0
					success = 0
				}
			} else {
				success++
				if success > h.opt.Readiness.SuccessThreshold {
					h.setReadinessStatus(true)
					failure = 0
					success = 0
				}
			}
		case <-h.termReady:
			ticker.Stop()
			return
		}
	}
}

func (h *healthcheck) setReadinessStatus(OK bool) {
	s := h.status
	s.mu.Lock()
	defer s.mu.Unlock()
	s.isReady = OK
}
