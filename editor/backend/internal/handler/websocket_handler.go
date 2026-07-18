package handler

import (
	"encoding/json"

	"github.com/gofiber/contrib/websocket"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/event"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/logger"
	"github.com/krishnaZawar/LevelCraft/editor/backend/internal/queue"
	"github.com/krishnaZawar/LevelCraft/utils/helper"
	"github.com/krishnaZawar/LevelCraft/utils/models"
	utilqueue "github.com/krishnaZawar/LevelCraft/utils/queue"
)

var ls = logger.Get()

func HandleCommandRequests(c *websocket.Conn) {
	ls.Info().Msg("client connected")

	defer func() {
		ls.Info().Msg("client disconnected")
		c.Close()
	}()

	errCh := make(chan error, 2)

	go readLoop(c, errCh)
	go writeLoop(c, errCh)

	err := <-errCh
	if err != nil {
		ls.ErrorWith(err).Msg("connection closed")
	}
}

// reads and ingests commandRequests to the queue for further processing
func readLoop(c *websocket.Conn, errCh chan<- error) {
	cmdQueue := queue.GetCommandQueue()
	for {
		err := readLoopIter(c, cmdQueue)
		if err != nil {
			errCh <- err
			return
		}
	}
}

// handles the read logic that runs every iteration of the read loop
//
// returns error on read failures only.
// Unmarshalling is ignored as a failure because
// this is due to a faulty request, not a connection issue
func readLoopIter(c *websocket.Conn, cmdQueue *utilqueue.CommandQueue) error {
	messageType, msg, err := c.ReadMessage()
	if err != nil {
		return err
	}
	if messageType != websocket.TextMessage {
		ls.Warn().Msg("messageType received is invalid")
		return nil
	}

	var req *models.CommandRequest
	err = json.Unmarshal(msg, &req)
	if err != nil {
		ls.ErrorWith(err).Msg("JSON unmarshalling failed")
		return nil
	}

	ls.Info().Msgf("CommandRequest of type %s received", req.RequestType)
	cmdQueue.Ingest(*req)
	return nil
}

// Writes out eventResponses for frontend to act upon
func writeLoop(c *websocket.Conn, errCh chan<- error) {
	respQueue := queue.GetRespQueue()
	for {
		err := writeLoopIter(c, respQueue)
		if err != nil {
			errCh <- err
			return
		}
	}
}

// handles the write logic that runs every iteration of the write loop
//
// returns error on write failures only.
// Marshalling is ignored as a failure because
// this is due to a faulty response, not a connection issue
func writeLoopIter(c *websocket.Conn, respQueue *helper.Queue[event.EventResponse]) error {
	resp, ok := respQueue.Pop()
	if !ok {
		return nil
	}
	data, err := json.Marshal(resp)
	if err != nil {
		ls.ErrorWith(err).Msg("failed to marshal response")
		return nil
	}
	err = c.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return err
	}
	return nil
}
