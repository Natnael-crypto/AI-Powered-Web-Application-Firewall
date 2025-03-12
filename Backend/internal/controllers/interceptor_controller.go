package controllers

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"path/filepath"

// 	"github.com/gin-gonic/gin"
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// 	"k8s.io/client-go/kubernetes"
// 	"k8s.io/client-go/tools/clientcmd"
// )

// func getKubeClient() *kubernetes.Clientset {
// 	kubeconfig := filepath.Join("/root/.kube/config") // Adjust path as necessary
// 	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
// 	if err != nil {
// 		log.Fatalf("❌ Failed to load kubeconfig: %v", err)
// 	}
// 	clientset, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		log.Fatalf("❌ Failed to create Kubernetes client: %v", err)
// 	}
// 	return clientset
// }

// func scaleInterceptor(replicas int32) error {
// 	clientset := getKubeClient()
// 	deploymentsClient := clientset.AppsV1().Deployments("default") // Change to your namespace

// 	scale, err := deploymentsClient.GetScale(context.Background(), "interceptor", metav1.GetOptions{})
// 	if err != nil {
// 		return fmt.Errorf("❌ Failed to get scale: %v", err)
// 	}
// 	scale.Spec.Replicas = replicas

// 	_, err = deploymentsClient.UpdateScale(context.Background(), "interceptor", scale, metav1.UpdateOptions{})
// 	if err != nil {
// 		return fmt.Errorf("❌ Failed to update scale: %v", err)
// 	}
// 	fmt.Printf("✅ Interceptor scaled to %d replicas\n", replicas)
// 	return nil
// }

// func StartInterceptor(c *gin.Context) {
// 	if err := scaleInterceptor(1); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "✅ Interceptor started"})
// }

// func StopInterceptor(c *gin.Context) {
// 	if err := scaleInterceptor(0); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "✅ Interceptor stopped"})
// }

// func RestartInterceptor(c *gin.Context) {
// 	if err := scaleInterceptor(0); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	if err := scaleInterceptor(1); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "✅ Interceptor restarted"})
// }

// func RepullInterceptor(c *gin.Context) {
// 	clientset := getKubeClient()
// 	deploymentsClient := clientset.AppsV1().Deployments("default") // Change to your namespace

// 	// Perform rolling update by changing the image
// 	interceptorDeployment, err := deploymentsClient.Get(context.Background(), "interceptor", metav1.GetOptions{})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("❌ Failed to get interceptor deployment: %v", err)})
// 		return
// 	}

// 	// Update the image to the latest version
// 	interceptorDeployment.Spec.Template.Spec.Containers[0].Image = "natnaelcrypto/interceptor:latest"
// 	_, err = deploymentsClient.Update(context.Background(), interceptorDeployment, metav1.UpdateOptions{})
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("❌ Failed to update deployment image: %v", err)})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "♻️ Interceptor image repulled and restarted"})
// }

// func main() {
// 	// Initialize Gin router
// 	r := gin.Default()

// 	// Define routes for controlling the interceptor
// 	r.POST("/interceptor/start", startInterceptor)
// 	r.POST("/interceptor/stop", stopInterceptor)
// 	r.POST("/interceptor/restart", restartInterceptor)
// 	r.POST("/interceptor/repull", repullInterceptor)

// 	// Run the server
// 	fmt.Println("Backend WAF Controller running...")
// 	if err := r.Run(":8080"); err != nil {
// 		log.Fatalf("❌ Failed to start server: %v", err)
// 	}
// }
