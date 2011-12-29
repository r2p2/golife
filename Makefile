all:
	8g life.go
	8g main.go
	8l main.8

clean:
	rm *.8
	rm 8.out
