# ScheduleFlow

This is a solution to a common problem at my worksite in which several of my coworkers fail to submit a weekly
Schedule before wednesday @ noon. Instead of having to remember to send their excel sheet every week, this project would allow coworkers to set a day and time, and have the app automatically send their schedule to our supervisoers, guaranteeing that everyone has their schedules submitted on time. This will reduce overhead for co-workers, and allow supervisors to get bulk schedules

### Tech Stack:
* [HTMX](https://htmx.org/)
* [Sass](https://sass-lang.com/)
* [Go](https://go.dev/)
* [go-chi](https://github.com/go-chi/chi)
* [MongoDB](https://www.mongodb.com/) 
* [Railway](https://railway.com)
* [Docker](https://www.docker.com/)

## Features

* Users can set Month, Day and Time for the email to be sent out, in which one email will be sent out every week at that specific date.
* Users can upload their schedule, and have that be sent out based on the time set, or sent manually with a "Send Now" button.
* Oauth using microsoft 365 to connect the app to our outlook accounts so the emails sent out on behalf of the coworkers work email.

### Requirements:

* Clone repo using `git clone https://github.com/darienmiller88/ScheduleFlow.git`
* Migrate the necessary information to your local `.env` as described in the `.env_sample` file
* Install the Go migrate CLI tool:
  * **Windows (using Scoop):** `scoop install migrate`
  * **macOS (using Homebrew):** `brew install golang-migrate`
  * **Linux:** Download from [releases](https://github.com/golang-migrate/migrate/releases) or use `curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add - && apt-get update && apt-get install -y migrate`
* Run database migrations: `migrate -path ./migrations -database "postgres://username:password@localhost:5432/dbname?sslmode=disable" up`
  * Replace connection string values with those from your `.env` file
* Run go build to create a root level `JeopardyScoreBoardV2.exe` file, and then run `.\JeopardyScoreBoard-V2` to run the executable. If an executable is not needed, simply input `go run main.go` instead, or `.\fresh` to enable a server restart on change.

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Feel free to leave suggestions as well, I'm always looking for ways to improve!

<p align="right">(<a href="#top">back to top</a>)</p>

## License
[MIT](https://choosealicense.com/licenses/mit/)