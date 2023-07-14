// packagesketch;
package main

import (
	"fmt"
	"math/rand"
)

// This Sketch is a struct that represents a sketch.
// It has two fields:
//   - HashFunc: A function that takes an integer and returns a hash value.
//   - Summary: A slice of integers that stores the cardinality of each bucket.
type DataSketch struct {
	HashFunc func(x int) int
	Summary  []int
}

// NewSketch creates a new sketch.
// The hash function is a function that takes an integer and returns a hash value.
// The length of the summary slice is 2^16.
func NewSketch(hashFunc func(x int) int) *DataSketch {
	// Create a new sketch.
	ts := &DataSketch{
		HashFunc: hashFunc,
		Summary:  make([]int, 1<<16),
	}

	// Return the sketch.
	return ts
}

// Add adds an element to the sketch.
// The element is hashed using the hash function and the corresponding bucket in the summary slice is incremented.
func (ts *DataSketch) Add(x int) {
	// Hash the element.
	hashValue := ts.HashFunc(x)

	// Increment the corresponding bucket in the summary slice.
	ts.Summary[hashValue]++
}

// EstimateCardinality estimates the cardinality of the sketch.
// The cardinality is the sum of the values in the summary slice.
func (ts *DataSketch) EstimateCardinality() int {
	// Initialize a variable to store the cardinality.
	count := 0

	// Iterate over the summary slice and add the values to the cardinality variable.
	for _, x := range ts.Summary {
		count += x
	}

	// Return the cardinality.
	return count
}

// Serialize serializes the sketch to a byte slice.
// The byte slice has the following format:
//   - The first byte is the length of the summary slice.
//   - The next `len(summary)` bytes are the values in the summary slice.
func (ts *DataSketch) Serialize() []byte {
	// Initialize a byte slice to store the serialized sketch.
	buf := make([]byte, 1+len(ts.Summary))

	// Write the length of the summary slice to the byte slice.
	buf[0] = byte(len(ts.Summary))

	// Write the values in the summary slice to the byte slice.
	i := 0
	for _, x := range ts.Summary {
		buf[1+i] = byte(x)
		i = i + 1
	}

	// Return the serialized sketch.
	return buf
}

// Deserialize deserializes a sketch from a byte slice.
// The byte slice must have the following format:
//   - The first byte is the length of the summary slice.
//   - The next `len(summary)` bytes are the values in the summary slice.
func (ts *DataSketch) Deserialize(data []byte) error {
	// Get the length of the summary slice.
	n := len(data)

	// Check if the byte slice has the correct format.
	if n < 1 {
		return fmt.Errorf("invalid data")
	}

	// Initialize the summary slice.
	ts.Summary = make([]int, n-1)

	// Deserialize the values in the summary slice.
	for i := 1; i < n; i++ {
		ts.Summary[i-1] = int(data[i])
	}

	// Return nil if the deserialization was successful.
	return nil
}

// return maximum value
func min(val1 int, val2 int) int {
	if val1 < val2 {
		return val1
	}
	return val2
}

// return minimum value
func max(val1 int, val2 int) int {
	if val1 > val2 {
		return val1
	}
	return val2
}

// Union returns a new sketch that is the union of the two given sketches.
// The union is created by adding the counts from the two sketches together.
func (ts *DataSketch) Union(other *DataSketch) *DataSketch {
	// Create a new sketch.
	newTs := NewSketch(ts.HashFunc)

	// Iterate over the summaries of the two sketches and add the counts together.
	for i := range ts.Summary {
		newTs.Summary[i] += ts.Summary[i] + other.Summary[i]
	}

	// Return the new sketch.
	return newTs
}

// Intersection returns a new sketch that is the intersection of the two given sketches.
// The intersection is created by taking the minimum of the counts from the two sketches.
func (ts *DataSketch) Intersection(other *DataSketch) *DataSketch {
	// Create a new sketch.
	newTs := NewSketch(ts.HashFunc)

	// Iterate over the summaries of the two sketches and take the minimum of the counts.
	for i := range ts.Summary {
		newTs.Summary[i] = min(ts.Summary[i], other.Summary[i])
	}

	// Return the new sketch.
	return newTs
}

// ANotB returns a new sketch that is the difference between the two given sketches.
// The difference is created by taking the difference of the counts from the two sketches.
func (ts *DataSketch) ANotB(other *DataSketch) *DataSketch {
	// Create a new sketch.
	newTs := NewSketch(ts.HashFunc)

	// Iterate over the summaries of the two sketches and take the difference of the counts.
	for i := range ts.Summary {
		newTs.Summary[i] = max(0, ts.Summary[i]-other.Summary[i])
	}

	// Return the new sketch.
	return newTs
}

// main is the entry point for the program.
func main() {

	// Create two sketches
	// Add elements to the  sketches
	ts1 := NewSketch(func(x int) int { return x % 10 })
	for i := 0; i < 100; i++ {
		ts1.Add(rand.Intn(10))
	}

	ts2 := NewSketch(func(x int) int { return x % 10 })
	for i := 0; i < 50; i++ {
		ts2.Add(rand.Intn(10))
	}

	// Print the cardinality of the  sketches
	fmt.Println("ts1:", ts1.EstimateCardinality())
	fmt.Println("ts2:", ts2.EstimateCardinality())

	// Create a new  sketch that is the union of ts1 and ts2
	newTs := ts1.Union(ts2)
	fmt.Println("ts1 union ts2:", newTs.EstimateCardinality())

	// Create a new  sketch that is the intersection of ts1 and ts2.
	newTs = ts1.Intersection(ts2)
	fmt.Println("ts1 intersection ts2:", newTs.EstimateCardinality())

	// Create a new  sketch that is the ANotB of ts1 and ts2.
	newTs = ts1.ANotB(ts2)
	fmt.Println("ts1 ANotB ts2:", newTs.EstimateCardinality())

}
