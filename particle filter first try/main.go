package main

import (
	"bohao"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"

	"github.com/go-echarts/go-echarts/charts"
)

var worldSize float64 = 100
var landmarks = [][]float64{[]float64{20.0, 20.0}, []float64{80.0, 80.0}, []float64{20.0, 80.0}, []float64{80.0, 20.0}}

type particle struct {
	x      float64
	y      float64
	orient float64

	forward_noise float64
	turn_noise    float64
	sense_noise   float64
}

func (p *particle) init() {
	p.x = rand.Float64() * worldSize
	p.y = rand.Float64() * worldSize
	p.orient = rand.Float64() * 2.0 * math.Pi

	p.forward_noise = 0.0
	p.turn_noise = 0.0
	p.sense_noise = 0.0
}

func (p *particle) set(new_x, new_y, new_orient float64) {
	p.x = new_x
	p.y = new_y
	p.orient = new_orient
}

func (p *particle) set_noise(new_f_noise, new_t_noise, new_s_noise float64) {
	p.forward_noise = new_f_noise
	p.turn_noise = new_t_noise
	p.sense_noise = new_s_noise
}

func (p *particle) sense() (Z []float64) {
	for i := 0; i < len(landmarks); i++ {
		dist := math.Sqrt(math.Pow(p.x-landmarks[i][0], 2)+math.Pow(p.y-landmarks[i][1], 2)) + (rand.NormFloat64() * p.sense_noise)
		Z = append(Z, dist)
	}
	return
}

func (p *particle) move(turn, forward float64) {
	if forward < 0 {
		log.Fatal("robot can't get backward")
	}

	//turn, and add randomness to the turning command.
	orient := p.orient + turn + rand.NormFloat64()*p.turn_noise
	for {
		if orient >= 2.0*math.Pi {
			orient = orient - 2.0*math.Pi
		} else if orient < 0 {
			orient = orient + 2.0*math.Pi
		} else {
			break
		}
	}

	//move, abd add rabdinbess to the motion command
	dist := forward + rand.NormFloat64()*p.forward_noise
	x := p.x + math.Cos(p.orient)*dist
	y := p.y + math.Sin(p.orient)*dist
	for {
		if x > worldSize {
			x -= worldSize
		} else if x < 0 {
			x += worldSize
		} else {
			break
		}
	}
	for {
		if y > worldSize {
			y -= worldSize
		} else if y < 0 {
			y += worldSize
		} else {
			break
		}
	}

	// set particle
	p.set(x, y, orient)
	p.set_noise(p.forward_noise, p.turn_noise, p.sense_noise)
}

func (p particle) Gaussian(mu, sigma, x float64) float64 {
	return math.Exp(-math.Pow(mu-x, 2)/math.Pow(sigma, 2)/2.0) / math.Sqrt(2.0*math.Pi*math.Pow(sigma, 2))
}

func (p particle) measurement_prob(measurement []float64) float64 {

	prob := 1.0
	for i := 0; i < len(landmarks); i++ {
		dist := math.Sqrt(math.Pow(p.x-landmarks[i][0], 2) + math.Pow(p.y-landmarks[i][1], 2))
		prob *= p.Gaussian(dist, p.sense_noise, measurement[i])

	}

	return prob
}

func single_particle_example() {
	// initial a particle this is an example for single particle.
	p := particle{}
	p.init()
	p.set_noise(5.0, 0.1, 5.0)
	p.set(30, 50, 0.5)
	fmt.Println(p)

	Z := p.sense()
	fmt.Println(Z)

	p.move(math.Pi/2.0, 10.0)
	fmt.Println(p)
	Z = p.sense()
	fmt.Println(Z)
}

func takecoordinate(particle_list []particle) (result [][]float64) {
	for i := 0; i < len(particle_list); i++ {
		result = append(result, []float64{particle_list[i].x, particle_list[i].y})
	}
	return
}

func real_work() {
	page := charts.NewPage()

	//initial N particles
	N := 1000
	var particle_list []particle

	for i := 0; i < N; i++ {
		robot := particle{}
		robot.init()
		robot.set_noise(0.05, 0.05, 5.0)
		particle_list = append(particle_list, robot)
	}

	for time := 0; time < 200; time++ {

		//move all particles
		for i := 0; i < N; i++ {
			particle_list[i].move(0.5, 2)
		}

		//compute the weight of each particle
		var weight_list []float64
		for i := 0; i < N; i++ {
			weight_list = append(weight_list, particle_list[i].measurement_prob(particle_list[i].sense()))
		}

		//resample the particle
		//random initial a list index
		indx := rand.Intn(1) * N
		beta := 0.0
		max_weight := bohao.MaxSlice_float64(weight_list)
		p := []particle{}

		//for N times
		for i := 0; i < N; i++ {

			beta += rand.Float64() * 2.0 * max_weight
			//beta variable increasing
			//if beta is smaller than weight_list[index] then break the cor loop.
			for beta > weight_list[indx] {
				beta -= weight_list[indx]
				indx = (indx + 1) % N
			}

			p = append(p, particle_list[indx])

		}
		particle_list = p

		if time%10 == 0 {
			coordinate := takecoordinate(particle_list)

			scatter := charts.NewScatter()
			scatter.SetGlobalOptions(charts.TitleOpts{Title: strconv.Itoa(time)})

			scatter.AddYAxis(strconv.Itoa(time), coordinate)
			scatter.AddYAxis("landmarks", landmarks)

			page.Add(scatter)
		}
	}

	h, err := os.Create("scatter.html")
	if err != nil {
		panic(err)
	}

	page.Render(h)

}

func main() {
	real_work()
}
