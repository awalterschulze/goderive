# Contributing

tl;dr make a pull request

## Contributing Examples

We can always use better examples.
Contributing examples for a plugin can be done by creating a folder if one does not already exists in the `example/plugin` folder.

First create your example go file `foldername.go` with your example code.

Please do not use words like foo, bar and MyStruct.
Examples are more useful if they are as close as possible to actual use cases.
Yes some current examples violate these rules, contributions that replace these examples will be highly appreciated.

Next create a Readme.md in the folder with the following content:

````
The <functionName> function ...

```go

```

goderive will generate the following code:

```go

```

````

The go blocks will be populated with your example code and the generated code using the following command:

```sh
cd .. # goderive/example
make
```

Finally go and add a link (if one does not already exist) in the top level Readme.md of the project.

You are ready to submit a pull request.
