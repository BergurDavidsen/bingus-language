# Bingus Language

Bingus is a tiny compiled language with simple syntax. This project is my playground for learning compilers — parsing code, building ASTs, and generating machine code. Just a fun way to explore how programming languages work under the hood.

## Update (Mon, 25/8-2025)

It can now compile a simple return statement with a return code into assembly and then machine code. It runs on a docker image for a linux x86 distro, generates the given code in a `.bng` file into assembly x86. Then it uses `nasm` to create the object file and finally it uses the linker command `ld` to create the machine code.

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
