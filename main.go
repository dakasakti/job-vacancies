package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/dakasakti/job-vacancies/config"
	"github.com/dakasakti/job-vacancies/entities"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xuri/excelize/v2"
)

func main() {
	db := config.Database()
	config.AutoMigrate(db)

	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())

	e.GET("/data/create", func(c echo.Context) error {
		url := config.GetConfig().URLFile

		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()

		file, err := os.Create("export.xlsx")
		if err != nil {
			return c.JSON(500, echo.Map{
				"message": "failed export file",
			})
		}

		size, _ := io.Copy(file, resp.Body)

		defer file.Close()

		return c.JSON(200, echo.Map{
			"message": "berhasil membuat file",
			"data":    fmt.Sprintf("%d Kb", size),
		})
	})

	e.GET("/data/update", func(c echo.Context) error {
		f, err := excelize.OpenFile("export.xlsx")
		if err != nil {
			return c.JSON(500, echo.Map{
				"message": "failed get File excel",
			})
		}

		defer func() {
			if err := f.Close(); err != nil {
				fmt.Println("gagal close sheet")
			}
		}()

		BES, err := f.GetRows("Back-end")
		if err != nil {
			fmt.Println(err)
			return c.JSON(500, echo.Map{
				"message": "failed get sheets Back-end",
			})
		}

		var waktuID uint
		var counter uint = config.GetConfig().BE_Index

		for i := int(counter); i < len(BES); i++ {
			counter++

			if len(BES[i]) > 1 {
				// setData
				data := entities.BackEnd{
					ID:            counter,
					CompanyName:   BES[i][0],
					JobPosition:   BES[i][1],
					WorkType:      BES[i][2],
					TechStack:     BES[i][3],
					LinkToJob:     BES[i][4],
					TimeBackendID: waktuID,
				}

				tx := db.Save(&data)
				if tx.Error != nil {
					fmt.Println("ERROR : ", tx.Error.Error())
				}
			} else {
				// setId
				waktuID = counter

				// setData
				data := entities.TimeBackend{
					ID:   counter,
					Name: BES[i][0],
				}

				tx := db.Save(&data)
				if tx.Error != nil {
					fmt.Println("ERROR : ", tx.Error.Error())
				}
			}

		}

		// resetUlang
		counter = config.GetConfig().FE_Index

		FES, err := f.GetRows("Front-end")
		if err != nil {
			fmt.Println(err)
			return c.JSON(500, echo.Map{
				"message": "failed get sheets Front-end",
			})
		}

		for i := int(counter); i < len(FES); i++ {
			counter++

			if len(FES[i]) > 1 {
				// setData
				data := entities.FrontEnd{
					ID:             counter,
					CompanyName:    FES[i][0],
					JobPosition:    FES[i][1],
					WorkType:       FES[i][2],
					TechStack:      FES[i][3],
					LinkToJob:      FES[i][4],
					TimeFrontendID: waktuID,
				}

				tx := db.Save(&data)
				if tx.Error != nil {
					fmt.Println("ERROR : ", tx.Error.Error())
				}
			} else {
				// setId
				waktuID = counter

				// setData
				data := entities.TimeFrontend{
					ID:   counter,
					Name: FES[i][0],
				}

				tx := db.Save(&data)
				if tx.Error != nil {
					fmt.Println("ERROR : ", tx.Error.Error())
				}
			}

		}

		// resetUlang
		counter = config.GetConfig().QA_Index

		QAS, err := f.GetRows("Quality")
		if err != nil {
			fmt.Println(err)
			return c.JSON(500, echo.Map{
				"message": "failed get sheets Quality",
			})
		}

		for i := int(counter); i < len(QAS); i++ {
			counter++

			if len(QAS[i]) > 1 {
				// setData
				data := entities.Quality{
					ID:            counter,
					CompanyName:   QAS[i][0],
					JobPosition:   QAS[i][1],
					WorkType:      QAS[i][2],
					TechStack:     QAS[i][3],
					LinkToJob:     QAS[i][4],
					TimeQualityID: waktuID,
				}

				tx := db.Save(&data)
				if tx.Error != nil {
					fmt.Println("ERROR : ", tx.Error.Error())
				}
			} else {
				// setId
				waktuID = counter

				// setData
				data := entities.TimeQuality{
					ID:   counter,
					Name: QAS[i][0],
				}

				tx := db.Save(&data)
				if tx.Error != nil {
					fmt.Println("ERROR : ", tx.Error.Error())
				}
			}

		}

		return c.JSON(200, echo.Map{
			"message": "data berhasil diupdate",
		})
	})

	e.GET("/data/back-end/last", func(c echo.Context) error {
		var result entities.TimeBackend

		db.Preload("BackEnds").Last(&result)

		return c.JSON(200, echo.Map{
			"data": result,
		})
	})

	e.GET("/data/front-end/last", func(c echo.Context) error {
		var result entities.TimeFrontend

		db.Preload("FrontEnds").Last(&result)

		return c.JSON(200, echo.Map{
			"data": result,
		})
	})

	e.GET("/data/quality/last", func(c echo.Context) error {
		var result entities.TimeQuality

		db.Preload("Qualitys").Last(&result)

		return c.JSON(200, echo.Map{
			"data": result,
		})
	})

	e.GET("/data/front-end", func(c echo.Context) error {
		ql := c.QueryParam("limit")
		setLimit, err := strconv.Atoi(ql)
		if err != nil {
			return c.Redirect(301, config.PageLimit("front-end"))
		}

		qp := c.QueryParam("page")
		setPage, err := strconv.Atoi(qp)
		if err != nil {
			return c.Redirect(301, config.PageLimit("front-end"))
		}

		var results []entities.FrontEnd

		Offset := (setPage - 1) * setLimit
		db.Limit(setLimit).Offset(Offset).Order("id desc").Find(&results)

		return c.JSON(200, echo.Map{
			"data": results,
		})
	})

	e.GET("/data/back-end", func(c echo.Context) error {
		ql := c.QueryParam("limit")
		setLimit, err := strconv.Atoi(ql)
		if err != nil {
			return c.Redirect(301, config.PageLimit("back-end"))
		}

		qp := c.QueryParam("page")
		setPage, err := strconv.Atoi(qp)
		if err != nil {
			return c.Redirect(301, config.PageLimit("back-end"))
		}

		var results []entities.BackEnd

		Offset := (setPage - 1) * setLimit
		db.Limit(setLimit).Offset(Offset).Order("id desc").Find(&results)

		return c.JSON(200, echo.Map{
			"data": results,
		})
	})

	e.GET("/data/quality", func(c echo.Context) error {
		ql := c.QueryParam("limit")
		setLimit, err := strconv.Atoi(ql)
		if err != nil {
			return c.Redirect(301, config.PageLimit("quality"))
		}

		qp := c.QueryParam("page")
		setPage, err := strconv.Atoi(qp)
		if err != nil {
			return c.Redirect(301, config.PageLimit("quality"))
		}

		var results []entities.Quality

		Offset := (setPage - 1) * setLimit
		db.Limit(setLimit).Offset(Offset).Order("id desc").Find(&results)

		return c.JSON(200, echo.Map{
			"data": results,
		})
	})

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.GetConfig().Port)))
}
