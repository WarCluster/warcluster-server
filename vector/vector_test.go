package vector

import "testing"

func TestNew(t *testing.T) {
	vector := New(2.0, 4.0)
	if vector.X != 2.0 && vector.Y != 4.0 {
		t.Error("Creating new vector has failed")
	}
}

func TestAngle(t *testing.T) {
	vector := New(2.0, 4.0)
	zero_vector := New(0.0, 0.0)

	if vector.Angle() != 63.434948822922 {
		t.Error("vector.Angle() appeared to be", vector.Angle())
	}

	if zero_vector.Angle() != 0.0 {
		t.Error("zero_vector.Angle() appeared to be", vector.Angle())
	}
}

func TestSetAngle(t *testing.T) {
	vector := New(2.0, 4.0)
	vector.SetAngle(45.0)

	if vector.Angle() != 45.0 {
		t.Error("After changing the angle vector.Angle() appeared to be", vector.Angle())
	}
	if vector.X != 3.1622776601683795 && vector.Y != 3.162277660168379 {
		t.Error("After changing the angle coordinates are X:", vector.X, " Y:", vector.Y)
	}
}

func TestLength(t *testing.T) {
	vector := New(2.0, 4.0)

	if vector.Length() != 4.47213595499958 {
		t.Error("vector.Length() appeared to be", vector.Length())
	}
}

func TestSetLength(t *testing.T) {
	vector := New(2.0, 4.0)
	vector.SetLength(5.0)

	if vector.Length() != 5.0 {
		t.Error("Vector's length is not the given one, but", vector.Length())
	}
	if vector.X != 2.23606797749979 && vector.Y != 4.47213595499958 {
		t.Error("After changing the length coordinates are X:", vector.X, " Y:", vector.Y)
	}
}

func TestAdd(t *testing.T) {
	vector := New(2.0, 4.0)
	other := New(3.0, 5.0)
	result := vector.Add(other)

	if result.X != 5.0 && result.Y != 9.0 {
		t.Error("Vector + Other vector gave X:", result.X, " Y:", result.Y)
	}
}

func TestAddToFloat64(t *testing.T) {
	vector := New(2.0, 4.0)
	result := vector.AddToFloat64(3.0)

	if result.X != 5.0 && result.Y != 6.0 {
		t.Error("Vector + float64 gave X:", result.X, " Y:", result.Y)
	}
}

func TestSubstitute(t *testing.T) {
	vector := New(4.0, 4.0)
	other := New(2.0, 3.0)
	result := vector.Substitute(other)

	if result.X != 2.0 && result.Y != 1.0 {
		t.Error("Vector - Other vector gave X:", result.X, " Y:", result.Y)
	}
}

func TestSubstituteToFloat64(t *testing.T) {
	vector := New(2.0, 4.0)
	result := vector.SubstituteToFloat64(1.0)

	if result.X != 1.0 && result.Y != 3.0 {
		t.Error("Vector - float64 gave X:", result.X, " Y:", result.Y)
	}
}

func TestMultiply(t *testing.T) {
	vector := New(2.0, 4.0)
	other := New(3.0, 5.0)
	result := vector.Multiply(other)

	if result.X != 6.0 && result.Y != 20.0 {
		t.Error("Vector * Other gave X:", result.X, " Y:", result.Y)
	}
}

func TestMultiplyToFloat64(t *testing.T) {
	vector := New(2.0, 4.0)
	result := vector.MultiplyToFloat64(3.0)

	if result.X != 6.0 && result.Y != 12.0 {
		t.Error("Vector * Other gave X:", result.X, " Y:", result.Y)
	}
}

func TestDivide(t *testing.T) {
	vector := New(4.0, 8.0)
	other := New(2.0, 4.0)
	result := vector.Divide(other)

	if result.X != 2.0 && result.Y != 2.0 {
		t.Error("Vector * Other gave X:", result.X, " Y:", result.Y)
	}
}

func TestDivideToFloat64(t *testing.T) {
	vector := New(2.0, 4.0)
	result := vector.DivideToFloat64(2.0)

	if result.X != 1.0 && result.Y != 2.0 {
		t.Error("Vector * Other gave X:", result.X, " Y:", result.Y)
	}
}

func TestRotate(t *testing.T) {
	vector := New(2.0, 4.0)
	vector.Rotate(2.0)

	if vector.X != 1.8591836672281876 && vector.Y != 4.067362301481385 {
		t.Error("Vector.Rotate(2.0) gave X:", vector.X, " Y:", vector.Y)
	}
}

func TestRotated(t *testing.T) {
	vector := New(2.0, 4.0)
	result := vector.Rotated(2.0)

	if result.X != 1.8591836672281876 && result.Y != 4.067362301481385 {
		t.Error("Vector.Rotate(2.0) gave X:", result.X, " Y:", result.Y)
	}

	if vector == result {
		t.Error("Vector.Rotated did not return a copy object")
	}
}

func TestGetAngleBetween(t *testing.T) {
	vector := New(2.0, 4.0)
	other := New(5.0, 12.0)
	angle_between := vector.GetAngleBetween(other)

	if angle_between != 3.945186229037563 {
		t.Error("Angle between vector and other happens to be:", angle_between)
	}
}

func TestGetDistance(t *testing.T) {
	vector := New(2.0, 4.0)
	other := New(5.0, 12.0)
	distance := vector.GetDistance(other)

	if distance != 8.54400374531753 {
		t.Error("Distance between vector and other happens to be:", distance)
	}
}
