package circularqueue

import (
	"testing"
)

type message struct {
	msgID int16
	msg   interface{}
}

type C2SLogin struct {
	UID   int64  `json:"uid"`
	Token string `json:"token"`
}

func TestPush(t *testing.T) {
	q1 := NewCircularQueue()
	if q1.Len() != 0 {
		t.Errorf("q1.Len() = %d, want %d", q1.Len(), 0)
	}
	if err := q1.Retrieve(); err != ErrEmptyQueue {
		t.Error("Retrieve empty queue should return ErrEmptyQueue")
	}

	q := NewCircularQueueWithSize(2)
	if !q.IsEmpty() {
		t.Errorf("q.IsEmpty() = %v, want %v", q.IsEmpty(), true)
	}

	m := &message{1, &C2SLogin{101, "sss"}}
	q.Push(m)
	if q.readableIndex != 0 {
		t.Error("q.readleIndex is wrong")
	}
	if q.writableIndex != 1 {
		t.Error("q.writableIndex is wrong")
	}
	if q.buffer[0] != m {
		t.Error("q.buffer[0] is wrong")
	}

	{
		m := &message{1, &C2SLogin{101, "sss"}}
		q.Push(m)
		if q.readableIndex != 0 {
			t.Error("q.readleIndex is wrong")
		}
		if q.writableIndex != 2 {
			t.Error("q.writableIndex is wrong")
		}
		if q.buffer[1] != m {
			t.Error("q.buffer[0] is wrong")
		}
		if !q.isFull() {
			t.Errorf("q.isFull() = %v, want %v", q.isFull(), true)
		}
	}

	{
		m := &message{1, &C2SLogin{101, "sss"}}
		q.Push(m)
		if q.readableIndex != 0 {
			t.Error("q.readleIndex is wrong")
		}
		if q.writableIndex != 3 {
			t.Error("q.writableIndex is wrong")
		}
		if q.buffer[2] != m {
			t.Error("q.buffer[0] is wrong")
		}
		if q.isFull() {
			t.Errorf("q.isFull() = %v, want %v", q.isFull(), false)
		}
		if len(q.buffer) != 5 {
			t.Errorf("len(q.buffer) = %v, want %v", len(q.buffer), 5)
		}
		if cap(q.buffer) != 5 {
			t.Errorf("cap(q.buffer) = %v, want %v", cap(q.buffer), 5)
		}
		if q.Len() != 3 {
			t.Errorf("q.Len() = %v, want %v", q.Len(), 3)
		}
	}

	{
		q := NewCircularQueueWithSize(5)
		m := &message{1, &C2SLogin{101, "sss"}}
		q.Push(m)
		q.Push(m)
		q.Push(m)
		q.Push(m)
		_, err := q.Pop()
		if err != nil {
			t.Fatalf("Pop error %v", err)
		}
		q.Push(m)
		if q.writableIndex != 0 || q.readableIndex != 1 {
			t.Error("index is wrong")
		}
		if q.Len() != 4 {
			t.Errorf("q.Len() = %d, want %d", q.Len(), 4)
		}

		q.Push(m)
	}
}

func TestPop(t *testing.T) {
	q := NewCircularQueueWithSize(2)

	for i := 0; i < 2; i++ {
		m := &message{1, &C2SLogin{101, "sss"}}
		q.Push(m)
	}

	m, err := q.Pop()
	if err != nil {
		t.Error("q.Pop() return err")
	}
	v, ok := m.(*message)
	if !ok {
		t.Error("m.(Type) is not *message")
	}
	if v.msgID != 1 {
		t.Errorf("m.msgID = %d, want %d", v.msgID, 1)
	}

	_, err = q.Pop()
	if err != nil {
		t.Error("q.Pop() return err")
	}
	if !q.IsEmpty() {
		t.Errorf("q.IsEmpty() = %v, want %v", q.IsEmpty(), true)
	}
	_, err = q.Pop()
	if err == nil {
		t.Error("q.Pop() should return err")
	}

	for i := 0; i < 2; i++ {
		m := &message{1, &C2SLogin{101, "sss"}}
		q.Push(m)
	}

	_, err = q.Pop()
	if err != nil {
		t.Fatalf("Pop error %v", err)
	}
	if q.readableIndex != 1 {
		t.Error("q.readableIndex != 1")
	}
	if q.writableIndex != 0 {
		t.Error("q.writableIndex != 0")
	}
	if !q.isFull() {
		t.Errorf("q.isFull() = %v, want %v", q.isFull(), true)
	}

	_, err = q.Pop()
	if err != nil {
		t.Fatalf("Pop error %v", err)
	}
	if !q.IsEmpty() {
		t.Errorf("q.IsEmpty() = %v, want %v", q.IsEmpty(), true)
	}
	if q.readableIndex != 0 {
		t.Error("q.readableIndex != 0")
	}
	if q.writableIndex != 0 {
		t.Error("q.writableIndex != 0")
	}
}

func BenchmarkPush(b *testing.B) {
	b.ReportAllocs()
	q := NewCircularQueue()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Push(i)
	}
}

func BenchmarkPushWithSizeQueue(b *testing.B) {
	b.ReportAllocs()
	q := NewCircularQueueWithSize(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Push(i)
	}
}

func BenchmarkPushAndPop(b *testing.B) {
	b.ReportAllocs()
	q := NewCircularQueue()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Push(i)
	}
	for i := 0; i < b.N; i++ {
		_, err := q.Pop()
		if err != nil {
			b.Fatalf("Pop error %v", err)
		}
	}
}

func BenchmarkPushAndPopWithSizeQueue(b *testing.B) {
	b.ReportAllocs()
	q := NewCircularQueueWithSize(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		q.Push(i)
	}
	for i := 0; i < b.N; i++ {
		_, err := q.Pop()
		if err != nil {
			b.Fatalf("Pop error %v", err)
		}
	}
}
