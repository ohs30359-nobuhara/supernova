# Browser CLI

This is a command-line interface (CLI) tool built with TypeScript to interact with a headless browser.   
It allows you to navigate to a specified URL, take a screenshot, and retrieve core web vital metrics.

## Installation
Ensure that you have Node.js and Yarn installed on your machine.  
```
cd browser
yarn install
yarn build
```

## Usage
The CLI provides the following options:

- --performance <performance>: Get core web vital metrics. Please specify true or false. Default is false.
- --screenshot <screenshot>: Take a screenshot. Please specify the file name.


## Example Usage
```shell
node ./dist/main.js --performance true --screenshot ./sample.png https://www.google.com/ 
```

## Dependencies
Puppeteer: Headless Chrome Node API.  
CAC: Simple yet powerful framework for building command-line applications.
