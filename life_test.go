/**
 * Created by IntelliJ IDEA.
 * User: r2p2
 * Date: 12/29/11
 * Time: 7:54 PM
 */
package life

import (
    "testing"
)

var testcases [][3]uint32

func init() {
    testcases = [][3]uint32{
        [3]uint32{0,0,0},
        [3]uint32{2,1,7},
        [3]uint32{3,2,13},
    }
}

func TestToArea(t *testing.T){
    field := NewField(5, 4)
    for _, testcase := range testcases {
        x := testcase[0]
        y := testcase[1]
        expected := testcase[2]

        result := field.toArea(x, y)
        if result != expected {
            t.Fatalf("x: %d y: %d should be %d and not %d\n", x, y, expected, result)
        }
    }
}

func TestToReal(t *testing.T){
    field := NewField(5, 4)
    for _, testcase := range testcases {
        expected_x := testcase[0]
        expected_y := testcase[1]
        index := testcase[2]

        x, y := field.toReal(index)
        if x != expected_x ||  y != expected_y {
            t.Fatalf("with index: %d I got x: %d y: %d but should be x: %d y: %d\n", index, x, y, expected_x, expected_y)
        }
    }
}
