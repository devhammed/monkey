t = "Hello, " + input("What is your name? ") + "\n"

f = open("/tmp/monkey.txt", "w+")

write(f, t)

seek(f, 0)

print(read(f, len(t)))

close(f)