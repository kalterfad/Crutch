package crutchService

import (
	types "crutch/internal/services"
	"log"
)

type mainFacade struct {
	cancelChanMap  map[types.ChannelID]chan struct{}
	dirEvent       chan types.DirPathData
	newWorkerCh    chan types.DirData
	cancelWorkerCh chan types.ChannelID
}

func newFacade(newWorkerCh chan types.DirData, cancelWorkerCh chan types.ChannelID) *mainFacade {
	facade := &mainFacade{
		cancelChanMap:  map[types.ChannelID]chan struct{}{},
		newWorkerCh:    newWorkerCh,
		cancelWorkerCh: cancelWorkerCh,
		dirEvent:       make(chan types.DirPathData, 10),
	}
	go facade.manage()

	return facade
}

func (f *mainFacade) restart(dirInfoList []types.DirData) {
	log.Println("restart observer")
	for _, dirData := range dirInfoList {
		f.startObserver(dirData)
	}
}

func (f *mainFacade) stopObserver(channelID types.ChannelID) {
	close(f.cancelChanMap[channelID])
	delete(f.cancelChanMap, channelID)
	log.Printf("stop '%v' worker", channelID)
}

func (f *mainFacade) startObserver(dirData types.DirData) {
	cancelChan := make(chan struct{})
	f.cancelChanMap[dirData.ChannelID] = cancelChan

	go newObserver(f.dirEvent, dirData, cancelChan).run()
}

func (f *mainFacade) startHandler(pathData types.DirPathData) {
	go newFileHandler(pathData).handle()
}

func (f *mainFacade) manage() {
	for {
		select {
		case dirData := <-f.newWorkerCh:
			f.startObserver(dirData)
		case channelID := <-f.cancelWorkerCh:
			f.stopObserver(channelID)
		case pathData := <-f.dirEvent:
			f.startHandler(pathData)
		}
	}
}
