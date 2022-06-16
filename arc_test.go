// Copyright 2015 Andrew Bursavich. All rights reserved.
// Use of this source code is governed by The MIT License
// which can be found in the LICENSE file.

package arc

import (
	"fmt"
	"reflect"
	"testing"

	"bursavich.dev/arc/internal/list"
)

func TestTable(t *testing.T) {
	type state struct {
		rl, rd []int
		fl, fd []int
	}
	tests := []struct {
		cmd   string
		key   int
		val   string
		state state
	}{
		{cmd: "set", key: 5, val: "five", state: state{rl: []int{5}, rd: []int{}, fl: []int{}, fd: []int{}}},
		{cmd: "set", key: 17, val: "seventeen", state: state{rl: []int{17, 5}, rd: []int{}, fl: []int{}, fd: []int{}}},
		{cmd: "set", key: 8, val: "eight", state: state{rl: []int{8, 17, 5}, rd: []int{}, fl: []int{}, fd: []int{}}},
		{cmd: "set", key: 0, val: "zero", state: state{rl: []int{0, 8, 17, 5}, rd: []int{}, fl: []int{}, fd: []int{}}},
		{cmd: "set", key: 13, val: "thirteen", state: state{rl: []int{13, 0, 8, 17, 5}, rd: []int{}, fl: []int{}, fd: []int{}}},
		{cmd: "set", key: 25, val: "twenty-five", state: state{rl: []int{25, 13, 0, 8, 17, 5}, rd: []int{}, fl: []int{}, fd: []int{}}},
		{cmd: "set", key: 27, val: "twenty-seven", state: state{rl: []int{27, 25, 13, 0, 8, 17, 5}, rd: []int{}, fl: []int{}, fd: []int{}}},
		{cmd: "set", key: 26, val: "twenty-six", state: state{rl: []int{26, 27, 25, 13, 0, 8, 17, 5}, rd: []int{}, fl: []int{}, fd: []int{}}},
		{cmd: "get", key: 8, val: "eight", state: state{rl: []int{26, 27, 25, 13, 0, 17, 5}, rd: []int{}, fl: []int{8}, fd: []int{}}},
		{cmd: "get", key: 13, val: "thirteen", state: state{rl: []int{26, 27, 25, 0, 17, 5}, rd: []int{}, fl: []int{13, 8}, fd: []int{}}},
		{cmd: "set", key: 19, val: "nineteen", state: state{rl: []int{19, 26, 27, 25, 0, 17, 5}, rd: []int{}, fl: []int{13, 8}, fd: []int{}}},
		{cmd: "get", key: 17, val: "seventeen", state: state{rl: []int{19, 26, 27, 25, 0, 5}, rd: []int{}, fl: []int{17, 13, 8}, fd: []int{}}},
		{cmd: "set", key: 15, val: "fifteen", state: state{rl: []int{15, 19, 26, 27, 25, 0, 5}, rd: []int{}, fl: []int{17, 13, 8}, fd: []int{}}},
		{cmd: "set", key: 28, val: "twenty-eight", state: state{rl: []int{28, 15, 19, 26, 27, 25, 0}, rd: []int{5}, fl: []int{17, 13, 8}, fd: []int{}}},
		{cmd: "set", key: 2, val: "two", state: state{rl: []int{2, 28, 15, 19, 26, 27, 25}, rd: []int{0, 5}, fl: []int{17, 13, 8}, fd: []int{}}},
		{cmd: "set", key: 3, val: "three", state: state{rl: []int{3, 2, 28, 15, 19, 26, 27}, rd: []int{25, 0, 5}, fl: []int{17, 13, 8}, fd: []int{}}},
		{cmd: "get", key: 2, val: "two", state: state{rl: []int{3, 28, 15, 19, 26, 27}, rd: []int{25, 0, 5}, fl: []int{2, 17, 13, 8}, fd: []int{}}},
		{cmd: "set", key: 11, val: "eleven", state: state{rl: []int{11, 3, 28, 15, 19, 26}, rd: []int{27, 25, 0, 5}, fl: []int{2, 17, 13, 8}, fd: []int{}}},
		{cmd: "set", key: 5, val: "five", state: state{rl: []int{11, 3, 28, 15, 19}, rd: []int{26, 27, 25, 0}, fl: []int{5, 2, 17, 13, 8}, fd: []int{}}},
		{cmd: "set", key: 14, val: "fourteen", state: state{rl: []int{14, 11, 3, 28, 15}, rd: []int{19, 26, 27, 25, 0}, fl: []int{5, 2, 17, 13, 8}, fd: []int{}}},
		{cmd: "set", key: 22, val: "twenty-two", state: state{rl: []int{22, 14, 11, 3, 28}, rd: []int{15, 19, 26, 27, 25}, fl: []int{5, 2, 17, 13, 8}, fd: []int{}}},
		{cmd: "set", key: 24, val: "twenty-four", state: state{rl: []int{24, 22, 14, 11, 3}, rd: []int{28, 15, 19, 26, 27}, fl: []int{5, 2, 17, 13, 8}, fd: []int{}}},
		{cmd: "set", key: 4, val: "four", state: state{rl: []int{4, 24, 22, 14, 11}, rd: []int{3, 28, 15, 19, 26}, fl: []int{5, 2, 17, 13, 8}, fd: []int{}}},
		{cmd: "set", key: 15, val: "fifteen", state: state{rl: []int{4, 24, 22, 14}, rd: []int{11, 3, 28, 19, 26}, fl: []int{15, 5, 2, 17, 13, 8}, fd: []int{}}},
		{cmd: "get", key: 4, val: "four", state: state{rl: []int{24, 22, 14}, rd: []int{11, 3, 28, 19, 26}, fl: []int{4, 15, 5, 2, 17, 13, 8}, fd: []int{}}},
		{cmd: "set", key: 11, val: "eleven", state: state{rl: []int{24, 22, 14}, rd: []int{3, 28, 19, 26}, fl: []int{11, 4, 15, 5, 2, 17, 13}, fd: []int{8}}},
		{cmd: "get", key: 2, val: "two", state: state{rl: []int{24, 22, 14}, rd: []int{3, 28, 19, 26}, fl: []int{2, 11, 4, 15, 5, 17, 13}, fd: []int{8}}},
		{cmd: "set", key: 19, val: "nineteen", state: state{rl: []int{24, 22, 14}, rd: []int{3, 28, 26}, fl: []int{19, 2, 11, 4, 15, 5, 17}, fd: []int{13, 8}}},
		{cmd: "set", key: 7, val: "seven", state: state{rl: []int{7, 24, 22, 14}, rd: []int{3, 28, 26}, fl: []int{19, 2, 11, 4, 15, 5}, fd: []int{17, 13, 8}}},
		{cmd: "set", key: 29, val: "twenty-nine", state: state{rl: []int{29, 7, 24, 22, 14}, rd: []int{3, 28, 26}, fl: []int{19, 2, 11, 4, 15}, fd: []int{5, 17, 13, 8}}},
		{cmd: "set", key: 12, val: "twelve", state: state{rl: []int{12, 29, 7, 24, 22}, rd: []int{14, 3, 28, 26}, fl: []int{19, 2, 11, 4, 15}, fd: []int{5, 17, 13, 8}}},
		{cmd: "get", key: 2, val: "two", state: state{rl: []int{12, 29, 7, 24, 22}, rd: []int{14, 3, 28, 26}, fl: []int{2, 19, 11, 4, 15}, fd: []int{5, 17, 13, 8}}},
		{cmd: "set", key: 25, val: "twenty-five", state: state{rl: []int{25, 12, 29, 7, 24}, rd: []int{22, 14, 3, 28, 26}, fl: []int{2, 19, 11, 4, 15}, fd: []int{5, 17, 13, 8}}},
		{cmd: "set", key: 0, val: "zero", state: state{rl: []int{0, 25, 12, 29, 7}, rd: []int{24, 22, 14, 3, 28}, fl: []int{2, 19, 11, 4, 15}, fd: []int{5, 17, 13, 8}}},
		{cmd: "get", key: 29, val: "twenty-nine", state: state{rl: []int{0, 25, 12, 7}, rd: []int{24, 22, 14, 3, 28}, fl: []int{29, 2, 19, 11, 4, 15}, fd: []int{5, 17, 13, 8}}},
		{cmd: "set", key: 28, val: "twenty-eight", state: state{rl: []int{0, 25, 12, 7}, rd: []int{24, 22, 14, 3}, fl: []int{28, 29, 2, 19, 11, 4}, fd: []int{15, 5, 17, 13, 8}}},
		{cmd: "set", key: 24, val: "twenty-four", state: state{rl: []int{0, 25, 12, 7}, rd: []int{22, 14, 3}, fl: []int{24, 28, 29, 2, 19, 11}, fd: []int{4, 15, 5, 17, 13, 8}}},
		{cmd: "get", key: 7, val: "seven", state: state{rl: []int{0, 25, 12}, rd: []int{22, 14, 3}, fl: []int{7, 24, 28, 29, 2, 19, 11}, fd: []int{4, 15, 5, 17, 13, 8}}},
		{cmd: "get", key: 12, val: "twelve", state: state{rl: []int{0, 25}, rd: []int{22, 14, 3}, fl: []int{12, 7, 24, 28, 29, 2, 19, 11}, fd: []int{4, 15, 5, 17, 13, 8}}},
		{cmd: "set", key: 17, val: "seventeen", state: state{rl: []int{0, 25}, rd: []int{22, 14, 3}, fl: []int{17, 12, 7, 24, 28, 29, 2, 19}, fd: []int{11, 4, 15, 5, 13, 8}}},
		{cmd: "set", key: 5, val: "five", state: state{rl: []int{0, 25}, rd: []int{22, 14, 3}, fl: []int{5, 17, 12, 7, 24, 28, 29, 2}, fd: []int{19, 11, 4, 15, 13, 8}}},
		{cmd: "get", key: 7, val: "seven", state: state{rl: []int{0, 25}, rd: []int{22, 14, 3}, fl: []int{7, 5, 17, 12, 24, 28, 29, 2}, fd: []int{19, 11, 4, 15, 13, 8}}},
		{cmd: "set", key: 14, val: "fourteen", state: state{rl: []int{0, 25}, rd: []int{22, 3}, fl: []int{14, 7, 5, 17, 12, 24, 28, 29}, fd: []int{2, 19, 11, 4, 15, 13, 8}}},
		{cmd: "set", key: 20, val: "twenty", state: state{rl: []int{20, 0, 25}, rd: []int{22, 3}, fl: []int{14, 7, 5, 17, 12, 24, 28}, fd: []int{29, 2, 19, 11, 4, 15, 13, 8}}},
		{cmd: "set", key: 26, val: "twenty-six", state: state{rl: []int{26, 20, 0, 25}, rd: []int{22, 3}, fl: []int{14, 7, 5, 17, 12, 24}, fd: []int{28, 29, 2, 19, 11, 4, 15, 13}}},
		{cmd: "set", key: 22, val: "twenty-two", state: state{rl: []int{26, 20, 0, 25}, rd: []int{3}, fl: []int{22, 14, 7, 5, 17, 12}, fd: []int{24, 28, 29, 2, 19, 11, 4, 15, 13}}},
		{cmd: "set", key: 24, val: "twenty-four", state: state{rl: []int{26, 20, 0, 25}, rd: []int{3}, fl: []int{24, 22, 14, 7, 5, 17}, fd: []int{12, 28, 29, 2, 19, 11, 4, 15, 13}}},
		{cmd: "set", key: 12, val: "twelve", state: state{rl: []int{26, 20, 0, 25}, rd: []int{3}, fl: []int{12, 24, 22, 14, 7, 5}, fd: []int{17, 28, 29, 2, 19, 11, 4, 15, 13}}},
		{cmd: "set", key: 9, val: "nine", state: state{rl: []int{9, 26, 20, 0, 25}, rd: []int{3}, fl: []int{12, 24, 22, 14, 7}, fd: []int{5, 17, 28, 29, 2, 19, 11, 4, 15}}},
		{cmd: "get", key: 9, val: "nine", state: state{rl: []int{26, 20, 0, 25}, rd: []int{3}, fl: []int{9, 12, 24, 22, 14, 7}, fd: []int{5, 17, 28, 29, 2, 19, 11, 4, 15}}},
		{cmd: "set", key: 28, val: "twenty-eight", state: state{rl: []int{26, 20, 0, 25}, rd: []int{3}, fl: []int{28, 9, 12, 24, 22, 14}, fd: []int{7, 5, 17, 29, 2, 19, 11, 4, 15}}},
		{cmd: "get", key: 22, val: "twenty-two", state: state{rl: []int{26, 20, 0, 25}, rd: []int{3}, fl: []int{22, 28, 9, 12, 24, 14}, fd: []int{7, 5, 17, 29, 2, 19, 11, 4, 15}}},
		{cmd: "set", key: 29, val: "twenty-nine", state: state{rl: []int{26, 20, 0, 25}, rd: []int{3}, fl: []int{29, 22, 28, 9, 12, 24}, fd: []int{14, 7, 5, 17, 2, 19, 11, 4, 15}}},
		{cmd: "set", key: 21, val: "twenty-one", state: state{rl: []int{21, 26, 20, 0, 25}, rd: []int{3}, fl: []int{29, 22, 28, 9, 12}, fd: []int{24, 14, 7, 5, 17, 2, 19, 11, 4}}},
		{cmd: "set", key: 2, val: "two", state: state{rl: []int{21, 26, 20, 0}, rd: []int{25, 3}, fl: []int{2, 29, 22, 28, 9, 12}, fd: []int{24, 14, 7, 5, 17, 19, 11, 4}}},
		{cmd: "set", key: 6, val: "six", state: state{rl: []int{6, 21, 26, 20, 0}, rd: []int{25, 3}, fl: []int{2, 29, 22, 28, 9}, fd: []int{12, 24, 14, 7, 5, 17, 19, 11}}},
		{cmd: "set", key: 18, val: "eighteen", state: state{rl: []int{18, 6, 21, 26, 20, 0}, rd: []int{25, 3}, fl: []int{2, 29, 22, 28}, fd: []int{9, 12, 24, 14, 7, 5, 17, 19}}},
		{cmd: "get", key: 21, val: "twenty-one", state: state{rl: []int{18, 6, 26, 20, 0}, rd: []int{25, 3}, fl: []int{21, 2, 29, 22, 28}, fd: []int{9, 12, 24, 14, 7, 5, 17, 19}}},
		{cmd: "get", key: 6, val: "six", state: state{rl: []int{18, 26, 20, 0}, rd: []int{25, 3}, fl: []int{6, 21, 2, 29, 22, 28}, fd: []int{9, 12, 24, 14, 7, 5, 17, 19}}},
		{cmd: "set", key: 17, val: "seventeen", state: state{rl: []int{18, 26, 20}, rd: []int{0, 25, 3}, fl: []int{17, 6, 21, 2, 29, 22, 28}, fd: []int{9, 12, 24, 14, 7, 5, 19}}},
		{cmd: "set", key: 4, val: "four", state: state{rl: []int{4, 18, 26, 20}, rd: []int{0, 25, 3}, fl: []int{17, 6, 21, 2, 29, 22}, fd: []int{28, 9, 12, 24, 14, 7, 5}}},
		{cmd: "get", key: 18, val: "eighteen", state: state{rl: []int{4, 26, 20}, rd: []int{0, 25, 3}, fl: []int{18, 17, 6, 21, 2, 29, 22}, fd: []int{28, 9, 12, 24, 14, 7, 5}}},
		{cmd: "get", key: 21, val: "twenty-one", state: state{rl: []int{4, 26, 20}, rd: []int{0, 25, 3}, fl: []int{21, 18, 17, 6, 2, 29, 22}, fd: []int{28, 9, 12, 24, 14, 7, 5}}},
		{cmd: "get", key: 26, val: "twenty-six", state: state{rl: []int{4, 20}, rd: []int{0, 25, 3}, fl: []int{26, 21, 18, 17, 6, 2, 29, 22}, fd: []int{28, 9, 12, 24, 14, 7, 5}}},
		{cmd: "set", key: 15, val: "fifteen", state: state{rl: []int{15, 4, 20}, rd: []int{0, 25, 3}, fl: []int{26, 21, 18, 17, 6, 2, 29}, fd: []int{22, 28, 9, 12, 24, 14, 7}}},
		{cmd: "set", key: 25, val: "twenty-five", state: state{rl: []int{15, 4, 20}, rd: []int{0, 3}, fl: []int{25, 26, 21, 18, 17, 6, 2}, fd: []int{29, 22, 28, 9, 12, 24, 14, 7}}},
		{cmd: "get", key: 26, val: "twenty-six", state: state{rl: []int{15, 4, 20}, rd: []int{0, 3}, fl: []int{26, 25, 21, 18, 17, 6, 2}, fd: []int{29, 22, 28, 9, 12, 24, 14, 7}}},
		{cmd: "set", key: 8, val: "eight", state: state{rl: []int{8, 15, 4, 20}, rd: []int{0, 3}, fl: []int{26, 25, 21, 18, 17, 6}, fd: []int{2, 29, 22, 28, 9, 12, 24, 14}}},
		{cmd: "set", key: 16, val: "sixteen", state: state{rl: []int{16, 8, 15, 4, 20}, rd: []int{0, 3}, fl: []int{26, 25, 21, 18, 17}, fd: []int{6, 2, 29, 22, 28, 9, 12, 24}}},
		{cmd: "get", key: 25, val: "twenty-five", state: state{rl: []int{16, 8, 15, 4, 20}, rd: []int{0, 3}, fl: []int{25, 26, 21, 18, 17}, fd: []int{6, 2, 29, 22, 28, 9, 12, 24}}},
		{cmd: "get", key: 16, val: "sixteen", state: state{rl: []int{8, 15, 4, 20}, rd: []int{0, 3}, fl: []int{16, 25, 26, 21, 18, 17}, fd: []int{6, 2, 29, 22, 28, 9, 12, 24}}},
		{cmd: "get", key: 17, val: "seventeen", state: state{rl: []int{8, 15, 4, 20}, rd: []int{0, 3}, fl: []int{17, 16, 25, 26, 21, 18}, fd: []int{6, 2, 29, 22, 28, 9, 12, 24}}},
		{cmd: "set", key: 0, val: "zero", state: state{rl: []int{8, 15, 4, 20}, rd: []int{3}, fl: []int{0, 17, 16, 25, 26, 21}, fd: []int{18, 6, 2, 29, 22, 28, 9, 12, 24}}},
		{cmd: "set", key: 27, val: "twenty-seven", state: state{rl: []int{27, 8, 15, 4, 20}, rd: []int{3}, fl: []int{0, 17, 16, 25, 26}, fd: []int{21, 18, 6, 2, 29, 22, 28, 9, 12}}},
		{cmd: "set", key: 10, val: "ten", state: state{rl: []int{10, 27, 8, 15, 4, 20}, rd: []int{3}, fl: []int{0, 17, 16, 25}, fd: []int{26, 21, 18, 6, 2, 29, 22, 28, 9}}},
		{cmd: "get", key: 4, val: "four", state: state{rl: []int{10, 27, 8, 15, 20}, rd: []int{3}, fl: []int{4, 0, 17, 16, 25}, fd: []int{26, 21, 18, 6, 2, 29, 22, 28, 9}}},
		{cmd: "set", key: 11, val: "eleven", state: state{rl: []int{11, 10, 27, 8, 15, 20}, rd: []int{3}, fl: []int{4, 0, 17, 16}, fd: []int{25, 26, 21, 18, 6, 2, 29, 22, 28}}},
		{cmd: "set", key: 14, val: "fourteen", state: state{rl: []int{14, 11, 10, 27, 8, 15, 20}, rd: []int{3}, fl: []int{4, 0, 17}, fd: []int{16, 25, 26, 21, 18, 6, 2, 29, 22}}},
		{cmd: "set", key: 3, val: "three", state: state{rl: []int{14, 11, 10, 27, 8, 15, 20}, rd: []int{}, fl: []int{3, 4, 0}, fd: []int{17, 16, 25, 26, 21, 18, 6, 2, 29, 22}}},
		{cmd: "set", key: 26, val: "twenty-six", state: state{rl: []int{14, 11, 10, 27, 8, 15, 20}, rd: []int{}, fl: []int{26, 3, 4}, fd: []int{0, 17, 16, 25, 21, 18, 6, 2, 29, 22}}},
		{cmd: "set", key: 9, val: "nine", state: state{rl: []int{9, 14, 11, 10, 27, 8, 15, 20}, rd: []int{}, fl: []int{26, 3}, fd: []int{4, 0, 17, 16, 25, 21, 18, 6, 2, 29}}},
		{cmd: "set", key: 23, val: "twenty-three", state: state{rl: []int{23, 9, 14, 11, 10, 27, 8, 15, 20}, rd: []int{}, fl: []int{26}, fd: []int{3, 4, 0, 17, 16, 25, 21, 18, 6, 2}}},
		{cmd: "set", key: 13, val: "thirteen", state: state{rl: []int{13, 23, 9, 14, 11, 10, 27, 8, 15, 20}, rd: []int{}, fl: []int{}, fd: []int{26, 3, 4, 0, 17, 16, 25, 21, 18, 6}}},
		{cmd: "set", key: 19, val: "nineteen", state: state{rl: []int{19, 13, 23, 9, 14, 11, 10, 27, 8, 15}, rd: []int{}, fl: []int{}, fd: []int{26, 3, 4, 0, 17, 16, 25, 21, 18, 6}}},
		{cmd: "set", key: 26, val: "twenty-six", state: state{rl: []int{19, 13, 23, 9, 14, 11, 10, 27, 8}, rd: []int{15}, fl: []int{26}, fd: []int{3, 4, 0, 17, 16, 25, 21, 18, 6}}},
		{cmd: "set", key: 4, val: "four", state: state{rl: []int{19, 13, 23, 9, 14, 11, 10, 27}, rd: []int{8, 15}, fl: []int{4, 26}, fd: []int{3, 0, 17, 16, 25, 21, 18, 6}}},
		{cmd: "set", key: 5, val: "five", state: state{rl: []int{5, 19, 13, 23, 9, 14, 11, 10}, rd: []int{27, 8}, fl: []int{4, 26}, fd: []int{3, 0, 17, 16, 25, 21, 18, 6}}},
		{cmd: "get", key: 26, val: "twenty-six", state: state{rl: []int{5, 19, 13, 23, 9, 14, 11, 10}, rd: []int{27, 8}, fl: []int{26, 4}, fd: []int{3, 0, 17, 16, 25, 21, 18, 6}}},
		{cmd: "get", key: 23, val: "twenty-three", state: state{rl: []int{5, 19, 13, 9, 14, 11, 10}, rd: []int{27, 8}, fl: []int{23, 26, 4}, fd: []int{3, 0, 17, 16, 25, 21, 18, 6}}},
		{cmd: "set", key: 28, val: "twenty-eight", state: state{rl: []int{28, 5, 19, 13, 9, 14, 11, 10}, rd: []int{27, 8}, fl: []int{23, 26}, fd: []int{4, 3, 0, 17, 16, 25, 21, 18}}},
		{cmd: "set", key: 22, val: "twenty-two", state: state{rl: []int{22, 28, 5, 19, 13, 9, 14, 11}, rd: []int{10, 27}, fl: []int{23, 26}, fd: []int{4, 3, 0, 17, 16, 25, 21, 18}}},
		{cmd: "set", key: 8, val: "eight", state: state{rl: []int{8, 22, 28, 5, 19, 13, 9, 14}, rd: []int{11, 10}, fl: []int{23, 26}, fd: []int{4, 3, 0, 17, 16, 25, 21, 18}}},
		{cmd: "set", key: 1, val: "one", state: state{rl: []int{1, 8, 22, 28, 5, 19, 13, 9}, rd: []int{14, 11}, fl: []int{23, 26}, fd: []int{4, 3, 0, 17, 16, 25, 21, 18}}},
		{cmd: "set", key: 27, val: "twenty-seven", state: state{rl: []int{27, 1, 8, 22, 28, 5, 19, 13}, rd: []int{9, 14}, fl: []int{23, 26}, fd: []int{4, 3, 0, 17, 16, 25, 21, 18}}},
		{cmd: "set", key: 3, val: "three", state: state{rl: []int{27, 1, 8, 22, 28, 5, 19}, rd: []int{13, 9, 14}, fl: []int{3, 23, 26}, fd: []int{4, 0, 17, 16, 25, 21, 18}}},
		{cmd: "set", key: 11, val: "eleven", state: state{rl: []int{11, 27, 1, 8, 22, 28, 5}, rd: []int{19, 13, 9}, fl: []int{3, 23, 26}, fd: []int{4, 0, 17, 16, 25, 21, 18}}},
		{cmd: "set", key: 10, val: "ten", state: state{rl: []int{10, 11, 27, 1, 8, 22, 28}, rd: []int{5, 19, 13}, fl: []int{3, 23, 26}, fd: []int{4, 0, 17, 16, 25, 21, 18}}},
		{cmd: "set", key: 25, val: "twenty-five", state: state{rl: []int{10, 11, 27, 1, 8, 22}, rd: []int{28, 5, 19, 13}, fl: []int{25, 3, 23, 26}, fd: []int{4, 0, 17, 16, 21, 18}}},
		{cmd: "get", key: 26, val: "twenty-six", state: state{rl: []int{10, 11, 27, 1, 8, 22}, rd: []int{28, 5, 19, 13}, fl: []int{26, 25, 3, 23}, fd: []int{4, 0, 17, 16, 21, 18}}},
		{cmd: "set", key: 12, val: "twelve", state: state{rl: []int{12, 10, 11, 27, 1, 8}, rd: []int{22, 28, 5, 19}, fl: []int{26, 25, 3, 23}, fd: []int{4, 0, 17, 16, 21, 18}}},
		{cmd: "set", key: 18, val: "eighteen", state: state{rl: []int{12, 10, 11, 27, 1}, rd: []int{8, 22, 28, 5, 19}, fl: []int{18, 26, 25, 3, 23}, fd: []int{4, 0, 17, 16, 21}}},
		{cmd: "set", key: 24, val: "twenty-four", state: state{rl: []int{24, 12, 10, 11, 27}, rd: []int{1, 8, 22, 28, 5}, fl: []int{18, 26, 25, 3, 23}, fd: []int{4, 0, 17, 16, 21}}},
		{cmd: "get", key: 18, val: "eighteen", state: state{rl: []int{24, 12, 10, 11, 27}, rd: []int{1, 8, 22, 28, 5}, fl: []int{18, 26, 25, 3, 23}, fd: []int{4, 0, 17, 16, 21}}},
		{cmd: "set", key: 5, val: "five", state: state{rl: []int{24, 12, 10, 11, 27}, rd: []int{1, 8, 22, 28}, fl: []int{5, 18, 26, 25, 3}, fd: []int{23, 4, 0, 17, 16, 21}}},
		{cmd: "set", key: 8, val: "eight", state: state{rl: []int{24, 12, 10, 11, 27}, rd: []int{1, 22, 28}, fl: []int{8, 5, 18, 26, 25}, fd: []int{3, 23, 4, 0, 17, 16, 21}}},
		{cmd: "get", key: 8, val: "eight", state: state{rl: []int{24, 12, 10, 11, 27}, rd: []int{1, 22, 28}, fl: []int{8, 5, 18, 26, 25}, fd: []int{3, 23, 4, 0, 17, 16, 21}}},
		{cmd: "get", key: 12, val: "twelve", state: state{rl: []int{24, 10, 11, 27}, rd: []int{1, 22, 28}, fl: []int{12, 8, 5, 18, 26, 25}, fd: []int{3, 23, 4, 0, 17, 16, 21}}},
		{cmd: "set", key: 17, val: "seventeen", state: state{rl: []int{24, 10, 11, 27}, rd: []int{1, 22, 28}, fl: []int{17, 12, 8, 5, 18, 26}, fd: []int{25, 3, 23, 4, 0, 16, 21}}},
		{cmd: "set", key: 0, val: "zero", state: state{rl: []int{24, 10, 11}, rd: []int{27, 1, 22, 28}, fl: []int{0, 17, 12, 8, 5, 18, 26}, fd: []int{25, 3, 23, 4, 16, 21}}},
		{cmd: "set", key: 7, val: "seven", state: state{rl: []int{7, 24, 10, 11}, rd: []int{27, 1, 22, 28}, fl: []int{0, 17, 12, 8, 5, 18}, fd: []int{26, 25, 3, 23, 4, 16}}},
		{cmd: "get", key: 8, val: "eight", state: state{rl: []int{7, 24, 10, 11}, rd: []int{27, 1, 22, 28}, fl: []int{8, 0, 17, 12, 5, 18}, fd: []int{26, 25, 3, 23, 4, 16}}},
		{cmd: "set", key: 1, val: "one", state: state{rl: []int{7, 24, 10, 11}, rd: []int{27, 22, 28}, fl: []int{1, 8, 0, 17, 12, 5}, fd: []int{18, 26, 25, 3, 23, 4, 16}}},
		{cmd: "set", key: 9, val: "nine", state: state{rl: []int{9, 7, 24, 10, 11}, rd: []int{27, 22, 28}, fl: []int{1, 8, 0, 17, 12}, fd: []int{5, 18, 26, 25, 3, 23, 4}}},
		{cmd: "set", key: 29, val: "twenty-nine", state: state{rl: []int{29, 9, 7, 24, 10, 11}, rd: []int{27, 22, 28}, fl: []int{1, 8, 0, 17}, fd: []int{12, 5, 18, 26, 25, 3, 23}}},
		{cmd: "get", key: 9, val: "nine", state: state{rl: []int{29, 7, 24, 10, 11}, rd: []int{27, 22, 28}, fl: []int{9, 1, 8, 0, 17}, fd: []int{12, 5, 18, 26, 25, 3, 23}}},
		{cmd: "set", key: 28, val: "twenty-eight", state: state{rl: []int{29, 7, 24, 10, 11}, rd: []int{27, 22}, fl: []int{28, 9, 1, 8, 0}, fd: []int{17, 12, 5, 18, 26, 25, 3, 23}}},
		{cmd: "set", key: 27, val: "twenty-seven", state: state{rl: []int{29, 7, 24, 10, 11}, rd: []int{22}, fl: []int{27, 28, 9, 1, 8}, fd: []int{0, 17, 12, 5, 18, 26, 25, 3, 23}}},
		{cmd: "set", key: 13, val: "thirteen", state: state{rl: []int{13, 29, 7, 24, 10, 11}, rd: []int{22}, fl: []int{27, 28, 9, 1}, fd: []int{8, 0, 17, 12, 5, 18, 26, 25, 3}}},
		{cmd: "set", key: 21, val: "twenty-one", state: state{rl: []int{21, 13, 29, 7, 24, 10, 11}, rd: []int{22}, fl: []int{27, 28, 9}, fd: []int{1, 8, 0, 17, 12, 5, 18, 26, 25}}},
		{cmd: "set", key: 8, val: "eight", state: state{rl: []int{21, 13, 29, 7, 24, 10, 11}, rd: []int{22}, fl: []int{8, 27, 28}, fd: []int{9, 1, 0, 17, 12, 5, 18, 26, 25}}},
		{cmd: "set", key: 3, val: "three", state: state{rl: []int{3, 21, 13, 29, 7, 24, 10, 11}, rd: []int{22}, fl: []int{8, 27}, fd: []int{28, 9, 1, 0, 17, 12, 5, 18, 26}}},
		{cmd: "set", key: 2, val: "two", state: state{rl: []int{2, 3, 21, 13, 29, 7, 24, 10, 11}, rd: []int{22}, fl: []int{8}, fd: []int{27, 28, 9, 1, 0, 17, 12, 5, 18}}},
		{cmd: "set", key: 0, val: "zero", state: state{rl: []int{2, 3, 21, 13, 29, 7, 24, 10}, rd: []int{11, 22}, fl: []int{0, 8}, fd: []int{27, 28, 9, 1, 17, 12, 5, 18}}},
		{cmd: "set", key: 18, val: "eighteen", state: state{rl: []int{2, 3, 21, 13, 29, 7, 24}, rd: []int{10, 11, 22}, fl: []int{18, 0, 8}, fd: []int{27, 28, 9, 1, 17, 12, 5}}},
		{cmd: "get", key: 13, val: "thirteen", state: state{rl: []int{2, 3, 21, 29, 7, 24}, rd: []int{10, 11, 22}, fl: []int{13, 18, 0, 8}, fd: []int{27, 28, 9, 1, 17, 12, 5}}},
		{cmd: "get", key: 21, val: "twenty-one", state: state{rl: []int{2, 3, 29, 7, 24}, rd: []int{10, 11, 22}, fl: []int{21, 13, 18, 0, 8}, fd: []int{27, 28, 9, 1, 17, 12, 5}}},
		{cmd: "set", key: 5, val: "five", state: state{rl: []int{2, 3, 29, 7, 24}, rd: []int{10, 11, 22}, fl: []int{5, 21, 13, 18, 0}, fd: []int{8, 27, 28, 9, 1, 17, 12}}},
		{cmd: "set", key: 14, val: "fourteen", state: state{rl: []int{14, 2, 3, 29, 7, 24}, rd: []int{10, 11, 22}, fl: []int{5, 21, 13, 18}, fd: []int{0, 8, 27, 28, 9, 1, 17}}},
		{cmd: "get", key: 3, val: "three", state: state{rl: []int{14, 2, 29, 7, 24}, rd: []int{10, 11, 22}, fl: []int{3, 5, 21, 13, 18}, fd: []int{0, 8, 27, 28, 9, 1, 17}}},
		{cmd: "set", key: 4, val: "four", state: state{rl: []int{4, 14, 2, 29, 7, 24}, rd: []int{10, 11, 22}, fl: []int{3, 5, 21, 13}, fd: []int{18, 0, 8, 27, 28, 9, 1}}},
		{cmd: "set", key: 8, val: "eight", state: state{rl: []int{4, 14, 2, 29, 7}, rd: []int{24, 10, 11, 22}, fl: []int{8, 3, 5, 21, 13}, fd: []int{18, 0, 27, 28, 9, 1}}},
		{cmd: "set", key: 19, val: "nineteen", state: state{rl: []int{19, 4, 14, 2, 29, 7}, rd: []int{24, 10, 11, 22}, fl: []int{8, 3, 5, 21}, fd: []int{13, 18, 0, 27, 28, 9}}},
		{cmd: "get", key: 14, val: "fourteen", state: state{rl: []int{19, 4, 2, 29, 7}, rd: []int{24, 10, 11, 22}, fl: []int{14, 8, 3, 5, 21}, fd: []int{13, 18, 0, 27, 28, 9}}},
		{cmd: "set", key: 11, val: "eleven", state: state{rl: []int{19, 4, 2, 29, 7}, rd: []int{24, 10, 22}, fl: []int{11, 14, 8, 3, 5}, fd: []int{21, 13, 18, 0, 27, 28, 9}}},
		{cmd: "get", key: 4, val: "four", state: state{rl: []int{19, 2, 29, 7}, rd: []int{24, 10, 22}, fl: []int{4, 11, 14, 8, 3, 5}, fd: []int{21, 13, 18, 0, 27, 28, 9}}},
		{cmd: "set", key: 17, val: "seventeen", state: state{rl: []int{17, 19, 2, 29, 7}, rd: []int{24, 10, 22}, fl: []int{4, 11, 14, 8, 3}, fd: []int{5, 21, 13, 18, 0, 27, 28}}},
		{cmd: "get", key: 19, val: "nineteen", state: state{rl: []int{17, 2, 29, 7}, rd: []int{24, 10, 22}, fl: []int{19, 4, 11, 14, 8, 3}, fd: []int{5, 21, 13, 18, 0, 27, 28}}},
		{cmd: "set", key: 12, val: "twelve", state: state{rl: []int{12, 17, 2, 29, 7}, rd: []int{24, 10, 22}, fl: []int{19, 4, 11, 14, 8}, fd: []int{3, 5, 21, 13, 18, 0, 27}}},
		{cmd: "set", key: 23, val: "twenty-three", state: state{rl: []int{23, 12, 17, 2, 29, 7}, rd: []int{24, 10, 22}, fl: []int{19, 4, 11, 14}, fd: []int{8, 3, 5, 21, 13, 18, 0}}},
		{cmd: "get", key: 29, val: "twenty-nine", state: state{rl: []int{23, 12, 17, 2, 7}, rd: []int{24, 10, 22}, fl: []int{29, 19, 4, 11, 14}, fd: []int{8, 3, 5, 21, 13, 18, 0}}},
		{cmd: "set", key: 25, val: "twenty-five", state: state{rl: []int{25, 23, 12, 17, 2, 7}, rd: []int{24, 10, 22}, fl: []int{29, 19, 4, 11}, fd: []int{14, 8, 3, 5, 21, 13, 18}}},
		{cmd: "get", key: 25, val: "twenty-five", state: state{rl: []int{23, 12, 17, 2, 7}, rd: []int{24, 10, 22}, fl: []int{25, 29, 19, 4, 11}, fd: []int{14, 8, 3, 5, 21, 13, 18}}},
		{cmd: "get", key: 19, val: "nineteen", state: state{rl: []int{23, 12, 17, 2, 7}, rd: []int{24, 10, 22}, fl: []int{19, 25, 29, 4, 11}, fd: []int{14, 8, 3, 5, 21, 13, 18}}},
		{cmd: "get", key: 23, val: "twenty-three", state: state{rl: []int{12, 17, 2, 7}, rd: []int{24, 10, 22}, fl: []int{23, 19, 25, 29, 4, 11}, fd: []int{14, 8, 3, 5, 21, 13, 18}}},
		{cmd: "get", key: 11, val: "eleven", state: state{rl: []int{12, 17, 2, 7}, rd: []int{24, 10, 22}, fl: []int{11, 23, 19, 25, 29, 4}, fd: []int{14, 8, 3, 5, 21, 13, 18}}},
		{cmd: "set", key: 18, val: "eighteen", state: state{rl: []int{12, 17, 2, 7}, rd: []int{24, 10, 22}, fl: []int{18, 11, 23, 19, 25, 29}, fd: []int{4, 14, 8, 3, 5, 21, 13}}},
		{cmd: "set", key: 24, val: "twenty-four", state: state{rl: []int{12, 17, 2, 7}, rd: []int{10, 22}, fl: []int{24, 18, 11, 23, 19, 25}, fd: []int{29, 4, 14, 8, 3, 5, 21, 13}}},
		{cmd: "set", key: 13, val: "thirteen", state: state{rl: []int{12, 17, 2, 7}, rd: []int{10, 22}, fl: []int{13, 24, 18, 11, 23, 19}, fd: []int{25, 29, 4, 14, 8, 3, 5, 21}}},
		{cmd: "get", key: 12, val: "twelve", state: state{rl: []int{17, 2, 7}, rd: []int{10, 22}, fl: []int{12, 13, 24, 18, 11, 23, 19}, fd: []int{25, 29, 4, 14, 8, 3, 5, 21}}},
		{cmd: "set", key: 27, val: "twenty-seven", state: state{rl: []int{27, 17, 2, 7}, rd: []int{10, 22}, fl: []int{12, 13, 24, 18, 11, 23}, fd: []int{19, 25, 29, 4, 14, 8, 3, 5}}},
		{cmd: "set", key: 19, val: "nineteen", state: state{rl: []int{27, 17, 2, 7}, rd: []int{10, 22}, fl: []int{19, 12, 13, 24, 18, 11}, fd: []int{23, 25, 29, 4, 14, 8, 3, 5}}},
		{cmd: "set", key: 16, val: "sixteen", state: state{rl: []int{16, 27, 17, 2, 7}, rd: []int{10, 22}, fl: []int{19, 12, 13, 24, 18}, fd: []int{11, 23, 25, 29, 4, 14, 8, 3}}},
		{cmd: "set", key: 3, val: "three", state: state{rl: []int{16, 27, 17, 2}, rd: []int{7, 10, 22}, fl: []int{3, 19, 12, 13, 24, 18}, fd: []int{11, 23, 25, 29, 4, 14, 8}}},
		{cmd: "set", key: 7, val: "seven", state: state{rl: []int{16, 27, 17, 2}, rd: []int{10, 22}, fl: []int{7, 3, 19, 12, 13, 24}, fd: []int{18, 11, 23, 25, 29, 4, 14, 8}}},
		{cmd: "set", key: 1, val: "one", state: state{rl: []int{1, 16, 27, 17, 2}, rd: []int{10, 22}, fl: []int{7, 3, 19, 12, 13}, fd: []int{24, 18, 11, 23, 25, 29, 4, 14}}},
		{cmd: "set", key: 26, val: "twenty-six", state: state{rl: []int{26, 1, 16, 27, 17, 2}, rd: []int{10, 22}, fl: []int{7, 3, 19, 12}, fd: []int{13, 24, 18, 11, 23, 25, 29, 4}}},
		{cmd: "get", key: 2, val: "two", state: state{rl: []int{26, 1, 16, 27, 17}, rd: []int{10, 22}, fl: []int{2, 7, 3, 19, 12}, fd: []int{13, 24, 18, 11, 23, 25, 29, 4}}},
		{cmd: "set", key: 15, val: "fifteen", state: state{rl: []int{15, 26, 1, 16, 27, 17}, rd: []int{10, 22}, fl: []int{2, 7, 3, 19}, fd: []int{12, 13, 24, 18, 11, 23, 25, 29}}},
		{cmd: "set", key: 23, val: "twenty-three", state: state{rl: []int{15, 26, 1, 16, 27}, rd: []int{17, 10, 22}, fl: []int{23, 2, 7, 3, 19}, fd: []int{12, 13, 24, 18, 11, 25, 29}}},
		{cmd: "get", key: 23, val: "twenty-three", state: state{rl: []int{15, 26, 1, 16, 27}, rd: []int{17, 10, 22}, fl: []int{23, 2, 7, 3, 19}, fd: []int{12, 13, 24, 18, 11, 25, 29}}},
		{cmd: "set", key: 18, val: "eighteen", state: state{rl: []int{15, 26, 1, 16}, rd: []int{27, 17, 10, 22}, fl: []int{18, 23, 2, 7, 3, 19}, fd: []int{12, 13, 24, 11, 25, 29}}},
		{cmd: "get", key: 1, val: "one", state: state{rl: []int{15, 26, 16}, rd: []int{27, 17, 10, 22}, fl: []int{1, 18, 23, 2, 7, 3, 19}, fd: []int{12, 13, 24, 11, 25, 29}}},
		{cmd: "set", key: 12, val: "twelve", state: state{rl: []int{15, 26}, rd: []int{16, 27, 17, 10, 22}, fl: []int{12, 1, 18, 23, 2, 7, 3, 19}, fd: []int{13, 24, 11, 25, 29}}},
		{cmd: "set", key: 20, val: "twenty", state: state{rl: []int{20, 15, 26}, rd: []int{16, 27, 17, 10, 22}, fl: []int{12, 1, 18, 23, 2, 7, 3}, fd: []int{19, 13, 24, 11, 25}}},
		{cmd: "set", key: 21, val: "twenty-one", state: state{rl: []int{21, 20, 15, 26}, rd: []int{16, 27, 17, 10, 22}, fl: []int{12, 1, 18, 23, 2, 7}, fd: []int{3, 19, 13, 24, 11}}},
		{cmd: "set", key: 6, val: "six", state: state{rl: []int{6, 21, 20, 15}, rd: []int{26, 16, 27, 17, 10, 22}, fl: []int{12, 1, 18, 23, 2, 7}, fd: []int{3, 19, 13, 24}}},
		{cmd: "set", key: 3, val: "three", state: state{rl: []int{6, 21, 20}, rd: []int{15, 26, 16, 27, 17, 10, 22}, fl: []int{3, 12, 1, 18, 23, 2, 7}, fd: []int{19, 13, 24}}},
		{cmd: "set", key: 19, val: "nineteen", state: state{rl: []int{6, 21}, rd: []int{20, 15, 26, 16, 27, 17, 10, 22}, fl: []int{19, 3, 12, 1, 18, 23, 2, 7}, fd: []int{13, 24}}},
		{cmd: "get", key: 19, val: "nineteen", state: state{rl: []int{6, 21}, rd: []int{20, 15, 26, 16, 27, 17, 10, 22}, fl: []int{19, 3, 12, 1, 18, 23, 2, 7}, fd: []int{13, 24}}},
		{cmd: "set", key: 15, val: "fifteen", state: state{rl: []int{6}, rd: []int{21, 20, 26, 16, 27, 17, 10, 22}, fl: []int{15, 19, 3, 12, 1, 18, 23, 2, 7}, fd: []int{13, 24}}},
		{cmd: "set", key: 8, val: "eight", state: state{rl: []int{8, 6}, rd: []int{21, 20, 26, 16, 27, 17, 10, 22}, fl: []int{15, 19, 3, 12, 1, 18, 23, 2}, fd: []int{7, 13}}},
		{cmd: "get", key: 3, val: "three", state: state{rl: []int{8, 6}, rd: []int{21, 20, 26, 16, 27, 17, 10, 22}, fl: []int{3, 15, 19, 12, 1, 18, 23, 2}, fd: []int{7, 13}}},
		{cmd: "get", key: 23, val: "twenty-three", state: state{rl: []int{8, 6}, rd: []int{21, 20, 26, 16, 27, 17, 10, 22}, fl: []int{23, 3, 15, 19, 12, 1, 18, 2}, fd: []int{7, 13}}},
		{cmd: "set", key: 29, val: "twenty-nine", state: state{rl: []int{29, 8}, rd: []int{6, 21, 20, 26, 16, 27, 17, 10}, fl: []int{23, 3, 15, 19, 12, 1, 18, 2}, fd: []int{7, 13}}},
		{cmd: "set", key: 7, val: "seven", state: state{rl: []int{29}, rd: []int{8, 6, 21, 20, 26, 16, 27, 17, 10}, fl: []int{7, 23, 3, 15, 19, 12, 1, 18, 2}, fd: []int{13}}},
		{cmd: "set", key: 10, val: "ten", state: state{rl: []int{29}, rd: []int{8, 6, 21, 20, 26, 16, 27, 17}, fl: []int{10, 7, 23, 3, 15, 19, 12, 1, 18}, fd: []int{2, 13}}},
		{cmd: "set", key: 9, val: "nine", state: state{rl: []int{9, 29}, rd: []int{8, 6, 21, 20, 26, 16, 27, 17}, fl: []int{10, 7, 23, 3, 15, 19, 12, 1}, fd: []int{18, 2}}},
		{cmd: "set", key: 28, val: "twenty-eight", state: state{rl: []int{28, 9}, rd: []int{29, 8, 6, 21, 20, 26, 16, 27}, fl: []int{10, 7, 23, 3, 15, 19, 12, 1}, fd: []int{18, 2}}},
		{cmd: "set", key: 6, val: "six", state: state{rl: []int{28, 9}, rd: []int{29, 8, 21, 20, 26, 16, 27}, fl: []int{6, 10, 7, 23, 3, 15, 19, 12}, fd: []int{1, 18, 2}}},
		{cmd: "get", key: 15, val: "fifteen", state: state{rl: []int{28, 9}, rd: []int{29, 8, 21, 20, 26, 16, 27}, fl: []int{15, 6, 10, 7, 23, 3, 19, 12}, fd: []int{1, 18, 2}}},
		{cmd: "get", key: 9, val: "nine", state: state{rl: []int{28}, rd: []int{29, 8, 21, 20, 26, 16, 27}, fl: []int{9, 15, 6, 10, 7, 23, 3, 19, 12}, fd: []int{1, 18, 2}}},
		{cmd: "get", key: 3, val: "three", state: state{rl: []int{28}, rd: []int{29, 8, 21, 20, 26, 16, 27}, fl: []int{3, 9, 15, 6, 10, 7, 23, 19, 12}, fd: []int{1, 18, 2}}},
		{cmd: "get", key: 6, val: "six", state: state{rl: []int{28}, rd: []int{29, 8, 21, 20, 26, 16, 27}, fl: []int{6, 3, 9, 15, 10, 7, 23, 19, 12}, fd: []int{1, 18, 2}}},
		{cmd: "set", key: 20, val: "twenty", state: state{rl: []int{28}, rd: []int{29, 8, 21, 26, 16, 27}, fl: []int{20, 6, 3, 9, 15, 10, 7, 23, 19}, fd: []int{12, 1, 18, 2}}},
		{cmd: "set", key: 8, val: "eight", state: state{rl: []int{28}, rd: []int{29, 21, 26, 16, 27}, fl: []int{8, 20, 6, 3, 9, 15, 10, 7, 23}, fd: []int{19, 12, 1, 18, 2}}},
		{cmd: "get", key: 6, val: "six", state: state{rl: []int{28}, rd: []int{29, 21, 26, 16, 27}, fl: []int{6, 8, 20, 3, 9, 15, 10, 7, 23}, fd: []int{19, 12, 1, 18, 2}}},
		{cmd: "get", key: 9, val: "nine", state: state{rl: []int{28}, rd: []int{29, 21, 26, 16, 27}, fl: []int{9, 6, 8, 20, 3, 15, 10, 7, 23}, fd: []int{19, 12, 1, 18, 2}}},
		{cmd: "get", key: 3, val: "three", state: state{rl: []int{28}, rd: []int{29, 21, 26, 16, 27}, fl: []int{3, 9, 6, 8, 20, 15, 10, 7, 23}, fd: []int{19, 12, 1, 18, 2}}},
		{cmd: "set", key: 26, val: "twenty-six", state: state{rl: []int{28}, rd: []int{29, 21, 16, 27}, fl: []int{26, 3, 9, 6, 8, 20, 15, 10, 7}, fd: []int{23, 19, 12, 1, 18, 2}}},
		{cmd: "get", key: 20, val: "twenty", state: state{rl: []int{28}, rd: []int{29, 21, 16, 27}, fl: []int{20, 26, 3, 9, 6, 8, 15, 10, 7}, fd: []int{23, 19, 12, 1, 18, 2}}},
		{cmd: "set", key: 17, val: "seventeen", state: state{rl: []int{17, 28}, rd: []int{29, 21, 16, 27}, fl: []int{20, 26, 3, 9, 6, 8, 15, 10}, fd: []int{7, 23, 19, 12, 1, 18}}},
		{cmd: "set", key: 22, val: "twenty-two", state: state{rl: []int{22, 17, 28}, rd: []int{29, 21, 16, 27}, fl: []int{20, 26, 3, 9, 6, 8, 15}, fd: []int{10, 7, 23, 19, 12, 1}}},
		{cmd: "set", key: 10, val: "ten", state: state{rl: []int{22, 17, 28}, rd: []int{29, 21, 16, 27}, fl: []int{10, 20, 26, 3, 9, 6, 8}, fd: []int{15, 7, 23, 19, 12, 1}}},
		{cmd: "set", key: 2, val: "two", state: state{rl: []int{2, 22, 17, 28}, rd: []int{29, 21, 16, 27}, fl: []int{10, 20, 26, 3, 9, 6}, fd: []int{8, 15, 7, 23, 19, 12}}},
		{cmd: "set", key: 7, val: "seven", state: state{rl: []int{2, 22, 17}, rd: []int{28, 29, 21, 16, 27}, fl: []int{7, 10, 20, 26, 3, 9, 6}, fd: []int{8, 15, 23, 19, 12}}},
		{cmd: "get", key: 7, val: "seven", state: state{rl: []int{2, 22, 17}, rd: []int{28, 29, 21, 16, 27}, fl: []int{7, 10, 20, 26, 3, 9, 6}, fd: []int{8, 15, 23, 19, 12}}},
		{cmd: "set", key: 15, val: "fifteen", state: state{rl: []int{2, 22}, rd: []int{17, 28, 29, 21, 16, 27}, fl: []int{15, 7, 10, 20, 26, 3, 9, 6}, fd: []int{8, 23, 19, 12}}},
		{cmd: "set", key: 1, val: "one", state: state{rl: []int{1, 2, 22}, rd: []int{17, 28, 29, 21, 16, 27}, fl: []int{15, 7, 10, 20, 26, 3, 9}, fd: []int{6, 8, 23, 19}}},
		{cmd: "get", key: 26, val: "twenty-six", state: state{rl: []int{1, 2, 22}, rd: []int{17, 28, 29, 21, 16, 27}, fl: []int{26, 15, 7, 10, 20, 3, 9}, fd: []int{6, 8, 23, 19}}},
		{cmd: "set", key: 21, val: "twenty-one", state: state{rl: []int{1, 2, 22}, rd: []int{17, 28, 29, 16, 27}, fl: []int{21, 26, 15, 7, 10, 20, 3}, fd: []int{9, 6, 8, 23, 19}}},
	}
	ints := func(l *list.List[item[int, string]]) []int {
		a := make([]int, l.Len())
		e := l.Front()
		for i := range a {
			a[i] = e.Value.key
			e = e.Next()
		}
		return a
	}
	c := New[int, string](10)
	for i, tt := range tests {
		var cmd string
		switch tt.cmd {
		case "get":
			cmd = fmt.Sprintf("Get(%d)", tt.key)
			if val, _ := c.Get(tt.key); tt.val != val {
				t.Fatalf("step %d: %s: unexpected value; got: %q; want: %q", i, cmd, val, tt.val)
			}
		case "set":
			cmd = fmt.Sprintf("Set(%d, %q)", tt.key, tt.val)
			c.Set(tt.key, tt.val)
		default:
			t.Fatalf("step %d: unexpected command: %q", i, tt.cmd)
		}
		got := state{rl: ints(c.rl), rd: ints(c.rd), fl: ints(c.fl), fd: ints(c.fd)}
		if !reflect.DeepEqual(got, tt.state) {
			var prev state
			if i > 0 {
				prev = tests[i-1].state
			}
			t.Fatalf("step %d: %s: unexpected state:\nprev %+v\ngot  %+v\nwant %+v", i, cmd, prev, got, tt.state)
		}
	}
}
