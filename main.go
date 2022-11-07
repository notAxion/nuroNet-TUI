package main

import (
	"fmt"
	"math"
	"math/rand"
)

type outputNeurons struct {
	sigVal  float32
	correct bool
}

type Network struct {
	WS, Bias []float32
	Output   [10][4]outputNeurons
	Total    float32
}

const (
	tick  = "✅"
	cross = "❌"
)

func main() {
	net := NewNetwork()

	net.showChange()
	for {
		choice, index, change := getStuff(net.WS, net.Bias)
		if choice == "w" {
			net.change(index, net.WS[index]+change)
		} else if choice == "b" {
			net.Bias[index] += change
		}

		net.showChange()
		fmt.Println("")
		fmt.Printf("weights: %v\n", net)
		fmt.Printf("weights.Bias: %v\n", net.Bias)
		fmt.Println()

	}
}

func getStuff(ws, bias []float32) (string, int, float32) {
	choice := ""
	index := 0
	fmt.Print("What to change: ")
	fmt.Scanf("%1s%d", &choice, &index)
	if choice == "w" {
		fmt.Println(ws[index])
	} else if choice == "b" {
		fmt.Println(bias[index])
	}

	fmt.Print("\nHow much Change: ")
	var input float32
	fmt.Scan(&input)

	return choice, index, input
}

// Tried different random numbers
// Have found this seed is a good starting point
// weights seed 1667659351435904
// bias seed 1667659351435945
func randFlArrn(input []float32, min, max int, seed int64) {
	// seed := time.Now().UnixMicro()
	r := rand.New(rand.NewSource(seed))
	// fmt.Println("seed", seed)
	if max-min <= 0 {
		panic("min is bigger than max :[ ")
	}
	for i := range input {
		input[i] = float32(r.Intn(max-min)+min) - r.Float32()
	}
}

func NewNetwork() *Network {
	net := Network{
		WS:     make([]float32, 10),
		Bias:   make([]float32, 4),
		Output: [10][4]outputNeurons{},
	}
	var seedw int64 = 1667659351435904
	var seedb int64 = 1667659351435945
	randFlArrn(net.WS, -30, 30, seedw)
	randFlArrn(net.Bias, -30, 30, seedb)
	net.calcTotal()

	return &net
}

func (net *Network) showChange() {
	accuracy := 0
	for i, w := range net.WS {
		fmt.Print(i, ": ")
		for j, b := range net.Bias {
			s := net.sigmoid(w, b)
			correct := false
			if (minOut[i][j] == 0.99 && s > 0.99) || (minOut[i][j] == 0.01 && s <= 0.01) {
				correct = true
				fmt.Printf("%s ", tick)
				accuracy++
			} else {
				fmt.Printf("%s ", cross)
			}
			net.Output[i][j] = outputNeurons{
				sigVal:  s,
				correct: correct,
			}
		}
		fmt.Println()
	}
	fmt.Println()
	fmt.Println("accuracy:", accuracy, "/ 40")
	fmt.Println()
	fmt.Printf("Weights: %v\nWeights Total: %v\n", net.WS, net.Total)
	fmt.Printf("Bias: %v\n", net.Bias)
	fmt.Println()
	// fmt.Printf("outs: %v\n", weights.Output)

}

func (net *Network) sigmoid(w, b float32) float32 {
	var x float32 = 0.99
	v := w*x + 0.01*(net.Total-w) + b
	return sigmoid(v)
}

// func sigmoid(w, x, b float32) float32 {
// v1 := w*x + b
func sigmoid(v1 float32) float32 {
	v := float64(v1)
	s1 := 1 / (1 + math.Exp(-v))
	s1 = math.Round(s1*1e7) / 1e7
	return float32(s1)
}

func (net *Network) calcTotal() {
	net.Total = 0
	for _, w := range net.WS {
		net.Total += w
	}
}

func (net *Network) change(i int, val float32) {
	// weight.total -= weight.ws[i] + val
	net.WS[i] = val
	net.calcTotal()
}
