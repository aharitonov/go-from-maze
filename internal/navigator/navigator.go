package navigator

import (
	"fmt"
	. "maze/internal/global"
	"maze/internal/navigator/routers/deer"
	"maze/internal/navigator/routers/fox"
	"maze/internal/navigator/routers/hare"
	"maze/internal/navigator/routers/hog"
	"maze/internal/navigator/routers/wolf"
	"maze/internal/world"
)

type NavRoute struct {

	// RouterName Поставщик маршрута
	RouterName string

	// Route Результирующий маршрут
	Route *Route

	// RecPointLists По каждой точке маршрута фиксируем множества доступных
	// следующих точек.
	// Нужно здесь для последующего воспроизведения построения и отладки.
	RecPointLists []PointList

	// RecRouteFrames По каждой точке маршрута фиксируем ограничительные окна.
	// Механизм окон может применяться для сужения выбора в принятии решения
	// о следующей точки маршрута.
	// Нужно здесь для последующего воспроизведения построения и отладки.
	RecRouteFrames []RouteFrame

	target PointOnMap
}

func (rr NavRoute) String() string {
	return fmt.Sprintf("%v%s", *rr.Route, rr.GetResultMarker(": "))
}

func (rr NavRoute) IsFoundTarget() bool {
	return rr.Route.IsFinished(rr.target)
}

func (rr NavRoute) GetResultMarker(prefixIfFound string) string {
	foundExit := ""
	if rr.Route.IsFinished(rr.target) {
		foundExit = prefixIfFound + "✅ (Exit)"
	}
	return foundExit
}

type RouterInterface interface {
	BuildRoutes(
		rsConstructor func(reverted bool) RoutingStruct,
		start, target PointOnMap,
		width, height int,
	) []RouterResult
}

const (
	RouterHare = "hare"
	RouterDeer = "deer"
	RouterHog  = "hog"
	RouterFox  = "fox"
	RouterWolf = "wolf"
)

func routerFactory(name string) RouterInterface {
	var r RouterInterface
	switch name {
	case RouterWolf:
		r = wolf.New()
	case RouterFox:
		r = fox.New()
	case RouterHog:
		r = hog.New()
	case RouterDeer:
		r = deer.New()
	case RouterHare:
		r = hare.New()
	default:
		panic("Unknown router")
	}
	return r
}

// FindBestRoute возвращает лучший маршрут
func FindBestRoute(w *world.World) NavRoute {
	return FindRoutes(w, RouterFox)[4]
}

// FindRoutes возвращает массив маршрутов
func FindRoutes(w *world.World, routerName string) []NavRoute {

	router := routerFactory(routerName)
	width, height := w.GetSizes()
	start := w.GetStart().ToArray()
	target := w.GetExit().ToArray()

	rsConstructor := func(reverted bool) RoutingStruct {
		a, b := start, target
		if reverted {
			a, b = b, a
		}
		rs := BuildRoutingTreeFor(w, a, b)
		rs.SortPointsInValues()
		return rs
	}

	allRoutes := router.BuildRoutes(rsConstructor, start, target, width, height)

	//// сверху короткие маршруты
	//sort.Slice(allRoutes, func(i, j int) bool {
	//	return allRoutes[i].Route.GetLength() < allRoutes[j].Route.GetLength()
	//})

	var results []NavRoute
	for _, route := range allRoutes {
		results = append(results, NavRoute{
			target:         target,
			Route:          &route.Route,
			RecRouteFrames: route.RecRouteFrames,
			RecPointLists:  route.RecPointLists,
			RouterName:     routerName,
		})
	}
	return results
}

// BuildRoutingTreeFor выполняет обход в ширину и возвращает дерево локаций
// для прокладывания маршрутов
func BuildRoutingTreeFor(w *world.World, start, target PointOnMap) RoutingStruct {

	tree := RoutingStruct{target: PointList{}}
	queue := append(PointList{}, start)
	queueRegistry := PointRegistry{start: true}
	pointReg := PointRegistry{}

	for i := 0; i < len(queue); i++ {
		point := queue[i]
		for _, mov := range w.FindNextMoves(point[0], point[1], target[0], target[1]) {
			nextPoint := PointOnMap{mov[0], mov[1]}
			pointReg[nextPoint] = true
			if _, exist := queueRegistry[nextPoint]; !exist {
				//pointReg[nextPoint] = true
				if target == nextPoint {
					continue
				}
				queue = append(queue, nextPoint)
				queueRegistry[nextPoint] = true
			}
		}

		if total := len(pointReg); total > 0 {
			pl := make(PointList, 0, total)
			for p := range pointReg {
				pl = append(pl, p)
			}
			tree[point] = pl
			clear(pointReg)
		}
	}

	return tree
}

func BuildRoutingTree(w *world.World) RoutingStruct {
	start := w.GetStart().ToArray()
	target := w.GetExit().ToArray()
	rs := BuildRoutingTreeFor(w, start, target)
	rs.SortPointsInValues() // reach stable routes
	return rs
}

func PrintRoutingTree(rs RoutingStruct) {
	keys := rs.ToKeys()
	keys.Sort()
	for _, mapKey := range keys {
		points := rs[mapKey]
		points.Sort()
		fmt.Println(" ", mapKey, "=>", points)
	}
}
