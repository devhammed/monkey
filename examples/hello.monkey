# Only symbols that starts with capital letter will be exported
# And also there will be a global boolean variable "MAIN" to help determine if the file is been imported
# Try to execute this file directly to see the change

if (!MAIN) {
  print(FILE + " is imported")
} else {
  print(FILE + " is executed directly")
}

hi = "hi world";

Hello = "Hello World";

Greet = fn(x) {
  print("Hello " + x + "!");
};
