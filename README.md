# Clean Architecture Template for Golang webservices
A webservice implementation using Golang and strictly following the clean architecture.

The clean architecture visualised:

![clean architecture](https://canobbioe.com/posts/cleaner-architecture/CleanArchitecture.jpeg)

The folder structure matches the rings in the above diagram. Starting from the innermost ring:

- The `domain` folder contains the Enterprise Business Rules - i.e. entities
- The `usecase` folder contains the Application Business Rules - i.e. data manipulation
- The `interface` folder contains the Interface Adapters - i.e. input/output transformation
- The `infrastructure` folder contains the Frameworks and Drivers - i.e. protocol-specific implementations

# TODO
- Propagate context
- Add list filters
- Add Docker
- Add CI/CD
- Add integration tests
