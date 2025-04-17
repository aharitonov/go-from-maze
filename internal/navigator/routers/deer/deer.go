package deer

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

	// recPointLists Фиксируем множество точек перед отбором
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
			nil,
		},
	}
	return allResults
}

func (rp *plan) doRoute(route Route, fromPoint PointOnMap) Route {

	route.Add(fromPoint)
	points := rp.graph[fromPoint]
	//delete(rp.graph, fromPoint)
	if len(points) > 0 {
		if nextPoint, ok := rp.findNextPoint(route, points); ok {
			return rp.doRoute(route, *nextPoint)
		}
	}
	return route
}

func (rp *plan) findNextPoint(route Route, points PointList) (*PointOnMap, bool) {

	rp.recPointLists = append(rp.recPointLists, points)

	points = points.FilterWhereNotIn(route)

	if len(points) > 0 {
		nextPoint := points.FindByMinDistanceTo(rp.target)
		return &nextPoint, true
	}
	return nil, false
}
