# supernova
supernova is a command-line interface (CLI) designed for executing various tests in a pipeline-based manner.

## usage 
In supernova, functionalities are provided in the form of "templates".   
Each template encapsulates a specific set of features.   
For instance, there are templates for making HTTP requests and validating HTML.

```
steps:
  - name: Register data via API
    template: curl
    option:
      url: http://example.com/user
      method: POST
      body: {"name": "test user"}
      expect:
        equal: {"result": "OK"}
        status: 200
  - name: Check if the registered user is visible in the admin dashboard
    template: html
    option:
      url: http://example.com/dashboard/user/xxx
      screenshot:
        waitSec: 1
```

To explore the available templates and their functionalities, please refer to the Wiki.
