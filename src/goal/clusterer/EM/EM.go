package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
)

//number to tags
var tags []string = []string{"c2", "c1"}

//data structure to store in lattice
type alphabeta struct {
	alpha float64
	beta  float64
}

var alphaEnd float64
var betaStart float64
var lattice [][]alphabeta

//file data
var fData []byte

//data length
var dataLen int

//map for storing probabilities
var probTTMap = map[string]float64{"c1Start": math.Log(0.5), "c2Start": math.Log(0.5), "c1c1": math.Log(0.6), "c1c2": math.Log(0.9), "c2c1": math.Log(0.38), "c2c2": math.Log(0.08), "Endc1": math.Log(0.02), "Endc2": math.Log(0.02)}
var probWTMap = make(map[string]float64)
var counts = make(map[string]float64)
var lemma = make(map[string]bool)

func main() {
	//Parameters
	filePath := "TRAIN"
	iteratn := 50

	//read data
	var err error
	fData, err = ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	//fill lemma
	for _, v := range fData {
		lemma[string(v)] = true
	}
	//generate P(W|T)
	for k, _ := range lemma {
		for _, tag := range tags {
			probWTMap[string(k)+tag] = math.Log(0.03)
		}
	}

	//initialize counts
	for _, itag := range tags {
		counts[itag+"Start"] = math.Log(0.0)
		counts["End"+itag] = math.Log(0.0)
		for _, v := range fData {
			counts[string(v)+itag] = math.Log(0.0)
			for _, ktag := range tags {
				counts[ktag+itag] = math.Log(0.0)
			}
		}
	}

	tagLen := len(tags)
	dataLen = len(fData)
	lattice = make([][]alphabeta, tagLen)
	for i := range lattice {
		//length of lattice = 2* number of lattice
		lattice[i] = make([]alphabeta, 2*dataLen)
		for j := 0; j < 2*dataLen; j++ {
			lattice[i][j].alpha = math.Log(0.0)
			lattice[i][j].beta = math.Log(0.0)
		}
	}

	for i := 0; i < iteratn; i++ {
		forward()
		backward()
		collectCounts()
		normalizeCounts()
		fmt.Println("================================================")
		if (i+1)%10 == 0 {
			fmt.Println(viterbiTag())

			fmt.Println("============P(W|T)==========")
			for k1, _ := range lemma {
				for _, v := range tags {
					fmt.Println("P(", k1, "|", v, ")=", math.Exp(probWTMap[k1+v]))
				}
			}

			fmt.Println("============P(T|T)==========")
			for _, v1 := range tags {
				for _, v := range tags {
					fmt.Println("P(", v1, "|", v, ")=", math.Exp(probTTMap[v1+v]))
				}
			}

		}

		//clear lattice
		for i := range lattice {
			for j := 0; j < 2*dataLen; j++ {
				lattice[i][j].alpha = math.Log(0.0)
				lattice[i][j].beta = math.Log(0.0)
			}
		}

		//clear counts
		for _, itag := range tags {
			counts[itag+"Start"] = math.Log(0.0)
			counts["End"+itag] = math.Log(0.0)
			for j, v := range fData {
				counts[string(v)+itag] = math.Log(0.0)
				if j < dataLen-1 {
					for _, ktag := range tags {
						counts[ktag+itag] = math.Log(0.0)
					}
				}
			}
		}

	}

}

//addition of log probabilities
func logAdd(x, y float64) float64 {
	m := math.Min(x, y)
	val := 16.0
	if y-x > val {
		return (y)
	} else if y > x {
		return (y + math.Log(1+math.Exp(x-y)))
	} else if x-y > val {
		return (x)
	} else if x > y {
		return (x + math.Log(1+math.Exp(y-x)))
	}
	return (m + math.Log(math.Exp(x-m)+math.Exp(y-m)))
}

//apply forward algorithm on lattice
func forward() {

	//populate first column
	for i, tag := range tags {
		lattice[i][0].alpha = probTTMap[tag+"Start"]
		lattice[i][1].alpha = lattice[i][0].alpha + probWTMap[string(fData[0])+tag]
	}

	//walk the lattice
	for i := 2; i < (2*dataLen)-1; i += 2 {
		for j, jtag := range tags {
			for k, ktag := range tags {
				if math.IsInf(lattice[j][i].alpha, -1) {
					lattice[j][i].alpha = probTTMap[jtag+ktag] + lattice[k][i-1].alpha
				} else {
					lattice[j][i].alpha = logAdd(lattice[j][i].alpha, probTTMap[jtag+ktag]+lattice[k][i-1].alpha)
				}
			}
			lattice[j][i+1].alpha = lattice[j][i].alpha + probWTMap[string(fData[i/2])+jtag]
		}
	}

	//compute alphaEnd
	alphaEnd = math.Log(0.0)
	for i, v := range tags {
		if math.IsInf(alphaEnd, -1) {
			alphaEnd = (lattice[i][(2*dataLen)-1].alpha + probTTMap["End"+v])
		} else {
			alphaEnd = logAdd(alphaEnd, (lattice[i][(2*dataLen)-1].alpha + probTTMap["End"+v]))
		}
	}

	fmt.Println("alphaEnd:", alphaEnd, ":", math.Exp(alphaEnd))
}

//apply backward algorithm on the lattice
func backward() {

	//start at the end
	for i, tag := range tags {
		lattice[i][2*dataLen-1].beta = probTTMap["End"+tag]
		lattice[i][2*dataLen-2].beta = lattice[i][2*dataLen-1].beta + probWTMap[string(fData[dataLen-1])+tag]

	}

	//walk the lattice
	for i := (2 * dataLen) - 3; i >= 0; i -= 2 {
		for j, jtag := range tags {
			for k, ktag := range tags {
				lattice[j][i].beta = logAdd(lattice[j][i].beta, (probTTMap[ktag+jtag] + lattice[k][i+1].beta))
			}
			lattice[j][i-1].beta = lattice[j][i].beta + probWTMap[string(fData[i/2])+jtag]

		}
	}

	//compute betaEnd
	betaStart = math.Log(0.0)
	for i, tag := range tags {
		if math.IsInf(betaStart, -1) {
			betaStart = lattice[i][0].beta + probTTMap[tag+"Start"]
		} else {
			betaStart = logAdd(betaStart, (lattice[i][0].beta + probTTMap[tag+"Start"]))
		}
	}

	fmt.Println("betaStart:", betaStart, ":", math.Exp(betaStart))

}

//collect all probabilities from the data file
func collectCounts() {
	for i, itag := range tags {
		//fmt.Println(probTTMap["End"+itag])
		if math.IsInf(counts[itag+"Start"], -1) {
			counts[itag+"Start"] = (probTTMap[itag+"Start"] + lattice[i][0].beta - alphaEnd)
		} else {
			counts[itag+"Start"] = logAdd(counts[itag+"Start"], (probTTMap[itag+"Start"] + lattice[i][0].beta - alphaEnd))
		}

		if math.IsInf(counts["End"+itag], -1) {
			counts["End"+itag] = (probTTMap["End"+itag] + lattice[i][2*dataLen-1].alpha - alphaEnd)
		} else {
			counts["End"+itag] = logAdd(counts["End"+itag], (probTTMap["End"+itag] + lattice[i][2*dataLen-1].alpha - alphaEnd))
		}

		for j, v := range fData {
			//substitution count
			if math.IsInf(counts[string(v)+itag], -1) {
				counts[string(v)+itag] = (lattice[i][2*j].alpha + probWTMap[string(v)+itag] + lattice[i][j*2+1].beta) - alphaEnd
			} else {
				counts[string(v)+itag] = logAdd(counts[string(v)+itag], (lattice[i][2*j].alpha+probWTMap[string(v)+itag]+lattice[i][j*2+1].beta)-alphaEnd)
			}

			//fmt.Println("--------------------")
			for k, ktag := range tags {
				if j < dataLen-1 {
					if math.IsInf(counts[ktag+itag], -1) {
						counts[ktag+itag] = (lattice[i][(j*2)+1].alpha + probTTMap[ktag+itag] + lattice[k][(j+1)*2].beta) - alphaEnd
					} else {
						counts[ktag+itag] = logAdd(counts[ktag+itag], (lattice[i][(j*2)+1].alpha + probTTMap[ktag+itag] + lattice[k][(j+1)*2].beta - alphaEnd))

						//fmt.Println(ktag+itag,":",counts[ktag+itag],":", math.Exp(counts[ktag+itag]))

					}
				}
			}
			//fmt.Println("--------------------")
		}

	}

}

/*
* M Step: Normalize the fractional counts to get probabilities
 */

func normalizeCounts() {

	var totalTransition = make(map[string]float64)
	var totalSubstitution = make(map[string]float64)
	//initialize counts for each tag
	totalTransition["Start"] = math.Log(0.0)
	for _, itag := range tags {
		totalTransition[itag] = math.Log(0.0)
		totalSubstitution[itag] = math.Log(0.0)
	}

	//get total counts for each tag
	for _, itag := range tags {
		if math.IsInf(totalTransition["Start"], -1) {
			totalTransition["Start"] = counts[itag+"Start"]
		} else {
			totalTransition["Start"] = logAdd(totalTransition["Start"], counts[itag+"Start"])
		}

		if math.IsInf(totalTransition["End"], -1) {
			totalTransition["End"] = counts["End"+itag]
		} else {
			totalTransition["End"] = logAdd(totalTransition["End"], counts["End"+itag])
		}
		for _, jtag := range tags {
			if math.IsInf(totalTransition[itag], -1) {
				totalTransition[itag] = counts[jtag+itag]
			} else {
				totalTransition[itag] = logAdd(totalTransition[itag], counts[jtag+itag])
			}
		}
		//totalTransition[itag] = logAdd(totalTransition[itag],counts["End"+itag])

		for k, _ := range lemma {
			if math.IsInf(totalSubstitution[itag], -1) {
				totalSubstitution[itag] = counts[string(k)+itag]
			} else {
				totalSubstitution[itag] = logAdd(totalSubstitution[itag], counts[string(k)+itag])
			}
		}
	}

	//divide fractional counts by totals
	for _, itag := range tags {
		probTTMap[itag+"Start"] = counts[itag+"Start"] - totalTransition["Start"]
		//	probTTMap["End"+itag] = counts["End"+itag] - totalTransition[itag]
		for _, jtag := range tags {
			probTTMap[jtag+itag] = counts[jtag+itag] - totalTransition[itag]
		}
		for k, _ := range lemma {
			probWTMap[string(k)+itag] = counts[string(k)+itag] - totalSubstitution[itag]
		}
	}

}

//viterbi algorithm for most probable sequence
func viterbiTag() string {
	tagSeq := ""
	for i := 0; i < 100; i += 2 {
		if lattice[0][i].alpha < lattice[1][i].alpha {
			//fmt.Println(lattice[0][i].alpha,":",lattice[1][i].alpha)
			tagSeq += tags[0]
		} else {
			//fmt.Println(lattice[0][i].alpha,":",lattice[1][i].alpha)
			tagSeq += tags[1]
		}
	}
	return tagSeq
}
