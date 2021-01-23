package os

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/hugorut/protop/portgw/internal"
)

// DecodedEntry holds information about a given decoded entry.
type DecodedEntry struct {
	Port  internal.Port
	Error error
}

// DecodeFunc defines a function that takes a input os file and decodes
// the file from a know format. It returns a channel of which it sends
// decoded entries to.
type DecodeFunc func(input io.Reader) (chan DecodedEntry, error)

// JSONDecode decodes the provides a json implementation of DecodeFunc.
// It decodes the input file and sends each JSON object onto the returned
// channel. If there is a problem with any entry it will send an error
// on the channel but continue to process the file.
func JSONDecode(input io.Reader) (chan DecodedEntry, error) {
	dec := json.NewDecoder(input)

	_, err := dec.Token()
	if err != nil {
		return nil, fmt.Errorf("could not read start token of json file, %w", err)
	}

	ch := make(chan DecodedEntry)

	go readPorts(dec, ch)

	return ch, nil
}

// readPorts uses the json.Decoder to read individual
// entries from json file.
func readPorts(dec *json.Decoder, ch chan DecodedEntry) {
	for dec.More() {
		// read the start token of the json structure as
		// the file is specified as a map.
		_, err := dec.Token()
		if err != nil {
			break
		}

		var p internal.Port
		err = dec.Decode(&p)

		ch <- DecodedEntry{
			Port:  p,
			Error: err,
		}
	}

	close(ch)
}
