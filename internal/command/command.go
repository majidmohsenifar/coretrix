package command

import "context"

type ConsoleCommand interface {
	Run(ctx context.Context, flags []string)
}
