package main

import (
	"gin-user-management/initializers"
	"gin-user-management/routes"
)


func init()  {
	initializers.LoadEnvVariables()
	initializers.ConnectToDB()
	initializers.SyncDatabase()
}

func main()  {
	routes.NewRoutes().Run()
}