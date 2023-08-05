//go:build !windows

// Package console sets console's behavior on init
package console

import (
	"fmt"

	"github.com/FloatTech/ZeroBot-Plugin/kanban/banner"
)

func init() {
	fmt.Print("\033]0;RemiliaBot " + banner.Version + " " + banner.Copyright + "\007")
}
