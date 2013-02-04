package vector

import "math"

type Vector struct {
	X, Y float64
}

func (self *Vector) Init(x, y float64) *Vector {
	self.X = x
	self.Y = y
	return self
}

func New(x, y float64) *Vector {
	return new(Vector).Init(x, y)
}

func (self *Vector) Angle() float64 {
	if self.getLengthSqrd() == 0 {
		return 0
	}

	return math.Atan2(self.Y, self.X) / (math.Pi / 180)
}

func (self *Vector) SetAngle(angle_degrees float64) {
	self.X = self.Length()
	self.Y = 0.0
	self.Rotate(angle_degrees)
}

func (self *Vector) Length() float64 {
	return math.Sqrt(math.Pow(self.X, 2) + math.Pow(self.Y, 2))
}

func (self *Vector) SetLength(value float64) {
	length := self.Length()
	self.X *= value / length
	self.Y *= value / length
}
func (self *Vector) Rotate(angle_degrees float64) {
	radians := (math.Pi / 180) * angle_degrees
	sin := math.Sin(radians)
	cos := math.Cos(radians)

	x := self.X*cos - self.Y*sin
	y := self.X*sin + self.Y*cos
	self.X = x
	self.Y = y
}

func (self *Vector) Rotated(angle_degrees float64) *Vector {
	vector := New(self.X, self.Y)
	vector.Rotate(angle_degrees)
	return vector
}

func (self *Vector) GetAngleBetween(other *Vector) float64 {
	cross := self.X*other.Y - self.Y*other.X
	dot := self.X*other.X + self.Y*other.Y
	return math.Atan2(cross, dot) / (math.Pi / 180)
}

func (self *Vector) GetDistance(other *Vector) float64 {
	return math.Sqrt(math.Pow(self.X-other.X, 2) + math.Pow(self.Y-other.Y, 2))
}

func (self *Vector) Add(other *Vector) *Vector {
	return New(self.X+other.X, self.Y+other.Y)
}

func (self *Vector) AddToFloat64(value float64) *Vector {
	return New(self.X+value, self.Y+value)
}

func (self *Vector) Substitute(other *Vector) *Vector {
	return New(self.X-other.X, self.Y-other.Y)
}

func (self *Vector) SubstituteToFloat64(value float64) *Vector {
	return New(self.X-value, self.Y-value)
}

func (self *Vector) Multiply(other *Vector) *Vector {
	return New(self.X*other.X, self.Y*other.Y)
}

func (self *Vector) MultiplyToFloat64(value float64) *Vector {
	return New(self.X*value, self.Y*value)
}

func (self *Vector) Divide(other *Vector) *Vector {
	return New(self.X/other.X, self.Y/other.Y)
}

func (self *Vector) DivideToFloat64(value float64) *Vector {
	return New(self.X/value, self.Y/value)
}

func (self *Vector) getLengthSqrd() float64 {
	return math.Pow(self.X, 2) + math.Pow(self.Y, 2)
}
