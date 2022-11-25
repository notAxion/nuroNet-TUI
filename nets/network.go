package nets

import (
	"encoding/json"
	"fmt"
	"io"
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
	tick     = "✅"
	cross    = "❌"
	saveFile = "log/network.log"
	logFile  = "log/error.log"
)

func NewNetwork() *Network {
	net, err := LoadLast()
	if err != nil {
		log.Fatalln(err)
	}

	return net
}

func (net *Network) ShowChange() string {
	var str strings.Builder
	accuracy := 0
	str.WriteByte('\n')
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
		str.WriteByte('\n')
	}
	str.WriteString(fmt.Sprintln("\naccuracy:", accuracy, "/ 40"))

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

func (net *Network) WSChange(i int, val float32) {
	// weight.total -= weight.ws[i] + val
	net.WS[i] = val
	net.calcTotal()
}

func (net *Network) BiasChange(i int, val float32) {
	net.Bias[i] = val
}

func (net *Network) Save() error {
	file, err := os.OpenFile(saveFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(net)
	if err != nil {
		return err
	}
	logger := log.New(file, "", log.LstdFlags)
	logger.Println(string(data))

	return nil
}

func (net *Network) LogError(tuiErr error) {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	logger := log.New(file, "", log.LstdFlags)
	logger.Println(tuiErr.Error())
}

func LoadLast() (*Network, error) {
	file, err := os.Open("log/network.log")

	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cursor int64 = -900
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	filesize := stat.Size()
	for {
		cursor -= 1
		file.Seek(cursor, io.SeekEnd)

		char := make([]byte, 1)
		file.Read(char)

		if cursor != -1 && (char[0] == 10 || char[0] == 13) { // stop if we find a line
			break
		}

		if cursor == -filesize { // stop if we are at the begining
			break
		}
	}

	line, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	// skipping the log data
	line = line[19:]

	net := Network{}
	err = json.Unmarshal(line, &net)
	if err != nil {
		return nil, err
	}

	return &net, err
}
