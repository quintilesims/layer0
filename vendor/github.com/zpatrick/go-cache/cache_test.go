package cache

import (
	"math/rand"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestSet(t *testing.T) {
	c := New()
	c.Set("1", 1)

	if result, expected := c.Get("1"), 1; !reflect.DeepEqual(result, expected) {
		t.Errorf("Result was %#v, expected %#v", result, expected)
	}
}

func TestSetf(t *testing.T) {
	c := New()
	c.Setf("1", 1, time.Millisecond)

	if _, exists := c.Getf("1"); !exists {
		t.Errorf("Entry for key '1' should not have expired yet")
	}

	time.Sleep(time.Millisecond * 2)

	if _, exists := c.Getf("1"); exists {
		t.Errorf("Entry for key '1' should have expired by now")
	}
}

func TestClear(t *testing.T) {
	c := New()
	for i := 0; i < 10; i++ {
		c.Set(strconv.Itoa(i), i)
	}

	c.Clear()

	if keys := c.Keys(); len(keys) != 0 {
		t.Errorf("Cache should have been empty, had keys: %v", keys)
	}
}

func TestDelete(t *testing.T) {
	c := New()
	c.Set("1", 1)
	c.Delete("1")

	if _, exists := c.Getf("1"); exists {
		t.Errorf("Entry for key '1' should not exist")
	}
}

func TestClearEvery(t *testing.T) {
	c := New()
	for i := 0; i < 10; i++ {
		c.Set(strconv.Itoa(i), i)
	}

	c.ClearEvery(time.Millisecond)

	if keys := c.Keys(); len(keys) != 10 {
		t.Errorf("Cache should have had 10 keys, but had keys: %v", keys)
	}

	time.Sleep(time.Millisecond * 2)

	if keys := c.Keys(); len(keys) != 0 {
		t.Errorf("Cache should have been empty, had keys: %v", keys)
	}
}

func TestGet(t *testing.T) {
	c := New()
	c.Set("1", 1)

	if result, expected := c.Get("1"), 1; !reflect.DeepEqual(result, expected) {
		t.Errorf("Result for entry '1' was %#v, expected %#v", result, expected)
	}

	if result := c.Get("2"); result != nil {
		t.Errorf("Result for entry '2' was %#v, expected nil", result)
	}
}

func TestGetf(t *testing.T) {
	c := New()
	c.Set("1", 1)

	result, exists := c.Getf("1")
	if !exists {
		t.Error("Entry for key '1' should exist")
	}

	if expected := 1; !reflect.DeepEqual(result, expected) {
		t.Errorf("Entry for key '1' was %#v, expected %#v", result, expected)
	}

	if _, exists := c.Getf("2"); exists {
		t.Errorf("Entry for key '2' should not exist")
	}
}

func TestItems(t *testing.T) {
	c := New()
	for i := 0; i < 5; i++ {
		c.Set(strconv.Itoa(i), i)
	}

	expected := map[string]T{
		"0": 0,
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 4,
	}

	if result := c.Items(); !reflect.DeepEqual(result, expected) {
		t.Errorf("Result was %#v, expected %#v", result, expected)
	}
}

func TestKeys(t *testing.T) {
	c := New()
	for i := 0; i < 5; i++ {
		c.Set(strconv.Itoa(i), i)
	}

	expected := []string{"0", "1", "2", "3", "4"}
	if result := c.Keys(); !reflect.DeepEqual(result, expected) {
		t.Errorf("Result was %#v, expected %#v", result, expected)
	}
}

func TestStressConcurrentAccess(t *testing.T) {
	c := New()
	c.ClearEvery(time.Nanosecond * 10)

	done := make(chan bool)
	for i := 0; i < 1000; i++ {
		go func() {
			key := strconv.Itoa(rand.Int())

			switch rand.Intn(8) {
			case 0:
				c.Set(key, rand.Int())
			case 1:
				c.Setf(key, rand.Int(), time.Nanosecond*5)
			case 2:
				c.Clear()
			case 3:
				c.Delete(key)
			case 4:
				c.Get(key)
			case 5:
				c.Getf(key)
			case 6:
				c.Items()
			case 7:
				c.Keys()
			}

			done <- true
		}()
	}

	for i := 0; i < 1000; i++ {
		<-done
	}
}

func benchmarkSet(count int, b *testing.B) {
	c := New()

	for n := 0; n < b.N; n++ {
		for i := 0; i < count; i++ {
			c.Set(strconv.Itoa(i), i)
		}
	}
}

func BenchmarkSet1(b *testing.B)     { benchmarkSet(1, b) }
func BenchmarkSet10(b *testing.B)    { benchmarkSet(10, b) }
func BenchmarkSet100(b *testing.B)   { benchmarkSet(100, b) }
func BenchmarkSet1000(b *testing.B)  { benchmarkSet(1000, b) }
func BenchmarkSet10000(b *testing.B) { benchmarkSet(10000, b) }

func benchmarkDelete(count int, b *testing.B) {
	c := New()

	for n := 0; n < b.N; n++ {
		for i := 0; i < count; i++ {
			c.Delete(strconv.Itoa(i))
		}
	}
}

func BenchmarkDelete1(b *testing.B)     { benchmarkDelete(1, b) }
func BenchmarkDelete10(b *testing.B)    { benchmarkDelete(10, b) }
func BenchmarkDelete100(b *testing.B)   { benchmarkDelete(100, b) }
func BenchmarkDelete1000(b *testing.B)  { benchmarkDelete(1000, b) }
func BenchmarkDelete10000(b *testing.B) { benchmarkDelete(10000, b) }

func benchmarkGet(count int, b *testing.B) {
	c := New()

	for n := 0; n < b.N; n++ {
		for i := 0; i < count; i++ {
			c.Get(strconv.Itoa(i))
		}
	}
}

func BenchmarkGet1(b *testing.B)     { benchmarkGet(1, b) }
func BenchmarkGet10(b *testing.B)    { benchmarkGet(10, b) }
func BenchmarkGet100(b *testing.B)   { benchmarkGet(100, b) }
func BenchmarkGet1000(b *testing.B)  { benchmarkGet(1000, b) }
func BenchmarkGet10000(b *testing.B) { benchmarkGet(10000, b) }
