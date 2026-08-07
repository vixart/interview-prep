// for-range по строке идет по РУНАМ, а индекс считается в БАЙТАХ:
// на многобайтовом π индекс перепрыгивает сразу на 2.
package main

import "fmt"

func main() {
	samples := []string{"hello", "apple_π!"}
	for _, sample := range samples {
		fmt.Println(sample)
		for i, r := range sample {
			// i — смещение в БАЙТАХ, r — руна; на π индекс прыгнет на 2
			fmt.Println(i, r, string(r))
		}
		fmt.Println()
	}

}
