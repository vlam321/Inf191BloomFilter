# Inf191BloomFilter

## Installation
To start working on this project, simply fork this repository and clone it using after replacing ¿?vlam321¿? with your user name, Like so: 
`git clone https://github.com/yourGithubUserName/Inf191BloomFilter`

## Dependencies
The following dependencies will be needed to run some of the files in this repo. Use the commands below to install them if needed.
### Go MySQL Driver
` go get -u github.com/go-sql-driver/mysql`
### 
`go get -u github.com/willf/bloom`

## Test Data Generator
Test Data can be generated using the bloomDataGenerator package. The function GenData will require 3 arguments: amount of users, and a minimum and maximum values to set the range for the number of email addresses to be generated.
To populate database with these data, ensure that your database, use the InsertDataSet function and a map[int][]string as instance as argument.
```
data := bloomDataGenerator.GenData(5, 100, 1000)
update := updatebloomData.New(dsn)
update.InsertDataSet(data)
```

