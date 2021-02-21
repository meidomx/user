del /Q user
del /Q user_linux.tgz
set GOOS=linux
go build
tar czf user_linux.tgz user static user.toml.example
del /Q user
