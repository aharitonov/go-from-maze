package hog

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
	p := tr.plan
	return RouterResult{
		Route:          route,
		RecPointLists:  p.recPointLists,
		RecRouteFrames: p.recRouteFrames,
	}
}

type plan struct {
	width  int
	height int
	start  PointOnMap
	target PointOnMap
	graph  RoutingStruct

	// recRouteFrames Фиксируем размеры фреймов для каждой точки маршрута
	recRouteFrames []RouteFrame
	// recPointLists Фиксируем множества точек перед фреймовым отбором
	recPointLists []PointList

	ex PointRegistry
}

func (tr *ThisRouter) BuildRoutes(
	rsProvider func(reverted bool) RoutingStruct,
	start, target PointOnMap,
	width, height int,
) []RouterResult {

	tr.plan = &plan{
		width:  width,
		height: height,
		start:  start,
		target: target,
		graph:  rsProvider(false),
		ex:     make(PointRegistry),
	}

	var allResults []RouterResult

	point := start       // first point for route
	routeRef := &Route{} // initial route

	for {
		tr.doRoute(routeRef, point)
		if routeRef.IsFinished(target) {
			allResults = append(allResults, tr.createResult(*routeRef))
		}

		if route := *routeRef; route.GetLength() > 1 {
			route.Pop()
			point = route.Pop()
			routeCopy := route.Copy()
			routeRef = &routeCopy

			newLen := routeCopy.GetLength()
			tr.plan.recRouteFrames = tr.plan.recRouteFrames[:newLen]
			tr.plan.recPointLists = tr.plan.recPointLists[:newLen]
			continue
		}
		break
	}

	return allResults
}

func (tr *ThisRouter) doRoute(routeRef *Route, point PointOnMap) {

	rp := tr.plan
	route := *routeRef

	route.Add(point)
	*routeRef = route

	nextPoints := rp.stepByPoint(point)
	if len(nextPoints) > 0 {
		if nextPoint, ok := tr.selectNextPoint(route, nextPoints); ok {
			tr.doRoute(routeRef, *nextPoint)
		}
	}

	rp.exclude(point)
}

func (tr *ThisRouter) selectNextPoint(route Route, points PointList) (*PointOnMap, bool) {

	rp := tr.plan
	rp.recPointLists = append(rp.recPointLists, points)

	pointsInFrame := rp.lastFrame().FilterList(points)
	pointList := pointsInFrame.FilterWhereNotIn(route)

	if len(pointList) > 0 {
		nextPoint := pointList.FindByMinDistanceTo(rp.target)
		rp.addFrame(route, nextPoint)
		return &nextPoint, true
	}
	return nil, false
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

// lastFrame возвращает фрейм, для сужения выбора точек маршрута
func (rp *plan) lastFrame() RouteFrame {
	frames := rp.recRouteFrames
	if len(frames) > 0 {
		framePoints := frames[len(frames)-1]
		return RouteFrame{framePoints[0], framePoints[1]}
	}
	return RouteFrame{
		PointOnMap{-1, -1},
		PointOnMap{rp.width, rp.height},
	}
}

// addFrame здесь не используются фреймы, поэтому всегда подсовываем фреймы
// во всё поле
func (rp *plan) addFrame(route Route, point PointOnMap) {
	_, _ = route, point
	minX, maxX := -1, rp.width
	minY, maxY := -1, rp.height
	rp.recRouteFrames = append(rp.recRouteFrames, RouteFrame{
		PointOnMap{minX, minY},
		PointOnMap{maxX, maxY},
	})
}
