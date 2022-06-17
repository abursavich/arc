// Copyright 2015 Andrew Bursavich. All rights reserved.
// Use of this source code is governed by The MIT License
// which can be found in the LICENSE file.

package arc

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"bursavich.dev/arc/internal/list"
)

type state[K comparable] [4][]K

func cacheState[K comparable, V any](c *Cache[K, V]) state[K] {
	return state[K]{
		reverse(keys(&c.dead.mru)),
		reverse(keys(&c.live.mru)),
		keys(&c.live.mfu),
		keys(&c.dead.mfu),
	}
}

func (s state[K]) String() string {
	return fmt.Sprintf(
		"[ %s ... %s | %s ... %s ]",
		strings.Join(toStrings(s[0]), ", "),
		strings.Join(toStrings(s[1]), ", "),
		strings.Join(toStrings(s[2]), ", "),
		strings.Join(toStrings(s[3]), ", "),
	)
}

func (s state[K]) GoString() string {
	var typ K
	return fmt.Sprintf(
		"state[%T]{{%s}, {%s}, {%s}, {%s}}",
		typ,
		strings.Join(toStrings(s[0]), ","),
		strings.Join(toStrings(s[1]), ","),
		strings.Join(toStrings(s[2]), ","),
		strings.Join(toStrings(s[3]), ","),
	)
}

func keys[K comparable, V any](l *list.List[item[K, V]]) []K {
	a := make([]K, l.Len())
	e := l.Front()
	for i := range a {
		a[i] = e.Value.key
		e = e.Next()
	}
	return a
}

func reverse[K any](s []K) []K {
	n := len(s)
	for i := 0; i < n/2; i++ {
		s[i], s[n-i-1] = s[n-i-1], s[i]
	}
	return s
}

func toStrings[K any](s []K) []string {
	r := make([]string, len(s))
	for i, v := range s {
		r[i] = fmt.Sprint(v)
	}
	return r
}

func TestTable(t *testing.T) {
	tests := []struct {
		cmd   string
		key   int
		val   string
		state state[int]
	}{
		{cmd: "set", key: 0, val: "0", state: state[int]{{}, {0}, {}, {}}},
		{cmd: "set", key: 2, val: "2", state: state[int]{{}, {0, 2}, {}, {}}},
		{cmd: "set", key: 0, val: "0", state: state[int]{{}, {2}, {0}, {}}},
		{cmd: "get", key: 3, val: "", state: state[int]{{}, {2}, {0}, {}}},
		{cmd: "set", key: 2, val: "2", state: state[int]{{}, {}, {2, 0}, {}}},
		{cmd: "set", key: 5, val: "5", state: state[int]{{}, {5}, {2, 0}, {}}},
		{cmd: "get", key: 7, val: "", state: state[int]{{}, {5}, {2, 0}, {}}},
		{cmd: "get", key: 5, val: "5", state: state[int]{{}, {}, {5, 2, 0}, {}}},
		{cmd: "set", key: 0, val: "0", state: state[int]{{}, {}, {0, 5, 2}, {}}},
		{cmd: "get", key: 3, val: "", state: state[int]{{}, {}, {0, 5, 2}, {}}},
		{cmd: "set", key: 2, val: "2", state: state[int]{{}, {}, {2, 0, 5}, {}}},
		{cmd: "set", key: 1, val: "1", state: state[int]{{}, {1}, {2, 0, 5}, {}}},
		{cmd: "get", key: 9, val: "", state: state[int]{{}, {1}, {2, 0, 5}, {}}},
		{cmd: "set", key: 2, val: "2", state: state[int]{{}, {1}, {2, 0, 5}, {}}},
		{cmd: "get", key: 7, val: "", state: state[int]{{}, {1}, {2, 0, 5}, {}}},
		{cmd: "get", key: 2, val: "2", state: state[int]{{}, {1}, {2, 0, 5}, {}}},
		{cmd: "get", key: 9, val: "", state: state[int]{{}, {1}, {2, 0, 5}, {}}},
		{cmd: "set", key: 6, val: "6", state: state[int]{{}, {1, 6}, {2, 0, 5}, {}}},
		{cmd: "set", key: 9, val: "9", state: state[int]{{}, {1, 6, 9}, {2, 0, 5}, {}}},
		{cmd: "get", key: 5, val: "5", state: state[int]{{}, {1, 6, 9}, {5, 2, 0}, {}}},
		{cmd: "get", key: 16, val: "", state: state[int]{{}, {1, 6, 9}, {5, 2, 0}, {}}},
		{cmd: "set", key: 5, val: "5", state: state[int]{{}, {1, 6, 9}, {5, 2, 0}, {}}},
		{cmd: "set", key: 9, val: "9", state: state[int]{{}, {1, 6}, {9, 5, 2, 0}, {}}},
		{cmd: "get", key: 16, val: "", state: state[int]{{}, {1, 6}, {9, 5, 2, 0}, {}}},
		{cmd: "set", key: 0, val: "0", state: state[int]{{}, {1, 6}, {0, 9, 5, 2}, {}}},
		{cmd: "set", key: 1, val: "1", state: state[int]{{}, {6}, {1, 0, 9, 5, 2}, {}}},
		{cmd: "set", key: 7, val: "7", state: state[int]{{6}, {7}, {1, 0, 9, 5, 2}, {}}},
		{cmd: "set", key: 7, val: "7", state: state[int]{{6}, {}, {7, 1, 0, 9, 5, 2}, {}}},
		{cmd: "set", key: 2, val: "2", state: state[int]{{6}, {}, {2, 7, 1, 0, 9, 5}, {}}},
		{cmd: "get", key: 3, val: "", state: state[int]{{6}, {}, {2, 7, 1, 0, 9, 5}, {}}},
		{cmd: "get", key: 0, val: "0", state: state[int]{{6}, {}, {0, 2, 7, 1, 9, 5}, {}}},
		{cmd: "get", key: 4, val: "", state: state[int]{{6}, {}, {0, 2, 7, 1, 9, 5}, {}}},
		{cmd: "set", key: 1, val: "1", state: state[int]{{6}, {}, {1, 0, 2, 7, 9, 5}, {}}},
		{cmd: "get", key: 2, val: "2", state: state[int]{{6}, {}, {2, 1, 0, 7, 9, 5}, {}}},
		{cmd: "set", key: 10, val: "10", state: state[int]{{6}, {10}, {2, 1, 0, 7, 9}, {5}}},
		{cmd: "get", key: 9, val: "9", state: state[int]{{6}, {10}, {9, 2, 1, 0, 7}, {5}}},
		{cmd: "get", key: 17, val: "", state: state[int]{{6}, {10}, {9, 2, 1, 0, 7}, {5}}},
		{cmd: "get", key: 12, val: "", state: state[int]{{6}, {10}, {9, 2, 1, 0, 7}, {5}}},
		{cmd: "set", key: 11, val: "11", state: state[int]{{6, 10}, {11}, {9, 2, 1, 0, 7}, {5}}},
		{cmd: "set", key: 4, val: "4", state: state[int]{{6, 10, 11}, {4}, {9, 2, 1, 0, 7}, {5}}},
		{cmd: "set", key: 8, val: "8", state: state[int]{{6, 10, 11, 4}, {8}, {9, 2, 1, 0, 7}, {5}}},
		{cmd: "get", key: 8, val: "8", state: state[int]{{6, 10, 11, 4}, {}, {8, 9, 2, 1, 0, 7}, {5}}},
		{cmd: "get", key: 7, val: "7", state: state[int]{{6, 10, 11, 4}, {}, {7, 8, 9, 2, 1, 0}, {5}}},
		{cmd: "get", key: 11, val: "", state: state[int]{{6, 10, 11, 4}, {}, {7, 8, 9, 2, 1, 0}, {5}}},
		{cmd: "get", key: 8, val: "8", state: state[int]{{6, 10, 11, 4}, {}, {8, 7, 9, 2, 1, 0}, {5}}},
		{cmd: "get", key: 5, val: "", state: state[int]{{6, 10, 11, 4}, {}, {8, 7, 9, 2, 1, 0}, {5}}},
		{cmd: "get", key: 0, val: "0", state: state[int]{{6, 10, 11, 4}, {}, {0, 8, 7, 9, 2, 1}, {5}}},
		{cmd: "set", key: 8, val: "8", state: state[int]{{6, 10, 11, 4}, {}, {8, 0, 7, 9, 2, 1}, {5}}},
		{cmd: "set", key: 16, val: "16", state: state[int]{{6, 10, 11, 4}, {16}, {8, 0, 7, 9, 2}, {1, 5}}},
		{cmd: "set", key: 3, val: "3", state: state[int]{{6, 10, 11, 4, 16}, {3}, {8, 0, 7, 9, 2}, {1}}},
		{cmd: "set", key: 2, val: "2", state: state[int]{{6, 10, 11, 4, 16}, {3}, {2, 8, 0, 7, 9}, {1}}},
		{cmd: "get", key: 13, val: "", state: state[int]{{6, 10, 11, 4, 16}, {3}, {2, 8, 0, 7, 9}, {1}}},
		{cmd: "get", key: 6, val: "", state: state[int]{{6, 10, 11, 4, 16}, {3}, {2, 8, 0, 7, 9}, {1}}},
		{cmd: "get", key: 6, val: "", state: state[int]{{6, 10, 11, 4, 16}, {3}, {2, 8, 0, 7, 9}, {1}}},
		{cmd: "get", key: 12, val: "", state: state[int]{{6, 10, 11, 4, 16}, {3}, {2, 8, 0, 7, 9}, {1}}},
		{cmd: "get", key: 1, val: "", state: state[int]{{6, 10, 11, 4, 16}, {3}, {2, 8, 0, 7, 9}, {1}}},
		{cmd: "get", key: 17, val: "", state: state[int]{{6, 10, 11, 4, 16}, {3}, {2, 8, 0, 7, 9}, {1}}},
		{cmd: "get", key: 9, val: "9", state: state[int]{{6, 10, 11, 4, 16}, {3}, {9, 2, 8, 0, 7}, {1}}},
		{cmd: "set", key: 10, val: "10", state: state[int]{{6, 11, 4, 16}, {3}, {10, 9, 2, 8, 0}, {7, 1}}},
		{cmd: "set", key: 9, val: "9", state: state[int]{{6, 11, 4, 16}, {3}, {9, 10, 2, 8, 0}, {7, 1}}},
		{cmd: "get", key: 4, val: "", state: state[int]{{6, 11, 4, 16}, {3}, {9, 10, 2, 8, 0}, {7, 1}}},
		{cmd: "get", key: 17, val: "", state: state[int]{{6, 11, 4, 16}, {3}, {9, 10, 2, 8, 0}, {7, 1}}},
		{cmd: "set", key: 9, val: "9", state: state[int]{{6, 11, 4, 16}, {3}, {9, 10, 2, 8, 0}, {7, 1}}},
		{cmd: "get", key: 3, val: "3", state: state[int]{{6, 11, 4, 16}, {}, {3, 9, 10, 2, 8, 0}, {7, 1}}},
		{cmd: "get", key: 4, val: "", state: state[int]{{6, 11, 4, 16}, {}, {3, 9, 10, 2, 8, 0}, {7, 1}}},
		{cmd: "set", key: 9, val: "9", state: state[int]{{6, 11, 4, 16}, {}, {9, 3, 10, 2, 8, 0}, {7, 1}}},
		{cmd: "set", key: 2, val: "2", state: state[int]{{6, 11, 4, 16}, {}, {2, 9, 3, 10, 8, 0}, {7, 1}}},
		{cmd: "set", key: 1, val: "1", state: state[int]{{6, 11, 4, 16}, {}, {1, 2, 9, 3, 10, 8}, {0, 7}}},
		{cmd: "set", key: 6, val: "6", state: state[int]{{11, 4, 16}, {}, {6, 1, 2, 9, 3, 10}, {8, 0, 7}}},
		{cmd: "get", key: 6, val: "6", state: state[int]{{11, 4, 16}, {}, {6, 1, 2, 9, 3, 10}, {8, 0, 7}}},
		{cmd: "get", key: 5, val: "", state: state[int]{{11, 4, 16}, {}, {6, 1, 2, 9, 3, 10}, {8, 0, 7}}},
		{cmd: "get", key: 4, val: "", state: state[int]{{11, 4, 16}, {}, {6, 1, 2, 9, 3, 10}, {8, 0, 7}}},
		{cmd: "set", key: 2, val: "2", state: state[int]{{11, 4, 16}, {}, {2, 6, 1, 9, 3, 10}, {8, 0, 7}}},
		{cmd: "set", key: 10, val: "10", state: state[int]{{11, 4, 16}, {}, {10, 2, 6, 1, 9, 3}, {8, 0, 7}}},
		{cmd: "get", key: 4, val: "", state: state[int]{{11, 4, 16}, {}, {10, 2, 6, 1, 9, 3}, {8, 0, 7}}},
		{cmd: "set", key: 12, val: "12", state: state[int]{{11, 4, 16}, {12}, {10, 2, 6, 1, 9}, {3, 8, 0}}},
		{cmd: "get", key: 15, val: "", state: state[int]{{11, 4, 16}, {12}, {10, 2, 6, 1, 9}, {3, 8, 0}}},
		{cmd: "set", key: 9, val: "9", state: state[int]{{11, 4, 16}, {12}, {9, 10, 2, 6, 1}, {3, 8, 0}}},
		{cmd: "get", key: 0, val: "", state: state[int]{{11, 4, 16}, {12}, {9, 10, 2, 6, 1}, {3, 8, 0}}},
		{cmd: "set", key: 15, val: "15", state: state[int]{{11, 4, 16}, {12, 15}, {9, 10, 2, 6}, {1, 3, 8}}},
		{cmd: "set", key: 10, val: "10", state: state[int]{{11, 4, 16}, {12, 15}, {10, 9, 2, 6}, {1, 3, 8}}},
		{cmd: "set", key: 13, val: "13", state: state[int]{{11, 4, 16, 12}, {15, 13}, {10, 9, 2, 6}, {1, 3}}},
		{cmd: "set", key: 15, val: "15", state: state[int]{{11, 4, 16, 12}, {13}, {15, 10, 9, 2, 6}, {1, 3}}},
		{cmd: "set", key: 12, val: "12", state: state[int]{{11, 4, 16}, {13}, {12, 15, 10, 9, 2}, {6, 1, 3}}},
		{cmd: "get", key: 10, val: "10", state: state[int]{{11, 4, 16}, {13}, {10, 12, 15, 9, 2}, {6, 1, 3}}},
		{cmd: "set", key: 3, val: "3", state: state[int]{{11, 4, 16, 13}, {}, {3, 10, 12, 15, 9, 2}, {6, 1}}},
		{cmd: "set", key: 3, val: "3", state: state[int]{{11, 4, 16, 13}, {}, {3, 10, 12, 15, 9, 2}, {6, 1}}},
		{cmd: "get", key: 14, val: "", state: state[int]{{11, 4, 16, 13}, {}, {3, 10, 12, 15, 9, 2}, {6, 1}}},
		{cmd: "get", key: 2, val: "2", state: state[int]{{11, 4, 16, 13}, {}, {2, 3, 10, 12, 15, 9}, {6, 1}}},
		{cmd: "get", key: 1, val: "", state: state[int]{{11, 4, 16, 13}, {}, {2, 3, 10, 12, 15, 9}, {6, 1}}},
		{cmd: "get", key: 16, val: "", state: state[int]{{11, 4, 16, 13}, {}, {2, 3, 10, 12, 15, 9}, {6, 1}}},
		{cmd: "set", key: 4, val: "4", state: state[int]{{11, 16, 13}, {}, {4, 2, 3, 10, 12, 15}, {9, 6, 1}}},
		{cmd: "get", key: 3, val: "3", state: state[int]{{11, 16, 13}, {}, {3, 4, 2, 10, 12, 15}, {9, 6, 1}}},
		{cmd: "set", key: 1, val: "1", state: state[int]{{11, 16, 13}, {}, {1, 3, 4, 2, 10, 12}, {15, 9, 6}}},
		{cmd: "get", key: 9, val: "", state: state[int]{{11, 16, 13}, {}, {1, 3, 4, 2, 10, 12}, {15, 9, 6}}},
		{cmd: "get", key: 17, val: "", state: state[int]{{11, 16, 13}, {}, {1, 3, 4, 2, 10, 12}, {15, 9, 6}}},
		{cmd: "get", key: 7, val: "", state: state[int]{{11, 16, 13}, {}, {1, 3, 4, 2, 10, 12}, {15, 9, 6}}},
		{cmd: "set", key: 15, val: "15", state: state[int]{{11, 16, 13}, {}, {15, 1, 3, 4, 2, 10}, {12, 9, 6}}},
		{cmd: "set", key: 12, val: "12", state: state[int]{{11, 16, 13}, {}, {12, 15, 1, 3, 4, 2}, {10, 9, 6}}},
		{cmd: "get", key: 1, val: "1", state: state[int]{{11, 16, 13}, {}, {1, 12, 15, 3, 4, 2}, {10, 9, 6}}},
		{cmd: "set", key: 17, val: "17", state: state[int]{{11, 16, 13}, {17}, {1, 12, 15, 3, 4}, {2, 10, 9}}},
		{cmd: "get", key: 3, val: "3", state: state[int]{{11, 16, 13}, {17}, {3, 1, 12, 15, 4}, {2, 10, 9}}},
		{cmd: "get", key: 3, val: "3", state: state[int]{{11, 16, 13}, {17}, {3, 1, 12, 15, 4}, {2, 10, 9}}},
		{cmd: "get", key: 1, val: "1", state: state[int]{{11, 16, 13}, {17}, {1, 3, 12, 15, 4}, {2, 10, 9}}},
		{cmd: "get", key: 10, val: "", state: state[int]{{11, 16, 13}, {17}, {1, 3, 12, 15, 4}, {2, 10, 9}}},
		{cmd: "get", key: 7, val: "", state: state[int]{{11, 16, 13}, {17}, {1, 3, 12, 15, 4}, {2, 10, 9}}},
		{cmd: "get", key: 9, val: "", state: state[int]{{11, 16, 13}, {17}, {1, 3, 12, 15, 4}, {2, 10, 9}}},
		{cmd: "set", key: 8, val: "8", state: state[int]{{11, 16, 13, 17}, {8}, {1, 3, 12, 15, 4}, {2, 10}}},
		{cmd: "set", key: 9, val: "9", state: state[int]{{11, 16, 13, 17, 8}, {9}, {1, 3, 12, 15, 4}, {2}}},
		{cmd: "get", key: 9, val: "9", state: state[int]{{11, 16, 13, 17, 8}, {}, {9, 1, 3, 12, 15, 4}, {2}}},
		{cmd: "set", key: 11, val: "11", state: state[int]{{16, 13, 17, 8}, {}, {11, 9, 1, 3, 12, 15}, {4, 2}}},
		{cmd: "set", key: 7, val: "7", state: state[int]{{16, 13, 17, 8}, {7}, {11, 9, 1, 3, 12}, {15, 4}}},
		{cmd: "get", key: 17, val: "", state: state[int]{{16, 13, 17, 8}, {7}, {11, 9, 1, 3, 12}, {15, 4}}},
		{cmd: "get", key: 7, val: "7", state: state[int]{{16, 13, 17, 8}, {}, {7, 11, 9, 1, 3, 12}, {15, 4}}},
		{cmd: "set", key: 16, val: "16", state: state[int]{{13, 17, 8}, {}, {16, 7, 11, 9, 1, 3}, {12, 15, 4}}},
		{cmd: "set", key: 17, val: "17", state: state[int]{{13, 8}, {}, {17, 16, 7, 11, 9, 1}, {3, 12, 15, 4}}},
		{cmd: "set", key: 8, val: "8", state: state[int]{{13}, {}, {8, 17, 16, 7, 11, 9}, {1, 3, 12, 15, 4}}},
		{cmd: "set", key: 6, val: "6", state: state[int]{{13}, {6}, {8, 17, 16, 7, 11}, {9, 1, 3, 12, 15}}},
		{cmd: "get", key: 5, val: "", state: state[int]{{13}, {6}, {8, 17, 16, 7, 11}, {9, 1, 3, 12, 15}}},
		{cmd: "get", key: 9, val: "", state: state[int]{{13}, {6}, {8, 17, 16, 7, 11}, {9, 1, 3, 12, 15}}},
		{cmd: "set", key: 0, val: "0", state: state[int]{{13}, {6, 0}, {8, 17, 16, 7}, {11, 9, 1, 3, 12}}},
		{cmd: "get", key: 0, val: "0", state: state[int]{{13}, {6}, {0, 8, 17, 16, 7}, {11, 9, 1, 3, 12}}},
		{cmd: "get", key: 8, val: "8", state: state[int]{{13}, {6}, {8, 0, 17, 16, 7}, {11, 9, 1, 3, 12}}},
		{cmd: "set", key: 6, val: "6", state: state[int]{{13}, {}, {6, 8, 0, 17, 16, 7}, {11, 9, 1, 3, 12}}},
		{cmd: "set", key: 11, val: "11", state: state[int]{{13}, {}, {11, 6, 8, 0, 17, 16}, {7, 9, 1, 3, 12}}},
		{cmd: "set", key: 6, val: "6", state: state[int]{{13}, {}, {6, 11, 8, 0, 17, 16}, {7, 9, 1, 3, 12}}},
		{cmd: "get", key: 5, val: "", state: state[int]{{13}, {}, {6, 11, 8, 0, 17, 16}, {7, 9, 1, 3, 12}}},
		{cmd: "get", key: 5, val: "", state: state[int]{{13}, {}, {6, 11, 8, 0, 17, 16}, {7, 9, 1, 3, 12}}},
		{cmd: "set", key: 3, val: "3", state: state[int]{{13}, {}, {3, 6, 11, 8, 0, 17}, {16, 7, 9, 1, 12}}},
		{cmd: "get", key: 0, val: "0", state: state[int]{{13}, {}, {0, 3, 6, 11, 8, 17}, {16, 7, 9, 1, 12}}},
		{cmd: "get", key: 13, val: "", state: state[int]{{13}, {}, {0, 3, 6, 11, 8, 17}, {16, 7, 9, 1, 12}}},
		{cmd: "get", key: 6, val: "6", state: state[int]{{13}, {}, {6, 0, 3, 11, 8, 17}, {16, 7, 9, 1, 12}}},
		{cmd: "get", key: 14, val: "", state: state[int]{{13}, {}, {6, 0, 3, 11, 8, 17}, {16, 7, 9, 1, 12}}},
		{cmd: "get", key: 14, val: "", state: state[int]{{13}, {}, {6, 0, 3, 11, 8, 17}, {16, 7, 9, 1, 12}}},
		{cmd: "get", key: 4, val: "", state: state[int]{{13}, {}, {6, 0, 3, 11, 8, 17}, {16, 7, 9, 1, 12}}},
		{cmd: "get", key: 0, val: "0", state: state[int]{{13}, {}, {0, 6, 3, 11, 8, 17}, {16, 7, 9, 1, 12}}},
		{cmd: "set", key: 0, val: "0", state: state[int]{{13}, {}, {0, 6, 3, 11, 8, 17}, {16, 7, 9, 1, 12}}},
		{cmd: "set", key: 16, val: "16", state: state[int]{{13}, {}, {16, 0, 6, 3, 11, 8}, {17, 7, 9, 1, 12}}},
		{cmd: "get", key: 16, val: "16", state: state[int]{{13}, {}, {16, 0, 6, 3, 11, 8}, {17, 7, 9, 1, 12}}},
		{cmd: "set", key: 14, val: "14", state: state[int]{{13}, {14}, {16, 0, 6, 3, 11}, {8, 17, 7, 9, 1}}},
		{cmd: "set", key: 6, val: "6", state: state[int]{{13}, {14}, {6, 16, 0, 3, 11}, {8, 17, 7, 9, 1}}},
		{cmd: "set", key: 8, val: "8", state: state[int]{{13, 14}, {}, {8, 6, 16, 0, 3, 11}, {17, 7, 9, 1}}},
		{cmd: "set", key: 16, val: "16", state: state[int]{{13, 14}, {}, {16, 8, 6, 0, 3, 11}, {17, 7, 9, 1}}},
		{cmd: "set", key: 17, val: "17", state: state[int]{{13, 14}, {}, {17, 16, 8, 6, 0, 3}, {11, 7, 9, 1}}},
		{cmd: "set", key: 10, val: "10", state: state[int]{{13, 14}, {10}, {17, 16, 8, 6, 0}, {3, 11, 7, 9}}},
		{cmd: "set", key: 17, val: "17", state: state[int]{{13, 14}, {10}, {17, 16, 8, 6, 0}, {3, 11, 7, 9}}},
		{cmd: "set", key: 3, val: "3", state: state[int]{{13, 14, 10}, {}, {3, 17, 16, 8, 6, 0}, {11, 7, 9}}},
		{cmd: "get", key: 7, val: "", state: state[int]{{13, 14, 10}, {}, {3, 17, 16, 8, 6, 0}, {11, 7, 9}}},
		{cmd: "set", key: 7, val: "7", state: state[int]{{13, 14, 10}, {}, {7, 3, 17, 16, 8, 6}, {0, 11, 9}}},
		{cmd: "set", key: 6, val: "6", state: state[int]{{13, 14, 10}, {}, {6, 7, 3, 17, 16, 8}, {0, 11, 9}}},
		{cmd: "get", key: 15, val: "", state: state[int]{{13, 14, 10}, {}, {6, 7, 3, 17, 16, 8}, {0, 11, 9}}},
		{cmd: "set", key: 4, val: "4", state: state[int]{{13, 14, 10}, {4}, {6, 7, 3, 17, 16}, {8, 0, 11}}},
		{cmd: "set", key: 11, val: "11", state: state[int]{{13, 14, 10, 4}, {}, {11, 6, 7, 3, 17, 16}, {8, 0}}},
		{cmd: "get", key: 8, val: "", state: state[int]{{13, 14, 10, 4}, {}, {11, 6, 7, 3, 17, 16}, {8, 0}}},
		{cmd: "set", key: 8, val: "8", state: state[int]{{13, 14, 10, 4}, {}, {8, 11, 6, 7, 3, 17}, {16, 0}}},
		{cmd: "get", key: 2, val: "", state: state[int]{{13, 14, 10, 4}, {}, {8, 11, 6, 7, 3, 17}, {16, 0}}},
		{cmd: "set", key: 15, val: "15", state: state[int]{{13, 14, 10, 4}, {15}, {8, 11, 6, 7, 3}, {17, 16}}},
		{cmd: "get", key: 17, val: "", state: state[int]{{13, 14, 10, 4}, {15}, {8, 11, 6, 7, 3}, {17, 16}}},
		{cmd: "set", key: 15, val: "15", state: state[int]{{13, 14, 10, 4}, {}, {15, 8, 11, 6, 7, 3}, {17, 16}}},
		{cmd: "set", key: 6, val: "6", state: state[int]{{13, 14, 10, 4}, {}, {6, 15, 8, 11, 7, 3}, {17, 16}}},
		{cmd: "set", key: 5, val: "5", state: state[int]{{13, 14, 10, 4}, {5}, {6, 15, 8, 11, 7}, {3, 17}}},
		{cmd: "set", key: 17, val: "17", state: state[int]{{13, 14, 10, 4, 5}, {}, {17, 6, 15, 8, 11, 7}, {3}}},
		{cmd: "get", key: 0, val: "", state: state[int]{{13, 14, 10, 4, 5}, {}, {17, 6, 15, 8, 11, 7}, {3}}},
		{cmd: "get", key: 12, val: "", state: state[int]{{13, 14, 10, 4, 5}, {}, {17, 6, 15, 8, 11, 7}, {3}}},
		{cmd: "set", key: 6, val: "6", state: state[int]{{13, 14, 10, 4, 5}, {}, {6, 17, 15, 8, 11, 7}, {3}}},
		{cmd: "set", key: 6, val: "6", state: state[int]{{13, 14, 10, 4, 5}, {}, {6, 17, 15, 8, 11, 7}, {3}}},
		{cmd: "set", key: 12, val: "12", state: state[int]{{13, 14, 10, 4, 5}, {12}, {6, 17, 15, 8, 11}, {7}}},
		{cmd: "get", key: 10, val: "", state: state[int]{{13, 14, 10, 4, 5}, {12}, {6, 17, 15, 8, 11}, {7}}},
		{cmd: "get", key: 15, val: "15", state: state[int]{{13, 14, 10, 4, 5}, {12}, {15, 6, 17, 8, 11}, {7}}},
		{cmd: "set", key: 10, val: "10", state: state[int]{{13, 14, 4, 5}, {12}, {10, 15, 6, 17, 8}, {11, 7}}},
		{cmd: "set", key: 17, val: "17", state: state[int]{{13, 14, 4, 5}, {12}, {17, 10, 15, 6, 8}, {11, 7}}},
		{cmd: "get", key: 16, val: "", state: state[int]{{13, 14, 4, 5}, {12}, {17, 10, 15, 6, 8}, {11, 7}}},
		{cmd: "get", key: 7, val: "", state: state[int]{{13, 14, 4, 5}, {12}, {17, 10, 15, 6, 8}, {11, 7}}},
		{cmd: "set", key: 14, val: "14", state: state[int]{{13, 4, 5}, {12}, {14, 17, 10, 15, 6}, {8, 11, 7}}},
		{cmd: "set", key: 0, val: "0", state: state[int]{{13, 4, 5}, {12, 0}, {14, 17, 10, 15}, {6, 8, 11}}},
		{cmd: "get", key: 0, val: "0", state: state[int]{{13, 4, 5}, {12}, {0, 14, 17, 10, 15}, {6, 8, 11}}},
		{cmd: "set", key: 2, val: "2", state: state[int]{{13, 4, 5}, {12, 2}, {0, 14, 17, 10}, {15, 6, 8}}},
		{cmd: "set", key: 17, val: "17", state: state[int]{{13, 4, 5}, {12, 2}, {17, 0, 14, 10}, {15, 6, 8}}},
		{cmd: "get", key: 16, val: "", state: state[int]{{13, 4, 5}, {12, 2}, {17, 0, 14, 10}, {15, 6, 8}}},
		{cmd: "set", key: 14, val: "14", state: state[int]{{13, 4, 5}, {12, 2}, {14, 17, 0, 10}, {15, 6, 8}}},
		{cmd: "set", key: 12, val: "12", state: state[int]{{13, 4, 5}, {2}, {12, 14, 17, 0, 10}, {15, 6, 8}}},
		{cmd: "set", key: 17, val: "17", state: state[int]{{13, 4, 5}, {2}, {17, 12, 14, 0, 10}, {15, 6, 8}}},
		{cmd: "get", key: 4, val: "", state: state[int]{{13, 4, 5}, {2}, {17, 12, 14, 0, 10}, {15, 6, 8}}},
		{cmd: "get", key: 5, val: "", state: state[int]{{13, 4, 5}, {2}, {17, 12, 14, 0, 10}, {15, 6, 8}}},
		{cmd: "get", key: 14, val: "14", state: state[int]{{13, 4, 5}, {2}, {14, 17, 12, 0, 10}, {15, 6, 8}}},
		{cmd: "get", key: 12, val: "12", state: state[int]{{13, 4, 5}, {2}, {12, 14, 17, 0, 10}, {15, 6, 8}}},
		{cmd: "set", key: 4, val: "4", state: state[int]{{13, 5}, {2}, {4, 12, 14, 17, 0}, {10, 15, 6, 8}}},
		{cmd: "get", key: 17, val: "17", state: state[int]{{13, 5}, {2}, {17, 4, 12, 14, 0}, {10, 15, 6, 8}}},
		{cmd: "set", key: 7, val: "7", state: state[int]{{13, 5}, {2, 7}, {17, 4, 12, 14}, {0, 10, 15, 6}}},
		{cmd: "set", key: 12, val: "12", state: state[int]{{13, 5}, {2, 7}, {12, 17, 4, 14}, {0, 10, 15, 6}}},
		{cmd: "get", key: 9, val: "", state: state[int]{{13, 5}, {2, 7}, {12, 17, 4, 14}, {0, 10, 15, 6}}},
		{cmd: "get", key: 12, val: "12", state: state[int]{{13, 5}, {2, 7}, {12, 17, 4, 14}, {0, 10, 15, 6}}},
		{cmd: "set", key: 0, val: "0", state: state[int]{{13, 5, 2}, {7}, {0, 12, 17, 4, 14}, {10, 15, 6}}},
		{cmd: "get", key: 2, val: "", state: state[int]{{13, 5, 2}, {7}, {0, 12, 17, 4, 14}, {10, 15, 6}}},
		{cmd: "get", key: 11, val: "", state: state[int]{{13, 5, 2}, {7}, {0, 12, 17, 4, 14}, {10, 15, 6}}},
		{cmd: "get", key: 1, val: "", state: state[int]{{13, 5, 2}, {7}, {0, 12, 17, 4, 14}, {10, 15, 6}}},
		{cmd: "get", key: 6, val: "", state: state[int]{{13, 5, 2}, {7}, {0, 12, 17, 4, 14}, {10, 15, 6}}},
		{cmd: "get", key: 10, val: "", state: state[int]{{13, 5, 2}, {7}, {0, 12, 17, 4, 14}, {10, 15, 6}}},
		{cmd: "get", key: 2, val: "", state: state[int]{{13, 5, 2}, {7}, {0, 12, 17, 4, 14}, {10, 15, 6}}},
		{cmd: "set", key: 2, val: "2", state: state[int]{{13, 5}, {7}, {2, 0, 12, 17, 4}, {14, 10, 15, 6}}},
	}
	c := New[int, string](6)
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
		if got := cacheState(c); !reflect.DeepEqual(got, tt.state) {
			var prev state[int]
			if i > 0 {
				prev = tests[i-1].state
			}
			t.Fatalf("step %d: %s: unexpected state:\nprev %s\ngot  %s\nwant %s", i, cmd, prev, got, tt.state)
		}
	}
}
