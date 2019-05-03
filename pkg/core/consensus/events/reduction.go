package events

import (
	"bytes"

	"gitlab.dusk.network/dusk-core/dusk-go/pkg/p2p/wire"
	"gitlab.dusk.network/dusk-core/dusk-go/pkg/p2p/wire/encoding"
)

type (
	// Reduction is a basic reduction event.
	Reduction struct {
		*Header
		VotedHash  []byte
		SignedHash []byte
	}

	ReductionUnMarshaller struct {
		*UnMarshaller
	}

	ReductionUnmarshaller interface {
		wire.EventMarshaller
		wire.EventDeserializer
		MarshalVoteSet(*bytes.Buffer, []wire.Event) error
		UnmarshalVoteSet(*bytes.Buffer) ([]wire.Event, error)
	}
)

// NewReduction returns and empty Reduction event.
func NewReduction() *Reduction {
	return &Reduction{
		Header: &Header{},
	}
}

// Equal as specified in the Event interface
func (e *Reduction) Equal(ev wire.Event) bool {
	other, ok := ev.(*Reduction)
	return ok && (bytes.Equal(e.PubKeyBLS, other.PubKeyBLS)) &&
		(e.Round == other.Round) && (e.Step == other.Step)
}

func NewReductionUnMarshaller() *ReductionUnMarshaller {
	return &ReductionUnMarshaller{NewUnMarshaller()}
}

func (r *ReductionUnMarshaller) NewEvent() wire.Event {
	return NewReduction()
}

// Unmarshal unmarshals the buffer into a Committee
func (a *ReductionUnMarshaller) Unmarshal(r *bytes.Buffer, ev wire.Event) error {
	bev := ev.(*Reduction)
	if err := a.HeaderUnmarshaller.Unmarshal(r, bev.Header); err != nil {
		return err
	}

	if err := encoding.Read256(r, &bev.VotedHash); err != nil {
		return err
	}

	if err := encoding.ReadBLS(r, &bev.SignedHash); err != nil {
		return err
	}

	return nil
}

// Marshal a Reduction into a buffer.
func (a *ReductionUnMarshaller) Marshal(r *bytes.Buffer, ev wire.Event) error {
	bev := ev.(*Reduction)
	if err := a.HeaderMarshaller.Marshal(r, bev.Header); err != nil {
		return err
	}

	if err := encoding.Write256(r, bev.VotedHash); err != nil {
		return err
	}

	if err := encoding.WriteBLS(r, bev.SignedHash); err != nil {
		return err
	}

	return nil
}

func (a *ReductionUnMarshaller) UnmarshalVoteSet(r *bytes.Buffer) ([]wire.Event, error) {
	length, err := encoding.ReadVarInt(r)
	if err != nil {
		return nil, err
	}

	evs := make([]wire.Event, length)
	for i := uint64(0); i < length; i++ {
		rev := &Reduction{
			Header: &Header{},
		}
		if err := a.Unmarshal(r, rev); err != nil {
			return nil, err
		}

		evs[i] = rev
	}

	return evs, nil
}

func (a *ReductionUnMarshaller) MarshalVoteSet(r *bytes.Buffer, evs []wire.Event) error {
	if err := encoding.WriteVarInt(r, uint64(len(evs))); err != nil {
		return err
	}

	for _, event := range evs {
		if err := a.Marshal(r, event); err != nil {
			return err
		}
	}

	return nil
}
