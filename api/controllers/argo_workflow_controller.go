package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/argoproj/argo/pkg/apiclient/workflow"
	"github.com/argoproj/argo/pkg/apis/workflow/v1alpha1"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/sub-rat/ArogciGrpcWrapper/api/responses"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/testdata"
	"io"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	v13 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func init() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}
}
var (
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containing the CA root cert file")
	serverAddr         = flag.String("server_addr", "localhost:8080", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name used to verify the hostname returned by the TLS handshake")
)

type Clients struct {
	*v13.CoreV1Client
}


func (server *Server) GRPCConnection() *grpc.ClientConn {

	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		log.Println("Successfully created TLS credentials")
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		log.Println("creates with insecure")
		opts = append(opts, grpc.WithInsecure())
	}
	opts = append(opts, grpc.WithBlock())
	log.Println("Dialing ...")
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	return conn
}

func kubeConnection() *rest.Config{
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}
	return config
}

func (server *Server) GetWorkFlowList(w http.ResponseWriter, r *http.Request){
	namespace := mux.Vars(r)["namespace"]

	conn := server.GRPCConnection()
	defer conn.Close()

	wf :=workflow.NewWorkflowServiceClient(conn)
	var resp *v1alpha1.WorkflowList
	resp, err := wf.ListWorkflows(context.Background(),&workflow.WorkflowListRequest{Namespace: namespace})
	if err != nil {
		log.Fatalf("fil to fetch list %v",err)
	}
	for i, s:= range resp.Items {
	  log.Println(string(i) + " " + s.Name)
	}
	responses.JSON(w, http.StatusOK, resp)
}

func (server *Server) GetWorkFlowNames(w http.ResponseWriter, r *http.Request){
	namespace := mux.Vars(r)["namespace"]

	conn := server.GRPCConnection()
	defer conn.Close()

	wf :=workflow.NewWorkflowServiceClient(conn)
	fmt.Println("Grpc connection Established")
	var resp *v1alpha1.WorkflowList
	resp, err := wf.ListWorkflows(context.Background(),&workflow.WorkflowListRequest{Namespace: namespace})
	if err != nil {
		log.Fatalf("fil to fetch list %v",err)
	}
	var name []string
	for _, s:= range resp.Items {
		name = append(name,  s.Name)
	}
	responses.JSON(w, http.StatusOK, name)
}

func (server *Server) GetWorkFlow(w http.ResponseWriter, r *http.Request){
	namespace := mux.Vars(r)["namespace"]
	name := mux.Vars(r)["name"]
	conn := server.GRPCConnection()
	defer conn.Close()

	wf :=workflow.NewWorkflowServiceClient(conn)
	var resp *v1alpha1.Workflow
	resp, err := wf.GetWorkflow(context.Background(),&workflow.WorkflowGetRequest{
		Namespace: namespace,
		Name: name,
	})
	if err != nil {
		log.Fatalf("fil to fetch list %v",err)
	}
	responses.JSON(w, http.StatusOK, resp)
}

func (server *Server) CreateWorkFlow(w http.ResponseWriter, r *http.Request){
	conn := server.GRPCConnection()
	defer conn.Close()
	wf := workflow.NewWorkflowServiceClient(conn)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	var req workflow.WorkflowCreateRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	resp, err := wf.CreateWorkflow(context.Background(), &req)
	if err != nil {
		fmt.Println(err)
		responses.ERROR(w, http.StatusUnprocessableEntity,err)
	}
	responses.JSON(w,http.StatusOK,resp)
}

func (server *Server) GetPodLog(w http.ResponseWriter, r *http.Request) {
	namespace := mux.Vars(r)["namespace"]
	name := mux.Vars(r)["name"]
	var pod  v1.Pod
	pod.Name = name
	pod.Namespace = namespace
	logs := getPodLogs(pod)
	responses.JSON(w, http.StatusOK, logs)
}

func getPodLogs(pod v1.Pod) string {
	podLogOpts := v1.PodLogOptions{}
	podLogOpts.Container = "main"
	clients, err := kubernetes.NewForConfig(kubeConnection())
	if err != nil {
		return "error in getting access to K8S"
	}
	req := clients.CoreV1().
		Pods(pod.Namespace).
		GetLogs(pod.Name, &podLogOpts)
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		log.Println(err)
		return "error in opening stream"
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "error in copy information from podLogs to buf"
	}
	str := buf.String()

	return str
}