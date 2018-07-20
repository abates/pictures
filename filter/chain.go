package filter

import (
	"log"
	"sync"
)

var BufLen = 1024

type NonfatalError struct {
	reason string
}

func (nfe *NonfatalError) Error() string { return nfe.reason }

type FatalError struct {
	reason string
}

func (fe *FatalError) Error() string { return fe.reason }

func IsFatalError(err error) bool {
	if err == nil {
		return false
	} else if _, ok := err.(*NonfatalError); ok {
		return false
	}

	return true
}

type ProcessingChain struct {
	channels []chan *ImageInfo
	wg       sync.WaitGroup
}

func NewProcessingChain() *ProcessingChain {
	return &ProcessingChain{
		channels: []chan *ImageInfo{make(chan *ImageInfo)},
	}
}

func (pc *ProcessingChain) Input() chan<- *ImageInfo {
	return pc.channels[0]
}

func (pc *ProcessingChain) Output() <-chan *ImageInfo {
	return pc.channels[len(pc.channels)-1]
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
	output := make(chan *ImageInfo, BufLen)
	if cont {
		pc.channels = append(pc.channels, output)
	}
	pc.wg.Add(1)
	go func() {
		for ii := range input {
			ii, err := filter.Process(ii)
			if IsFatalError(err) {
				log.Fatalf("%T:Fatal Error:%v", filter, err)
			} else if err == nil && ii != nil && cont {
				output <- ii
			} else if err != nil {
				log.Printf("%v\n", err)
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
