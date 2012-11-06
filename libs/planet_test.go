package libs

import (
    "testing"
    "strconv"
)

func TestGeneratePlanet(t *testing.T) {
    player := CreatePlayer("gophie", GenerateHash("gophie"), new(Planet))
    sun_position := []int{500, 300}
    expected_planets := []Planet{
        Planet{[]int{375, 247}, 6, 3, 0, 0, nil},
        Planet{[]int{694, 300}, 8, 5, 0, 0, nil},
        Planet{[]int{271, 203}, 3, 1, 0, 0, nil},
        Planet{[]int{209, 365}, 2, 8, 0, 0, nil},
        Planet{[]int{671,  -6}, 3, 1, 0, 0, nil},
        Planet{[]int{907, 300}, 6, 8, 0, 0, nil},
        Planet{[]int{918, 101}, 9, 6, 0, 0, nil},
        Planet{[]int{352, 798}, 5, 4, 0, 0, nil},
        Planet{[]int{686, 841}, 1, 1, 0, 0, nil},
    }
    generated_planets, _ := GeneratePlanets(player.hash, sun_position)

    if len(generated_planets) != 9 {
        t.Error("Wrong planets count")
    }
    for i:=0; i<9; i++ {
        if generated_planets[i].coords[0] != expected_planets[i].coords[0] {
            t.Error("X coordinate missmatch on Planet[" + strconv.Itoa(i) + "]")
        }
        if generated_planets[i].coords[1] != expected_planets[i].coords[1] {
            t.Error("Y coordinate missmatch on Planet[" + strconv.Itoa(i) + "]")
        }
        if generated_planets[i].texture != expected_planets[i].texture {
            t.Error("Ring offset missmatch on Planet[" + strconv.Itoa(i) + "]")
        }
        if generated_planets[i].size != expected_planets[i].size {
            t.Error("Ring offset missmatch on Planet[" + strconv.Itoa(i) + "]")
        }
    }
}
