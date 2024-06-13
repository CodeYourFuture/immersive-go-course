run: out
	./out

clean:
	rm -f out
	find . -name '*.o' -exec rm {} \;

out: main.o formatting/formatting.o
	clang++ main.o formatting/formatting.o -o out

formatting/formatting.o: formatting/formatting.cpp
	clang++ formatting/formatting.cpp -I "." -c -o formatting/formatting.o

main.o: main.cpp gen/constants.h formatting/formatting.h
	clang++ main.cpp -std=c++11 -I "." -c -o main.o

gen/constants.h: scripts/generate_constants.sh
	./scripts/generate_constants.sh "Daniel" "Laura"
