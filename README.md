# PROJECT KALLAXY
VincNT21's Capstone project   
End of backend-development learning path on [boot.dev](https://www.boot.dev/)

## PRESENTATION
Project Kallaxy is a virtual shelf that allows users to :
- Store information about various cultural media (books, movies, series, boardgames, video games...), as well as your personal reading/watching/playing record
- Use 3rd party API to get precise details online about a specific medium
- Review your collections as a shelf, organized by media type.


And more to come in future update !
- Searching, sorting, filtering your collection based on criteria
- Get "by-period" stats (weekly, monthly, yearly, custom period) and review
- Share medias to other users


## PROJECT DETAILS
Project Kallaxy, in its current version,  is made up of four parts : 
- a PostgreSQL database
- an HTTP server RESTful API
- an HTTP client REPL API 
- a front-end simple GUI, using Fyne
  
Server and client are coded in Go.  
Everything is written by myself (around 125 hours of work), *certified with no vibe coding ;)*

The idea of this project is based on my previous personal project : [Book Shelf and Stats](https://github.com/VincNT21/books_shelf_and_stats)

## LEARNED SKILLS

For this project I used and improved some of the skills I've learned on boot.dev.  
Here's a short summary :
- Using Go (Golang) programming language and its packages
- Designing a SQL database and writing queries to interact with it  
(including the use of tools to make it easier: Goose and SQLC generator)
- Building an HTTP web server to handle requests, interact with the database and make JSON responses  
(including integration testing, middleware, routing, logging, webhooks, authentification, authorization, JWTs, proxy serving)
- Writing clear, concise server's documentation
- Building an HTTP Client to make request, handle response and manage user's input  
(including config and appState, caching, making requests with headers, query parameters)
- Learning quickly how to use external libraries/toolkit with online documentation
- Using Docker to create, run and publish images, containers and compose
- Using git and github for publishing commits and proper versioning

## INSTALLATION
### For Windows users
1. Download and install [Docker for desktop](https://desktop.docker.com/win/main/amd64/Docker%20Desktop%20Installer.exe?utm_source=docker&utm_medium=webreferral&utm_campaign=dd-smartbutton&utm_location=module)  
1b. Once installed, you do not need to login in Docker.   
Just got to settings and uncheck "Open Docker Dashboard when Docker Desktop starts"
2. Download Project Kallaxy [last release](https://github.com/VincNT21/project_kallaxy/releases)
3. Unzip the file 
4. Double-click kallaxy_start.bat to run the start script  
4b. If you don't want to be bothered by the command window opening :  
Create a shortcut for kallaxy_start.bat, right-click on it, "Properties",   
in "Run" list, choose "Minimized"
5. Use the app !
6. Close it when you're done and server/database will be properly stopped  
6b. For now, Docker Desktop won't stop automatically.  
You can manually stop it by right-clicking on its icon in taskbar and chose "Quit Docker Desktop"

> A word about updates :
> - If a new version of the Server is online, it will be automatically get by Docker
> - If a new version of the Client is online, you'll have a pop-up when launching the app with a link to download it
>   
> In both cases, **your local data (users, media and personal records) will be stored!**

### For Linux and MacOs users
*Coming soon!*

## DOCUMENTATION

### Server API

All documentation about server endpoints and resources can be found in the */server/documentation* folder

### Server Database

Database schemas and queries can be found in the */database* folder

### How to use Client GUI
1. First, create a new user (a valid email is not required in local version)
2. Login using your username and password and you are on the Home page !
3. When creating new media, the easiest way is to provide a title and click the "Get Info Online" button
Once medium find online, you can edit whatever field you want and add your personal records at the bottom of the page
4. When in your Shelf, you can access media details of a certain type by clicking on the corresponding compartment's title
5. From there, you can edit or delete medium's info and/or your personal record about it. 


*Screenshots coming soon!*

## App icon attribution
[bookshelf icons](https://www.flaticon.com/free-icons/bookshelf)Bookshelf icons created by Freepik - Flaticon