package main

import (
    "sort"
)

type Point struct {
    x, y float32
}

type Points []Point

func (p Points) Len() int {
    return len(p)
}
func (p Points) Swap(i, j int) {
    p[i], p[j] = p[j], p[i]
}
func (p Points) Less(i, j int) bool {
    if p[i].x == p[j].x {
        return p[i].y < p[j].y
    }
    return p[i].x < p[j].x
}

func fastConvexHull(points Points) (Points){
    if len(points) < 2 {
        return points
    }

    sort.Sort(points)
    u := Points{points[0], points[1]}
    for _, p := range points[2:] {
        u = append(u, p)
        for len(u) > 2 && !rightAngle(u[len(u)-3], u[len(u)-2], u[len(u)-1]) {
            u = append(u[:len(u)-2], u[len(u)-1:]...)
        }
    }

    sort.Sort(sort.Reverse(points))
    l := Points{points[0], points[1]}
    for _, p := range points[2:] {
        l = append(l, p)
        for len(l) > 2 && !rightAngle(l[len(l)-3], l[len(l)-2], l[len(l)-1]) {
            l = append(l[:len(l)-2], l[len(l)-1:]...)
        }
    }
    return append(u[:len(u)-1], l[:len(l)-1]...)
}

func rightAngle(o, a, b Point) bool {
    cross := (a.x - o.x) * (b.y - o.y) - (a.y - o.y) * (b.x - o.x)
    println(cross)
    if cross > 0 {
        return false
    } else {
        return true
    }
}
