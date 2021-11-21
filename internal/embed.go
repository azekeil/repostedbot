// While this is unusual, it's required to be able to embed content as the
// content must be at or below the current directory level, and it can't
// be part of the package otherwise that causes an import cycle.
//
package internal

import "embed"

//go:embed commands/*.go
var CommandsEmbedFS embed.FS
