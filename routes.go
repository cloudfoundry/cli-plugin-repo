package main

import "github.com/tedsuo/rata"

const (
	Index = "Index"
	List  = "List"
)

var Routes = rata.Routes([]rata.Route{
	{Path: "/", Method: "GET", Name: Index},
	{Path: "/js/:file", Method: "GET", Name: Index},
	{Path: "/css/:file", Method: "GET", Name: Index},
	{Path: "/font/:file", Method: "GET", Name: Index},
	{Path: "/images/:file", Method: "GET", Name: Index},
	{Path: "/list", Method: "GET", Name: List},
})
