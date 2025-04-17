package fox

import (
	. "maze/internal/global"
)

type ThisRouter struct {
	plan *plan
}

func New() *ThisRouter {
	return &ThisRouter{}
}

func (tr *ThisRouter) createResult(route Route) RouterResult {
	return RouterResult{
		Route:          route,
		RecPointLists:  tr.plan.recPointLists,
		RecRouteFrames: nil,
	}
}

type plan struct {
	width  int
	height int
	start  PointOnMap
	target PointOnMap
	graph  RoutingStruct

	// recPointLists Фиксируем множества точек перед фреймовым отбором
	recPointLists []PointList

	ex PointRegistry
}

func (tr *ThisRouter) BuildRoutes(
	rsProvider func(reverted bool) RoutingStruct,
	start, target PointOnMap,
	width, height int,
) []RouterResult {

	var allResults []RouterResult

	tr.plan = &plan{
		width:  width,
		height: height,
		start:  start,
		target: target,
		graph:  rsProvider(false),
		ex:     make(PointRegistry),
	}
	r1 := tr.build()

	tr.plan = &plan{
		width:  width,
		height: height,
		start:  target, // correct
		target: start,  // correct
		graph:  rsProvider(true),
		ex:     make(PointRegistry),
	}
	r2 := tr.build()

	if r1 != nil {
		// #0 прямой маршрут
		allResults = append(allResults, *r1)
	}

	if r2 != nil {
		r2.Reverse()
		// #1 обратный маршрут
		allResults = append(allResults, *r2)
	}

	if r1 != nil && r2 != nil {
		r3 := shortByIntersection((*r2).Route, (*r1).Route)
		// #2 маршрут построенный по наложению прямого и обратного
		allResults = append(allResults, RouterResult{
			Route:          r3,
			RecPointLists:  nil,
			RecRouteFrames: nil,
		})
		// #3 пересечение прямого и обратного маршрута без петляний
		allResults = append(allResults, RouterResult{
			Route:          unlooping(r3),
			RecPointLists:  nil,
			RecRouteFrames: nil,
		})
		// #4 пересечение прямого и обратного маршрута без петляний по одной прямой
		allResults = append(allResults, RouterResult{
			Route:          unloopingLine(r3),
			RecPointLists:  nil,
			RecRouteFrames: nil,
		})
	}

	return allResults
}

func (tr *ThisRouter) build() *RouterResult {

	point := tr.plan.start // first point for route
	routeRef := &Route{}   // initial route

	for {
		tr.doRoute(routeRef, point)
		if routeRef.IsFinished(tr.plan.target) {
			ret := tr.createResult(*routeRef)
			return &ret
		}

		if route := *routeRef; route.GetLength() > 1 {
			route.Pop()
			point = route.Pop()
			routeCopy := route.Copy()
			routeRef = &routeCopy

			newLen := routeCopy.GetLength()
			tr.plan.recPointLists = tr.plan.recPointLists[:newLen]
			continue
		}
		break
	}
	return nil
}

func (tr *ThisRouter) doRoute(routeRef *Route, point PointOnMap) {

	rp := tr.plan
	route := *routeRef

	route.Add(point)
	*routeRef = route

	nextPoints := rp.stepByPoint(point)
	if len(nextPoints) > 0 {
		if nextPoint := tr.selectNextPoint(route, nextPoints); nextPoint != nil {
			tr.doRoute(routeRef, *nextPoint)
		}
	}

	rp.exclude(point)
}

func (tr *ThisRouter) selectNextPoint(route Route, points PointList) *PointOnMap {

	rp := tr.plan
	rp.recPointLists = append(rp.recPointLists, points)

	pointList := points.FilterWhereNotIn(route)
	if len(pointList) > 0 {
		nextPoint := pointList.FindByMinDistanceTo(rp.target)
		return &nextPoint
	}
	return nil
}

func (rp *plan) stepByPoint(point PointOnMap) PointList {
	preNextPoints := rp.graph[point]
	//delete(rp.graph, point
	list := make(PointList, 0, len(preNextPoints))
	for _, prePoint := range preNextPoints {
		if !rp.isExcluded(prePoint) {
			list = append(list, prePoint)
		}
	}
	return list
}

func (rp *plan) exclude(point PointOnMap) {
	rp.ex[point] = true
}

func (rp *plan) isExcluded(point PointOnMap) bool {
	_, ok := rp.ex[point]
	return ok
}

// shortByIntersection сравнивает длину участков между пересечениями и собирает
// кратчайший маршрут. Вернёт кротчайший маршрут, независимо от наличия или
// отсутствия пересечений
func shortByIntersection(routeA, routeB Route) Route {

	if routeA.GetLength() > routeB.GetLength() {
		routeA, routeB = routeB, routeA
	}

	from, to := 1, routeA.GetLength()-1 // skip start and finish
	pointListA, pointListB := routeA.GetItems(), routeB.GetItems()
	lastIndexA, lastIndexB := 0, 0
	pointToIndexB := routeB.ToPointMap()

	var segments []PointList

	for indexA := from; indexA < to; indexA++ {
		pointA := routeA.Get(indexA)
		if indexB, ok := pointToIndexB[pointA]; ok {
			if indexB < lastIndexB {
				break
			}
			if indexB-lastIndexB < indexA-lastIndexA {
				segments = append(segments, pointListB[lastIndexB:indexB])
			} else {
				segments = append(segments, pointListA[lastIndexA:indexA])
			}
			lastIndexA, lastIndexB = indexA, indexB
		}
	}

	if len(segments) > 0 {
		tailA := pointListA[lastIndexA:]
		tailB := pointListB[lastIndexB:]

		if len(tailA) < len(tailB) {
			segments = append(segments, tailA)
		} else {
			segments = append(segments, tailB)
		}

		pointList := make(PointList, 0, routeA.GetLength())
		for _, segment := range segments {
			for _, point := range segment {
				pointList = append(pointList, point)
			}
		}
		return pointList.ToRoute()
	}

	// no intersection
	return routeA
}

// unlooping устраняет петляния
func unlooping(route Route) Route {

	var result = make(PointList, 0, route.GetLength())
	var items = route.GetItems()

	for i := 0; i < len(items); i++ {
		node := items[i]
		for j := len(items) - 1; i < j; j-- {
			if items[i] == items[j] {
				i = j
				break
			}
		}
		result = append(result, node)
	}

	return result.ToRoute()
}

// unloopingLine устраняет петляния: заменяет множество точек на одной вертикали
// или горизонтали парой крайних точек
func unloopingLine(route Route) Route {

	var result = make(PointList, 0, route.GetLength())
	var items = route.GetItems()

	yLast := -1
	yBeginIndex := 0
	yRepeat := 0

	xLast := -1
	xBeginIndex := 0
	xRepeat := 0

	lastAddedIndex := 0
	for i := 0; i < len(items); i++ {

		y := items[i][1]
		if yLast != y {
			if yRepeat > 2 {
				result = append(result, items[lastAddedIndex:yBeginIndex+1]...)
				result = append(result, items[i-1])
				lastAddedIndex = i
			}
			yRepeat = 1
			yBeginIndex = i
		} else {
			yRepeat++
		}
		yLast = y

		x := items[i][0]
		if xLast != x {
			if xRepeat > 2 {
				result = append(result, items[lastAddedIndex:xBeginIndex+1]...)
				result = append(result, items[i-1])
				lastAddedIndex = i
			}
			xRepeat = 1
			xBeginIndex = i
		} else {
			xRepeat++
		}
		xLast = x
	}

	tail := items[lastAddedIndex:]

	if repeat := max(xRepeat, yRepeat); repeat > 2 {
		length := len(tail)
		offset := length - repeat
		lastPair := append(tail[offset:offset+1], tail[length-1:length]...)
		tail = append(tail[:offset], lastPair...)
	}

	result = append(result, tail...)
	return result.ToRoute()
}
