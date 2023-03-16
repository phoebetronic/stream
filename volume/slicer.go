package volume

import "github.com/phoebetronic/stream"

type Slicer struct {
	his int
	lis []stream.Volume
}

func (s *Slicer) Add(v stream.Volume) {
	s.lis = append(s.lis, v)

	if len(s.lis) > s.his {
		s.lis = append(s.lis[:0], s.lis[len(s.lis)-s.his:]...)
	}
}

func (s *Slicer) Lis() []stream.Volume {
	return s.lis
}
