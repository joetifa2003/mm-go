package mm

import "testing"

const LOOP_TIMES = 500

func BenchmarkSlice(b *testing.B) {
	for i := 0; i < b.N; i++ {
		numbers := []int{}

		for j := 0; j < LOOP_TIMES; j++ {
			numbers = append(numbers, j)
		}

		for j := 0; j < LOOP_TIMES; j++ {
			// Pop
			numbers = numbers[:len(numbers)-1]
		}
	}
}

func BenchmarkVector(b *testing.B) {
	for i := 0; i < b.N; i++ {
		numbersVec := NewVector[int]()

		for j := 0; j < LOOP_TIMES; j++ {
			numbersVec.Push(j)
		}

		for j := 0; j < LOOP_TIMES; j++ {
			numbersVec.Pop()
		}

		numbersVec.Free()
	}
}

func BenchmarkLinkedList(b *testing.B) {
	for i := 0; i < b.N; i++ {
		numbersVec := NewLinkedList[int]()

		for j := 0; j < LOOP_TIMES; j++ {
			numbersVec.PushBack(j)
		}

		for j := 0; j < LOOP_TIMES; j++ {
			numbersVec.PopBack()
		}

		numbersVec.Free()
	}
}
