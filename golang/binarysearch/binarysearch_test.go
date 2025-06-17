package binarysearch_test

import (
	"algorithm/binarysearch"
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestReadFileToSlice(t *testing.T) {
	filename := "test_number.txt"
	content := "30, 10, 50, 20, 40\n5, 45, 25, 35, 15\n"
	err := os.WriteFile(filename, []byte(content), 0644)

	if err != nil {
		t.Fatalf("无法创建文件: %v", err)
	}

	defer os.Remove(filename)

	got, err := binarysearch.ReadFileToSlice(filename)
	if err != nil {
		t.Fatalf("ReadFileToSlice执行异常: %v", err)
	}

	want := []int{30, 10, 50, 20, 40, 5, 45, 25, 35, 15}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("ReadFileToSlice(%q) = %v; want %v", filename, got, want)
	} else {
		t.Logf("\n测试通过, ReadFileToSlice(%q) = %v", filename, got)
	}
}

func TestArrSorted(t *testing.T) {
	arr, err := binarysearch.ReadFileToSlice("numbers.txt")

	if err != nil {
		t.Fatalf("\nReadFileToSlice函数执行失败: %v", err)
	}

	sort.Ints(arr)
	t.Logf("\n测试通过,排序后的arr数组: %v", arr)
}

func TestRank(t *testing.T) {
	arr, err := binarysearch.ReadFileToSlice("numbers.txt")

	if err != nil {
		t.Fatalf("ReadFileToSlice函数执行失败: %v", err)
	}

	sort.Ints(arr)
	t.Logf("排序后的arr数组: %v", arr)

	tests := []struct {
		target int
		want   int
	}{
		{10, 1},
		{25, 4},
		{50, 9},
		{5, 0},
		{35, 6},
		{100, -1}, // target not in array
		{-10, -1}, // target not in array
	}

	for _, tt := range tests {
		got := binarysearch.Rank(arr, tt.target)
		if got != tt.want {
			t.Errorf("Rank(%v, %d) = %d; want %d", arr, tt.target, got, tt.want)
		} else {
			t.Logf("测试通过, Rank(%v, %d) = %d", arr, tt.target, got)
		}
	}
}
