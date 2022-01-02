# Remote JobTracker

In some cases it is required to submit jobs through the network (like to a sidecar).
Not all JobTracker implementations support that functionality. Hence a general
_JobTracker_ implementation which splits in a server part and a client part can
be useful. The server part should be implementable on top of any existing _JobTracker_
implementation, while the client part should not bother about the actual
server implementation, and can be used in any Go DRMAA2 application.

In order to provide a well-defined interface between the client and server, the
_JobTracker_ interface was re-implemented using the OpenAPI v3 specification.
The current specification file inside this directory. _oapi-codegen_ created client
and server stubs out of the specification. Using them resulted in the remote
client and remote server packages. The OpenAPI specification might be generally
useful for other language bindings or later implementations using protobuf as
more efficient protocol.
