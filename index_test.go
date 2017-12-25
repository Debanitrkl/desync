package desync

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestIndexLoad(t *testing.T) {
	f, err := os.Open("testdata/index.caibx")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	index, err := IndexFromReader(f)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("min", index.Index.ChunkSizeMin)
	fmt.Println("avg", index.Index.ChunkSizeAvg)
	fmt.Println("max", index.Index.ChunkSizeMax)
	fmt.Printf("flags %x", index.Index.FeatureFlags)

	type chunk struct {
		chunk string
		start uint64
		size  uint64
	}
	expected := []chunk{
		{"437884da2d1e61cf50b43b263ff15f25a870b0eae84bc22e4b5c307a0428764d", 0, 242168},
		{"985462e6b3293bbe61e43882686b481751ecf4b285bae4dffc2dfa8829f971ac", 242168, 75740},
		{"fadff4b303624f2be3d0e04c2f105306118a9f608ef1e4f83c1babbd23a2315f", 317908, 20012},
	}
	for i := range expected {
		id, _ := ChunkIDFromString(expected[i].chunk)
		exp := IndexChunk{ID: id, Start: expected[i].start, Size: expected[i].size}
		got := index.Chunks[i]
		if !reflect.DeepEqual(exp, got) {
			t.Fatalf("expected %v, got %v", exp, got)
		}
	}
}

func TestIndexWrite(t *testing.T) {
	in, err := ioutil.ReadFile("testdata/index.caibx")
	if err != nil {
		t.Fatal(err)
	}

	idx, err := IndexFromReader(bytes.NewReader(in))
	if err != nil {
		t.Fatal(err)
	}

	out := new(bytes.Buffer)
	n, err := idx.WriteTo(out)
	if err != nil {
		t.Fatal(err)
	}

	// in/out should match
	if !bytes.Equal(in, out.Bytes()) {
		t.Fatalf("decoded/encoded don't match")
	}
	if n != int64(out.Len()) {
		t.Fatalf("unexpected length")
	}
}