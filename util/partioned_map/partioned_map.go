package partioned_map

import "errors"

type partitioner interface {
	Find(key string) (uint, error)
}

type PartitionedMap[T any] struct {
	partsNum   uint            // partitions number
	partitions []*partition[T] // partitions slice
	finder     partitioner     // abstract partitions finder
}

func NewPartitionedMap[T any](partitioner partitioner, partsNum uint) *PartitionedMap[T] {
	partitions := make([]*partition[T], 0, partsNum)
	for i := 0; i < int(partsNum); i++ {
		m := make(map[string]T)
		partitions = append(partitions, &partition[T]{storage: m})
	}
	return &PartitionedMap[T]{partsNum: partsNum, partitions: partitions, finder: partitioner}
}

func (c *PartitionedMap[T]) Set(key string, value T) error {
	// find partition index
	partitionIndex, err := c.finder.Find(key)
	if err != nil {
		return err
	}

	// get partition from slice
	p := c.partitions[partitionIndex]
	// write key to the partition
	p.set(key, value)

	return nil
}
func (c *PartitionedMap[T]) Get(key string) (any, error) {
	partitionIndex, err := c.finder.Find(key)
	if err != nil {
		return nil, err
	}

	p := c.partitions[partitionIndex]
	value, ok := p.get(key)
	if !ok {
		return nil, errors.New("no such data")
	}

	return value, nil
}
