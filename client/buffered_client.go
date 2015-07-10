package client

import "time"

type BufferConfig struct {
	BufferMaxSize    int
	FlushMaxWaitTime time.Duration
}

func NewBufferedClient(clientConfig Config, bufferConfig BufferConfig) (bufferedClient *BufferedClient, err error) {
	client, err := NewClient(c)
	if err != nil {
		return
	}
	bufferedClient = &BufferedClient{client, bufferConfig}
	go bufferedClient.ingestLoop()
	go bufferedClient.flushLoop()
	return bufferedClient
}

type BufferedClient struct {
	*Client
	bufferConfig BufferConfig
	batch        *BatchPoints
	ingestChan   chan Point
	flushChan    chan bool
}

func (b *BufferedClient) BatchWrite(point) {

}

// Async ingest and flush loops
///////////////////////////////

func (b *BufferedClient) ingestLoop() {
	for { // loop indefinitely
		point := <-b.ingestChan
		b.batch.Points = append(b.batch.Points, point)
		if len(b.batch.Points) > b.bufferConfig.BufferMaxSize {
			b.initiateFlush()
		}
	}
}

func (b *BufferedClient) initiateFlush() {
	b.flushChan <- b.batch
	b.batch = b.newBatch()
}

func (b *BufferedClient) flushLoop() {
	for { // loop indefinitely
		select {
		case <-time.After(b.bufferConfig.FlushMaxWaitTime):

		case <-b.flushChan:

		}

	}
}
