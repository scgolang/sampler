// Package sampler is a sampler that uses SuperCollider.
package sampler

import (
	"time"

	"github.com/mkb218/gosndfile/sndfile"
	"github.com/pkg/errors"
	"github.com/scgolang/sc"
)

// Sampler is a sample playback engine based on SuperCollider.
// It can load samples into 128 different slots.
// Each slot will usually correspond to a midi note that is used to trigger the sample.
// A slot can play multiple samples.
type Sampler struct {
	client  *sc.Client
	group   *sc.Group
	samples [128][]sample
}

// New creates a new sampler.
func New(scsynthAddr string) (*Sampler, error) {
	client, err := sc.NewClient("udp", "0.0.0.0:0", scsynthAddr, 5*time.Second)
	if err != nil {
		return nil, err
	}
	group, err := client.AddDefaultGroup()
	if err != nil {
		return nil, err
	}
	for _, def := range []*sc.Synthdef{
		sc.NewSynthdef("sampler_simple_mono", simpleDef(1)),
		sc.NewSynthdef("sampler_simple_stereo", simpleDef(2)),
	} {
		if err := client.SendDef(def); err != nil {
			return nil, err
		}
	}
	return &Sampler{
		client: client,
		group:  group,
	}, nil
}

// Add adds a sample at the provided path to the specified slot
func (s *Sampler) Add(audioFile string, slot int) error {
	var info sndfile.Info
	if _, err := sndfile.Open(audioFile, sndfile.Read, &info); err != nil {
		return err
	}
	if info.Channels != 1 && info.Channels != 2 {
		return errors.New("only samples with 1 or 2 channels are supported")
	}
	if err := validateSlot(slot); err != nil {
		return err
	}
	s.samples[slot] = append(s.samples[slot], sample{numChannels: int(info.Channels)})
	return nil
}

type sample struct {
	numChannels int
}

var (
	defSimpleMono   = sc.NewSynthdef("sampler_simple_mono", simpleDef(1))
	defSimpleStereo = sc.NewSynthdef("sampler_simple_stereo", simpleDef(2))
)

func simpleDef(numChannels int) sc.UgenFunc {
	return func(params sc.Params) sc.Ugen {
		sig := sc.PlayBuf{
			NumChannels: numChannels,
			BufNum:      params.Add("bufnum", 0),
			Done:        sc.FreeEnclosing,
		}.Rate(sc.AR)

		return sc.Out{
			Bus:      sc.C(0),
			Channels: sc.Multi(sig, sig),
		}.Rate(sc.AR)
	}
}

func validateSlot(slot int) error {
	if slot < 0 || slot > 127 {
		return errors.Errorf("slot (%d) must be >= 0 and <= 127")
	}
	return nil
}
