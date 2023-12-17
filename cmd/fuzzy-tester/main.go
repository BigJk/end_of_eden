package main

import (
	"flag"
	"fmt"
	"github.com/BigJk/end_of_eden/game"
	"github.com/samber/lo"
	"math/rand"
	"os"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

var endTime int64
var seed int64
var mods []string

func main() {
	routines := flag.Int("n", 1, "number of goroutines")
	timeout := flag.Duration("timeout", time.Minute, "length of testing")
	baseSeed := flag.Int64("seed", 0, "random seed")
	modsString := flag.String("mods", "", "mods to load and test, separated by ',' (e.g. mod1,mod2,mod3)")
	help := flag.Bool("help", false, "show help")
	flag.Parse()

	if *help {
		fmt.Println("End Of Eden :: Fuzzy Tester")
		fmt.Println("The fuzzy tester hits a game session with a random number of operations and tries to trigger a panic.")
		fmt.Println()
		flag.PrintDefaults()
		return
	}

	seed = *baseSeed
	endTime = time.Now().Add(*timeout).Unix()

	if len(*modsString) > 0 {
		mods = strings.Split(*modsString, ",")
	}

	if *baseSeed == 0 {
		seed = rand.Int63()
	}

	fmt.Println("N    :", *routines)
	fmt.Println("Seed :", seed)
	fmt.Println("\nWorking...")

	wg := &sync.WaitGroup{}
	for i := 0; i < *routines; i++ {
		wg.Add(1)
		tester(i, wg)
	}

	wg.Wait()
}

func tester(index int, wg *sync.WaitGroup) {
	rnd := rand.New(rand.NewSource(seed + int64(index)))
	opKeys := lo.Keys(Operations)
	stack := [][]string{}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(stack)
			fmt.Println(r)
			fmt.Println(string(debug.Stack()))
			os.Exit(-1)
		}
	}()

	for time.Now().Unix() < endTime {
		s := game.NewSession(game.WithMods(mods))
		ops := 5 + rand.Intn(1000)
		stack = [][]string{}
		s.SetOnLuaError(func(file string, line int, callback string, typeId string, err error) {
			fmt.Println("File     :", file)
			fmt.Println("Line     :", line)
			fmt.Println("Callback :", callback)
			fmt.Println("TypeId   :", typeId)
			fmt.Println("Err      :", err)
			panic("lua error")
		})

		for i := 0; i < ops; i++ {
			next := lo.Shuffle(opKeys)[0]
			stack = append(stack, []string{next, Operations[next](rnd, s)})
		}
	}
}
