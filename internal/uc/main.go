package uc

import (
	"github.com/cfichtmueller/stor/internal/bus"
	"github.com/cfichtmueller/stor/internal/domain/archive"
)

func Configure() {
	bus.SubscribeCE(archive.EventCompleted, onArchiveCompleted)
}
