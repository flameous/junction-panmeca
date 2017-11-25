package main

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	_ "github.com/lib/pq"
	"github.com/flameous/junction-panmeca/backend/models"
	"strconv"
	"net/http"
	"runtime"
)

func getUserHandler(isDoc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if id, err := strconv.Atoi(c.Param(`id`)); err == nil {
			u := getUser(id, isDoc)
			if u == nil {
				c.String(http.StatusNotFound, `user not found`)
				return
			}
			c.IndentedJSON(http.StatusOK, id)
		} else {
			c.String(http.StatusBadRequest, `getPatient method; invalid url path; error: `+err.Error())
		}
	}
}

func getPatchUserHandler(isDoc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if id, err := strconv.Atoi(c.Param(`id`)); err == nil {
			u := getUser(id, isDoc)
			if u == nil {
				c.String(http.StatusNotFound, `user not found`)
				return
			}
			c.IndentedJSON(http.StatusOK, id)
		} else {
			c.String(http.StatusBadRequest, `getPatient method; invalid url path; error: `+err.Error())
		}
	}
}

func getUser(id int, isDoc bool) models.User {
	var user models.User
	if isDoc {
		user = new(models.Doctor)
	} else {
		user = new(models.Patient)
	}
	result := d.First(user, id)
	if result.RecordNotFound() {
		return nil
	}

	var projects []models.Project
	d.Model(user).Related(&projects)
	user.SetProjects(projects)

	for k := range projects {
		tasks := new([]models.Task)
		d.Model(&projects[k]).Related(tasks)
		projects[k].RelatedTasks = *tasks
	}
	return user
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	patients := router.Group("/patient")
	patients.GET("/:id", getUserHandler(false))
	patients.PATCH("/:id", getPatchUserHandler(false))

	doctors := router.Group("/doctor")
	doctors.GET("/:id", getUserHandler(true))
	doctors.PATCH("/:id", getPatchUserHandler(true))

	projects := router.Group("/project")
	projects.GET("/:id", nil)
	projects.POST("/:id/add_task", nil)

	router.POST("/new_project", nil)

	tasks := router.Group("/task")
	tasks.GET("/:id", nil)
	tasks.PATCH("/:id", nil)
	tasks.POST("/:id/add_image", nil)

	log.Fatal(router.Run(":80"))
}

var (
	d *gorm.DB
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	args := `host=db user=demo dbname=demo sslmode=disable password=demo`
	if runtime.GOOS == "darwin" {
		args = `host=localhost user=flameous dbname=models sslmode=disable`
	}
	var err error
	d, err = gorm.Open(`postgres`, args)
	if err != nil {
		log.Fatal(err)
	}

	d.AutoMigrate(&models.Patient{}, &models.Doctor{}, &models.Task{}, &models.Project{})

	pat := models.Patient{
		FirstName: "patient_name",
		LastName:  "patient_last_name",
		BirthDate: "12-10-1996",
		ExtraData: models.PatientExtraData{String: "patient!!!"},
	}

	log.Println(d.FirstOrCreate(&pat).Error, pat)

	doc := models.Doctor{
		FirstName: "doctor_name",
		LastName:  "doctor_last_name",
		BirthDate: "01-05-1990",
		ExtraData: models.DoctorExtraData{String: "doctor!!!", IsCoolDoctor: true},
	}
	log.Println(d.FirstOrCreate(&doc).Error, doc)

	project := models.Project{
		PatientID:   pat.ID,
		DoctorID:    doc.ID,
		Description: "project description 2",
	}

	log.Println(d.Create(&project).Error, project)

	var pr []models.Project
	log.Println(d.Model(&pat).Related(&pr, "PatientID").Error, pr)

	pr = []models.Project{}
	log.Println(d.Model(&models.Doctor{ID: 1}).Related(&pr, "DoctorID").Error, pr)

	var task = models.Task{Description: "task descr", StartDate: "123", EndDate: "233", ProjectID: project.ID}
	log.Println(d.Create(&task).Error)

	var tasks []models.Task
	log.Println(d.Model(&project).Related(&tasks, "ProjectID").Error, tasks)
}
