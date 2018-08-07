package pictures

import (
	"log"
	"sync"
)

var BufLen = 1024

type NonFatalError struct {
	reason string
	cont   bool
}

func (nfe *NonFatalError) Error() string { return nfe.reason }

func IsNonFatalError(err error) bool {
	if _, ok := err.(*NonFatalError); ok {
		return true
	}
	return false
}

type payload struct {
	info  *ImageInfo
	errCh chan error
}

type ProcessingChain struct {
	channels []chan *payload
	wg       sync.WaitGroup
	debug    bool
}

func NewProcessingChain() *ProcessingChain {
	return &ProcessingChain{
		channels: []chan *payload{make(chan *payload)},
	}
}

func (pc *ProcessingChain) Process(ii *ImageInfo) error {
	pl := &payload{ii, make(chan error)}
	if len(pc.channels) > 0 {
		pc.channels[0] <- pl
	}
	return <-pl.errCh
}

func (pc *ProcessingChain) AppendLast(filter ImageFilter) {
	pc.append(filter, false)
}

func (pc *ProcessingChain) Append(filter ImageFilter) *ProcessingChain {
	pc.append(filter, true)
	return pc
}

func (pc *ProcessingChain) append(filter ImageFilter, cont bool) {
	input := pc.channels[len(pc.channels)-1]
	output := make(chan *payload, BufLen)
	if cont {
		pc.channels = append(pc.channels, output)
	}
	pc.wg.Add(1)
	go func() {
		var err error
		for pl := range input {
			if pc.debug {
				log.Printf("%T processing", filter)
			}
			pl.info, err = filter.Process(pl.info)

			// Success/Failure scenarios
			// 1) No error, keep going if necessary
			//
			// 2) Something non-fatal happened, processing should continue (if necessary)
			//    e.g. no exif data, no need to notify the user
			//
			// 3) Something non-fatal happened, processing should stop
			//    e.g. duplicate image found, notify the user
			//
			// 4) Something un-expected happened, processing should stop
			//    log the error and notify the user of an internal failure
			if err == nil {
				if cont {
					output <- pl
				} else {
					pl.errCh <- nil
				}
			} else {
				if nfe, ok := err.(*NonFatalError); ok {
					if nfe.cont {
						if cont {
							output <- pl
						}
					} else {
						pl.errCh <- err
					}
				} else {
					pl.errCh <- err
				}
			}
		}
		close(output)
		pc.wg.Done()
	}()
}

func (pc *ProcessingChain) Close() error {
	close(pc.channels[0])
	pc.wg.Wait()
	return nil
}
