package nets

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

type OutputNeurons struct {
	SigVal  float32
	Correct bool
}

type Network struct {
	WS, Bias []float32
	Output   [10][4]OutputNeurons
	Total    float32
}

const (
	tick    = "✅"
	cross   = "❌"
	logFile = "log/network.log"
)

func NewNetwork() *Network {
	// this is the weights and bias after trying to balance
	// the network by hand with the prev. technique
	// weights: [-0.878692 9.761792 0.3720191 12.552786 -25.448135 -17.989666 -4.451405 -26.450663 25.954193 26.252296]
	// total: -0.3254738
	// bias: [-20.432455 -30.319965 -30.324236 -3.7374773]

	net := Network{
		WS:     make([]float32, 10),
		Bias:   make([]float32, 4),
		Output: [10][4]OutputNeurons{},
	}
	// var seedw int64 = 1667659351435904
	// var seedb int64 = 1667659351435945
	// randFlArrn(net.WS, -30, 30, seedw)
	// randFlArrn(net.Bias, -30, 30, seedb)
	net.WS = []float32{-0.878692, 9.761792, 0.3720191, 12.552786, -25.448135, -17.989666, -4.451405, -26.450663, 25.954193, 26.252296}
	net.Bias = []float32{-20.432455, -30.319965, -30.324236, -3.7374773}
	net.calcTotal()

	return &net
}

func (net *Network) ShowChange() string {
	var str strings.Builder
	accuracy := 0
	for i, w := range net.WS {
		str.WriteString(fmt.Sprint(i, ": "))
		for j, b := range net.Bias {
			s := net.sigmoid(w, b)
			correct := false
			if (minOut[i][j] == 0.99 && s > 0.99) || (minOut[i][j] == 0.01 && s <= 0.01) {
				correct = true
				str.WriteString(fmt.Sprintf("%s ", tick))
				accuracy++
			} else {
				str.WriteString(fmt.Sprintf("%s ", cross))
			}
			net.Output[i][j] = OutputNeurons{
				SigVal:  s,
				Correct: correct,
			}
		}
		// fmt.Println()
		str.WriteByte('\n')
	}
	// fmt.Println()
	str.WriteString(fmt.Sprintln("\naccuracy:", accuracy, "/ 40"))
	// fmt.Println()
	// fmt.Printf("Weights: %v\nWeights Total: %v\n", net.WS, net.Total)
	// fmt.Printf("Bias: %v\n", net.Bias)
	// fmt.Println()
	// fmt.Printf("outs: %v\n", weights.Output)

	return str.String()
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

func (net *Network) Change(i int, val float32) {
	// weight.total -= weight.ws[i] + val
	net.WS[i] = val
	net.calcTotal()
}

func (net *Network) Save() error {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(net)
	if err != nil {
		return err
	}
	// fmt.Printf("%s \nLength: %d\n", log, len(log))
	logger := log.New(file, "", log.LstdFlags)
	logger.Println(string(data))

	return nil
}

/*
	file, err := os.Open("log/network.log")

    if err != nil {
        panic("Cannot open file")
        os.Exit(1)
    }
    defer file.Close()

    line := ""
    var cursor int64 = 0
    stat, _ := file.Stat()
    filesize := stat.Size()
    for {
        cursor -= 1
        file.Seek(cursor, io.SeekEnd)

        char := make([]byte, 1)
        file.Read(char)

        if cursor != -1 && (char[0] == 10 || char[0] == 13) { // stop if we find a line
            break
        }

        line = fmt.Sprintf("%s%s", string(char), line) // there is more efficient way

        if cursor == -filesize { // stop if we are at the begining
            break
        }
    }

    return line
*/
