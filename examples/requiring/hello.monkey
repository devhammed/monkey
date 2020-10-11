# Only symbols that starts with capital letter will be exported
# And also there will be a global boolean variable "MAIN" to help determine if the file is been imported
# Try to execute this file directly to see the change

if (!MAIN) {
  puts(FILE + " is imported")
} else {
  puts(FILE + " is executed directly")
}

let hi = "hi world";

let Hello = "Hello World";

let Print = fn(x) {
  puts(x);
};
