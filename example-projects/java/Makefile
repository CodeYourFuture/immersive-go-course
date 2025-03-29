run: com/example/Main.class
	java -cp .:third_party/guava-33.2.1-jre.jar com.example.Main

clean:
	find . -name '*.class' -exec rm {} \;
	rm -rf third_party com/example/gen/Constants.java:

com/example/Main.class: com/example/gen/Constants.class com/example/fmt/Formatting.class
	javac com/example/Main.java

com/example/gen/Constants.java: scripts/generate_constants.sh
	./scripts/generate_constants.sh "Daniel" "Laura"

com/example/gen/Constants.class: com/example/gen/Constants.java
	javac com/example/gen/Constants.java

com/example/fmt/Formatting.class: third_party/guava-33.2.1-jre.jar
	javac -cp third_party/guava-33.2.1-jre.jar com/example/fmt/Formatting.java

third_party/guava-33.2.1-jre.jar:
	mkdir -p third_party
	curl -o third_party/guava-33.2.1-jre.jar https://repo1.maven.org/maven2/com/google/guava/guava/33.2.1-jre/guava-33.2.1-jre.jar
