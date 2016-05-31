// by setting package as main, Go will compile this as an executable file.
// Any other package turns this into a library
package main

// These are your imports / libraries / frameworks
import (
	// this is Go's built-in sql library
	"database/sql"
	"log"
	"net/http"
	"os"
	"strconv"

	// this allows us to run our web server
	"github.com/gin-gonic/gin"
	// this lets us connect to Postgres DB's
	_ "github.com/lib/pq"
)

var (
	// this is the pointer to the database we will be working with
	// this is a "global" variable (sorta kinda, but you can use it as such)
	db *sql.DB
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("$PORT must be set")
	}

	var errd error
	// here we want to open a connection to the database using an environemnt variable.
	// This isn't the best technique, but it is the simplest one for heroku
	db, errd = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if errd != nil {
		log.Fatalf("Error opening database: %q", errd)
	}
	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("html/*")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/ping", func(c *gin.Context) {
		ping := db.Ping()
		if ping != nil {
			// our site can't handle http status codes, but I'll still put them in cause why not
			c.JSON(http.StatusOK, gin.H{"error": "true", "message": "db was not created. Check your DATABASE_URL"})
		} else {
			c.JSON(http.StatusOK, gin.H{"error": "false", "message": "db created"})
		}
	})

	router.GET("/query1", func(c *gin.Context) {
		table := "<table class='table'><thead><tr>"
		// put your query here
		rows, err := db.Query("SELECT * FROM Patient") // <--- EDIT THIS LINE
		if err != nil {
			// careful about returning errors to the user!
			c.AbortWithError(http.StatusInternalServerError, err)
		}
		// foreach loop over rows.Columns, using value
		cols, _ := rows.Columns()
		if len(cols) == 0 {
			c.AbortWithStatus(http.StatusNoContent)
		}
		for _, value := range cols {
			table += "<th class='text-center'>" + value + "</th>"
		}
		// once you've added all the columns in, close the header
		table += "</thead><tbody>"
		// declare all your RETURNED columns here
		var patient_id int
		var address_id int
		var firstname string
		var lastname string
		var email string 
		var dob string
		
		for rows.Next() {
			// assign each of them, in order, to the parameters of rows.Scan.
			// preface each variable with &
			rows.Scan(&patient_id,&address_id,&firstname,&lastname,&email,&dob) // <--- EDIT THIS LINE
			// can't combine ints and strings in Go. Use strconv.Itoa(int) instead
			table += "<tr><td>" + strconv.Itoa(patient_id) + "</td><td>"+strconv.Itoa(address_id)+"</td><td>"+firstname+"</td><td>"+lastname+"</td><td>"+email+"</td><td>"+dob+"</td></tr>" // <--- EDIT THIS LINE
		}
		// finally, close out the body and table
		table += "</tbody></table>"
		c.Data(http.StatusOK, "text/html", []byte(table))
	})

	// router.GET("/query2", func(c *gin.Context) {
	// 	table := "<table class='table'><thead><tr>"
	// 	// put your query here
	// 	rows, err := db.Query("SELECT Song.title AS title, Artist.firstName AS firstName, Artist.lastName AS lastName FROM Song JOIN Album ON Song.albumId = Album.albumId JOIN Artist ON Album.artistId = Artist.artistId WHERE Album.genre = (SELECT favoriteGenre FROM Person WHERE Person.firstName = 'Katie')") // <--- EDIT THIS LINE
	// 	if err != nil {
	// 		// careful about returning errors to the user!
	// 		c.AbortWithError(http.StatusInternalServerError, err)
	// 	}
	// 	// foreach loop over rows.Columns, using value
	// 	cols, _ := rows.Columns()
	// 	if len(cols) == 0 {
	// 		c.AbortWithStatus(http.StatusNoContent)
	// 	}
	// 	for _, value := range cols {
	// 		table += "<th class='text-center'>" + value + "</th>"
	// 	}
	// 	// once you've added all the columns in, close the header
	// 	table += "</thead><tbody>"
	// 	// columns
	// 	var title string
	// 	var firstName string
	// 	var lastName string
	// 	for rows.Next() {
	// 		rows.Scan(&title,&firstName,&lastName) // put columns here prefaced with &
	// 		table += "<tr><td>"+title+"</td><td>"+firstName+"</td><td>"+lastName+"</td></tr>" // <--- EDIT THIS LINE
	// 	}
	// 	// finally, close out the body and table
	// 	table += "</tbody></table>"
	// 	c.Data(http.StatusOK, "text/html", []byte(table))
	// })

	// router.GET("/query3", func(c *gin.Context) {
	// 	table := "<table class='table'><thead><tr>"
	// 	// put your query here
	// 	rows, err := db.Query("SELECT Song.title AS title, Song.length AS seconds FROM Song JOIN Album ON Song.albumId = Album.albumId WHERE Album.releaseDate > date('2010-01-01') ORDER BY seconds ASC") // <--- EDIT THIS LINE
	// 	if err != nil {
	// 		// careful about returning errors to the user!
	// 		c.AbortWithError(http.StatusInternalServerError, err)
	// 	}
	// 	// foreach loop over rows.Columns, using value
	// 	cols, _ := rows.Columns()
	// 	if len(cols) == 0 {
	// 		c.AbortWithStatus(http.StatusNoContent)
	// 	}
	// 	for _, value := range cols {
	// 		table += "<th class='text-center'>" + value + "</th>"
	// 	}
	// 	// once you've added all the columns in, close the header
	// 	table += "</thead><tbody>"
	// 	// columns
	// 	var title string
	// 	var seconds int
	// 	for rows.Next() {
	// 		rows.Scan(&title,&seconds) // put columns here prefaced with &
	// 		table += "<tr><td>"+title+"</td><td>"+strconv.Itoa(seconds)+"</td><</tr>" // <--- EDIT THIS LINE
	// 	}
	// 	// finally, close out the body and table
	// 	table += "</tbody></table>"
	// 	c.Data(http.StatusOK, "text/html", []byte(table))
	// })

	// NO code should go after this line. it won't ever reach that point
	router.Run(":" + port)
}

/*
Example of processing a GET request

// this will run whenever someone goes to last-first-lab7.herokuapp.com/EXAMPLE
router.GET("/EXAMPLE", func(c *gin.Context) {
    // process stuff
    // run queries
    // do math
    //decide what to return
    c.JSON(http.StatusOK, gin.H{
        "key": "value"
        }) // this returns a JSON file to the requestor
    // look at https://godoc.org/github.com/gin-gonic/gin to find other return types. JSON will be the most useful for this
})

*/
