// ImageStore Service API
//
// .
//
//     Schemes: http, https
//     Host: localhost:8000
//     Version: 0.1.0
//     basePath: /
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
// swagger:meta
package main

import (
	"database/sql"
	_ "database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

//constants can be set as env variables as well. I am hardcoding it for now

const (
	hostname     = "localhost"
	hostport     = 5432
	username     = "postgres"
	password     = "root"
	databasename = "ImageStore"
	uploadPath   = "./images"
)

//startupServer will start the server with all the routes
func startupServer() {
	r := mux.NewRouter()
	staticpath := strings.TrimPrefix(uploadPath, ".")
	staticpath = staticpath + "/"
	log.Println(staticpath)
	r.PathPrefix(staticpath).Handler(http.StripPrefix(staticpath, http.FileServer(http.Dir(uploadPath))))
	//hello server
	// swagger:operation Get / helloServer
	//
	// ---
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: successful operation
	r.HandleFunc("/", helloServer)
	//CreateAlbum
	// swagger:operation POST /createAlbum createAlbum
	//
	// ---
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: successful operation
	r.HandleFunc("/createAlbum", createAlbum).Methods("POST")
	//DeleteAlbum
	// swagger:operation DELETE /deleteAlbum DeleteAlbum
	//
	// ---
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: successful operation
	r.HandleFunc("/deleteAlbum", deleteAlbum).Methods("DELETE")
	//CreateImage
	// swagger:operation POST /createImage createImage
	//
	// ---
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: successful operation
	r.HandleFunc("/createImage", createImage).Methods("POST")
	//DeleteImage
	// swagger:operation DELETE /deleteImage deleteImage
	//
	// ---
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: successful operation
	r.HandleFunc("/deleteImage", deleteimage).Methods("DELETE")

	//GetImage
	// swagger:operation GET /getImage getImage
	//
	// ---
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: successful operation
	r.HandleFunc("/getImage", getImage).Methods("GET")
	//GetAlbumImages

	// swagger:operation GET /getAlbumImage AlbumImage
	//
	// ---
	// produces:
	// - application/json
	// responses:
	//   '200':
	//     description: successful operation
	r.HandleFunc("/getAlbumImage", getAlbumImage).Methods("GET")

	log.Fatal(http.ListenAndServe("localhost:8000", r))

}

//Album represents the structure of album
//Album Json request payload is as follows,
//{
//  "Name": "image1"
//  "ID": "1",
//}
type Album struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

//CreateAlbum for creation of album
type CreateAlbum struct {
	Name string `json:"name"`
}

//ImageStruct represents the structure of Image
//ImageStruct Json request payload is as follows,
//{
//  "Name": "image1"
//  "ID": "1",
//	"AlbumID": "123"
//}
type ImageStruct struct {
	AlbumID string `json:"albumid"`
}

//PushImageDB represents the structure of image while pushing it into database
//PushImageDB will have values similar to as follows,
//{
//  "Name": "image1"
//  "ID": "1",
//	"AlbumID": "123"
//	"Imagepath":"./images/123/{filename}"
//}
type PushImageDB struct {
	Name      string `json:"name"`
	AlbumID   string `json:"albumid"`
	Imagepath string `json:"imagepath"`
}

// @title ImageStore service api
// @version 1.0
// @description This api is written to create/delete images and albums.
// @termsOfService http://swagger.io/terms/
// @contact.name Shikhar Kannoje
// @contact.email shikharkannoje@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8000
// @BasePath /

//main is calling startupServer
func main() {

	fmt.Println("The server is up and running")
	startupServer()
}

// WriteJSONResponse represents a utility function which writes status code and JSON to response
func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

//Basic Hello Server message
func helloServer(w http.ResponseWriter, r *http.Request) {

	//w.WriteHeader(http.StatusInternalServerError)
	log.Println("Server said Hello")
	fmt.Fprintf(w, "Hello")
	return
}

// GetImage godoc
// @Summary Get an image directly as a response
// @Description Get an image by providing the imageID and the albumID of the image.
// @Tags images
// @Accept  json
// @Produce  image
// @Success 200 image
// @Router /getImage [get]

//GetImageStruct used for getting an image
type GetImageStruct struct {
	ImageID string `json:"imageid"`
	AlbumID string `json:"albumid"`
}

//getImage method is called to get an individual image using image ID
func getImage(w http.ResponseWriter, r *http.Request) {

	imageid, ok := r.URL.Query()["imageid"]
	if !ok || len(imageid[0]) < 1 {
		log.Println("imageid is missing")
		return
	}
	albumid, ok := r.URL.Query()["albumid"]
	if !ok || len(albumid[0]) < 1 {
		log.Println("imageid is missing")
		return
	}

	log.Println(imageid[0], albumid[0])
	var getimage GetImageStruct
	getimage.ImageID = imageid[0]
	getimage.AlbumID = albumid[0]
	imagepath, err := gettingImage(getimage)
	if err != nil {
		w.WriteHeader(403)
		WriteJSONResponse(w, 403, "Invalid image or albumid")
		fmt.Println("Invalid image or albumid")
		fmt.Println(err)
		return
	} else {
		filebytes, err := ioutil.ReadFile(imagepath)
		if err != nil {
			log.Println("Error while reading the file")
		}
		_, err = w.Write(filebytes)
		if err != nil {
			log.Println("Error while writing the file")
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		return
	}

}

//getImage is calling this gettingImage internally to handle DB operations
func gettingImage(getimage GetImageStruct) (string, error) {
	conStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", hostname, hostport, username, password, databasename)

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		//log.Println(err)
		return "", err
	} else {
		log.Println("Database connected")
	}
	defer db.Close()

	statement := `SELECT imagepath FROM public.image WHERE imageid = $1 AND albumid = $2`
	imagepath := ""
	imageintID, err := strconv.Atoi(getimage.ImageID)
	albumintID, err := strconv.Atoi(getimage.AlbumID)
	err = db.QueryRow(statement, imageintID, albumintID).Scan(&imagepath)
	if err != nil {
		return "", err
	}

	return imagepath, err
}

// GetAlbumImage godoc
// @Summary Get all the images of an album.
// @Description Get all the images by providing the albumID. Images will be statically served.
// @Tags images
// @Accept  json
// @Produce  image
// @Param albumid
// @Success static served images
// @Router /getAlbumImage [get]

//getAlbumImage is called to get all the images of an album, this method is only redirecting the request
//to correct file location where the files are being served.
func getAlbumImage(w http.ResponseWriter, r *http.Request) {

	albumid, ok := r.URL.Query()["albumid"]
	if !ok || len(albumid[0]) < 1 {
		log.Println("Album is missing")
		return
	}
	conStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", hostname, hostport, username, password, databasename)

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Database connected")
	}
	defer db.Close()

	albumidInt, err := strconv.Atoi(albumid[0])
	statement := `SELECT albumid FROM public.album WHERE albumid = $1`
	var id int
	err = db.QueryRow(statement, albumidInt).Scan(&id)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Album does not exist"))
		WriteJSONResponse(w, 403, http.StatusForbidden)
		return
	} else {
		url := uploadPath + "/" + albumid[0] + "/"
		log.Println(url)
		http.Redirect(w, r, url, http.StatusSeeOther)
	}
	return
}

// CreateAlbum godoc
// @Summary Create a new album
// @Description Create a new album with the album id and album name
// @Tags Album
// @Accept  json
// @Produce  json
// @Success 200
// @Router /createAlbum [post]
//createAlbum is called when album creation request is made.
func createAlbum(w http.ResponseWriter, r *http.Request) {

	var newAlbum CreateAlbum

	err := json.NewDecoder(r.Body).Decode(&newAlbum)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		log.Println(err)
	}

	//log.Println("Album ID is : ", newAlbum.ID)
	log.Println("Album Name is : ", newAlbum.Name)

	album, err := creatingAlbum(newAlbum)
	if err != nil {
		log.Println("Album creation failed")
		log.Println(err)
		WriteJSONResponse(w, 403, "Cannot have multiple albums with same ID")
		return
	} else {
		log.Printf("Album %s created", album.Name)
		err = os.Mkdir(uploadPath+"/"+string(album.ID), 755)
		if err != nil {
			log.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("AlbumID " + string(album.ID) + " created"))
		WriteJSONResponse(w, 200, "Successfull")
		//fmt.Fprintf(w, "Successfully Album Created\n")
		return
	}

}

//creatingAlbum is being called from createAlbum to handle the DB operation
func creatingAlbum(newAlbum CreateAlbum) (Album, error) {

	conStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", hostname, hostport, username, password, databasename)
	var album Album
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		log.Println(err)
		return album, err
	} else {
		log.Println("Database connected")
	}

	defer db.Close()

	statement := `INSERT into public.album(albumname) VALUES($1) RETURNING albumname, albumid`
	name := ""
	var id int
	log.Println(newAlbum.Name)
	err = db.QueryRow(statement, newAlbum.Name).Scan(&name, &id)
	if err != nil {
		//log.Println(err)
		return album, err
	}
	log.Println("New Album is : ", name)

	album.ID = strconv.Itoa(id)
	album.Name = name
	return album, nil
}

// DeleteAlbum godoc
// @Summary Delete an existing album
// @Description Delete an existing album by the album id
// @Tags Album
// @Accept  json
// @Produce  json
// @Param albumID albumName
// @Success 200 {object} Order
// @Router /deleteAlbum [post]

type deleteAlbumStruct struct {
	AlbumID string `json:"albumid"`
}

//deleteAlbum is called to delete an album by providing the albumid.
func deleteAlbum(w http.ResponseWriter, r *http.Request) {

	var DelAlbum deleteAlbumStruct
	err := json.NewDecoder(r.Body).Decode(&DelAlbum)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		log.Println(err)
	}

	log.Println("Album ID needs to be deleted : ", DelAlbum)

	delalbum, err := deletingAlbum(DelAlbum)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		WriteJSONResponse(w, 403, "No Album Found with the given id")
		return
	} else {
		err = os.RemoveAll(uploadPath + "/" + delalbum.ID)
		if err != nil {
			log.Println(err)
		}
		log.Printf("Album %s Deleted", delalbum.Name)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("Album " + delalbum.Name + " is deleted"))
		WriteJSONResponse(w, 200, "Album deletion Successfull")
		return
	}

}

//deletingAlbum is being called from deleteAlbum to handle DB operation
func deletingAlbum(DelAlbum deleteAlbumStruct) (Album, error) {

	conStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", hostname, hostport, username, password, databasename)
	var delalbum Album
	db, err := sql.Open("postgres", conStr)
	if err != nil {
		//log.Println(err)
		return delalbum, err
	} else {
		log.Println("Database connected")
	}

	defer db.Close()

	statement := `Delete FROM public.album WHERE albumid = $1 RETURNING albumname, albumid`
	name := ""
	var id int
	albumIDInt, err := strconv.Atoi(DelAlbum.AlbumID)
	err = db.QueryRow(statement, albumIDInt).Scan(&name, &id)
	if err != nil {
		log.Println(err)
		return delalbum, err
	}
	log.Println("Deleted Album is : ", name)
	delalbum.ID = strconv.Itoa(id)
	delalbum.Name = name
	return delalbum, nil

}

// CreateImage godoc
// @Summary Create/upload a new image
// @Description Create a new image with the data(imageID, imageName, albumid) imagepath(file) in payload.
// @Tags Album
// @Accept  json
// @Produce  json
// @Success 200
// @Router /createImage [post]
//createImage is called when a new image is uploaded on a specific albumid
func createImage(w http.ResponseWriter, r *http.Request) {

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		fmt.Printf("Could not parse multipart form: %v\n", err)
		WriteJSONResponse(w, 500, http.StatusInternalServerError)
		return
	}
	// FormFile returns the first file for the given key `imagepath`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file

	file, handler, err := r.FormFile("imagepath")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		w.Header().Set("Content-Type", "application/json")
		WriteJSONResponse(w, 403, "Error Retrieving the File")
		return
	}
	fileBytes, err := ioutil.ReadAll(file)

	filetype := http.DetectContentType(fileBytes)
	if filetype != "image/jpeg" && filetype != "image/jpg" &&
		filetype != "image/gif" && filetype != "image/png" {
		w.Header().Set("Content-Type", "application/json")
		WriteJSONResponse(w, 401, http.StatusBadRequest)
		return
	}
	//log.Println(filetype[0])
	var creatImage ImageStruct

	err = json.Unmarshal([]byte(r.FormValue("data")), &creatImage)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	checkAlbumID, err := checkAlbumIDMeth(creatImage)

	if checkAlbumID == false {
		log.Println("Failed")
		if err != nil {
			log.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		WriteJSONResponse(w, 403, "Album does not exist")
		return
	} else {
		var ext string
		switch filetype {
		case "image/jpeg":
			ext = "jpeg"
		case "image/jpg":
			ext = "jpg"
		case "image/gif":
			ext = "gif"
		case "image/png":
			ext = "png"

		}
		tempFile, err := ioutil.TempFile(uploadPath+"/"+creatImage.AlbumID, "upload-*."+ext)
		log.Println("Name of tempfile", tempFile.Name())
		if err != nil {
			fmt.Println(err)
		}
		defer tempFile.Close()
		filestat, err := os.Stat(tempFile.Name())
		log.Println("Name of filestat", filestat.Name())

		// read all of the contents of our uploaded file into a
		// byte array
		var pushImageDB PushImageDB

		pushImageDB.Name = filestat.Name()
		pushImageDB.AlbumID = creatImage.AlbumID
		pushImageDB.Imagepath = uploadPath + "/" + creatImage.AlbumID + "/" + filestat.Name()

		fmt.Println(pushImageDB.Imagepath)
		imageidreturn, err := creatingImage(pushImageDB)
		if err != nil {
			fmt.Println("Error while saving in db")
			WriteJSONResponse(w, 403, "Cannot have multiple images with same ID")
			tempFile.Close()
			err := os.Remove(pushImageDB.Imagepath)
			if err != nil {
				log.Println(err)
			}
			return
		}
		// log.Println("Image ID is : ", creatImage.ID)
		// log.Println("Image Name is : ", creatImage.Name)
		// log.Println("Image AlbumID is : ", creatImage.AlbumID)

		log.Printf("Image %s saved in db", pushImageDB.Name)
		log.Printf("ImageID %s", imageidreturn)
		// write this byte array to our file
		tempFile.Write(fileBytes)
		// return that we have successfully uploaded our file!
		//fmt.Fprintf(w, "Successfully Uploaded File\n")
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("ImageID " + imageidreturn + "uploaded"))
		WriteJSONResponse(w, 200, "Successfully Uploaded File")
	}

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	return
}

//checkAlbumIDMeth is called to check whether the albumID exists or not.
func checkAlbumIDMeth(creatImage ImageStruct) (bool, error) {
	conStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", hostname, hostport, username, password, databasename)

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		log.Println(err)
		return false, err
	} else {
		log.Println("Database connected here")
	}

	statement1 := `SELECT albumid FROM public.album where albumid = $1`
	var id int
	albumidint, err := strconv.Atoi(creatImage.AlbumID)
	err = db.QueryRow(statement1, albumidint).Scan(&id)
	if err != nil {
		log.Println(err)
		log.Println("AlbumID does not exist")
		return false, err
	}
	return true, nil
}

//creatingImage is called to handle the DB operation of image creation
func creatingImage(pushImageDB PushImageDB) (string, error) {

	conStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", hostname, hostport, username, password, databasename)

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		log.Println(err)
		return "", err
	} else {
		log.Println("Database connected")
	}

	defer db.Close()
	albumidInt, err := strconv.Atoi(pushImageDB.AlbumID)
	statement := `INSERT into public.image(imagename, albumid, imagepath) VALUES($1, $2, $3) RETURNING imageid`
	var id int
	log.Println(albumidInt, pushImageDB.Imagepath)
	err = db.QueryRow(statement, pushImageDB.Name, albumidInt, pushImageDB.Imagepath).Scan(&id)
	if err != nil {
		log.Println(err)
		return "", err
	}
	rid := strconv.Itoa(id)
	log.Println("New image saved: ")
	return rid, nil

}

// DeleteImage godoc
// @Summary Delete an existing image
// @Description Delete an existing image by imageID, imageName, albumid) in payload.
// @Tags Album
// @Accept  json
// @Produce  json
// @Param data(imageID, imageName, albumid)
// @Success 200
// @Router /deleteImage [post]

type DelImageStruct struct {
	ImageID string `json:"imageid"`
	AlbumID string `json:"albumid"`
}

//deleteimage is called to delete a specific image, by given imageid
func deleteimage(w http.ResponseWriter, r *http.Request) {

	var DelImage DelImageStruct
	err := json.NewDecoder(r.Body).Decode(&DelImage)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		log.Println(err)
	}

	log.Println("Image ID needs to be deleted : ", DelImage.ImageID)
	log.Println("Image needs to be delted from Album ", DelImage.AlbumID)

	imageName, err := deletingImage(DelImage)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		WriteJSONResponse(w, 403, "No image found with the given imageID")
		return
	} else {
		err = os.Remove(uploadPath + "/" + DelImage.AlbumID + "/" + imageName)
		if err != nil {
			log.Println(err)
		}
		log.Printf("Image %s Deleted", imageName)
		w.Header().Set("Content-Type", "application/json")
		WriteJSONResponse(w, 200, "Successfully deleted the image")
		return
	}

}

//deletingImage is called to handle the DB operation
func deletingImage(DelImage DelImageStruct) (string, error) {
	conStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", hostname, hostport, username, password, databasename)

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		log.Println(err)
		return "Database connection Failed", err
	} else {
		log.Println("Database connected")
	}

	defer db.Close()

	statement := `Delete FROM public.image WHERE imageid = $1 AND albumid = $2 RETURNING imagename`
	name := ""
	imageidInt, err := strconv.Atoi(DelImage.ImageID)
	albumidInt, err := strconv.Atoi(DelImage.AlbumID)
	log.Println(DelImage.ImageID, DelImage.AlbumID)
	err = db.QueryRow(statement, imageidInt, albumidInt).Scan(&name)
	if err != nil {
		log.Println(err)
		return "Query Failed", err
	}
	log.Println("Deleted Image is : ", name)
	return name, nil
}
