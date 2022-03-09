package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	orderIdSymbolMap := make(map[string]string)
	symbolVolMap := make(map[string]int)
	for scanner.Scan() {
		message := string(scanner.Bytes())

		fieldMap := make(map[string]string)
		if len(message) > 40 {
			fieldMap["initChar"] = message[0:1]
			fieldMap["timeStamp"] = message[1:9]
			fieldMap["messageType"] = message[9:10]
		} else {
			continue
		}

		switch fieldMap["messageType"] {
		case "A": //Add Order message (short)
			if len(message) == 46 {
				fieldMap["orderID"] = message[10:22]
				fieldMap["sideIndicator"] = message[22:23]
				fieldMap["shares"] = message[23:29]
				fieldMap["stockSymbol"] = strings.TrimSpace(message[29:35])
				fieldMap["price"] = message[35:45]
				fieldMap["display"] = message[45:46]

				orderID := fieldMap["orderID"]
				orderIdSymbolMap[orderID] = fieldMap["stockSymbol"]
			}
			break
		case "E": //Order Executed
			if len(message) == 40 {
				fieldMap["orderID"] = message[10:22]
				fieldMap["executedShares"] = message[22:28]
				fieldMap["executionID"] = message[28:40]

				orderID := fieldMap["orderID"]
				stockSymbol := orderIdSymbolMap[orderID]

				if len(stockSymbol) > 0 {
					shares, _ := strconv.Atoi(fieldMap["executedShares"])
					symbolVolMap[stockSymbol] += shares
				}
			}
			break
		case "P": //Trade message (short)
			if len(message) == 57 {
				fieldMap["orderID"] = message[10:22]
				fieldMap["sideIndicator"] = message[22:23]
				fieldMap["shares"] = message[23:29]
				fieldMap["stockSymbol"] = strings.TrimSpace(message[29:35])
				fieldMap["price"] = message[35:45]
				fieldMap["executionID"] = message[45:57]

				stockSymbol := fieldMap["stockSymbol"]
				shares, _ := strconv.Atoi(fieldMap["shares"])
				symbolVolMap[stockSymbol] += shares
				break
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	shareCount := sortSymbolsByVol(symbolVolMap)
	if len(shareCount) > 0 {

		n := 10
		// checks if total number of symbols is less than n
		if len(shareCount) < n {
			n = len(shareCount)
		}
		fmt.Println("----------------------------------")
		fmt.Printf("Top %d Symbols by Executed Volume\n", n)
		fmt.Println("----------------------------------")
		for i := 0; i < n; i++ {
			fmt.Printf("%s\t%d\n", shareCount[i].Symbol, shareCount[i].Vol)
		}
	}
}

type Pair struct {
	Symbol string
	Vol    int
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Vol < p[j].Vol }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// sort symbolVolMap by volume
func sortSymbolsByVol(symbolVolMap map[string]int) PairList {
	pl := make(PairList, len(symbolVolMap))
	i := 0
	for k, v := range symbolVolMap {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}
