package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
)

var wtPairFreqMap = make(map[string]float64)
var ttPairFreqMap = make(map[string]float64)
var wtPairProbMap = make(map[string]float64)
var ttPairProbMap = make(map[string]float64)
var wordTypesFreqMap = make(map[string]float64)
var tagTypesFreqMap = make(map[string]float64)
var capWTPairProbMap = make(map[string]float64)
var capWTPairFreqMap = make(map[string]float64)
var capTypesFreqMap = make(map[string]float64)
var startPairList list.List
var endPairList list.List
var wordList list.List
var wList []string
var tList []string

type cell struct {
	value   float64
	currTag string
	prevTag string
	word    string
}

var lattice [][]cell

func main() {

	/*
	* Question 1.
	* Read file and create P(W|T) and P(T|T)
	 */
	//read file
	fStart, err := os.Open("E:\\projects\\csci562\\hw4\\hw4\\train.tags")
	if err != nil {
		log.Fatal(err)
	}

	//use buffer
	bfStart := bufio.NewReader(fStart)

	//read file line by line
	for {
		bLine, _, err := bfStart.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		//fill probability matrices
		FillPairFreqValues(string(bLine))
	}
	//close the file to use later
	fStart.Close()

	for k, _ := range wtPairFreqMap {
		pairSplit := strings.Split(k, "|")
		if len(pairSplit) == 2 {
			wordTypesFreqMap[pairSplit[0]] = wordTypesFreqMap[pairSplit[0]] + 1.0
		}
	}

	for k, _ := range ttPairFreqMap {
		pairSplit := strings.Split(k, "|")
		if len(pairSplit) == 2 {
			tagTypesFreqMap[pairSplit[1]] = tagTypesFreqMap[pairSplit[1]] + 1.0
		}
	}
	
	for k, _ := range capWTPairFreqMap {
		pairSplit := strings.Split(k, "|")
		if len(pairSplit) == 2 {
			capTypesFreqMap[pairSplit[1]] = capTypesFreqMap[pairSplit[1]] + 1.0
		}
	}

	fmt.Println("Number of character types:", len(wordTypesFreqMap))
	fmt.Println("Number of tag types:", len(tagTypesFreqMap))
	FindMaxTagged(wordTypesFreqMap, tagTypesFreqMap)
	tList = make([]string, len(tagTypesFreqMap)-1)
	i := 0
	for k, _ := range tagTypesFreqMap {
		if !strings.EqualFold(k, "start") {
			tList[i] = k
			i++
		}
	}

	//turn counts into log probabilities
	for k, v := range wtPairFreqMap {
		split := strings.Split(k, "|")
		wtPairProbMap[k] = math.Log(v / tagTypesFreqMap[split[1]])
	}

	for k, v := range ttPairFreqMap {
		split := strings.Split(k, "|")
		ttPairProbMap[k] = math.Log(v / tagTypesFreqMap[split[1]])
	}
	
	for x, v := range capWTPairFreqMap {
		split := strings.Split(x, "|")
		capWTPairProbMap[x] = math.Log(v / tagTypesFreqMap[split[1]])
	}

	for i := 0; i < 25; i++ {
		trainResult, trainGoldStd:= RunViterbi("E:\\projects\\csci562\\hw4\\hw4\\train.tags", "E:\\projects\\csci562\\hw4\\hw4\\trainResult.tags")
		fmt.Println("--------TRAIN--------")
		TestAccuracy(trainResult, trainGoldStd)
		
		devResult, devGoldStd := RunViterbi("E:\\projects\\csci562\\hw4\\hw4\\dev.tags", "E:\\projects\\csci562\\hw4\\hw4\\devResult.tags")
		fmt.Println("--------DEV--------")
		TestAccuracy(devResult, devGoldStd)
		UpdateLambda(devResult, devGoldStd)
	}
}

func UpdateLambda(rData, gData string) {
	jumpRate:= 0.5

	//generate W|T and T|T freq table for result
	rWTPairFreqMap := make(map[string]float64)
	rTTPairFreqMap := make(map[string]float64)
	rcapWTPairFreqMap := make(map[string]float64)
	rTagTypesFreqMap := make(map[string]float64)
	rWordFreqMap := make(map[string]float64)
	rcapFreqMap := make(map[string]float64)
	var rWordTagList list.List

	rSent := strings.Split(rData, "\n")
	rSent = rSent[:len(rSent)-1]
	for _, sent := range rSent {
		rTok := strings.Split(sent, " ")
		//generate word related frequencies
		for _, pair := range rTok {
			pairTok := strings.Split(pair, "_")
			if len(pairTok) == 2 {
				rWordTagList.PushBack(pair)
				rWTPairFreqMap[pairTok[0]+"|"+pairTok[1]] = rWTPairFreqMap[pairTok[0]+"|"+pairTok[1]] + 1.0
				rWordFreqMap[pairTok[0]] = rWordFreqMap[pairTok[0]] + 1.0
				if byte(pairTok[0][0])>64 && byte(pairTok[0][0])<91{
					rcapWTPairFreqMap[pairTok[0]+"|"+pairTok[1]] = rcapWTPairFreqMap[pairTok[0]+"|"+pairTok[1]] + 1.0
					rcapFreqMap[pairTok[0]] = rcapFreqMap[pairTok[0]] + 1.0
				}
			}
		}

		//generate tag related frequencies
		for i := 0; i < len(rTok)-1; i++ {
			firstPairSplit := strings.Split(rTok[i], "_")
			secondPairSplit := strings.Split(rTok[i+1], "_")
			if len(firstPairSplit) == 2 && len(secondPairSplit) == 2 {
				if i == 0 {
					rTTPairFreqMap[firstPairSplit[1]+"|start"] = rTTPairFreqMap[firstPairSplit[1]+"|start"] + 1.0
				}
				if i == len(rTok)-2 {
					rTTPairFreqMap["end|"+secondPairSplit[1]] = rTTPairFreqMap["end|"+secondPairSplit[1]] + 1.0
				}
				rTTPairFreqMap[secondPairSplit[1]+"|"+firstPairSplit[1]] = rTTPairFreqMap[secondPairSplit[1]+"|"+firstPairSplit[1]] + 1.0
			}
		}
		for k, _ := range rTTPairFreqMap {
			pairSplit := strings.Split(k, "|")
			rTagTypesFreqMap[pairSplit[1]] = rTagTypesFreqMap[pairSplit[1]] + 1.0
		}

	}

	//generate W|T and T|T freq table for gold
	gWTPairFreqMap := make(map[string]float64)
	gTTPairFreqMap := make(map[string]float64)
	gTagTypesFreqMap := make(map[string]float64)
	gWordFreqMap := make(map[string]float64)
	gcapWTPairFreqMap := make(map[string]float64)
	gcapFreqMap := make(map[string]float64)
	var gWordTagList list.List

	gSent := strings.Split(gData, "\n")
	gSent = gSent[:len(gSent)-1]
	for _, sent := range gSent {
		gTok := strings.Split(sent, " ")
		//generate word related frequencies
		for _, pair := range gTok {
			pairTok := strings.Split(pair, "_")
			if len(pairTok) == 2 {
				gWordTagList.PushBack(pair)
				gWTPairFreqMap[pairTok[0]+"|"+pairTok[1]] = gWTPairFreqMap[pairTok[0]+"|"+pairTok[1]] + 1.0
				gWordFreqMap[pairTok[0]] = gWordFreqMap[pairTok[0]] + 1.0
				if byte(pairTok[0][0])>64 && byte(pairTok[0][0])<91{
					gcapWTPairFreqMap[pairTok[0]+"|"+pairTok[1]] = gcapWTPairFreqMap[pairTok[0]+"|"+pairTok[1]] + 1.0
					gcapFreqMap[pairTok[0]] = gcapFreqMap[pairTok[0]] + 1.0
				}
			}
		}

		//generate tag related frequencies
		for i := 0; i < len(gTok)-1; i++ {
			firstPairSplit := strings.Split(gTok[i], "_")
			secondPairSplit := strings.Split(gTok[i+1], "_")
			if len(firstPairSplit) == 2 && len(secondPairSplit) == 2 {
				if i == 0 {
					gTTPairFreqMap[firstPairSplit[1]+"|start"] = gTTPairFreqMap[firstPairSplit[1]+"|start"] + 1.0
				}
				if i == len(gTok)-2 {
					gTTPairFreqMap["end|"+secondPairSplit[1]] = gTTPairFreqMap["end|"+secondPairSplit[1]] + 1.0
				}
				gTTPairFreqMap[secondPairSplit[1]+"|"+firstPairSplit[1]] = gTTPairFreqMap[secondPairSplit[1]+"|"+firstPairSplit[1]] + 1.0
			}
		}
		for k, _ := range gTTPairFreqMap {
			pairSplit := strings.Split(k, "|")
			gTagTypesFreqMap[pairSplit[1]] = gTagTypesFreqMap[pairSplit[1]] + 1.0
		}

	}
	
	//calculate worse tag in the column for a feature capital(W)|T
	for e:= gWordTagList.Front();e!=nil;e=e.Next(){
		gPairSplit := strings.Split(e.Value.(string), "_")
		if len(gPairSplit) == 2{
			sGold:= math.Log(gcapWTPairFreqMap[gPairSplit[0]+"|"+gPairSplit[1]]/gTagTypesFreqMap[gPairSplit[1]])*gcapWTPairFreqMap[gPairSplit[0]+"|"+gPairSplit[1]]
			var worstTag string
			worstTagVal:=math.Log(0)
			for t,_:=range rTagTypesFreqMap{
				sRslt:= math.Log(rcapWTPairFreqMap[gPairSplit[0]+"|"+t]/rTagTypesFreqMap[t])*rcapWTPairFreqMap[gPairSplit[0]+"|"+t]
				tVal:= Loss(gPairSplit[1], t)- (sGold-sRslt)
				if worstTagVal<tVal{
					worstTagVal = tVal
					worstTag = t
				}
			}
			capWTPairFreqMap[gPairSplit[0]+"|"+ worstTag] = capWTPairFreqMap[gPairSplit[0]+"|"+ worstTag] + jumpRate*(gcapWTPairFreqMap[gPairSplit[0]+"|"+gPairSplit[1]] - rcapWTPairFreqMap[gPairSplit[0]+"|"+worstTag])
		}
	}
	
	//calculate worse tag in the column for a feature W|T
	for e:= gWordTagList.Front();e!=nil;e=e.Next(){
		gPairSplit := strings.Split(e.Value.(string), "_")
		if len(gPairSplit) == 2{
			sGold:= math.Log(gWTPairFreqMap[gPairSplit[0]+"|"+gPairSplit[1]]/gTagTypesFreqMap[gPairSplit[1]])*gWTPairFreqMap[gPairSplit[0]+"|"+gPairSplit[1]]
			var worstTag string
			worstTagVal:=math.Log(0)
			for t,_:=range rTagTypesFreqMap{
				sRslt:= math.Log(rWTPairFreqMap[gPairSplit[0]+"|"+t]/rTagTypesFreqMap[t])*rWTPairFreqMap[gPairSplit[0]+"|"+t]
				tVal:= Loss(gPairSplit[1], t)- (sGold-sRslt)
				if worstTagVal<tVal{
					worstTagVal = tVal
					worstTag = t
				}
			}
			wtPairFreqMap[gPairSplit[0]+"|"+ worstTag] = wtPairFreqMap[gPairSplit[0]+"|"+ worstTag] + jumpRate*(gWTPairFreqMap[gPairSplit[0]+"|"+gPairSplit[1]] - rWTPairFreqMap[gPairSplit[0]+"|"+worstTag])
		}
	}
	
	//calculate worse tag in the column for a feature T|T
	prevT:="start"
	firstCol:=true
	for e:= gWordTagList.Front();e!=nil;e=e.Next(){
		gPairSplit := strings.Split(e.Value.(string), "_")
		if len(gPairSplit) == 2{
			var sGold float64
			if firstCol{
				sGold= math.Log(gTTPairFreqMap[gPairSplit[1]+"|start"]/gTagTypesFreqMap["start"])*gTTPairFreqMap[gPairSplit[1]+"|start"]
				firstCol = false
			}else{
				sGold= math.Log(gTTPairFreqMap[gPairSplit[1]+"|"+prevT]/gTagTypesFreqMap[gPairSplit[1]])*gTTPairFreqMap[gPairSplit[1]+"|"+prevT]
			}
			var worstTag string
			worstTagVal:=math.Log(0)
			for t,_:=range rTagTypesFreqMap{
				sRslt:= math.Log(rTTPairFreqMap[gPairSplit[1]+"|"+t]/rTagTypesFreqMap[t])*rTTPairFreqMap[gPairSplit[1]+"|"+t]
				tVal:= Loss(gPairSplit[1], t)- (sGold-sRslt)
				if worstTagVal<tVal{
					worstTagVal = tVal
					worstTag = t
				}
			}
			
			ttPairFreqMap[gPairSplit[1]+"|"+ worstTag] = ttPairFreqMap[gPairSplit[1]+"|"+ worstTag] + jumpRate*(gTTPairFreqMap[gPairSplit[1]+"|"+prevT] - rTTPairFreqMap[gPairSplit[1]+"|"+worstTag])
		}
		
		prevT = gPairSplit[1]
	}

}

func Loss(gTag, rTag string) float64 {
	l := 1.0
	if strings.EqualFold(gTag, rTag) {
		l = 0
	}
	return l
}

func TestAccuracy(rslt, gold string) (rsltTokTemp, goldTokTemp list.List) {
	var m, n float64
	rSent := strings.Split(rslt, "\n")
	for _, sent := range rSent {
		rsltTokCand := strings.Split(sent, " ")
		for i := 0; i < len(rsltTokCand); i++ {
			if len(strings.TrimSpace(rsltTokCand[i])) > 0 {
				rsltTokTemp.PushBack(strings.TrimSpace(rsltTokCand[i]))
			}
		}
	}

	gSent := strings.Split(gold, "\n")
	for _, sent := range gSent {
		goldTokCand := strings.Split(sent, " ")
		//fmt.Println(goldTokCand)
		for i := 0; i < len(goldTokCand); i++ {
			if len(strings.TrimSpace(goldTokCand[i])) > 0 {
				goldTokTemp.PushBack(strings.TrimSpace(goldTokCand[i]))
			}
		}
	}

	if goldTokTemp.Len() != rsltTokTemp.Len() {
		fmt.Println("Accuracy: Length mis-match:", goldTokTemp.Len(), rsltTokTemp.Len())
	} else {
		r := rsltTokTemp.Front()
		g := goldTokTemp.Front()
		gLen := goldTokTemp.Len()
		for i := 0; i < gLen; i++ {
			rTok := strings.Split(r.Value.(string), "_")
			gTok := strings.Split(g.Value.(string), "_")

			if !strings.EqualFold(rTok[0], gTok[0]) {
				fmt.Println("Accuracy: word mis-match - ", rTok[0], gTok[0])
				fmt.Println("Accuracy: word mis-match Index- ", i)
				os.Exit(1)
			} else {
				n++
				if strings.EqualFold(rTok[1], gTok[1]) {
					m++
				}

			}
			r = r.Next()
			g = g.Next()
		}

	}

	fmt.Println("total words:  ", n)
	fmt.Println("correct tags: ", m)
	fmt.Println("accuracy:     ", float64(m/n))
	return
}

func RunViterbi(inFile, outFile string) (result string, gold string) {

	//read file
	f1, err := os.Open(inFile)
	defer f1.Close()
	if err != nil {
		log.Fatal(err)
	}
	//use buffer
	bf1 := bufio.NewReader(f1)

	//create file to write
	fOut, err := os.Create(outFile)
	defer fOut.Close()
	if err != nil {
		log.Fatal(err)
	}
	//read file line by line
	for {
		bLine, _, err := bf1.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		//run VITERBI
		gold = gold + strings.TrimSpace(string(bLine)) + "\n"
		result = result + ViterbiTagged(string(bLine)) + "\n"
	}

	if _, err := fOut.WriteString(result); err != nil {
		log.Fatal(err)
	}
	return
}
func FindMaxTagged(w, t map[string]float64) {
	var tagPerWordFreqMap = make(map[string]float64)
	for kw, _ := range w {
		for kt, _ := range t {
			value := wtPairFreqMap[kw+"|"+kt]
			if value != 0 {
				tagPerWordFreqMap[kw] = tagPerWordFreqMap[kw] + 1.0
			}
		}
	}

	var max float64
	for _, v := range tagPerWordFreqMap {
		if max < v {
			max = v
		}
	}
	fmt.Println("max # of tags:", max)
	var maxTaggedWords list.List
	for k, v := range tagPerWordFreqMap {
		if v == max {
			maxTaggedWords.PushBack(k)
		}
	}

	//print words with most tags
	fmt.Println("words with max # of tags:")
	for e := maxTaggedWords.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value.(string))
	}

}

func ViterbiTagged(l string) string {
	var wList list.List
	strTokCand := strings.Split(l, " ")

	var strTokTemp list.List
	for i := 0; i < len(strTokCand); i++ {
		if len(strings.TrimSpace(strTokCand[i])) > 0 {
			strTokTemp.PushBack(strTokCand[i])
		}
	}

	j := 0
	var strTok = make([]string, strTokTemp.Len())
	for e := strTokTemp.Front(); e != nil; e = e.Next() {
		strTok[j] = e.Value.(string)
		j++
	}

	//create W frequency table
	for _, v := range strTok {
		pairSplit := strings.Split(v, "_")
		wList.PushBack(pairSplit[0])
	}

	tagLen := len(tList)
	dataLen := wList.Len()
	lattice = make([][]cell, tagLen)
	for i := range lattice {
		lattice[i] = make([]cell, dataLen)
		for j := 0; j < dataLen; j++ {
			lattice[i][j].value = math.Log(0.0)
		}
	}

	i := 0
	for e := wList.Front(); e != nil; e = e.Next() {
		w := e.Value.(string)
		for j, t := range tList {
			lattice[j][i].word = e.Value.(string)
			lattice[j][i].currTag = t
			if i == 0 {
				ttPair, okTT := ttPairProbMap[t+"|start"]
				ttPairC := ttPairFreqMap[t+"|start"]
				wtPair, okWT := wtPairProbMap[w+"|"+t]
				wtPairC := wtPairFreqMap[w+"|"+t]
				capWTPair, _ := capWTPairProbMap[w+"|"+t]
				capPairC:= capWTPairFreqMap[w+"|"+t]
				if okTT && okWT {
					lattice[j][i].value = ttPairC*ttPair + wtPairC*wtPair + capPairC*capWTPair
				} else if !okTT && okWT {
					lattice[j][i].value = float64(len(ttPairFreqMap))*math.Log(float64(1)/float64(len(ttPairFreqMap))) + wtPairC*wtPair + capPairC*capWTPair
				} else if okTT && !okWT {
					lattice[j][i].value = ttPairC*ttPair + float64(len(wtPairFreqMap))*math.Log(float64(1)/float64(len(wtPairFreqMap))) + float64(len(capWTPairFreqMap))*math.Log(float64(1)/float64(len(capWTPairFreqMap)))
				} else if !okTT && !okWT {
					lattice[j][i].value = float64(len(ttPairFreqMap))*math.Log(float64(1)/float64(len(ttPairFreqMap))) + float64(len(wtPairFreqMap))*math.Log(float64(1)/float64(len(wtPairFreqMap))) + float64(len(capWTPairFreqMap))*math.Log(float64(1)/float64(len(capWTPairFreqMap)))
				}
				lattice[j][i].prevTag = "start"

			} else {
				maxVal := math.Log(0)
				maxValTag := tList[0]

				for k, kTag := range tList {
					ttPair, okTT := ttPairProbMap[t+"|"+kTag]
					ttPairC := ttPairFreqMap[t+"|"+kTag]
					wtPair, okWT := wtPairProbMap[w+"|"+t]
					wtPairC := wtPairFreqMap[w+"|"+t]
					capWTPair, _ := capWTPairProbMap[w+"|"+t]
					capPairC:= capWTPairFreqMap[w+"|"+t]
					var val float64
					if okTT && okWT {
						val = lattice[k][i-1].value + ttPairC*ttPair + wtPairC*wtPair + capPairC*capWTPair
					} else if okTT && !okWT {
						val = lattice[k][i-1].value + ttPairC*ttPair + float64(len(wtPairFreqMap))*math.Log(float64(1)/float64(len(wtPairFreqMap))) + float64(len(capWTPairFreqMap))*math.Log(float64(1)/float64(len(capWTPairFreqMap)))
					} else if !okTT && okWT {
						val = lattice[k][i-1].value + float64(len(ttPairFreqMap))*math.Log(float64(1)/float64(len(ttPairFreqMap))) + wtPairC*wtPair + capPairC*capWTPair
					} else if !okTT && !okWT {
						val = lattice[k][i-1].value + float64(len(ttPairFreqMap))*math.Log(float64(1)/float64(len(ttPairFreqMap))) + float64(len(wtPairFreqMap))*math.Log(float64(1)/float64(len(wtPairFreqMap)))+ float64(len(capWTPairFreqMap))*math.Log(float64(1)/float64(len(capWTPairFreqMap)))
					}

					if val >= maxVal {
						maxVal = val
						maxValTag = kTag
					}

				}
				lattice[j][i].value = maxVal
				lattice[j][i].prevTag = maxValTag

			}

		}
		i++
	}

	return BackTrack(lattice, wList)

}

func BackTrack(l [][]cell, wList list.List) (result string) {
	t := len(l)
	w := len(l[0])
	var rslt list.List
	//find max in the last column
	maxLast := l[0][w-1].value
	maxLastTagIdx := 0
	for i := 1; i < t; i++ {
		if maxLast < l[i][w-1].value {
			maxLast = l[i][w-1].value
			maxLastTagIdx = i
		}
	}

	//fmt.Println(maxLastTag, maxLastTagIdx)
	var wordList = make([]string, wList.Len())
	i := 0
	for e := wList.Front(); e != nil; e = e.Next() {
		wordList[i] = e.Value.(string)
		i++
	}

	//fmt.Println(wordList)
	rIdx := maxLastTagIdx
	for i := w - 1; i >= 0; i-- {
		rslt.PushBack(wordList[i] + "_" + strings.ToUpper(l[rIdx][i].currTag))

		for j, v := range tList {
			if strings.EqualFold(v, l[rIdx][i].prevTag) {
				rIdx = j
			}
		}
	}

	x := false
	for e := rslt.Back(); e != nil; e = e.Prev() {

		if x {
			result = result + " "

		}
		result = result + e.Value.(string)
		x = true
	}
	return
}

func FillPairFreqValues(line string) {

	strTokCand := strings.Split(line, " ")
	var strTokTemp list.List
	for i := 0; i < len(strTokCand); i++ {
		if len(strings.TrimSpace(strTokCand[i])) > 0 {
			strTokTemp.PushBack(strTokCand[i])
		}
	}

	j := 0
	var strTok = make([]string, strTokTemp.Len())
	for e := strTokTemp.Front(); e != nil; e = e.Next() {
		strTok[j] = e.Value.(string)
		j++
	}
	startPairList.PushBack(strTok[0])
	endPairList.PushBack(strTok[len(strTok)-1])
	//create W|T frequency table
	for _, v := range strTok {
		pairSplit := strings.Split(v, "_")
		if len(pairSplit) == 2 {
			wtPairFreqMap[pairSplit[0]+"|"+pairSplit[1]] = wtPairFreqMap[pairSplit[0]+"|"+pairSplit[1]] + 1.0
			if byte(pairSplit[0][0])>64 && byte(pairSplit[0][0])<91{
				capWTPairFreqMap[pairSplit[0]+"|"+pairSplit[1]] = capWTPairFreqMap[pairSplit[0]+"|"+pairSplit[1]] + 1.0
			}
		}
	}
	//remove space key from freq table
	delete(wtPairFreqMap, " ")

	//create TT pair freq table
	for i := 0; i < len(strTok)-1; i = i + 1 {
		firstPairSplit := strings.Split(strTok[i], "_")
		secondPairSplit := strings.Split(strTok[i+1], "_")
		if len(firstPairSplit) == 2 && len(secondPairSplit) == 2 {
			if i == 0 {
				ttPairFreqMap[firstPairSplit[1]+"|start"] = ttPairFreqMap[firstPairSplit[1]+"|start"] + 1.0
			}
			if i == len(strTok)-2 {
				ttPairFreqMap["end|"+secondPairSplit[1]] = ttPairFreqMap["end|"+secondPairSplit[1]] + 1.0
			}
			ttPairFreqMap[secondPairSplit[1]+"|"+firstPairSplit[1]] = ttPairFreqMap[secondPairSplit[1]+"|"+firstPairSplit[1]] + 1.0
		}
	}
	return
}