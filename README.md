# Bingus Language

Bingus is a tiny compiled language with simple syntax. This project is my playground for learning compilers — parsing code, building ASTs, and generating machine code. Just a fun way to explore how programming languages work under the hood.

## Updates

### Update (Mon, 25/8-2025)

It can now compile a simple return statement with a return code into assembly and then machine code. It runs on a docker image for a linux x86 distro, generates the given code in a `.bng` file into assembly x86. Then it uses `nasm` to create the object file and finally it uses the linker command `ld` to create the machine code.

### Update (Tue, 26/8-2025)

It can now store and retrieve variables, aswell as do arithmetic operations with numbers and variables🥳. You define a variable using the syntax `let <variable> = <value>`. It currently only supports integers and number literals. The variable can be any string.

To do arithmetic, you chose to do addition (`+`), subtraction (`-`), multiplication (`*`), and integer division (`/`).
Examples for valid programs would be:

```bash
# Defining variables
let x = 8;
let y = 2;

# Arithmetic operations using variables
let add = x+y;
let sub = x-y;
let mul = x*y;
let div = x/y;

# Returning variables
return add+sub+mul+div;
```

The return statement will evaluate the result to be `36` and will set that as the exit code. You can verify this by compiling the code, running the compiled code and then by checking the exit code of the compiled program with `echo $?`.

This syntax is temporary and will most likely be changed to something else, but I have still not really decided exactly what I want it to be.

### Update (Saturday, 30/8-2025)

The language now has print features. It can now print numbers literals and variables to the terminal using the syntax `print <number or variable>`. This makes it way more usable as you now do not need to check the return code to see the result of the written code, because you can now just print directly to the terminal.

Some example code, building on the last example, that compiles and works is:

### Update (Sunday, 31/8-2025)

The language now has boolean literals and can do boolean comparisons that result in boolean literals. An example that compiles is:

```bash
let x = 3 <= 6;
print x;
return 0;
```

It supports `==`, `<=`, `>=`, `<` and `>` as comparators and also `true` and `false` as boolean constants. If the literal is `true` it will result in the integer value `1` and `0` otherwise.

```bash
let x = 8;
let y = 2;
let add = x+y;
let sub = x-y;
let mul = x*y;
let div = x/y;
let all = add+sub+mul+div;
print all;
```

It also now checks if there is a return statement at the end of the file, and if there is not, then it prints a warning at compile time and adds a default `return 0` to the generated `assembly` code. This is was to make the syntax a bit simpler and also to avoid segmentation fault when a user does not add a return statement to the code file.

Now it also has slightly better stack space management by allocating the appropriate space when storing variables on the stack, and as well as now not pushing and popping of the stack for every binary operation. Instead it simply writes the right side and left side to each of their respective registers.

## Check it out

To compile a `.bng` file, can follow these steps:

### 1. Create a file bingus file

Create a new file bingus ending with the `.bng` extension or use the `test.bng`

### 2. Build the project with make

Now, build the project using the configured Makefile by running:

```bash
make build LINUX=1
```

This builds an exacutable for the linux amd64 architecture, because the project will use assembly x86 and will be running on a docker image with that specific architecture.

### 3. Run the docker container

Start the docker container by running the following command:

```bash
docker compose run --rm bingus-dev    
```

This will give you a linux shell with the necessary tools like `nasm` and `ld`, aswell as make sure you use the correct architecture.

### 4. Compile the file

To compile the file, simply run this command in the shell:

```bash
./bin/bingus <your-filename>.bng
```

This will compile your bingus file if it is correct. It will give an error otherwise.

### 5. See the result

To see and run the compiled `.bng` file, you can run the created executable in the output folder like so:

```bash
./output/test
```

This will show you the resulting code execution.
As of Monday, 25/8-2025, it only compiles return statements, so to see the result of the execution, check the exit code like so:

```bash
echo $?
```

This will show you what the exit code of the last call. For example, if the code is ```return 5;```, then it will echo `5` in the shell.
