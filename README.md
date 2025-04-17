# GO-FROM-MAZE

Looking for a way out from the 2D maze. More precisely, we are investigating
and solving an algorithmic problem about finding the shortest path in a maze.

To do this, we will need to design the map model in text form, where:
- `w` = wall,
- `@` = start or current location,
- `Q` = target (exit)

In dynamic, we can see:
- `*` = node
- `.` = trace

Examples of map see in directory "./maps"

## Quick guide

```shell
go run main.go -t hog -f maps/04.txt
```

```
Map size: 12x9
Initial position: [11,1]
Exit position: [5,5]

 w w w w w w w w w w w w
 w                     @
 w   w w w w w w w w   w
 w                 w   w
 w   w   w w w w   w   w
 w   w w w Q   w   w   w
 w   w   w w   w   w   w
 w             w       w
 w w w w w w w w w w w w

Nodes: 13
Router: hare
Routes:
 # 0 : {[[11 1] [1 1] [1 3] [1 7] [6 7] [6 5] [5 5]] 7}: ✅ (Exit) Validation: TRUE

```

## Route animation in debug

```shell
go run main.go -t hare -f maps/01.txt -v 5 -r 0 -D
```

## Exit codes

- 0 - route for exit found
- 1 - no route
- 2 - error
