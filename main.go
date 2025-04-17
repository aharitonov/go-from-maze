package main

import (
	"bufio"
	"flag"
	"fmt"
	"maze/internal/cli"
	"maze/internal/navigator"
	"maze/internal/world"
	"os"
	"strings"
	"time"
)

const (
	ExitTargetNotFound = 1
	ExitError          = 2
)

const (
	debugMode             = 0
	defaultRouter         = navigator.RouterFox
	defaultAnimationSpeed = 3
	useTraceOnMove        = true
)

type configParams struct {
	animateRoute        int
	animationSpeed      int
	routerType          string
	stdinFlag           bool
	debugAnimationFlag  bool
	revertDirectionFlag bool
	showRoutingTreeFlag bool
	fromFile            string
}

var params configParams

func initMain() {
	animateRouteArg := flag.Int("r", -1, "animate route by number (if presented)")
	animateSpeedArg := flag.Int("v", defaultAnimationSpeed, "set animation speed")
	routerTypeArg := flag.String("t", defaultRouter, "router type: hare,deer,hog")
	stdinFlagArg := flag.Bool("i", false, "read world from stdin")
	debugAnimationArg := flag.Bool("D", false, "use debug animation (if provided by router)")
	showRoutingTreeArg := flag.Bool("T", false, "show routing tree")
	revertDirectionArg := flag.Bool("R", false, "swap start and finish")
	fromFileArg := flag.String("f", "", "read world from file")

	flag.Parse()

	params = configParams{
		animateRoute:        *animateRouteArg,
		animationSpeed:      *animateSpeedArg,
		routerType:          *routerTypeArg,
		stdinFlag:           *stdinFlagArg,
		debugAnimationFlag:  *debugAnimationArg,
		revertDirectionFlag: *revertDirectionArg,
		showRoutingTreeFlag: *showRoutingTreeArg,
		fromFile:            *fromFileArg,
	}

	if isDebug() { // DEBUG
		params.showRoutingTreeFlag = true
		params.routerType = navigator.RouterFox
		params.fromFile = "./maps/09.txt"
		//params.revertDirectionFlag = true
		//params.animateRoute = 0
	}
}

func main() {

	initMain()
	cli.ClearScreen()
	w := constructWorld(params)
	if isDebug() {
		fmt.Println(cli.WarnStyle("DEBUG MODE: ON"))
	}
	world.PrintMe(w)

	foundRoutes := navigator.FindRoutes(w, params.routerType)

	if params.animateRoute != -1 {
		if params.animateRoute < 0 || params.animateRoute >= len(foundRoutes) {
			fatalExit("Route not found")
		}
		result := foundRoutes[params.animateRoute]
		animate(w, result, params.animationSpeed, params.debugAnimationFlag)
		fmt.Println()
		fmt.Println("Routes:")
		showRoutes(foundRoutes)
		if !hasExit(foundRoutes) {
			os.Exit(ExitTargetNotFound)
		}
		return
	}

	tree := navigator.BuildRoutingTree(w)

	if params.showRoutingTreeFlag {
		fmt.Println("Routing tree:")
		navigator.PrintRoutingTree(tree)
		if err := tree.Validate(); err != nil {
			fmt.Println(" ", cli.ErrorStyle("Validation: FALSE"))
			fmt.Println(err)
		} else {
			fmt.Println(" ", cli.ShadowStyle("Validation: TRUE"))
		}
		fmt.Println()
	}

	fmt.Println("Nodes:", len(tree))
	fmt.Println("Router:", params.routerType)
	fmt.Println("Routes:")
	if len(foundRoutes) > 0 {
		showRoutes(foundRoutes)
		fmt.Print(cli.ShadowStyle("HELP: -r for animate route. Example:"))
		fmt.Println(cli.ShadowStyle(" ./main -f path/to/map3.txt -v 5 -r 1"))
	} else {
		fmt.Println()
		fmt.Println(" ", "No route!!")
		fmt.Println()
	}

	if !hasExit(foundRoutes) {
		os.Exit(ExitTargetNotFound)
	}
}

func hasExit(items []navigator.NavRoute) bool {
	for _, n := range items {
		if n.IsFoundTarget() {
			return true
		}
	}
	return false
}

func constructWorld(params configParams) *world.World {

	fromFile := params.fromFile
	fromStdin := params.stdinFlag
	revertDirectionFlag := params.revertDirectionFlag

	var textWorld string

	if fromStdin {
		textWorld = loadTextWorldFromStdin()
	} else if fromFile != "" {
		byteContent, err := os.ReadFile(fromFile)
		if err != nil {
			fatalExit(err)
		}
		textWorld = string(byteContent)
	} else {
		fatalExit("No world content. Use -f for load from file or -i for load from stdin")
	}

	w, err := world.Construct(textWorld)
	if err != nil {
		fatalExit(err)
	}

	if revertDirectionFlag {
		start, finish := w.GetStart(), w.GetExit()
		w.SetStart(finish)
		w.SetExit(start)
	}
	return w
}

func loadTextWorldFromStdin() string {

	if !cli.UsedStdin() {
		fatalExit("data input from STDIN was expected")
	}

	sc := bufio.NewScanner(os.Stdin)
	lines := make([]string, 0, 25)
	for sc.Scan() {
		line := sc.Text()
		lines = append(lines, line)
	}
	if err := sc.Err(); err != nil {
		fatalExit(err)
	}

	if len(lines) < 1 {
		fatalExit("no lines found")
	}
	return strings.Join(lines, "\n")
}

func fatalExit(e interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, "FATAL: %v\n", e)
	os.Exit(ExitError)
}

func animate(w *world.World, result navigator.NavRoute, speed int, useFrameAnimation bool) {

	cli.ClearScreen()
	speedValue := time.Duration(speed)
	cmdString := cli.GetExecutedCommand()

	route := *result.Route

	for i, node := range route.GetItems() {
		if err := w.Move(node[0], node[1], useTraceOnMove); err != nil {
			panic(node)
		}

		lineForShow := 1
		cli.SetCursorPosition(0, 0)
		if isDebug() {
			fmt.Println(cli.WarnStyle("DEBUG MODE: ON"))
			lineForShow++
		}
		fmt.Println("Executed:", cmdString)
		world.PrintMe(w)
		fmt.Println("Router:", result.RouterName)
		fmt.Println(" ", route.GetItems()[:i+1])
		time.Sleep(time.Millisecond * 1500 / speedValue)

		if !useFrameAnimation {
			continue
		}

		func(iFrame int) {
			wData := w.Pack()
			for k := 0; k < 2; k++ {

				if l := len(result.RecRouteFrames); 0 < l && iFrame < l-1 {
					frm := result.RecRouteFrames[iFrame]
					frmMin, frmMax := frm[0], frm[1]
					frmMaxX, frmMaxY := frmMax[0]-1, frmMax[1]-1
					frmMinX, frmMinY := frmMin[0]+1, frmMin[1]+1
					w.SetRectangle(frmMinX, frmMinY, frmMaxX, frmMaxY)
				}

				if l := len(result.RecPointLists); 0 < l && iFrame < l-1 {
					for _, p := range result.RecPointLists[iFrame] {
						w.SetPoint(p[0], p[1], '?')
					}

					cli.SetCursorPosition(0, lineForShow)
					world.PrintMapOnly(w)
					time.Sleep(time.Millisecond * 1000 / speedValue)

					w.Unpack(wData)
					cli.SetCursorPosition(0, lineForShow)
					world.PrintMapOnly(w)
					time.Sleep(time.Millisecond * 1000 / speedValue)
				}
			}
			w.Unpack(wData)
		}(i)
	}
}

func showRoutes(results []navigator.NavRoute) {
	for n, route := range results {
		validationSign := cli.ShadowStyle("Validation: TRUE")
		if err := route.Route.Validate(); err != nil {
			validationSign = cli.ErrorStyle("Validation: FALSE")
		}
		fmt.Println()
		fmt.Println(" #", n, ":", route, validationSign)
	}
	fmt.Println()
}

func isDebug() bool {
	return debugMode != 0
}
