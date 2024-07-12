package main

import (
	"context"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/worming004/top10tick/common"
)

// Tick server represent a single instance that track all transaction for single tick, and publish it to stream
type TickServer struct {
	TickName    string
	writer      *kafka.Writer
	currentTick common.TickValue
	duration    time.Duration
}

func NewTickServer(tickName string, writer *kafka.Writer, duration time.Duration) *TickServer {
	return &TickServer{
		TickName:    tickName,
		writer:      writer,
		currentTick: common.TickValue{TickName: tickName, Value: 100},
		duration:    duration,
	}
}

func (ts TickServer) Start() {
	for {
		time.Sleep(ts.duration)

		go func() {
			newTick, err := ts.currentTick.GetNextTransaction()
			if err != nil {
				slog.Error("Failed to generate next tick", "name", ts.TickName)
				return
			}
			ts.currentTick = newTick

			value, err := ts.currentTick.SerializeJson()

			if err != nil {
				slog.Error("Failed to serialize tick value", "name", ts.TickName)
				return
			}

			err = ts.writer.WriteMessages(
				context.Background(),
				kafka.Message{
					Key:   []byte(ts.TickName),
					Value: value,
				})

			if err != nil {
				slog.Error("Failed to write message", "error", err.Error())
			} else {
				slog.Debug("Successfully message write")
			}
		}()
	}
}
