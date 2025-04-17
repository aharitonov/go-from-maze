package global

import (
	"fmt"
	"math"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type (
	PointOnMap    [2]int
	PointList     []PointOnMap
	RoutingStruct map[PointOnMap]PointList
	PointRegistry map[PointOnMap]bool

	Route struct {
		items  []PointOnMap
		length int
	}

	RouteFrame [2]PointOnMap
	FrameList  []RouteFrame

	RouterResult struct {
		Route          Route
		RecPointLists  []PointList
		RecRouteFrames []RouteFrame
	}
)

func (p PointOnMap) String() string {
	return fmt.Sprintf("(%d,%d)", p[0], p[1])
}

func (p PointOnMap) CalcDistance(toPoint PointOnMap) float64 {
	x1, y1 := p[0], p[1]
	x2, y2 := toPoint[0], toPoint[1]
	return math.Sqrt(float64((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1)))
}

func (pl PointList) ToRoute() Route {
	return Route{
		items:  pl,
		length: len(pl),
	}
}

func (r *Route) Add(p ...PointOnMap) {
	r.items = append(r.items, p...)
	r.length = len(r.items)
}

func (r *Route) IsFinished(item PointOnMap) bool {
	if r.length > 0 && r.items[r.length-1] == item {
		return true
	}
	return false
}

func (r *Route) ToPointMap() map[PointOnMap]int {
	reg := make(map[PointOnMap]int, r.length)
	for i, v := range r.items {
		reg[v] = i
	}
	return reg
}

func (r *Route) ToFrames() FrameList {

	if r.GetLength() == 0 {
		return FrameList(nil)
	}

	frames := make(FrameList, r.GetLength()-1)
	items := r.GetItems()
	for i := 1; i < r.GetLength(); i++ {
		frames[i-1] = RouteFrame{items[i-1], items[i]}

	}
	return frames
}

func RouteFrameNew(x, y, x2, y2 int) RouteFrame {
	return RouteFrame{PointOnMap{x, y}, PointOnMap{x2, y2}}
}

func (rf RouteFrame) Eq(rf2 RouteFrame) bool {
	return rf[0] == rf2[0] && rf[1] == rf2[1]
}

func (fl FrameList) Eq(fl2 FrameList) bool {

	if len(fl) != len(fl2) {
		return false
	}

	for i := 0; i < len(fl); i++ {
		if !fl[i].Eq(fl2[i]) {
			return false
		}
	}
	return true
}

func (rf RouteFrame) Overlap(point PointOnMap) bool {

	x, y := point[0], point[1]
	first, last := rf[0], rf[1]

	if first[0] <= x && x <= last[0] {
		if first[1] <= y && y <= last[1] {
			return true
		}
	}
	return false
}

func (rf RouteFrame) IsHorizontalLine() bool {
	return rf[0][0] == rf[1][0]
}

func (rf RouteFrame) IsVerticalLine() bool {
	return rf[0][0] == rf[1][0]
}

func (r *Route) Copy() Route {
	dst := make([]PointOnMap, len(r.items))
	copy(dst, r.items)
	return Route{items: dst, length: r.length}
}

func (r *Route) GetLength() int {
	return r.length
}

func (r *Route) GetItems() []PointOnMap {
	return r.items
}

func (r *Route) Get(index int) PointOnMap {
	return r.items[index]
}

func (r *Route) Last() PointOnMap {
	if r.length == 0 {
		panic("Empty route on last()")
	}
	return r.items[r.length-1]
}

func (r *Route) Pop() PointOnMap {
	if r.length == 0 {
		panic("Empty route on pop()")
	}
	r.length--
	item := r.items[r.length]
	r.items = r.items[:r.length]
	return item
}

func (r *Route) Reverse() Route {
	slices.Reverse(r.items)
	return *r
}

func (r *Route) Unserialize(src string) *Route {

	if src == "" || src == "[]" {
		r.items = nil
		r.length = len(r.items)
		return r
	}

	a := strings.Split(src, "]")
	items := make([]PointOnMap, 0, len(a))
	for _, v := range a {
		if v == "" {
			break // last empty
		}

		strPair := strings.Trim(v, "] [")
		pair := strings.Split(strPair, " ")

		x, err := strconv.Atoi(pair[0])
		if err != nil {
			panic(err)
		}

		y, err := strconv.Atoi(pair[1])
		if err != nil {
			panic(err)
		}
		items = append(items, PointOnMap{x, y})
	}

	r.items = items
	r.length = len(r.items)
	return r
}

func (r *Route) Serialize() string {
	items := make([]string, r.GetLength())
	for i, point := range r.GetItems() {
		x := strconv.Itoa(point[0])
		y := strconv.Itoa(point[1])
		items[i] = x + " " + y
	}
	return "[" + strings.Join(items, "] [") + "]"
}

func (r *Route) Eq(r2 *Route) bool {
	if r.length != r2.length {
		return false
	}
	for i := range r.items {
		if r.items[i] != r2.items[i] {
			return false
		}
	}
	return true
}

func (rs RoutingStruct) ToKeys() PointList {
	lst := make(PointList, 0, len(rs))
	for pointAsKey := range rs {
		lst = append(lst, pointAsKey)
	}
	return lst
}

func (pl PointList) Sort() {
	slice := pl
	sort.Slice(slice, func(i, j int) bool {
		return slice[i][0]*10+slice[i][1] < slice[j][0]*10+slice[j][1]
	})
}

// FilterWhereNotIn вернёт только те точки, которых нет в маршруте из аргумента
func (pl PointList) FilterWhereNotIn(route Route) PointList {
	pointsInRoute := route.ToPointMap()
	resultPoints := make(PointList, 0, len(pl))
	for _, point := range pl {
		if _, ok := pointsInRoute[point]; !ok {
			resultPoints = append(resultPoints, point)
		}
	}
	return resultPoints
}

// FindByMinDistanceTo вернёт точку с наименьшим расстоянием до точки из аргумента
func (pl PointList) FindByMinDistanceTo(target PointOnMap) PointOnMap {

	if len(pl) == 0 {
		panic("Empty list on findByMinDistanceTo()")
	}

	point := pl[0]
	minDistance := point.CalcDistance(target)
	for i := 1; i < len(pl); i++ {
		dist := pl[i].CalcDistance(target)
		if dist < minDistance {
			minDistance = dist
			point = pl[i]
		}
	}
	return point
}

// FilterList вернёт только те точки из аргумента, которые попадают внутрь фрейма
func (rf RouteFrame) FilterList(points PointList) PointList {

	minPoint, maxPoint := rf[0], rf[1]
	minX, minY := minPoint[0], minPoint[1]
	maxX, maxY := maxPoint[0], maxPoint[1]

	pointsInFrame := make([]PointOnMap, 0, len(points))
	for _, point := range points {
		x, y := point[0], point[1]
		if minX < x && x < maxX && minY < y && y < maxY {
			pointsInFrame = append(pointsInFrame, point)
		}
	}
	return pointsInFrame
}

func (r *Route) Validate() error {
	if r.length < 1 {
		return nil
	}
	items := r.items
	x0, y0 := items[0][0], items[0][1]
	for i := 1; i < r.length; i++ {
		x, y := items[i][0], items[i][1]
		if x != x0 && y != y0 {
			return fmt.Errorf(
				"bad node address in route (%d,%d) -> (%d,%d), step: %d",
				x0, y0, x, y, i,
			)
		}
		x0, y0 = x, y
	}
	return nil
}

func (rs RoutingStruct) Validate() error {
	for node, toNodes := range rs {
		x0, y0 := node[0], node[1]
		for i := 0; i < len(toNodes); i++ {
			x, y := toNodes[i][0], toNodes[i][1]
			if x != x0 && y != y0 {
				return fmt.Errorf(
					"bad node address (%d,%d) -> (%d,%d), node: %v",
					x0, y0, x, y, node,
				)
			}
		}
	}
	return nil
}

func (rs RoutingStruct) SortPointsInValues() {
	for _, points := range rs {
		points.Sort()
	}
}

func (rr *RouterResult) Reverse() {
	rr.Route.Reverse()
	//TODO
	//rr.RecRouteFrames.Reverse()
	//rr.RecRouteFrames.Reverse()
}

func NewRouteFromSlice(slice [][2]int) Route {

	items := make([]PointOnMap, len(slice))
	for i, point := range slice {
		items[i] = point
	}

	return Route{
		items:  items,
		length: len(slice),
	}
}
