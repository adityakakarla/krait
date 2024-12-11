# Krait

Krait helps Go developers build Go CLI tools in seconds using LLMs.

To install, run the command ```go install github.com/adityakakarla/krait@latest```. Go will automatically install it in your $GOPATH/bin directory which should be in your $PATH.

Run ```krait``` to verify installation.

## Configuration

To configure Krait, grab your OpenAI API key and run the following command: ```krait set -k [OpenAI API key goes here]```.

You can verify that Krait has successfully stored your key by running ```krait view```.

In case you need to change your key, you can do so by simply running ```krait set -k [new key]```.

## Generation

To generate your CLI tool, first navigate to a parent folder. Then run ```krait generate```.

You will be asked to enter in a tool name and tool description.

Once generated, cd into the tool folder (```cd [tool name]```). Then, run ```go install```.

Once installed, you should be able to invoke your CLI tool anywhere. To learn how it works, check the code or simply run ```[tool name]```.

## About

I built this project as a way to learn Go and procrastinate on studying for my finals. Enjoy.

Feel free to contact me at adi[at]adikakarla[dot]com for any feedback. This is my first real Go project, and I would love to know if you run into any errors/bugs or have constructive criticism.

## Future Steps

I want to turn this into a v0-esque project. Being able to generate multiple files would be awesome and is something I'm currently working on.

Inspiration taken from [cobra-cli](https://github.com/spf13/cobra-cli).