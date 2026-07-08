# Game Editor Design

Before getting into the editors detailed design, let us first understand what all components exist and what are the tasks of each of these components

On a very high level, the game editor has 2 main components:
- Frontend / UI
    - displays the game scene and all the editor components visually to the user
    - allows selection of existing gameObjects in the editor
    - implements drag and drop functionality for the user
    - propagate state update requests to the backend for validation and changes
- Backend
    - accepts and processes command requests from the frontend
    - validates/updates any changes done to the game scene

## Frontend

The panel updates and all the UI are handled in the frontend.
These updates include:
- Selection/Deselection of gameObjects in one panel and updating the others respectively.
- drag and drop of objects in the workspace panel.
- Command Request dispatch according to the event occured
    - These requests wont be sent for every small operation, it will be sent only on notable events that change the game scene state

> Note: For the drag and drop, once the operation is complete there should be a commandRequest propagated to update the new state of the object.

## Backend

The backend should not worry about the styling and placement of the panels nor about the selection/deselection or drag and drop events performed on the frontend. The backend should mainly listen for the following things:
- CommandRequests from the frontend to process validations/ updates
- Fetch requests to fetch state of the scene or particular objects, like:
    - fetch the game scene from the backend on initial boot
    - fetch components of an object on selection

> Note: Any state updates to the game scene should and must propagate through the backend

## Communication between backend and frontend

### Frontend -> Backend
The frontend communicates with the backend via commandRequests. These commandRequests are usually JSON objects propagated to a commandQueue in the backend for processing. 

Also the frontend does not communicate every event to the backend

Let us take examples use cases here to understand better:
- Say we have an Object whose **Transform** has been **updated from pos(x, y) -> (100, 200) to (400, 200)**. Now as this is an event that can update the state of the game scene, this should be sent to the backend for validation and updation. So a commandRequest for this would look something like this
    ```json
    {
        "requestType": "UpdateTransform",
        "objectId": "obj123",
        "x": 400,
        "y": 200
    }
    ```
    This is a sample commandRequest and can be refined further as and how the need arises

- Say the user just selects and deselects the object in the hierarchy panel. Here there was no such event triggered that could update the game scene. Hence, this need not be communicated to the backend.

Now the communication can happen through various ways here:
- http requests
    - The commandRequest is sent by hitting a `/ingest` endpoint
        - The drawback here is this can become slow when there are multiple requests to be sent.
- Web sockets
    - We open an event stream and send request data which makes things faster as it avoids overhead of TCP handshake with every ingest call

The most probable choice would be `Websockets` here.

### Backend -> Frontend

Now that the commandRequest has been submitted by the frontend, we should process it and send a response back saying what was the outcome of the request here.

Let us continue the same example of location update above:
- Now that the backend validated the response and updated the state, the frontend should be notified saying this the update was made. This can be achieved via events. Just send out an event stating the update was complete and these are the new values
    ```json
    {
        "eventStatus": "success",
        "eventType": "UpdateTransform",
        "objectId": "obj123",
        "x": 400,
        "y": 200
    }
    ```

- There can also be requests that are invalid and should not change state, like, say the user sets the movement speed for the component as -ve, this should not be allowed as speed cannot be -ve anytime, velocity can. In such cases the event sent would not be successful and return previous config.
    ```json
    {
        "eventStatus": "failure",
        "eventType": "UpdateComponent",
        "objectId": "obj123",
        "componentName": "movement",
        "speed": 100
    }
    ```
    Notice how the speed was not changed to -ve

These events are propagated back to the frontend via the same `Websocket`.

### Request Response ordering

Here each request is processed sequentially and the response are sent accordingly. But the request and response are decoupled because the response need not know the corresponding request that called it and the frontend just cares about whether the event returned a success or failure and what were the updated states.

> Note: Not all the possible cases are listed here as there can be many. Also the commandRequest and eventResponse can vary depending on the cases. These are just examples to define communication possibilities.