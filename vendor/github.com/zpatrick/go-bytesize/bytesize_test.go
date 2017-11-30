package bytesize

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

const MAX_INT64 = 9223372036854775807

func init() {
	rand.Seed(time.Now().Unix())
}

func TestBytesizeConversions(t *testing.T) {
	for i := 0; i < 1000; i++ {
		b := Bytesize(rand.Int63n(MAX_INT64))
		t.Logf("Testing %f", b)

		if r, e := b.Bytes(), float64(b); r != e {
			t.Errorf("Bytes %f: result %f, expected %f", b, r, e)
		}

		if r, e := b.Kilobytes(), float64(b/KB); r != e {
			t.Errorf("Kilobytes %f: result %f, expected %f", b, r, e)
		}

		if r, e := b.Megabytes(), float64(b/MB); r != e {
			t.Errorf("Megabytes %f: result %f, expected %f", b, r, e)
		}

		if r, e := b.Gigabytes(), float64(b/GB); r != e {
			t.Errorf("Gigabytes %f: result %f, expected %f", b, r, e)
		}

		if r, e := b.Terabytes(), float64(b/TB); r != e {
			t.Errorf("Terabytes %f: result %f, expected %f", b, r, e)
		}

		if r, e := b.Petabytes(), float64(b/PB); r != e {
			t.Errorf("Petabytes %f: result %f, expected %f", b, r, e)
		}

		if r, e := b.Exabytes(), float64(b/EB); r != e {
			t.Errorf("Exabytes %f: result %f, expected %f", b, r, e)
		}

		if r, e := b.Kibibytes(), float64(b/KiB); r != e {
			t.Errorf("Kibibytes %f: result %f, expected %f", b, r, e)
		}

		if r, e := b.Mebibytes(), float64(b/MiB); r != e {
			t.Errorf("Mebibytes %f: result %f, expected %f", b, r, e)
		}

		if r, e := b.Gibibytes(), float64(b/GiB); r != e {
			t.Errorf("Gibibytes %f: result %f, expected %f", b, r, e)
		}

		if r, e := b.Tebibytes(), float64(b/TiB); r != e {
			t.Errorf("Tebibytes %f: result %f, expected %f", b, r, e)
		}

		if r, e := b.Pebibytes(), float64(b/PiB); r != e {
			t.Errorf("Pebibytes %f: result %f, expected %f", b, r, e)
		}

		if r, e := b.Exbibytes(), float64(b/EiB); r != e {
			t.Errorf("Exbibytes %f: result %f, expected %f", b, r, e)
		}
	}
}

func TestBytesizeConstants(t *testing.T) {
	constants := map[string]struct {
		Constant Bytesize
		Expected Bytesize
	}{
		"B":   {B, 1},
		"KB":  {KB, 1000},
		"MB":  {MB, 1000000},
		"GB":  {GB, 1000000000},
		"TB":  {TB, 1000000000000},
		"PB":  {PB, 1000000000000000},
		"EB":  {EB, 1000000000000000000},
		"KiB": {KiB, 1024},
		"MiB": {MiB, 1048576},
		"GiB": {GiB, 1073741824},
		"TiB": {TiB, 1099511627776},
		"PiB": {PiB, 1125899906842624},
		"EiB": {EiB, 1152921504606846976},
	}

	for abbreviation, s := range constants {
		if s.Constant != s.Expected {
			t.Errorf("%s: result %f, expected %f", abbreviation, s.Constant, s.Expected)
		}
	}
}

func Example() {
	b := Bytesize(10000)
	fmt.Printf("%g bytes is: %g KB and %g MB\n", b, b.Kilobytes(), b.Megabytes())

	// Output:
	// 10000 bytes is: 10 KB and 0.01 MB
}

func ExampleBytesize_Format() {
	b := Bytesize(100000)
	fmt.Println(b.Format("b"))
	fmt.Println(b.Format("kb"))
	fmt.Println(b.Format("gb"))
	fmt.Println(b.Format("tb"))
	fmt.Println(b.Format("pb"))
	fmt.Println(b.Format("eb"))
	fmt.Println(b.Format("kib"))
	fmt.Println(b.Format("gib"))
	fmt.Println(b.Format("tib"))
	fmt.Println(b.Format("pib"))
	fmt.Println(b.Format("eib"))

	// Output:
	// 100000B
	// 100KB
	// 0.0001GB
	// 1e-07TB
	// 1e-10PB
	// 1e-13EB
	// 97.65625KiB
	// 9.313225746154785e-05GiB
	// 9.094947017729282e-08TiB
	// 8.881784197001252e-11PiB
	// 8.673617379884035e-14EiB
}
