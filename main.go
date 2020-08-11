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

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

const (
	hostname     = "localhost"
	hostport     = 5432
	username     = "postgres"
	password     = "root"
	databasename = "ImageStore"
	uploadPath   = "./images"
)

type Album struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type ImageStuct struct {
	Name    string `json:"name"`
	ID      string `json:"id"`
	AlbumID string `json:"albumid"`
}

type PushImageDB struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
	AlbumID   string `json:"albumid"`
	Imagepath string `json:"imagepath"`
}

func startupServer() {
	r := mux.NewRouter()
	r.HandleFunc("/", helloServer)
	r.HandleFunc("/createAlbum", createAlbum).Methods("POST")
	r.HandleFunc("/deleteAlbum", deleteAlbum).Methods("DELETE")
	r.HandleFunc("/createImage", createImage).Methods("POST")
	r.HandleFunc("/deleteImage", deleteimage).Methods("DELETE")
	r.HandleFunc("/getImage", getImage).Methods("GET")
	log.Fatal(http.ListenAndServe("localhost:8000", r))

}

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

func helloServer(w http.ResponseWriter, r *http.Request) {

	//w.WriteHeader(http.StatusInternalServerError)
	log.Println("Server said Hello")
	//fmt.Fprintf(w, "Hello")
	return
}

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
	var getimage ImageStuct
	getimage.ID = imageid[0]
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
		w.WriteHeader(200)
		return
	}

}

func gettingImage(getimage ImageStuct) (string, error) {
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
	err = db.QueryRow(statement, getimage.ID, getimage.AlbumID).Scan(&imagepath)
	if err != nil {
		return "", err
	}

	return imagepath, err
}

func createAlbum(w http.ResponseWriter, r *http.Request) {

	var newAlbum Album

	err := json.NewDecoder(r.Body).Decode(&newAlbum)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		log.Println(err)
	}

	log.Println("Album ID is : ", newAlbum.ID)
	log.Println("Album Name is : ", newAlbum.Name)

	check, err := creatingAlbum(newAlbum)
	if check == false {
		log.Println("Album creation failed")
		log.Println(err)
		WriteJSONResponse(w, 403, "Cannot have multiple albums with same ID")
		return
	} else {
		log.Printf("Album %s created", newAlbum.Name)
		err = os.Mkdir(uploadPath+"/"+newAlbum.ID, 755)
		if err != nil {
			log.Println(err)
		}

		WriteJSONResponse(w, 200, "Successfully Album Created")
		//fmt.Fprintf(w, "Successfully Album Created\n")
		return
	}

}

func creatingAlbum(newAlbum Album) (bool, error) {

	conStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", hostname, hostport, username, password, databasename)

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		log.Println(err)
		return false, err
	} else {
		log.Println("Database connected")
	}

	defer db.Close()

	statement := `INSERT into public.album(albumid, albumname) VALUES($1, $2) RETURNING albumname`
	name := ""
	log.Println(newAlbum.ID, newAlbum.Name)
	err = db.QueryRow(statement, newAlbum.ID, newAlbum.Name).Scan(&name)
	if err != nil {
		//log.Println(err)
		return false, err
	}
	log.Println("New record is : ", name)

	return true, nil
}

func deleteAlbum(w http.ResponseWriter, r *http.Request) {

	var DelAlbum Album
	err := json.NewDecoder(r.Body).Decode(&DelAlbum)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		log.Println(err)
	}

	log.Println("Album ID needs to be deleted : ", DelAlbum.ID)
	log.Println("Album Name needs to be deleted : ", DelAlbum.Name)

	check, err := DeletingAlbum(DelAlbum)
	if check == false {
		log.Println(err)
		WriteJSONResponse(w, 403, "No Album Found")
		return
	} else {
		err = os.RemoveAll(uploadPath + "/" + DelAlbum.ID)
		if err != nil {
			log.Println(err)
		}
		log.Printf("Album %s Deleted", DelAlbum.Name)
		WriteJSONResponse(w, 200, "Album deletion Successfull")
		return
	}

}

func DeletingAlbum(DelAlbum Album) (bool, error) {

	conStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", hostname, hostport, username, password, databasename)

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		//log.Println(err)
		return false, err
	} else {
		log.Println("Database connected")
	}

	defer db.Close()

	statement := `Delete FROM public.album WHERE albumid = $1 RETURNING albumname`
	name := ""
	log.Println(DelAlbum.ID, DelAlbum.Name)
	err = db.QueryRow(statement, DelAlbum.ID).Scan(&name)
	if err != nil {
		log.Println(err)
		return false, err
	}
	log.Println("Deleted Album is : ", name)
	return true, nil

}

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
		WriteJSONResponse(w, 403, "Error Retrieving the File")
		return
	}
	fileBytes, err := ioutil.ReadAll(file)

	filetype := http.DetectContentType(fileBytes)
	if filetype != "image/jpeg" && filetype != "image/jpg" &&
		filetype != "image/gif" && filetype != "image/png" {
		WriteJSONResponse(w, 401, http.StatusBadRequest)
		return
	}
	//log.Println(filetype[0])
	var creatImage ImageStuct

	err = json.Unmarshal([]byte(r.FormValue("data")), &creatImage)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	checkAlbumID, err := CheckAlbumID(creatImage)

	if checkAlbumID == false {
		log.Println("Failed")
		if err != nil {
			log.Println(err)
		}
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

		if err != nil {
			fmt.Println(err)
		}
		defer tempFile.Close()
		filestat, err := os.Stat(tempFile.Name())
		creatImage.Name = filestat.Name()
		// read all of the contents of our uploaded file into a
		// byte array
		var pushImageDB PushImageDB

		pushImageDB.Name = creatImage.Name
		pushImageDB.ID = creatImage.ID
		pushImageDB.AlbumID = creatImage.AlbumID
		pushImageDB.Imagepath = uploadPath + "/" + creatImage.AlbumID + "/" + creatImage.Name

		fmt.Println(pushImageDB.Imagepath)
		check, err := CreatingImage(pushImageDB)
		if check == false {
			fmt.Println("Error while saving in db")
			WriteJSONResponse(w, 403, "Cannot have multiple images with same ID")
			tempFile.Close()
			err := os.Remove(pushImageDB.Imagepath)
			if err != nil {
				log.Println(err)
			}
			return
		}
		log.Println("Image ID is : ", creatImage.ID)
		log.Println("Image Name is : ", creatImage.Name)
		log.Println("Image AlbumID is : ", creatImage.AlbumID)

		log.Printf("Image %s saved in db", pushImageDB.Name)
		// write this byte array to our file
		tempFile.Write(fileBytes)
		// return that we have successfully uploaded our file!
		//fmt.Fprintf(w, "Successfully Uploaded File\n")
		WriteJSONResponse(w, 200, "Successfully Uploaded File")
	}

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	return
}

func CheckAlbumID(creatImage ImageStuct) (bool, error) {
	conStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", hostname, hostport, username, password, databasename)

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		log.Println(err)
		return false, err
	} else {
		log.Println("Database connected here")
	}

	statement1 := `SELECT albumid FROM public.album where albumid = $1`
	id := ""
	err = db.QueryRow(statement1, creatImage.AlbumID).Scan(&id)
	if err != nil {
		log.Println(err)
		log.Println("AlbumID does not exist")
		return false, err
	}
	return true, nil
}

func CreatingImage(pushImageDB PushImageDB) (bool, error) {

	conStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", hostname, hostport, username, password, databasename)

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		log.Println(err)
		return false, err
	} else {
		log.Println("Database connected")
	}

	defer db.Close()

	statement := `INSERT into public.image(imageid, imagename, albumid, imagepath) VALUES($1, $2, $3, $4) RETURNING imageid`
	id := ""
	log.Println(pushImageDB.ID, pushImageDB.Name, pushImageDB.AlbumID, pushImageDB.Imagepath)
	err = db.QueryRow(statement, pushImageDB.ID, pushImageDB.Name, pushImageDB.AlbumID, pushImageDB.Imagepath).Scan(&id)
	if err != nil {
		log.Println(err)
		return false, err
	}

	log.Println("New image saved: ")
	return true, nil

}

func deleteimage(w http.ResponseWriter, r *http.Request) {

	var DelImage ImageStuct
	err := json.NewDecoder(r.Body).Decode(&DelImage)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		log.Println(err)
	}

	log.Println("Image ID needs to be deleted : ", DelImage.ID)
	log.Println("Image Name needs to be deleted : ", DelImage.Name)
	log.Println("Image needs to be delted from Album ", DelImage.AlbumID)

	imageName, err := DeletingImage(DelImage)
	if err != nil {
		log.Println(err)
		WriteJSONResponse(w, 403, "No image found")
		return
	} else {
		err = os.Remove(uploadPath + "/" + DelImage.AlbumID + "/" + imageName)
		if err != nil {
			log.Println(err)
		}
		log.Printf("Image %s Deleted", imageName)
		WriteJSONResponse(w, 200, "Successfully deleted the image")
		return
	}

}

func DeletingImage(DelImage ImageStuct) (string, error) {
	conStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", hostname, hostport, username, password, databasename)

	db, err := sql.Open("postgres", conStr)
	if err != nil {
		log.Println(err)
		return "Database connection Failed", err
	} else {
		log.Println("Database connected")
	}

	defer db.Close()

	statement := `Delete FROM public.image WHERE imageid = $1 RETURNING imagename`
	name := ""
	log.Println(DelImage.ID, DelImage.Name)
	err = db.QueryRow(statement, DelImage.ID).Scan(&name)
	if err != nil {
		log.Println(err)
		return "Query Failed", err
	}
	log.Println("Deleted Image is : ", name)
	return name, nil
}
