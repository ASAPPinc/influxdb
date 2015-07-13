package client

import "time"

type BufferConfig struct {
	FlushMaxPoints   int
	FlushMaxWaitTime time.Duration
}

func NewBufferedClient(clientConfig Config, bufferConfig BufferConfig) (bufferedClient *BufferedClient, err error) {
	client, err := NewClient(clientConfig)
	if err != nil {
		return
	}
	ingestChan := make(chan Point, 10000)
	flushTimer := time.NewTimer(bufferConfig.FlushMaxWaitTime)
	bufferedClient = &BufferedClient{client, bufferConfig, ingestChan, flushTimer, BatchPoints{}}
	go bufferedClient.ingestAndFlushLoop()
	return
}

type BufferedClient struct {
	*Client
	bufferConfig BufferConfig
	ingestChan   chan Point
	flushTimer   *time.Timer
	batch        BatchPoints
}

func (b *BufferedClient) Add(measurement string, val interface{}, tags map[string]string) {
	b.ingestChan <- Point{
		Measurement: measurement,
		Tags:        tags,
		Fields: map[string]interface{}{
			"value": val,
		},
		// Time        time.Time
		// Precision   string
		// Raw         string
	}
}

// Async ingest and flush loops
///////////////////////////////

func (b *BufferedClient) ingestAndFlushLoop() {
	for { // loop indefinitely
		select {
		case point := <-b.ingestChan:
			b.batch.Points = append(b.batch.Points, point)
			if len(b.batch.Points) > b.bufferConfig.FlushMaxPoints {
				b.initiateFlush()
			}
		case <-b.flushTimer.C:
			b.initiateFlush()
		}
	}
}

func (b *BufferedClient) initiateFlush() {
	b.flushTimer.Stop()
	b.Client.Write(b.batch)
	b.batch.Points = nil
	b.flushTimer.Reset(b.bufferConfig.FlushMaxWaitTime)
}
