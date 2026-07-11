package shutdown

import (
	"io"
	"sync"

	"github.com/art-es/b2b-usage-based-billing-platform/services/service-auth/internal/pkg/log"
)

// graceful shutdown manager
type Shutdowner struct {
	logger      log.Logger
	openedConns []io.Closer
}

func New(logger log.Logger) *Shutdowner {
	return &Shutdowner{logger: logger}
}

func (m *Shutdowner) Add(conn io.Closer) {
	m.openedConns = append(m.openedConns, conn)
}

func (m *Shutdowner) AddFunc(fn func() error) {
	m.openedConns = append(m.openedConns, closerFunc(fn))
}

func (m *Shutdowner) Shutdown() {
	var wg sync.WaitGroup
	wg.Add(len(m.openedConns))

	for _, openedConn := range m.openedConns {
		go func() {
			defer wg.Done()
			if err := openedConn.Close(); err != nil {
				m.logger.Log(log.Error).
					Set("message", "graceful shutdown error").
					Set("error", err.Error()).
					Write()
			}
		}()
	}

	wg.Wait()

	m.logger.Log(log.Info).
		Set("message", "graceful shutdown completed").
		Write()
}

type closerFunc func() error

func (f closerFunc) Close() error { return f() }
