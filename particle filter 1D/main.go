package main

import (
	"bohao"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/go-echarts/go-echarts/charts"
)

var worldSize float64 = 200

type particle struct {
	x float64

	move_noise  float64
	sense_noise float64
}

func (p *particle) init() {
	p.x = rand.Float64() * worldSize

	p.move_noise = 0.0
	p.sense_noise = 0.0
}

func (p *particle) set(new_x float64) {
	p.x = new_x
}

func (p *particle) set_noise(new_f_noise, new_s_noise float64) {
	p.move_noise = new_f_noise
	p.sense_noise = new_s_noise
}

func (p *particle) sense(target float64) (Z float64) {
	Z = math.Abs(p.x - target + (rand.NormFloat64() * p.sense_noise))
	return
}

func (p *particle) move(forward float64) {
	//turn, and add randomness to the turning command.
	turn := math.Floor(rand.Float64()+0.5)*2.0 - 1.0

	//move, abd add rabdinbess to the motion command
	dist := forward + rand.NormFloat64()*p.move_noise
	x := p.x + dist*turn

	// don't let it across the limit
	for {
		if x > worldSize {
			x = worldSize - (x - worldSize)
		} else if x < 0 {
			x = -x
		} else {
			break
		}
	}

	// set particle
	p.set(x)
	p.set_noise(p.move_noise, p.sense_noise)
}

func (p particle) Gaussian(sigma, x float64) float64 {
	//return math.Exp(-math.Pow(mu-x, 2)/math.Pow(sigma, 2)/2.0) / math.Sqrt(2.0*math.Pi*math.Pow(sigma, 2))
	return 1 / sigma / math.Sqrt(2.0*math.Pi) * math.Exp(-math.Pow(x, 2)/2/math.Pow(sigma, 2))
}

func (p particle) measurement_prob(target float64) float64 {

	prob := 1.0
	prob *= p.Gaussian(p.sense_noise, p.sense(target))
	return prob
}

func takecoordinate(particle_list []particle) (result []float64) {
	for i := 0; i < len(particle_list); i++ {
		result = append(result, particle_list[i].x)
	}
	return
}

func where_center(coordinate []float64) (center float64) {
	center = bohao.Sum_float(coordinate)
	center /= float64(len(coordinate))
	return
}

func target_move(t int) (value float64) {
	value = 0.8*float64(t) + 40 + math.Sin(float64(t)*math.Pi/5.0)
	return
}

func add_state(coordinate []float64, state float64) (result [][]float64) {
	for i := 0; i < len(coordinate); i++ {
		result = append(result, []float64{state, coordinate[i]})
	}
	return
}

func real_work() {
	rand.Seed(time.Now().Unix())
	fmt.Println("rand seed success!")

	// number of particle
	N := 100
	times := 2

	h, _ := os.Create("scatter.html")
	scatter := charts.NewScatter()

	for state := 1; state < 20; state++ {
		// find target
		target := target_move(state)

		// initial N particles
		var particle_list []particle
		for i := 0; i < N; i++ {
			robot := particle{}
			robot.init()
			robot.set_noise(0.5, 2.5)
			particle_list = append(particle_list, robot)
		}

		for time := 0; time < times; time++ {
			// move all particles
			for i := 0; i < N; i++ {
				particle_list[i].move(0.5)
			}

			//compute the weight of each particle
			var weight_list []float64
			for i := 0; i < N; i++ {
				weight_list = append(weight_list, particle_list[i].measurement_prob(target))
			}

			// normalized the weights
			total_weight := bohao.Sum_float(weight_list)
			for i := 0; i < len(weight_list); i++ {
				weight_list[i] /= total_weight
			}

			//resample the particle
			
			indx := rand.Intn(1) * N //random initial a list index
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

		} // time
		fmt.Println(where_center(takecoordinate(particle_list)), target)
		scatter.AddYAxis("particle", add_state(takecoordinate(particle_list), float64(state)))
		scatter.AddYAxis("target", [][]float64{{float64(state), target_move(state)}})

	} // state
	scatter.Render(h)
	fmt.Println("program done!")
}

func maybe_better() {
	rand.Seed(time.Now().Unix())
	fmt.Println("rand seed success!")

	// number of particle
	N := 100
	times := 2

	h, _ := os.Create("scatter.html")
	scatter := charts.NewScatter()

	for state := 1; state < 20; state++ {
		// find target
		target := target_move(state)

		// initial N particles
		var particle_list []particle
		for i := 0; i < N; i++ {
			robot := particle{}
			robot.init()
			robot.set_noise(0.5, 2.5)
			particle_list = append(particle_list, robot)
		}

		for time := 0; time < times; time++ {
			// move all particles
			for i := 0; i < N; i++ {
				particle_list[i].move(0.5)
			}

			//compute the weight of each particle
			var weight_list []float64
			for i := 0; i < N; i++ {
				weight_list = append(weight_list, particle_list[i].measurement_prob(target))
			}

			// normalized the weights
			total_weight := bohao.Sum_float(weight_list)
			for i := 0; i < len(weight_list); i++ {
				weight_list[i] /= total_weight
			}

			//resample the particle
			
			indx := rand.Intn(1) * N //random initial a list index
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

		} // time
		fmt.Println(where_center(takecoordinate(particle_list)), target)
		scatter.AddYAxis("particle", add_state(takecoordinate(particle_list), float64(state)))
		scatter.AddYAxis("target", [][]float64{{float64(state), target_move(state)}})

	} // state
	scatter.Render(h)
	fmt.Println("program done!")
}

func main() {
	maybe_better()

}
