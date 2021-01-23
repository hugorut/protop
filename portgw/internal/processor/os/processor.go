package os

import (
	"fmt"
	"os"
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/hugorut/protop/portgw/internal/store"
)

// OpenFunc defines a function that opens a file at the given name.
type OpenFunc func(name string) (*os.File, error)

// Processor handles processing a file in the machine os.
// It satisfies the FileProcessor interface.
type Processor struct {
	Processors int

	Open   OpenFunc
	Decode DecodeFunc
	Store  store.PortStore

	Logger logrus.FieldLogger
}

// Process takes the location of a given json
func (p Processor) Process(loc string) (int, error) {
	f, err := p.Open(loc)
	if err != nil {
		return 0, fmt.Errorf("could not open provided file at loc %s, %w", loc, err)
	}

	stream, err := p.Decode(f)
	if err != nil {
		p.closeFile(f)
		return 0, fmt.Errorf("failed to initalize decoding file, %w", err)
	}

	// set the default number of processors to just a single goroutine
	processors := 1
	if p.Processors > 1 {
		processors = p.Processors
	}

	// We batch the processing of individual ports into a set number of workers.
	// This allows us to constrain the number of downstream go routine so as to
	// limit machine resource usage.
	wg := &sync.WaitGroup{}
	for i := 0; i < processors; i++ {
		wg.Add(1)
		go p.processStream(stream, f.Name(), wg)
	}

	go p.waitForCompletion(wg, f)

	// right now we return an arbitrary identifier and just let the file processing
	// happen in the background. In future we would allow API clients to peek on the
	// status of the file processing, e.g. % complete. This id could be used to return
	// a status. This is not implemented currently.
	return 1, nil
}

// processStream handles a stream of decode messages
func (p Processor) processStream(stream chan DecodedEntry, name string, wg *sync.WaitGroup) {
	for msg := range stream {
		if msg.Error != nil {
			p.Logger.Errorf("os processor, could not process record in file %s, err: %s", name, msg.Error)
			continue
		}

		err := p.Store.Store(msg.Port)
		if err != nil {
			p.Logger.Errorf("os processor, could not store port record %s, err: %s", msg.Port.Name, err)
			continue
		}
	}

	wg.Done()
}

func (p Processor) waitForCompletion(wg *sync.WaitGroup, f *os.File) {
	wg.Wait()

	p.closeFile(f)
}

func (p Processor) closeFile(f *os.File) {
	err := f.Close()
	if err != nil {
		p.Logger.Errorf("os processor, could not close processing file , %s", err)
	}
}
