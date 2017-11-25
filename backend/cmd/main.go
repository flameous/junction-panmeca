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
	"encoding/json"
)

func makeUserHandler(isDoc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if id, err := strconv.Atoi(c.Param("id")); err == nil {
			u := getUser(id, isDoc)
			if u == nil {
				c.String(http.StatusNotFound, "user not found")
				return
			}
			c.IndentedJSON(http.StatusOK, u)
		} else {
			c.String(http.StatusBadRequest, "getPatient method; invalid url path; error: "+err.Error())
		}
	}
}

func makePatchUserHandler(isDoc bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if id, err := strconv.Atoi(c.Param("id")); err == nil {
			u := getUser(id, isDoc)
			if u == nil {
				c.String(http.StatusNotFound, "user not found")
				return
			}
			c.IndentedJSON(http.StatusOK, id)
		} else {
			c.String(http.StatusBadRequest, "getPatient method; invalid url path; error: "+err.Error())
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

func getProject(c *gin.Context) {
	project := new(models.Project)
	d.Find(project, c.Param("id"))

	if project == nil {
		c.String(http.StatusNotFound, "project not found")
	} else {
		tasks := new([]models.Task)
		d.Model(project).Related(tasks, "RelatedTasks")
		project.RelatedTasks = *tasks
		c.IndentedJSON(http.StatusOK, project)
	}
}

func addTaskToProject(c *gin.Context) {
	taskData, ok := c.GetPostForm("task")
	if !ok {
		c.String(http.StatusBadRequest, "where is 'task' ?")
		return
	}

	task := new(models.Task)
	if err := json.Unmarshal([]byte(taskData), task); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	task.ProjectID = uint(id)

	if err := d.Create(task).Error; err != nil {
		c.String(http.StatusInternalServerError, err.Error())
	} else {
		c.String(http.StatusOK, "task added")
	}
}

func newProject(c *gin.Context) {
	data, ok := c.GetPostForm("project")
	if !ok {
		c.String(http.StatusBadRequest, "where is 'project' ?")
		return
	}

	project := new(models.Project)
	if err := json.Unmarshal([]byte(data), project); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	d.Create(project)
	c.String(http.StatusOK, "new project")
}

func getTask(c *gin.Context) {
	task := new(models.Task)
	d.First(task, c.Param("id"))
	if task == nil {
		c.String(http.StatusNotFound, "task not found")
	} else {
		c.IndentedJSON(http.StatusOK, task)
	}

}

func editTask(c *gin.Context) {
	c.String(http.StatusOK, "editTask")
}

func uploadFile(c *gin.Context) {
	ff, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	task := new(models.Task)
	d.Find(task, c.Param("id"))
	task.Image = ff.Filename
	d.Save(task)
	c.String(http.StatusOK, "ok")
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(CORSMiddleware())

	router.GET("/", func(c *gin.Context) { c.String(http.StatusOK, "hello!") })
	router.Static("/static", "../data")

	patients := router.Group("/patient/:id")
	{
		patients.GET("/", makeUserHandler(false))
		patients.PATCH("/", makePatchUserHandler(false))
	}

	doctors := router.Group("/doctor/:id")
	{
		doctors.GET("/", makeUserHandler(true))
		doctors.PATCH("/", makePatchUserHandler(true))
	}

	projects := router.Group("/project/:id")
	{
		projects.GET("/", getProject)
		projects.POST("/add_task", addTaskToProject)
	}

	router.POST("/new_project", newProject)

	tasks := router.Group("/task/:id")
	{
		tasks.GET("/", getTask)
		tasks.PATCH("/", editTask)
		tasks.POST("/upload_image", uploadFile)
	}

	log.Fatal(router.Run(":8080"))
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

var (
	d *gorm.DB
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	args := "host=db user=demo dbname=demo sslmode=disable password=demo"
	if runtime.GOOS == "darwin" {
		args = "host=localhost user=flameous dbname=models sslmode=disable"
	}
	var err error
	d, err = gorm.Open("postgres", args)
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

	log.Printf("create patient: error: %v, data: %v\n", d.FirstOrCreate(&pat).Error, pat)

	doc := models.Doctor{
		FirstName: "doctor_name",
		LastName:  "doctor_last_name",
		BirthDate: "01-05-1990",
		ExtraData: models.DoctorExtraData{String: "doctor!!!", IsCoolDoctor: true},
	}
	log.Printf("create doctor: error: %v, data: %v\n", d.FirstOrCreate(&doc).Error, doc)

	project := models.Project{
		PatientID:   pat.ID,
		DoctorID:    doc.ID,
		Description: "project description",
	}
	log.Printf("create project: err: %v, data: %v\n", d.Create(&project).Error, project)

	var task = models.Task{
		Description: "task description 1",
		StartDate:   "1511827200000",
		EndDate:     "1511913600000",
		ProjectID:   project.ID,
		Image:       "/static/Lower.stl",
	}
	log.Printf("create task 1: error: %v\n", d.Create(&task).Error)

	task = models.Task{
		Description: "task description 2",
		StartDate:   "1512432000000",
		EndDate:     "1512518400000",
		ProjectID:   project.ID,
		Image:       "/static/Upper.stl",
	}
	log.Printf("create task 2, error: %v\n", d.Create(&task).Error)

	task = models.Task{
		Description: "task description 3",
		StartDate:   "1512864000000",
		EndDate:     "1512950400000",
		ProjectID:   project.ID,
		Image:       "/static/Buccal.stl",
	}
	log.Printf("create task 3, error: %v\n", d.Create(&task).Error)
	log.Println("started!")
}
