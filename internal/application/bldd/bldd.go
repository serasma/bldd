package bldd

import (
	"context"
	"fmt"

	"github.com/serrasma/bldd/internal/domain"
)

type BLDD struct {
	walker   *domain.Walker
	reporter *domain.Reporter
	errs     chan error
}

func New(walker *domain.Walker, reporter *domain.Reporter, errs chan error) *BLDD {
	return &BLDD{
		walker:   walker,
		reporter: reporter,
		errs:     errs,
	}
}

func (b *BLDD) Run(ctx context.Context) {
	go func() {
		librariesUsage, err := b.walker.Walk(ctx)
		if err != nil {
			b.errs <- fmt.Errorf("bldd: walker walk error: %w", err)
			return
		}

		if err = b.reporter.Render(librariesUsage); err != nil {
			b.errs <- fmt.Errorf("bldd: reporter render error: %w", err)
			return
		}

		b.errs <- nil
	}()
}
