package shutdown

import (
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mercor/payment-service/pkg/log"
)

const maxTerminationTime = 25 * time.Second

type Callback interface {
	Close() error
}

type namedCallback struct {
	name     string
	callback Callback
}

type Server struct {
	shutdownCallbacks []*namedCallback
	drainCallbacks    []*namedCallback
	doneClosure       chan bool
}

var globalServer *Server

func init() {
	globalServer = &Server{
		doneClosure:       make(chan bool),
		shutdownCallbacks: make([]*namedCallback, 0),
		drainCallbacks:    make([]*namedCallback, 0),
	}
	globalServer.waitForTermination()
}

func RegisterShutdownCallback(name string, callback Callback) {
	globalServer.registerShutdownCallback(name, callback)
}

func RegisterDrainCallback(name string, callback Callback) {
	globalServer.registerDrainCallback(name, callback)
}

func GetWaitChannel() <-chan bool {
	return globalServer.doneClosure
}

func (s *Server) registerShutdownCallback(name string, callback Callback) {
	s.shutdownCallbacks = append(s.shutdownCallbacks, &namedCallback{name, callback})
}

func (s *Server) registerDrainCallback(name string, callback Callback) {
	s.drainCallbacks = append(s.drainCallbacks, &namedCallback{name, callback})
}

func (s *Server) waitForTermination() {
	var signals = make(chan os.Signal, 2)
	signal.Notify(signals, syscall.SIGTERM)
	signal.Notify(signals, syscall.SIGINT)
	go func() {
		s.gracefulShutdown(<-signals)
	}()
}

func (s *Server) gracefulShutdown(sig os.Signal) {
	log.Infof("Received signal %s, shutting down the application", sig)
	err := s.runWithTimeout(func() {
		s.drain()
		s.shutdown()
	}, maxTerminationTime)

	// If timeout occurred, log and exit with non-zero exit code
	if err != nil {
		log.Warn("Failed to shutdown gracefully, exiting forcefully. Err:", err)
		os.Exit(1)
	}

	log.Info("Completed graceful shutdown, exiting")
	s.doneClosure <- true
}

func (s *Server) drain() {
	for _, callback := range s.drainCallbacks {
		log.Info("Starting drain: ", callback.name)
		err := callback.callback.Close()
		if err != nil {
			log.Errorf("Error while draining %d: %v", callback.name, err)
		}
		log.Info("Completed drain for: ", callback.name)
	}
}

func (s *Server) shutdown() {
	for _, callback := range s.shutdownCallbacks {
		log.Info("Starting shutdown for ", callback.name)
		err := callback.callback.Close()
		if err != nil {
			log.Errorf("Error while shutdown %d: %v", callback.name, err)
		}
		log.Info("Completed shutdown for ", callback.name)
	}
}

func (s *Server) runWithTimeout(f func(), timeout time.Duration) (err error) {
	doneChan := make(chan bool, 1)
	go func() {
		f()
		doneChan <- true
	}()
	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case <-doneChan:
		err = nil
	case <-timer.C:
		err = errors.New("max time exceeded in running the functions")
	}
	return err
}
