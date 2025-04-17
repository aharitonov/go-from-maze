package hare

import (
	. "maze/internal/global"
)

type ThisRouter struct {
}

func New() *ThisRouter {
	return &ThisRouter{}
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
}

func (r ThisRouter) BuildRoutes(
	rsProvider func(reverted bool) RoutingStruct,
	start, target PointOnMap,
	width, height int,
) []RouterResult {

	rp := &plan{
		width:  width,
		height: height,
		start:  start,
		target: target,
		graph:  rsProvider(false),
	}

	allResults := []RouterResult{
		{
			rp.doRoute(Route{}, rp.start),
			rp.recPointLists,
			rp.recRouteFrames,
		},
	}
	return allResults
}

func (rp *plan) doRoute(route Route, fromPoint PointOnMap) Route {

	route.Add(fromPoint)
	points := rp.graph[fromPoint]

	// Предотвращаем зацикливание/наложение маршрута с помощью фреймов
	//delete(rp.graph, fromPoint) // we use frames !

	if len(points) > 0 {
		if nextPoint, ok := rp.findNextPoint(route, points); ok {
			return rp.doRoute(route, *nextPoint)
		}
	}
	return route
}

func (rp *plan) findNextPoint(route Route, points PointList) (*PointOnMap, bool) {

	rp.recPointLists = append(rp.recPointLists, points)

	pointsInFrame := PointList{}

	frmMinPoint, frmMaxPoint := rp.lastFrame()
	minX, minY := frmMinPoint[0], frmMinPoint[1]
	maxX, maxY := frmMaxPoint[0], frmMaxPoint[1]
	for _, point := range points {
		x, y := point[0], point[1]
		if minX < x && x < maxX && minY < y && y < maxY {
			pointsInFrame = append(pointsInFrame, point)
		}
	}

	if len(pointsInFrame) > 0 {
		nextPoint := pointsInFrame.FindByMinDistanceTo(rp.target)
		rp.addFrame(route, nextPoint)
		return &nextPoint, true
	}
	return nil, false
}

func (rp *plan) lastFrame() (PointOnMap, PointOnMap) {
	frames := rp.recRouteFrames
	if len(frames) > 0 {
		framePoints := frames[len(frames)-1]
		return framePoints[0], framePoints[1]
	}
	return PointOnMap{-1, -1}, PointOnMap{rp.width, rp.height}
}

func (rp *plan) addFrame(route Route, point PointOnMap) {
	minX, maxX := rp.findRouteFrameOnX(route, point)
	minY, maxY := rp.findRouteFrameOnY(route, point)
	rp.recRouteFrames = append(rp.recRouteFrames, [2]PointOnMap{
		{minX, minY},
		{maxX, maxY},
	})
}

func (rp *plan) findRouteFrameOnX(nodes Route, fromNode PointOnMap) (min, max int) {
	fromX, fromY := fromNode[0], fromNode[1]
	min, max = -1, rp.width
	for _, point := range nodes.GetItems() {
		x, y := point[0], point[1]
		if y == fromY {
			if x == fromX {
				continue
			}
			if x > fromX {
				if max > x {
					max = x
				}
			} else {
				if min < x {
					min = x
				}
			}
		}
	}
	return min, max
}

func (rp *plan) findRouteFrameOnY(nodes Route, fromNode PointOnMap) (min, max int) {
	fromX, fromY := fromNode[0], fromNode[1]
	min, max = -1, rp.height
	for _, point := range nodes.GetItems() {
		x, y := point[0], point[1]
		if x == fromX {
			if y == fromY {
				continue
			}
			if y > fromY {
				if max > y {
					max = y
				}
			} else {
				if min < y {
					min = y
				}
			}
		}
	}
	return min, max
}
