# Only symbols that starts with capital letter will be exported
# And also there will be a global boolean variable "MAIN" to help determine if the file is been imported
# Try to execute this file directly to see the change

if (!MAIN) {
  println(FILE + " is imported")
} else {
  println(FILE + " is executed directly")
}

hi = "hi world";

Hello = "Hello World";

Greet = fn(x) {
  println("Hello " + x + "!");
};
