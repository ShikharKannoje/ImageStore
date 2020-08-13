"# ImageStore" 


ImageStore service is based on Golang and PostgreSQL. I have used kafka for producing notification on a topic
everytime when a image or album is created to deleted. 

To setup the application, follow the below procedures

NOTE: I am assuming you have the golang environment setup with a text edior of your choice.
    I am also assuming that you have PostgreSQL setup too on the system

1. Place the the code on the go workspace folder, (go/src/)
2. Create a Database on postgreSQL by the name ImageStore on the running machine. I have attached the SQL scripts to create the tables. Basically there are 2 tables, one for keeping the record of album ID and album name having albumID as the primary key. Second table has the record of imageID, imageName, albumID (foreingkey) and the location of the image file.
3. Once you created the DB and Go project setup, you can start running the project.
4. Refer the swagger documentation to understand the API requests.


Key elements of the service
1. Every Image needs to be attached to a album.
2. Every image should have unique id.
3. Name of the image file can be common as the service will change the image file name.
4. A new album needs to be created first to insert image in that album.
5. Deleting a single image will delete the image.
6. Deleting the album will remove the complete album including all the images associated to that album.
7. You can get a single image from a get request which will basically send the image in the response.
8. If the whole album needs to be fetched then I have served them statically, every image will be served on that location.


## Development

To make changes to the cc-iot-user-management-service, clone the source code and
download all the required development dependencies:



These are the environement variables, kindly set the environment variable first to run the serices

    hostname     = "localhost"
    hostport     = 5432
	username     = "postgres"
	password     = "root"
	databasename = "ImageStore"
	uploadPath   = "./images"


### Steps to Build the project

Run `go get` command to download all dependencies.
Run `go build` to create the binary.