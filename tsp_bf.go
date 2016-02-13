package main

import (
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"

	"image"
	"image/color"

	"math"
	"fmt"

	"strconv"
	"strings"
	"bufio"
	"os"
	"time"
)


type Point struct {
	x	float64
	y	float64
	distances []float64
}


func Factorial(n int) int {
    if n == 0 {
        return 1
    }
    return n * Factorial(n - 1)
}


func calculate_score(way *[]int, cities *[]Point) float64 {
	score := 0.0
	for idx, k := range *way {
		var j int
		if idx+1 < len(*way) {
			j = (*way)[idx+1]
		} else {
			j = (*way)[0]
		}
		score += (*cities)[k].distances[j]
	}
	return score
}


func calc_distances(cities * []Point) {
	for i, v := range *cities {
		(*cities)[i].distances = make([]float64, len(*cities))
		for k, v2 := range *cities {
			dx, dy := v.x - v2.x, v.y - v2.y
			(*cities)[i].distances[k] = math.Sqrt(dx * dx + dy * dy)
		}
	}
}


func load_points() []Point {
	file, err := os.Open("./test.txt")
	if err != nil {
		fmt.Fprint(os.Stderr, "Cannot open points file")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	points := make([]Point, 0)
	for scanner.Scan() {
		p_data := strings.Split(scanner.Text(), " ")
		x, _ := strconv.ParseFloat(p_data[1], 64)
		y, _ := strconv.ParseFloat(p_data[2], 64)
		points = append(points, Point{
			x,
			y,
			nil,
		})
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprint(os.Stderr, "Points file reading error!")
	}
	return points
}


func draw_way(way []int, points []Point) {

	biggest_x, biggest_y := points[0].x, points[0].y
	for _, p := range points[1:] {
		if biggest_x < p.x {
			biggest_x = p.x
		}
		if biggest_y < p.y {
			biggest_y = p.y
		}
	}
	padding := 40.0
	resize_by := 10.0

	dest := image.NewRGBA(image.Rect(0, 0, int(biggest_x * resize_by + padding), int(biggest_y * resize_by + 2*padding)))
	gc := draw2dimg.NewGraphicContext(dest)


	gc.SetStrokeColor(color.RGBA{0x44, 0x44, 0x44, 0xff})
	gc.SetLineWidth(1)


	draw2dkit.Circle(gc, padding + points[way[0]].x * resize_by, padding + points[way[0]].y * resize_by, 3)
	gc.MoveTo(padding + points[way[0]].x * resize_by, padding + points[way[0]].y * resize_by)
	for _, v := range way[1:] {
		gc.LineTo(padding + points[v].x * resize_by, padding + points[v].y * resize_by)
		gc.MoveTo(padding + points[v].x * resize_by, padding + points[v].y * resize_by)
		gc.Close()
		draw2dkit.Circle(gc, padding + points[v].x * resize_by, padding + points[v].y * resize_by, 3)
	}
	gc.LineTo(padding + points[way[0]].x * resize_by , padding + points[way[0]].y * resize_by )
	gc.FillStroke()
	draw2dimg.SaveToPngFile("way.png", dest)
}


func permutations(p *[]int, c chan []int){
	n := len(*p)
	indices := make([]int, n)
	for i := 0; i < n; i++ {
		indices[i] = i
	}

	cycles := make([]int, n)
	for i := 0; i < n; i++ {
		cycles[i] = n - i
	}

	c <- *p
	for n > 0 {
		stop_loop := true
		for i := n-1; i>=0; i-- {
			stop_loop = true
			cycles[i]--
			if cycles[i] == 0 {
				x := make([]int, 0)
				x = append(x, indices[i+1:]...)
				x = append(x, indices[i:i+1]...)
				backup := indices[:i]
				indices = append(backup, x...)
				cycles[i] = n - i
			} else {
				j := cycles[i]
				indices[i], indices[n-j] = indices[n-j], indices[i]
				c <- indices
				stop_loop = false
				break
			}
		}
		if stop_loop {
			break
		}
	}
	close(c)
}


func main() {

	cities := load_points()
	calc_distances(&cities)
	points := make([]int, 0)

	for i, _ := range cities {
		points = append(points, i)
	}

	fmt.Printf("Ilość miast: %d \n", len(cities))
	fmt.Printf("Ilość możliwych tras: %d \n", Factorial(len(cities)))

	ch := make(chan []int)
	go permutations(&points, ch)
	best_score := -1.0
	best_way := make([]int, len(cities))
	score := -1.0
	idx := 0
	start := time.Now()
	for i := range ch {

		if idx % 5000000 == 0 && idx != 0 {
			fmt.Printf("%d iteracja\n", idx)
		}
		idx++

		score = calculate_score(&i, &cities)

		if best_score < 0 || best_score > score {
			fmt.Printf("Nowa najlepsza trasa! Wynik: %f \n", score)
			best_score = score
			copy(best_way, i)
		}

	}
	elapsed := time.Since(start)
 	fmt.Printf("Czas wyknonania: %s", elapsed)

	draw_way(best_way, 	cities)
}
