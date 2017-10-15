# Go 'Serverless'

This is a guide to help you build OpenFaaS functions with [Golang server](http://golangserver.com)

### Requirements
1. [Golang server](http://golangserver.com). CLI installed.
2. [OpenFaaS Gateway](https://github.com/openfaas/faas). Running & accessible OpenFaaS gateway
3. [OpenFaaS CLI](https://github.com/openfaas/faas).


# Step 0
Create a new repository on github. I used the Desktop client and later switched to terminal.

# Step 1 - Setup
Change your working directory to your github repository folder. Initialize a Golang server project within the repository with following command :

	gos --make

# Step 2 - Setup
Set FaaS gateway of project. Open the `gos.gxml` within your repository folder and update the second line of the file. Add the attribute `gateway` to specify the location of your OpenFaaS gateway. The example below is set to the default setting of OpenFaaS :

	<gos gateway="http://localhost:8080">

# Step 3 - Setup
Update your `gos.gxml` deploy tag content from `webapp` to `faas`.

# Step 4 - Add
#### Serverless `<end>` tag
Add a new `<end>` tag within `<endpoints>` and voila you have a serverless function with full access to `http.Request` to retrieve body data. `<end>` tags with type attribute set to `f` or left blank will not be processed.

#### Templates
Templates are also processed into OpenFaas functions. More information about adding new templates with [Golang server](http://golangserver.com/docs/markup.html#templates).


# Step 5 - Build
Build and deploy functions with command 

	gos --run

# Step 6 - Access

#### `<end>` tags :
Your `<end>`'s path attribute will be stripped of `/` (forward slashes) and the letter following it will be converted to uppercase. This will become the name of your function on OpenFaaS. The Stdin of this function will be converted into a request body. 


#### `<template>` tags :
Your `<template>`'s name attribute will be the name of your OpenFaaS function. The Stdin of this function is a JSON string which will be converted to the tag's specified `struct` attribute.


#### Notes
With end tags write your response to Stdout directly.
I used [this](https://github.com/cheikhshift/TestFaas/blob/master/gos.gxml) Go server repository for this guide.
Template functions compile their linked `.tmpl` (with internal package `html/template`) file then write the output to Stdout. 
