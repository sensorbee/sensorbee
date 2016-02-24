FORMAT: 1A

# SensorBee API Version 1

This is a document for SensorBee API version 1.

# Group Topologies

This resource allows clients to manage topologies to create sources and sinks
through BQL.

## Topology Collection [/api/v1/topologies]

### List All Topologies [GET]

This action returns a list of all topologies in the server. It doesn't support
pagination yet.

+ Response 200 (application/json)
    + Attributes (object)
        + topologies (array[Topology]) - A list of topologies

+ Response 500 (application/json)

    500 is returned when the server failed to process the request properly and
    the request did not have any problem.

    + Attributes (Error Response)

### Create a New Topology [POST]

This action creates a new topology on the server.

+ Request (application/json)

    + Body

            {
                "name": "some_topology"
            }

    + Attributes (object)
        + name: `some_topology` (string) - The name of the topology to be created

+ Response 200 (application/json)

    200 OK is returned on success with the information of the newly created
    topology.

    + Attributes (object)
        + topology (Topology)

+ Response 400 (application/json)

    400 is returned when the following cases happened: (1) a topology having
    the same name already exists on the server, (2) request body has a bad
    value.

    + Attributes (Error Response)

+ Response 500 (application/json)

    500 is returned when the server failed to process the request properly and
    the request did not have any problem.

    + Attributes (Error Response)

## Topology [/api/v1/topologies/{topology_name}]

### View a Topology Detail [GET]

This action returns detailed information of a topology having `topology_name`.

+ Response 200 (application/json)
    + Attributes (object)
        + topology (Topology) - Information of a topology

+ Response 404 (application/json)

    404 is returned when the topology having `topology_name` does not exist
    on the server.

    + Attributes (Error Response)

+ Response 500 (application/json)

    500 is returned when the server failed to process the request properly and
    the request did not have any problem.

    + Attributes (Error Response)

### Destroy a Topology [DELETE]

This action destroys a topology having `topology_name`. It also stops the
topology before destroying it. This action may take time to stop all nodes in
the topology. This action does not return 404 when the topology does not exist.

+ Response 200 (application/json)

    An empty object is currently returned on success.

    + Attributes (object)

+ Response 500 (application/json)

    500 is returned when the server failed to process the request properly and
    the request did not have any problem.

    + Attributes (Error Response)

## Queries [/api/v1/topologies/{topology_name}/queries]

### Send Queries [POST]

This action accepts BQL queries. As a result, new nodes may be created or some
existing nodes are changed or dropped. A client can send multiple queries at
once. However, SELECT statements cannot be mixed with other statements including
SELECT statements themselves. In other words, only one SELECT statement can be
issued in a request and the request must only have one statement.

A response of a SELECT statement differs from other statements' responses. It's
returned as a `multipart/mixed` response having multiple `application/json`
contents. Other statements return `application/json` content as described below.

+ Request (application/json)
    + Attributes (object)
        + queries: `CREATE SOURCE s TYPE my_source WITH param="value";` (string) - Multiple BQL statements to be executed

+ Response 200 (application/json)

    This is the regular response of statements other than SELECT statements.

    + Attributes (object)
        + responses (array[Topology Query Response]) - An array having a response of each statement

+ Response 200 (multipart/mixed)

    This is the response of a SELECT statement containing multiple
    `application/json` split by boundaries. Each part contains a tuple emitted
    from the SELECT statement.

    + Body

            --boundary
            Content-Type: application/json

            {"id":1,"price":100,"name":"book1"}
            --boundary
            Content-Type: application/json

            {"id":2,"price":150,"name":"book3"}
            --boundary--

+ Response 400 (application/json)

    400 is returned when one of the given statements has a syntax error or
    fails to be executed. It's also returned when a SELECT statement is issued
    with other statements.

    + Attributes (Error Response)

+ Response 500 (application/json)

    500 is returned when the server failed to process the request properly and
    the request did not have any problem.

    + Attributes (Error Response)

# Data Structures

## Topology (object)

+ name: `some_topology` (string) - The name of the topology

## Node (object)

+ name: `node_name` (string) - The name of the node
+ type: `source` (string) - The type name of the node
+ status (object) - Status information of the node
+ path: `/api/v1/topologies/topology_name/source/node_name` (string) - The path at which the node is located

## Topology Query Response (object)

+ statement: `CREATE SOURCE s TYPE my_source WITH param="value";` (string) - A BQL statement which has been executed
+ nodes (object) - Nodes in the topology which the statement affected
    + created (array[Node]) - Nodes created by the statement
    + dropped (array[Node]) - Nodes dropped by the statement
    + updated (array[Node]) - Nodes updated by the statement

## Error (object)

+ code: `E0123` (string) - Error code
+ message: `something went wrong` (string) - A error message describing what happened
+ request_id: `123` (number) - ID of a request which caused the error
+ meta (object) - Meta information of the error

## Error Response (object)

+ error (Error) - An error information
