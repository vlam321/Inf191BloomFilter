# Inf191BloomFilter

## About
This project is a proof of concept. Our problem assumption is that the current system queries userid:email pairs from a giant database. Though this provides an accurate membership result, the network overhead of querying the database is too large. Therefore, we will build a bloom filter to improve data management performance. The current scope will be to integrate and test an existing bloom filter and expand it to a distributed node implementation. We will also build utilities for testing our bloom filter locally. These tools include a data generator and visualization metrics.

## Contribution
To start contributing to this project, simply fork this repository and clone it using after replacing "vlam321" with your user name, Like so: 
`git clone https://github.com/yourGithubUserName/Inf191BloomFilter`

## Dependencies
The following dependencies will be needed to run some of the files in this repo. Use the commands below to install them if needed.
### Go MySQL Driver
` go get -u github.com/go-sql-driver/mysql`
### Bloom Filter
`go get -u github.com/willf/bloom`
### Testify
`go get github.com/stretchr/testify`

## Test Data Generator
Test Data can be generated using the bloomDataGenerator package. The function GenData will require 3 arguments: amount of users, and a minimum and maximum values to set the range for the number of email addresses to be generated.
To populate database with these data, ensure that your database, use the InsertDataSet function and a map[int][]string as instance as argument.
```
data := bloomDataGenerator.GenData(5, 100, 1000)
update := updatebloomData.New(dsn)
update.InsertDataSet(data)
```
