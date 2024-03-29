package strMetInt

import "testing"

func TestPerimeter(t *testing.T){
	rectangle := Rectangle{10.0,10.0}
	got := Perimeter(rectangle)
	want := 40.0
	if got != want {
		t.Errorf("got %.2f want %.2f", got, want)
	}
}

//func TestArea(t *testing.T){
//
//	checkArea := func(t testing.TB, shape Shape, want float64){
//		t.Helper()
//		got := shape.Area()
//		if got != want {
//			t.Errorf("got %.2f want %.2f", got, want)
//		}
//	}
//
//	t.Run("rectangle", func(t *testing.T){
//		rectangle := Rectangle{10.0,3}
//		checkArea(t, rectangle, 30.0)
//	})
//	t.Run("circle", func(t *testing.T){
//		circle := Circle{10.0}
//		checkArea(t, circle, 314.1592653589793)
//	})
//}

func TestArea(t *testing.T){
	areaTests := []struct{
		name string
		shape Shape
		want float64
	}{
		{name: "Rectangle", shape: Rectangle{Width: 12, Height: 6}, want: 72.0},
		{name: "Circle", shape: Circle{Radius: 10}, want: 314.1592653589793},
		{name: "Triangle", shape: Triangle{Base: 12, Height: 6}, want: 36.0},
	}
	for _, test := range areaTests {
		t.Run(test.name, func(t *testing.T){
			got := test.shape.Area()
			want := test.want
			if got != want {
				t.Errorf("%#v got %g want %g",test.shape ,got, want)
			}
		})
	}
}
