package main

import (
	"fmt"
)

func addTwoNumbers(l1 []int, l2 []int) *[]int {
	lr1 := make([]int, len(l1))
	lr2 := make([]int, len(l2))
	for n1, n1value := range l1 {
		lr1[len(lr1)-1-n1] = n1value
	}
	for n2, n2value := range l2 {
		lr2[len(lr2)-1-n2] = n2value
	}

	var num2 int
	var num1 int
	for _, m1 := range lr1 {
		num1 = num1*10 + m1
	}
	for _, m2 := range lr2 {
		num2 = num2*10 + m2
	}
	var result int = num1 + num2

	count := 0
	for temp := result; temp != 0; temp /= 10 {
		count++
	}

	resultSlice := make([]int, count)

	for i := 0; i < count; i++ {
		resultSlice[i] = result % 10
		result /= 10
	}

	return &resultSlice
}

func main() {

	arr1 := []int{7, 7, 4, 22}

	arr2 := []int{1, 2, 9}

	fmt.Println(*addTwoNumbers(arr1, arr2))
}
