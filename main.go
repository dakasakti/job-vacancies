package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jaevor/go-nanoid"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type BEJob struct {
	gorm.Model  `json:"-"`
	DataID      string    `json:"dataId" gorm:"type:varchar(50);unique"`
	Date        time.Time `json:"date" gorm:"type:date"`
	CompanyName string    `json:"companyName" gorm:"type:varchar(255)"`
	JobPosition string    `json:"jobPosition" gorm:"type:varchar(255)"`
	WorkType    string    `json:"workType,omitempty" gorm:"type:varchar(50)"`
	TechStack   string    `json:"techStack" gorm:"type:varchar(255)"`
	LinkToJob   string    `json:"linkToJob" gorm:"type:varchar(255)"`
	Industry    string    `json:"industry,omitempty" gorm:"type:varchar(5)"`
	Status      string    `json:"status" gorm:"type:varchar(10);default:'open'"`
	JobID       uint      `json:"-"`
}

type FEJob struct {
	gorm.Model  `json:"-"`
	DataID      string    `json:"dataId" gorm:"type:varchar(50);unique"`
	Date        time.Time `json:"date" gorm:"type:date"`
	CompanyName string    `json:"companyName" gorm:"type:varchar(255)"`
	JobPosition string    `json:"jobPosition" gorm:"type:varchar(255)"`
	WorkType    string    `json:"workType,omitempty" gorm:"type:varchar(50)"`
	TechStack   string    `json:"techStack" gorm:"type:varchar(255)"`
	LinkToJob   string    `json:"linkToJob" gorm:"type:varchar(255)"`
	Industry    string    `json:"industry,omitempty" gorm:"type:varchar(5)"`
	Status      string    `json:"status" gorm:"type:varchar(10);default:'open'"`
	JobID       uint      `json:"-"`
}

type QEJob struct {
	gorm.Model  `json:"-"`
	DataID      string    `json:"dataId" gorm:"type:varchar(50);unique"`
	Date        time.Time `json:"date" gorm:"type:date"`
	CompanyName string    `json:"companyName" gorm:"type:varchar(255)"`
	JobPosition string    `json:"jobPosition" gorm:"type:varchar(255)"`
	WorkType    string    `json:"workType,omitempty" gorm:"type:varchar(50)"`
	TechStack   string    `json:"techStack" gorm:"type:varchar(255)"`
	LinkToJob   string    `json:"linkToJob" gorm:"type:varchar(255)"`
	Industry    string    `json:"industry,omitempty" gorm:"type:varchar(5)"`
	Status      string    `json:"status" gorm:"type:varchar(10);default:'open'"`
	JobID       uint      `json:"-"`
}

var (
	db *gorm.DB
)

func main() {
	// Membuat koneksi ke database MySQL
	dsn := "root:@tcp(localhost:3306)/job_vacancies_v2?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Migrasi tabel
	db.AutoMigrate(&Job{}, &BEJob{}, &FEJob{}, &QEJob{})

	// Inisialisasi Echo framework
	e := echo.New()

	// Mengatur route dan handler untuk endpoint tertentu
	repo := NewRepo(db)
	service := NewService(repo)

	e.GET("/jobs", getJobs)
	e.GET("/jobs/:id", getJob)
	e.POST("/jobs", service.createJob)

	// Menjalankan server pada port 8000
	e.Start(":8000")
}

type Response struct {
	Status    int         `json:"status"`
	Message   string      `json:"message"`
	TotalData int         `json:"totalData,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

type Information struct {
	Key  string `json:"key"`
	Data string `json:"data"`
}

type repo struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) *repo {
	return &repo{db}
}

type Repo interface {
	dataBE(first, totalData int, data [][]string, author string) error
	dataFE(first, totalData int, data [][]string, author string) error
	dataQE(first, totalData int, data [][]string, author string) error
}

func parseTime(input string) (time.Time, error) {
	var filterDate string

	// checkInput
	if input == "" {
		return time.Now(), nil
	}

	filterDate = strings.ReplaceAll(input, " ", "-")
	filterDate = strings.ReplaceAll(filterDate, "\\", "")

	if len(input) == 10 {
		filterDate = fmt.Sprintf("0%s", filterDate)
	}

	res, err := time.Parse("02-Jan-2006", filterDate)
	if err != nil {
		return time.Now(), err
	}

	return res, nil
}

func generateCode(input string) string {
	canonicID, err := nanoid.Standard(21)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s-%s-%s", input, uuid.New().String(), canonicID())
}

func (r *repo) dataBE(first, totalData int, req [][]string, author string) error {
	var jobs []BEJob

	for i := first; i < totalData; i++ {
		var job BEJob

		date, err := parseTime(req[i][0])
		if err != nil {
			return fmt.Errorf("parse time : %v", err.Error())
		}

		job.DataID = generateCode("BE")
		job.Date = date
		job.CompanyName = req[i][1]
		job.JobPosition = req[i][2]
		job.WorkType = req[i][3]
		job.TechStack = req[i][4]
		job.LinkToJob = req[i][5]

		if len(req[i]) == 10 {
			job.Industry = req[i][9]
		}

		jobs = append(jobs, job)
	}

	data := Job{
		Author:   author,
		Backends: jobs,
	}

	err := r.db.Create(&data).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) dataFE(first, totalData int, req [][]string, author string) error {
	var jobs []FEJob

	for i := first; i < totalData; i++ {
		var job FEJob

		date, err := parseTime(req[i][0])
		if err != nil {
			return fmt.Errorf("parse time : %v", err.Error())
		}

		job.DataID = generateCode("FE")
		job.Date = date
		job.CompanyName = req[i][1]
		job.JobPosition = req[i][2]
		job.WorkType = req[i][3]
		job.TechStack = req[i][4]
		job.LinkToJob = req[i][5]

		if len(req[i]) == 10 {
			job.Industry = req[i][9]
		}

		jobs = append(jobs, job)
	}

	data := Job{
		Author:    author,
		Frontends: jobs,
	}

	err := r.db.Create(&data).Error
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) dataQE(first, totalData int, req [][]string, author string) error {
	var jobs []QEJob

	for i := first; i < totalData; i++ {
		var job QEJob

		date, err := parseTime(req[i][0])
		if err != nil {
			return fmt.Errorf("parse time : %v", err.Error())
		}

		job.DataID = generateCode("QE")
		job.Date = date
		job.CompanyName = req[i][1]
		job.JobPosition = req[i][2]
		job.WorkType = req[i][3]
		job.TechStack = req[i][4]
		job.LinkToJob = req[i][5]

		if len(req[i]) == 10 {
			job.Industry = req[i][9]
		}

		jobs = append(jobs, job)
	}

	data := Job{
		Author:    author,
		Qualities: jobs,
	}

	err := r.db.Create(&data).Error
	if err != nil {
		return err
	}

	return nil
}

type service struct {
	r Repo
}

func NewService(r Repo) *service {
	return &service{r}
}

func (s *service) createJob(c echo.Context) error {
	var req Information

	err := c.Bind(&req)
	if err != nil {
		log.Println(err.Error())

		return c.JSON(http.StatusBadRequest, Response{
			Status:  400,
			Message: "Bad Request. Name and Key is Required",
		})
	}

	var dataStruct interface{}

	switch req.Key {
	case "2023 BE":
		dataStruct = &BEJob{}
	case "2023 FE":
		dataStruct = &FEJob{}
	case "2023 QE":
		dataStruct = &QEJob{}
	default:
		return c.JSON(http.StatusBadRequest, Response{
			Status:  400,
			Message: "Bad Request. Name and Key is Required",
		})
	}

	ctx := context.Background()
	sheetsService, err := sheets.NewService(ctx, option.WithAPIKey("API_KEY"))
	if err != nil {
		log.Println(err.Error())

		return c.JSON(http.StatusInternalServerError, Response{
			Status:  500,
			Message: "Internal Server Error",
		})
	}

	// ID spreadsheet yang ingin diakses
	spreadsheetID := "SHEET_ID"

	// Range dari data yang ingin diambil (misalnya, "Sheet1!A1:E10")
	readRange := fmt.Sprintf("'%s'!A:J", req.Key)

	// Mengambil data dari spreadsheet
	resp, err := sheetsService.Spreadsheets.Values.Get(spreadsheetID, readRange).Do()
	if err != nil {
		log.Println(err.Error())

		return c.JSON(http.StatusInternalServerError, Response{
			Status:  500,
			Message: "Internal Server Error",
		})
	}

	jsonData, err := resp.MarshalJSON()
	if err != nil {
		log.Println(err.Error())

		return c.JSON(http.StatusInternalServerError, Response{
			Status:  500,
			Message: "Internal Server Error",
		})
	}

	var result Data
	err = json.Unmarshal(jsonData, &result)
	if err != nil {
		log.Println(err.Error())

		return c.JSON(http.StatusInternalServerError, Response{
			Status:  500,
			Message: "Internal Server Error",
		})
	}

	// total Input
	var dataNow int64
	db.Find(dataStruct).Count(&dataNow)

	dataNow += 2

	// changeType
	first := int(dataNow)

	totalData := len(result.Data)

	if first == totalData {
		return c.JSON(http.StatusNotFound, Response{
			Status:    404,
			Message:   "Data already Updated",
			TotalData: first,
		})
	}

	switch req.Key {
	case "2023 BE":
		err = s.r.dataBE(first, totalData, result.Data, req.Data)
	case "2023 FE":
		err = s.r.dataFE(first, totalData, result.Data, req.Data)
	case "2023 QE":
		err = s.r.dataQE(first, totalData, result.Data, req.Data)
	}

	if err != nil {
		log.Println(err.Error())

		return c.JSON(http.StatusInternalServerError, Response{
			Status:  500,
			Message: "Internal Server Error",
		})
	}

	return c.JSON(http.StatusCreated, Response{
		Status:    201,
		Message:   "Created. Success Seeder",
		TotalData: totalData,
	})
}

type Data struct {
	MajorDimension string     `json:"majorDimension"`
	Range          string     `json:"range"`
	Data           [][]string `json:"values"`
}

type Model interface {
	Find() []interface{}
}

type Job struct {
	gorm.Model `json:"-"`
	Author     string  `json:"author"`
	Backends   []BEJob `json:"backends"`
	Frontends  []FEJob `json:"frontends"`
	Qualities  []QEJob `json:"qualities"`
}

func getJobs(c echo.Context) error {
	name := c.QueryParam("type")

	var data interface{}
	var errors error

	switch name {
	case "backend":
		data = []BEJob{}

		err := db.Order("date DESC").Find(&data).Error
		if err != nil {
			errors = err
		}
	case "frontend":
		data = []FEJob{}

		err := db.Order("date DESC").Find(&data).Error
		if err != nil {
			errors = err
		}
	case "quality":
		data = []QEJob{}

		err := db.Order("date DESC").Find(&data).Error
		if err != nil {
			errors = err
		}
	default:
		var all []Job

		err := db.Preload("Backends", func(db *gorm.DB) *gorm.DB {
			return db.Order("date DESC")
		}).Preload("Frontends", func(db *gorm.DB) *gorm.DB {
			return db.Order("date DESC")
		}).Preload("Qualities", func(db *gorm.DB) *gorm.DB {
			return db.Order("date DESC")
		}).Find(&all).Error

		if err != nil {
			errors = err
		}

		data = all
	}

	if errors != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Status:  500,
			Message: "Internal Server Error",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Status:  200,
		Message: "OK. Get All Data",
		Data:    data,
	})
}

func getJob(c echo.Context) error {
	id := c.Param("id")
	name := c.QueryParam("type")

	var data interface{}

	switch name {
	case "backend":
		data = &BEJob{}
	case "frontend":
		data = &FEJob{}
	case "quality":
		data = &QEJob{}
	default:
		return c.JSON(http.StatusNotFound, Response{
			Status:  404,
			Message: "Data Not Found",
		})
	}

	err := db.First(&data, "data_id = ?", id).Error
	if err != nil {
		log.Println(err.Error())

		return c.JSON(http.StatusInternalServerError, Response{
			Status:  500,
			Message: "Internal Server Error",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Status:  200,
		Message: "OK. Get Data by Id",
		Data:    data,
	})
}
