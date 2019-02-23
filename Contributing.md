# Contributing

## Introduction

Thank you so much for considering to contribute.
I don't have time to contribute myself, so I appreciate the help to maintain and build this project out further.
I will be here to discuss your issue, help with designing a solution and review your code.
I will do my best to guide you, whatever your skill level.

*In all cases please open an issue* to discuss your idea or issue, before doing any work.
I don't have time to do the work myself, so don't worry about losing the chance to implement your idea.

If it is your first time contributing and don't have an idea of your own, then you can start with the list of [good first issues](https://github.com/awalterschulze/goderive/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22)

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
