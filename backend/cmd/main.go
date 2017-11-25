package main

import (
	"github.com/flameous/backend-meta/backend/models"
	"github.com/gin-gonic/gin"
	"strconv"
	"net/http"
	"github.com/jinzhu/gorm"
	"log"
	_ "github.com/lib/pq"
	"fmt"
)

func getPatient(c *gin.Context) {
	fmt.Println(c.Param("id"))
	if num, err := strconv.Atoi(c.Param(`id`)); err == nil {
		pat := models.Patient{}
		d.First(&pat, num)

		var project []models.Project
		d.Model(&pat).Related(&project)

		pat.RelatedProjects = project

		for k := range project {
			tasks := new([]models.Task)
			d.Model(&project[k]).Related(tasks)
			project[k].RelatedTasks = *tasks
		}
		c.IndentedJSON(http.StatusOK, pat)
	} else {
		c.String(http.StatusBadRequest, `getPatients method; invalid url path; error: `+err.Error())
	}
}

func main() {
	_ = models.Patient{}

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	patients := router.Group("/patient")
	patients.GET("/:id", getPatient)
	patients.PATCH("/:id", nil)

	doctors := router.Group("/doctor")
	doctors.GET("/:id", nil)
	doctors.PATCH("/:id", nil)

	projects := router.Group("/project")
	projects.GET("/:id", nil)
	projects.POST("/:id/add_task", nil)

	router.POST("/new_project", nil)

	tasks := router.Group("/task")
	tasks.GET("/:id", nil)
	tasks.PATCH("/:id", nil)
	tasks.POST("/:id/add_image", nil)

	router.Run(":80")
}

var (
	d *gorm.DB
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var err error
	d, err = gorm.Open(`postgres`, `host=localhost user=flameous dbname=models sslmode=disable`)
	if err != nil {
		log.Fatal(err)
	}
	d.AutoMigrate(&models.Patient{}, &models.Task{}, &models.Project{})

	//p := models.Patient{
	//	FirstName: "Test",
	//	LastName:  "Patient",
	//	BirthDate: "12-10-1996",
	//	RelatedProjects: []models.Project{
	//		{
	//			Description: "Test Project", RelatedTasks: []models.Task{
	//			{Description: "task desc1", StartDate: "10-10-2016", EndDate: "13-10-2016"},
	//			{Description: "task desc2", StartDate: "01-11-2016", EndDate: "02-11-2016"},
	//		},
	//		},
	//	},
	//}
}
