fileName = 'hello.txt';

file_write(fileName, 'Hello World', 0644);

puts(file_read(fileName));
