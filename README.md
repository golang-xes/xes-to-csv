# XES to CSV Converter

## Introduction
This is a simple command-line tool written in Go (Golang) that converts XES files into CSV files. The tool reads an input XES file, parses the event log data, and writes it out to a CSV file.

## Requirements
- Go 1.18 or higher

## Installation
To install the converter, you need to have Go installed on your system. Then, you can clone this repository and build the project:

```bash
$ git clone https://github.com/golang-xes/xes-to-csv
$ cd xes-to-csv-converter
```

## Usage
- The basic usage of the converter is as follows:
```bash
$ ./xes2csv -input <input-file>.xes -output <output-file>.csv
```
- input: Path to the input XES file.
- output: Path where the output CSV file should be saved.
Basic Example
Convert an XES file to CSV:
```bash
$ ./xes2csv -input example.xes -output example.csv
```

## License
This project is licensed under the MIT License - see the LICENSE file for details.