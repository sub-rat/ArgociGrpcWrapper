package controllers

import "github.com/sub-rat/ArogciGrpcWrapper/api/middleware"

func (server *Server) initializeRoutes(){
	//HomeRoute
	server.Router.HandleFunc("/", middleware.SetMiddleWareJSON(server.Home)).Methods("GET")
	//Login route
	server.Router.HandleFunc("/", middleware.SetMiddleWareJSON(server.Login)).Methods("POST")
	//Users routes
	server.Router.HandleFunc("/users", middleware.SetMiddleWareJSON(server.CreateUser)).Methods("POST")
	server.Router.HandleFunc("/users", middleware.SetMiddleWareJSON(server.GetUsers)).Methods("GET")
	server.Router.HandleFunc("/users/{id}", middleware.SetMiddleWareJSON(server.GetUser)).Methods("GET")
	server.Router.HandleFunc("/users/{id}", middleware.SetMiddleWareJSON(middleware.SetMiddlewareAuthentication(server.UpdateUser))).Methods("PUT")
	server.Router.HandleFunc("/users/{id}", middleware.SetMiddlewareAuthentication(server.DeleteUser)).Methods("DELETE")

	// workflow routes
	server.Router.HandleFunc("/getWorkflow/{name}/{namespace}", middleware.SetMiddleWareJSON(middleware.SetMiddlewareAuthentication(server.GetWorkFlow))).Methods("GET")
	server.Router.HandleFunc("/createWorkflow", middleware.SetMiddleWareJSON(middleware.SetMiddlewareAuthentication(server.CreateWorkFlow))).Methods("POST")
	server.Router.HandleFunc("/getWorkflowList/{namespace}", middleware.SetMiddleWareJSON(middleware.SetMiddlewareAuthentication(server.GetWorkFlowList))).Methods("GET")
	server.Router.HandleFunc("/getWorkflowNames/{namespace}", middleware.SetMiddleWareJSON(middleware.SetMiddlewareAuthentication(server.GetWorkFlowNames))).Methods("GET")
	server.Router.HandleFunc("/getPodLog/{name}/{namespace}", middleware.SetMiddleWareJSON(middleware.SetMiddlewareAuthentication(server.GetPodLog))).Methods("GET")
}
