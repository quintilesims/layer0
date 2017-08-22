package series

type IntSeries struct {
	Start, End int
}

func Ints(start, end int) IntSeries {
	return IntSeries{start, end}
}

func (s IntSeries) Contains(v int) bool {
	return v >= s.Start && v <= s.End
}

func (s IntSeries) Slice() []int {
	slice := make([]int, s.End-s.Start+1)
	for i := 0; i < len(slice); i++ {
		slice[i] = i + s.Start
	}

	return slice
}

type Int64Series struct {
	Start, End int64
}

func Int64s(start, end int64) Int64Series {
	return Int64Series{start, end}
}

func (s Int64Series) Contains(v int64) bool {
	return v >= s.Start && v <= s.End
}

func (s Int64Series) Slice() []int64 {
	slice := make([]int64, int(s.End-s.Start)+1)
	for i := 0; i < len(slice); i++ {
		slice[i] = int64(i) + s.Start
	}

	return slice
}

type Float64Series struct {
	Start, End float64
}

func Float64s(start, end float64) Float64Series {
	return Float64Series{start, end}
}

func (s Float64Series) Contains(v float64) bool {
	return v >= s.Start && v <= s.End
}

func (s Float64Series) Slice() []float64 {
	slice := make([]float64, int(s.End-s.Start)+1)
	for i := 0; i < len(slice); i++ {
		slice[i] = float64(i) + s.Start
	}

	return slice
}
