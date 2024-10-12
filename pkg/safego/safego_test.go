package safego

import (
	"fmt"
	"io"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestCall(t *testing.T) {
	PanicCatchFunc = func(name string, p interface{}) {
		t.Log(name, p)
		t.Log(string(debug.Stack()))
		StackOnce(5)
		StackOnce(4)
		StackOnce(3)
		StackOnce(2)
		StackOnce(1)
		StackOnce(0)
	}
	CallError("1", func() error {
		panic("xx")
	})
	CallError("2", func() error {
		n := 0
		_ = 1 / n
		return nil
	})
	CallError("3", func() error {
		w := io.Writer(nil)
		w.Write([]byte{55, 65})
		return nil
	})

	Go("4", func() {
		n := 0
		_ = 1 / n
	})
	runtime.Gosched()
	time.Sleep(time.Millisecond)

	fmt.Println(string(debug.Stack()))

	StackOnce(0)
	return
}

func BenchmarkStackOnce(b *testing.B) {
	b.Run("once", func(b *testing.B) {
		for k := 0; k < b.N; k++ {
			StackOnce(2)
		}
	})
	b.Run("debug", func(b *testing.B) {
		for k := 0; k < b.N; k++ {
			debug.Stack()
		}
	})
}

func StackOnce(skip int) string {

	var pcs = []uintptr{0}
	numFrames := runtime.Callers(skip+2, pcs)
	_ = numFrames
	//fmt.Println("numFrames", numFrames)

	frames := runtime.CallersFrames(pcs)
	frame, _ := frames.Next()
	buf := strings.Builder{}
	buf.Grow(1024)
	buf.WriteByte(' ')
	buf.WriteString(strconv.FormatInt(int64(skip), 10))
	buf.WriteByte(' ')
	buf.WriteString(frame.File)
	buf.WriteByte(':')
	buf.WriteString(strconv.FormatInt(int64(frame.Line), 10))
	buf.WriteByte(' ')
	buf.WriteString(frame.Func.Name())
	return buf.String() // fmt.Sprintf("\t%d %s:%d %s\n", skip, frame.File, frame.Line, frame.Func.Name())
	//return ""
}

func PanicStack(p interface{}) {
	switch p.(string) {
	case "index out of range",
		"slice bounds out of range",
		"negative shift amount",
		"slice length too short to convert to pointer to array":
		// panicCheck1
	case "integer divide by zero",
		"integer overflow",
		"floating point error",
		"invalid memory address or nil pointer dereference":
		// panicCheck2
	}
}
