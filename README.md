Gator is a terminal-based RSS feed aggregator.

Requirements: Go > 1.24 and Postrges

Installation: Download, navigate to folder, and run go install to install gator.

Config: Create a config file ~/.gatorconfig.json to record the db location and current user.
    {"db_url":"postgres://<user>:<host>:<port>/gator?sslmode=disable","current_user_name":""}
